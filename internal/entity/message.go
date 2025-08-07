package entity

import (
	"time"
)

type Message struct {
	ID        uint64
	ThreadID  uint64
	AccountID uint64
	CreatedAt time.Time
	Mentions  Mentions
	Content   string
}
