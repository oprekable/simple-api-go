package cronjob

type TaskArg struct {
	Value string `mapstructure:"value"`
}

type Signature struct {
	UUID         string    `mapstructure:"uuid"`
	Name         string    `mapstructure:"name"`
	TaskArgs     []TaskArg `mapstructure:"task_arg"`
	RetryCount   int       `mapstructure:"retry_count"`
	RetryTimeout int       `mapstructure:"retry_timeout"`
}

type CronJob struct {
	Name      string    `mapstructure:"name"`
	Spec      string    `mapstructure:"spec"`
	Signature Signature `mapstructure:"signature"`
	IsEnabled bool      `mapstructure:"is_enabled"`
}
