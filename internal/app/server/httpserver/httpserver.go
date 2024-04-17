package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/app/entity"
	logHttpMiddleware "simple-api-go/internal/app/middleware/http_middleware/log"
	"simple-api-go/internal/app/router"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/pkg/phttp"
	"simple-api-go/internal/pkg/phttp/variable"
	"simple-api-go/internal/pkg/utils/log"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	App                appContext.IAppContext
	ChiMux             *chi.Mux
	HTTPHandlerContext phttp.HTTPHandlerContext
	Handlers           []server.HTTPHandler
}

var _ server.IServer = (*HTTPServer)(nil)

// NewHTTPServer create object server
func NewHTTPServer(appCtx appContext.IAppContext) server.IServer {
	return &HTTPServer{
		App: appCtx,
	}
}

// AddHandlers ...
func (h *HTTPServer) AddHandlers(handlers ...interface{}) {
	hs := make([]server.HTTPHandler, 0)
	for k := range handlers {
		if v, ok := handlers[k].(server.HTTPHandler); ok {
			hs = append(hs, v)
		}
	}

	h.Handlers = hs
}

func (h *HTTPServer) initHTTPHandlerContext(
	ctx context.Context,
	errors []*error,
	responseCodes []variable.ResponseCode,
	allowedHeaders []string,
	allowedOrigins []string,
	logger zerolog.Logger,
	pprofPath string,
) {
	h.HTTPHandlerContext = phttp.NewContextHandler(variable.Meta{
		Version:   entity.Version,
		APIEnv:    entity.Environment,
		GitCommit: entity.GitCommit,
		BuildDate: entity.BuildDate,
	})

	h.HTTPHandlerContext.AddAndMapErrors(errors, responseCodes)

	corsMid := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowCredentials: true,
		AllowedHeaders:   allowedHeaders,
		AllowedOrigins:   allowedOrigins,
	})

	h.ChiMux = chi.NewRouter()
	h.ChiMux.Use(middleware.RequestID)
	h.ChiMux.Use(middleware.RealIP)
	h.ChiMux.Use(middleware.Recoverer)
	h.ChiMux.Use(logHttpMiddleware.NewHTTPRequestLogger(ctx, logger))
	h.ChiMux.Use(corsMid.Handler)
	h.ChiMux.Use(middleware.Timeout(60 * time.Second))

	router.Router(h.App, h.ChiMux, h.HTTPHandlerContext, h.Handlers...)

	var registerRoutes []string
	_ = chi.Walk(h.ChiMux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if !strings.HasPrefix(route, fmt.Sprintf("%s/", pprofPath)) {
			registerRoutes = append(registerRoutes, fmt.Sprintf("[%s] %s", method, route))
		}
		return nil
	})

	h.App.SetRegisteredRoutes(registerRoutes)
}

// StartApp ...
func (h *HTTPServer) StartApp() {
	var srv http.Server

	if h.App.GetCtx() == nil ||
		h.App.GetErrors() == nil ||
		h.App.GetConfig() == nil {
		return
	}

	ctx := h.App.GetCtx()
	config := h.App.GetConfig()

	h.initHTTPHandlerContext(
		ctx,
		h.App.GetErrors(),
		config.ResponseCodes,
		config.Cors.AllowedHeaders,
		config.Cors.AllowedOrigins,
		h.App.GetLogger(),
		config.App.PprofPath,
	)

	srv.Addr = fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)
	srv.Handler = h.ChiMux
	srv.ReadTimeout = 1 * time.Minute
	srv.WriteTimeout = 2 * time.Minute
	srv.IdleTimeout = 5 * time.Minute
	srv.ReadHeaderTimeout = 20 * time.Second

	log.Msg(ctx, fmt.Sprintf("[start] HTTP server at %v", srv.Addr))

	h.App.GetEg().Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic(ctx, "[start] HTTP server failed", err)
		}
		return err
	})

	h.App.GetEg().Go(func() (err error) {
		//lint:ignore S1000 please skip
		select {
		case <-ctx.Done():
			if err = srv.Shutdown(context.Background()); err != nil {
				log.Err(context.Background(), "[shutdown] HTTP server failed to shutting down", err)
			}

			log.Msg(ctx, "[shutdown] HTTP server")
			return
		}
	})
}

// StopApp ...
func (h *HTTPServer) StopApp() {
}
