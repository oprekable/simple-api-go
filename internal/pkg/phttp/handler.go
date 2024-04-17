package phttp

import "net/http"

// HTTPHandler ...
type HTTPHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)
	CustomWriter
}

// NewHTTPHandler ...
func NewHTTPHandler(c HTTPHandlerContext) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)) HTTPHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, string, int, string, error)) HTTPHandler {
		return HTTPHandler{H: handler, CustomWriter: CustomWriter{C: c}}
	}
}

// ServeHTTP ...
func (hh HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, pageToken, httpStatus, responseCode, err := hh.H(w, r)
	if err != nil {
		hh.WriteError(w, err)
		return
	}

	hh.Write(w, data, pageToken, httpStatus, responseCode)
}
