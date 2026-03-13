package trello

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func overrideWaitForRetry(t *testing.T, fn retryWaitFunc) {
	t.Helper()
	prev := waitForRetryFunc
	waitForRetryFunc = fn
	t.Cleanup(func() {
		waitForRetryFunc = prev
	})
}

func useImmediateWait(t *testing.T) {
	overrideWaitForRetry(t, func(context.Context, int) error {
		return nil
	})
}

func useBlockingWait(t *testing.T) {
	overrideWaitForRetry(t, func(ctx context.Context, _ int) error {
		<-ctx.Done()
		return ctx.Err()
	})
}

func TestRetryOnRateLimit(t *testing.T) {
	useImmediateWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"id": "b1"})
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 3
	client := NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)
	if err != nil {
		t.Fatalf("Get() should succeed after retries, got: %v", err)
	}
	if result["id"] != "b1" {
		t.Errorf("id = %q, want %q", result["id"], "b1")
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("attempts = %d, want 3", atomic.LoadInt32(&attempts))
	}
}

func TestRetryExhausted(t *testing.T) {
	useImmediateWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 2
	client := NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Get(context.Background(), "/1/boards/b1", nil, &result)
	if err == nil {
		t.Fatal("Get() should fail after exhausting retries")
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("attempts = %d, want 3", atomic.LoadInt32(&attempts))
	}
}

func TestNoRetryOnMutations(t *testing.T) {
	useImmediateWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 3
	opts.RetryMutations = false
	client := NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	client.Post(context.Background(), "/1/lists", map[string]string{"name": "test"}, &result)

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on mutations)", atomic.LoadInt32(&attempts))
	}
}

func TestRetryOnMutationsWhenEnabled(t *testing.T) {
	useImmediateWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"id": "new1"})
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 3
	opts.RetryMutations = true
	client := NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	err := client.Post(context.Background(), "/1/lists", map[string]string{"name": "test"}, &result)
	if err != nil {
		t.Fatalf("Post() should succeed after retry, got: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("attempts = %d, want 2", atomic.LoadInt32(&attempts))
	}
}

func TestNoRetryOnNon429Errors(t *testing.T) {
	useImmediateWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 3
	client := NewClient(server.URL, "k", "t", opts)

	var result map[string]string
	client.Get(context.Background(), "/1/boards/nope", nil, &result)

	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on 404)", atomic.LoadInt32(&attempts))
	}
}

func TestRetryWaitRespectsContextCancellation(t *testing.T) {
	useBlockingWait(t)

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := DefaultClientOptions()
	opts.MaxRetries = 3
	client := NewClient(server.URL, "k", "t", opts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.Get(ctx, "/1/boards/b1", nil, nil)
	}()

	time.AfterFunc(5*time.Millisecond, cancel)

	if err := <-errCh; !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("attempts = %d, want 1 (cancelled before retry)", atomic.LoadInt32(&attempts))
	}
}
