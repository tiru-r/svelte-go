package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type NATSService struct {
	server *server.Server
	conn   *nats.Conn
	js     nats.JetStreamContext
}

type Event struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

func NewNATSService() (*NATSService, error) {
	// Configure embedded NATS server
	opts := &server.Options{
		Port:      4223, // Use different port to avoid conflicts
		Host:      "127.0.0.1",
		JetStream: true,
		StoreDir:  "./data/jetstream",
	}

	// Start embedded NATS server
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	// Start the server in a goroutine
	go ns.Start()

	// Wait for server to be ready
	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, fmt.Errorf("NATS server failed to start")
	}

	// Connect to the embedded server
	conn, err := nats.Connect("nats://127.0.0.1:4223")
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

	// Create default streams
	service := &NATSService{
		server: ns,
		conn:   conn,
		js:     js,
	}

	if err := service.setupStreams(); err != nil {
		service.Close()
		return nil, fmt.Errorf("failed to setup streams: %w", err)
	}

	log.Println("NATS embedded server started on port 4223")
	return service, nil
}

func (ns *NATSService) setupStreams() error {
	// Create application events stream
	_, err := ns.js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.>"},
		Storage:  nats.FileStorage,
		MaxAge:   time.Hour * 24, // Keep events for 24 hours
	})
	if err != nil {
		return fmt.Errorf("failed to create EVENTS stream: %w", err)
	}

	// Create user actions stream
	_, err = ns.js.AddStream(&nats.StreamConfig{
		Name:     "USER_ACTIONS",
		Subjects: []string{"user.>"},
		Storage:  nats.FileStorage,
		MaxAge:   time.Hour * 24 * 7, // Keep user actions for 7 days
	})
	if err != nil {
		return fmt.Errorf("failed to create USER_ACTIONS stream: %w", err)
	}

	return nil
}

func (ns *NATSService) PublishEvent(subject string, event Event) error {
	event.Timestamp = time.Now()
	
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = ns.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event: %s to subject: %s", event.Type, subject)
	return nil
}

func (ns *NATSService) Subscribe(subject string, handler func(*Event) error) (*nats.Subscription, error) {
	sub, err := ns.js.Subscribe(subject, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			msg.Nak()
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Event handler failed: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	log.Printf("Subscribed to subject: %s", subject)
	return sub, nil
}

func (ns *NATSService) SubscribeQueue(subject, queue string, handler func(*Event) error) (*nats.Subscription, error) {
	sub, err := ns.js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			msg.Nak()
			return
		}

		if err := handler(&event); err != nil {
			log.Printf("Event handler failed: %v", err)
			msg.Nak()
			return
		}

		msg.Ack()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to queue subscribe to %s: %w", subject, err)
	}

	log.Printf("Queue subscribed to subject: %s (queue: %s)", subject, queue)
	return sub, nil
}

func (ns *NATSService) Close() error {
	if ns.conn != nil {
		ns.conn.Close()
	}
	if ns.server != nil {
		ns.server.Shutdown()
		ns.server.WaitForShutdown()
	}
	log.Println("NATS service shut down")
	return nil
}

func (ns *NATSService) GetStats() map[string]interface{} {
	if ns.server != nil {
		varz, _ := ns.server.Varz(nil)
		return map[string]interface{}{
			"connections":    varz.Connections,
			"subscriptions":  varz.Subscriptions,
			"in_msgs":        varz.InMsgs,
			"out_msgs":       varz.OutMsgs,
			"in_bytes":       varz.InBytes,
			"out_bytes":      varz.OutBytes,
			"uptime":         varz.Now.Sub(varz.Start),
		}
	}
	return nil
}