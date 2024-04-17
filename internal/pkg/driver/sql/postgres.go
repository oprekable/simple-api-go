package sql

import (
	"fmt"
	"time"

	"github.com/go-gorp/gorp/v3"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

// DBPostgresOption options for postgres connection
type DBPostgresOption struct {
	LogPrefix         string        `deepcopier:"skip"`
	Host              string        `deepcopier:"field:Host"`
	DB                string        `deepcopier:"field:DB"`
	Username          string        `deepcopier:"field:Username"`
	Password          string        `deepcopier:"field:Password"`
	SSLMode           string        `deepcopier:"field:sslmode"`
	Port              int           `deepcopier:"field:Port"`
	MaxPoolSize       int           `deepcopier:"field:MaxPoolSize"`
	ConnMaxLifetime   time.Duration `deepcopier:"field:ConnMaxLifetime"`
	MaxIdleConnection int           `deepcopier:"field:MaxIdleConnection"`
}

// NewPostgresDatabase return gorp dbmap object with postgres options param
func NewPostgresDatabase(option DBPostgresOption, logger zerolog.Logger, isDoLogging bool) (g *gorp.DbMap, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		option.Host,
		option.Port,
		option.Username,
		option.DB,
		option.Password,
		option.SSLMode,
	)

	loggerAdapter := NewNoopLog()

	if isDoLogging {
		loggerAdapter = zerologadapter.New(logger)
	}

	db := sqldblogger.OpenDriver(dsn, &pq.Driver{}, loggerAdapter)
	db.SetMaxOpenConns(option.MaxPoolSize)
	db.SetConnMaxLifetime(option.ConnMaxLifetime)
	db.SetMaxIdleConns(option.MaxIdleConnection)

	g = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	return
}
