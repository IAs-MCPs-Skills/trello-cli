package trello

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests
}

func isMutation(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return false
	default:
		return true
	}
}

func backoff(attempt int) time.Duration {
	base := math.Pow(2, float64(attempt)) * 500
	jitter := rand.Float64() * base * 0.5
	return time.Duration(base+jitter) * time.Millisecond
}

type retryWaitFunc func(context.Context, int) error

var waitForRetryFunc retryWaitFunc = defaultWaitForRetry

func waitForRetry(ctx context.Context, attempt int) error {
	return waitForRetryFunc(ctx, attempt)
}

func defaultWaitForRetry(ctx context.Context, attempt int) error {
	d := backoff(attempt)
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
