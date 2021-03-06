package openapistyle

// subset of openAPI definition

// Paths : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#paths-object
type Paths map[string]PathItem

// PathItem : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#paths-object
type PathItem struct {
	// Ref string
	// Summary string
	// Description string
	Get     *Operation `json:"get,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Trace   *Operation `json:"trace,omitempty"`
	// Servers []Server
	// Parameters []Parameter
}

// Operation : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#operation-object
type Operation struct {
	// Tags string
	// Summary string
	// Description string
	// ExternalDocs ExternalDocumentation
	// OperationID string
	Parameters  []Parameter  `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"requestBody,omitempty"`
	Responses   []Response   `json:"responses,omitempty"`
	// Callbacks map[string]Callback
	// Deprecated bool
	// Security SecurityRequirement
	// Servers []Server
}

// Parameter : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#parameter-object
type Parameter struct {
	Name string `json:"name"`
	In   string `json:"in"` // query,header,path,cookie
	// Description string
	// Required bool
	// Deprecated bool
	// AllowEmptyValue bool
	// Style string
	// Explode bool
	// AllowReserved bool
	// Schema Schema
	Example  interface{}   `json:"example,omitempty"`
	Examples []interface{} `json:"examples,omitempty"` // cookies, headers
}

// RequestBody : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#request-body-object
type RequestBody struct {
	// Description string               `json:"description"`
	Content map[string]MediaType `json:"content"`
	// Required bool
}

// MediaType : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#request-body-object
type MediaType struct {
	// Schema Schema
	Example  interface{}   `json:"example,omitempty"`
	Examples []interface{} `json:"examples,omitempty"`
	// Encoding map[string]EncodingType
}

// Responses : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#responses-object
type Responses map[string]Response

// Response : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#response-object
type Response struct {
	// Description string
	Headers map[string]Header    `json:"headers"`
	Content map[string]MediaType `json:"content"`
	// Links []Link
}

// Header : https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.1.0.md#header-object
type Header Parameter
