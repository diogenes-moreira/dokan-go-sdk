package products

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diogenes-moreira/dokan-go-sdk/types"
	"github.com/diogenes-moreira/dokan-go-sdk/utils"
)

// Service provides methods for interacting with the Dokan Products API
type Service struct {
	client ClientInterface
}

// ClientInterface defines the interface for making HTTP requests
type ClientInterface interface {
	MakeRequest(ctx context.Context, opts utils.RequestOptions) (*utils.Response, error)
}

// NewService creates a new products service
func NewService(client ClientInterface) *Service {
	return &Service{client: client}
}

// Create creates a new product in the Dokan marketplace
func (s *Service) Create(ctx context.Context, product *types.Product) (*types.Product, error) {
	opts := utils.RequestOptions{
		Method: http.MethodPost,
		Path:   "/wp-json/dokan/v1/products/",
		Body:   product,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	var createdProduct types.Product
	if err := utils.ParseJSON(resp.Body, &createdProduct); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &createdProduct, nil
}

// Get retrieves a single product by ID
func (s *Service) Get(ctx context.Context, id int) (*types.Product, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/products/%d", id),
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	var product types.Product
	if err := utils.ParseJSON(resp.Body, &product); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &product, nil
}

// List retrieves a list of products with optional filtering
func (s *Service) List(ctx context.Context, params *types.ProductListParams) (*ProductListResponse, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/wp-json/dokan/v1/products/",
		Query:  params,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	
	var products []types.Product
	if err := utils.ParseJSON(resp.Body, &products); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract pagination info from headers
	listResponse := &ProductListResponse{
		Products: products,
		ListResponse: types.ListResponse{
			TotalItems: extractIntHeader(resp.Headers, "X-WP-Total"),
			TotalPages: extractIntHeader(resp.Headers, "X-WP-TotalPages"),
		},
	}
	
	if params != nil {
		listResponse.Page = params.Page
		listResponse.PerPage = params.PerPage
	}
	
	return listResponse, nil
}

// Update updates an existing product
func (s *Service) Update(ctx context.Context, id int, product *types.Product) (*types.Product, error) {
	opts := utils.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/products/%d", id),
		Body:   product,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	var updatedProduct types.Product
	if err := utils.ParseJSON(resp.Body, &updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &updatedProduct, nil
}

// Delete deletes a product by ID
func (s *Service) Delete(ctx context.Context, id int) error {
	opts := utils.RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/products/%d", id),
	}
	
	_, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	
	return nil
}

// GetSummary retrieves a summary of products
func (s *Service) GetSummary(ctx context.Context) (*ProductSummary, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/wp-json/dokan/v1/products/summary",
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get product summary: %w", err)
	}
	
	var summary ProductSummary
	if err := utils.ParseJSON(resp.Body, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &summary, nil
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Products []types.Product `json:"products"`
	types.ListResponse
}

// ProductSummary represents a summary of products
type ProductSummary struct {
	Total     int `json:"total"`
	Published int `json:"published"`
	Draft     int `json:"draft"`
	Pending   int `json:"pending"`
	Featured  int `json:"featured"`
}

// extractIntHeader extracts an integer value from HTTP headers
func extractIntHeader(headers http.Header, key string) int {
	value := headers.Get(key)
	if value == "" {
		return 0
	}
	
	// Simple conversion, in a real implementation you might want better error handling
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

