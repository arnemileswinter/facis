package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/eventtype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"time"
)

// CreateContractTemplate is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type CreateContractTemplate struct {
	DID            string         `json:"did"`
	DocumentNumber int            `json:"document_number"`
	Version        int            `json:"version"`
	CreatedBy      string         `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Name           *string        `json:"name"`
	Description    *string        `json:"description"`
	TemplateData   *datatype.JSON `json:"templatedata"`
	OccurredAt     time.Time      `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e CreateContractTemplate) EventType() string {
	return eventtype.CreateContractTemplate.String()
}

// GetDID implements the Event interface.
func (e CreateContractTemplate) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e CreateContractTemplate) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e CreateContractTemplate) GetVersion() int {
	return e.Version
}

// SubmitContractTemplateEvent is emitted when a template is submitted for review.
// This event signals state transition and includes reviewer comments.
type SubmitContractTemplateEvent struct {
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
func (e SubmitContractTemplateEvent) EventType() string {
	return eventtype.SubmitContractTemplate.String()
}

// GetDID implements the Event interface.
func (e SubmitContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e SubmitContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e SubmitContractTemplateEvent) GetVersion() int {
	return e.Version
}

// ApproveContractTemplateEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type ApproveContractTemplateEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int       `json:"version"`
	ApprovedBy     string    `json:"approved_by"`
	DecisionNotes  []string  `json:"decision_notes,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ApproveContractTemplateEvent) EventType() string {
	return eventtype.ApproveContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ApproveContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ApproveContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ApproveContractTemplateEvent) GetVersion() int {
	return e.Version
}

// RejectContractTemplateEvent is emitted when a template is rejected.
// This event includes rejection reason and rejector information.
type RejectContractTemplateEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
	RejectedBy     string    `json:"rejected_by"`
	Reason         string    `json:"reason"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RejectContractTemplateEvent) EventType() string {
	return eventtype.RejectContractTemplate.String()
}

// GetDID implements the Event interface.
func (e RejectContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RejectContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RejectContractTemplateEvent) GetVersion() int {
	return e.Version
}

// UpdateContractTemplateEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type UpdateContractTemplateEvent struct {
	DID             string         `json:"did"`
	DocumentNumber  int            `json:"document_number"`
	Version         int            `json:"version"`
	UpdatedBy       string         `json:"updated_by"`
	OldName         *string        `json:"old_name,omitempty"`
	NewName         *string        `json:"new_name,omitempty"`
	OldDescription  *string        `json:"old_description,omitempty"`
	NewDescription  *string        `json:"new_description,omitempty"`
	OldTemplateData *datatype.JSON `json:"old_template_data,omitempty"`
	NewTemplateData *datatype.JSON `json:"new_metadata,omitempty"`
	OccurredAt      time.Time      `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e UpdateContractTemplateEvent) EventType() string {
	return eventtype.UpdateContractTemplate.String()
}

// GetDID implements the Event interface.
func (e UpdateContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e UpdateContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e UpdateContractTemplateEvent) GetVersion() int {
	return e.Version
}

// SearchContractTemplateEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type SearchContractTemplateEvent struct {
	RetrievedBy    string    `json:"updated_by"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e SearchContractTemplateEvent) EventType() string {
	return eventtype.SearchContractTemplate.String()
}

// GetDID implements the Event interface.
func (e SearchContractTemplateEvent) GetDID() string {
	return "*"
}

// GetDocumentNumber implements the Event interface.
func (e SearchContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e SearchContractTemplateEvent) GetVersion() int {
	return e.Version
}

// RetrieveAllContractTemplatesEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveAllContractTemplatesEvent struct {
	RetrievedBy string    `json:"updated_by"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveAllContractTemplatesEvent) EventType() string {
	return eventtype.RetrieveAllContractTemplates.String()
}

// GetDID implements the Event interface.
func (e RetrieveAllContractTemplatesEvent) GetDID() string {
	return "*"
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveAllContractTemplatesEvent) GetDocumentNumber() int {
	return 0
}

// GetVersion implements the Event interface.
func (e RetrieveAllContractTemplatesEvent) GetVersion() int {
	return 0
}

// RetrieveContractTemplateByIdEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveContractTemplateByIdEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveContractTemplateByIdEvent) EventType() string {
	return eventtype.RetrieveContractTemplateById.String()
}

// GetDID implements the Event interface.
func (e RetrieveContractTemplateByIdEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveContractTemplateByIdEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RetrieveContractTemplateByIdEvent) GetVersion() int {
	return e.Version
}

// CreateContractTemplateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type CreateContractTemplateReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Reviewer       string    `json:"reviewer"`
	CreatedBy      string    `json:"created_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e CreateContractTemplateReviewTaskEvent) EventType() string {
	return eventtype.CreateContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e CreateContractTemplateReviewTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e CreateContractTemplateReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e CreateContractTemplateReviewTaskEvent) GetVersion() int {
	return e.Version
}

// CreateContractTemplateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveAllContractTemplateReviewTasksEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveAllContractTemplateReviewTasksEvent) EventType() string {
	return eventtype.RetrieveAllContractTemplateReviewTasks.String()
}

// GetDID implements the Event interface.
func (e RetrieveAllContractTemplateReviewTasksEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveAllContractTemplateReviewTasksEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RetrieveAllContractTemplateReviewTasksEvent) GetVersion() int {
	return e.Version
}

// ReopenContractTemplateReviewAndApprovalTasksEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ReopenContractTemplateReviewAndApprovalTasksEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	CreatedBy      string    `json:"created_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ReopenContractTemplateReviewAndApprovalTasksEvent) EventType() string {
	return eventtype.ReopenContractTemplateReviewAndApprovalTasks.String()
}

// GetDID implements the Event interface.
func (e ReopenContractTemplateReviewAndApprovalTasksEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ReopenContractTemplateReviewAndApprovalTasksEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ReopenContractTemplateReviewAndApprovalTasksEvent) GetVersion() int {
	return e.Version
}

// DeleteContractTemplateReviewTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type DeleteContractTemplateReviewTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	DeletedBy      string    `json:"deleted_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e DeleteContractTemplateReviewTaskEvent) EventType() string {
	return eventtype.DeleteContractTemplateReviewTask.String()
}

// GetDID implements the Event interface.
func (e DeleteContractTemplateReviewTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e DeleteContractTemplateReviewTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e DeleteContractTemplateReviewTaskEvent) GetVersion() int {
	return e.Version
}

// CreateContractTemplateApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type CreateContractTemplateApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `db:"document_number"`
	Version        int       `db:"version"`
	Approver       string    `json:"approver"`
	CreatedBy      string    `json:"created_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e CreateContractTemplateApprovalTaskEvent) EventType() string {
	return eventtype.CreateContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e CreateContractTemplateApprovalTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e CreateContractTemplateApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e CreateContractTemplateApprovalTaskEvent) GetVersion() int {
	return e.Version
}

type RetrieveAllContractTemplateApprovalTasksEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveAllContractTemplateApprovalTasksEvent) EventType() string {
	return eventtype.RetrieveAllContractTemplateApprovalTasks.String()
}

// GetDID implements the Event interface.
func (e RetrieveAllContractTemplateApprovalTasksEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveAllContractTemplateApprovalTasksEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RetrieveAllContractTemplateApprovalTasksEvent) GetVersion() int {
	return e.Version
}

// DeleteContractTemplateApprovalTaskEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type DeleteContractTemplateApprovalTaskEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	DeletedBy      string    `json:"deleted_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e DeleteContractTemplateApprovalTaskEvent) EventType() string {
	return eventtype.DeleteContractTemplateApprovalTask.String()
}

// GetDID implements the Event interface.
func (e DeleteContractTemplateApprovalTaskEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e DeleteContractTemplateApprovalTaskEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e DeleteContractTemplateApprovalTaskEvent) GetVersion() int {
	return e.Version
}
