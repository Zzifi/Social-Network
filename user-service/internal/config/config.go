package config

import "time"

type Config struct {
	HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad() Config {
	var config Config
	config.Address = "http://user_service:8081"
	config.Timeout = 4 * time.Second
	config.Timeout = 60 * time.Second

	return config
}
