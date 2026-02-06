package cache

import (
	"log"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*node

	head *node // most recently used
	tail *node // least recently used

	mu sync.Mutex
}

type node struct {
	key   string
	value any
	prev  *node
	next  *node
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		log.Printf("INVALID: capacity must be > 0")
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*node),
	}
}

func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.moveToHead(n)
	return n.value, true
}

func (c *LRUCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.cache[key]; ok {
		n.value = value
		c.moveToHead(n)
		return
	}

	n := &node{
		key:   key,
		value: value,
	}

	c.cache[key] = n
	c.addToHead(n)

	if len(c.cache) > c.capacity {
		c.evict()
	}
}

func (c *LRUCache) addToHead(n *node) {
	n.prev = nil
	n.next = c.head

	if c.head != nil {
		c.head.prev = n
	}
	c.head = n

	if c.tail == nil {
		c.tail = n
	}
}

func (c *LRUCache) removeNode(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

func (c *LRUCache) moveToHead(n *node) {
	c.removeNode(n)
	c.addToHead(n)
}

func (c *LRUCache) evict() {
	if c.tail == nil {
		return
	}

	lru := c.tail
	c.removeNode(lru)
	delete(c.cache, lru.key)
}
