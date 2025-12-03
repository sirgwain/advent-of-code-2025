package advent

// Options holds the configurable parameters for a service or feature.
type Options struct {
	Delay int
}

// Option is a functional option type that modifies the Options.
type Option func(*Options)

// WithDelay sets the Delay option.
func WithDelay(delay int) Option {
	return func(o *Options) {
		o.Delay = delay
	}
}

func NewRun(opts ...Option) *Options {
	// Default options
	options := &Options{
		Delay: 0,
	}

	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	return options
}
