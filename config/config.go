package config

import (
	"flag"
)

type Config struct {
	ParseOnly   bool
	Token       string
	Endpoint    string
	ShowVersion bool
}

func New() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.ParseOnly, "parse-only", false, "Returns the parsed email without sending")
	flag.StringVar(&cfg.Token, "token", "", "API service authentication token")
	flag.StringVar(&cfg.Endpoint, "endpoint", "", "API endpoint")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Show version number")
	flag.Parse()
	return cfg
}
