package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/skamenetskiy/messages/internal/database"
	"golang.org/x/sync/errgroup"
)

func (r *Repo) Delete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	shards := r.db.Shards(ids...)
	eg, ctx := errgroup.WithContext(ctx)
	for shard, list := range shards {
		eg.Go(func() error {
			return r.DeleteOnShard(ctx, shard, list)
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}
	return nil
}

const deleteQuery = "delete from messages"

func (r *Repo) DeleteOnShard(ctx context.Context, shard database.Shard, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	in := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		in[i] = "$" + strconv.Itoa(i+1)
		args[i] = strconv.FormatUint(id, 10)
	}
	query := deleteQuery + fmt.Sprintf(" where id in (%s)", strings.Join(in, ","))
	if _, err := shard.Conn().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete messages on shard %d: %w", shard.ShardID(), err)
	}
	return nil
}
