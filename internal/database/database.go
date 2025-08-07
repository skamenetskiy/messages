package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/goccy/go-yaml"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

func New(c Config) (DB, error) {
	if len(c.Shards) == 0 {
		return nil, errors.New("no shards specified")
	}
	var (
		shards   = make([]*shard, len(c.Shards))
		writable = make([]*shard, 0, len(c.Shards))
		byId     = make(map[uint16]*shard, len(c.Shards))
		dedupe   = make(map[uint16]struct{}, len(c.Shards))
	)
	for i, s := range c.Shards {
		if _, ok := dedupe[s.ID]; ok {
			return nil, fmt.Errorf("duplicate shard id: %d", s.ID)
		}
		dedupe[s.ID] = struct{}{}
		conn, err := sql.Open("postgres", s.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to shard: %w", err)
		}
		if err = conn.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping shard: %w", err)
		}
		sh := &shard{
			id:       s.ID,
			writable: s.Writable,
			conn:     conn,
			gen:      newGenerator(s.ID),
		}
		shards[i] = sh
		byId[s.ID] = sh
		if s.Writable {
			writable = append(writable, sh)
		}
	}
	if len(writable) == 0 {
		return nil, errors.New("there should be at least one writable shard")
	}
	return (&db{
		all:      shards,
		writable: writable,
		byID:     byId,
		next:     0,
	}).migrate()
}

func NewFromFile(path string) (DB, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read database config file: %w", err)
	}
	var c Config
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
	}
	return New(c)
}

func NewDefault() (DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}
	return NewFromFile(filepath.Join(wd, "db.yml"))
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type DB interface {
	// Close all shards' connections.
	Close() error

	// Next writable shard.
	Next() Shard

	// Shard by entity id.
	Shard(id uint64) Shard

	// Shards by entity ids.
	Shards(ids ...uint64) map[Shard][]uint64
}

type db struct {
	all      []*shard
	writable []*shard
	byID     map[uint16]*shard
	next     uint64
}

func (d *db) Close() error {
	eg := errgroup.Group{}
	for _, s := range d.all {
		eg.Go(s.conn.Close)
	}
	return eg.Wait()
}

func (d *db) Next() Shard {
	n := atomic.AddUint64(&d.next, 1)
	return d.all[(int(n)-1)%len(d.all)]
}

func (d *db) Shard(id uint64) Shard {
	shardID := ShardID(id)
	if s, ok := d.byID[shardID]; ok {
		return s
	}
	return nil
}

func (d *db) Shards(ids ...uint64) map[Shard][]uint64 {
	shards := make(map[Shard][]uint64, len(ids))
	for _, id := range ids {
		s := d.Shard(id)
		if s == nil {
			continue
		}
		shards[s] = append(shards[s], id)
	}
	return shards
}

func (d *db) migrate() (DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	for _, s := range d.all {
		eg.Go(migrate(ctx, s.id, s.conn))
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return d, nil
}

func migrate(ctx context.Context, id uint16, conn *sql.DB) func() error {
	return func() error {
		_, err := conn.ExecContext(ctx, schema)
		if err != nil {
			return fmt.Errorf("failed to migrate shard %d: %w", id, err)
		}
		return nil
	}
}
