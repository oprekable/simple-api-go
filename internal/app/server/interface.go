package server

import (
	"net/http"
	appContext "simple-api-go/internal/app/context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// IServer interface for server
type IServer interface {
	AddHandlers(handlers ...interface{})
	StartApp()
	StopApp()
}

// HTTPHandler ...
type HTTPHandler interface {
	Process(w http.ResponseWriter, r *http.Request) (data interface{}, pageToken string, httpStatus int, responseCode string, err error)
	IsGroup(groupName string) bool
	GetGroup() string
	GetMethod() string
	GetPattern() string
	SetPattern(pattern string)
	SetAppContext(appCtx appContext.IAppContext)
}

// WatermillHandler ...
type WatermillHandler interface {
	SetApp(app appContext.IAppContext)
	AddHandler(r *message.Router, p *gochannel.GoChannel)
}
