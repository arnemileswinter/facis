package templatetype

import (
	"fmt"
	"strings"
)

type TemplateType string

const (
	FrameContract TemplateType = "FRAME_CONTRACT"
	SubContract   TemplateType = "SUB_CONTRACT"
)

var validFlag = map[TemplateType]bool{
	FrameContract: true,
	SubContract:   true,
}

func NewTemplateType(s string) (TemplateType, error) {
	flag := TemplateType(strings.ToUpper(s))
	if !flag.IsValid() {
		return "", fmt.Errorf("invalid template type: %s", s)
	}
	return flag, nil
}

// IsValid checks if the ActionFlag is a valid role
func (f TemplateType) IsValid() bool {
	upper := TemplateType(strings.ToUpper(string(f)))
	return validFlag[upper]
}

// String returns the string representation of the ActionFlag
func (f TemplateType) String() string {
	return string(f)
}
