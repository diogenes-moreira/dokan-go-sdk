# Ejemplo: Sincronización de Inventario

Este ejemplo demuestra cómo sincronizar un inventario externo (desde un archivo CSV) con los productos de Dokan.

## Características

- Lee inventario desde archivo CSV
- Compara con productos existentes en Dokan
- Crea productos nuevos automáticamente
- Actualiza productos existentes si hay cambios
- Genera reporte detallado de la sincronización
- Maneja rate limiting y errores

## Configuración

### Variables de Entorno

Configura las siguientes variables de entorno:

```bash
export DOKAN_BASE_URL="https://tu-sitio-dokan.com"
export DOKAN_USERNAME="tu-usuario"
export DOKAN_PASSWORD="tu-contraseña"
```

### Formato del Archivo CSV

El archivo `inventory.csv` debe tener las siguientes columnas:

| Columna     | Descripción                    | Ejemplo                                    |
|-------------|--------------------------------|--------------------------------------------|
| SKU         | Código único del producto      | SHIRT-001                                  |
| Name        | Nombre del producto            | Camiseta Básica Blanca                     |
| Price       | Precio regular                 | 19.99                                      |
| Stock       | Cantidad en inventario         | 50                                         |
| Description | Descripción del producto       | Camiseta de algodón 100% en color blanco  |
| CategoryID  | ID de categoría en Dokan       | 15                                         |
| ImageURL    | URL de la imagen del producto  | https://ejemplo.com/camiseta-blanca.jpg    |
| Featured    | Si es producto destacado       | true/false                                 |

## Uso

1. Prepara tu archivo CSV con el inventario:
   ```bash
   cp inventory.csv mi_inventario.csv
   # Edita mi_inventario.csv con tus productos
   ```

2. Ejecuta la sincronización:
   ```bash
   go run main.go
   ```

3. Revisa el reporte generado:
   ```bash
   cat sync_report_*.txt
   ```

## Ejemplo de Salida

```
Sincronizando 8 productos del inventario...
Procesando item 1/8: SHIRT-001
✓ Creado: SHIRT-001
Procesando item 2/8: SHIRT-002
✓ Creado: SHIRT-002
Procesando item 3/8: PANTS-001
✓ Actualizado: PANTS-001
Procesando item 4/8: SHOES-001
- Sin cambios: SHOES-001
...

=== Resumen de Sincronización ===
Productos creados: 5
Productos actualizados: 2
Productos sin cambios: 1
Total procesados: 8

Reporte generado: sync_report_2024-01-15_14-30-25.txt
```

## Personalización

### Modificar Lógica de Actualización

Puedes personalizar qué campos se comparan para determinar si un producto necesita actualización:

```go
func needsUpdate(existing *dokan.Product, item InventoryItem) bool {
    // Agregar más campos a comparar
    if existing.ShortDescription != item.Description {
        return true
    }
    // Tu lógica personalizada aquí
    return false
}
```

### Agregar Validaciones

```go
func validateInventoryItem(item InventoryItem) error {
    if item.SKU == "" {
        return fmt.Errorf("SKU es requerido")
    }
    
    if price, err := strconv.ParseFloat(item.Price, 64); err != nil || price <= 0 {
        return fmt.Errorf("precio inválido: %s", item.Price)
    }
    
    return nil
}
```

### Configurar Categorías Automáticamente

```go
func mapCategory(categoryName string) int {
    categories := map[string]int{
        "ropa":      15,
        "zapatos":   17,
        "accesorios": 18,
    }
    
    if id, exists := categories[strings.ToLower(categoryName)]; exists {
        return id
    }
    
    return 1 // Categoría por defecto
}
```

## Consideraciones

### Rate Limiting

El ejemplo incluye pausas entre operaciones para evitar rate limiting:

```go
time.Sleep(500 * time.Millisecond) // Entre productos
time.Sleep(200 * time.Millisecond) // Entre páginas
```

Ajusta estos valores según los límites de tu servidor.

### Manejo de Errores

Los errores se registran pero no detienen el proceso completo:

```go
if err := createProduct(client, ctx, item); err != nil {
    log.Printf("Error creando producto %s: %v", item.SKU, err)
    continue // Continúa con el siguiente producto
}
```

### Backup

Antes de ejecutar sincronizaciones masivas, considera hacer un backup de tu base de datos.

## Extensiones Posibles

- Soporte para productos variables
- Sincronización de stock en tiempo real
- Integración con sistemas ERP
- Notificaciones por email de cambios
- Interfaz web para monitoreo
- Programación automática (cron jobs)

## Troubleshooting

### Error: "authentication error"
- Verifica las credenciales en las variables de entorno
- Asegúrate de que el usuario tenga permisos de API

### Error: "rate limit exceeded"
- Aumenta las pausas entre operaciones
- Reduce el tamaño de lotes

### Error: "validation error"
- Verifica el formato del archivo CSV
- Asegúrate de que todos los campos requeridos estén presentes

### Productos no se actualizan
- Verifica la lógica en `needsUpdate()`
- Confirma que los SKUs coincidan exactamente

