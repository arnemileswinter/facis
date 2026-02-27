package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/eventtype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"time"
)

// CreateEvent is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type CreateEvent struct {
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
func (e CreateEvent) EventType() string {
	return eventtype.Create.String()
}

// GetDID implements the Event interface.
func (e CreateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e CreateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e CreateEvent) GetVersion() int {
	return e.Version
}

// SubmitEvent is emitted when a template is submitted for review.
// This event signals state transition and includes reviewer comments.
type SubmitEvent struct {
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
func (e SubmitEvent) EventType() string {
	return eventtype.Submit.String()
}

// GetDID implements the Event interface.
func (e SubmitEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e SubmitEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e SubmitEvent) GetVersion() int {
	return e.Version
}

// ApproveEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type ApproveEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int       `json:"version"`
	ApprovedBy     string    `json:"approved_by"`
	DecisionNotes  []string  `json:"decision_notes,omitempty"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ApproveEvent) EventType() string {
	return eventtype.Approve.String()
}

// GetDID implements the Event interface.
func (e ApproveEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ApproveEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ApproveEvent) GetVersion() int {
	return e.Version
}

// RejectEvent is emitted when a template is rejected.
// This event includes rejection reason and rejector information.
type RejectEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
	RejectedBy     string    `json:"rejected_by"`
	Reason         string    `json:"reason"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RejectEvent) EventType() string {
	return eventtype.Reject.String()
}

// GetDID implements the Event interface.
func (e RejectEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RejectEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RejectEvent) GetVersion() int {
	return e.Version
}

// VerifyEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type VerifyEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int       `json:"version"`
	VerifiedBy     string    `json:"verified_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e VerifyEvent) EventType() string {
	return eventtype.VerifyContractTemplate.String()
}

// GetDID implements the Event interface.
func (e VerifyEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e VerifyEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e VerifyEvent) GetVersion() int {
	return e.Version
}

// UpdateEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type UpdateEvent struct {
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
func (e UpdateEvent) EventType() string {
	return eventtype.Update.String()
}

// GetDID implements the Event interface.
func (e UpdateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e UpdateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e UpdateEvent) GetVersion() int {
	return e.Version
}

// UpdateManageEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type UpdateManageEvent struct {
	DID             string                       `json:"did"`
	DocumentNumber  int                          `json:"document_number"`
	Version         int                          `json:"version"`
	UpdatedBy       string                       `json:"updated_by"`
	OldState        *templatestate.TemplateState `json:"old_state"`
	NewState        *templatestate.TemplateState `json:"new_state"`
	OldName         *string                      `json:"old_name,omitempty"`
	NewName         *string                      `json:"new_name,omitempty"`
	OldDescription  *string                      `json:"old_description,omitempty"`
	NewDescription  *string                      `json:"new_description,omitempty"`
	OldTemplateData *datatype.JSON               `json:"old_template_data,omitempty"`
	NewTemplateData *datatype.JSON               `json:"new_metadata,omitempty"`
	OccurredAt      time.Time                    `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e UpdateManageEvent) EventType() string {
	return eventtype.Update.String()
}

// GetDID implements the Event interface.
func (e UpdateManageEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e UpdateManageEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e UpdateManageEvent) GetVersion() int {
	return e.Version
}

// SearchEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type SearchEvent struct {
	RetrievedBy    string    `json:"updated_by"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e SearchEvent) EventType() string {
	return eventtype.SearchContractTemplate.String()
}

// GetDID implements the Event interface.
func (e SearchEvent) GetDID() string {
	return "*"
}

// GetDocumentNumber implements the Event interface.
func (e SearchEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e SearchEvent) GetVersion() int {
	return e.Version
}

// RetrieveAllEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveAllEvent struct {
	RetrievedBy string    `json:"updated_by"`
	OccurredAt  time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveAllEvent) EventType() string {
	return eventtype.RetrieveAll.String()
}

// GetDID implements the Event interface.
func (e RetrieveAllEvent) GetDID() string {
	return "*"
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveAllEvent) GetDocumentNumber() int {
	return 0
}

// GetVersion implements the Event interface.
func (e RetrieveAllEvent) GetVersion() int {
	return 0
}

// RetrieveByIDEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveByIDEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveByIDEvent) EventType() string {
	return eventtype.RetrieveByID.String()
}

// GetDID implements the Event interface.
func (e RetrieveByIDEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveByIDEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RetrieveByIDEvent) GetVersion() int {
	return e.Version
}

// ArchiveEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ArchiveEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	ArchivedBy     string    `json:"archived_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ArchiveEvent) EventType() string {
	return eventtype.Archive.String()
}

// GetDID implements the Event interface.
func (e ArchiveEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ArchiveEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ArchiveEvent) GetVersion() int {
	return e.Version
}

// RegisterEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RegisterEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RegisteredBy   string    `json:"registered_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RegisterEvent) EventType() string {
	return eventtype.Register.String()
}

// GetDID implements the Event interface.
func (e RegisterEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RegisterEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RegisterEvent) GetVersion() int {
	return e.Version
}
