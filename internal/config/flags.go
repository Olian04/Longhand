package config

// FlagOverrides are applied last (highest precedence). Nil fields leave values unchanged.
type FlagOverrides struct {
	LoggingLevel  *string
	LoggingFormat *string
	LoggingStream *string
}

func applyFlags(cfg *Config, f FlagOverrides) {

	if f.LoggingLevel != nil {
		cfg.Logging.Level = *f.LoggingLevel
	}
	if f.LoggingFormat != nil {
		cfg.Logging.Format = *f.LoggingFormat
	}
	if f.LoggingStream != nil {
		cfg.Logging.Stream = *f.LoggingStream
	}
}
