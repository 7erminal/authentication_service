package api

type RequestMethod int

const (
	GET RequestMethod = iota
	POST
	PUT
	DELETE
	HEAD
)

func (r RequestMethod) String() string {
	switch r {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	default:
		panic("Invalid request method")
	}
}

type Request struct {
	BaseURL     string            ``
	Path        string            ``
	Method      RequestMethod     ``
	Params      map[string]string ``
	HeaderField map[string]string ``
	// ForwardTraceHeaders forwards the inbound X-Request-Id, X-Trace-Id and
	// X-Session-Id to the target. Set it (via NewInternalRequest) only for
	// calls to our own services; third parties (e.g. Google) must not
	// receive internal correlation IDs.
	ForwardTraceHeaders bool ``
}

func NewRequest(baseURL string, path string, method RequestMethod) *Request {
	return &Request{
		BaseURL:     baseURL,
		Path:        path,
		Method:      method,
		Params:      map[string]string{},
		HeaderField: map[string]string{},
	}
}

// NewInternalRequest is NewRequest for calls to other internal services: the
// tracing headers of the current request are forwarded unchanged.
func NewInternalRequest(baseURL string, path string, method RequestMethod) *Request {
	r := NewRequest(baseURL, path, method)
	r.ForwardTraceHeaders = true
	return r
}
