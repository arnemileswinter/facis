package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/template_repository/datatype/action_flag"
	"digital-contracting-service/internal/template_repository/datatype/event_type"
	"digital-contracting-service/internal/template_repository/datatype/template_state"
	"time"
)

// ContractTemplateCreatedEvent is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type ContractTemplateCreatedEvent struct {
	DID            string         `json:"did"`
	DocumentNumber int            `json:"document_number"`
	Version        int            `json:"version"`
	CreatedBy      string         `json:"created_by"`
	Name           *string        `json:"name"`
	Description    *string        `json:"description"`
	MetaData       *datatype.JSON `json:"metadata"`
	OccurredAt     time.Time      `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateCreatedEvent) EventType() string {
	return event_type.CreateContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateCreatedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateSubmittedEvent is emitted when a template is submitted for review.
// This event signals state transition and includes reviewer comments.
type ContractTemplateSubmittedEvent struct {
	DID            string                       `json:"did"`
	DocumentNumber int                          `json:"document_number"`
	Version        int                          `json:"version"`
	PreviousState  template_state.TemplateState `json:"previous_state"`
	NewState       template_state.TemplateState `json:"new_state"`
	SubmittedBy    string                       `json:"submitted_by"`
	ActionFlag     *action_flag.ActionFlag      `json:"action_flag"`
	Comments       []string                     `json:"comments,omitempty"`
	OccurredAt     time.Time                    `json:"occurred_at"`
}

// GetDID implements the Event interface.
func (e ContractTemplateSubmittedEvent) GetDID() string {
	return e.DID
}

// EventType implements the Event interface.
func (e ContractTemplateSubmittedEvent) EventType() string {
	return event_type.SubmitContractTemplate.String()
}

// ContractTemplateApprovedEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type ContractTemplateApprovedEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	ApprovedBy     string    `json:"approved_by"`
	DecisionNotes  []string  `json:"decision_notes,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateApprovedEvent) EventType() string {
	return event_type.ApproveContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateApprovedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateRejectedEvent is emitted when a template is rejected.
// This event includes rejection reason and rejector information.
type ContractTemplateRejectedEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RejectedBy     string    `json:"rejected_by"`
	Reason         string    `json:"reason"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRejectedEvent) EventType() string {
	return event_type.RejectContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRejectedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateUpdatedEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateUpdatedEvent struct {
	DID            string         `json:"did"`
	DocumentNumber int            `json:"document_number"`
	Version        int            `json:"version"`
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
	return event_type.UpdateContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateUpdatedEvent) GetDID() string {
	return e.DID
}

// ContractTemplateRetrievedAllEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateRetrievedAllEvent struct {
	RetrievedBy string    `json:"updated_by"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) EventType() string {
	return event_type.RetrieveAllContractTemplates.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) GetDID() string {
	return "*"
}

// ContractTemplateRetrievedByIdEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateRetrievedByIdEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"updated_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrievedByIdEvent) EventType() string {
	return event_type.RetrieveContractTemplateById.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrievedByIdEvent) GetDID() string {
	return e.DID
}

// ContractTemplateCreateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateCreateReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Reviewer       string    `json:"reviewer"`
	CreatedBy      string    `json:"updated_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateCreateReviewTaskEvent) EventType() string {
	return event_type.CreateContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateCreateReviewTaskEvent) GetDID() string {
	return e.DID
}

// ContractTemplateCreateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateRetrieveAllReviewTasksEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrieveAllReviewTasksEvent) EventType() string {
	return event_type.CreateContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrieveAllReviewTasksEvent) GetDID() string {
	return e.DID
}

// ContractTemplateCreateApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateCreateApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Approver       string    `json:"approver"`
	CreatedBy      string    `json:"updated_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateCreateApprovalTaskEvent) EventType() string {
	return event_type.CreateContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateCreateApprovalTaskEvent) GetDID() string {
	return e.DID
}
