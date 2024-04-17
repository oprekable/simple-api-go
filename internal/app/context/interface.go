package context

import (
	"context"
	"embed"
	"simple-api-go/internal/app/config"
	"simple-api-go/internal/app/repository"
	"simple-api-go/internal/app/service"
	"simple-api-go/internal/pkg/driver/redis"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/RichardKnop/machinery/v2"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-gorp/gorp/v3"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type IAppContext interface {
	InitConfig(cfg *config.Data)
	InitLogger()
	InitErrors(err ...*error)
	InitRedis()
	InitPostgresWrite() (err error)
	StopPostgresWrite() (err error)
	InitPostgresRead() (err error)
	StopPostgresRead() (err error)
	InitMysqlWrite() (err error)
	StopMysqlWrite() (err error)
	InitMysqlRead() (err error)
	StopMysqlRead() (err error)
	InitMachineryServer(defaultQueue string, resultsExpireIn int)
	AddToEg(in func())
	GetEmbedFS() *embed.FS
	GetIsProduction() bool
	GetDefaultQueueName() string
	GetEnvironment() string
	SetEnvironment(env string)
	GetCtx() context.Context
	GetCtxCancel() context.CancelFunc
	GetEg() *errgroup.Group
	GetConfig() *config.Data
	GetLogger() zerolog.Logger
	GetErrors() []*error
	GetRedis() redis.IRedis
	GetDBPostgresWrite() *gorp.DbMap
	GetDBPostgresRead() *gorp.DbMap
	GetDBMysqlWrite() *gorp.DbMap
	GetDBMysqlRead() *gorp.DbMap
	GetMachineryServer() *machinery.Server
	GetWorkDirPath() string
	GetMachineID() string
	WiringRepositories()
	SetRepositories(r *repository.Repositories)
	GetRepositories() *repository.Repositories
	WiringServices()
	SetServices(s *service.Services)
	GetServices() *service.Services
	GetRegisteredRoutes() []string
	SetRegisteredRoutes([]string)
	WiringMachineryTaskMap()
	SetMachineryTaskMap(m map[string]interface{})
	GetMachineryTaskMap() map[string]interface{}
	IsMachineryActive() bool
	InitWatermill() error
	GetWatermillRouter() *message.Router
	GetWatermillPubSub() *gochannel.GoChannel
	GetTimeLocation() *time.Location
	GetTimeZone() string
	GetTimeFormat() string
	GetTimePostgresFriendlyFormat() string
	GetValidator() *validator.Validate
	GetSingleFlight() *singleflight.Group
}
