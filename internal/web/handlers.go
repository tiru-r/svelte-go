package web

import (
	"encoding/json"
	"net/http"

	"datastar-go/internal/modules/auth"
	"datastar-go/internal/modules/client"
	"datastar-go/internal/modules/expense"
	"datastar-go/internal/modules/invoice"
	timemodule "datastar-go/internal/modules/time"
	"datastar-go/templates"
)

type Handlers struct {
	clientService  *client.Service
	expenseService *expense.Service
	invoiceService *invoice.Service
	timeService    *timemodule.Service
}

func NewHandlers(
	clientService *client.Service,
	expenseService *expense.Service,
	invoiceService *invoice.Service,
	timeService *timemodule.Service,
) *Handlers {
	return &Handlers{
		clientService:  clientService,
		expenseService: expenseService,
		invoiceService: invoiceService,
		timeService:    timeService,
	}
}

func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	// Main application page routes  
	mux.HandleFunc("/", h.Dashboard)
	mux.HandleFunc("/clients", h.Clients)
	mux.HandleFunc("/invoices", h.Invoices)
	mux.HandleFunc("/expenses", h.Expenses)
	mux.HandleFunc("/timer", h.Timer)

	// API routes for Datastar
	mux.HandleFunc("/api/dashboard/stats", h.DashboardStats)
	
	// Client API routes (frontend compatibility)
	mux.HandleFunc("/api/clients", h.GetClients)
	mux.HandleFunc("/api/clients/form", h.ClientForm)
	
	// Invoice API routes (frontend compatibility)
	mux.HandleFunc("/api/invoices", h.GetInvoices)
	mux.HandleFunc("/api/invoices/form", h.InvoiceForm)
	
	// Expense API routes (frontend compatibility)
	mux.HandleFunc("/api/expenses", h.GetExpenses)
	mux.HandleFunc("/api/expenses/form", h.ExpenseForm)
	
	// Timer API routes (frontend compatibility)
	mux.HandleFunc("/api/timer/entries", h.GetTimeEntries)
	mux.HandleFunc("/api/timer/start", h.StartTimer)
	mux.HandleFunc("/api/timer/stop", h.StopTimer)
	mux.HandleFunc("/api/timer/pause", h.PauseTimer)
}

func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Dashboard().Render(r.Context(), w)
}

func (h *Handlers) Clients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Clients().Render(r.Context(), w)
}

func (h *Handlers) Invoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Invoices().Render(r.Context(), w)
}

func (h *Handlers) Expenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Expenses().Render(r.Context(), w)
}

func (h *Handlers) Timer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.Timer().Render(r.Context(), w)
}


// DashboardStats provides statistics for the dashboard
func (h *Handlers) DashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	// Get actual counts from services
	clients, err := h.clientService.GetClients(userID)
	clientsCount := 0
	if err == nil {
		clientsCount = len(clients)
	}
	
	invoices, err := h.invoiceService.GetInvoices(userID)
	invoicesCount := 0
	if err == nil {
		invoicesCount = len(invoices)
	}
	
	expenses, err := h.expenseService.GetExpenses(userID)
	expensesCount := 0
	if err == nil {
		expensesCount = len(expenses)
	}
	
	// Check for active timers
	activeTimers := 0
	if h.timeService.GetActiveTimer(userID) != nil {
		activeTimers = 1
	}

	// Create response data
	stats := map[string]interface{}{
		"clientsCount":  clientsCount,
		"invoicesCount": invoicesCount,
		"expensesCount": expensesCount,
		"activeTimers":  activeTimers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetClients returns clients data for frontend
func (h *Handlers) GetClients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	clients, err := h.clientService.GetClients(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clients": clients,
	})
}

// ClientForm returns form HTML for adding/editing clients
func (h *Handlers) ClientForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return simple form for now
	form := `<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
		<div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
			<h3 class="text-lg font-bold text-gray-900 mb-4">Add Client</h3>
			<form data-on-submit="$$post('/api/client/create')">
				<input type="text" name="name" placeholder="Client Name" class="w-full p-2 border rounded mb-4" required>
				<input type="email" name="email" placeholder="Email" class="w-full p-2 border rounded mb-4" required>
				<input type="tel" name="phone" placeholder="Phone" class="w-full p-2 border rounded mb-4">
				<div class="flex justify-end">
					<button type="button" class="mr-2 px-4 py-2 bg-gray-300 rounded" data-on-click="$$remove('#modal-container')">
						Cancel
					</button>
					<button type="submit" class="px-4 py-2 bg-blue-500 text-white rounded">
						Save
					</button>
				</div>
			</form>
		</div>
	</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(form))
}

// GetInvoices returns invoices data for frontend
func (h *Handlers) GetInvoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	invoices, err := h.invoiceService.GetInvoices(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"invoices": invoices,
	})
}

// InvoiceForm returns form HTML for adding/editing invoices
func (h *Handlers) InvoiceForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	form := `<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
		<div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
			<h3 class="text-lg font-bold text-gray-900 mb-4">Create Invoice</h3>
			<form data-on-submit="$$post('/api/invoice/create')">
				<input type="text" name="client_id" placeholder="Client ID" class="w-full p-2 border rounded mb-4" required>
				<input type="number" name="amount" placeholder="Amount" step="0.01" class="w-full p-2 border rounded mb-4" required>
				<textarea name="description" placeholder="Description" class="w-full p-2 border rounded mb-4" rows="3"></textarea>
				<div class="flex justify-end">
					<button type="button" class="mr-2 px-4 py-2 bg-gray-300 rounded" data-on-click="$$remove('#modal-container')">
						Cancel
					</button>
					<button type="submit" class="px-4 py-2 bg-blue-500 text-white rounded">
						Create
					</button>
				</div>
			</form>
		</div>
	</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(form))
}

// GetExpenses returns expenses data for frontend
func (h *Handlers) GetExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	expenses, err := h.expenseService.GetExpenses(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expenses": expenses,
	})
}

// ExpenseForm returns form HTML for adding/editing expenses
func (h *Handlers) ExpenseForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	form := `<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
		<div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
			<h3 class="text-lg font-bold text-gray-900 mb-4">Add Expense</h3>
			<form data-on-submit="$$post('/api/expense/create')">
				<input type="text" name="description" placeholder="Description" class="w-full p-2 border rounded mb-4" required>
				<input type="text" name="category" placeholder="Category" class="w-full p-2 border rounded mb-4" required>
				<input type="number" name="amount" placeholder="Amount" step="0.01" class="w-full p-2 border rounded mb-4" required>
				<div class="flex justify-end">
					<button type="button" class="mr-2 px-4 py-2 bg-gray-300 rounded" data-on-click="$$remove('#modal-container')">
						Cancel
					</button>
					<button type="submit" class="px-4 py-2 bg-blue-500 text-white rounded">
						Save
					</button>
				</div>
			</form>
		</div>
	</div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(form))
}

// GetTimeEntries returns time entries for frontend
func (h *Handlers) GetTimeEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	// TODO: Implement GetTimeEntries method in time service
	// For now, return empty array
	entries := []interface{}{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"time_entries": entries,
	})
}

// StartTimer starts a timer for frontend
func (h *Handlers) StartTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	projectID := "default-project"
	description := "Work session"
	
	entry, err := h.timeService.StartTimer(userID, projectID, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"timer_running": true,
		"entry": entry,
	})
}

// StopTimer stops a timer for frontend
func (h *Handlers) StopTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	entry, err := h.timeService.StopTimer(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"timer_running": false,
		"entry": entry,
	})
}

// PauseTimer pauses a timer for frontend
func (h *Handlers) PauseTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	entry, err := h.timeService.PauseTimer(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"timer_running": false,
		"entry": entry,
	})
}