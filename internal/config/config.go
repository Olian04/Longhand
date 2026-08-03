// Package config holds root config loading and subsection types (`labels`, `logging`).
//
// Load precedence when multiple sources are enabled: YAML → ENV → flags (later wins).
package config

// Config is the aggregated application configuration.
type Config struct {
	Labels Labels `yaml:"labels,omitempty"`

	Logging LoggingSection `yaml:"logging,omitempty"`
}

func (c Config) WithDefaults() Config {
	c.Labels = c.Labels.WithDefaults()

	c.Logging = c.Logging.WithDefaults()
	return c
}

func (c Config) Validate() error {
	if err := c.Labels.Validate(); err != nil {
		return err
	}

	return c.Logging.Validate()
}
