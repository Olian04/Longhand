package config

const (
	DefaultLoggingLevel  = "info"
	DefaultLoggingFormat = "json"

	// A CLI's stdout is its data channel: keep logs on stderr so results stay
	// pipeable. Override with logging.stream or APP_LOGGING_STREAM.
	DefaultLoggingStream = "stderr"
)
