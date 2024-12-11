package advent

// Options holds the configurable parameters for a service or feature.
type Options struct {
	Delay            int
	UpdateOnNumMoves int
	RedactSolution   bool
}

// Option is a functional option type that modifies the Options.
type Option func(*Options)

// WithDelay sets the Delay option.
func WithDelay(delay int) Option {
	return func(o *Options) {
		o.Delay = delay
	}
}

// WithUpdateOnNumMoves sets the UpdateOnNumMoves option.
func WithUpdateOnNumMoves(numMoves int) Option {
	return func(o *Options) {
		o.UpdateOnNumMoves = numMoves
	}
}

// WithRedactSolution sets the RedactSolution option.
func WithRedactSolution(redact bool) Option {
	return func(o *Options) {
		o.RedactSolution = redact
	}
}

func newRun(opts ...Option) *Options {
	// Default options
	options := &Options{
		Delay:            0,
		UpdateOnNumMoves: 0,
		RedactSolution:   false,
	}

	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	return options
}
