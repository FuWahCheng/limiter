# limiter

_Simple limiter package for Go._

## Installation
```bash
$ go get github.com/FuWahCheng/limiter
```

## Usage
1. Choose suitable limiter implement
2. New limiter instance with limit config
3. Use it according to `Limiter` interface function

### Implements

- CounterLimiter
- SlidingWindowLimiter
- LeakyBucketLimiter
- TokenBucketLimiter

**Example:**
```go
func main() {
	l := CounterLimiterNew(10, time.Second)
	for true {
		if l.Take() {
			log.Println("take")
		} else {
			log.Println("not take")
		}
		time.Sleep(time.Millisecond * 200)
	}
}
```