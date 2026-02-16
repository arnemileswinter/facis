package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/eventtype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"time"
)

// ContractTemplateCreatedEvent is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type ContractTemplateCreatedEvent struct {
	DID            string         `json:"did"`
	DocumentNumber int            `json:"document_number"`
	Version        int            `json:"version"`
	CreatedBy      string         `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Name           *string        `json:"name"`
	Description    *string        `json:"description"`
	MetaData       *datatype.JSON `json:"metadata"`
	OccurredAt     time.Time      `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateCreatedEvent) EventType() string {
	return eventtype.CreateContractTemplate.String()
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
	DID            string                      `json:"did"`
	DocumentNumber int                         `json:"document_number"`
	Version        int                         `json:"version"`
	PreviousState  templatestate.TemplateState `json:"previous_state"`
	NewState       templatestate.TemplateState `json:"new_state"`
	SubmittedBy    string                      `json:"submitted_by"`
	ActionFlag     *actionflag.ActionFlag      `json:"actionflag"`
	Comments       []string                    `json:"comments,omitempty"`
	OccurredAt     time.Time                   `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateSubmittedEvent) EventType() string {
	return eventtype.SubmitContractTemplate.String()
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
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int       `json:"version"`
	ApprovedBy     string    `json:"approved_by"`
	DecisionNotes  []string  `json:"decision_notes,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateApprovedEvent) EventType() string {
	return eventtype.ApproveContractTemplate.String()
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
	UpdatedAt      time.Time `json:"updated_at"`
	RejectedBy     string    `json:"rejected_by"`
	Reason         string    `json:"reason"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateRejectedEvent) EventType() string {
	return eventtype.RejectContractTemplate.String()
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
	return eventtype.UpdateContractTemplate.String()
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
	RetrievedBy    string                 `json:"updated_by"`
	DocumentNumber int                    `json:"document_number"`
	Version        int                    `json:"version"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Filter         map[string]interface{} `json:"filter"`
}

// EventType implements the Event interface.
func (e ContractTemplateRetrievedAllEvent) EventType() string {
	return eventtype.RetrieveAllContractTemplates.String()
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
	return eventtype.RetrieveContractTemplateById.String()
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
	return eventtype.CreateContractTemplateReviewTask.String()
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
	return eventtype.RetrieveAllContractTemplateReviewTasks.String()
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
	return eventtype.ReopenContractTemplateReviewTask.String()
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
	return eventtype.UpdateContractTemplateReviewTask.String()
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

// ContractTemplateDeleteReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateDeleteReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	DeletedBy      string    `json:"deleted_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateDeleteReviewTaskEvent) EventType() string {
	return eventtype.DeleteContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateDeleteReviewTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateDeleteReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateDeleteReviewTaskEvent) GetVersion() int {
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
	return eventtype.CreateContractTemplateApprovalTask.String()
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
	return eventtype.RetrieveAllContractTemplateApprovalTasks.String()
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
	return eventtype.ReopenContractTemplateApprovalTask.String()
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
	return eventtype.UpdateContractTemplateApprovalTask.String()
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

// ContractTemplateDeleteApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ContractTemplateDeleteApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	DeletedBy      string    `json:"deleted_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ContractTemplateDeleteApprovalTaskEvent) EventType() string {
	return eventtype.DeleteContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e ContractTemplateDeleteApprovalTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ContractTemplateDeleteApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ContractTemplateDeleteApprovalTaskEvent) GetVersion() int {
	return e.Version
}
