package main

import (
	"container/list"
	"time"
)

type SlidingWindowLimiter struct {
	limit     int
	period    time.Duration
	windowNum int
	queue     *list.List
	isStop    bool
}

type queueElement struct {
	unix   int64
	weight int
}

func SlidingWindowLimiterNew(limit int, period time.Duration) *SlidingWindowLimiter {
	l := &SlidingWindowLimiter{
		limit:     limit,
		period:    period,
		windowNum: 0,
		queue:     list.New(),
		isStop:    false,
	}
	return l
}

func (s *SlidingWindowLimiter) Wait() bool {
	return s.WaitN(1)
}

func (s *SlidingWindowLimiter) WaitN(weight int) bool {
	if s.limit < weight {
		panic("weight too big, WaitN will be never return")
	}
	for !s.isStop {
		if s.TakeN(weight) {
			break
		}
	}
	return true
}

func (s *SlidingWindowLimiter) Take() bool {
	return s.TakeN(1)
}

func (s *SlidingWindowLimiter) TakeN(weight int) bool {
	s.clearExpire()
	unix := time.Now().UnixNano()
	if s.windowNum+weight <= s.limit {
		ele := queueElement{
			unix:   unix,
			weight: weight,
		}
		s.queue.PushBack(ele)
		s.windowNum += weight
		return true
	}
	return false
}

func (s *SlidingWindowLimiter) Start() error {
	s.isStop = false
	return nil
}

func (s *SlidingWindowLimiter) Stop() error {
	s.isStop = true
	return nil
}

func (s *SlidingWindowLimiter) clearExpire() {
	for true {
		front := s.queue.Front()
		if front != nil && (time.Now().UnixNano()-front.Value.(queueElement).unix) > s.period.Nanoseconds() {
			s.queue.Remove(front)
			s.windowNum -= front.Value.(queueElement).weight
		} else {
			break
		}
	}
}
