package main

type Limiter interface {
	Wait() bool
	WaitN(int) bool
	Take() bool
	TakeN(int) bool
	Start() error
	Stop() error
}
