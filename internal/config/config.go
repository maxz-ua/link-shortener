package config

import "time"

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer 	string `yaml:"http_server"`
}

type HTTPServer struct {
	Address 	string `yaml:"address" end-default:"localhost:8080"`
	Timeout 	time.Duration    `yaml:"timeout" end-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" end-default:"60s"`
}
