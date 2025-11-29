package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Client represents an HTTP client for inter-service communication
// Inspired by Loki's approach to HTTP client management
type Client struct {
	httpClient     *http.Client
	baseURL        string
	retryConfig    RetryConfig
	circuitBreaker *CircuitBreaker
}

// Config holds HTTP client configuration
type Config struct {
	BaseURL         string
	Timeout         time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
	RetryConfig     RetryConfig
	CircuitBreaker  CircuitBreakerConfig
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
	Jitter         bool
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled             bool
	MaxConsecutiveFails int
	Timeout             time.Duration
	HalfOpenRequests    int
}

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu                  sync.RWMutex
	state               CircuitBreakerState
	consecutiveFails    int
	maxConsecutiveFails int
	lastFailTime        time.Time
	timeout             time.Duration
	halfOpenRequests    int
	halfOpenSuccess     int
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		state:               StateClosed,
		maxConsecutiveFails: cfg.MaxConsecutiveFails,
		timeout:             cfg.Timeout,
		halfOpenRequests:    cfg.HalfOpenRequests,
	}
}

// CanAttempt checks if a request can be attempted
func (cb *CircuitBreaker) CanAttempt() error {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == StateClosed {
		return nil
	}

	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.timeout {
			return nil // Allow transition to half-open
		}
		return errors.New("circuit breaker is open")
	}

	// StateHalfOpen - allow limited requests
	return nil
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.halfOpenSuccess++
		if cb.halfOpenSuccess >= cb.halfOpenRequests {
			cb.state = StateClosed
			cb.consecutiveFails = 0
			cb.halfOpenSuccess = 0
		}
	} else {
		cb.consecutiveFails = 0
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.consecutiveFails++
	cb.lastFailTime = time.Now()

	if cb.state == StateHalfOpen {
		cb.state = StateOpen
		cb.halfOpenSuccess = 0
	} else if cb.consecutiveFails >= cb.maxConsecutiveFails {
		cb.state = StateOpen
	}
}

// DefaultConfig returns default HTTP client configuration
func DefaultConfig(baseURL string) Config {
	return Config{
		BaseURL:         baseURL,
		Timeout:         10 * time.Second,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		RetryConfig: RetryConfig{
			MaxRetries:     3,
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     5 * time.Second,
			Multiplier:     2.0,
			Jitter:         true,
		},
		CircuitBreaker: CircuitBreakerConfig{
			Enabled:             true,
			MaxConsecutiveFails: 5,
			Timeout:             30 * time.Second,
			HalfOpenRequests:    3,
		},
	}
}

// NewClient creates a new HTTP client with the given configuration
func NewClient(cfg Config) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        cfg.MaxIdleConns,
				MaxIdleConnsPerHost: cfg.MaxIdleConns,
				IdleConnTimeout:     cfg.IdleConnTimeout,
			},
		},
		baseURL:     cfg.BaseURL,
		retryConfig: cfg.RetryConfig,
	}

	if cfg.CircuitBreaker.Enabled {
		client.circuitBreaker = NewCircuitBreaker(cfg.CircuitBreaker)
	}

	return client
}

// Response represents a standard API response
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// Get performs a GET request to the specified path
func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request to the specified path with the given body
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// Put performs a PUT request to the specified path with the given body
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// Delete performs a DELETE request to the specified path
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil)
}

// doRequest executes an HTTP request with retry and circuit breaker
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	// Check circuit breaker
	if c.circuitBreaker != nil {
		if err := c.circuitBreaker.CanAttempt(); err != nil {
			return nil, err
		}
	}

	var lastErr error
	maxAttempts := c.retryConfig.MaxRetries + 1

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Apply backoff delay for retries
		if attempt > 0 {
			backoff := c.calculateBackoff(attempt)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := c.executeRequest(ctx, method, path, body)
		if err == nil {
			// Success - record in circuit breaker
			if c.circuitBreaker != nil {
				c.circuitBreaker.RecordSuccess()
			}
			return resp, nil
		}

		lastErr = err

		// Don't retry on context errors
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			break
		}

		// Don't retry on 4xx errors (client errors)
		if isClientError(err) {
			break
		}
	}

	// All retries failed - record in circuit breaker
	if c.circuitBreaker != nil {
		c.circuitBreaker.RecordFailure()
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxAttempts, lastErr)
}

// executeRequest performs a single HTTP request attempt
func (c *Client) executeRequest(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	var response Response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// calculateBackoff calculates exponential backoff with jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	backoff := float64(c.retryConfig.InitialBackoff) * math.Pow(c.retryConfig.Multiplier, float64(attempt-1))

	if backoff > float64(c.retryConfig.MaxBackoff) {
		backoff = float64(c.retryConfig.MaxBackoff)
	}

	if c.retryConfig.Jitter {
		// Add random jitter between 0% and 25%
		jitter := rand.Float64() * 0.25 * backoff
		backoff += jitter
	}

	return time.Duration(backoff)
}

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// isClientError checks if error is a 4xx client error
func isClientError(err error) bool {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode >= 400 && httpErr.StatusCode < 500
	}
	return false
}

// GetBaseURL returns the base URL configured for this client
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
