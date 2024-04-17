package core

// Log ..
type Log struct {
	Directory             string `default:"-"    mapstructure:"directory"`
	Filename              string `default:"-"    mapstructure:"filename"`
	MaxAge                int    `default:"28"   mapstructure:"max_age"`
	MaxSize               int    `default:"500"  mapstructure:"max_size"`
	MaxBackups            int    `default:"60"   mapstructure:"max_backups"`
	FileLoggingEnabled    bool   `default:"true" mapstructure:"file_logging_enabled"`
	ConsoleLoggingEnabled bool   `default:"true" mapstructure:"console_logging_enabled"`
	EnableLoggingCompress bool   `default:"true" mapstructure:"enable_compress"`
}
