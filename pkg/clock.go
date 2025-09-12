package pkg

import (
	"sync"
)

type LamportClock struct {
	value int64
	mu    sync.Mutex
}

func NewLamportClock() *LamportClock {
	return &LamportClock{
		value: 0,
	}
}

func (c *LamportClock) Now() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *LamportClock) Tick() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// clock = max(local, remote) + 1
func (c *LamportClock) OnReceive(remote int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if remote > c.value {
		c.value = remote
	}
	c.value++

	return c.value
}
