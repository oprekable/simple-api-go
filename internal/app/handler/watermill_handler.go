package handler

import (
	"simple-api-go/internal/app/server"
)

// Register new watermill handlers here!
var (
	// WatermillHandlers ...
	WatermillHandlers = append(commonWatermillHandlers, applicationWatermillHandlers...)

	commonWatermillHandlers = []server.WatermillHandler{}

	applicationWatermillHandlers = []server.WatermillHandler{}
)
