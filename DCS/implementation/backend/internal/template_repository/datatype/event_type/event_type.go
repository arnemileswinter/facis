package event_type

import (
	"fmt"
	"strings"
)

type EventType string

const (
	CreateContractTemplate                   EventType = "CREATE_CONTRACT_TEMPLATE"
	SubmitContractTemplate                   EventType = "SUBMIT_CONTRACT_TEMPLATE"
	ApproveContractTemplate                  EventType = "APPROVE_CONTRACT_TEMPLATE"
	RejectContractTemplate                   EventType = "REJECT_CONTRACT_TEMPLATE"
	UpdateContractTemplate                   EventType = "UPDATE_CONTRACT_TEMPLATE"
	RetrieveAllContractTemplates             EventType = "RETRIEVE_ALL_CONTRACT_TEMPLATES"
	RetrieveContractTemplateById             EventType = "RETRIEVE_CONTRACT_TEMPLATE_BY_ID"
	CreateContractTemplateReviewTask         EventType = "CREATE_CONTRACT_TEMPLATE_REVIEW_TASK"
	RetrieveAllContractTemplateReviewTasks   EventType = "RETRIEVE_ALL_CONTRACT_TEMPLATE_REVIEW_TASKS"
	CreateContractTemplateApprovalTask       EventType = "CREATE_CONTRACT_TEMPLATE_APPROVAL_TASK"
	ReopenContractTemplateReviewTask         EventType = "REOPEN_CONTRACT_TEMPLATE_REVIEW_TASK"
	UpdateContractTemplateReviewTask         EventType = "UPDATE_CONTRACT_TEMPLATE_REVIEW_TASK"
	RetrieveAllContractTemplateApprovalTasks EventType = "RETRIEVE_ALL_CONTRACT_TEMPLATE_APPROVAL_TASKS"
	ReopenContractTemplateApprovalTask       EventType = "REOPEN_CONTRACT_TEMPLATE_APPROVAL_TASK"
	UpdateContractTemplateApprovalTask       EventType = "UPDATE_CONTRACT_TEMPLATE_APPROVAL_TASK"
)

var validStates = map[EventType]bool{
	CreateContractTemplate:                   true,
	SubmitContractTemplate:                   true,
	ApproveContractTemplate:                  true,
	RejectContractTemplate:                   true,
	UpdateContractTemplate:                   true,
	RetrieveAllContractTemplates:             true,
	RetrieveContractTemplateById:             true,
	CreateContractTemplateReviewTask:         true,
	CreateContractTemplateApprovalTask:       true,
	ReopenContractTemplateReviewTask:         true,
	UpdateContractTemplateReviewTask:         true,
	RetrieveAllContractTemplateApprovalTasks: true,
	ReopenContractTemplateApprovalTask:       true,
	UpdateContractTemplateApprovalTask:       true,
}

func NewEventType(s string) (EventType, error) {
	ts := EventType(strings.ToUpper(s))
	if !ts.IsValid() {
		return "", fmt.Errorf("invalid event type: %s", s)
	}
	return ts, nil
}

// IsValid checks if the EventType is a valid role
func (s EventType) IsValid() bool {
	upper := EventType(strings.ToUpper(string(s)))
	return validStates[upper]
}

// String returns the string representation of the EventType
func (s EventType) String() string {
	return string(s)
}
