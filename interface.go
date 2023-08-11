package rotatelogs

import (
	"os"
	"sync"
	"time"

	strftime "github.com/lestrrat-go/strftime"
)

type Handler interface {
	Handle(Event)
}

type HandlerFunc func(Event)

type Event interface {
	Type() EventType
}

type EventType int

const (
	InvalidEventType EventType = iota
	FileRotatedEventType
)

type FileRotatedEvent struct {
	prev    string // previous filename
	current string // current, new filename
}

// Rotate represents a log file that gets
// automatically rotated as you write to it.
type Rotate struct {
	clock         Clock
	curFn         string
	curBaseFn     string
	globPattern   string
	generation    int
	linkName      string // 软链接名称
	maxAge        time.Duration
	mutex         *sync.RWMutex      // 读写锁
	eventHandler  Handler            // 事件处理
	outFh         *os.File           // 文件句柄
	pattern       *strftime.Strftime // 时间格式
	rotationTime  time.Duration      // 旋转时间
	rotationSize  int64              // 旋转大小
	rotationCount uint               // 旋转次数
	forceNewFile  bool               // 强制新文件
}

// Clock is the interface used by the Rotate
// object to determine the current time
type Clock interface {
	Now() time.Time
}
type clockFn func() time.Time

// UTC is an object satisfying the Clock interface, which
// returns the current time in UTC
var UTC = clockFn(func() time.Time { return time.Now().UTC() })

// Local is an object satisfying the Clock interface, which
// returns the current time in the local timezone
var Local = clockFn(time.Now)
