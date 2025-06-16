package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk"
)

func main() {
	// Example of JWT Authentication
	client, err := dokan.NewClientBuilder().
		BaseURL("https://your-dokan-site.com").
		JWTAuth("your-jwt-token").
		Timeout(60 * time.Second).
		RetryCount(5).
		Debug(true).
		Build()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Advanced error handling
	fmt.Println("=== Advanced Error Handling ===")
	_, err = client.Products.Get(ctx, 99999) // Non-existent product
	if err != nil {
		handleError(err)
	}

	// Example 2: Working with stores and their products
	fmt.Println("\n=== Store Management ===")
	stores, err := client.Stores.List(ctx, &dokan.StoreListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 10,
		},
		Enabled: &[]bool{true}[0],
	})
	if err != nil {
		log.Printf("Failed to list stores: %v", err)
		return
	}

	if len(stores.Stores) > 0 {
		store := stores.Stores[0]
		fmt.Printf("Working with store: %s (ID: %d)\n", store.StoreName, store.ID)

		// Get products for this store
		storeProducts, err := client.Stores.GetProducts(ctx, store.ID, &dokan.ProductListParams{
			ListParams: dokan.ListParams{
				Page:    1,
				PerPage: 5,
			},
		})
		if err != nil {
			log.Printf("Failed to get store products: %v", err)
		} else {
			fmt.Printf("Store has %d products\n", len(storeProducts.Products))
			for _, product := range storeProducts.Products {
				fmt.Printf("- %s (Price: %s)\n", product.Name, product.Price)
			}
		}

		// Get reviews for this store
		storeReviews, err := client.Stores.GetReviews(ctx, store.ID, &dokan.ReviewListParams{
			ListParams: dokan.ListParams{
				Page:    1,
				PerPage: 3,
			},
		})
		if err != nil {
			log.Printf("Failed to get store reviews: %v", err)
		} else {
			fmt.Printf("Store has %d reviews\n", len(storeReviews.Reviews))
			for _, review := range storeReviews.Reviews {
				fmt.Printf("- Rating: %d/5 - %s\n", review.Rating, review.Review)
			}
		}
	}

	// Example 3: Bulk product operations
	fmt.Println("\n=== Bulk Product Operations ===")
	products := []*dokan.Product{
		{
			Name:         "Bulk Product 1",
			Type:         dokan.ProductTypeSimple,
			RegularPrice: "19.99",
			Status:       dokan.ProductStatusDraft,
			SKU:          "BULK-001",
		},
		{
			Name:         "Bulk Product 2",
			Type:         dokan.ProductTypeSimple,
			RegularPrice: "29.99",
			Status:       dokan.ProductStatusDraft,
			SKU:          "BULK-002",
		},
		{
			Name:         "Bulk Product 3",
			Type:         dokan.ProductTypeSimple,
			RegularPrice: "39.99",
			Status:       dokan.ProductStatusDraft,
			SKU:          "BULK-003",
		},
	}

	var createdProducts []*dokan.Product
	for i, product := range products {
		fmt.Printf("Creating product %d/%d: %s\n", i+1, len(products), product.Name)
		
		created, err := client.Products.Create(ctx, product)
		if err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
			continue
		}
		
		createdProducts = append(createdProducts, created)
		fmt.Printf("✓ Created product ID: %d\n", created.ID)
		
		// Small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	// Example 4: Update products in bulk
	fmt.Println("\n=== Bulk Product Updates ===")
	for i, product := range createdProducts {
		fmt.Printf("Updating product %d/%d: %s\n", i+1, len(createdProducts), product.Name)
		
		// Update to published status and add description
		product.Status = dokan.ProductStatusPublish
		product.Description = fmt.Sprintf("This is product #%d created in bulk using the Dokan Go SDK", i+1)
		
		updated, err := client.Products.Update(ctx, product.ID, product)
		if err != nil {
			log.Printf("Failed to update product %d: %v", product.ID, err)
			continue
		}
		
		fmt.Printf("✓ Updated product ID: %d (Status: %s)\n", updated.ID, updated.Status)
		
		// Small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	// Example 5: Advanced filtering and search
	fmt.Println("\n=== Advanced Product Filtering ===")
	
	// Search for products with specific criteria
	searchParams := &dokan.ProductListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 20,
			Search:  "Bulk",
			OrderBy: "date",
			Order:   "desc",
		},
		Status:   []dokan.ProductStatus{dokan.ProductStatusPublish, dokan.ProductStatusDraft},
		Type:     []dokan.ProductType{dokan.ProductTypeSimple},
		MinPrice: &[]float64{10.0}[0],
		MaxPrice: &[]float64{50.0}[0],
	}

	searchResults, err := client.Products.List(ctx, searchParams)
	if err != nil {
		log.Printf("Failed to search products: %v", err)
	} else {
		fmt.Printf("Search found %d products matching criteria\n", len(searchResults.Products))
		for _, product := range searchResults.Products {
			fmt.Printf("- %s (ID: %d, Price: %s, Status: %s)\n", 
				product.Name, product.ID, product.Price, product.Status)
		}
	}

	// Example 6: Order management
	fmt.Println("\n=== Order Management ===")
	
	// Get recent orders with detailed filtering
	orderParams := &dokan.OrderListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 10,
			OrderBy: "date",
			Order:   "desc",
		},
		Status: []dokan.OrderStatus{
			dokan.OrderStatusProcessing,
			dokan.OrderStatusOnHold,
		},
		After: &[]time.Time{time.Now().AddDate(0, -1, 0)}[0], // Last month
	}

	orders, err := client.Orders.List(ctx, orderParams)
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
	} else {
		fmt.Printf("Found %d recent orders\n", len(orders.Orders))
		
		for _, order := range orders.Orders {
			fmt.Printf("Order #%s:\n", order.Number)
			fmt.Printf("  Status: %s\n", order.Status)
			fmt.Printf("  Total: %s %s\n", order.Total, order.Currency)
			fmt.Printf("  Customer: %s %s\n", order.Billing.FirstName, order.Billing.LastName)
			fmt.Printf("  Items: %d\n", len(order.LineItems))
			
			// Show line items
			for _, item := range order.LineItems {
				fmt.Printf("    - %s (Qty: %d, Total: %s)\n", 
					item.Name, item.Quantity, item.Total)
			}
			fmt.Println()
		}
	}

	// Example 7: Cleanup - Delete created products
	fmt.Println("\n=== Cleanup ===")
	for i, product := range createdProducts {
		fmt.Printf("Deleting product %d/%d: %s\n", i+1, len(createdProducts), product.Name)
		
		err := client.Products.Delete(ctx, product.ID)
		if err != nil {
			log.Printf("Failed to delete product %d: %v", product.ID, err)
			continue
		}
		
		fmt.Printf("✓ Deleted product ID: %d\n", product.ID)
		
		// Small delay to avoid rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n=== Advanced example completed ===")
}

// handleError demonstrates advanced error handling
func handleError(err error) {
	switch e := err.(type) {
	case *dokan.DokanError:
		fmt.Printf("Dokan API Error: %s (Code: %s, Status: %d)\n", 
			e.Message, e.Code, e.StatusCode)
		if e.Data != nil {
			fmt.Printf("Additional data: %+v\n", e.Data)
		}
	case *dokan.NetworkError:
		fmt.Printf("Network Error: %v\n", e.Err)
	case *dokan.AuthenticationError:
		fmt.Printf("Authentication Error: %s\n", e.Message)
	case *dokan.NotFoundError:
		fmt.Printf("Resource Not Found: %s with ID %v\n", e.Resource, e.ID)
	case *dokan.RateLimitError:
		fmt.Printf("Rate Limit Exceeded: retry after %d seconds\n", e.RetryAfter)
	case *dokan.ValidationError:
		fmt.Printf("Validation Error on field '%s': %s (Code: %s)\n", 
			e.Field, e.Message, e.Code)
	default:
		fmt.Printf("Unknown Error: %v\n", err)
	}
}

