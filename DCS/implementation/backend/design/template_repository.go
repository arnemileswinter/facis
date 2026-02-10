package design

import (
	. "goa.design/goa/v3/dsl"
)

var ContractTemplateCreateRequest = Type("ContractTemplateCreateRequest", func() {
	Description("Contract template create request")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("meta_data", Any, "The metadata of the contract template")
})

var ContractTemplateCreateResponse = Type("ContractTemplateCreateResponse", func() {
	Description("Result for creating a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Required("did")
})

var ContractTemplateSubmitRequest = Type("ContractTemplateSubmitRequest", func() {
	Description("Contract template submit request")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("forward_to", String, "Action flag: approval | draft")
	Attribute("review_comments", ArrayOf(String), "Optional review comments")

	Required("did")
})

var ContractTemplateSubmitResponse = Type("ContractTemplateSubmitResponse", func() {
	Description("Result for submitting a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Required("did")
})

var ContractTemplateUpdateRequest = Type("ContractTemplateUpdateRequest", func() {
	Description("Contract template update request")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")
	Attribute("meta_data", Any, "The metadata of the contract template")

	Required("did")
})

var ContractTemplateUpdateResponse = Type("ContractTemplateUpdateResponse", func() {
	Description("Result for updating a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Required("did")
})

var ContractTemplateRetrieveResponse = Type("ContractTemplateRetrieveResponse", func() {
	Description("Result for retrieving a contract template")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Attribute("document_number", Int, "The document number of the contract template")
	Attribute("version", Int, "The version number of the contract template")

	Attribute("state", String, "The state of the contract template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")

	Attribute("created_by", String, "Identifier of who created the contract template")
	Attribute("created_at", String, "The timestamp when the contract template was created")

	Attribute("meta_data", Any, "The metadata of the contract template")

	Required("did", "document_number", "version", "state", "created_by", "created_at", "meta_data")
})

var ContractTemplateRetrieveByIdRequest = Type("ContractTemplateRetrieveByIdRequest", func() {
	Description("Contract template retrieve by id request")

	Attribute("template_id", String, "ID of the contract template")

	Required("template_id")
})

var ContractTemplateRetrieveByIdResponse = Type("ContractTemplateRetrieveByIdResponse", func() {
	Description("Result for retrieving a contract template by id")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Attribute("document_number", Int, "The document number of the contract template")
	Attribute("version", Int, "The version number of the contract template")

	Attribute("state", String, "The state of the contract template")

	Attribute("name", String, "The name of the contract template")
	Attribute("description", String, "A description for that template")

	Attribute("created_by", String, "Identifier of who created the contract template")
	Attribute("created_at", String, "The timestamp when the contract template was created")

	Attribute("meta_data", Any, "The metadata of the contract template")

	Required("did", "document_number", "version", "state", "created_by", "created_at", "meta_data")
})

var ContractTemplateApproveRequest = Type("ContractTemplateApproveRequest", func() {
	Description("Contract template approve request")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("decision_notes", ArrayOf(String), "A list of decision notes")

	Required("did")
})

var ContractTemplateApproveResponse = Type("ContractTemplateApproveResponse", func() {
	Description("Result for retrieving a contract template by id")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Required("did")
})

var ContractTemplateRejectRequest = Type("ContractTemplateRejectRequest", func() {
	Description("Contract template retrieve by id request")

	Attribute("did", String, "Decentralized Identifier of the contract template")
	Attribute("reason", String, "Reason for rejecting the contract template")

	Required("did", "reason")
})

var ContractTemplateRejectResponse = Type("ContractTemplateRejectResponse", func() {
	Description("Result for retrieving a contract template by id")

	Attribute("did", String, "Decentralized Identifier of the contract template")

	Required("did")
})

// Template Repository Service  (/template/...)
var _ = Service("TemplateRepository", func() {
	Description("Template Repository APIs (/template/...)")

	// POST /template/create
	Method("create", func() {
		Description("Create a new template.")
		Meta("dcs:requirements", "DCS-IR-TR-01")
		Meta("dcs:roles", "Template Creator")
		Meta("dcs:tr:components", "Single- or multi-tiered template generation")
		Meta("dcs:ui", "Template Builder")

		Payload(ContractTemplateCreateRequest)
		Result(ContractTemplateCreateResponse)

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
		Meta("dcs:roles", "Template Creator", "Template Reviewer", "Template Approver")
		Meta("dcs:tr:components", "Single- or multi-tiered template generation")
		Meta("dcs:ui", "Template Builder, Template Review, Template Approver")

		Payload(ContractTemplateSubmitRequest)
		Result(ContractTemplateSubmitResponse)

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
		Description("persist reviewer edits (metadata/clauses/semantics).")
		Meta("dcs:requirements", "DCS-IR-TR-03")
		Meta("dcs:roles", "Template Creator", "Template Reviewer")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Review")

		Payload(ContractTemplateUpdateRequest)
		Result(ContractTemplateUpdateResponse)

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
		Description("update metadata or status.")
		Meta("dcs:requirements", "DCS-IR-TR-07")
		Meta("dcs:roles", "Template Manager")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Management Dashboard")

		HTTP(func() {
			POST("/template/update")
			Response(StatusOK)
		})

		Result(Int)
	})

	// GET /template/search
	Method("search", func() {
		Description("perform filtered searches.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-07")
		Meta("dcs:roles", "Template Creator", "Template Manager")
		Meta("dcs:tr:components", "Search Capabilities")
		Meta("dcs:ui", "Template Builder, Template Management Dashboard")

		HTTP(func() {
			GET("/template/search")
			Response(StatusOK)
		})

		Result(ArrayOf(Any))
	})

	// GET /template/retrieve
	Method("retrieve", func() {
		Description("load submitted template and history/provenance summary. fetch reviewed template with metadata, review history, and validation results. fetch all template entries for dashboard view.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-03", "DCS-IR-TR-05", "DCS-IR-TR-08")
		Meta("dcs:roles", "Template Reviewer", "Template Approver", "Template Manager")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Approver, Template Management Dashboard")

		Result(ArrayOf(ContractTemplateRetrieveResponse))

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/retrieve")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// GET /template/retrieve/{template_id}
	Method("retrieve_by_id", func() {
		Description("Retrieve a template by template id.")
		Meta("dcs:requirements", "DCS-IR-TR-02", "DCS-IR-TR-03", "DCS-FR-TR-19")
		Meta("dcs:roles", "Template Reviewer", "Template Approver", "Template Manager")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Builder, Template Approver, Template Management Dashboard")

		Payload(ContractTemplateRetrieveByIdRequest)
		Result(ContractTemplateRetrieveByIdResponse)

		Error("bad_request", ErrorResult, "Bad request")
		Error("internal_error", ErrorResult, "Internal server error")

		HTTP(func() {
			GET("/template/retrieve/{template_id}")
			Param("template_id")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_error", StatusInternalServerError)
		})
	})

	// POST /template/verify
	Method("verify", func() {
		Description("run policy, schema, and semantic validations; return findings.")
		Meta("dcs:requirements", "DCS-IR-TR-03")
		Meta("dcs:roles", "Template Reviewer")
		Meta("dcs:tr:components", "Semantic Hub")
		Meta("dcs:ui", "Template Review")

		HTTP(func() {
			POST("/template/verify")
			Response(StatusOK)
		})

		Result(Any)
	})

	// POST /template/approve
	Method("approve", func() {
		Description("mark template as approved, with optional decision notes.")
		Meta("dcs:requirements", "DCS-IR-TR-05", "DCS-IR-TR-06")
		Meta("dcs:roles", "Template Approver")
		Meta("dcs:tr:components", "Template Versioning")
		Meta("dcs:ui", "Template Approver")

		Payload(ContractTemplateApproveRequest)
		Result(ContractTemplateApproveResponse)

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
		Meta("dcs:roles", "Template Approver")
		Meta("dcs:tr:components", "")
		Meta("dcs:ui", "Template Approver")

		Payload(ContractTemplateRejectRequest)
		Result(ContractTemplateRejectResponse)

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
		Meta("dcs:roles", "Template Manager")
		Meta("dcs:tr:components", "Contract Templates Storage & Provenance")
		Meta("dcs:ui", "Template Management Dashboard")

		HTTP(func() {
			POST("/template/register")
			Response(StatusOK)
		})

		Result(Any)
	})

	// POST /template/archive
	Method("archive", func() {
		Description("archive obsolete template.")
		Meta("dcs:requirements", "DCS-IR-TR-07")
		Meta("dcs:roles", "Template Manager")
		Meta("dcs:tr:components", "Contract Templates Storage & Provenance")
		Meta("dcs:ui", "Template Management Dashboard")

		HTTP(func() {
			POST("/template/archive")
			Response(StatusOK)
		})

		Result(Int)
	})

	// GET /template/audit
	Method("audit", func() {
		Description("retrieve audit history of template actions.")
		Meta("dcs:requirements", "DCS-IR-TR-07", "DCS-IR-TR-08")
		Meta("dcs:roles", "Template Manager")
		Meta("dcs:tr:components", "")
		Meta("dcs:ui", "Template Management Dashboard")

		HTTP(func() {
			GET("/template/audit")
			Response(StatusOK)
		})

		Result(ArrayOf(String))
	})
})
