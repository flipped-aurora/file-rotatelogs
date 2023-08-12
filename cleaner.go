package rotatelogs

import "sync"

type Cleaner struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *Cleaner) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}

func (g *Cleaner) Run() {
	if g.enable {
		g.mutex.Lock()
		defer g.mutex.Unlock()
		g.fn()
	}
}
