package router

import (
	"fmt"
	"net/http"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/pkg/phttp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	HeaderKeyXClientName = "X-Client-Name"
	HeaderKeyXClientKey  = "X-Client-Key"
)

func EmbedStaticPath(h http.Handler, staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add header(s)
		w.Header().Set("Cache-Control", "no-cache")

		if r.URL.Path == "/" {
			r.URL.Path = fmt.Sprintf("/%s/", staticDir)
		} else {
			b := strings.Split(r.URL.Path, "/")[0]
			if b != staticDir {
				r.URL.Path = fmt.Sprintf("/%s%s", staticDir, r.URL.Path)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func Router(
	appCtx appContext.IAppContext,
	chiMux *chi.Mux,
	httpHandlerContext phttp.HTTPHandlerContext,
	httpHandlers ...server.HTTPHandler,
) {
	ph := phttp.NewHTTPHandler(httpHandlerContext)
	config := appCtx.GetConfig()

	// Do pprof check via http://localhost:3000/dbg/pprof/ where "/dbg" from config
	chiMux.Mount(config.App.PprofPath, middleware.Profiler())

	// Fulfill Option
	for i := range httpHandlers {
		httpHandlers[i].SetAppContext(appCtx)
	}

	proxyHandler := phttp.NewHTTPProxyHandler(httpHandlerContext)
	// For proxy type route
	chiMux.Group(func(r chi.Router) {
		for i := range httpHandlers {
			if httpHandlers[i].IsGroup("proxy") {
				r.Handle(
					httpHandlers[i].GetPattern(),
					proxyHandler(httpHandlers[i].Process),
				)
			}
		}
	})

	// Public Access Group
	chiMux.Group(func(r chi.Router) {
		for i := range httpHandlers {
			if httpHandlers[i].IsGroup("") {
				r.Method(httpHandlers[i].GetMethod(), httpHandlers[i].GetPattern(), ph(httpHandlers[i].Process))
			}
		}
	})

}
