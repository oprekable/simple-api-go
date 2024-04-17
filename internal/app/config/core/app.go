package core

// App ..
type App struct {
	Host           string `default:"0.0.0.0" mapstructure:"host"`
	Name           string `default:"-"       mapstructure:"name"`
	Secret         string `default:"-"       mapstructure:"secret"`
	PprofPath      string `default:"/dgb"    mapstructure:"pprof_path"`
	IsMockDeActive string `default:"false"   mapstructure:"is_mock_deactive"`
	Port           int    `default:"3000"    mapstructure:"port"`
}
