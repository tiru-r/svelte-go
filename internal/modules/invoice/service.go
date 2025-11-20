package invoice

import (
	"fmt"
	"log"
	"time"

	"svelte-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

type Service struct {
	eventBus types.EventBus
	repo     *Repository
}

func NewService(eventBus types.EventBus, db *badger.DB) *Service {
	service := &Service{
		eventBus: eventBus,
		repo:     NewRepository(db),
	}

	service.setupEventSubscriptions()
	return service
}

func (s *Service) setupEventSubscriptions() {
	s.eventBus.SubscribeQueue("time.entry.completed", "invoice_service", s.handleTimeEntryCompleted)
	s.eventBus.SubscribeQueue("expense.created", "invoice_service", s.handleExpenseCreated)

	log.Println("Invoice service event subscriptions configured")
}

func (s *Service) CreateInvoice(userID, clientID, projectID string, items []types.InvoiceItem) (*types.Invoice, error) {
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Amount
	}

	invoice := &types.Invoice{
		ID:          types.GenerateID(),
		UserID:      userID,
		ClientID:    clientID,
		ProjectID:   projectID,
		Number:      s.generateInvoiceNumber(),
		Items:       items,
		TotalAmount: totalAmount,
		Currency:    "USD",
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	event := types.NewEvent("invoice_created", "invoice_service", map[string]any{
		"invoice_id":   invoice.ID,
		"client_id":    invoice.ClientID,
		"project_id":   invoice.ProjectID,
		"total_amount": invoice.TotalAmount,
		"status":       invoice.Status,
	})

	s.eventBus.Publish("invoice.created", event)
	log.Printf("📄 Invoice created: %s ($%.2f)", invoice.Number, invoice.TotalAmount)

	return invoice, nil
}

func (s *Service) GenerateFromTimeEntries(userID, clientID, projectID string, hourlyRate float64, timeEntries []*types.TimeEntry) (*types.Invoice, error) {

	var items []types.InvoiceItem
	var totalHours float64

	for _, entry := range timeEntries {
		if entry.EndTime != nil {
			hours := entry.EndTime.Sub(entry.StartTime).Hours()
			totalHours += hours

			items = append(items, types.InvoiceItem{
				Description: entry.Description,
				Quantity:    hours,
				Rate:        hourlyRate,
				Amount:      hours * hourlyRate,
			})
		}
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no completed time entries found for the specified period")
	}

	invoice, err := s.CreateInvoice(userID, clientID, projectID, items)
	if err != nil {
		return nil, err
	}

	event := types.NewEvent("invoice_generated", "invoice_service", map[string]any{
		"invoice_id":   invoice.ID,
		"client_id":    invoice.ClientID,
		"project_id":   invoice.ProjectID,
		"total_hours":  totalHours,
		"total_amount": invoice.TotalAmount,
		"entry_count":  len(timeEntries),
	})

	s.eventBus.Publish("invoice.generated", event)
	log.Printf("🧾 Invoice generated from %.1f hours of work", totalHours)

	return invoice, nil
}

func (s *Service) GetInvoices(userID string) ([]*types.Invoice, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) GetInvoicesByClient(clientID string) ([]*types.Invoice, error) {
	return s.repo.GetByClientID(clientID)
}

func (s *Service) UpdateInvoiceStatus(invoiceID, status string) (*types.Invoice, error) {
	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	oldStatus := invoice.Status
	invoice.Status = status
	invoice.UpdatedAt = time.Now()

	if status == "sent" && invoice.SentAt == nil {
		now := time.Now()
		invoice.SentAt = &now
	}

	if status == "paid" && invoice.PaidAt == nil {
		now := time.Now()
		invoice.PaidAt = &now
	}

	err = s.repo.Update(invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	event := types.NewEvent("invoice_status_updated", "invoice_service", map[string]any{
		"invoice_id": invoice.ID,
		"old_status": oldStatus,
		"new_status": status,
		"client_id":  invoice.ClientID,
		"project_id": invoice.ProjectID,
	})

	s.eventBus.Publish("invoice.status_updated", event)
	return invoice, nil
}

func (s *Service) DeleteInvoice(invoiceID string) error {
	invoice, err := s.repo.GetByID(invoiceID)
	if err != nil {
		return fmt.Errorf("invoice not found: %w", err)
	}

	if invoice.Status == "paid" {
		return fmt.Errorf("cannot delete paid invoice")
	}

	err = s.repo.Delete(invoiceID)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	event := types.NewEvent("invoice_deleted", "invoice_service", map[string]any{
		"invoice_id": invoice.ID,
		"client_id":  invoice.ClientID,
		"project_id": invoice.ProjectID,
	})

	s.eventBus.Publish("invoice.deleted", event)
	return nil
}

func (s *Service) generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", time.Now().Unix())
}

func (s *Service) handleTimeEntryCompleted(event *types.Event) error {
	log.Printf("⏱️ Time entry completed - ready for invoicing: %v", event.Data["time_entry_id"])
	return nil
}

func (s *Service) handleExpenseCreated(event *types.Event) error {
	log.Printf("💰 Expense created - consider adding to next invoice: %v", event.Data["expense_id"])
	return nil
}
