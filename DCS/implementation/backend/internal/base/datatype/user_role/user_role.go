package user_role

type UserRole string

const (
	// Human User Roles - Template Management
	TemplateCreator  UserRole = "TEMPLATE_CREATOR"
	TemplateReviewer          = "TEMPLATE_REVIEWER"
	TemplateApprover          = "TEMPLATE_APPROVER"
	TemplateManager           = "TEMPLATE_MANAGER"

	// Human User Roles - Contract Management
	ContractCreator  UserRole = "CONTRACT_CREATOR"
	ContractReviewer          = "CONTRACT_REVIEWER"
	ContractApprover          = "CONTRACT_APPROVER"
	ContractManager           = "CONTRACT_MANAGER"
	ContractSigner            = "CONTRACT_SIGNER"
	ContractObserver          = "CONTRACT_OBSERVER"

	// Human User Roles - System Administration
	ArchiveManager     UserRole = "ARCHIVE_MANAGER"
	Auditor                     = "AUDITOR"
	SystemAdmin                 = "SYSTEM_ADMINISTRATOR"
	ComplianceOfficer           = "COMPLIANCE_OFFICER"
	IntegrationManager          = "INTEGRATION_MANAGER"

	// Human User Roles - Process Management
	ProcessOrchestrator UserRole = "PROCESS_ORCHESTRATOR"
	Validator                    = "VALIDATOR"

	// System User Roles - API/Automated
	SystemContractCreator  UserRole = "SYSTEM_CONTRACT_CREATOR"
	SystemContractReviewer          = "SYSTEM_CONTRACT_REVIEWER"
	SystemContractApprover          = "SYSTEM_CONTRACT_APPROVER"
	SystemContractManager           = "SYSTEM_CONTRACT_MANAGER"
	SystemContractSigner            = "SYSTEM_CONTRACT_SIGNER"
	ContractTargetSystem            = "CONTRACT_TARGET_SYSTEM"
)

// IsValid checks if the UserRole is a valid role
func (r UserRole) IsValid() bool {
	switch r {
	case TemplateCreator, TemplateReviewer, TemplateApprover, TemplateManager,
		ContractCreator, ContractReviewer, ContractApprover, ContractManager,
		ContractSigner, ContractObserver,
		ArchiveManager, Auditor, SystemAdmin, ComplianceOfficer, IntegrationManager,
		ProcessOrchestrator, Validator,
		SystemContractCreator, SystemContractReviewer, SystemContractApprover,
		SystemContractManager, SystemContractSigner, ContractTargetSystem:
		return true
	}
	return false
}

// String returns the string representation of the UserRole
func (r UserRole) String() string {
	return string(r)
}

// IsSystemRole returns true if the role is a system/automated role
func (r UserRole) IsSystemRole() bool {
	switch r {
	case SystemContractCreator, SystemContractReviewer, SystemContractApprover,
		SystemContractManager, SystemContractSigner, ContractTargetSystem:
		return true
	}
	return false
}

// IsHumanRole returns true if the role is a human user role
func (r UserRole) IsHumanRole() bool {
	return r.IsValid() && !r.IsSystemRole()
}
