package responsesDTOs

// RequestIDCarrier is implemented by the legacy response DTOs. Their error
// bodies (StatusCode != 200) get an additive top-level "request_id" field so
// clients can quote it; the gateways decode these DTOs into typed structs and
// ignore unknown fields, so this is backwards-compatible. Success bodies are
// unchanged.
type RequestIDCarrier interface {
	WithRequestIDOnError(requestID string) interface{}
}

func isLegacyError(statusCode int) bool { return statusCode != 200 }

// WithRequestIDOnError implements RequestIDCarrier.
func (r StringResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r InviteHashResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r InviteDecodeResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r UserResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r CustomerResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r CustomerAccessTokenResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r LoginTokenResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}

// WithRequestIDOnError implements RequestIDCarrier.
func (r LoginResponseDTO) WithRequestIDOnError(id string) interface{} {
	if isLegacyError(r.StatusCode) {
		r.RequestId = id
	}
	return r
}
