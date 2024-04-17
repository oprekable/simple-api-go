package root

import (
	"context"
	"expvar"
	"fmt"
	"os"
	"runtime"
	"simple-api-go/internal/app/config"
	appErr "simple-api-go/internal/app/error"
	"simple-api-go/internal/app/handler"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/app/server/httpserver"
	"simple-api-go/internal/app/server/mworker"
	"simple-api-go/internal/app/server/watermillserver"
	"simple-api-go/internal/pkg/metrics"
	"simple-api-go/internal/pkg/utils/atexit"
	"simple-api-go/internal/pkg/utils/configloader"
	"simple-api-go/internal/pkg/utils/log"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/creasty/defaults"
	godotenvFS "github.com/driftprogramming/godotenv"
	"github.com/joho/godotenv"
	"github.com/samber/lo/parallel"
)

// InitErrors ...
func InitErrors() {
	App.InitErrors(appErr.Errors...)
}

// InitConfig ...
func InitConfig() {
	var cfg config.Data
	_, err := hunch.Waterfall(
		App.GetCtx(),
		// Set env from embedFS file
		func(c context.Context, _ interface{}) (interface{}, error) {
			fileEnvPath := fmt.Sprintf("embeds/envs/%s/.env", App.GetEnvironment())
			_ = godotenvFS.Load(*App.GetEmbedFS(), fileEnvPath)
			return nil, nil
		},
		// Set env from regular file
		func(c context.Context, _ interface{}) (interface{}, error) {
			_, e := os.Stat(App.GetWorkDirPath() + "/params/.env")

			if e != nil {
				_ = godotenv.Overload("params/.env")
			} else {
				_ = godotenv.Overload(App.GetWorkDirPath() + "/params/.env")
			}

			return nil, nil
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			cPFS := append(ConfigPaths[:], fmt.Sprintf("embeds/params/%s/*.toml", App.GetEnvironment()))
			return nil, configloader.FromFS(App.GetEmbedFS(), cPFS, &cfg)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			cP := append(ConfigPaths[:], fmt.Sprintf("%s/params/*.toml", App.GetWorkDirPath()))
			return nil, configloader.FromFiles(cP, &cfg)
		},
		func(c context.Context, _ interface{}) (interface{}, error) {
			return nil, defaults.Set(&cfg)
		},
	)

	if err != nil {
		panic(fmt.Errorf("failed to init config %v", err))
	}

	App.InitConfig(&cfg)
}

func InitLogger() {
	App.InitLogger()
	log.Msg(App.GetCtx(), "Success call InitLogger")
}

// InitExpVar ...
func InitExpVar() {
	// stats = NewStats("stats")
	LastUpdate := &metrics.TimeVar{}

	LastUpdate.Set(time.Now())
	expvar.Publish("last_update", LastUpdate)
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return fmt.Sprintf("%d", runtime.NumGoroutine())
	}))

	expvar.Publish("cgocall", expvar.Func(func() interface{} {
		return fmt.Sprintf("%d", runtime.NumCgoCall())
	}))

	expvar.Publish("cpu", expvar.Func(func() interface{} {
		return fmt.Sprintf("%d", runtime.NumCPU())
	}))

	log.Msg(App.GetCtx(), "Success call InitExpVar")
}

// InitDependencies init all DB, persistence's, etc connections
func InitDependencies() {
	_, err := hunch.All(
		App.GetCtx(),
		func(c context.Context) (interface{}, error) {
			App.InitRedis()
			return nil, nil
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.InitPostgresWrite()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.InitPostgresRead()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.InitMysqlWrite()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.InitMysqlRead()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.InitWatermill()
		},
		func(c context.Context) (interface{}, error) {
			App.InitMachineryServer(App.GetDefaultQueueName(), 600) // 600 seconds = 10 minutes
			return nil, nil
		},
	)

	if err != nil {
		panic(err)
	}

	App.WiringRepositories()
	App.WiringServices()
	App.WiringMachineryTaskMap()
}

func PreStart() {
	InitErrors()
	InitConfig()
	InitLogger()
	InitExpVar()
	InitDependencies()
}

func Start() {
	atexit.Add(Shutdown)

	App.GetEg().Go(func() error {
		//lint:ignore S1000 please skip
		select {
		case <-App.GetCtx().Done():
			atexit.AtExit()
			log.Msg(App.GetCtx(), "[shutdown] application ended")
		}

		return nil
	})

	servers := []server.IServer{
		func() server.IServer {
			s := httpserver.NewHTTPServer(
				App,
			)

			i := make([]interface{}, len(handler.HTTPHandlers))
			for k := range handler.HTTPHandlers {
				i[k] = handler.HTTPHandlers[k]
			}

			s.AddHandlers(i...)
			return s
		}(),
		func() server.IServer {
			s := watermillserver.NewWatermill(
				App,
			)

			i := make([]interface{}, len(handler.WatermillHandlers))
			for k := range handler.WatermillHandlers {
				i[k] = handler.WatermillHandlers[k]
			}

			s.AddHandlers(i...)
			return s
		}(),
		func() server.IServer {
			mRepo := mworker.NewRepo(App.GetMachineryServer())
			return mworker.NewMachineryWorker(
				App,
				mRepo,
				mworker.NewSvc(
					App.GetEnvironment(),
					App.GetMachineID(),
					App.GetConfig(),
					mRepo,
				),
			)
		}(),
	}

	parallel.ForEach(servers, func(s server.IServer, _ int) {
		atexit.Add(s.StopApp)
		s.StartApp()
	})
}

func Shutdown() {
	ctx := App.GetCtx()
	_, err := hunch.All(
		ctx,
		func(c context.Context) (interface{}, error) {
			return nil, App.StopPostgresWrite()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.StopPostgresRead()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.StopMysqlWrite()
		},
		func(c context.Context) (interface{}, error) {
			return nil, App.StopMysqlRead()
		},
	)

	log.Err(ctx, "failed to shutdown", err)
}
