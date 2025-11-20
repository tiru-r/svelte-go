package expense

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
	s.eventBus.SubscribeQueue("client.project.started", "expense_service", s.handleProjectStarted)
	s.eventBus.SubscribeQueue("time.entry.completed", "expense_service", s.handleTimeEntryCompleted)

	log.Println("Expense service event subscriptions configured")
}

func (s *Service) CreateExpense(userID, projectID, category, description string, amount float64) (*types.Expense, error) {
	expense := &types.Expense{
		ID:          types.GenerateID(),
		UserID:      userID,
		ProjectID:   projectID,
		Category:    category,
		Description: description,
		Amount:      amount,
		Currency:    "USD",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(expense)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	event := types.NewEvent("expense_created", "expense_service", map[string]any{
		"expense_id": expense.ID,
		"user_id":    expense.UserID,
		"project_id": expense.ProjectID,
		"amount":     expense.Amount,
		"category":   expense.Category,
	})

	s.eventBus.Publish("expense.created", event)
	log.Printf("💰 Expense created: $%.2f for %s", expense.Amount, expense.Description)

	return expense, nil
}

func (s *Service) GetExpenses(userID string) ([]*types.Expense, error) {
	return s.repo.GetByUserID(userID)
}

func (s *Service) GetExpensesByProject(projectID string) ([]*types.Expense, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *Service) UpdateExpense(expenseID string, updates map[string]any) (*types.Expense, error) {
	expense, err := s.repo.GetByID(expenseID)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %w", err)
	}

	if description, ok := updates["description"].(string); ok {
		expense.Description = description
	}
	if amount, ok := updates["amount"].(float64); ok {
		expense.Amount = amount
	}
	if category, ok := updates["category"].(string); ok {
		expense.Category = category
	}

	expense.UpdatedAt = time.Now()

	err = s.repo.Update(expense)
	if err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	event := types.NewEvent("expense_updated", "expense_service", map[string]any{
		"expense_id": expense.ID,
		"user_id":    expense.UserID,
		"updates":    updates,
	})

	s.eventBus.Publish("expense.updated", event)
	return expense, nil
}

func (s *Service) DeleteExpense(expenseID string) error {
	expense, err := s.repo.GetByID(expenseID)
	if err != nil {
		return fmt.Errorf("expense not found: %w", err)
	}

	err = s.repo.Delete(expenseID)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	event := types.NewEvent("expense_deleted", "expense_service", map[string]any{
		"expense_id": expense.ID,
		"user_id":    expense.UserID,
		"project_id": expense.ProjectID,
	})

	s.eventBus.Publish("expense.deleted", event)
	return nil
}

func (s *Service) handleProjectStarted(event *types.Event) error {
	log.Printf("📋 Project started - expense tracking enabled for project: %v", event.Data["project_id"])
	return nil
}

func (s *Service) handleTimeEntryCompleted(event *types.Event) error {
	log.Printf("⏱️ Time entry completed - consider adding related expenses for: %v", event.Data["project_id"])
	return nil
}
