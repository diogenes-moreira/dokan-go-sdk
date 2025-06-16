package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestNewBasicAuth(t *testing.T) {
	auth := NewBasicAuth("user", "pass")
	
	if auth == nil {
		t.Fatal("NewBasicAuth() returned nil")
	}
	
	if auth.username != "user" {
		t.Errorf("Expected username 'user', got '%s'", auth.username)
	}
	
	if auth.password != "pass" {
		t.Errorf("Expected password 'pass', got '%s'", auth.password)
	}
}

func TestBasicAuth_Type(t *testing.T) {
	auth := NewBasicAuth("user", "pass")
	
	if auth.Type() != AuthTypeBasic {
		t.Errorf("Expected type %v, got %v", AuthTypeBasic, auth.Type())
	}
}

func TestBasicAuth_Authenticate(t *testing.T) {
	auth := NewBasicAuth("testuser", "testpass")
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	
	err := auth.Authenticate(req)
	if err != nil {
		t.Fatalf("Authenticate() returned error: %v", err)
	}
	
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		t.Error("Authorization header not set")
	}
	
	// Basic auth should start with "Basic "
	if len(authHeader) < 6 || authHeader[:6] != "Basic " {
		t.Errorf("Expected Authorization header to start with 'Basic ', got '%s'", authHeader)
	}
}

func TestBasicAuth_IsValid(t *testing.T) {
	// Valid auth
	auth := NewBasicAuth("user", "pass")
	if !auth.IsValid() {
		t.Error("Valid basic auth should return true")
	}
	
	// Invalid auth - empty username
	auth = NewBasicAuth("", "pass")
	if auth.IsValid() {
		t.Error("Basic auth with empty username should return false")
	}
	
	// Invalid auth - empty password
	auth = NewBasicAuth("user", "")
	if auth.IsValid() {
		t.Error("Basic auth with empty password should return false")
	}
}

func TestNewJWTAuth(t *testing.T) {
	token := "test-jwt-token"
	expiry := time.Now().Add(time.Hour)
	
	auth := NewJWTAuth(token, expiry)
	
	if auth == nil {
		t.Fatal("NewJWTAuth() returned nil")
	}
	
	if auth.token != token {
		t.Errorf("Expected token '%s', got '%s'", token, auth.token)
	}
	
	if !auth.expiresAt.Equal(expiry) {
		t.Errorf("Expected expiry %v, got %v", expiry, auth.expiresAt)
	}
}

func TestJWTAuth_Type(t *testing.T) {
	auth := NewJWTAuth("token", time.Now().Add(time.Hour))
	
	if auth.Type() != AuthTypeJWT {
		t.Errorf("Expected type %v, got %v", AuthTypeJWT, auth.Type())
	}
}

func TestJWTAuth_Authenticate(t *testing.T) {
	token := "test-jwt-token"
	auth := NewJWTAuth(token, time.Now().Add(time.Hour))
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	
	err := auth.Authenticate(req)
	if err != nil {
		t.Fatalf("Authenticate() returned error: %v", err)
	}
	
	authHeader := req.Header.Get("Authorization")
	expectedHeader := "Bearer " + token
	
	if authHeader != expectedHeader {
		t.Errorf("Expected Authorization header '%s', got '%s'", expectedHeader, authHeader)
	}
}

func TestJWTAuth_IsValid(t *testing.T) {
	// Valid auth
	auth := NewJWTAuth("valid-token", time.Now().Add(time.Hour))
	if !auth.IsValid() {
		t.Error("Valid JWT auth should return true")
	}
	
	// Invalid auth - empty token
	auth = NewJWTAuth("", time.Now().Add(time.Hour))
	if auth.IsValid() {
		t.Error("JWT auth with empty token should return false")
	}
	
	// Invalid auth - expired token
	auth = NewJWTAuth("valid-token", time.Now().Add(-time.Hour))
	if auth.IsValid() {
		t.Error("JWT auth with expired token should return false")
	}
}

func TestNewAuthenticator(t *testing.T) {
	config := Config{
		Type:     AuthTypeBasic,
		Username: "user",
		Password: "pass",
	}
	
	auth, err := NewAuthenticator(config)
	if err != nil {
		t.Fatalf("NewAuthenticator() returned error: %v", err)
	}
	
	if auth == nil {
		t.Fatal("NewAuthenticator() returned nil")
	}
	
	if auth.Type() != AuthTypeBasic {
		t.Errorf("Expected type %v, got %v", AuthTypeBasic, auth.Type())
	}
}

func TestNewAuthenticator_JWT(t *testing.T) {
	config := Config{
		Type:  AuthTypeJWT,
		Token: "test-token",
	}
	
	auth, err := NewAuthenticator(config)
	if err != nil {
		t.Fatalf("NewAuthenticator() returned error: %v", err)
	}
	
	if auth == nil {
		t.Fatal("NewAuthenticator() returned nil")
	}
	
	if auth.Type() != AuthTypeJWT {
		t.Errorf("Expected type %v, got %v", AuthTypeJWT, auth.Type())
	}
}

func TestNewAuthenticator_InvalidType(t *testing.T) {
	config := Config{
		Type: "invalid",
	}
	
	_, err := NewAuthenticator(config)
	if err == nil {
		t.Error("NewAuthenticator() should return error for invalid type")
	}
}

func TestNewAuthenticator_MissingCredentials(t *testing.T) {
	// Basic auth without credentials
	config := Config{
		Type: AuthTypeBasic,
	}
	
	_, err := NewAuthenticator(config)
	if err == nil {
		t.Error("NewAuthenticator() should return error for basic auth without credentials")
	}
	
	// JWT auth without token
	config = Config{
		Type: AuthTypeJWT,
	}
	
	_, err = NewAuthenticator(config)
	if err == nil {
		t.Error("NewAuthenticator() should return error for JWT auth without token")
	}
}

