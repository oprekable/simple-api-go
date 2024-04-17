package core

import (
	"simple-api-go/internal/pkg/driver/sql"
	"time"

	"github.com/ulule/deepcopier"
)

// MysqlParameters ..
type MysqlParameters struct {
	Password             string        `mapstructure:"password"`
	Host                 string        `mapstructure:"host"`
	DB                   string        `mapstructure:"db"`
	Username             string        `mapstructure:"username"`
	AdditionalParameters string        `mapstructure:"additional_parameters"`
	ConnOpenMax          int           `mapstructure:"conn_open_max"`
	Port                 int           `mapstructure:"port"`
	ConnLifetimeMax      time.Duration `mapstructure:"conn_lifetime_max"`
	ConnIdleMax          int           `mapstructure:"conn_idle_max"`
	IsEnabled            bool          `default:"false"                      mapstructure:"is_enabled"`
	IsMigrationEnable    bool          `mapstructure:"is_migration_enable"`
}

func (pp *MysqlParameters) Options(logPrefix string) (returnData sql.DBMysqlOption) {
	_ = deepcopier.Copy(pp).To(&returnData)
	returnData.LogPrefix = logPrefix
	return
}

// Mysql ..
type Mysql struct {
	Write     MysqlParameters `mapstructure:"write"`
	Read      MysqlParameters `mapstructure:"read"`
	IsEnabled bool            `mapstructure:"is_enabled"`
}
