# Dokan Go SDK

Un SDK completo en Go para la API REST de Dokan Multivendor Marketplace.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Descripción

Dokan es un plugin de WordPress que permite crear marketplaces multivendor basados en WooCommerce. Este SDK proporciona una interfaz idiomática en Go para interactuar con todas las funcionalidades de la API REST de Dokan.

## Características

- ✅ **Cobertura completa de la API**: Soporte para productos, órdenes, tiendas, cupones, retiros, reseñas y más
- ✅ **Autenticación múltiple**: HTTP Basic Auth y JWT
- ✅ **Manejo robusto de errores**: Tipos de error específicos y manejo de reintentos
- ✅ **Paginación automática**: Soporte completo para listados paginados
- ✅ **Context support**: Cancelación y timeouts usando context.Context
- ✅ **Type safety**: Tipos fuertemente tipados para todas las estructuras de datos
- ✅ **Retry logic**: Reintentos automáticos con backoff exponencial
- ✅ **Rate limiting**: Protección contra límites de tasa
- ✅ **Documentación completa**: GoDoc y ejemplos de uso

## Instalación

```bash
go get github.com/diogenes-moreira/dokan-go-sdk
```

## Uso Básico

### Configuración del Cliente

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/diogenes-moreira/dokan-go-sdk"
)

func main() {
    // Crear cliente con autenticación básica
    client, err := dokan.NewClientBuilder().
        BaseURL("https://tu-sitio-dokan.com").
        BasicAuth("usuario", "contraseña").
        Timeout(30 * time.Second).
        RetryCount(3).
        Build()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Tu código aquí...
}
```

### Crear un Producto

```go
product := &dokan.Product{
    Name:         "Producto Ejemplo",
    Type:         dokan.ProductTypeSimple,
    RegularPrice: "29.99",
    SalePrice:    "24.99",
    Description:  "Descripción del producto",
    Status:       dokan.ProductStatusPublish,
    SKU:          "PROD-001",
}

createdProduct, err := client.Products.Create(ctx, product)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Producto creado con ID: %d\n", createdProduct.ID)
```

### Listar Productos

```go
params := &dokan.ProductListParams{
    ListParams: dokan.ListParams{
        Page:    1,
        PerPage: 10,
        OrderBy: "date",
        Order:   "desc",
    },
    Status:   []dokan.ProductStatus{dokan.ProductStatusPublish},
    Featured: &[]bool{true}[0], // Solo productos destacados
}

products, err := client.Products.List(ctx, params)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Encontrados %d productos\n", len(products.Products))
for _, p := range products.Products {
    fmt.Printf("- %s (ID: %d, Precio: %s)\n", p.Name, p.ID, p.Price)
}
```

### Gestión de Tiendas

```go
// Listar tiendas
stores, err := client.Stores.List(ctx, &dokan.StoreListParams{
    ListParams: dokan.ListParams{Page: 1, PerPage: 10},
    Enabled:    &[]bool{true}[0],
})
if err != nil {
    log.Fatal(err)
}

// Obtener productos de una tienda específica
if len(stores.Stores) > 0 {
    storeID := stores.Stores[0].ID
    storeProducts, err := client.Stores.GetProducts(ctx, storeID, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("La tienda tiene %d productos\n", len(storeProducts.Products))
}
```

### Gestión de Órdenes

```go
// Listar órdenes recientes
orders, err := client.Orders.List(ctx, &dokan.OrderListParams{
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
})
if err != nil {
    log.Fatal(err)
}

for _, order := range orders.Orders {
    fmt.Printf("Orden #%s - Estado: %s - Total: %s\n", 
        order.Number, order.Status, order.Total)
}
```

## Autenticación

### HTTP Basic Auth

```go
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    BasicAuth("usuario", "contraseña").
    Build()
```

### JWT Authentication

```go
client, err := dokan.NewClientBuilder().
    BaseURL("https://tu-sitio.com").
    JWTAuth("tu-jwt-token").
    Build()
```

### Autenticación Personalizada

```go
// Crear autenticador personalizado
auth := dokan.NewJWTAuth("token", time.Now().Add(time.Hour))

client, err := dokan.NewClient(&dokan.Config{
    BaseURL: "https://tu-sitio.com",
    Auth: dokan.AuthConfig{
        Type:  dokan.AuthTypeJWT,
        Token: "tu-token",
    },
})
```

## Manejo de Errores

El SDK proporciona tipos de error específicos para diferentes situaciones:

```go
product, err := client.Products.Get(ctx, 123)
if err != nil {
    switch e := err.(type) {
    case *dokan.DokanError:
        fmt.Printf("Error de API: %s (Código: %s)\n", e.Message, e.Code)
    case *dokan.NetworkError:
        fmt.Printf("Error de red: %v\n", e.Err)
    case *dokan.AuthenticationError:
        fmt.Printf("Error de autenticación: %s\n", e.Message)
    case *dokan.NotFoundError:
        fmt.Printf("Recurso no encontrado: %s\n", e.Resource)
    case *dokan.RateLimitError:
        fmt.Printf("Límite de tasa excedido, reintentar en %d segundos\n", e.RetryAfter)
    default:
        fmt.Printf("Error desconocido: %v\n", err)
    }
}
```

## Configuración Avanzada

### Cliente HTTP Personalizado

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
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
    RetryCount(5).
    Timeout(120 * time.Second).
    Build()
```

## Tipos de Datos Principales

### Producto

```go
type Product struct {
    ID                int                `json:"id,omitempty"`
    Name              string             `json:"name"`
    Type              ProductType        `json:"type"`
    Status            ProductStatus      `json:"status"`
    RegularPrice      string             `json:"regular_price"`
    SalePrice         string             `json:"sale_price,omitempty"`
    Description       string             `json:"description"`
    ShortDescription  string             `json:"short_description"`
    SKU               string             `json:"sku"`
    Categories        []ProductCategory  `json:"categories,omitempty"`
    Images            []ProductImage     `json:"images,omitempty"`
    // ... más campos
}
```

### Orden

```go
type Order struct {
    ID          int           `json:"id,omitempty"`
    Number      string        `json:"number,omitempty"`
    Status      OrderStatus   `json:"status"`
    Currency    string        `json:"currency"`
    Total       string        `json:"total,omitempty"`
    LineItems   []LineItem    `json:"line_items,omitempty"`
    Billing     *Address      `json:"billing,omitempty"`
    Shipping    *Address      `json:"shipping,omitempty"`
    // ... más campos
}
```

### Tienda

```go
type Store struct {
    ID        int      `json:"id"`
    StoreName string   `json:"store_name"`
    FirstName string   `json:"first_name"`
    LastName  string   `json:"last_name"`
    Email     string   `json:"email"`
    Address   *Address `json:"address,omitempty"`
    Rating    *Rating  `json:"rating,omitempty"`
    // ... más campos
}
```

## Constantes

### Estados de Producto

```go
const (
    ProductStatusDraft   ProductStatus = "draft"
    ProductStatusPending ProductStatus = "pending"
    ProductStatusPublish ProductStatus = "publish"
)
```

### Tipos de Producto

```go
const (
    ProductTypeSimple   ProductType = "simple"
    ProductTypeGrouped  ProductType = "grouped"
    ProductTypeExternal ProductType = "external"
    ProductTypeVariable ProductType = "variable"
)
```

### Estados de Orden

```go
const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusProcessing OrderStatus = "processing"
    OrderStatusCompleted  OrderStatus = "completed"
    OrderStatusCancelled  OrderStatus = "cancelled"
    // ... más estados
)
```

## Ejemplos

Consulta la carpeta `examples/` para ver ejemplos completos:

- [`basic_usage.go`](examples/basic_usage.go) - Uso básico del SDK
- [`advanced_usage.go`](examples/advanced_usage.go) - Funcionalidades avanzadas

## Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.

## Soporte

Si encuentras algún problema o tienes preguntas:

1. Revisa la [documentación de la API de Dokan](https://getdokan.github.io/dokan/)
2. Busca en los [issues existentes](https://github.com/diogenes-moreira/dokan-go-sdk/issues)
3. Crea un nuevo issue si no encuentras una solución

## Roadmap

- [ ] Soporte para webhooks
- [ ] Cliente para GraphQL API (si está disponible)
- [ ] Más ejemplos y casos de uso
- [ ] Métricas y logging integrado
- [ ] Soporte para upload de archivos
- [ ] CLI tool para operaciones comunes

---

Desarrollado con ❤️ para la comunidad de Go y Dokan.

