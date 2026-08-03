package config_test

import (
	"testing"

	"github.com/Olian04/Longhand/internal/config"
)

// Defaults must survive Validate, otherwise a zero-config start fails.
func TestDefaultsValidate(t *testing.T) {
	cfg := config.Config{}.WithDefaults()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("defaults invalid: %v", err)
	}
}

func TestLoggingDefaults(t *testing.T) {
	cfg := config.Config{}.WithDefaults()
	if cfg.Logging.Level != config.DefaultLoggingLevel {
		t.Fatalf("level = %q, want %q", cfg.Logging.Level, config.DefaultLoggingLevel)
	}
	if cfg.Logging.Format != config.DefaultLoggingFormat {
		t.Fatalf("format = %q, want %q", cfg.Logging.Format, config.DefaultLoggingFormat)
	}
}

func TestLoggingRejectsBadLevel(t *testing.T) {
	cfg := config.Config{}.WithDefaults()
	cfg.Logging.Level = "verbose"
	if err := cfg.Validate(); err == nil {
		t.Fatal("want error for logging.level=verbose")
	}
}
