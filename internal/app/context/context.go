package context

import (
	ctx "context"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"simple-api-go/internal/app/config"
	"simple-api-go/internal/app/repository"
	"simple-api-go/internal/app/service"
	"simple-api-go/internal/pkg/driver/redis"
	"simple-api-go/internal/pkg/driver/sql"
	"simple-api-go/internal/pkg/utils/log"
	"simple-api-go/internal/pkg/zerowater"
	"simple-api-go/variable"
	"strconv"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/RichardKnop/machinery/v2"
	redisBackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisBroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	mConfig "github.com/RichardKnop/machinery/v2/config"
	eagerLock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/denisbrodbeck/machineid"
	"github.com/go-gorp/gorp/v3"
	"github.com/go-playground/validator/v10"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Prod  = "production"
	Local = "local"
)

// appContext the app context struct
type appContext struct {
	logger                     zerolog.Logger
	ctx                        ctx.Context
	appRedis                   redis.IRedis
	repositories               *repository.Repositories
	dbMysqlRead                *gorp.DbMap
	eg                         *errgroup.Group
	watermillRouter            *message.Router
	embedFS                    *embed.FS
	validate                   *validator.Validate
	singleflight               *singleflight.Group
	dbMysqlWrite               *gorp.DbMap
	timeLocation               *time.Location
	dbPostgresWrite            *gorp.DbMap
	ctxCancel                  ctx.CancelFunc
	machineryServer            *machinery.Server
	machineryTaskMap           map[string]interface{}
	services                   *service.Services
	appConfig                  *config.Data
	watermillPubSub            *gochannel.GoChannel
	dbPostgresRead             *gorp.DbMap
	machineID                  string
	timePostgresFriendlyFormat string
	defaultQueueName           string
	timeZone                   string
	workDirPath                string
	timeFormat                 string
	environment                string
	appErrors                  []*error
	registeredRoutes           []string
	dBPostgresWriteConnOnce    sync.Once
	dBPostgresReadConnOnce     sync.Once
	dBMysqlReadConnOnce        sync.Once
	dBMysqlWriteConnOnce       sync.Once
	isProduction               bool
	isLocal                    bool
}

var _ IAppContext = (*appContext)(nil)

// NewAppContext initiate AppContext object
func NewAppContext(
	c ctx.Context,
	cCancel ctx.CancelFunc,
	e *errgroup.Group,
	emFS *embed.FS,
	tZone string,
	tFormat string,
	tPostgresFriendlyFormat string,
) IAppContext {
	wDirPath, _ := os.UserHomeDir()

	if ex, er := os.Executable(); er == nil {
		wDirPath = filepath.Dir(ex)
	}

	mID, er := os.Hostname()
	if mID == "" || er != nil {
		mID = "localhost"
	}

	if mid, er := machineid.ID(); er == nil {
		mID = mid
	}

	tLocation, _ := time.LoadLocation(tZone)

	return &appContext{
		ctx:                        c,
		ctxCancel:                  cCancel,
		eg:                         e,
		embedFS:                    emFS,
		workDirPath:                wDirPath,
		machineID:                  mID,
		timeLocation:               tLocation,
		timeZone:                   tZone,
		timeFormat:                 tFormat,
		timePostgresFriendlyFormat: tPostgresFriendlyFormat,
		validate:                   validator.New(),
		singleflight:               &singleflight.Group{},
	}
}

// InitConfig ...
func (a *appContext) InitConfig(cfg *config.Data) {
	a.appConfig = cfg
	a.isProduction = a.environment == Prod
	//a.isLocal = a.environment == Local
	a.isLocal = false
	a.defaultQueueName = fmt.Sprintf("%s_%s_%s", variable.AppName, "queue", a.environment)
}

// InitLogger ...
func (a *appContext) InitLogger() {
	var writers []io.Writer
	isConsoleLoggingEnabled := a.appConfig != nil && a.appConfig.Log.ConsoleLoggingEnabled
	isFileLoggingEnabled := a.appConfig != nil && a.appConfig.Log.FileLoggingEnabled

	if isConsoleLoggingEnabled {
		writers = append(writers, os.Stdout)
	}

	if isFileLoggingEnabled {
		if a.appConfig.Log.Filename == "" {
			a.appConfig.Log.Filename = variable.AppName + "-" + a.environment + ".log"
		}

		logDir := a.appConfig.Log.Directory
		if logDir == "" {
			logDir = path.Join(a.workDirPath, "/logs/")
		}

		l := &lumberjack.Logger{
			Filename:   path.Join(logDir, a.appConfig.Filename),
			MaxBackups: a.appConfig.Log.MaxBackups, // files
			MaxSize:    a.appConfig.Log.MaxSize,    // megabytes
			MaxAge:     a.appConfig.Log.MaxAge,     // days
			Compress:   a.appConfig.Log.EnableLoggingCompress,
		}

		writers = append(
			writers,
			l,
		)

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)

		go func() {
			for {
				<-c
				_ = l.Rotate()
			}
		}()
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	mw := io.MultiWriter(writers...)
	a.logger = zerolog.New(mw).
		With().
		Timestamp().
		Stack().
		Str("env", a.environment).
		Str("service", variable.AppName).
		Logger()

	a.logger = a.logger.Hook(log.UptimeHook{})
	a.ctx = a.logger.WithContext(a.ctx)
}

// InitErrors ...
func (a *appContext) InitErrors(err ...*error) {
	a.appErrors = append(a.appErrors, err...)
}

// InitRedis ...
func (a *appContext) InitRedis() {
	if a.appConfig == nil || !a.appConfig.RedisRing.IsEnabled {
		return
	}

	logs := log.Zero{
		LogZ:   a.logger,
		LogCtx: a.ctx,
	}

	goRedis.SetLogger(&logs)
	a.appRedis = redis.NewRedis(a.appConfig.RedisRing.Options())
	a.appRedis.SetRedSync()
	a.appRedis.SetCache()
	a.appRedis.SetLimiter()
}

// InitPostgresWrite ...
func (a *appContext) InitPostgresWrite() (err error) {
	if a.appConfig == nil || !(a.appConfig.Postgres.IsEnabled && a.appConfig.Postgres.Write.IsEnabled) {
		return
	}

	a.dBPostgresWriteConnOnce.Do(func() {
		if a.dbPostgresWrite, err = sql.NewPostgresDatabase(
			a.appConfig.Postgres.Write.Options("postgres_write"),
			a.logger,
			a.isLocal,
		); err != nil {
			log.AddStrOrAddErr(a.ctx, err, "db postgres failed to initialize (write)", "db postgres initialized (write)")
		}
	})

	return
}

// StopPostgresWrite ...
func (a *appContext) StopPostgresWrite() (err error) {
	if !(a.dbPostgresWrite != nil && a.dbPostgresWrite.Db != nil) {
		return
	}

	err = a.dbPostgresWrite.Db.Close()
	log.AddStrOrAddErr(a.ctx, err, "[shutdown] failed to close postgres (write)", "[shutdown] postgres closed (write)")

	return
}

// InitPostgresRead ...
func (a *appContext) InitPostgresRead() (err error) {
	if a.appConfig == nil || !(a.appConfig.Postgres.IsEnabled && a.appConfig.Postgres.Read.IsEnabled) {
		return
	}

	a.dBPostgresReadConnOnce.Do(func() {
		if a.dbPostgresRead, err = sql.NewPostgresDatabase(
			a.appConfig.Postgres.Read.Options("postgres_read"),
			a.logger,
			a.isLocal,
		); err != nil {
			log.AddStrOrAddErr(a.ctx, err, "db postgres failed to initialize (read)", "db postgres initialized (read)")
		}
	})

	return
}

// StopPostgresRead ...
func (a *appContext) StopPostgresRead() (err error) {
	if !(a.dbPostgresRead != nil && a.dbPostgresRead.Db != nil) {
		return
	}

	err = a.dbPostgresRead.Db.Close()
	log.AddStrOrAddErr(a.ctx, err, "[shutdown] failed to close postgres (read)", "[shutdown] postgres closed (read)")

	return
}

// InitMysqlWrite ...
func (a *appContext) InitMysqlWrite() (err error) {
	if a.appConfig == nil || !(a.appConfig.Mysql.IsEnabled && a.appConfig.Mysql.Write.IsEnabled) {
		return
	}

	a.dBMysqlWriteConnOnce.Do(func() {
		if a.dbMysqlWrite, err = sql.NewMysqlDatabase(
			a.appConfig.Mysql.Write.Options("mysql_write"),
			a.logger,
			a.isLocal,
		); err != nil {
			log.AddStrOrAddErr(a.ctx, err, "db mysql failed to initialize (write)", "db mysql initialized (write)")
		}
	})

	return
}

// StopMysqlWrite ...
func (a *appContext) StopMysqlWrite() (err error) {
	if !(a.dbMysqlWrite != nil && a.dbMysqlWrite.Db != nil) {
		return
	}

	err = a.dbMysqlWrite.Db.Close()
	log.AddStrOrAddErr(a.ctx, err, "[shutdown] failed to close mysql (write)", "[shutdown] mysql closed (write)")

	return
}

// InitMysqlRead ...
func (a *appContext) InitMysqlRead() (err error) {
	if a.appConfig == nil || !(a.appConfig.Mysql.IsEnabled && a.appConfig.Mysql.Read.IsEnabled) {
		return
	}

	a.dBMysqlReadConnOnce.Do(func() {
		if a.dbMysqlRead, err = sql.NewMysqlDatabase(
			a.appConfig.Mysql.Read.Options("mysql_read"),
			a.logger,
			a.isLocal,
		); err != nil {
			log.AddStrOrAddErr(a.ctx, err, "db mysql failed to initialize (read)", "db mysql initialized (read)")
		}
	})

	return
}

// StopMysqlRead ...
func (a *appContext) StopMysqlRead() (err error) {
	if !(a.dbMysqlRead != nil && a.dbMysqlRead.Db != nil) {
		return
	}

	err = a.dbMysqlRead.Db.Close()
	log.AddStrOrAddErr(a.ctx, err, "[shutdown] failed to close mysql (read)", "[shutdown] mysql closed (read)")

	return
}

// IsMachineryActive ...
func (a *appContext) IsMachineryActive() bool {
	if a.appConfig == nil {
		return false
	}

	isEnabled, _ := strconv.ParseBool(a.appConfig.Machinery.IsEnabled)
	return isEnabled
}

// InitMachineryServer ...
func (a *appContext) InitMachineryServer(defaultQueue string, resultsExpireIn int) {
	if !a.IsMachineryActive() {
		return
	}

	if resultsExpireIn == 0 {
		resultsExpireIn = 3600 // 1 Hours
	}

	cnf := &mConfig.Config{
		DefaultQueue:    defaultQueue,
		ResultsExpireIn: resultsExpireIn,
		Redis: &mConfig.RedisConfig{
			MaxIdle:                a.appConfig.Machinery.Redis.MaxIdle,
			IdleTimeout:            a.appConfig.Machinery.Redis.IdleTimeout,
			ReadTimeout:            a.appConfig.Machinery.Redis.ReadTimeout,
			WriteTimeout:           a.appConfig.Machinery.Redis.WriteTimeout,
			ConnectTimeout:         a.appConfig.Machinery.Redis.ConnectTimeout,
			NormalTasksPollPeriod:  a.appConfig.Machinery.Redis.NormalTasksPollPeriod,
			DelayedTasksPollPeriod: a.appConfig.Machinery.Redis.DelayedTasksPollPeriod,
		},
	}

	redisAddress := fmt.Sprintf("%s@%s", a.appConfig.Machinery.Redis.Password, a.appConfig.Machinery.Redis.Host)
	// Create server instance
	broker := redisBroker.NewGR(cnf, []string{redisAddress}, a.appConfig.Machinery.Redis.DB)
	backend := redisBackend.NewGR(cnf, []string{redisAddress}, a.appConfig.Machinery.Redis.DB)
	lock := eagerLock.New()
	a.machineryServer = machinery.NewServer(cnf, broker, backend, lock)
}

// AddToEg ...
func (a *appContext) AddToEg(in func()) {
	a.eg.Go(func() error {
		select {
		case <-a.ctx.Done():
			log.Err(a.ctx, "error group process killed", errors.New("error group process killed"))
		default:
			in()
		}

		return nil
	})
}

// GetEmbedFS ...
func (a *appContext) GetEmbedFS() *embed.FS {
	return a.embedFS
}

// GetIsProduction ...
func (a *appContext) GetIsProduction() bool {
	return a.isProduction
}

// GetDefaultQueueName ...
func (a *appContext) GetDefaultQueueName() string {
	return a.defaultQueueName
}

// GetEnvironment ...
func (a *appContext) GetEnvironment() string {
	return a.environment
}

// SetEnvironment ...
func (a *appContext) SetEnvironment(env string) {
	a.environment = env
}

// GetCtx ...
func (a *appContext) GetCtx() ctx.Context {
	return a.ctx
}

// GetCtxCancel ...
func (a *appContext) GetCtxCancel() ctx.CancelFunc {
	return a.ctxCancel
}

// GetEg ...
func (a *appContext) GetEg() *errgroup.Group {
	return a.eg
}

// GetConfig ...
func (a *appContext) GetConfig() *config.Data {
	return a.appConfig
}

// GetLogger ...
func (a *appContext) GetLogger() zerolog.Logger {
	return a.logger
}

// GetErrors ...
func (a *appContext) GetErrors() []*error {
	return a.appErrors
}

// GetRedis ...
func (a *appContext) GetRedis() redis.IRedis {
	return a.appRedis
}

// GetDBPostgresWrite ...
func (a *appContext) GetDBPostgresWrite() *gorp.DbMap {
	return a.dbPostgresWrite
}

// GetDBPostgresRead ...
func (a *appContext) GetDBPostgresRead() *gorp.DbMap {
	return a.dbPostgresRead
}

// GetDBMysqlWrite ...
func (a *appContext) GetDBMysqlWrite() *gorp.DbMap {
	return a.dbMysqlWrite
}

// GetDBMysqlRead ...
func (a *appContext) GetDBMysqlRead() *gorp.DbMap {
	return a.dbMysqlRead
}

// GetMachineryServer ...
func (a *appContext) GetMachineryServer() *machinery.Server {
	return a.machineryServer
}

// GetWorkDirPath ...
func (a *appContext) GetWorkDirPath() string {
	return a.workDirPath
}

// SetRepositories ...
func (a *appContext) SetRepositories(r *repository.Repositories) {
	a.repositories = r
}

// GetRepositories ...
func (a *appContext) GetRepositories() *repository.Repositories {
	return a.repositories
}

// SetServices ...
func (a *appContext) SetServices(s *service.Services) {
	a.services = s
}

// GetServices ...
func (a *appContext) GetServices() *service.Services {
	return a.services
}

// GetMachineID ...
func (a *appContext) GetMachineID() string {
	return a.machineID
}

// GetRegisteredRoutes ...
func (a *appContext) GetRegisteredRoutes() []string {
	return a.registeredRoutes
}

// SetRegisteredRoutes ...
func (a *appContext) SetRegisteredRoutes(r []string) {
	a.registeredRoutes = r
}

// SetMachineryTaskMap ...
func (a *appContext) SetMachineryTaskMap(m map[string]interface{}) {
	a.machineryTaskMap = m
}

// GetMachineryTaskMap ...
func (a *appContext) GetMachineryTaskMap() map[string]interface{} {
	return a.machineryTaskMap
}

// InitWatermill ...
func (a *appContext) InitWatermill() (err error) {
	logger := watermill.NewStdLogger(true, false)
	if !reflect.DeepEqual(a.logger, zerolog.Logger{}) {
		logger = zerowater.NewZerologLoggerAdapter(
			a.ctx,
			a.logger.With().Str("component", "watermill").Logger(),
		)
	}

	router, e := message.NewRouter(message.RouterConfig{}, logger)
	if e != nil {
		return e
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	a.watermillRouter = router
	a.watermillPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
	return
}

// GetWatermillRouter ...
func (a *appContext) GetWatermillRouter() *message.Router {
	return a.watermillRouter
}

// GetWatermillPubSub ...
func (a *appContext) GetWatermillPubSub() *gochannel.GoChannel {
	return a.watermillPubSub
}

// GetTimeLocation ...
func (a *appContext) GetTimeLocation() *time.Location {
	return a.timeLocation
}

// GetTimeZone ...
func (a *appContext) GetTimeZone() string {
	return a.timeZone
}

// GetTimeFormat ...
func (a *appContext) GetTimeFormat() string {
	return a.timeFormat
}

// GetTimePostgresFriendlyFormat ...
func (a *appContext) GetTimePostgresFriendlyFormat() string {
	return a.timePostgresFriendlyFormat
}

// GetValidator ...
func (a *appContext) GetValidator() *validator.Validate {
	return a.validate
}

func (a *appContext) GetSingleFlight() *singleflight.Group {
	return a.singleflight
}
