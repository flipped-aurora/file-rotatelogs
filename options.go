package rotatelogs

import "time"

// Option is used to pass optional arguments to
// the Rotate constructor
type Option func(*Rotate)

// WithClock creates a new Option that sets a clock
// that the Rotate object will use to determine
// the current time.
//
// By default, rotatelogs.Local, which returns the
// current time in the local time zone, is used. If you
// would rather use UTC, use rotatelogs.UTC as the argument
// to this option, and pass it to the constructor.
func WithClock(c Clock) Option {
	return func(rotate *Rotate) {
		rotate.clock = c
	}
}

// WithLocation creates a new Option that sets up a
// "Clock" interface that the Rotate object will use
// to determine the current time.
//
// This optin works by always returning the in the given
// location.
func WithLocation(location *time.Location) Option {
	return func(rotate *Rotate) {
		rotate.clock = clock(func() time.Time {
			return time.Now().In(location)
		})
	}
}

// WithMaxAge creates a new Option that sets the
// max age of a log file before it gets purged from
// the file system.
func WithMaxAge(age time.Duration) Option {
	return func(rotate *Rotate) {
		rotate.maxAge = age
	}
}

// WithRotationTime creates a new Option that sets the
// time between rotation.
func WithRotationTime(time time.Duration) Option {
	return func(rotate *Rotate) {
		rotate.rotationTime = time
	}
}
