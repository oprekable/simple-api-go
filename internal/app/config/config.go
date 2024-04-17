package config

import (
	"simple-api-go/internal/app/config/core"
	"simple-api-go/internal/app/config/cronjob"
	"simple-api-go/internal/app/config/machinery"
	"simple-api-go/internal/pkg/phttp/variable"
)

type Data struct {
	core.App
	core.Cors
	ResponseCodes []variable.ResponseCode `mapstructure:"response_code"`
	CronJob       []cronjob.CronJob       `mapstructure:"cron_job"`
	core.Log
	core.Mysql
	core.Postgres
	core.RedisRing
	machinery.Machinery
}
