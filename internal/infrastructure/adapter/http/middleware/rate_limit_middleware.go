package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// rateLimitErrorResponse represents a rate limit error response.
type rateLimitErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Path      string `json:"path"`
	TraceID   string `json:"trace_id"`
}

// RateLimitMiddleware implements token bucket rate limiting per client.
type RateLimitMiddleware struct {
	clients    map[string]*rateClient
	mu         sync.RWMutex
	public     int
	login      int
	auth       int
}

type rateClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimitMiddleware creates a new RateLimitMiddleware.
// rates are in requests per minute.
func NewRateLimitMiddleware(public, login, auth int) *RateLimitMiddleware {
	m := &RateLimitMiddleware{
		clients: make(map[string]*rateClient),
		public:  public,
		login:   login,
		auth:    auth,
	}
	m.startCleanup(10 * time.Minute)
	return m
}

// startCleanup runs a background goroutine that removes stale entries.
func (m *RateLimitMiddleware) startCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			m.mu.Lock()
			for key, client := range m.clients {
				if time.Since(client.lastSeen) > 30*time.Minute {
					delete(m.clients, key)
				}
			}
			m.mu.Unlock()
		}
	}()
}

// getClient returns or creates a rate limiter for the given key.
func (m *RateLimitMiddleware) getClient(key string, ratePerMin int) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[key]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(ratePerMin)/60, ratePerMin)
		m.clients[key] = &rateClient{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	client.lastSeen = time.Now()
	return client.limiter
}

// RateLimit returns middleware that applies rate limiting.
func (m *RateLimitMiddleware) RateLimit(rateType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ratePerMin int
			var key string

			switch rateType {
			case "login":
				ratePerMin = m.login
				key = "login:" + getClientIP(r)
			case "public":
				ratePerMin = m.public
				key = "public:" + getClientIP(r)
			case "auth":
				ratePerMin = m.auth
				key = "auth:" + getClientIP(r)
			default:
				ratePerMin = m.public
				key = "public:" + getClientIP(r)
			}

			limiter := m.getClient(key, ratePerMin)
			if !limiter.Allow() {
				traceID := ""
				if rid, ok := r.Context().Value(RequestIDKey).(string); ok {
					traceID = rid
				}
				resp := rateLimitErrorResponse{
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Status:    http.StatusTooManyRequests,
					Error:     http.StatusText(http.StatusTooManyRequests),
					Code:      "RATE_LIMIT_EXCEEDED",
					Message:   "Rate limit exceeded. Please try again later.",
					Path:      r.URL.Path,
					TraceID:   traceID,
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(resp)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Extract IP from RemoteAddr (strip port)
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}
