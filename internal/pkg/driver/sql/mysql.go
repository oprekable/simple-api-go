package sql

import (
	"fmt"
	"time"

	"github.com/go-gorp/gorp/v3"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

// DBMysqlOption options for mysql connection
type DBMysqlOption struct {
	LogPrefix            string        `deepcopier:"skip"`
	Host                 string        `deepcopier:"field:Host"`
	DB                   string        `deepcopier:"field:DB"`
	Username             string        `deepcopier:"field:Username"`
	Password             string        `deepcopier:"field:Password"`
	AdditionalParameters string        `deepcopier:"field:AdditionalParameters"`
	Port                 int           `deepcopier:"field:Port"`
	ConnOpenMax          int           `deepcopier:"field:ConnOpenMax"`
	ConnLifetimeMax      time.Duration `deepcopier:"field:ConnLifetimeMax"`
	ConnIdleMax          int           `deepcopier:"field:ConnIdleMax"`
}

// NewMysqlDatabase return gorp dbmap object with postgres options param
func NewMysqlDatabase(option DBMysqlOption, logger zerolog.Logger, isDoLogging bool) (g *gorp.DbMap, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		option.Username,
		option.Password,
		option.Host,
		option.Port,
		option.DB,
		option.AdditionalParameters,
	)

	loggerAdapter := NewNoopLog()

	if isDoLogging {
		loggerAdapter = zerologadapter.New(logger)
	}

	db := sqldblogger.OpenDriver(dsn, &mysql.MySQLDriver{}, loggerAdapter)
	db.SetMaxOpenConns(option.ConnOpenMax)
	db.SetConnMaxLifetime(option.ConnLifetimeMax)
	db.SetMaxIdleConns(option.ConnIdleMax)

	g = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{},
	}

	return
}
