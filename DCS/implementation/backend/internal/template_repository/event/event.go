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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateCreatedEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateCreatedEvent) GetVersion() int {
	return e.Version
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

// EventType implements the Event interface.
func (e ContractTemplateSubmittedEvent) EventType() string {
	return event_type.SubmitContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateSubmittedEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateSubmittedEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateSubmittedEvent) GetVersion() int {
	return e.Version
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateApprovedEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateApprovedEvent) GetVersion() int {
	return e.Version
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateRejectedEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateRejectedEvent) GetVersion() int {
	return e.Version
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateUpdatedEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateUpdatedEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateRetrievedAllEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateRetrievedAllEvent struct {
	RetrievedBy    string    `json:"updated_by"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) EventType() string {
	return event_type.RetrieveAllContractTemplates.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) GetDID() string {
	return "*"
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateRetrievedByIdEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateRetrievedByIdEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateRetrievedByIdEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateRetrievedByIdEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateCreateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateCreateReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Reviewer       string    `json:"reviewer"`
	CreatedBy      string    `json:"created_by"`
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateCreateReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateCreateReviewTaskEvent) GetVersion() int {
	return e.Version
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
	return event_type.RetrieveAllContractTemplateReviewTasks.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrieveAllReviewTasksEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateRetrieveAllReviewTasksEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateRetrieveAllReviewTasksEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateReopenReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateReopenReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	CreatedBy      string    `json:"created_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateReopenReviewTaskEvent) EventType() string {
	return event_type.ReopenContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateReopenReviewTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateReopenReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateReopenReviewTaskEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateUpdatedApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateUpdatedReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	UpdatedBy      string    `json:"updated_by"`
	OldState       *string   `json:"old_name,omitempty"`
	NewState       *string   `json:"new_name,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateUpdatedReviewTaskEvent) EventType() string {
	return event_type.UpdateContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateUpdatedReviewTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateUpdatedReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateUpdatedReviewTaskEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateCreateApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateCreateApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Approver       string    `json:"approver"`
	CreatedBy      string    `json:"created_by"`
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

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateCreateApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateCreateApprovalTaskEvent) GetVersion() int {
	return e.Version
}

type ContractTemplateRetrieveAllApprovalTasksEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrieveAllApprovalTasksEvent) EventType() string {
	return event_type.RetrieveAllContractTemplateApprovalTasks.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateRetrieveAllApprovalTasksEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateRetrieveAllApprovalTasksEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateRetrieveAllApprovalTasksEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateReopenApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateReopenApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	CreatedBy      string    `json:"created_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateReopenApprovalTaskEvent) EventType() string {
	return event_type.ReopenContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateReopenApprovalTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateReopenApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateReopenApprovalTaskEvent) GetVersion() int {
	return e.Version
}

// ContractTemplateUpdatedApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateUpdatedApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	UpdatedBy      string    `json:"updated_by"`
	OldState       *string   `json:"old_name,omitempty"`
	NewState       *string   `json:"new_name,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateUpdatedApprovalTaskEvent) EventType() string {
	return event_type.UpdateContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateUpdatedApprovalTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateUpdatedApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateUpdatedApprovalTaskEvent) GetVersion() int {
	return e.Version
}
