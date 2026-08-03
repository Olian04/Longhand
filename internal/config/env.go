package config

import (
	"os"

	"strings"
)

// ENV keys (applied after YAML, before flags):
//

//	APP_LOGGING_LEVEL
//	APP_LOGGING_FORMAT
//	APP_LOGGING_STREAM

func applyEnv(cfg *Config) error {

	if v, ok := lookupEnv("APP_LOGGING_LEVEL"); ok {
		cfg.Logging.Level = v
	}
	if v, ok := lookupEnv("APP_LOGGING_FORMAT"); ok {
		cfg.Logging.Format = v
	}
	if v, ok := lookupEnv("APP_LOGGING_STREAM"); ok {
		cfg.Logging.Stream = v
	}
	return nil
}

func lookupEnv(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", false
	}
	return v, true
}
