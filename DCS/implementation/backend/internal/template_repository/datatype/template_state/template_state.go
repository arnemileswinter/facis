package template_state

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// TemplateState represents the lifecycle state of a contract template
type TemplateState string

const (
	Draft     TemplateState = "DRAFT"
	Submitted TemplateState = "SUBMITTED"
	Rejected  TemplateState = "REJECTED"
	Reviewed  TemplateState = "REVIEWED"
	Approved  TemplateState = "APPROVED"
)

var validStates = map[TemplateState]bool{
	Draft:     true,
	Submitted: true,
	Rejected:  true,
	Reviewed:  true,
	Approved:  true,
}

func NewTemplateState(s string) (TemplateState, error) {
	ts := TemplateState(strings.ToUpper(s))
	if !ts.IsValid() {
		return "", fmt.Errorf("invalid template state: %s", s)
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

// Scan implements the sql.Scanner interface
func (s *TemplateState) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("template state cannot be null")
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("unsupported type for TemplateState: %T", value)
	}

	state, err := NewTemplateState(str)
	if err != nil {
		return err
	}

	*s = state
	return nil
}

// Value implements the driver.Valuer interface
func (s TemplateState) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid template state: %s", s)
	}
	return string(s), nil
}
