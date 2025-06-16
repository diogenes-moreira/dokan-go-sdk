# Documentación Técnica del SDK de Dokan para Go

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Arquitectura del SDK](#arquitectura-del-sdk)
3. [Guía de Instalación](#guía-de-instalación)
4. [Configuración Avanzada](#configuración-avanzada)
5. [API Reference](#api-reference)
6. [Ejemplos de Uso](#ejemplos-de-uso)
7. [Manejo de Errores](#manejo-de-errores)
8. [Testing](#testing)
9. [Mejores Prácticas](#mejores-prácticas)
10. [Troubleshooting](#troubleshooting)

## Introducción

El SDK de Dokan para Go es una biblioteca completa que proporciona una interfaz idiomática para interactuar con la API REST de Dokan Multivendor Marketplace. Este SDK está diseñado siguiendo las mejores prácticas de Go y proporciona type safety, manejo robusto de errores y una experiencia de desarrollo excelente.

### Características Principales

- **Type Safety**: Todos los tipos de datos están fuertemente tipados
- **Context Support**: Soporte completo para context.Context para cancelación y timeouts
- **Retry Logic**: Reintentos automáticos con backoff exponencial
- **Rate Limiting**: Protección integrada contra límites de tasa
- **Múltiples Métodos de Autenticación**: HTTP Basic Auth y JWT
- **Paginación Automática**: Manejo transparente de respuestas paginadas
- **Error Handling**: Sistema robusto de manejo de errores con tipos específicos

## Arquitectura del SDK

### Estructura de Paquetes

```
dokan-go-sdk/
├── auth/            # Sistema de autenticación
├── client/          # Cliente principal
├── errors/          # Manejo de errores
├── products/        # Servicio de productos
├── orders/          # Servicio de órdenes
├── stores/          # Servicio de tiendas
├── types/           # Tipos de datos
├── utils/           # Utilidades comunes
└── examples/        # Ejemplos de uso
```

### Componentes Principales

#### Cliente Principal

El cliente principal (`client.Client`) actúa como punto de entrada para todas las operaciones del SDK. Maneja la configuración global, autenticación y coordinación entre servicios.

#### Servicios

Cada endpoint principal de la API de Dokan está implementado como un servicio independiente:

- **ProductsService**: Gestión de productos
- **OrdersService**: Gestión de órdenes
- **StoresService**: Gestión de tiendas

#### Sistema de Autenticación

El sistema de autenticación es modular y extensible, soportando múltiples métodos:

- **BasicAuth**: Autenticación HTTP básica
- **JWTAuth**: Autenticación con JSON Web Tokens

## Guía de Instalación

### Requisitos

- Go 1.21 o superior
- Acceso a una instalación de Dokan con API REST habilitada

### Instalación

```bash
go get github.com/diogenes-moreira/dokan-go-sdk
```

### Verificación de Instalación

```go
package main

import (
    "fmt"
    "github.com/diogenes-moreira/dokan-go-sdk"
)

func main() {
    client, err := dokan.NewClientBuilder().
        BaseURL("https://tu-sitio.com").
        BasicAuth("usuario", "contraseña").
        Build()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Cliente creado exitosamente: %s\n", client.GetBaseURL())
}
```

## Configuración Avanzada

### Configuración del Cliente HTTP

```go
import (
    "net/http"
    "time"
)

// Cliente HTTP personalizado con configuración avanzada
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    TLSHandshakeTimeout: 10 * time.Second,
}

httpClient := &http.Client{
    Transport: transport,
    Timeout:   60 * time.Second,
}

client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    HTTPClient(httpClient).
    Build()
```

### Configuración de Reintentos

```go
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    RetryCount(5).                    // Máximo 5 reintentos
    Timeout(120 * time.Second).       // Timeout de 2 minutos
    Build()
```

### Configuración con Variables de Entorno

```go
import "os"

baseURL := os.Getenv("DOKAN_BASE_URL")
username := os.Getenv("DOKAN_USERNAME")
password := os.Getenv("DOKAN_PASSWORD")

client, err := dokan.NewClientBuilder().
    BaseURL(baseURL).
    BasicAuth(username, password).
    Build()
```

## API Reference

### Cliente Principal

#### NewClientBuilder()

Crea un nuevo builder para configurar el cliente.

```go
func NewClientBuilder() *ClientBuilder
```

#### ClientBuilder Methods

```go
func (b *ClientBuilder) BaseURL(url string) *ClientBuilder
func (b *ClientBuilder) BasicAuth(username, password string) *ClientBuilder
func (b *ClientBuilder) JWTAuth(token string) *ClientBuilder
func (b *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder
func (b *ClientBuilder) RetryCount(count int) *ClientBuilder
func (b *ClientBuilder) HTTPClient(client *http.Client) *ClientBuilder
func (b *ClientBuilder) Build() (*Client, error)
```

### Servicio de Productos

#### Create

Crea un nuevo producto.

```go
func (s *ProductsService) Create(ctx context.Context, product *Product) (*Product, error)
```

**Ejemplo:**

```go
product := &dokan.Product{
    Name:         "Producto Ejemplo",
    Type:         dokan.ProductTypeSimple,
    RegularPrice: "29.99",
    Status:       dokan.ProductStatusPublish,
}

created, err := client.Products.Create(ctx, product)
```

#### Get

Obtiene un producto por ID.

```go
func (s *ProductsService) Get(ctx context.Context, id int) (*Product, error)
```

#### List

Lista productos con filtros opcionales.

```go
func (s *ProductsService) List(ctx context.Context, params *ProductListParams) (*ProductListResponse, error)
```

#### Update

Actualiza un producto existente.

```go
func (s *ProductsService) Update(ctx context.Context, id int, product *Product) (*Product, error)
```

#### Delete

Elimina un producto.

```go
func (s *ProductsService) Delete(ctx context.Context, id int) error
```

### Servicio de Órdenes

#### Get

Obtiene una orden por ID.

```go
func (s *OrdersService) Get(ctx context.Context, id int) (*Order, error)
```

#### List

Lista órdenes con filtros opcionales.

```go
func (s *OrdersService) List(ctx context.Context, params *OrderListParams) (*OrderListResponse, error)
```

#### Update

Actualiza una orden existente.

```go
func (s *OrdersService) Update(ctx context.Context, id int, order *OrderUpdate) (*Order, error)
```

### Servicio de Tiendas

#### Get

Obtiene información de una tienda por ID de vendedor.

```go
func (s *StoresService) Get(ctx context.Context, vendorID int) (*Store, error)
```

#### List

Lista tiendas con filtros opcionales.

```go
func (s *StoresService) List(ctx context.Context, params *StoreListParams) (*StoreListResponse, error)
```

#### GetProducts

Obtiene productos de una tienda específica.

```go
func (s *StoresService) GetProducts(ctx context.Context, vendorID int, params *ProductListParams) (*StoreProductsResponse, error)
```

#### GetReviews

Obtiene reseñas de una tienda específica.

```go
func (s *StoresService) GetReviews(ctx context.Context, vendorID int, params *ReviewListParams) (*StoreReviewsResponse, error)
```


## Ejemplos de Uso

### Gestión Básica de Productos

#### Crear un Producto Simple

```go
func createSimpleProduct(client *dokan.Client) error {
    ctx := context.Background()
    
    product := &dokan.Product{
        Name:              "Camiseta Premium",
        Type:              dokan.ProductTypeSimple,
        RegularPrice:      "49.99",
        SalePrice:         "39.99",
        Description:       "Camiseta de algodón 100% orgánico",
        ShortDescription:  "Camiseta premium de algodón orgánico",
        Status:            dokan.ProductStatusPublish,
        Featured:          true,
        CatalogVisibility: dokan.CatalogVisibilityVisible,
        SKU:               "SHIRT-PREMIUM-001",
        Categories: []dokan.ProductCategory{
            {ID: 15, Name: "Ropa"},
            {ID: 23, Name: "Camisetas"},
        },
        Images: []dokan.ProductImage{
            {
                Src: "https://ejemplo.com/camiseta-frente.jpg",
                Alt: "Camiseta vista frontal",
                Position: 0,
            },
            {
                Src: "https://ejemplo.com/camiseta-espalda.jpg",
                Alt: "Camiseta vista trasera",
                Position: 1,
            },
        },
    }
    
    created, err := client.Products.Create(ctx, product)
    if err != nil {
        return fmt.Errorf("error creando producto: %w", err)
    }
    
    fmt.Printf("Producto creado exitosamente:\n")
    fmt.Printf("ID: %d\n", created.ID)
    fmt.Printf("Nombre: %s\n", created.Name)
    fmt.Printf("URL: %s\n", created.Permalink)
    
    return nil
}
```

#### Buscar Productos con Filtros Avanzados

```go
func searchProducts(client *dokan.Client) error {
    ctx := context.Background()
    
    params := &dokan.ProductListParams{
        ListParams: dokan.ListParams{
            Page:    1,
            PerPage: 20,
            Search:  "camiseta",
            OrderBy: "popularity",
            Order:   "desc",
        },
        Status:      []dokan.ProductStatus{dokan.ProductStatusPublish},
        Type:        []dokan.ProductType{dokan.ProductTypeSimple},
        Featured:    &[]bool{true}[0],
        Category:    []int{15, 23}, // Ropa y Camisetas
        MinPrice:    &[]float64{20.0}[0],
        MaxPrice:    &[]float64{100.0}[0],
        StockStatus: "instock",
    }
    
    result, err := client.Products.List(ctx, params)
    if err != nil {
        return fmt.Errorf("error buscando productos: %w", err)
    }
    
    fmt.Printf("Encontrados %d productos (página %d de %d):\n", 
        len(result.Products), result.Page, result.TotalPages)
    
    for _, product := range result.Products {
        fmt.Printf("- %s (ID: %d)\n", product.Name, product.ID)
        fmt.Printf("  Precio: %s\n", product.Price)
        fmt.Printf("  Estado: %s\n", product.Status)
        fmt.Printf("  SKU: %s\n", product.SKU)
        fmt.Println()
    }
    
    return nil
}
```

### Gestión de Órdenes

#### Procesar Órdenes Pendientes

```go
func processOrders(client *dokan.Client) error {
    ctx := context.Background()
    
    // Buscar órdenes pendientes
    params := &dokan.OrderListParams{
        ListParams: dokan.ListParams{
            Page:    1,
            PerPage: 50,
            OrderBy: "date",
            Order:   "asc",
        },
        Status: []dokan.OrderStatus{dokan.OrderStatusPending},
    }
    
    orders, err := client.Orders.List(ctx, params)
    if err != nil {
        return fmt.Errorf("error obteniendo órdenes: %w", err)
    }
    
    fmt.Printf("Procesando %d órdenes pendientes...\n", len(orders.Orders))
    
    for _, order := range orders.Orders {
        fmt.Printf("Procesando orden #%s...\n", order.Number)
        
        // Validar la orden
        if err := validateOrder(order); err != nil {
            fmt.Printf("Error validando orden #%s: %v\n", order.Number, err)
            continue
        }
        
        // Actualizar estado a "processing"
        update := &orders.OrderUpdate{
            Status: &[]dokan.OrderStatus{dokan.OrderStatusProcessing}[0],
        }
        
        updated, err := client.Orders.Update(ctx, order.ID, update)
        if err != nil {
            fmt.Printf("Error actualizando orden #%s: %v\n", order.Number, err)
            continue
        }
        
        fmt.Printf("✓ Orden #%s actualizada a estado: %s\n", 
            updated.Number, updated.Status)
        
        // Pequeña pausa para evitar rate limiting
        time.Sleep(500 * time.Millisecond)
    }
    
    return nil
}

func validateOrder(order dokan.Order) error {
    if len(order.LineItems) == 0 {
        return fmt.Errorf("orden sin productos")
    }
    
    if order.Total == "" || order.Total == "0" {
        return fmt.Errorf("orden sin total válido")
    }
    
    if order.Billing == nil {
        return fmt.Errorf("orden sin información de facturación")
    }
    
    return nil
}
```

### Gestión de Tiendas

#### Analizar Rendimiento de Tiendas

```go
func analyzeStorePerformance(client *dokan.Client) error {
    ctx := context.Background()
    
    // Obtener todas las tiendas activas
    stores, err := client.Stores.List(ctx, &dokan.StoreListParams{
        ListParams: dokan.ListParams{
            Page:    1,
            PerPage: 100,
        },
        Enabled: &[]bool{true}[0],
    })
    if err != nil {
        return fmt.Errorf("error obteniendo tiendas: %w", err)
    }
    
    fmt.Printf("Analizando rendimiento de %d tiendas...\n\n", len(stores.Stores))
    
    for _, store := range stores.Stores {
        fmt.Printf("=== Tienda: %s (ID: %d) ===\n", store.StoreName, store.ID)
        
        // Obtener productos de la tienda
        products, err := client.Stores.GetProducts(ctx, store.ID, &dokan.ProductListParams{
            ListParams: dokan.ListParams{
                Page:    1,
                PerPage: 1000, // Obtener todos los productos
            },
        })
        if err != nil {
            fmt.Printf("Error obteniendo productos: %v\n", err)
            continue
        }
        
        // Analizar productos
        var publishedCount, draftCount, featuredCount int
        var totalValue float64
        
        for _, product := range products.Products {
            switch product.Status {
            case dokan.ProductStatusPublish:
                publishedCount++
            case dokan.ProductStatusDraft:
                draftCount++
            }
            
            if product.Featured {
                featuredCount++
            }
            
            if price, err := strconv.ParseFloat(product.RegularPrice, 64); err == nil {
                totalValue += price
            }
        }
        
        // Obtener reseñas
        reviews, err := client.Stores.GetReviews(ctx, store.ID, &dokan.ReviewListParams{
            ListParams: dokan.ListParams{
                Page:    1,
                PerPage: 1000,
            },
        })
        if err != nil {
            fmt.Printf("Error obteniendo reseñas: %v\n", err)
        }
        
        // Calcular rating promedio
        var totalRating, reviewCount int
        if reviews != nil {
            for _, review := range reviews.Reviews {
                totalRating += review.Rating
                reviewCount++
            }
        }
        
        avgRating := 0.0
        if reviewCount > 0 {
            avgRating = float64(totalRating) / float64(reviewCount)
        }
        
        // Mostrar estadísticas
        fmt.Printf("Productos totales: %d\n", len(products.Products))
        fmt.Printf("- Publicados: %d\n", publishedCount)
        fmt.Printf("- Borradores: %d\n", draftCount)
        fmt.Printf("- Destacados: %d\n", featuredCount)
        fmt.Printf("Valor total del inventario: $%.2f\n", totalValue)
        fmt.Printf("Reseñas: %d (Rating promedio: %.1f/5)\n", reviewCount, avgRating)
        fmt.Printf("Email: %s\n", store.Email)
        fmt.Println()
    }
    
    return nil
}
```

## Manejo de Errores

### Tipos de Error

El SDK define varios tipos de error específicos para diferentes situaciones:

#### DokanError

Error específico de la API de Dokan.

```go
type DokanError struct {
    Code       string      `json:"code"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    StatusCode int         `json:"-"`
}
```

#### NetworkError

Error de conectividad de red.

```go
type NetworkError struct {
    Err error
}
```

#### AuthenticationError

Error de autenticación.

```go
type AuthenticationError struct {
    Message string
}
```

#### ValidationError

Error de validación de datos.

```go
type ValidationError struct {
    Field   string `json:"field"`
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Manejo Robusto de Errores

```go
func handleAPICall(client *dokan.Client) {
    ctx := context.Background()
    
    product, err := client.Products.Get(ctx, 123)
    if err != nil {
        switch e := err.(type) {
        case *dokan.DokanError:
            switch e.StatusCode {
            case 400:
                log.Printf("Solicitud inválida: %s", e.Message)
            case 401:
                log.Printf("No autorizado: %s", e.Message)
                // Renovar token o solicitar nuevas credenciales
            case 404:
                log.Printf("Producto no encontrado: %s", e.Message)
            case 429:
                log.Printf("Rate limit excedido: %s", e.Message)
                // Implementar backoff
            case 500:
                log.Printf("Error del servidor: %s", e.Message)
                // Reintentar después de un tiempo
            default:
                log.Printf("Error de API: %s (código: %s)", e.Message, e.Code)
            }
            
        case *dokan.NetworkError:
            log.Printf("Error de red: %v", e.Err)
            // Verificar conectividad
            
        case *dokan.AuthenticationError:
            log.Printf("Error de autenticación: %s", e.Message)
            // Renovar credenciales
            
        case *dokan.NotFoundError:
            log.Printf("Recurso no encontrado: %s con ID %v", e.Resource, e.ID)
            
        case *dokan.RateLimitError:
            log.Printf("Rate limit excedido, esperar %d segundos", e.RetryAfter)
            time.Sleep(time.Duration(e.RetryAfter) * time.Second)
            
        default:
            log.Printf("Error desconocido: %v", err)
        }
        return
    }
    
    // Procesar producto exitosamente obtenido
    fmt.Printf("Producto obtenido: %s\n", product.Name)
}
```

### Retry con Backoff Personalizado

```go
func retryWithCustomBackoff(client *dokan.Client) error {
    ctx := context.Background()
    maxRetries := 5
    baseDelay := time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        product, err := client.Products.Get(ctx, 123)
        if err == nil {
            fmt.Printf("Éxito en intento %d: %s\n", attempt+1, product.Name)
            return nil
        }
        
        // Verificar si es un error que vale la pena reintentar
        if dokanErr, ok := err.(*dokan.DokanError); ok {
            if dokanErr.StatusCode >= 400 && dokanErr.StatusCode < 500 {
                // Error del cliente, no reintentar
                return err
            }
        }
        
        if attempt < maxRetries-1 {
            delay := time.Duration(attempt+1) * baseDelay
            fmt.Printf("Intento %d falló, reintentando en %v...\n", attempt+1, delay)
            time.Sleep(delay)
        }
    }
    
    return fmt.Errorf("falló después de %d intentos", maxRetries)
}
```


## Testing

### Unit Testing

El SDK está diseñado para ser fácilmente testeable usando interfaces mock.

#### Ejemplo de Test Unitario

```go
package main

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/diogenes-moreira/dokan-go-sdk"
    "github.com/diogenes-moreira/dokan-go-sdk/products"
    "github.com/diogenes-moreira/dokan-go-sdk/utils"
)

// MockClient implementa la interfaz ClientInterface para testing
type MockClient struct {
    mock.Mock
}

func (m *MockClient) MakeRequest(ctx context.Context, opts utils.RequestOptions) (*utils.Response, error) {
    args := m.Called(ctx, opts)
    return args.Get(0).(*utils.Response), args.Error(1)
}

func TestProductsService_Create(t *testing.T) {
    // Arrange
    mockClient := new(MockClient)
    service := products.NewService(mockClient)
    
    product := &dokan.Product{
        Name:         "Test Product",
        Type:         dokan.ProductTypeSimple,
        RegularPrice: "29.99",
        Status:       dokan.ProductStatusDraft,
    }
    
    expectedResponse := &utils.Response{
        StatusCode: 201,
        Body:       []byte(`{"id": 123, "name": "Test Product", "type": "simple", "regular_price": "29.99", "status": "draft"}`),
    }
    
    mockClient.On("MakeRequest", mock.Anything, mock.MatchedBy(func(opts utils.RequestOptions) bool {
        return opts.Method == "POST" && opts.Path == "/wp-json/dokan/v1/products/"
    })).Return(expectedResponse, nil)
    
    // Act
    result, err := service.Create(context.Background(), product)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, 123, result.ID)
    assert.Equal(t, "Test Product", result.Name)
    
    mockClient.AssertExpectations(t)
}

func TestProductsService_Get_NotFound(t *testing.T) {
    // Arrange
    mockClient := new(MockClient)
    service := products.NewService(mockClient)
    
    expectedError := &dokan.DokanError{
        Code:       "product_invalid_id",
        Message:    "Invalid product ID",
        StatusCode: 404,
    }
    
    mockClient.On("MakeRequest", mock.Anything, mock.MatchedBy(func(opts utils.RequestOptions) bool {
        return opts.Method == "GET" && opts.Path == "/wp-json/dokan/v1/products/999"
    })).Return(nil, expectedError)
    
    // Act
    result, err := service.Get(context.Background(), 999)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.IsType(t, &dokan.DokanError{}, err)
    
    dokanErr := err.(*dokan.DokanError)
    assert.Equal(t, "product_invalid_id", dokanErr.Code)
    assert.Equal(t, 404, dokanErr.StatusCode)
    
    mockClient.AssertExpectations(t)
}
```

### Integration Testing

Para tests de integración, puedes usar una instancia real de Dokan en un entorno de testing.

#### Configuración de Test de Integración

```go
package integration

import (
    "context"
    "os"
    "testing"
    "time"
    "github.com/diogenes-moreira/dokan-go-sdk"
)

func setupTestClient(t *testing.T) *dokan.Client {
    baseURL := os.Getenv("DOKAN_TEST_BASE_URL")
    username := os.Getenv("DOKAN_TEST_USERNAME")
    password := os.Getenv("DOKAN_TEST_PASSWORD")
    
    if baseURL == "" || username == "" || password == "" {
        t.Skip("Variables de entorno de test no configuradas")
    }
    
    client, err := dokan.NewClientBuilder().
        BaseURL(baseURL).
        BasicAuth(username, password).
        Timeout(30 * time.Second).
        Build()
    
    if err != nil {
        t.Fatalf("Error creando cliente de test: %v", err)
    }
    
    return client
}

func TestIntegration_ProductLifecycle(t *testing.T) {
    if testing.Short() {
        t.Skip("Saltando test de integración en modo short")
    }
    
    client := setupTestClient(t)
    ctx := context.Background()
    
    // Crear producto
    product := &dokan.Product{
        Name:         "Test Integration Product",
        Type:         dokan.ProductTypeSimple,
        RegularPrice: "19.99",
        Status:       dokan.ProductStatusDraft,
        SKU:          "TEST-INTEGRATION-001",
    }
    
    created, err := client.Products.Create(ctx, product)
    if err != nil {
        t.Fatalf("Error creando producto: %v", err)
    }
    
    defer func() {
        // Cleanup: eliminar producto al final del test
        if err := client.Products.Delete(ctx, created.ID); err != nil {
            t.Logf("Error eliminando producto de test: %v", err)
        }
    }()
    
    // Verificar que se creó correctamente
    assert.NotZero(t, created.ID)
    assert.Equal(t, product.Name, created.Name)
    assert.Equal(t, product.SKU, created.SKU)
    
    // Obtener producto
    retrieved, err := client.Products.Get(ctx, created.ID)
    if err != nil {
        t.Fatalf("Error obteniendo producto: %v", err)
    }
    
    assert.Equal(t, created.ID, retrieved.ID)
    assert.Equal(t, created.Name, retrieved.Name)
    
    // Actualizar producto
    retrieved.Description = "Descripción actualizada por test de integración"
    retrieved.Status = dokan.ProductStatusPublish
    
    updated, err := client.Products.Update(ctx, retrieved.ID, retrieved)
    if err != nil {
        t.Fatalf("Error actualizando producto: %v", err)
    }
    
    assert.Equal(t, "Descripción actualizada por test de integración", updated.Description)
    assert.Equal(t, dokan.ProductStatusPublish, updated.Status)
}
```

## Mejores Prácticas

### 1. Manejo de Context

Siempre usa context.Context para operaciones que pueden ser canceladas o que tienen timeout.

```go
// ✅ Correcto
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

product, err := client.Products.Get(ctx, 123)

// ❌ Incorrecto
product, err := client.Products.Get(context.Background(), 123)
```

### 2. Manejo de Errores

Siempre verifica y maneja los errores apropiadamente.

```go
// ✅ Correcto
product, err := client.Products.Get(ctx, 123)
if err != nil {
    if dokanErr, ok := err.(*dokan.DokanError); ok && dokanErr.StatusCode == 404 {
        log.Printf("Producto no encontrado: %d", 123)
        return nil
    }
    return fmt.Errorf("error obteniendo producto: %w", err)
}

// ❌ Incorrecto
product, _ := client.Products.Get(ctx, 123)
```

### 3. Reutilización de Cliente

Reutiliza la instancia del cliente en lugar de crear una nueva para cada operación.

```go
// ✅ Correcto - Crear una vez, usar muchas veces
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    Build()

// Usar el mismo cliente para múltiples operaciones
product1, _ := client.Products.Get(ctx, 1)
product2, _ := client.Products.Get(ctx, 2)
orders, _ := client.Orders.List(ctx, nil)

// ❌ Incorrecto - Crear cliente para cada operación
client1, _ := dokan.NewClientBuilder().BaseURL("...").Build()
product1, _ := client1.Products.Get(ctx, 1)

client2, _ := dokan.NewClientBuilder().BaseURL("...").Build()
product2, _ := client2.Products.Get(ctx, 2)
```

### 4. Paginación

Para listados grandes, usa paginación para evitar timeouts y problemas de memoria.

```go
// ✅ Correcto - Procesar en páginas
func getAllProducts(client *dokan.Client) ([]dokan.Product, error) {
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
            return nil, err
        }
        
        allProducts = append(allProducts, result.Products...)
        
        if page >= result.TotalPages {
            break
        }
        
        page++
    }
    
    return allProducts, nil
}
```

### 5. Rate Limiting

Respeta los límites de tasa de la API.

```go
// ✅ Correcto - Agregar delays entre operaciones masivas
func bulkUpdateProducts(client *dokan.Client, products []dokan.Product) error {
    for i, product := range products {
        _, err := client.Products.Update(ctx, product.ID, &product)
        if err != nil {
            return err
        }
        
        // Pausa entre operaciones para evitar rate limiting
        if i < len(products)-1 {
            time.Sleep(500 * time.Millisecond)
        }
    }
    return nil
}
```

### 6. Configuración Segura

No hardcodees credenciales en el código.

```go
// ✅ Correcto - Usar variables de entorno
baseURL := os.Getenv("DOKAN_BASE_URL")
username := os.Getenv("DOKAN_USERNAME")
password := os.Getenv("DOKAN_PASSWORD")

client, err := dokan.NewClientBuilder().
    BaseURL(baseURL).
    BasicAuth(username, password).
    Build()

// ❌ Incorrecto - Credenciales hardcodeadas
client, err := dokan.NewClientBuilder().
    BaseURL("https://mi-sitio.com").
    BasicAuth("admin", "password123").
    Build()
```

## Troubleshooting

### Problemas Comunes

#### 1. Error de Autenticación

**Síntoma:** `authentication error: unauthorized access`

**Soluciones:**
- Verificar que las credenciales sean correctas
- Asegurar que el usuario tenga permisos para acceder a la API
- Verificar que la API REST esté habilitada en WordPress
- Comprobar que el plugin Dokan esté activo

```go
// Verificar credenciales
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    Build()

if err != nil {
    log.Printf("Error creando cliente: %v", err)
    return
}

// Test simple para verificar autenticación
_, err = client.Products.List(context.Background(), &dokan.ProductListParams{
    ListParams: dokan.ListParams{Page: 1, PerPage: 1},
})

if err != nil {
    if authErr, ok := err.(*dokan.AuthenticationError); ok {
        log.Printf("Error de autenticación: %s", authErr.Message)
        // Verificar credenciales
    }
}
```

#### 2. Timeouts

**Síntoma:** `context deadline exceeded` o `network error: timeout`

**Soluciones:**
- Aumentar el timeout del cliente
- Verificar la conectividad de red
- Reducir el tamaño de las solicitudes (usar paginación)

```go
// Aumentar timeout
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    Timeout(120 * time.Second).  // Aumentar a 2 minutos
    Build()

// Usar context con timeout personalizado
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

products, err := client.Products.List(ctx, params)
```

#### 3. Rate Limiting

**Síntoma:** `rate limit exceeded, retry after X seconds`

**Soluciones:**
- Implementar delays entre solicitudes
- Usar el campo RetryAfter del error
- Reducir la frecuencia de solicitudes

```go
func handleRateLimit(client *dokan.Client, productID int) (*dokan.Product, error) {
    for {
        product, err := client.Products.Get(ctx, productID)
        if err == nil {
            return product, nil
        }
        
        if rateLimitErr, ok := err.(*dokan.RateLimitError); ok {
            log.Printf("Rate limit excedido, esperando %d segundos...", rateLimitErr.RetryAfter)
            time.Sleep(time.Duration(rateLimitErr.RetryAfter) * time.Second)
            continue
        }
        
        return nil, err
    }
}
```

#### 4. Errores de Validación

**Síntoma:** `validation error on field 'X': message`

**Soluciones:**
- Verificar que todos los campos requeridos estén presentes
- Validar el formato de los datos
- Consultar la documentación de la API de Dokan

```go
func validateProduct(product *dokan.Product) error {
    if product.Name == "" {
        return fmt.Errorf("el nombre del producto es requerido")
    }
    
    if product.Type == "" {
        product.Type = dokan.ProductTypeSimple // Valor por defecto
    }
    
    if product.RegularPrice == "" {
        return fmt.Errorf("el precio regular es requerido")
    }
    
    // Validar formato de precio
    if _, err := strconv.ParseFloat(product.RegularPrice, 64); err != nil {
        return fmt.Errorf("formato de precio inválido: %s", product.RegularPrice)
    }
    
    return nil
}
```

### Debugging

#### Habilitar Modo Debug

```go
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    Debug(true).  // Habilitar logging detallado
    Build()
```

#### Logging Personalizado

```go
import "log"

func logAPICall(method, path string, err error) {
    if err != nil {
        log.Printf("API Call Failed: %s %s - Error: %v", method, path, err)
    } else {
        log.Printf("API Call Success: %s %s", method, path)
    }
}

// Usar en tu código
product, err := client.Products.Get(ctx, 123)
logAPICall("GET", "/products/123", err)
```

### Recursos Adicionales

- [Documentación oficial de Dokan API](https://getdokan.github.io/dokan/)
- [WordPress REST API Handbook](https://developer.wordpress.org/rest-api/)
- [Go Context Package](https://golang.org/pkg/context/)
- [Go HTTP Client Best Practices](https://golang.org/pkg/net/http/)

---

Esta documentación técnica proporciona una guía completa para usar el SDK de Dokan para Go de manera efectiva y siguiendo las mejores prácticas. Para ejemplos adicionales y casos de uso específicos, consulta la carpeta `examples/` del repositorio.

