package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Event interface must be implemented by all event types.
// This interface ensures consistency across event handling.
type Event interface {
	// EventType returns the name of the event (used as NATS subject).
	EventType() string

	// GetDID returns the entity DID for event reference and correlation.
	GetDID() string

	// GetDocumentNumber returns the entity DocumentNumber for event reference and correlation.
	GetDocumentNumber() int

	// GetVersion returns the entity Version for event reference and correlation.
	GetVersion() int
}

// CreateNewEvent persists an event to the outbox table.
// This function is called by command handlers to store events durably
// before they are published to NATS.
//
// IMPORTANT: Must be called within the same database transaction as the
// command that triggered the event. This ensures atomicity: either both
// the command and event succeed, or both are rolled back.
//
// Usage in a command handler:
//
//	evt := TemplateCreatedEvent{
//	    TemplateID: "123",
//	    CreatedBy:  "user@example.com",
//	    OccurredAt: time.Now(),
//	}
//	if err := event.CreateNewEvent(ctx, tx, evt); err != nil {
//	    return err
//	}
func CreateNewEvent(ctx context.Context, tx *sqlx.Tx, evt Event) error {
	if evt == nil {
		return errors.New("event cannot be nil")
	}

	eventType := evt.EventType()
	if eventType == "" {
		return errors.New("event type cannot be empty")
	}

	did := evt.GetDID()
	if did == "" {
		return errors.New("template did cannot be empty")
	}

	documentNumber := evt.GetDocumentNumber()
	version := evt.GetVersion()

	// Serialize event to JSON
	eventJSON, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Insert into outbox (MUST be in same transaction as DB update!)
	// The outbox table ensures events are never lost, even if NATS is down.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events 
		 (event_type, event_data, did, document_number, version, processed)
		 VALUES ($1, $2, $3, $4, $5, FALSE)`,
		eventType,
		eventJSON,
		did,
		documentNumber,
		version,
	)

	if err != nil {
		return fmt.Errorf("failed to insert event into outbox: %w", err)
	}

	return nil
}

// CreateNewEvents persists multiple events to the outbox in a single operation.
// This is useful when a single command handler needs to emit multiple events.
//
// Events are inserted sequentially within the same transaction.
//
// Usage:
//
//	evts := []event.Event{
//	    TemplateCreatedEvent{...},
//	    TemplateInitializedEvent{...},
//	    TemplateReadyEvent{...},
//	}
//	if err := event.CreateNewEvents(ctx, tx, evts...); err != nil {
//	    return err
//	}
func CreateNewEvents(ctx context.Context, tx *sqlx.Tx, events ...Event) error {
	if len(events) == 0 {
		return nil // Nothing to store
	}

	for _, evt := range events {
		if err := CreateNewEvent(ctx, tx, evt); err != nil {
			return err
		}
	}

	return nil
}
