package rotatelogs

import (
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/strftime"
)

// Rotate represents a log file that gets
// automatically rotated as you write to it.
type Rotate struct {
	clock        Clock              // 时间
	out          *os.File           // 文件句柄
	business     *os.File           // 文件句柄
	mutex        *sync.RWMutex      // 读写锁
	maxAge       time.Duration      // 最大保存时间
	pattern      *strftime.Strftime // 时间格式
	rotationTime time.Duration      // 旋转时间
}

// New creates a new Rotate object. A log filename pattern
// must be passed. Optional `Option` parameters may be passed
func New(p string, options ...Option) (*Rotate, error) {
	pattern, err := strftime.New(p)
	if err != nil {
		return nil, err
	}
	rotate := &Rotate{
		clock:        Local,
		mutex:        new(sync.RWMutex),
		pattern:      pattern,
		rotationTime: 24 * time.Hour,
	}
	for i := 0; i < len(options); i++ {
		options[i](rotate)
	}
	return rotate, nil
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (r *Rotate) Write(bytes []byte) (n int, err error) {
	r.mutex.Lock() // Guard against concurrent writes
	defer func() {
		r.mutex.Unlock()
		r.Close()
	}()
	var out io.Writer
	if strings.Contains(string(bytes), "business") {
		var compile *regexp.Regexp
		compile, err = regexp.Compile(`{"business": "([^,]+)"}`)
		if err != nil {
			return 0, err
		}
		if compile.Match(bytes) {
			finds := compile.FindSubmatch(bytes)
			business := string(finds[len(finds)-1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
			out, err = r.getBusinessWriter(business)
			if err != nil {
				return 0, err
			}
			return out.Write(bytes)
		}
		compile, err = regexp.Compile(`"business": "([^,]+)"`)
		if err != nil {
			return 0, err
		}
		if compile.Match(bytes) {
			finds := compile.FindSubmatch(bytes)
			business := string(finds[len(finds)-1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
			out, err = r.getBusinessWriter(business)
			if err != nil {
				return 0, err
			}
			return out.Write(bytes)
		}
	}
	out, err = r.getWriter()
	if err != nil {
		return 0, err
	}
	return out.Write(bytes)
}

// getBusinessWriter 获取 business io.Writer
func (r *Rotate) getBusinessWriter(business string) (io.Writer, error) {
	var pattern *strftime.Strftime
	slice := strings.Split(r.pattern.Pattern(), "/")
	if slice[len(slice)-2] != business {
		slice = append(slice[:len(slice)-1], business, slice[len(slice)-1])
		pattern, _ = strftime.New(strings.Join(slice, "/"))
	}
	filename := GenerateFile(pattern, r.clock, r.rotationTime)
	out, err := CreateFile(filename)
	if err != nil {
		return nil, err
	}
	r.business = out
	return out, nil
}

// getWriter 获取 io.Writer
func (r *Rotate) getWriter() (io.Writer, error) {
	filename := GenerateFile(r.pattern, r.clock, r.rotationTime)
	out, err := CreateFile(filename)
	if err != nil {
		return nil, err
	}
	_ = r.out.Close()
	r.out = out
	return out, nil
}

func (r *Rotate) Close() {
	r.mutex.Lock() // Guard against concurrent writes
	defer r.mutex.Unlock()
	if r.out != nil {
		_ = r.out.Close()
		r.out = nil
	}
	if r.business != nil {
		_ = r.business.Close()
		r.business = nil
	}
}
