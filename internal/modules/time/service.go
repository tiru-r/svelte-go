package time

import (
	"fmt"
	"log"
	"time"

	"svelte-go/internal/shared/types"

	"github.com/dgraph-io/badger/v4"
)

// Service handles time tracking business logic
type Service struct {
	eventBus types.EventBus
	repo     *Repository
}

// NewService creates a new time tracking service
func NewService(eventBus types.EventBus, db *badger.DB) *Service {
	service := &Service{
		eventBus: eventBus,
		repo:     NewRepository(db),
	}

	// Subscribe to relevant events
	service.setupEventSubscriptions()

	return service
}

// setupEventSubscriptions configures event handlers
func (s *Service) setupEventSubscriptions() {
	// Listen for client project events
	s.eventBus.SubscribeQueue("client.project.started", "time_service", s.handleProjectStarted)

	// Listen for invoice events
	s.eventBus.SubscribeQueue("invoice.generated", "time_service", s.handleInvoiceGenerated)

	// Listen for system events
	s.eventBus.SubscribeQueue("system.user.logout", "time_service", s.handleUserLogout)

	log.Println("Time service event subscriptions configured")
}

// StartTimer starts a new time tracking session
func (s *Service) StartTimer(userID, projectID, description string) (*types.TimeEntry, error) {
	// Stop any existing timer for this user
	err := s.stopActiveTimer(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to stop existing timer: %w", err)
	}

	// Create new time entry
	entry := types.NewTimeEntry(userID, projectID, description)

	// Save to Badger
	err = s.repo.Save(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to save time entry: %w", err)
	}

	// Publish event
	event := types.NewEvent("session_started", "time_service", map[string]any{
		"time_entry_id": entry.ID,
		"user_id":       userID,
		"project_id":    projectID,
		"description":   description,
		"start_time":    entry.StartTime,
	}).WithAggregateID(entry.ID)

	err = s.eventBus.Publish("time.session.started", event)
	if err != nil {
		return nil, fmt.Errorf("failed to publish start event: %w", err)
	}

	log.Printf("⏱️  Timer started: %s for project %s", userID, projectID)
	return entry, nil
}

// StopTimer stops the active timer for a user
func (s *Service) StopTimer(userID string) (*types.TimeEntry, error) {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active timer: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("no active timer for user %s", userID)
	}

	// Stop the timer
	entry.Stop()

	// Calculate amount based on hourly rate (would come from project/client)
	entry.HourlyRate = 75.0 // Demo rate
	amount := entry.CalculateAmount()

	// Save updated entry
	err = s.repo.Save(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to save stopped timer: %w", err)
	}

	// Publish event
	event := types.NewEvent("session_stopped", "time_service", map[string]any{
		"time_entry_id": entry.ID,
		"user_id":       userID,
		"project_id":    entry.ProjectID,
		"duration":      entry.Duration,
		"amount":        amount,
		"start_time":    entry.StartTime,
		"end_time":      entry.EndTime,
	}).WithAggregateID(entry.ID)

	err = s.eventBus.Publish("time.session.stopped", event)
	if err != nil {
		return nil, fmt.Errorf("failed to publish stop event: %w", err)
	}

	log.Printf("⏹️  Timer stopped: %s (Duration: %d seconds, Amount: $%.2f)",
		userID, entry.Duration, amount)

	return entry, nil
}

// PauseTimer pauses the active timer
func (s *Service) PauseTimer(userID string) (*types.TimeEntry, error) {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active timer: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("no active timer for user %s", userID)
	}

	if !entry.IsRunning {
		return entry, nil // Already paused
	}

	// Calculate current duration
	now := time.Now()
	entry.Duration += int64(now.Sub(entry.StartTime).Seconds())
	entry.StartTime = now // Reset start time for resume
	entry.IsRunning = false
	entry.UpdatedAt = now

	// Save updated entry
	err = s.repo.Save(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to save paused timer: %w", err)
	}

	// Publish event
	event := types.NewEvent("session_paused", "time_service", map[string]any{
		"time_entry_id":    entry.ID,
		"user_id":          userID,
		"current_duration": entry.Duration,
		"paused_at":        now,
	}).WithAggregateID(entry.ID)

	err = s.eventBus.Publish("time.session.paused", event)
	if err != nil {
		return nil, fmt.Errorf("failed to publish pause event: %w", err)
	}

	log.Printf("⏸️  Timer paused: %s (Current duration: %d seconds)", userID, entry.Duration)
	return entry, nil
}

// ResumeTimer resumes a paused timer
func (s *Service) ResumeTimer(userID string) (*types.TimeEntry, error) {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active timer: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("no active timer for user %s", userID)
	}

	if entry.IsRunning {
		return entry, nil // Already running
	}

	entry.StartTime = time.Now() // Reset start time
	entry.IsRunning = true
	entry.UpdatedAt = time.Now()

	// Save updated entry
	err = s.repo.Save(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to save resumed timer: %w", err)
	}

	// Publish event
	event := types.NewEvent("session_resumed", "time_service", map[string]any{
		"time_entry_id":  entry.ID,
		"user_id":        userID,
		"resumed_at":     entry.StartTime,
		"total_duration": entry.Duration,
	}).WithAggregateID(entry.ID)

	err = s.eventBus.Publish("time.session.resumed", event)
	if err != nil {
		return nil, fmt.Errorf("failed to publish resume event: %w", err)
	}

	log.Printf("▶️  Timer resumed: %s", userID)
	return entry, nil
}

// GetActiveTimer returns the active timer for a user
func (s *Service) GetActiveTimer(userID string) *types.TimeEntry {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		log.Printf("Error getting active timer: %v", err)
		return nil
	}
	return entry
}

// GetCurrentDuration returns the current duration of an active timer
func (s *Service) GetCurrentDuration(userID string) int64 {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil || entry == nil || !entry.IsRunning {
		return 0
	}

	elapsed := int64(time.Since(entry.StartTime).Seconds())
	return entry.Duration + elapsed
}

// stopActiveTimer stops any active timer for a user
func (s *Service) stopActiveTimer(userID string) error {
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		return err
	}
	if entry != nil && entry.IsRunning {
		_, err = s.StopTimer(userID)
		return err
	}
	return nil
}

// Event Handlers

func (s *Service) handleProjectStarted(event *types.Event) error {
	projectID, _ := event.Data["project_id"].(string)
	userID, _ := event.Data["user_id"].(string)

	log.Printf("🎯 Project started event received: %s for user %s", projectID, userID)

	// Auto-start timer suggestion event
	suggestionEvent := types.NewEvent("timer_suggestion", "time_service", map[string]any{
		"user_id":    userID,
		"project_id": projectID,
		"suggestion": "start_timer",
		"reason":     "project_started",
	})

	return s.eventBus.Publish("time.suggestion.start", suggestionEvent)
}

func (s *Service) handleInvoiceGenerated(event *types.Event) error {
	invoiceID, _ := event.Data["invoice_id"].(string)
	userID, _ := event.Data["user_id"].(string)

	log.Printf("🧾 Invoice generated: %s - marking time entries as billed", invoiceID)

	// Mark relevant time entries as billed
	billedEvent := types.NewEvent("time_entries_billed", "time_service", map[string]any{
		"invoice_id": invoiceID,
		"user_id":    userID,
		"billed_at":  time.Now(),
	})

	return s.eventBus.Publish("time.entries.billed", billedEvent)
}

func (s *Service) handleUserLogout(event *types.Event) error {
	userID, _ := event.Data["user_id"].(string)

	log.Printf("👋 User logout: %s - stopping active timer", userID)

	// Stop any active timer
	entry, err := s.repo.GetActiveTimer(userID)
	if err != nil {
		log.Printf("Error getting active timer on logout: %v", err)
		return nil
	}
	if entry != nil {
		_, err := s.StopTimer(userID)
		if err != nil {
			log.Printf("Error stopping timer on logout: %v", err)
		}
	}

	return nil
}

// UpdateTimeEntry updates a time entry description or tags
func (s *Service) UpdateTimeEntry(userID, timeEntryID, description string, tags []string) error {
	entry, err := s.repo.Get(timeEntryID)
	if err != nil {
		return fmt.Errorf("failed to get time entry: %w", err)
	}
	if entry == nil || entry.UserID != userID {
		return fmt.Errorf("time entry not found")
	}

	entry.Description = description
	entry.Tags = tags
	entry.UpdatedAt = time.Now()

	// Save updated entry
	err = s.repo.Save(entry)
	if err != nil {
		return fmt.Errorf("failed to save updated time entry: %w", err)
	}

	// Publish update event
	event := types.NewEvent("session_updated", "time_service", map[string]any{
		"time_entry_id": timeEntryID,
		"user_id":       userID,
		"description":   description,
		"tags":          tags,
		"updated_at":    entry.UpdatedAt,
	}).WithAggregateID(timeEntryID)

	return s.eventBus.Publish("time.session.updated", event)
}
