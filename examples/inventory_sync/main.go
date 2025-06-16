package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk"
)

// InventoryItem representa un item del inventario externo
type InventoryItem struct {
	SKU          string
	Name         string
	Price        string
	Stock        int
	Description  string
	CategoryID   int
	ImageURL     string
	Featured     bool
}

func main() {
	// Configurar cliente
	client, err := dokan.NewClientBuilder().
		BaseURL(os.Getenv("DOKAN_BASE_URL")).
		BasicAuth(os.Getenv("DOKAN_USERNAME"), os.Getenv("DOKAN_PASSWORD")).
		Timeout(60 * time.Second).
		RetryCount(3).
		Build()
	if err != nil {
		log.Fatalf("Error creando cliente: %v", err)
	}

	ctx := context.Background()

	// Leer inventario desde archivo CSV
	inventoryItems, err := readInventoryFromCSV("inventory.csv")
	if err != nil {
		log.Fatalf("Error leyendo inventario: %v", err)
	}

	fmt.Printf("Sincronizando %d productos del inventario...\n", len(inventoryItems))

	// Obtener productos existentes de Dokan
	existingProducts, err := getAllProducts(client, ctx)
	if err != nil {
		log.Fatalf("Error obteniendo productos existentes: %v", err)
	}

	// Crear mapa de productos existentes por SKU
	existingBySKU := make(map[string]*dokan.Product)
	for _, product := range existingProducts {
		if product.SKU != "" {
			existingBySKU[product.SKU] = &product
		}
	}

	var created, updated, skipped int

	// Sincronizar cada item del inventario
	for i, item := range inventoryItems {
		fmt.Printf("Procesando item %d/%d: %s\n", i+1, len(inventoryItems), item.SKU)

		if existingProduct, exists := existingBySKU[item.SKU]; exists {
			// Producto existe, verificar si necesita actualización
			if needsUpdate(existingProduct, item) {
				if err := updateProduct(client, ctx, existingProduct, item); err != nil {
					log.Printf("Error actualizando producto %s: %v", item.SKU, err)
					continue
				}
				updated++
				fmt.Printf("✓ Actualizado: %s\n", item.SKU)
			} else {
				skipped++
				fmt.Printf("- Sin cambios: %s\n", item.SKU)
			}
		} else {
			// Producto no existe, crear nuevo
			if err := createProduct(client, ctx, item); err != nil {
				log.Printf("Error creando producto %s: %v", item.SKU, err)
				continue
			}
			created++
			fmt.Printf("✓ Creado: %s\n", item.SKU)
		}

		// Pausa para evitar rate limiting
		time.Sleep(500 * time.Millisecond)
	}

	// Resumen de la sincronización
	fmt.Printf("\n=== Resumen de Sincronización ===\n")
	fmt.Printf("Productos creados: %d\n", created)
	fmt.Printf("Productos actualizados: %d\n", updated)
	fmt.Printf("Productos sin cambios: %d\n", skipped)
	fmt.Printf("Total procesados: %d\n", created+updated+skipped)

	// Generar reporte de sincronización
	if err := generateSyncReport(created, updated, skipped, inventoryItems); err != nil {
		log.Printf("Error generando reporte: %v", err)
	}
}

// readInventoryFromCSV lee el inventario desde un archivo CSV
func readInventoryFromCSV(filename string) ([]InventoryItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("archivo CSV debe tener al menos una fila de datos")
	}

	var items []InventoryItem
	// Saltar header (primera fila)
	for i, record := range records[1:] {
		if len(record) < 8 {
			log.Printf("Fila %d: datos insuficientes, saltando", i+2)
			continue
		}

		stock, _ := strconv.Atoi(record[3])
		categoryID, _ := strconv.Atoi(record[5])
		featured, _ := strconv.ParseBool(record[7])

		item := InventoryItem{
			SKU:         record[0],
			Name:        record[1],
			Price:       record[2],
			Stock:       stock,
			Description: record[4],
			CategoryID:  categoryID,
			ImageURL:    record[6],
			Featured:    featured,
		}

		items = append(items, item)
	}

	return items, nil
}

// getAllProducts obtiene todos los productos de Dokan
func getAllProducts(client *dokan.Client, ctx context.Context) ([]dokan.Product, error) {
	var allProducts []dokan.Product
	page := 1
	perPage := 100

	for {
		params := &dokan.ProductListParams{
			ListParams: dokan.ListParams{
				Page:    page,
				PerPage: perPage,
			},
		}

		result, err := client.Products.List(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo productos página %d: %w", page, err)
		}

		allProducts = append(allProducts, result.Products...)

		if page >= result.TotalPages {
			break
		}

		page++
		time.Sleep(200 * time.Millisecond) // Pausa entre páginas
	}

	return allProducts, nil
}

// needsUpdate verifica si un producto necesita ser actualizado
func needsUpdate(existing *dokan.Product, item InventoryItem) bool {
	if existing.Name != item.Name {
		return true
	}
	if existing.RegularPrice != item.Price {
		return true
	}
	if existing.Description != item.Description {
		return true
	}
	if existing.Featured != item.Featured {
		return true
	}
	return false
}

// createProduct crea un nuevo producto en Dokan
func createProduct(client *dokan.Client, ctx context.Context, item InventoryItem) error {
	product := &dokan.Product{
		Name:              item.Name,
		Type:              dokan.ProductTypeSimple,
		RegularPrice:      item.Price,
		Description:       item.Description,
		ShortDescription:  fmt.Sprintf("Producto %s - Stock: %d", item.Name, item.Stock),
		Status:            dokan.ProductStatusPublish,
		Featured:          item.Featured,
		CatalogVisibility: dokan.CatalogVisibilityVisible,
		SKU:               item.SKU,
	}

	// Agregar categoría si está especificada
	if item.CategoryID > 0 {
		product.Categories = []dokan.ProductCategory{
			{ID: item.CategoryID},
		}
	}

	// Agregar imagen si está especificada
	if item.ImageURL != "" {
		product.Images = []dokan.ProductImage{
			{
				Src:      item.ImageURL,
				Alt:      item.Name,
				Position: 0,
			},
		}
	}

	_, err := client.Products.Create(ctx, product)
	return err
}

// updateProduct actualiza un producto existente
func updateProduct(client *dokan.Client, ctx context.Context, existing *dokan.Product, item InventoryItem) error {
	// Actualizar campos que han cambiado
	existing.Name = item.Name
	existing.RegularPrice = item.Price
	existing.Description = item.Description
	existing.Featured = item.Featured

	// Actualizar imagen si es diferente
	if item.ImageURL != "" {
		hasImage := false
		for _, img := range existing.Images {
			if img.Src == item.ImageURL {
				hasImage = true
				break
			}
		}

		if !hasImage {
			existing.Images = []dokan.ProductImage{
				{
					Src:      item.ImageURL,
					Alt:      item.Name,
					Position: 0,
				},
			}
		}
	}

	_, err := client.Products.Update(ctx, existing.ID, existing)
	return err
}

// generateSyncReport genera un reporte de la sincronización
func generateSyncReport(created, updated, skipped int, items []InventoryItem) error {
	filename := fmt.Sprintf("sync_report_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de reporte: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Reporte de Sincronización de Inventario\n")
	fmt.Fprintf(file, "Fecha: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "========================================\n\n")
	
	fmt.Fprintf(file, "Resumen:\n")
	fmt.Fprintf(file, "- Productos creados: %d\n", created)
	fmt.Fprintf(file, "- Productos actualizados: %d\n", updated)
	fmt.Fprintf(file, "- Productos sin cambios: %d\n", skipped)
	fmt.Fprintf(file, "- Total procesados: %d\n", len(items))
	fmt.Fprintf(file, "\nDetalles de productos procesados:\n")
	fmt.Fprintf(file, "==================================\n")
	
	for _, item := range items {
		fmt.Fprintf(file, "SKU: %s | Nombre: %s | Precio: %s\n", 
			item.SKU, item.Name, item.Price)
	}

	fmt.Printf("Reporte generado: %s\n", filename)
	return nil
}

