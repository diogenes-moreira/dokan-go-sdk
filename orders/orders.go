package orders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/diogenes-moreira/dokan-go-sdk/types"
	"github.com/diogenes-moreira/dokan-go-sdk/utils"
)

// Service provides methods for interacting with the Dokan Orders API
type Service struct {
	client ClientInterface
}

// ClientInterface defines the interface for making HTTP requests
type ClientInterface interface {
	MakeRequest(ctx context.Context, opts utils.RequestOptions) (*utils.Response, error)
}

// NewService creates a new orders service
func NewService(client ClientInterface) *Service {
	return &Service{client: client}
}

// Get retrieves a single order by ID
func (s *Service) Get(ctx context.Context, id int) (*types.Order, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/orders/%d", id),
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	var order types.Order
	if err := utils.ParseJSON(resp.Body, &order); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &order, nil
}

// List retrieves a list of orders with optional filtering
func (s *Service) List(ctx context.Context, params *types.OrderListParams) (*OrderListResponse, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/wp-json/dokan/v1/orders/",
		Query:  params,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	var orders []types.Order
	if err := utils.ParseJSON(resp.Body, &orders); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Extract pagination info from headers
	listResponse := &OrderListResponse{
		Orders: orders,
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

// Update updates an existing order
func (s *Service) Update(ctx context.Context, id int, order *OrderUpdate) (*types.Order, error) {
	opts := utils.RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/wp-json/dokan/v1/orders/%d", id),
		Body:   order,
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	
	var updatedOrder types.Order
	if err := utils.ParseJSON(resp.Body, &updatedOrder); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &updatedOrder, nil
}

// GetSummary retrieves a summary of orders
func (s *Service) GetSummary(ctx context.Context) (*OrderSummary, error) {
	opts := utils.RequestOptions{
		Method: http.MethodGet,
		Path:   "/wp-json/dokan/v1/orders/summary",
	}
	
	resp, err := s.client.MakeRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get order summary: %w", err)
	}
	
	var summary OrderSummary
	if err := utils.ParseJSON(resp.Body, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &summary, nil
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	Orders []types.Order `json:"orders"`
	types.ListResponse
}

// OrderUpdate represents fields that can be updated in an order
type OrderUpdate struct {
	Status       *types.OrderStatus `json:"status,omitempty"`
	CustomerNote *string           `json:"customer_note,omitempty"`
	Billing      *types.Address    `json:"billing,omitempty"`
	Shipping     *types.Address    `json:"shipping,omitempty"`
	LineItems    []types.LineItem  `json:"line_items,omitempty"`
	ShippingLines []types.ShippingLine `json:"shipping_lines,omitempty"`
	FeeLines     []types.FeeLine   `json:"fee_lines,omitempty"`
	CouponLines  []types.CouponLine `json:"coupon_lines,omitempty"`
	MetaData     []types.MetaData  `json:"meta_data,omitempty"`
}

// OrderSummary represents a summary of orders
type OrderSummary struct {
	Total      int                        `json:"total"`
	Totals     map[string]int            `json:"totals"`
	StatusCounts map[types.OrderStatus]int `json:"status_counts"`
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

