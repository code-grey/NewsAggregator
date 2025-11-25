package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	// Create a dummy handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap it with the middleware
	handlerToTest := securityHeadersMiddleware(nextHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	handlerToTest.ServeHTTP(rr, req)

	// Check the headers
	assert.Equal(t, "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';", rr.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "max-age=63072000; includeSubDomains", rr.Header().Get("Strict-Transport-Security"))
}

func TestRateLimitMiddleware(t *testing.T) {
	// Reset the global limiter for this test to ensure isolation.
	limiter = rate.NewLimiter(2, 10)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := rateLimitMiddleware(nextHandler)

	// 1. Test that the first 10 requests are allowed
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/some-path", nil)
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "request %d should be allowed", i+1)
	}

	// 2. Test that the 11th request is rate-limited
	req := httptest.NewRequest("GET", "/some-path", nil)
	rr := httptest.NewRecorder()
	handlerToTest.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code, "11th request should be rate-limited")

	// 3. Test that the /healthz endpoint is never rate-limited, even when the limiter is exhausted
	healthzReq := httptest.NewRequest("GET", "/healthz", nil)
	healthzRr := httptest.NewRecorder()
	handlerToTest.ServeHTTP(healthzRr, healthzReq)
	assert.Equal(t, http.StatusOK, healthzRr.Code, "/healthz endpoint should not be rate-limited")
}

// Note: Testing loggingMiddleware typically involves capturing log output,
// which can be complex. A simple smoke test is to ensure it calls the next handler.
func TestLoggingMiddleware(t *testing.T) {
	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handlerToTest := loggingMiddleware(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handlerToTest.ServeHTTP(rr, req)

	assert.True(t, called, "next handler was not called")
	assert.Equal(t, http.StatusOK, rr.Code)
}
