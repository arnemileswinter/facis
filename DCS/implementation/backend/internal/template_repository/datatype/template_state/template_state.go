package template_state

// TemplateState represents the lifecycle state of a contract template
type TemplateState string

const (
	Draft      TemplateState = "DRAFT"
	Submitted                = "SUBMITTED"
	Reviewed                 = "REVIEWED"
	Created                  = "CREATED"
	Changed                  = "CHANGED"
	Approved                 = "APPROVED"
	Archived                 = "ARCHIVED"
	Deprecated               = "DEPRECATED"
	Deleted                  = "DELETED"
)

var validStates = map[TemplateState]bool{
	Draft:      true,
	Submitted:  true,
	Reviewed:   true,
	Created:    true,
	Changed:    true,
	Approved:   true,
	Archived:   true,
	Deprecated: true,
	Deleted:    true,
}

// IsValid checks if the TemplateState is a valid role
func (s TemplateState) IsValid() bool {
	return validStates[s]
}

// String returns the string representation of the TemplateState
func (s TemplateState) String() string {
	return string(s)
}
