package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"datastar-go/internal/shared/types"
)

// EventBus implements the types.EventBus interface using NATS
type EventBus struct {
	server *server.Server
	conn   *nats.Conn
	js     nats.JetStreamContext
}

// NewEventBus creates a new NATS-based event bus
func NewEventBus() (*EventBus, error) {
	// Try different ports to avoid conflicts
	for port := 4223; port <= 4230; port++ {
		bus, err := NewEventBusWithPort(port)
		if err == nil {
			return bus, nil
		}
		log.Printf("Port %d unavailable, trying next...", port)
	}
	return nil, fmt.Errorf("no available ports for NATS server")
}

// NewEventBusWithPort creates a new NATS-based event bus on specific port
func NewEventBusWithPort(port int) (*EventBus, error) {
	// Configure embedded NATS server
	opts := &server.Options{
		Port:      port,
		Host:      "127.0.0.1",
		JetStream: true,
		StoreDir:  fmt.Sprintf("./data/jetstream_%d", port),
	}

	// Start embedded NATS server
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, fmt.Errorf("NATS server failed to start")
	}

	// Connect to the embedded server
	conn, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", port))
	if err != nil {
		ns.Shutdown()
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		ns.Shutdown()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	bus := &EventBus{
		server: ns,
		conn:   conn,
		js:     js,
	}

	if err := bus.setupStreams(); err != nil {
		bus.Close()
		return nil, fmt.Errorf("failed to setup streams: %w", err)
	}

	log.Printf("Event bus initialized with NATS on port %d", port)
	return bus, nil
}

// setupStreams creates the event streams for different domains
func (eb *EventBus) setupStreams() error {
	streams := []struct {
		name     string
		subjects []string
		maxAge   time.Duration
	}{
		{
			name:     "TIME_EVENTS",
			subjects: []string{"time.>"},
			maxAge:   time.Hour * 24 * 30, // 30 days
		},
		{
			name:     "EXPENSE_EVENTS",
			subjects: []string{"expense.>"},
			maxAge:   time.Hour * 24 * 90, // 90 days for tax records
		},
		{
			name:     "INVOICE_EVENTS",
			subjects: []string{"invoice.>"},
			maxAge:   time.Hour * 24 * 365, // 1 year for financial records
		},
		{
			name:     "CLIENT_EVENTS",
			subjects: []string{"client.>"},
			maxAge:   time.Hour * 24 * 365, // 1 year
		},
		{
			name:     "ANALYTICS_EVENTS",
			subjects: []string{"analytics.>"},
			maxAge:   time.Hour * 24 * 7, // 7 days
		},
		{
			name:     "SYSTEM_EVENTS",
			subjects: []string{"system.>", "auth.>"},
			maxAge:   time.Hour * 24, // 24 hours
		},
	}

	for _, stream := range streams {
		_, err := eb.js.AddStream(&nats.StreamConfig{
			Name:     stream.name,
			Subjects: stream.subjects,
			Storage:  nats.FileStorage,
			MaxAge:   stream.maxAge,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream %s: %w", stream.name, err)
		}
		log.Printf("Created event stream: %s", stream.name)
	}

	return nil
}

// Publish publishes an event to the specified subject
func (eb *EventBus) Publish(subject string, event *types.Event) error {
	data, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = eb.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	log.Printf("📤 Event published: %s -> %s (ID: %s)", event.Type, subject, event.ID)
	return nil
}

// Subscribe subscribes to events on the specified subject
func (eb *EventBus) Subscribe(subject string, handler types.EventHandler) error {
	_, err := eb.js.Subscribe(subject, func(msg *nats.Msg) {
		eb.handleMessage(msg, handler)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	log.Printf("📥 Subscribed to: %s", subject)
	return nil
}

// SubscribeQueue subscribes to events with a queue group for load balancing
func (eb *EventBus) SubscribeQueue(subject, queue string, handler types.EventHandler) error {
	_, err := eb.js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		eb.handleMessage(msg, handler)
	})

	if err != nil {
		return fmt.Errorf("failed to queue subscribe to %s: %w", subject, err)
	}

	log.Printf("📥 Queue subscribed to: %s (queue: %s)", subject, queue)
	return nil
}

// handleMessage processes incoming messages and calls the handler
func (eb *EventBus) handleMessage(msg *nats.Msg, handler types.EventHandler) {
	var event types.Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("❌ Failed to unmarshal event: %v", err)
		msg.Nak()
		return
	}

	log.Printf("📨 Processing event: %s (ID: %s) from %s", event.Type, event.ID, event.Source)

	if err := handler(&event); err != nil {
		log.Printf("❌ Event handler failed for %s: %v", event.Type, err)
		msg.Nak()
		return
	}

	msg.Ack()
	log.Printf("✅ Event processed: %s (ID: %s)", event.Type, event.ID)
}

// Close shuts down the event bus
func (eb *EventBus) Close() error {
	if eb.conn != nil {
		eb.conn.Close()
	}
	if eb.server != nil {
		eb.server.Shutdown()
		eb.server.WaitForShutdown()
	}
	log.Println("Event bus shut down")
	return nil
}

// GetStats returns NATS server statistics
func (eb *EventBus) GetStats() map[string]any {
	if eb.server != nil {
		varz, _ := eb.server.Varz(nil)
		return map[string]any{
			"connections":   varz.Connections,
			"subscriptions": varz.Subscriptions,
			"in_msgs":       varz.InMsgs,
			"out_msgs":      varz.OutMsgs,
			"in_bytes":      varz.InBytes,
			"out_bytes":     varz.OutBytes,
			"uptime":        varz.Now.Sub(varz.Start),
		}
	}
	return nil
}
