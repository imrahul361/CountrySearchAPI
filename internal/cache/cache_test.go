package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLRUCache(t *testing.T) {
	t.Run("creates cache with valid capacity", func(t *testing.T) {
		cache := NewLRUCache(5)
		require.NotNil(t, cache)
		assert.Equal(t, 5, cache.capacity)
		assert.NotNil(t, cache.cache)
		assert.Equal(t, 0, len(cache.cache))
	})

	t.Run("creates cache with capacity 1", func(t *testing.T) {
		cache := NewLRUCache(1)
		require.NotNil(t, cache)
		assert.Equal(t, 1, cache.capacity)
	})

	t.Run("creates cache with negative capacity (logs warning)", func(t *testing.T) {
		cache := NewLRUCache(-1)
		require.NotNil(t, cache)
		assert.Equal(t, -1, cache.capacity)
	})

	t.Run("creates cache with zero capacity (logs warning)", func(t *testing.T) {
		cache := NewLRUCache(0)
		require.NotNil(t, cache)
		assert.Equal(t, 0, cache.capacity)
	})
}

func TestLRUCache_Set_And_Get(t *testing.T) {
	t.Run("set and get single item", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")

		value, ok := cache.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "value1", value)
	})

	t.Run("get non-existent key returns false", func(t *testing.T) {
		cache := NewLRUCache(5)
		value, ok := cache.Get("nonexistent")
		assert.False(t, ok)
		assert.Nil(t, value)
	})

	t.Run("set multiple items", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		v1, ok1 := cache.Get("key1")
		v2, ok2 := cache.Get("key2")
		v3, ok3 := cache.Get("key3")

		assert.True(t, ok1)
		assert.True(t, ok2)
		assert.True(t, ok3)
		assert.Equal(t, "value1", v1)
		assert.Equal(t, "value2", v2)
		assert.Equal(t, "value3", v3)
	})

	t.Run("set with different value types", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("int", 42)
		cache.Set("string", "hello")
		cache.Set("map", map[string]int{"a": 1})

		intVal, _ := cache.Get("int")
		stringVal, _ := cache.Get("string")
		mapVal, _ := cache.Get("map")

		assert.Equal(t, 42, intVal)
		assert.Equal(t, "hello", stringVal)
		assert.Equal(t, map[string]int{"a": 1}, mapVal)
	})
}

func TestLRUCache_Update(t *testing.T) {
	t.Run("update existing key updates value", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")
		cache.Set("key1", "updated")

		value, ok := cache.Get("key1")
		assert.True(t, ok)
		assert.Equal(t, "updated", value)
	})

	t.Run("update does not increase cache size", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")
		initialSize := len(cache.cache)

		cache.Set("key1", "updated")
		assert.Equal(t, initialSize, len(cache.cache))
	})

	t.Run("update moves item to head (most recently used)", func(t *testing.T) {
		cache := NewLRUCache(3)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key1", "updated")

		assert.Equal(t, "key1", cache.head.key)
	})
}

func TestLRUCache_Eviction(t *testing.T) {
	t.Run("evicts least recently used when capacity exceeded", func(t *testing.T) {
		cache := NewLRUCache(2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		_, ok := cache.Get("key1")
		assert.False(t, ok)
		assert.Equal(t, 2, len(cache.cache))
	})

	t.Run("evicts oldest item when filled to capacity", func(t *testing.T) {
		cache := NewLRUCache(3)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)
		cache.Set("d", 4)

		_, okA := cache.Get("a")
		_, okD := cache.Get("d")

		assert.False(t, okA)
		assert.True(t, okD)
	})

	t.Run("maintains capacity after multiple evictions", func(t *testing.T) {
		cache := NewLRUCache(2)
		for i := 1; i <= 5; i++ {
			cache.Set("key"+string(rune('0'+i)), i)
		}

		assert.Equal(t, 2, len(cache.cache))
	})

	t.Run("evicts correct item with mixed access patterns", func(t *testing.T) {
		cache := NewLRUCache(3)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)
		cache.Get("a")
		cache.Set("d", 4)

		_, okB := cache.Get("b")
		_, okD := cache.Get("d")

		assert.False(t, okB)
		assert.True(t, okD)
	})
}

func TestLRUCache_LinkedListIntegrity(t *testing.T) {
	t.Run("head and tail are correct after single set", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")

		assert.Equal(t, "key1", cache.head.key)
		assert.Equal(t, "key1", cache.tail.key)
	})

	t.Run("head and tail are correct after multiple sets", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		assert.Equal(t, "key3", cache.head.key)
		assert.Equal(t, "key1", cache.tail.key)
	})

	t.Run("linked list maintains order", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)

		// Traverse from head
		keys := []string{}
		current := cache.head
		for current != nil {
			keys = append(keys, current.key)
			current = current.next
		}

		assert.Equal(t, []string{"c", "b", "a"}, keys)
	})

	t.Run("linked list maintains reverse order from tail", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)

		// Traverse from tail
		keys := []string{}
		current := cache.tail
		for current != nil {
			keys = append(keys, current.key)
			current = current.prev
		}

		assert.Equal(t, []string{"a", "b", "c"}, keys)
	})
}

func TestLRUCache_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent sets and gets", func(t *testing.T) {
		cache := NewLRUCache(100)
		done := make(chan bool, 10)

		// Concurrent writers
		for i := 0; i < 5; i++ {
			go func(id int) {
				for j := 0; j < 10; j++ {
					cache.Set("key", id*10+j)
				}
				done <- true
			}(i)
		}

		// Concurrent readers
		for i := 0; i < 5; i++ {
			go func() {
				for j := 0; j < 10; j++ {
					cache.Get("key")
				}
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		assert.True(t, len(cache.cache) <= 100)
	})
}

func TestLRUCache_EdgeCases(t *testing.T) {
	t.Run("cache with capacity 1", func(t *testing.T) {
		cache := NewLRUCache(1)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")

		_, ok := cache.Get("key1")
		assert.False(t, ok)

		v, ok := cache.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, "value2", v)
	})

	t.Run("set empty string as key", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("", "empty_key")

		value, ok := cache.Get("")
		assert.True(t, ok)
		assert.Equal(t, "empty_key", value)
	})

	t.Run("set nil as value", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("key", nil)

		value, ok := cache.Get("key")
		assert.True(t, ok)
		assert.Nil(t, value)
	})
}

func TestLRUCache_MoveToHead(t *testing.T) {
	t.Run("accessing item moves it to head", func(t *testing.T) {
		cache := NewLRUCache(5)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)

		cache.Get("a") // Move "a" to head

		assert.Equal(t, "a", cache.head.key)
		assert.Equal(t, "b", cache.tail.key) // Fixed: tail is now "b", not "c"
	})

	t.Run("head changes after get of old item", func(t *testing.T) {
		cache := NewLRUCache(3)
		cache.Set("a", 1)
		cache.Set("b", 2)
		cache.Set("c", 3)

		oldHead := cache.head.key
		cache.Get("a")
		newHead := cache.head.key

		assert.NotEqual(t, oldHead, newHead)
		assert.Equal(t, "a", newHead)
		assert.Equal(t, "b", cache.tail.key) // Added: verify tail
	})
}
