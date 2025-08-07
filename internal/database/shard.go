package database

import (
	"database/sql"
)

type Shard interface {
	ShardID() uint16
	Writable() bool
	Conn() *sql.DB
	NextID() uint64
}

type shard struct {
	id       uint16
	writable bool
	conn     *sql.DB
	gen      *generator
}

func (s *shard) ShardID() uint16 {
	return s.id
}

func (s *shard) Writable() bool {
	return s.writable
}

func (s *shard) Conn() *sql.DB {
	return s.conn
}

func (s *shard) NextID() uint64 {
	return s.gen.generate()
}
