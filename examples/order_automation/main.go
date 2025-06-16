package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/diogenes-moreira/dokan-go-sdk"
)

// OrderProcessor maneja el procesamiento automático de órdenes
type OrderProcessor struct {
	client *dokan.Client
	config ProcessorConfig
}

// ProcessorConfig contiene la configuración del procesador
type ProcessorConfig struct {
	CheckInterval    time.Duration
	MaxOrdersPerRun  int
	AutoApprove      bool
	NotifyCustomers  bool
	LogFile          string
}

// OrderAction representa una acción a realizar en una orden
type OrderAction struct {
	OrderID    int
	Action     string
	Reason     string
	NewStatus  dokan.OrderStatus
	Timestamp  time.Time
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

	// Configurar procesador
	config := ProcessorConfig{
		CheckInterval:   5 * time.Minute,
		MaxOrdersPerRun: 50,
		AutoApprove:     true,
		NotifyCustomers: true,
		LogFile:         "order_processing.log",
	}

	processor := &OrderProcessor{
		client: client,
		config: config,
	}

	// Configurar logging
	if err := setupLogging(config.LogFile); err != nil {
		log.Fatalf("Error configurando logging: %v", err)
	}

	log.Printf("Iniciando procesador automático de órdenes...")
	log.Printf("Intervalo de verificación: %v", config.CheckInterval)
	log.Printf("Máximo órdenes por ejecución: %d", config.MaxOrdersPerRun)

	// Ejecutar procesamiento inicial
	ctx := context.Background()
	if err := processor.ProcessOrders(ctx); err != nil {
		log.Printf("Error en procesamiento inicial: %v", err)
	}

	// Configurar ticker para procesamiento periódico
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	// Loop principal
	for {
		select {
		case <-ticker.C:
			log.Printf("Iniciando verificación periódica de órdenes...")
			if err := processor.ProcessOrders(ctx); err != nil {
				log.Printf("Error en procesamiento periódico: %v", err)
			}
		}
	}
}

// ProcessOrders procesa las órdenes pendientes
func (p *OrderProcessor) ProcessOrders(ctx context.Context) error {
	log.Printf("Buscando órdenes para procesar...")

	// Obtener órdenes pendientes
	pendingOrders, err := p.getPendingOrders(ctx)
	if err != nil {
		return fmt.Errorf("error obteniendo órdenes pendientes: %w", err)
	}

	if len(pendingOrders) == 0 {
		log.Printf("No hay órdenes pendientes para procesar")
		return nil
	}

	log.Printf("Encontradas %d órdenes pendientes", len(pendingOrders))

	var processed, failed int
	var actions []OrderAction

	// Procesar cada orden
	for i, order := range pendingOrders {
		if i >= p.config.MaxOrdersPerRun {
			log.Printf("Límite de órdenes por ejecución alcanzado (%d)", p.config.MaxOrdersPerRun)
			break
		}

		log.Printf("Procesando orden #%s (ID: %d)...", order.Number, order.ID)

		action, err := p.processOrder(ctx, order)
		if err != nil {
			log.Printf("Error procesando orden #%s: %v", order.Number, err)
			failed++
			continue
		}

		actions = append(actions, action)
		processed++

		// Pausa entre órdenes para evitar rate limiting
		time.Sleep(1 * time.Second)
	}

	// Generar resumen
	log.Printf("Procesamiento completado: %d exitosas, %d fallidas", processed, failed)

	// Generar reporte
	if err := p.generateReport(actions); err != nil {
		log.Printf("Error generando reporte: %v", err)
	}

	return nil
}

// getPendingOrders obtiene las órdenes pendientes de procesamiento
func (p *OrderProcessor) getPendingOrders(ctx context.Context) ([]dokan.Order, error) {
	params := &dokan.OrderListParams{
		ListParams: dokan.ListParams{
			Page:    1,
			PerPage: p.config.MaxOrdersPerRun,
			OrderBy: "date",
			Order:   "asc",
		},
		Status: []dokan.OrderStatus{dokan.OrderStatusPending},
	}

	result, err := p.client.Orders.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return result.Orders, nil
}

// processOrder procesa una orden individual
func (p *OrderProcessor) processOrder(ctx context.Context, order dokan.Order) (OrderAction, error) {
	action := OrderAction{
		OrderID:   order.ID,
		Timestamp: time.Now(),
	}

	// Validar la orden
	if err := p.validateOrder(order); err != nil {
		action.Action = "rejected"
		action.Reason = fmt.Sprintf("Validación fallida: %v", err)
		action.NewStatus = dokan.OrderStatusCancelled

		// Cancelar orden
		if err := p.updateOrderStatus(ctx, order.ID, dokan.OrderStatusCancelled); err != nil {
			return action, fmt.Errorf("error cancelando orden: %w", err)
		}

		log.Printf("Orden #%s cancelada: %s", order.Number, action.Reason)
		return action, nil
	}

	// Verificar inventario
	if err := p.checkInventory(order); err != nil {
		action.Action = "on_hold"
		action.Reason = fmt.Sprintf("Inventario insuficiente: %v", err)
		action.NewStatus = dokan.OrderStatusOnHold

		// Poner en espera
		if err := p.updateOrderStatus(ctx, order.ID, dokan.OrderStatusOnHold); err != nil {
			return action, fmt.Errorf("error poniendo orden en espera: %w", err)
		}

		log.Printf("Orden #%s en espera: %s", order.Number, action.Reason)
		return action, nil
	}

	// Verificar pago
	if err := p.verifyPayment(order); err != nil {
		action.Action = "payment_pending"
		action.Reason = fmt.Sprintf("Pago pendiente: %v", err)
		action.NewStatus = dokan.OrderStatusPending

		log.Printf("Orden #%s mantiene estado pendiente: %s", order.Number, action.Reason)
		return action, nil
	}

	// Si todo está bien, aprobar automáticamente
	if p.config.AutoApprove {
		action.Action = "approved"
		action.Reason = "Orden aprobada automáticamente"
		action.NewStatus = dokan.OrderStatusProcessing

		if err := p.updateOrderStatus(ctx, order.ID, dokan.OrderStatusProcessing); err != nil {
			return action, fmt.Errorf("error aprobando orden: %w", err)
		}

		log.Printf("Orden #%s aprobada automáticamente", order.Number)

		// Enviar notificación al cliente si está habilitado
		if p.config.NotifyCustomers {
			if err := p.notifyCustomer(order, "Su orden ha sido aprobada y está siendo procesada"); err != nil {
				log.Printf("Error enviando notificación para orden #%s: %v", order.Number, err)
			}
		}
	}

	return action, nil
}

// validateOrder valida una orden antes de procesarla
func (p *OrderProcessor) validateOrder(order dokan.Order) error {
	// Verificar que tenga productos
	if len(order.LineItems) == 0 {
		return fmt.Errorf("orden sin productos")
	}

	// Verificar que tenga total válido
	if order.Total == "" || order.Total == "0" {
		return fmt.Errorf("orden sin total válido")
	}

	// Verificar información de facturación
	if order.Billing == nil {
		return fmt.Errorf("orden sin información de facturación")
	}

	if order.Billing.Email == "" {
		return fmt.Errorf("orden sin email de contacto")
	}

	// Verificar que el email sea válido
	if !strings.Contains(order.Billing.Email, "@") {
		return fmt.Errorf("email inválido: %s", order.Billing.Email)
	}

	// Verificar información de envío si es necesaria
	if order.Shipping != nil && order.Shipping.Address1 == "" {
		return fmt.Errorf("dirección de envío incompleta")
	}

	return nil
}

// checkInventory verifica que haya suficiente inventario
func (p *OrderProcessor) checkInventory(order dokan.Order) error {
	ctx := context.Background()

	for _, item := range order.LineItems {
		// Obtener información del producto
		product, err := p.client.Products.Get(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("error obteniendo producto %d: %w", item.ProductID, err)
		}

		// Verificar que el producto esté publicado
		if product.Status != dokan.ProductStatusPublish {
			return fmt.Errorf("producto %s no está disponible", product.Name)
		}

		// Aquí podrías agregar lógica adicional para verificar stock
		// Por ejemplo, si tienes un campo de stock en el producto
		log.Printf("Verificando inventario para producto: %s (cantidad: %d)", product.Name, item.Quantity)
	}

	return nil
}

// verifyPayment verifica el estado del pago
func (p *OrderProcessor) verifyPayment(order dokan.Order) error {
	// Verificar método de pago
	if order.PaymentMethod == "" {
		return fmt.Errorf("método de pago no especificado")
	}

	// Para pagos que requieren verificación manual
	manualPaymentMethods := []string{"bacs", "cheque", "cod"}
	for _, method := range manualPaymentMethods {
		if order.PaymentMethod == method {
			return fmt.Errorf("método de pago requiere verificación manual: %s", order.PaymentMethodTitle)
		}
	}

	// Verificar si hay ID de transacción para pagos electrónicos
	electronicMethods := []string{"stripe", "paypal", "square"}
	for _, method := range electronicMethods {
		if strings.Contains(order.PaymentMethod, method) && order.TransactionID == "" {
			return fmt.Errorf("falta ID de transacción para pago electrónico")
		}
	}

	return nil
}

// updateOrderStatus actualiza el estado de una orden
func (p *OrderProcessor) updateOrderStatus(ctx context.Context, orderID int, status dokan.OrderStatus) error {
	update := &dokan.OrderUpdate{
		Status: &status,
	}

	_, err := p.client.Orders.Update(ctx, orderID, update)
	return err
}

// notifyCustomer envía una notificación al cliente
func (p *OrderProcessor) notifyCustomer(order dokan.Order, message string) error {
	// Aquí implementarías la lógica de notificación
	// Por ejemplo, envío de email, SMS, etc.
	log.Printf("Notificación enviada a %s: %s", order.Billing.Email, message)
	return nil
}

// generateReport genera un reporte de las acciones realizadas
func (p *OrderProcessor) generateReport(actions []OrderAction) error {
	if len(actions) == 0 {
		return nil
	}

	filename := fmt.Sprintf("order_processing_report_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo de reporte: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "Reporte de Procesamiento de Órdenes\n")
	fmt.Fprintf(file, "Fecha: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "=====================================\n\n")

	// Resumen por acción
	actionCounts := make(map[string]int)
	for _, action := range actions {
		actionCounts[action.Action]++
	}

	fmt.Fprintf(file, "Resumen:\n")
	for action, count := range actionCounts {
		fmt.Fprintf(file, "- %s: %d\n", action, count)
	}

	fmt.Fprintf(file, "\nDetalle de acciones:\n")
	fmt.Fprintf(file, "===================\n")

	for _, action := range actions {
		fmt.Fprintf(file, "Orden ID: %d\n", action.OrderID)
		fmt.Fprintf(file, "Acción: %s\n", action.Action)
		fmt.Fprintf(file, "Razón: %s\n", action.Reason)
		fmt.Fprintf(file, "Nuevo estado: %s\n", action.NewStatus)
		fmt.Fprintf(file, "Timestamp: %s\n", action.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(file, "---\n")
	}

	log.Printf("Reporte generado: %s", filename)
	return nil
}

// setupLogging configura el logging a archivo
func setupLogging(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}

