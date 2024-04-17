package httpserver

import (
	"net/http"
	appContext "simple-api-go/internal/app/context"
)

// HandlerVars ...
type HandlerVars struct {
	http.Handler
	App     appContext.IAppContext
	Method  string
	Pattern string
	Group   string
}
