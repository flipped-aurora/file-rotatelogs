package rotatelogs

import "time"

// Clock is the interface used by the Rotate
// object to determine the current time
type Clock interface {
	Now() time.Time
}
type clock func() time.Time

// UTC is an object satisfying the Clock interface, which
// returns the current time in UTC
var UTC = clock(func() time.Time { return time.Now().UTC() })

// Local is an object satisfying the Clock interface, which
// returns the current time in the local timezone
var Local = clock(time.Now)

func (c clock) Now() time.Time {
	return c()
}
