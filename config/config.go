package config

import (
	"flag"
)

type Config struct {
	ParseOnly bool
	Token     string
	Endpoint  string
}

func New() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.ParseOnly, "parse-only", false, "Returns the parsed email without sending")
	flag.StringVar(&cfg.Token, "token", "", "API service authentication token")
	flag.StringVar(&cfg.Endpoint, "endpoint", "", "API endpoint")
	flag.Parse()
	return cfg
}
