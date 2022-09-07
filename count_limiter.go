package main

import (
	"sync"
	"time"
)

type CounterLimiter struct {
	Limit    int
	Period   time.Duration
	count    int
	timeMark int64
	isStop   bool
	mu       sync.Mutex
}

func CounterLimiterNew(limit int, period time.Duration) *CounterLimiter {
	l := &CounterLimiter{
		Limit:    limit,
		Period:   period,
		count:    0,
		timeMark: time.Now().UnixNano(),
		isStop:   false,
	}
	return l
}

func (c *CounterLimiter) Wait() bool {
	return c.WaitN(1)
}

func (c *CounterLimiter) WaitN(weight int) bool {
	if c.Limit < weight {
		panic("weight too big, WaitN will be never return")
	}
	for !c.isStop {
		if c.TakeN(weight) {
			break
		}
	}
	return true
}

func (c *CounterLimiter) Take() bool {
	return c.TakeN(1)
}

func (c *CounterLimiter) TakeN(weight int) bool {
	if c.isStop {
		return true
	}
	c.checkTime()
	if c.count+weight <= c.Limit {
		c.count += weight
		return true
	}
	return false
}

func (c *CounterLimiter) Start() error {
	c.isStop = false
	return nil
}

func (c *CounterLimiter) Stop() error {
	c.isStop = true
	return nil
}

func (c *CounterLimiter) checkTime() {
	c.mu.Lock()
	defer c.mu.Unlock()
	unix := time.Now().UnixNano()
	if unix-c.timeMark > c.Period.Nanoseconds() {
		c.count = 0
		c.timeMark = unix
	}
}
