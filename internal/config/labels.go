package config

// Labels attach to slog.
type Labels map[string]string

func (l Labels) WithDefaults() Labels {
	if l == nil {
		return Labels{}
	}
	out := make(Labels, len(l))
	for k, v := range l {
		out[k] = v
	}
	return out
}

func (l Labels) Validate() error {

	return nil
}
