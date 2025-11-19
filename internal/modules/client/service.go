package client

import (
	"fmt"
	"log"
	"time"

	"svelte-go/internal/shared/database"
	"svelte-go/internal/shared/types"
)

type Service struct {
	eventBus types.EventBus
	db       *database.BadgerDB
}

func NewService(eventBus types.EventBus, db *database.BadgerDB) *Service {
	service := &Service{
		eventBus: eventBus,
		db:       db,
	}

	service.setupEventSubscriptions()
	return service
}

func (s *Service) setupEventSubscriptions() {
	s.eventBus.SubscribeQueue("invoice.generated", "client_service", s.handleInvoiceGenerated)
	s.eventBus.SubscribeQueue("time.entry.completed", "client_service", s.handleTimeEntryCompleted)
	
	log.Println("Client service event subscriptions configured")
}

func (s *Service) CreateClient(userID, name, email, company string) (*types.Client, error) {
	client := &types.Client{
		ID:        types.GenerateID(),
		UserID:    userID,
		Name:      name,
		Email:     email,
		Company:   company,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.db.ClientRepo.Create(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	event := types.NewEvent("client_created", "client_service", map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"name":      client.Name,
		"email":     client.Email,
		"company":   client.Company,
	})
	
	s.eventBus.Publish("client.created", event)
	log.Printf("👤 Client created: %s (%s)", client.Name, client.Company)
	
	return client, nil
}

func (s *Service) CreateProject(clientID, userID, name, description string, hourlyRate float64) (*types.Project, error) {
	project := &types.Project{
		ID:          types.GenerateID(),
		ClientID:    clientID,
		UserID:      userID,
		Name:        name,
		Description: description,
		HourlyRate:  hourlyRate,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.db.ProjectRepo.Create(project)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	event := types.NewEvent("project_created", "client_service", map[string]interface{}{
		"project_id":   project.ID,
		"client_id":    project.ClientID,
		"user_id":      project.UserID,
		"name":         project.Name,
		"hourly_rate":  project.HourlyRate,
	})
	
	s.eventBus.Publish("client.project.started", event)
	log.Printf("📋 Project created: %s (Rate: $%.2f/hr)", project.Name, project.HourlyRate)
	
	return project, nil
}

func (s *Service) GetClients(userID string) ([]*types.Client, error) {
	return s.db.ClientRepo.GetByUserID(userID)
}

func (s *Service) GetProjects(userID string) ([]*types.Project, error) {
	return s.db.ProjectRepo.GetByUserID(userID)
}

func (s *Service) GetProjectsByClient(clientID string) ([]*types.Project, error) {
	return s.db.ProjectRepo.GetByClientID(clientID)
}

func (s *Service) UpdateClient(clientID string, updates map[string]interface{}) (*types.Client, error) {
	client, err := s.db.ClientRepo.GetByID(clientID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	if name, ok := updates["name"].(string); ok {
		client.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		client.Email = email
	}
	if company, ok := updates["company"].(string); ok {
		client.Company = company
	}
	if phone, ok := updates["phone"].(string); ok {
		client.Phone = phone
	}
	if address, ok := updates["address"].(string); ok {
		client.Address = address
	}
	
	client.UpdatedAt = time.Now()

	err = s.db.ClientRepo.Update(client)
	if err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	event := types.NewEvent("client_updated", "client_service", map[string]interface{}{
		"client_id": client.ID,
		"user_id":   client.UserID,
		"updates":   updates,
	})
	
	s.eventBus.Publish("client.updated", event)
	return client, nil
}

func (s *Service) UpdateProject(projectID string, updates map[string]interface{}) (*types.Project, error) {
	project, err := s.db.ProjectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if name, ok := updates["name"].(string); ok {
		project.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		project.Description = description
	}
	if hourlyRate, ok := updates["hourly_rate"].(float64); ok {
		project.HourlyRate = hourlyRate
	}
	if status, ok := updates["status"].(string); ok {
		project.Status = status
	}
	
	project.UpdatedAt = time.Now()

	err = s.db.ProjectRepo.Update(project)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	event := types.NewEvent("project_updated", "client_service", map[string]interface{}{
		"project_id": project.ID,
		"user_id":    project.UserID,
		"updates":    updates,
	})
	
	s.eventBus.Publish("client.project.updated", event)
	return project, nil
}

func (s *Service) handleInvoiceGenerated(event *types.Event) error {
	log.Printf("📧 Invoice generated for client: %v", event.Data["client_id"])
	return nil
}

func (s *Service) handleTimeEntryCompleted(event *types.Event) error {
	log.Printf("⏱️ Time tracked for project: %v", event.Data["project_id"])
	return nil
}