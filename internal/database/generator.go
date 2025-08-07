package database

import (
	"sync/atomic"
	"time"
)

// ShardID from id.
func ShardID(id uint64) uint16 {
	shardID := id >> 10
	shardID = shardID & (uint64(1<<10) - 1)
	return uint16(shardID)
}

const epoch uint64 = 1474661123

func newGenerator(shardID uint16) *generator {
	return &generator{
		counter: 0,
		shardID: shardID,
	}
}

type generator struct {
	counter uint32
	shardID uint16
}

// Generate new unique id.
func (g *generator) generate() uint64 {
	var id uint64
	id = (uint64(time.Now().UnixMilli()*1000) - epoch) << 23
	id = id | uint64(g.shardID<<10)
	id = id | (g.next() % 1024)
	return id
}

func (g *generator) next() uint64 {
	return uint64(atomic.AddUint32(&g.counter, 1))
}
