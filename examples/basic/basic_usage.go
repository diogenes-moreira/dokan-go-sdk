package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk"
)

func main() {
	// Create a new Dokan client with Basic Authentication
	client, err := dokan.NewClientBuilder().
		BaseURL("https://your-dokan-site.com").
		BasicAuth("your-username", "your-password").
		Timeout(30 * time.Second).
		RetryCount(3).
		Build()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Create a new product
	fmt.Println("=== Creating a new product ===")
	product := &dokan.Product{
		Name:         "Amazing Go SDK Product",
		Type:         dokan.ProductTypeSimple,
		RegularPrice: "99.99",
		SalePrice:    "79.99",
		Description:  "This product was created using the Dokan Go SDK!",
		ShortDescription: "Created with Go SDK",
		Status:       dokan.ProductStatusPublish,
		Featured:     true,
		CatalogVisibility: dokan.CatalogVisibilityVisible,
		SKU:          "GO-SDK-001",
		Categories: []dokan.ProductCategory{
			{ID: 1, Name: "Electronics"},
		},
		Images: []dokan.ProductImage{
			{
				Src: "https://example.com/image.jpg",
				Alt: "Product image",
			},
		},
	}

	createdProduct, err := client.Products.Create(ctx, product)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
	} else {
		fmt.Printf("Created product with ID: %d\n", createdProduct.ID)
		fmt.Printf("Product URL: %s\n", createdProduct.Permalink)
	}

	// Example 2: List products with filtering
	fmt.Println("\n=== Listing products ===")
	listParams := &dokan.ProductListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 10,
			OrderBy: "date",
			Order:   "desc",
		},
		Status:   []dokan.ProductStatus{dokan.ProductStatusPublish},
		Featured: &[]bool{true}[0], // Pointer to true
	}

	productList, err := client.Products.List(ctx, listParams)
	if err != nil {
		log.Printf("Failed to list products: %v", err)
	} else {
		fmt.Printf("Found %d products (page %d of %d)\n", 
			len(productList.Products), 
			productList.Page, 
			productList.TotalPages)
		
		for _, p := range productList.Products {
			fmt.Printf("- %s (ID: %d, Price: %s)\n", p.Name, p.ID, p.Price)
		}
	}

	// Example 3: Get a specific product
	if createdProduct != nil {
		fmt.Println("\n=== Getting specific product ===")
		retrievedProduct, err := client.Products.Get(ctx, createdProduct.ID)
		if err != nil {
			log.Printf("Failed to get product: %v", err)
		} else {
			fmt.Printf("Retrieved product: %s\n", retrievedProduct.Name)
			fmt.Printf("Status: %s, Featured: %t\n", retrievedProduct.Status, retrievedProduct.Featured)
		}
	}

	// Example 4: List stores
	fmt.Println("\n=== Listing stores ===")
	storeParams := &dokan.StoreListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 5,
		},
		Enabled: &[]bool{true}[0], // Pointer to true
	}

	storeList, err := client.Stores.List(ctx, storeParams)
	if err != nil {
		log.Printf("Failed to list stores: %v", err)
	} else {
		fmt.Printf("Found %d stores\n", len(storeList.Stores))
		
		for _, store := range storeList.Stores {
			fmt.Printf("- %s (ID: %d, Email: %s)\n", 
				store.StoreName, 
				store.ID, 
				store.Email)
		}
	}

	// Example 5: Get orders
	fmt.Println("\n=== Listing recent orders ===")
	orderParams := &dokan.OrderListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: 5,
			OrderBy: "date",
			Order:   "desc",
		},
		Status: []dokan.OrderStatus{
			dokan.OrderStatusProcessing,
			dokan.OrderStatusCompleted,
		},
	}

	orderList, err := client.Orders.List(ctx, orderParams)
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
	} else {
		fmt.Printf("Found %d orders\n", len(orderList.Orders))
		
		for _, order := range orderList.Orders {
			fmt.Printf("- Order #%s (Status: %s, Total: %s)\n", 
				order.Number, 
				order.Status, 
				order.Total)
		}
	}

	// Example 6: Update a product (if we created one)
	if createdProduct != nil {
		fmt.Println("\n=== Updating product ===")
		createdProduct.Description = "Updated description using the Dokan Go SDK!"
		createdProduct.RegularPrice = "109.99"
		
		updatedProduct, err := client.Products.Update(ctx, createdProduct.ID, createdProduct)
		if err != nil {
			log.Printf("Failed to update product: %v", err)
		} else {
			fmt.Printf("Updated product price to: %s\n", updatedProduct.RegularPrice)
		}
	}

	// Example 7: Get product summary
	fmt.Println("\n=== Getting product summary ===")
	summary, err := client.Products.GetSummary(ctx)
	if err != nil {
		log.Printf("Failed to get product summary: %v", err)
	} else {
		fmt.Printf("Product Summary:\n")
		fmt.Printf("- Total: %d\n", summary.Total)
		fmt.Printf("- Published: %d\n", summary.Published)
		fmt.Printf("- Draft: %d\n", summary.Draft)
		fmt.Printf("- Pending: %d\n", summary.Pending)
		fmt.Printf("- Featured: %d\n", summary.Featured)
	}

	fmt.Println("\n=== Example completed ===")
}

