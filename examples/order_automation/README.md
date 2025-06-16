# Ejemplo: Automatización de Órdenes

Este ejemplo demuestra cómo crear un sistema de procesamiento automático de órdenes que verifica y procesa órdenes pendientes de forma periódica.

## Características

- Procesamiento automático de órdenes pendientes
- Validación completa de órdenes (productos, pagos, inventario)
- Actualización automática de estados
- Notificaciones a clientes
- Generación de reportes detallados
- Ejecución periódica configurable
- Logging completo de actividades

## Configuración

### Variables de Entorno

```bash
export DOKAN_BASE_URL="https://tu-sitio-dokan.com"
export DOKAN_USERNAME="tu-usuario"
export DOKAN_PASSWORD="tu-contraseña"
```

### Configuración del Procesador

```go
config := ProcessorConfig{
    CheckInterval:   5 * time.Minute,  // Verificar cada 5 minutos
    MaxOrdersPerRun: 50,               // Máximo 50 órdenes por ejecución
    AutoApprove:     true,             // Aprobar automáticamente
    NotifyCustomers: true,             // Enviar notificaciones
    LogFile:         "order_processing.log",
}
```

## Flujo de Procesamiento

### 1. Obtención de Órdenes

El sistema busca órdenes con estado `pending` ordenadas por fecha de creación.

### 2. Validación

Para cada orden se verifica:

- **Productos válidos**: La orden debe tener al menos un producto
- **Total válido**: Debe tener un total mayor a cero
- **Información de facturación**: Email y datos de contacto
- **Dirección de envío**: Si es requerida para el tipo de producto

### 3. Verificación de Inventario

- Verifica que los productos estén disponibles
- Comprueba el estado de publicación
- Valida cantidades disponibles (si está implementado)

### 4. Verificación de Pago

- Valida el método de pago
- Verifica transacciones para pagos electrónicos
- Identifica pagos que requieren verificación manual

### 5. Acciones Automáticas

Según el resultado de las validaciones:

- **Aprobar**: Cambia estado a `processing`
- **Rechazar**: Cambia estado a `cancelled`
- **En espera**: Cambia estado a `on-hold`
- **Mantener pendiente**: Para pagos que requieren verificación

## Uso

### Ejecución Simple

```bash
go run main.go
```

### Ejecución como Servicio

```bash
# Compilar
go build -o order-processor main.go

# Ejecutar en background
nohup ./order-processor > /dev/null 2>&1 &

# Verificar logs
tail -f order_processing.log
```

### Ejecución con Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o order-processor main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-processor .
CMD ["./order-processor"]
```

## Ejemplo de Salida

```
2024/01/15 14:30:00 Iniciando procesador automático de órdenes...
2024/01/15 14:30:00 Intervalo de verificación: 5m0s
2024/01/15 14:30:00 Máximo órdenes por ejecución: 50
2024/01/15 14:30:00 Buscando órdenes para procesar...
2024/01/15 14:30:01 Encontradas 3 órdenes pendientes
2024/01/15 14:30:01 Procesando orden #1001 (ID: 123)...
2024/01/15 14:30:02 Orden #1001 aprobada automáticamente
2024/01/15 14:30:03 Procesando orden #1002 (ID: 124)...
2024/01/15 14:30:04 Orden #1002 en espera: Inventario insuficiente: producto Camiseta XL no está disponible
2024/01/15 14:30:05 Procesando orden #1003 (ID: 125)...
2024/01/15 14:30:06 Orden #1003 cancelada: Validación fallida: email inválido: cliente@
2024/01/15 14:30:06 Procesamiento completado: 2 exitosas, 1 fallidas
2024/01/15 14:30:06 Reporte generado: order_processing_report_2024-01-15_14-30-06.txt
```

## Personalización

### Validaciones Personalizadas

```go
func (p *OrderProcessor) validateOrder(order dokan.Order) error {
    // Validación de monto mínimo
    if total, err := strconv.ParseFloat(order.Total, 64); err == nil {
        if total < 10.0 {
            return fmt.Errorf("monto mínimo de orden es $10.00")
        }
    }
    
    // Validación de productos restringidos
    for _, item := range order.LineItems {
        if strings.Contains(strings.ToLower(item.Name), "restringido") {
            return fmt.Errorf("producto restringido: %s", item.Name)
        }
    }
    
    return nil
}
```

### Lógica de Inventario Personalizada

```go
func (p *OrderProcessor) checkInventory(order dokan.Order) error {
    // Conectar con sistema de inventario externo
    inventoryAPI := NewInventoryAPI()
    
    for _, item := range order.LineItems {
        available, err := inventoryAPI.CheckStock(item.SKU)
        if err != nil {
            return fmt.Errorf("error verificando stock para %s: %w", item.SKU, err)
        }
        
        if available < item.Quantity {
            return fmt.Errorf("stock insuficiente para %s: disponible %d, requerido %d", 
                item.SKU, available, item.Quantity)
        }
    }
    
    return nil
}
```

### Notificaciones Personalizadas

```go
func (p *OrderProcessor) notifyCustomer(order dokan.Order, message string) error {
    // Envío por email
    emailService := NewEmailService()
    if err := emailService.SendEmail(order.Billing.Email, "Estado de su orden", message); err != nil {
        return fmt.Errorf("error enviando email: %w", err)
    }
    
    // Envío por SMS (opcional)
    if order.Billing.Phone != "" {
        smsService := NewSMSService()
        if err := smsService.SendSMS(order.Billing.Phone, message); err != nil {
            log.Printf("Error enviando SMS: %v", err) // No fallar por SMS
        }
    }
    
    return nil
}
```

## Integración con Sistemas Externos

### ERP/Inventario

```go
type InventoryAPI struct {
    baseURL string
    apiKey  string
}

func (api *InventoryAPI) CheckStock(sku string) (int, error) {
    // Implementar llamada a API de inventario
    resp, err := http.Get(fmt.Sprintf("%s/stock/%s?api_key=%s", api.baseURL, sku, api.apiKey))
    // ... procesar respuesta
    return stock, nil
}

func (api *InventoryAPI) ReserveStock(sku string, quantity int) error {
    // Implementar reserva de stock
    return nil
}
```

### Sistema de Pagos

```go
func (p *OrderProcessor) verifyPaymentWithGateway(order dokan.Order) error {
    switch order.PaymentMethod {
    case "stripe":
        return p.verifyStripePayment(order.TransactionID)
    case "paypal":
        return p.verifyPayPalPayment(order.TransactionID)
    default:
        return nil
    }
}

func (p *OrderProcessor) verifyStripePayment(transactionID string) error {
    // Verificar con API de Stripe
    return nil
}
```

## Monitoreo y Alertas

### Métricas

```go
type ProcessorMetrics struct {
    OrdersProcessed   int64
    OrdersApproved    int64
    OrdersRejected    int64
    OrdersOnHold      int64
    ProcessingErrors  int64
    LastRunTime       time.Time
}

func (p *OrderProcessor) updateMetrics(action string) {
    // Actualizar métricas
    // Enviar a sistema de monitoreo (Prometheus, etc.)
}
```

### Alertas

```go
func (p *OrderProcessor) checkAlerts() {
    // Verificar si hay muchos errores
    if p.metrics.ProcessingErrors > 10 {
        p.sendAlert("Alto número de errores en procesamiento de órdenes")
    }
    
    // Verificar si no se han procesado órdenes en mucho tiempo
    if time.Since(p.metrics.LastRunTime) > 30*time.Minute {
        p.sendAlert("Procesador de órdenes no ha ejecutado recientemente")
    }
}
```

## Consideraciones de Producción

### Concurrencia

Para evitar procesamiento duplicado en múltiples instancias:

```go
func (p *OrderProcessor) acquireLock(orderID int) (bool, error) {
    // Implementar lock distribuido (Redis, base de datos, etc.)
    return true, nil
}

func (p *OrderProcessor) releaseLock(orderID int) error {
    // Liberar lock
    return nil
}
```

### Recuperación de Errores

```go
func (p *OrderProcessor) handleCriticalError(err error, order dokan.Order) {
    // Log detallado
    log.Printf("Error crítico procesando orden #%s: %v", order.Number, err)
    
    // Notificar administradores
    p.notifyAdmins(fmt.Sprintf("Error crítico en orden #%s: %v", order.Number, err))
    
    // Marcar orden para revisión manual
    p.flagForManualReview(order.ID, err.Error())
}
```

### Escalabilidad

- Usar colas de mensajes para procesamiento asíncrono
- Implementar sharding por vendedor o región
- Usar cache para datos frecuentemente accedidos
- Implementar circuit breakers para APIs externas

## Troubleshooting

### Órdenes no se procesan

1. Verificar logs de errores
2. Comprobar conectividad con API de Dokan
3. Validar credenciales y permisos
4. Revisar configuración de filtros

### Errores de rate limiting

1. Aumentar pausas entre operaciones
2. Reducir `MaxOrdersPerRun`
3. Implementar backoff exponencial

### Notificaciones no se envían

1. Verificar configuración de email/SMS
2. Comprobar formato de direcciones
3. Revisar logs de servicios externos

### Alto uso de memoria

1. Procesar órdenes en lotes más pequeños
2. Liberar recursos después de cada orden
3. Implementar garbage collection manual si es necesario

