package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("queue cache logic", func(t *testing.T) {
		c := NewCache(2)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		_, ok := c.Get("aaa")
		require.False(t, ok)
		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, val, 200)
		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, val, 300)

		c.Set("ddd", 400)

		_, ok = c.Get("bbb")
		require.False(t, ok)
		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, val, 300)
		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, val, 400)

		c.Set("eee", 500)

		_, ok = c.Get("ccc")
		require.False(t, ok)
		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, val, 400)
		val, ok = c.Get("eee")
		require.True(t, ok)
		require.Equal(t, val, 500)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 300)
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("cache size", func(t *testing.T) {
		c := NewCache(2)

		c.Set("11", 11)
		c.Set("12", 12)
		c.Set("13", 13)

		_, ok := c.Get("11")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
