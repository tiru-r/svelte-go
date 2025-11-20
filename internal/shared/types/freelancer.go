package types

import (
	"time"

	"github.com/google/uuid"
)

// User represents a freelancer user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Client represents a client/customer
type Client struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Company    string    `json:"company"`
	Phone      string    `json:"phone"`
	HourlyRate float64   `json:"hourly_rate"`
	Currency   string    `json:"currency"`
	Address    string    `json:"address"`
	Notes      string    `json:"notes"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Project represents a project for a client
type Project struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ClientID    string     `json:"client_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	HourlyRate  float64    `json:"hourly_rate"`
	Currency    string     `json:"currency"`
	Status      string     `json:"status"` // active, paused, completed
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TimeEntry represents a time tracking entry
type TimeEntry struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ProjectID   string     `json:"project_id"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Duration    int64      `json:"duration"` // seconds
	IsRunning   bool       `json:"is_running"`
	IsBilled    bool       `json:"is_billed"`
	HourlyRate  float64    `json:"hourly_rate"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Expense represents an expense entry
type Expense struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProjectID   string    `json:"project_id,omitempty"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Receipt     string    `json:"receipt,omitempty"` // file path/URL
	IsBillable  bool      `json:"is_billable"`
	IsBilled    bool      `json:"is_billed"`
	TaxCategory string    `json:"tax_category"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
}

// Invoice represents an invoice
type Invoice struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	ClientID    string        `json:"client_id"`
	ProjectID   string        `json:"project_id,omitempty"`
	Number      string        `json:"number"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"` // draft, sent, paid, overdue
	Items       []InvoiceItem `json:"items"`
	Amount      float64       `json:"amount"`
	Currency    string        `json:"currency"`
	TaxRate     float64       `json:"tax_rate"`
	TaxAmount   float64       `json:"tax_amount"`
	TotalAmount float64       `json:"total_amount"`
	IssueDate   time.Time     `json:"issue_date"`
	DueDate     time.Time     `json:"due_date"`
	SentAt      *time.Time    `json:"sent_at,omitempty"`
	PaidAt      *time.Time    `json:"paid_at,omitempty"`
	TimeEntries []TimeEntry   `json:"time_entries,omitempty"`
	Expenses    []Expense     `json:"expenses,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewTimeEntry creates a new time entry with generated ID
func NewTimeEntry(userID, projectID, description string) *TimeEntry {
	return &TimeEntry{
		ID:          uuid.New().String(),
		UserID:      userID,
		ProjectID:   projectID,
		Description: description,
		StartTime:   time.Now(),
		IsRunning:   true,
		IsBilled:    false,
		Tags:        make([]string, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Stop stops the time entry and calculates duration
func (te *TimeEntry) Stop() {
	if te.IsRunning {
		now := time.Now()
		te.EndTime = &now
		te.Duration = int64(now.Sub(te.StartTime).Seconds())
		te.IsRunning = false
		te.UpdatedAt = now
	}
}

// CalculateAmount calculates the billable amount for the time entry
func (te *TimeEntry) CalculateAmount() float64 {
	if te.Duration > 0 && te.HourlyRate > 0 {
		hours := float64(te.Duration) / 3600.0 // convert seconds to hours
		return hours * te.HourlyRate
	}
	return 0
}

// GenerateID generates a new UUID string
func GenerateID() string {
	return uuid.New().String()
}
