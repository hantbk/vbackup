package config

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Data       Data             `yaml:"data"`
	Logger     LoggerConfig     `yaml:"logger"`
	Jwt        JwtConfig        `yaml:"jwt"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

type PrometheusConfig struct {
	Enabled bool `yaml:"enabled"`
}

type ServerConfig struct {
	Name  string     `yaml:"name"`
	Bind  BindConfig `yaml:"bind"`
	Debug bool       `yaml:"debug"` // Enable or disable debug mode
}

type JwtConfig struct {
	Key    string `yaml:"key"`    // Token key
	MaxAge int    `yaml:"maxAge"` // Token max age
}

type BindConfig struct {
	Host string `yaml:"host"` // ip address
	Port int    `yaml:"port"` // port
}

type Data struct {
	NoCache  bool   `yaml:"noCache"`  // Disable cache
	CacheDir string `yaml:"cacheDir"` // Cache directory
	DbDir    string `yaml:"dbDir"`    // Database directory
}

type LoggerConfig struct {
	Level   string `yaml:"level"`   // log level
	LogPath string `yaml:"logPath"` // log path
}
