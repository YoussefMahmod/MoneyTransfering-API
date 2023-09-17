package customtypes

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type ShardMap struct {
	shards     []*ConcurrentMapShard
	shardCount int
}

type ConcurrentMapShard struct {
	items map[interface{}]interface{}
	sync.RWMutex
}

func NewShardMap(shardCount int) *ShardMap {
	m := &ShardMap{
		shards:     make([]*ConcurrentMapShard, shardCount),
		shardCount: shardCount,
	}

	for i := 0; i < shardCount; i++ {
		m.shards[i] = &ConcurrentMapShard{
			items: make(map[interface{}]interface{}),
		}
	}
	return m
}

func (m *ShardMap) GetShard(key interface{}) *ConcurrentMapShard {
	var hash uint32

	switch key := key.(type) {
	case int:
		hash = uint32(key)
	case string:
		hash = fnv32(key)
	default:
		hash = fnv32(fmt.Sprintf("%s", key))
	}

	return m.shards[hash%uint32(m.shardCount)]
}

func fnv32(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}

func (m *ShardMap) Set(key interface{}, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (m *ShardMap) Get(key interface{}) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m *ShardMap) Count() int {
	count := 0
	for i := 0; i < m.shardCount; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m *ShardMap) Del(key interface{}) (interface{}, bool) {
	shard := m.GetShard(key)
	item, ok := m.Get(key)
	shard.Lock()

	if !ok {
		shard.Unlock()
		return item, false
	}

	delete(shard.items, key)
	shard.Unlock()

	return item, true
}

func (m *ShardMap) GetAll() []interface{} {
	values := make([]interface{}, 0)
	for i := 0; i < m.shardCount; i++ {
		shard := m.shards[i]
		shard.RLock()
		for _, value := range shard.items {
			values = append(values, value)
		}
		shard.RUnlock()
	}
	return values
}
