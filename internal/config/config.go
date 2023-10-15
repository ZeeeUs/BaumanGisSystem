package config

import (
	"io"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const formatJSON = "json"

type Config struct {
	Server struct {
		Host        string `envconfig:"SERVER_HOST" default:":9000"`
		MetricsBind string `envconfig:"BIND_METRICS" default:":9090"`
		HealthHost  string `envconfig:"BIND_HEALTH" default:":9091"`
	}

	Service struct {
		LogLevel  string `envconfig:"LOGGER_LEVEL" default:"debug"`
		LogFormat string `envconfig:"LOGGER_FORMAT" default:"console"`
	}
}

func Parse() (*Config, error) {
	var cfg = &Config{}
	err := envconfig.Process("", cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.Service.LogLevel); err == nil {
		level = newLevel
	}

	var out io.Writer = os.Stdout
	if cfg.Service.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Caller().Timestamp().Logger()
}
