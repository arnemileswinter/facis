package template_state

import (
	"fmt"
	"strings"
)

// TemplateState represents the lifecycle state of a contract template
type TemplateState string

const (
	Draft     TemplateState = "DRAFT"
	Submitted               = "SUBMITTED"
	Reviewed                = "REVIEWED"
	Approved                = "APPROVED"

	Changed    = "CHANGED"
	Archived   = "ARCHIVED"
	Deprecated = "DEPRECATED"
	Deleted    = "DELETED"
)

var validStates = map[TemplateState]bool{
	Draft:     true,
	Submitted: true,
	Reviewed:  true,
	Approved:  true,

	Changed:    true,
	Archived:   true,
	Deprecated: true,
	Deleted:    true,
}

func NewTemplateState(s string) (TemplateState, error) {
	ts := TemplateState(strings.ToUpper(s))
	if !ts.IsValid() {
		return "", fmt.Errorf(fmt.Sprintf("invalid template state: %s", s))
	}
	return ts, nil
}

// IsValid checks if the TemplateState is a valid role
func (s TemplateState) IsValid() bool {
	upper := TemplateState(strings.ToUpper(string(s)))
	return validStates[upper]
}

// String returns the string representation of the TemplateState
func (s TemplateState) String() string {
	return string(s)
}
