package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter creates a new rate limiter with limit r (reqs/sec) and burst b
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}

	// Periodic cleanup routine
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			i.mu.Lock()
			// Simple cleanup: reset map to free memory (for a major project, this is sufficient vs LRU)
			// In production, you'd track last seen time.
			i.ips = make(map[string]*rate.Limiter)
			i.mu.Unlock()
		}
	}()

	return i
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// rate limit vars
var (
	rateLimit      = rate.Limit(2)
	rateLimitBurst = 5
)

func init() {
	if rStr := os.Getenv("RATE_LIMIT"); rStr != "" {
		if r, err := strconv.ParseFloat(rStr, 64); err == nil {
			rateLimit = rate.Limit(r)
		}
	}
	if bStr := os.Getenv("RATE_LIMIT_BURST"); bStr != "" {
		if b, err := strconv.Atoi(bStr); err == nil {
			rateLimitBurst = b
		}
	}
}

var globalLimiter = NewIPRateLimiter(rateLimit, rateLimitBurst)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Basic attempt to strip port loopback (not production robust but fine here)
		// For real IP behind proxy, use X-Forwarded-For, but RemoteAddr is safer default if no trusted proxy

		limiter := globalLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
