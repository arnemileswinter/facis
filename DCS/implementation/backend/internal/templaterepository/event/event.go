package event

import (
	"digital-contracting-service/internal/base/datatype"
	"digital-contracting-service/internal/templaterepository/datatype/actionflag"
	"digital-contracting-service/internal/templaterepository/datatype/eventtype"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"time"
)

// CreateContractTemplateEvent is emitted when a new contract template is created.
// This event signals initial template creation with metadata.
type CreateContractTemplateEvent struct {
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
func (e CreateContractTemplateEvent) EventType() string {
	return eventtype.CreateContractTemplate.String()
}

// GetDID implements the Event interface.
func (e CreateContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e CreateContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e CreateContractTemplateEvent) GetVersion() int {
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

// VerifyContractTemplateEvent is emitted when a template is approved.
// This event signals successful approval with optional decision notes.
type VerifyContractTemplateEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	UpdatedAt      time.Time `json:"updated_at"`
	Version        int       `json:"version"`
	VerifiedBy     string    `json:"verified_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e VerifyContractTemplateEvent) EventType() string {
	return eventtype.VerifyContractTemplate.String()
}

// GetDID implements the Event interface.
func (e VerifyContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e VerifyContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e VerifyContractTemplateEvent) GetVersion() int {
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

// RetrieveContractTemplateByIDEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RetrieveContractTemplateByIDEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RetrievedBy    string    `json:"retrieved_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RetrieveContractTemplateByIDEvent) EventType() string {
	return eventtype.RetrieveContractTemplateByID.String()
}

// GetDID implements the Event interface.
func (e RetrieveContractTemplateByIDEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RetrieveContractTemplateByIDEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RetrieveContractTemplateByIDEvent) GetVersion() int {
	return e.Version
}

// ArchiveContractTemplateEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type ArchiveContractTemplateEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	ArchivedBy     string    `json:"archived_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e ArchiveContractTemplateEvent) EventType() string {
	return eventtype.ArchiveContractTemplate.String()
}

// GetDID implements the Event interface.
func (e ArchiveContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e ArchiveContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e ArchiveContractTemplateEvent) GetVersion() int {
	return e.Version
}

// RegisterContractTemplateEvent is emitted when template metadata is updated.
// This event is used for audit and synchronization purposes.
type RegisterContractTemplateEvent struct {
	DID            string    `json:"did"`
	DocumentNumber int       `json:"document_number"`
	Version        int       `json:"version"`
	RegisteredBy   string    `json:"registered_by"`
	OccurredAt     time.Time `json:"occurred_at"`
}

// EventType implements the Event interface.
func (e RegisterContractTemplateEvent) EventType() string {
	return eventtype.RegisterContractTemplate.String()
}

// GetDID implements the Event interface.
func (e RegisterContractTemplateEvent) GetDID() string {
	return e.DID
}

// GetDocumentNumber implements the Event interface.
func (e RegisterContractTemplateEvent) GetDocumentNumber() int {
	return e.DocumentNumber
}

// GetVersion implements the Event interface.
func (e RegisterContractTemplateEvent) GetVersion() int {
	return e.Version
}
