package api

import (
	"net/http"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/app/server/httpserver"
)

// Handler ...
type Handler struct {
	httpserver.HandlerVars
}

var _ server.HTTPHandler = (*Handler)(nil)

// NewHandler ...
func NewHandler(method string, pattern string, group string) server.HTTPHandler {
	data := Handler{}
	data.HandlerVars.Method = method
	data.HandlerVars.Pattern = pattern
	data.HandlerVars.Group = group
	return &data
}

// IsGroup ...
func (h *Handler) IsGroup(groupName string) bool {
	return h.HandlerVars.Group == groupName
}

// GetGroup ...
func (h *Handler) GetGroup() string {
	return h.HandlerVars.Group
}

// GetMethod ...
func (h *Handler) GetMethod() string {
	return h.HandlerVars.Method
}

// GetPattern ...
func (h *Handler) GetPattern() string {
	return h.HandlerVars.Pattern
}

// SetPattern ...
func (h *Handler) SetPattern(pattern string) {
	h.HandlerVars.Pattern = pattern
}

// SetAppContext ...
func (h *Handler) SetAppContext(appCtx appContext.IAppContext) {
	h.HandlerVars.App = appCtx
}

// Process ...
func (h *Handler) Process(_ http.ResponseWriter, r *http.Request) (data interface{}, pageToken string, httpStatus int, responseCode string, err error) {
	motor_brand := r.URL.Query().Get("brand")
	motor_type := r.URL.Query().Get("type")
	motor_transmission := r.URL.Query().Get("transmission")
	data, err = h.App.GetServices().Vehicle.GetVehicle(
		r.Context(),
		motor_brand,
		motor_type,
		motor_transmission,
	)

	return
}
