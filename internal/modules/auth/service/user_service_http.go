package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edalferes/monetics/internal/contracts"
	"github.com/edalferes/monetics/pkg/httpclient"
)

// UserServiceHTTP implements UserService interface using HTTP calls
// This allows the auth module to be deployed separately as a microservice
type UserServiceHTTP struct {
	client *httpclient.Client
}

// NewUserServiceHTTP creates a new HTTP-based user service client
func NewUserServiceHTTP(baseURL string) contracts.UserService {
	cfg := httpclient.DefaultConfig(baseURL)
	return &UserServiceHTTP{
		client: httpclient.NewClient(cfg),
	}
}

// GetUserByID retrieves user information via HTTP call
func (s *UserServiceHTTP) GetUserByID(ctx context.Context, userID uint) (*contracts.UserInfo, error) {
	path := fmt.Sprintf("/v1/auth/users/%d", userID)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("auth service returned error: %s", resp.Error)
	}

	var userInfo contracts.UserInfo
	if err := json.Unmarshal(resp.Data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// ValidateUserExists checks if a user exists via HTTP call
func (s *UserServiceHTTP) ValidateUserExists(ctx context.Context, userID uint) (bool, error) {
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		// User not found or service error - treat as not exists
		return false, nil
	}
	return true, nil
}

// GetUserPermissions retrieves user permissions via HTTP call
func (s *UserServiceHTTP) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	path := fmt.Sprintf("/v1/auth/users/%d/permissions", userID)

	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("auth service returned error: %s", resp.Error)
	}

	var permissions []string
	if err := json.Unmarshal(resp.Data, &permissions); err != nil {
		return nil, fmt.Errorf("failed to decode permissions: %w", err)
	}

	return permissions, nil
}
