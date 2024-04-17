package machinery

type Machinery struct {
	IsEnabled string `default:"true"       mapstructure:"is_enabled"`
	Redis     Redis  `mapstructure:"redis"`
}

type Redis struct {
	Host                   string `mapstructure:"host"`
	Password               string `mapstructure:"password"`
	DB                     int    `mapstructure:"db"`
	MaxIdle                int    `mapstructure:"max_idle"`
	IdleTimeout            int    `mapstructure:"idle_timeout"`
	ReadTimeout            int    `mapstructure:"read_timeout"`
	WriteTimeout           int    `mapstructure:"write_timeout"`
	ConnectTimeout         int    `mapstructure:"connect_timeout"`
	NormalTasksPollPeriod  int    `mapstructure:"normal_tasks_poll_period"`
	DelayedTasksPollPeriod int    `mapstructure:"delayed_tasks_poll_period"`
}
