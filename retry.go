package duffel

import (
	"math"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

// RetryCond is a condition that applies only to retry backoff mechanism.
type RetryCond func(resp *http.Response, err error) bool

// RetryFunc takes attemps number, minimal and maximal wait time for backoff.
// Returns duration that mechanism have to wait before making a request.
type RetryFunc func(n int, min, max time.Duration) time.Duration

// backoff is a thread-safe retry backoff mechanism.
// Currently supported only ExponentalBackoff retry algorithm.
type backoff struct {
	minWaitTime time.Duration
	maxWaitTime time.Duration
	maxAttempts int32
	attempts    int32
	f           RetryFunc
}

const stopBackoff time.Duration = -1

func (b *backoff) next() time.Duration {
	if atomic.LoadInt32(&b.attempts) >= b.maxAttempts {
		return stopBackoff
	}
	atomic.AddInt32(&b.attempts, 1)
	return b.f(int(atomic.LoadInt32(&b.attempts)), b.minWaitTime, b.maxWaitTime)
}

func (b *backoff) reset() {
	atomic.SwapInt32(&b.attempts, 0)
}

func ExponentalBackoff(attemptNum int, min, max time.Duration) time.Duration {
	const factor = 2.0
	rand.Seed(time.Now().UnixNano())
	delay := time.Duration(math.Pow(factor, float64(attemptNum)) * float64(min))
	jitter := time.Duration(rand.Float64() * float64(min) * float64(attemptNum))

	delay = delay + jitter
	if delay > max {
		delay = max
	}

	return delay
}
