package config

import "time"

type AppConfig struct {
	LogLevel string
	Rest     Rest
}

type Rest struct {
	ListenAddress string        `envconfig:"PORT"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT"`
	ServerName    string        `envconfig:"SERVER_NAME"`
	Token         string        `envconfig:"TOKEN"`
}
