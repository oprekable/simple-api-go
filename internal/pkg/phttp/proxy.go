package phttp

import "net/http"

// HTTPProxyHandler ...
type HTTPProxyHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)
	CustomWriter
}

// NewHTTPProxyHandler ...
func NewHTTPProxyHandler(c HTTPHandlerContext) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)) HTTPProxyHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)) HTTPProxyHandler {
		return HTTPProxyHandler{H: handler, CustomWriter: CustomWriter{C: c}}
	}
}

// ServeHTTP ...
func (h HTTPProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _, _, _, err := h.H(w, r)
	if err != nil {
		h.WriteError(w, err)
		return
	}
}
