package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk/auth"
	"github.com/diogenes-moreira/dokan-go-sdk/errors"
	"github.com/diogenes-moreira/dokan-go-sdk/utils"
)

func TestNewClientBuilder(t *testing.T) {
	builder := NewClientBuilder()
	if builder == nil {
		t.Fatal("NewClientBuilder() returned nil")
	}
}

func TestClientBuilder_BaseURL(t *testing.T) {
	builder := NewClientBuilder()
	result := builder.BaseURL("https://example.com")
	
	if result != builder {
		t.Error("BaseURL() should return the same builder instance")
	}
	
	if builder.config.BaseURL != "https://example.com" {
		t.Errorf("Expected BaseURL to be 'https://example.com', got '%s'", builder.config.BaseURL)
	}
}

func TestClientBuilder_BasicAuth(t *testing.T) {
	builder := NewClientBuilder()
	result := builder.BasicAuth("user", "pass")
	
	if result != builder {
		t.Error("BasicAuth() should return the same builder instance")
	}
	
	if builder.config.Auth.Type != auth.AuthTypeBasic {
		t.Errorf("Expected auth type to be Basic, got %v", builder.config.Auth.Type)
	}
	
	if builder.config.Auth.Username != "user" {
		t.Errorf("Expected username to be 'user', got '%s'", builder.config.Auth.Username)
	}
	
	if builder.config.Auth.Password != "pass" {
		t.Errorf("Expected password to be 'pass', got '%s'", builder.config.Auth.Password)
	}
}

func TestClientBuilder_JWTAuth(t *testing.T) {
	builder := NewClientBuilder()
	result := builder.JWTAuth("test-token")
	
	if result != builder {
		t.Error("JWTAuth() should return the same builder instance")
	}
	
	if builder.config.Auth.Type != auth.AuthTypeJWT {
		t.Errorf("Expected auth type to be JWT, got %v", builder.config.Auth.Type)
	}
	
	if builder.config.Auth.Token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got '%s'", builder.config.Auth.Token)
	}
}

func TestClientBuilder_Timeout(t *testing.T) {
	builder := NewClientBuilder()
	timeout := 30 * time.Second
	result := builder.Timeout(timeout)
	
	if result != builder {
		t.Error("Timeout() should return the same builder instance")
	}
	
	if builder.config.Timeout != timeout {
		t.Errorf("Expected timeout to be %v, got %v", timeout, builder.config.Timeout)
	}
}

func TestClientBuilder_RetryCount(t *testing.T) {
	builder := NewClientBuilder()
	result := builder.RetryCount(5)
	
	if result != builder {
		t.Error("RetryCount() should return the same builder instance")
	}
	
	if builder.config.RetryCount != 5 {
		t.Errorf("Expected retry count to be 5, got %d", builder.config.RetryCount)
	}
}

func TestClientBuilder_Build_Success(t *testing.T) {
	builder := NewClientBuilder()
	client, err := builder.
		BaseURL("https://example.com").
		BasicAuth("user", "pass").
		Build()
	
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}
	
	if client == nil {
		t.Fatal("Build() returned nil client")
	}
	
	if client.baseURL != "https://example.com" {
		t.Errorf("Expected client baseURL to be 'https://example.com', got '%s'", client.baseURL)
	}
}

func TestClientBuilder_Build_MissingBaseURL(t *testing.T) {
	builder := NewClientBuilder()
	_, err := builder.
		BasicAuth("user", "pass").
		Build()
	
	if err == nil {
		t.Error("Build() should return error when BaseURL is missing")
	}
}

func TestClientBuilder_Build_MissingAuth(t *testing.T) {
	builder := NewClientBuilder()
	_, err := builder.
		BaseURL("https://example.com").
		Build()
	
	if err == nil {
		t.Error("Build() should return error when authentication is missing")
	}
}

func TestClient_MakeRequest_Integration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Return a simple JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "Test Product"}`))
	}))
	defer server.Close()
	
	// Create client
	client, err := NewClientBuilder().
		BaseURL(server.URL).
		BasicAuth("testuser", "testpass").
		Build()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Make request
	ctx := context.Background()
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	}
	
	resp, err := client.MakeRequest(ctx, opts)
	if err != nil {
		t.Fatalf("MakeRequest() returned error: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	
	expectedBody := `{"id": 123, "name": "Test Product"}`
	if string(resp.Body) != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, string(resp.Body))
	}
}

func TestClient_MakeRequest_Unauthorized(t *testing.T) {
	// Create a test server that always returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": "unauthorized", "message": "Invalid credentials"}`))
	}))
	defer server.Close()
	
	// Create client
	client, err := NewClientBuilder().
		BaseURL(server.URL).
		BasicAuth("wronguser", "wrongpass").
		Build()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Make request
	ctx := context.Background()
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	}
	
	_, err = client.MakeRequest(ctx, opts)
	if err == nil {
		t.Error("MakeRequest() should return error for unauthorized request")
	}
	
	// Check if it's a Dokan error with 401 status
	if dokanErr, ok := err.(*errors.DokanError); ok {
		if dokanErr.StatusCode != 401 {
			t.Errorf("Expected status code 401, got %d", dokanErr.StatusCode)
		}
	} else {
		t.Errorf("Expected DokanError, got %T", err)
	}
}

