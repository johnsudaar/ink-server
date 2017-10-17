package config

import "github.com/johnsudaar/envconfig"

type Config struct {
	PrinterIP string `env:"PRINTER_IP"`
	SinkURL   string `env:"SINK_URL"`
	Sink      string `env:"SINK"`
}

func InitConfig() Config {
	c := Config{
		PrinterIP: "192.168.1.47",
		SinkURL:   "127.0.0.1:1234",
		Sink:      "http",
	}
	envconfig.Build(&c)
	return c
}
