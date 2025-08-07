package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/skamenetskiy/messages/internal/entity"
)

const insertQuery = `insert into messages (id, thread_id, account_id, created_at, mentions, content)  values ($1, $2, $3, $4, $5, $6)`

func (r *Repo) Insert(ctx context.Context, msg entity.Message) (entity.Message, error) {
	shard := r.db.Next()

	msg.ID = shard.NextID()
	msg.CreatedAt = time.Now()

	var (
		threadID  *string
		accountID *string
	)
	if msg.ThreadID > 0 {
		s := strconv.FormatUint(msg.ThreadID, 10)
		threadID = &s
	}
	if msg.AccountID > 0 {
		s := strconv.FormatUint(msg.AccountID, 10)
		accountID = &s
	}

	if _, err := shard.Conn().ExecContext(
		ctx,
		insertQuery,
		strconv.FormatUint(msg.ID, 10),
		threadID,
		accountID,
		msg.CreatedAt,
		msg.Mentions,
		msg.Content,
	); err != nil {
		return entity.Message{}, fmt.Errorf("failed to insert: %w", err)
	}

	return msg, nil
}
