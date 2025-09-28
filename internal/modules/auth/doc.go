/*
Package auth provides authentication and authorization functionality for the application.

This package implements a Clean Architecture approach with separate layers for:
- Domain entities (User, Role, Permission, AuditLog)
- Use cases with segregated interfaces (Reader/Writer pattern)
- HTTP handlers for REST API endpoints
- Repository implementations for data persistence

# Architecture Overview

The auth module follows Clean Architecture principles:

	Domain Layer (entities):
		- User: represents a system user with authentication data
		- Role: represents user roles for authorization
		- Permission: represents granular permissions
		- AuditLog: tracks sensitive operations for security

	Use Case Layer:
		- User operations: authentication, registration, profile management
		- Role operations: role assignment and management
		- Permission operations: permission checks and management
		- Audit operations: security event logging

	Interface Layer:
		- HTTP handlers for REST API endpoints
		- Request/response DTOs for data transfer
		- Middleware for JWT authentication

	Infrastructure Layer:
		- GORM repositories for database operations
		- JWT token management
		- Password hashing and validation

# Authentication Flow

The authentication process follows these steps:

 1. User submits credentials via POST /auth/login
 2. LoginWithAuditUseCase validates credentials
 3. JWT token is generated and returned
 4. Subsequent requests include JWT in Authorization header
 5. JWTMiddleware validates tokens and extracts user context
 6. Audit logs track all authentication events

# Authorization

Role-based access control (RBAC) is implemented through:

  - Users have multiple Roles
  - Roles have multiple Permissions
  - Permissions define granular access rights
  - Middleware checks permissions for protected endpoints

# Usage Examples

Basic authentication:

	// Initialize auth module
	authModule := auth.NewModule(db, logger, jwtConfig)
	authModule.RegisterRoutes(echoGroup)

	// Login request
	POST /v1/auth/login
	{
		"username": "admin",
		"password": "password123"
	}

	// Response
	{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"user": {
			"id": 1,
			"username": "admin",
			"roles": ["admin"]
		}
	}

Protected endpoint access:

	GET /v1/admin/users
	Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Security Features

- Password hashing using bcrypt
- JWT tokens with configurable expiration
- Audit logging for security events
- Role-based access control
- Request rate limiting (planned)
- Account lockout after failed attempts (planned)

# Database Schema

The module requires the following database tables:

	users: id, username, password_hash, created_at, updated_at
	roles: id, name, description, created_at, updated_at
	permissions: id, name, description, created_at, updated_at
	user_roles: user_id, role_id
	role_permissions: role_id, permission_id
	audit_logs: id, user_id, username, action, status, ip, details, created_at

# Configuration

Required configuration:

	jwt:
	  secret: "your-secret-key"
	  expiry_hour: 24

# Error Handling

Common authentication errors:

  - ErrInvalidCredentials: username/password mismatch
  - ErrUserNotFound: user does not exist
  - ErrInvalidToken: JWT token is invalid or expired
  - ErrInsufficientPermissions: user lacks required permissions
  - ErrInvalidData: request data validation failed

All errors are logged with appropriate context for debugging and security monitoring.
*/
package auth
