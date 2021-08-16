package doc

// Generic Error Response
//
// Error responses are sent when an error (e.g. unauthorized, bad request, ...) occurred.
//
// swagger:model jsonError
type jsonError struct {
	// Name is the error name.
	//
	// example: The requested resource could not be found
	Name string `json:"error"`

	// Description contains further information on the nature of the error.
	//
	// example: Object with ID 12345 does not exist
	Description string `json:"error_description"`

	// Code represents the error status code (404, 403, 401, ...).
	//
	// example: 404
	Code int `json:"status_code"`

	// Debug contains debug information. This is usually not available and has to be enabled.
	//
	// example: The database adapter was unable to find the element
	Debug string `json:"error_debug"`
}
