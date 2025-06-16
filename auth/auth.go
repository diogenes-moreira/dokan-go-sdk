package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// AuthType represents the type of authentication
type AuthType string

const (
	AuthTypeBasic AuthType = "basic"
	AuthTypeJWT   AuthType = "jwt"
)

// Authenticator interface defines methods for authentication
type Authenticator interface {
	Authenticate(req *http.Request) error
	IsValid() bool
	Refresh() error
	Type() AuthType
}

// BasicAuth implements HTTP Basic Authentication
type BasicAuth struct {
	username string
	password string
}

// NewBasicAuth creates a new BasicAuth authenticator
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

// Authenticate adds Basic Auth header to the request
func (b *BasicAuth) Authenticate(req *http.Request) error {
	if b.username == "" || b.password == "" {
		return fmt.Errorf("username and password are required for basic auth")
	}
	
	auth := b.username + ":" + b.password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encoded)
	return nil
}

// IsValid checks if the basic auth credentials are valid
func (b *BasicAuth) IsValid() bool {
	return b.username != "" && b.password != ""
}

// Refresh is a no-op for basic auth
func (b *BasicAuth) Refresh() error {
	return nil
}

// Type returns the authentication type
func (b *BasicAuth) Type() AuthType {
	return AuthTypeBasic
}

// JWTAuth implements JWT Authentication
type JWTAuth struct {
	token     string
	expiresAt time.Time
	refreshToken string
	refreshFunc func(refreshToken string) (string, time.Time, error)
}

// NewJWTAuth creates a new JWTAuth authenticator
func NewJWTAuth(token string, expiresAt time.Time) *JWTAuth {
	return &JWTAuth{
		token:     token,
		expiresAt: expiresAt,
	}
}

// NewJWTAuthWithRefresh creates a new JWTAuth authenticator with refresh capability
func NewJWTAuthWithRefresh(token string, expiresAt time.Time, refreshToken string, refreshFunc func(string) (string, time.Time, error)) *JWTAuth {
	return &JWTAuth{
		token:       token,
		expiresAt:   expiresAt,
		refreshToken: refreshToken,
		refreshFunc: refreshFunc,
	}
}

// Authenticate adds JWT Bearer token to the request
func (j *JWTAuth) Authenticate(req *http.Request) error {
	if j.token == "" {
		return fmt.Errorf("JWT token is required")
	}
	
	// Check if token is expired and try to refresh
	if !j.IsValid() && j.refreshFunc != nil {
		if err := j.Refresh(); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
	}
	
	if !j.IsValid() {
		return fmt.Errorf("JWT token is expired and cannot be refreshed")
	}
	
	req.Header.Set("Authorization", "Bearer "+j.token)
	return nil
}

// IsValid checks if the JWT token is still valid
func (j *JWTAuth) IsValid() bool {
	if j.token == "" {
		return false
	}
	
	// If no expiration time is set, assume it's valid
	if j.expiresAt.IsZero() {
		return true
	}
	
	// Add a 5-minute buffer before expiration
	return time.Now().Add(5 * time.Minute).Before(j.expiresAt)
}

// Refresh refreshes the JWT token using the refresh token
func (j *JWTAuth) Refresh() error {
	if j.refreshFunc == nil {
		return fmt.Errorf("no refresh function provided")
	}
	
	if j.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}
	
	newToken, expiresAt, err := j.refreshFunc(j.refreshToken)
	if err != nil {
		return err
	}
	
	j.token = newToken
	j.expiresAt = expiresAt
	return nil
}

// Type returns the authentication type
func (j *JWTAuth) Type() AuthType {
	return AuthTypeJWT
}

// SetToken updates the JWT token and expiration time
func (j *JWTAuth) SetToken(token string, expiresAt time.Time) {
	j.token = token
	j.expiresAt = expiresAt
}

// GetToken returns the current JWT token
func (j *JWTAuth) GetToken() string {
	return j.token
}

// SetRefreshToken sets the refresh token
func (j *JWTAuth) SetRefreshToken(refreshToken string) {
	j.refreshToken = refreshToken
}

// Config represents authentication configuration
type Config struct {
	Type         AuthType `json:"type"`
	Username     string   `json:"username,omitempty"`
	Password     string   `json:"password,omitempty"`
	Token        string   `json:"token,omitempty"`
	RefreshToken string   `json:"refresh_token,omitempty"`
}

// NewAuthenticator creates a new authenticator based on the config
func NewAuthenticator(config Config) (Authenticator, error) {
	switch config.Type {
	case AuthTypeBasic:
		if config.Username == "" || config.Password == "" {
			return nil, fmt.Errorf("username and password are required for basic auth")
		}
		return NewBasicAuth(config.Username, config.Password), nil
	case AuthTypeJWT:
		if config.Token == "" {
			return nil, fmt.Errorf("token is required for JWT auth")
		}
		return NewJWTAuth(config.Token, time.Time{}), nil
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", config.Type)
	}
}

