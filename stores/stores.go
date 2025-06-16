package stores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diogenes-moreira/dokan-go-sdk/types"
	"github.com/diogenes-moreira/dokan-go-sdk/utils"
)

// Service provides methods for interacting with the Dokan Stores API
type Service struct {
	client ClientInterface
}

// ClientInterface defines the interface for making HTTP requests
type ClientInterface interface {
	MakeRequest(ctx context.Context, opts utils.RequestOptions) (*utils.Response, error)
}

// NewService creates a new stores service
func NewService(client ClientInterface) *Service {
	return &Service{client: client}
}

// Get retrieves a single store by vendor ID
func (s *Service) Get(ctx context.Context, vendorID int) (*types.Store, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/stores/%d", vendorID),
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	
	var store types.Store
	if err := utils.ParseJSON(resp.Body, &store); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &store, nil
}

// List retrieves a list of stores with optional filtering
func (s *Service) List(ctx context.Context, params *types.StoreListParams) (*StoreListResponse, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/wp-json/dokan/v1/stores/",
		Query:  params,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}
	
	var stores []types.Store
	if err := utils.ParseJSON(resp.Body, &stores); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract pagination info from headers
	listResponse := &StoreListResponse{
		Stores: stores,
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

// GetProducts retrieves products for a specific store
func (s *Service) GetProducts(ctx context.Context, vendorID int, params *types.ProductListParams) (*StoreProductsResponse, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/stores/%d/products", vendorID),
		Query:  params,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get store products: %w", err)
	}
	
	var products []types.Product
	if err := utils.ParseJSON(resp.Body, &products); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract pagination info from headers
	listResponse := &StoreProductsResponse{
		Products: products,
		VendorID: vendorID,
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

// GetReviews retrieves reviews for a specific store
func (s *Service) GetReviews(ctx context.Context, vendorID int, params *ReviewListParams) (*StoreReviewsResponse, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/stores/%d/reviews", vendorID),
		Query:  params,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get store reviews: %w", err)
	}
	
	var reviews []Review
	if err := utils.ParseJSON(resp.Body, &reviews); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract pagination info from headers
	listResponse := &StoreReviewsResponse{
		Reviews:  reviews,
		VendorID: vendorID,
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

// StoreListResponse represents a paginated list of stores
type StoreListResponse struct {
	Stores []types.Store `json:"stores"`
	types.ListResponse
}

// StoreProductsResponse represents a paginated list of products for a store
type StoreProductsResponse struct {
	Products []types.Product `json:"products"`
	VendorID int            `json:"vendor_id"`
	types.ListResponse
}

// StoreReviewsResponse represents a paginated list of reviews for a store
type StoreReviewsResponse struct {
	Reviews  []Review `json:"reviews"`
	VendorID int      `json:"vendor_id"`
	types.ListResponse
}

// Review represents a store review
type Review struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	Status       string `json:"status"`
	Reviewer     string `json:"reviewer"`
	ReviewerEmail string `json:"reviewer_email"`
	Review       string `json:"review"`
	Rating       int    `json:"rating"`
	Verified     bool   `json:"verified"`
	DateCreated  string `json:"date_created"`
	DateCreatedGMT string `json:"date_created_gmt"`
}

// ReviewListParams represents parameters for listing reviews
type ReviewListParams struct {
	types.ListParams
	Product  int    `url:"product,omitempty"`
	Status   string `url:"status,omitempty"`
	Reviewer string `url:"reviewer,omitempty"`
	Rating   int    `url:"rating,omitempty"`
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

