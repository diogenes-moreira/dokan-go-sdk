// Package dokan provides a Go SDK for the Dokan Multivendor Marketplace API.
//
// Dokan is a WordPress plugin that allows you to create multivendor marketplaces
// based on WooCommerce. This SDK provides a comprehensive Go interface to interact
// with the Dokan REST API.
//
// Basic usage:
//
//	client, err := dokan.NewClientBuilder().
//		BaseURL("https://example.com").
//		BasicAuth("username", "password").
//		Build()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create a product
//	product := &dokan.Product{
//		Name:         "Example Product",
//		Type:         dokan.ProductTypeSimple,
//		RegularPrice: "29.99",
//		Status:       dokan.ProductStatusPublish,
//	}
//
//	created, err := client.Products.Create(context.Background(), product)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Created product with ID: %d\n", created.ID)
package dokan

import (
	"github.com/diogenes-moreira/dokan-go-sdk/auth"
	"github.com/diogenes-moreira/dokan-go-sdk/client"
	"github.com/diogenes-moreira/dokan-go-sdk/errors"
	"github.com/diogenes-moreira/dokan-go-sdk/orders"
	"github.com/diogenes-moreira/dokan-go-sdk/stores"
	"github.com/diogenes-moreira/dokan-go-sdk/types"
)

// Re-export main types for easier access
type (
	// Client types
	Client        = client.Client
	Config        = client.Config
	ClientBuilder = client.ClientBuilder

	// Product types
	Product           = types.Product
	ProductType       = types.ProductType
	ProductStatus     = types.ProductStatus
	CatalogVisibility = types.CatalogVisibility
	ProductCategory   = types.ProductCategory
	ProductTag        = types.ProductTag
	ProductImage      = types.ProductImage
	ProductAttribute  = types.ProductAttribute
	ProductListParams = types.ProductListParams

	// Order types
	Order           = types.Order
	OrderStatus     = types.OrderStatus
	OrderListParams = types.OrderListParams
	OrderUpdate     = orders.OrderUpdate
	Address         = types.Address
	LineItem        = types.LineItem
	TaxLine         = types.TaxLine
	ShippingLine    = types.ShippingLine
	FeeLine         = types.FeeLine
	CouponLine      = types.CouponLine
	Refund          = types.Refund

	// Store types
	Store           = types.Store
	StoreListParams = types.StoreListParams
	Rating          = types.Rating

	// Common types
	MetaData     = types.MetaData
	ListParams   = types.ListParams
	ListResponse = types.ListResponse

	// Auth types
	AuthType      = auth.AuthType
	Authenticator = auth.Authenticator
	BasicAuth     = auth.BasicAuth
	JWTAuth       = auth.JWTAuth
	AuthConfig    = auth.Config

	// Review types
	ReviewListParams = stores.ReviewListParams
	Review           = stores.Review

	// Error types
	DokanError          = errors.DokanError
	NetworkError        = errors.NetworkError
	AuthenticationError = errors.AuthenticationError
	ValidationError     = errors.ValidationError
	NotFoundError       = errors.NotFoundError
	RateLimitError      = errors.RateLimitError
)

// Re-export constants
const (
	// Product types
	ProductTypeSimple   = types.ProductTypeSimple
	ProductTypeGrouped  = types.ProductTypeGrouped
	ProductTypeExternal = types.ProductTypeExternal
	ProductTypeVariable = types.ProductTypeVariable

	// Product statuses
	ProductStatusDraft   = types.ProductStatusDraft
	ProductStatusPending = types.ProductStatusPending
	ProductStatusPublish = types.ProductStatusPublish

	// Catalog visibility
	CatalogVisibilityVisible = types.CatalogVisibilityVisible
	CatalogVisibilityCatalog = types.CatalogVisibilityCatalog
	CatalogVisibilitySearch  = types.CatalogVisibilitySearch
	CatalogVisibilityHidden  = types.CatalogVisibilityHidden

	// Order statuses
	OrderStatusPending    = types.OrderStatusPending
	OrderStatusProcessing = types.OrderStatusProcessing
	OrderStatusOnHold     = types.OrderStatusOnHold
	OrderStatusCompleted  = types.OrderStatusCompleted
	OrderStatusCancelled  = types.OrderStatusCancelled
	OrderStatusRefunded   = types.OrderStatusRefunded
	OrderStatusFailed     = types.OrderStatusFailed

	// Auth types
	AuthTypeBasic = auth.AuthTypeBasic
	AuthTypeJWT   = auth.AuthTypeJWT
)

// Re-export main functions
var (
	// Client functions
	NewClient        = client.NewClient
	NewClientBuilder = client.NewClientBuilder
	DefaultConfig    = client.DefaultConfig

	// Auth functions
	NewBasicAuth     = auth.NewBasicAuth
	NewJWTAuth       = auth.NewJWTAuth
	NewAuthenticator = auth.NewAuthenticator

	// Error functions
	NewDokanError          = errors.NewDokanError
	NewNetworkError        = errors.NewNetworkError
	NewAuthenticationError = errors.NewAuthenticationError
	NewValidationError     = errors.NewValidationError
	NewNotFoundError       = errors.NewNotFoundError
	NewRateLimitError      = errors.NewRateLimitError
	IsDokanError           = errors.IsDokanError
	HandleHTTPError        = errors.HandleHTTPError
)
