package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/skamenetskiy/messages/internal/database"
	"github.com/skamenetskiy/messages/internal/entity"
	"golang.org/x/sync/errgroup"
)

func (r *Repo) Get(ctx context.Context, ids []uint64) ([]entity.Message, error) {
	shards := r.db.Shards(ids...)
	mu := sync.Mutex{}
	res := make([]entity.Message, 0, len(ids))
	eg, ctx := errgroup.WithContext(ctx)
	for shard, list := range shards {
		eg.Go(func() error {
			items, err := r.GetFromShards(ctx, shard, list)
			if err != nil {
				if database.IsNotFound(err) {
					return nil
				}
				return err
			}
			if items == nil {
				return nil
			}
			mu.Lock()
			res = append(res, items...)
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return res, nil
}

const getQuery = "select id, thread_id, account_id, created_at, mentions, content from messages"

func (r *Repo) GetFromShards(ctx context.Context, shard database.Shard, ids []uint64) ([]entity.Message, error) {
	in := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		in[i] = "$" + strconv.Itoa(i+1)
		args[i] = strconv.FormatUint(id, 10)
	}
	query := getQuery + fmt.Sprintf(" where id in (%s)", strings.Join(in, ","))
	rows, err := shard.Conn().QueryContext(ctx, query, args...)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get messages from shard %d: %w", shard.ShardID(), err)
	}
	defer func() { _ = rows.Close() }()
	res := make([]entity.Message, 0, len(ids))
	for rows.Next() {
		var (
			m         entity.Message
			id        string
			threadID  *string
			accountID *string
		)
		if err = rows.Scan(&id, &threadID, &accountID, &m.CreatedAt, &m.Mentions, &m.Content); err != nil {
			return nil, fmt.Errorf("failed to scan messages on shard %d: %w", shard.ShardID(), err)
		}
		m.ID, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse message id on shard %d: %w", shard.ShardID(), err)
		}
		if threadID != nil {
			m.ThreadID, err = strconv.ParseUint(*threadID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse thread_id on shard %d: %w", shard.ShardID(), err)
			}
		}
		if accountID != nil {
			m.AccountID, err = strconv.ParseUint(*accountID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse account id on shard %d: %w", shard.ShardID(), err)
			}
		}
		res = append(res, m)
	}
	return res, nil
}
