package rotatelogs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat-go/strftime"
)

// GenerateFile creates a file name based on the pattern, the current time, and the
// rotation time.
//
// The bsase time that is used to generate the filename is truncated based
// on the rotation time.
func GenerateFile(pattern *strftime.Strftime, clock Clock, rotationTime time.Duration) string {
	now := clock.Now()
	// XXX HACK: Truncate only happens in UTC semantics, apparently.
	// observed values for truncating given time with 86400 secs:
	//
	// before truncation: 2018/06/01 03:54:54 2018-06-01T03:18:00+09:00
	// after  truncation: 2018/06/01 03:54:54 2018-05-31T09:00:00+09:00
	//
	// This is really annoying when we want to truncate in local time
	// so we hack: we take the apparent local time in the local zone,
	// and pretend that it's in UTC. do our math, and put it back to
	// the local zone
	if now.Location() != time.UTC {
		base := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
		base = base.Truncate(rotationTime)
		base = time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), base.Nanosecond(), base.Location())
		return pattern.FormatString(base)
	}
	return pattern.FormatString(now.Truncate(rotationTime))
}

// CreateFile creates a new file in the given path, creating parent directories  as necessary
func CreateFile(filename string) (*os.File, error) {
	dirname := filepath.Dir(filename)
	err := os.MkdirAll(dirname, 0755) // make sure the dir is existed, eg: ./foo/bar/baz/hello.log must make sure ./foo/bar/baz is existed
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // if we got here, then we need to create a file
	if err != nil {
		return nil, err
	}
	return file, nil
}
