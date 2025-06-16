package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk/auth"
	"github.com/diogenes-moreira/dokan-go-sdk/products"
	"github.com/diogenes-moreira/dokan-go-sdk/orders"
	"github.com/diogenes-moreira/dokan-go-sdk/stores"
	"github.com/diogenes-moreira/dokan-go-sdk/utils"
)

// Client is the main Dokan API client
type Client struct {
	baseURL    string
	httpClient utils.HTTPClient
	auth       auth.Authenticator
	retryConfig utils.RetryConfig
	
	// Services
	Products *products.Service
	Orders   *orders.Service
	Stores   *stores.Service
}

// Config represents client configuration
type Config struct {
	BaseURL     string
	Timeout     time.Duration
	RetryCount  int
	UserAgent   string
	Debug       bool
	Auth        auth.Config
	HTTPClient  *http.Client
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		UserAgent:  "dokan-go-sdk/1.0.0",
		Debug:      false,
	}
}

// NewClient creates a new Dokan API client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	// Validate required fields
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	
	// Create authenticator
	authenticator, err := auth.NewAuthenticator(config.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}
	
	// Create HTTP client if not provided
	var httpClient utils.HTTPClient
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	} else {
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}
		
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		}
	}
	
	// Create retry config
	retryConfig := utils.RetryConfig{
		MaxRetries: config.RetryCount,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
	}
	
	client := &Client{
		baseURL:     config.BaseURL,
		httpClient:  httpClient,
		auth:        authenticator,
		retryConfig: retryConfig,
	}
	
	// Initialize services
	client.Products = products.NewService(client)
	client.Orders = orders.NewService(client)
	client.Stores = stores.NewService(client)
	
	return client, nil
}

// MakeRequest makes an authenticated HTTP request
func (c *Client) MakeRequest(ctx context.Context, opts utils.RequestOptions) (*utils.Response, error) {
	var lastResponse *utils.Response
	var lastError error
	
	err := utils.WithRetry(ctx, c.retryConfig, func() error {
		// Create a new request for each retry attempt
		resp, err := utils.MakeRequest(ctx, &authenticatedClient{
			client: c.httpClient,
			auth:   c.auth,
		}, c.baseURL, opts)
		
		lastResponse = resp
		lastError = err
		return err
	})
	
	if err != nil {
		return lastResponse, err
	}
	
	return lastResponse, lastError
}

// authenticatedClient wraps an HTTP client with authentication
type authenticatedClient struct {
	client utils.HTTPClient
	auth   auth.Authenticator
}

// Do implements utils.HTTPClient interface with authentication
func (ac *authenticatedClient) Do(req *http.Request) (*http.Response, error) {
	// Add authentication to the request
	if err := ac.auth.Authenticate(req); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	
	return ac.client.Do(req)
}

// ClientBuilder provides a fluent interface for building clients
type ClientBuilder struct {
	config *Config
}

// NewClientBuilder creates a new client builder
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		config: DefaultConfig(),
	}
}

// BaseURL sets the base URL for the API
func (b *ClientBuilder) BaseURL(url string) *ClientBuilder {
	b.config.BaseURL = url
	return b
}

// Timeout sets the HTTP timeout
func (b *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder {
	b.config.Timeout = timeout
	return b
}

// RetryCount sets the number of retries for failed requests
func (b *ClientBuilder) RetryCount(count int) *ClientBuilder {
	b.config.RetryCount = count
	return b
}

// UserAgent sets the user agent string
func (b *ClientBuilder) UserAgent(userAgent string) *ClientBuilder {
	b.config.UserAgent = userAgent
	return b
}

// Debug enables or disables debug mode
func (b *ClientBuilder) Debug(debug bool) *ClientBuilder {
	b.config.Debug = debug
	return b
}

// BasicAuth configures HTTP Basic Authentication
func (b *ClientBuilder) BasicAuth(username, password string) *ClientBuilder {
	b.config.Auth = auth.Config{
		Type:     auth.AuthTypeBasic,
		Username: username,
		Password: password,
	}
	return b
}

// JWTAuth configures JWT Authentication
func (b *ClientBuilder) JWTAuth(token string) *ClientBuilder {
	b.config.Auth = auth.Config{
		Type:  auth.AuthTypeJWT,
		Token: token,
	}
	return b
}

// HTTPClient sets a custom HTTP client
func (b *ClientBuilder) HTTPClient(client *http.Client) *ClientBuilder {
	b.config.HTTPClient = client
	return b
}

// Build creates the client with the configured options
func (b *ClientBuilder) Build() (*Client, error) {
	return NewClient(b.config)
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// GetAuth returns the authenticator
func (c *Client) GetAuth() auth.Authenticator {
	return c.auth
}

// SetAuth sets a new authenticator
func (c *Client) SetAuth(authenticator auth.Authenticator) {
	c.auth = authenticator
}

