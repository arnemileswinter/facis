package design

import (
	. "goa.design/goa/v3/dsl"
)

var CreateRequest = Type("CreateRequest", func() {
	Description("Contract template create request")

	Token("token", String, "JWT token")

	Attribute("template_type", String, "The type of the template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("template_data", Any, "The template data of the contract template")

	Required("template_type")
})

var CreateResponse = Type("CreateResponse", func() {
	Description("Result for creating a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var SubmitRequest = Type("SubmitRequest", func() {
	Description("Contract template submit request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("forward_to", String, "Action flag: approval | draft")
	Attribute("comments", ArrayOf(String), "Optional comments")

	Required("did", "document_number", "version", "updated_at")
})

var SubmitResponse = Type("SubmitResponse", func() {
	Description("Result for submitting a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var UpdateRequest = Type("UpdateRequest", func() {
	Description("Contract template update request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("template_type", String, "The type of the template")
	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("template_data", Any, "The template data of the contract template")

	Required("did", "document_number", "version", "updated_at")
})

var UpdateResponse = Type("UpdateResponse", func() {
	Description("Result for updating a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var UpdateManageRequest = Type("UpdateManageRequest", func() {
	Description("Contract template update manage request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("state", String, "The state of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("template_type", String, "The type of the template")
	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("template_data", Any, "The template data of the contract template")

	Required("did", "document_number", "version", "updated_at")
})

var UpdateManageResponse = Type("UpdateManageResponse", func() {
	Description("Result for updating a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var SearchRequest = Type("SearchRequest", func() {
	Description("Contract template search request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")
	Attribute("template_type", String, "The type of the template")
	Attribute("state", String, "The state of the contract template")
	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("filter", String, "Search value for full text search in template data")
})

var SearchResponse = Type("SearchResponse", func() {
	Description("Result for searching a contract templates by filter")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The document number of the contract template")
	Attribute("version", Int, "The version number of the contract template")

	Attribute("state", String, "The state of the contract template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")

	Attribute("created_at", String, "The timestamp when the contract template was created")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Required("did", "document_number", "version", "state", "created_at", "updated_at")
})

var RetrieveRequest = Type("RetrieveRequest", func() {
	Description("Contract template retrieve request")

	Token("token", String, "JWT token")
})

var ContractTemplateItem = Type("ContractTemplateItem", func() {
	Attribute("did", String, "DID of the contract template")
	Attribute("document_number", Int, "Document number")
	Attribute("version", Int, "Version")
	Attribute("state", String, "State")
	Attribute("template_type", String, "The type of the template")
	Attribute("name", String, "Name")
	Attribute("description", String, "Description")
	Attribute("created_at", String, "Created at")
	Attribute("updated_at", String, "Updated at")

	Required("did", "document_number", "version", "state", "created_at", "updated_at")
})

var ReviewTaskItem = Type("ReviewTaskItem", func() {
	Attribute("did", String, "DID of the contract template")
	Attribute("document_number", Int, "Document number")
	Attribute("version", Int, "Version")
	Attribute("state", String, "State of the review task")
	Attribute("reviewer", String, "The reviewer of the contract template")
	Attribute("created_at", String, "Created at")

	Required("did", "document_number", "version", "state", "reviewer", "created_at")
})

var ApprovalTaskItem = Type("ApprovalTaskItem", func() {
	Attribute("did", String, "DID of the contract template")
	Attribute("document_number", Int, "Document number")
	Attribute("version", Int, "Version")
	Attribute("state", String, "State of the approval task")
	Attribute("approver", String, "The approver for the contract template")
	Attribute("created_at", String, "Created at")

	Required("did", "document_number", "version", "state", "approver", "created_at")
})

var RetrieveResponse = Type("RetrieveResponse", func() {
	Description("Result for retrieving a contract template by id")

	Attribute("contract_templates", ArrayOf(ContractTemplateItem), "A list of contract templates")

	Attribute("review_tasks", ArrayOf(ReviewTaskItem), "A list of review tasks")

	Attribute("approval_tasks", ArrayOf(ApprovalTaskItem), "A list of approval tasks")

	Required("contract_templates", "review_tasks", "approval_tasks")
})

var RetrieveByIDRequest = Type("RetrieveByIDRequest", func() {
	Description("Contract template retrieve by id request")

	Token("token", String, "JWT token")

	Attribute("did", String, "DID of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var RetrieveByIDResponse = Type("RetrieveByIDResponse", func() {
	Description("Result for retrieving a contract template by id")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The document number of the contract template")
	Attribute("version", Int, "The version number of the contract template")

	Attribute("state", String, "The state of the contract template")
	Attribute("template_type", String, "The type of the template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")

	Attribute("created_by", String, "Identifier of who created the contract template")
	Attribute("created_at", String, "The timestamp when the contract template was created")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("template_data", Any, "The template data of the contract template")

	Required("did", "document_number", "version", "state", "created_by", "created_at", "updated_at", "template_data")
})

var ApproveRequest = Type("ApproveRequest", func() {
	Description("Contract template approve request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("decision_notes", ArrayOf(String), "A list of decision notes")

	Required("did", "document_number", "version", "updated_at")
})

var ApproveResponse = Type("ApproveResponse", func() {
	Description("Result for approving a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var RejectRequest = Type("RejectRequest", func() {
	Description("Contract template retrieve by id request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("reason", String, "Reason for rejecting the contract template")

	Required("did", "document_number", "version", "updated_at", "reason")
})

var RejectResponse = Type("RejectResponse", func() {
	Description("Result for rejecting a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var VerifyRequest = Type("VerifyRequest", func() {
	Description("Contract template verify request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Attribute("decision_notes", ArrayOf(String), "A list of decision notes")

	Required("did", "document_number", "version", "updated_at")
})

var VerifyResponse = Type("VerifyResponse", func() {
	Description("Result for verifying a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var ArchiveRequest = Type("ArchiveRequest", func() {
	Description("Contract template archive request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Required("did", "document_number", "version", "updated_at")
})

var ArchiveResponse = Type("ArchiveResponse", func() {
	Description("Result for archiving a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var RegisterRequest = Type("RegisterRequest", func() {
	Description("Contract template register request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Attribute("updated_at", String, "The timestamp when the contract template was updated")

	Required("did", "document_number", "version", "updated_at")
})

var RegisterResponse = Type("RegisterResponse", func() {
	Description("Result for register a contract template")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var AuditRequest = Type("AuditRequest", func() {
	Description("Contract template audit request")

	Token("token", String, "JWT token")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

var AuditResponse = Type("AuditResponse", func() {
	Description("Result for auditing a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("document_number", Int, "The number of the contract template")
	Attribute("version", Int, "The version of the contract template")

	Required("did", "document_number", "version")
})

// Template Repository Service  (/template/...)
var _ = Service("TemplateRepository", func() {
	Description("Template Repository APIs (/template/...)")

	// POST /template/create
	Method("create", func() {
		Description("Create a new template.")
		Meta("dcs:requirements", "DCS-IR-TR-01")
		Meta("dcs:tr:components", "Single- or multi-tiered template generation")
		Meta("dcs:ui", "Template Builder")

		Security(JWTAuth, func() {
			Scope("Template Creator")
		})

		Payload(CreateRequest)
		Result(CreateResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/create")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/submit
	Method("submit", func() {
		Description(`with action flag { forwardTo: "approval" | "draft" } and optional reviewComments. allow resubmission path with approver comments.`)
		Meta("dcs:requirements", "DCS-IR-TR-03", "DCS-IR-TR-04", "DCS-IR-TR-05")
		Meta("dcs:tr:components", "Single- or multi-tiered template generation")
		Meta("dcs:ui", "Template Builder, Template Review, Template Approver")

		Security(JWTAuth, func() {
			Scope("Template Creator")
			Scope("Template Reviewer")
			Scope("Template Approver")
		})

		Payload(SubmitRequest)
		Result(SubmitResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/submit")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// PUT /template/update
	Method("update", func() {
		Description("persist reviewer edits (template data/clauses/semantics).")
		Meta("dcs:requirements", "DCS-IR-TR-03")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Review")

		Security(JWTAuth, func() {
			Scope("Template Creator")
			Scope("Template Reviewer")
			Scope("Template Approver")
		})

		Payload(UpdateRequest)
		Result(UpdateResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			PUT("/template/update")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/update
	Method("update_manage", func() {
		Description("update template data or status.")
		Meta("dcs:requirements", "DCS-IR-TR-07")
		Meta("dcs:roles", "Template Manager")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Manager")
		})

		Payload(UpdateManageRequest)
		Result(UpdateManageResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/update")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// GET /template/search
	Method("search", func() {
		Description("perform filtered searches.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-07")
		Meta("dcs:tr:components", "Search Capabilities")
		Meta("dcs:ui", "Template Builder, Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Creator")
			Scope("Template Manager")
		})

		Payload(SearchRequest)
		Result(ArrayOfRequired(SearchResponse))

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/search")
			Param("filter")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// GET /template/retrieve
	Method("retrieve", func() {
		Description("load submitted template and history/provenance summary. fetch reviewed template with metadata, review history, and validation results. fetch all template entries for dashboard view.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-03", "DCS-IR-TR-05", "DCS-IR-TR-08")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Approver, Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Reviewer")
			Scope("Template Approver")
			Scope("Template Manager")
		})

		Payload(RetrieveRequest)
		Result(RetrieveResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/retrieve")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// GET /template/retrieve/{did}
	Method("retrieve_by_id", func() {
		Description("Retrieve a template by template id.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-03", "DCS-FR-TR-19")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Approver, Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Reviewer")
			Scope("Template Approver")
			Scope("Template Manager")
		})

		Payload(RetrieveByIDRequest)
		Result(RetrieveByIDResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/retrieve/{did}")
			Param("did")
			Param("document_number")
			Param("version")

			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/verify
	Method("verify", func() {
		Description("run policy, schema, and semantic validations; return findings.")
		Meta("dcs:requirements", "DCS-IR-TR-03")
		Meta("dcs:tr:components", "Semantic Hub")
		Meta("dcs:ui", "Template Review")

		Security(JWTAuth, func() {
			Scope("Template Reviewer")
		})

		Payload(VerifyRequest)
		Result(VerifyResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/verify")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/approve
	Method("approve", func() {
		Description("mark template as approved, with optional decision notes.")
		Meta("dcs:requirements", "DCS-IR-TR-05", "DCS-IR-TR-06")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Approver")

		Security(JWTAuth, func() {
			Scope("Template Approver")
		})

		Payload(ApproveRequest)
		Result(ApproveResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/approve")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/reject
	Method("reject", func() {
		Description("mark template as rejected, requiring reason field.")
		Meta("dcs:requirements", "DCS-IR-TR-05")
		Meta("dcs:tr:components", "")
		Meta("dcs:ui", "Template Approver")

		Security(JWTAuth, func() {
			Scope("Template Approver")
		})

		Payload(RejectRequest)
		Result(RejectResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/reject")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/register
	Method("register", func() {
		Description("register new template into the repository.")
		Meta("dcs:requirements", "DCS-IR-TR-07")
		Meta("dcs:tr:components", "Contract Templates Storage & Provenance")
		Meta("dcs:ui", "Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Reviewer")
		})

		Payload(RegisterRequest)
		Result(RegisterResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/register")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/archive
	Method("archive", func() {
		Description("archive obsolete template.")
		Meta("dcs:requirements", "DCS-IR-TR-07")
		Meta("dcs:tr:components", "Contract Templates Storage & Provenance")
		Meta("dcs:ui", "Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Reviewer")
		})

		Payload(ArchiveRequest)
		Result(ArchiveResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			POST("/template/archive")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// GET /template/audit
	Method("audit", func() {
		Description("retrieve audit history of template actions.")
		Meta("dcs:requirements", "DCS-IR-TR-07", "DCS-IR-TR-08")
		Meta("dcs:tr:components", "")
		Meta("dcs:ui", "Template Management Dashboard")

		Security(JWTAuth, func() {
			Scope("Template Manager")
		})

		Payload(AuditRequest)
		Result(AuditResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/audit")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})
})
