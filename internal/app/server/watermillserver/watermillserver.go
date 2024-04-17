package watermillserver

import (
	"errors"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/pkg/utils/log"
)

type Watermill struct {
	App appContext.IAppContext
}

var _ server.IServer = (*Watermill)(nil)

// NewWatermill create object server
func NewWatermill(appCtx appContext.IAppContext) server.IServer {
	return &Watermill{
		App: appCtx,
	}
}

func (s *Watermill) AddHandlers(handlers ...interface{}) {
	for k := range handlers {
		if v, ok := handlers[k].(server.WatermillHandler); ok {
			v.SetApp(s.App)
			v.AddHandler(
				s.App.GetWatermillRouter(),
				s.App.GetWatermillPubSub(),
			)
		}
	}
}

func (s *Watermill) StartApp() {
	s.App.GetEg().Go(func() error {
		ctx := s.App.GetCtx()
		if s.App.GetWatermillRouter() == nil || s.App.GetWatermillPubSub() == nil {
			err := errors.New("watermill was nil")
			log.AddErr(ctx, err)
			return err
		}

		log.Msg(ctx, "[start] watermill")
		if err := s.App.GetWatermillRouter().Run(ctx); err != nil {
			log.Panic(ctx, "[shutdown] watermill", err)
		}

		return nil
	})
}

func (s *Watermill) StopApp() {
}
