package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"time"
)

// ContractTemplateCreatedEvent is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type ContractTemplateCreatedEvent struct {
	DID         string    `json:"did"`
	CreatedBy   string    `json:"created_by"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateCreatedEvent) EventType() string {
	return "TemplateCreatedEvent"
}

// GetDID implements the Event interface.
func (e ContractTemplateCreatedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateSubmittedEvent is emitted when a template is submitted for review.
// This event signals state transition and includes reviewer comments.
type ContractTemplateSubmittedEvent struct {
	DID            string                       `json:"did"`
	PreviousState  template_state.TemplateState `json:"previous_state"`
	NewState       template_state.TemplateState `json:"new_state"`
	SubmittedBy    string                       `json:"submitted_by"`
	ActionFlag     *action_flag.ActionFlag      `json:"action_flag"`
	ReviewComments []string                     `json:"review_comments,omitempty"`
	OccurredAt     time.Time                    `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateSubmittedEvent) EventType() string {
	return "TemplateSubmittedEvent"
}

// GetDID implements the Event interface.
func (e ContractTemplateSubmittedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateApprovedEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type ContractTemplateApprovedEvent struct {
	DID           string    `json:"did"`
	ApprovedBy    string    `json:"approved_by"`
	DecisionNotes []string  `json:"decision_notes,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateApprovedEvent) EventType() string {
	return "TemplateApprovedEvent"
}

// GetDID implements the Event interface.
func (e ContractTemplateApprovedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateRejectedEvent is emitted when a template is rejected.
// This event includes rejection reason and rejector information.
type ContractTemplateRejectedEvent struct {
	DID        string    `json:"did"`
	RejectedBy string    `json:"rejected_by"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRejectedEvent) EventType() string {
	return "TemplateRejectedEvent"
}

// GetDID implements the Event interface.
func (e ContractTemplateRejectedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateUpdatedEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateUpdatedEvent struct {
	DID            string         `json:"did"`
	UpdatedBy      string         `json:"updated_by"`
	OldName        *string        `json:"old_name,omitempty"`
	NewName        *string        `json:"new_name,omitempty"`
	OldDescription *string        `json:"old_description,omitempty"`
	NewDescription *string        `json:"new_description,omitempty"`
	OldMetaData    *datatype.JSON `json:"old_meta_data,omitempty"`
	NewMetaData    *datatype.JSON `json:"new_metadata,omitempty"`
	OccurredAt     time.Time      `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateUpdatedEvent) EventType() string {
	return "TemplateUpdatedEvent"
}

// GetDID implements the Event interface.
func (e ContractTemplateUpdatedEvent) GetDID() string {
	return e.DID
}
