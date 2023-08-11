package rotatelogs

import (
	"time"
)

// Option is used to pass optional arguments to
// the Rotate constructor
type Option func(*Rotate)

const (
	optkeyClock         = "clock"
	optkeyHandler       = "handler"
	optkeyLinkName      = "link-name"
	optkeyMaxAge        = "max-age"
	optkeyRotationTime  = "rotation-time"
	optkeyRotationSize  = "rotation-size"
	optkeyRotationCount = "rotation-count"
	optkeyForceNewFile  = "force-new-file"
)

// WithClock creates a new Option that sets a clock
// that the Rotate object will use to determine
// the current time.
//
// By default rotatelogs.Local, which returns the
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
		rotate.clock = clockFn(func() time.Time {
			return time.Now().In(location)
		})
	}
}

// WithLinkName creates a new Option that sets the
// symbolic link name that gets linked to the current
// file name being used.
func WithLinkName(name string) Option {
	return func(rotate *Rotate) {
		rotate.linkName = name
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

// WithRotationSize creates a new Option that sets the
// log file size between rotation.
func WithRotationSize(size int64) Option {
	return func(rotate *Rotate) {
		rotate.rotationSize = size
	}
}

// WithRotationCount creates a new Option that sets the
// number of files should be kept before it gets
// purged from the file system.
func WithRotationCount(count uint) Option {
	return func(rotate *Rotate) {
		rotate.rotationCount = count
	}
}

// WithHandler creates a new Option that specifies the
// Handler object that gets invoked when an event occurs.
// Currently `FileRotated` event is supported
func WithHandler(h Handler) Option {
	return func(rotate *Rotate) {
		rotate.eventHandler = h
	}
}

// ForceNewFile ensures a new file is created every time New()
// is called. If the base file name already exists, an implicit
// rotation is performed
func ForceNewFile() Option {
	return func(rotate *Rotate) {
		rotate.forceNewFile = true
	}
}
