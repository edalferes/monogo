package httpclient

import (
	"testing"
	"time"
)

func TestCircuitBreaker_StateClosed(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:             true,
		MaxConsecutiveFails: 3,
		Timeout:             1 * time.Second,
		HalfOpenRequests:    2,
	}
	cb := NewCircuitBreaker(cfg)

	// Initially should be closed
	if err := cb.CanAttempt(); err != nil {
		t.Errorf("Expected circuit to be closed, got error: %v", err)
	}
}

func TestCircuitBreaker_StateOpen(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:             true,
		MaxConsecutiveFails: 3,
		Timeout:             100 * time.Millisecond,
		HalfOpenRequests:    2,
	}
	cb := NewCircuitBreaker(cfg)

	// Record failures to open circuit
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// Circuit should be open
	if err := cb.CanAttempt(); err == nil {
		t.Error("Expected circuit to be open")
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Should allow attempt after timeout (transitioning to half-open)
	if err := cb.CanAttempt(); err != nil {
		t.Errorf("Expected circuit to allow attempt after timeout, got error: %v", err)
	}
}

func TestCircuitBreaker_StateHalfOpen(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:             true,
		MaxConsecutiveFails: 2,
		Timeout:             50 * time.Millisecond,
		HalfOpenRequests:    2,
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Circuit should allow half-open attempts
	cb.CanAttempt()

	// Record successful requests
	cb.RecordSuccess()
	cb.RecordSuccess()

	// Circuit should be closed again
	if err := cb.CanAttempt(); err != nil {
		t.Errorf("Expected circuit to be closed after successful half-open requests, got error: %v", err)
	}
}

func TestRetryLogic(t *testing.T) {
	cfg := DefaultConfig("http://test.example.com")
	cfg.RetryConfig.MaxRetries = 2
	cfg.RetryConfig.InitialBackoff = 10 * time.Millisecond

	client := NewClient(cfg)

	backoff1 := client.calculateBackoff(1)
	backoff2 := client.calculateBackoff(2)

	// Second backoff should be larger due to exponential increase
	if backoff2 <= backoff1 {
		t.Errorf("Expected exponential backoff, got backoff1=%v, backoff2=%v", backoff1, backoff2)
	}
}

func TestHTTPError(t *testing.T) {
	err := &HTTPError{
		StatusCode: 404,
		Body:       "Not Found",
	}

	expected := "HTTP 404: Not Found"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestClientConfig(t *testing.T) {
	cfg := DefaultConfig("http://localhost:8080")

	if cfg.RetryConfig.MaxRetries != 3 {
		t.Errorf("Expected 3 max retries, got %d", cfg.RetryConfig.MaxRetries)
	}

	if cfg.CircuitBreaker.Enabled != true {
		t.Error("Expected circuit breaker to be enabled by default")
	}

	if cfg.CircuitBreaker.MaxConsecutiveFails != 5 {
		t.Errorf("Expected 5 max consecutive fails, got %d", cfg.CircuitBreaker.MaxConsecutiveFails)
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{400, true},
		{404, true},
		{499, true},
		{500, false},
		{503, false},
	}

	for _, tt := range tests {
		err := &HTTPError{StatusCode: tt.statusCode}
		result := isClientError(err)
		if result != tt.expected {
			t.Errorf("For status %d, expected isClientError=%v, got %v", tt.statusCode, tt.expected, result)
		}
	}
}

func TestGetBaseURL(t *testing.T) {
	baseURL := "http://localhost:8080"
	cfg := DefaultConfig(baseURL)
	client := NewClient(cfg)

	if client.GetBaseURL() != baseURL {
		t.Errorf("Expected base URL %q, got %q", baseURL, client.GetBaseURL())
	}
}

func TestClientWithoutCircuitBreaker(t *testing.T) {
	cfg := DefaultConfig("http://localhost:8080")
	cfg.CircuitBreaker.Enabled = false

	client := NewClient(cfg)

	if client.circuitBreaker != nil {
		t.Error("Expected circuit breaker to be nil when disabled")
	}
}

func BenchmarkCircuitBreakerCanAttempt(b *testing.B) {
	cfg := CircuitBreakerConfig{
		Enabled:             true,
		MaxConsecutiveFails: 5,
		Timeout:             30 * time.Second,
		HalfOpenRequests:    3,
	}
	cb := NewCircuitBreaker(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.CanAttempt()
	}
}

func BenchmarkCalculateBackoff(b *testing.B) {
	cfg := DefaultConfig("http://test.example.com")
	client := NewClient(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.calculateBackoff(3)
	}
}
