package config

type App struct {
	Database Database `mapstructure:"database"`
	Cache    Cache    `mapstructure:"cache"`
	LogLevel string   `mapstructure:"log_level" default:"info"`
	Version  string   `mapstructure:"version" default:"0.0.1"`
	BotToken string   `mapstructure:"bot_token"`
}

type Database struct {
	Driver   string `mapstructure:"driver" default:"mysql"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"3306"`
	Name     string `mapstructure:"name" default:"tgbot"`
	User     string `mapstructure:"username" default:"root"`
}

type Cache struct {
	Driver   string `mapstructure:"driver" default:"redis"`
	Password string `mapstructure:"password" default:""`
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     string `mapstructure:"port" default:"3306"`
	//these are for redis
	DB     int    `mapstructure:"db" default:"0"`
	Prefix string `mapstructure:"prefix" default:"bot_cache_prefix"`
}
