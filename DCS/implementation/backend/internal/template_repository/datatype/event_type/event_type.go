package event_type

import (
	"fmt"
	"strings"
)

type EventType string

const (
	CreatedContractTemplate       EventType = "CREATE_CONTRACT_TEMPLATE"
	SubmittedContractTemplate     EventType = "SUBMITTED_CONTRACT_TEMPLATE"
	ApprovedContractTemplate      EventType = "APPROVED_CONTRACT_TEMPLATE"
	RejectedContractTemplate      EventType = "REJECTED_CONTRACT_TEMPLATE"
	UpdatedContractTemplate       EventType = "UPDATE_CONTRACT_TEMPLATE"
	RetrievedAllContractTemplates EventType = "RETRIEVED_ALL_CONTRACT_TEMPLATES"
	RetrievedContractTemplateById EventType = "RETRIEVED_CONTRACT_TEMPLATE_BY_ID"
)

var validStates = map[EventType]bool{
	CreatedContractTemplate:       true,
	SubmittedContractTemplate:     true,
	ApprovedContractTemplate:      true,
	RejectedContractTemplate:      true,
	UpdatedContractTemplate:       true,
	RetrievedAllContractTemplates: true,
	RetrievedContractTemplateById: true,
}

func NewEventType(s string) (EventType, error) {
	ts := EventType(strings.ToUpper(s))
	if !ts.IsValid() {
		return "", fmt.Errorf(fmt.Sprintf("invalid template state: %s", s))
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
