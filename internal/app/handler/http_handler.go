package handler

import (
	"net/http"
	"simple-api-go/internal/app/handler/http_handler/public/api"
	"simple-api-go/internal/app/handler/http_handler/public/home"
	"simple-api-go/internal/app/server"
)

// Register new handlers here!
var (
	// HTTPHandlers ...
	HTTPHandlers = append(publicHTTPHandlers, applicationHTTPHandlers...)

	publicHTTPHandlers = []server.HTTPHandler{
		home.NewHandler(http.MethodGet, "/", ""),
		api.NewHandler(http.MethodGet, "/api", ""),
	}

	applicationHTTPHandlers = []server.HTTPHandler{}
)
