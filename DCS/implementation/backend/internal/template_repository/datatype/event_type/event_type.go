package event_type

import (
	"fmt"
	"strings"
)

type EventType string

const (
	CreateContractTemplate EventType = "CreateContractTemplate"
	SubmitContractTemplate EventType = "SubmitContractTemplate"
)

var validStates = map[EventType]bool{
	CreateContractTemplate: true,
	SubmitContractTemplate: true,
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
