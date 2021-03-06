package loop

import (
	"github.com/dop251/goja"
	"time"
)

type _timer struct {
	timer    *time.Timer
	duration time.Duration
	interval bool
	call     goja.FunctionCall
}

var registry = make(map[*_timer]bool)
var ready = make(chan *_timer)

func (m *TimerModule) newTimer(
	call goja.FunctionCall,
	interval bool,
) (*_timer, goja.Value) {
	delay := call.Argument(1).ToInteger()
	if 0 >= delay {
		delay = 1
	}

	timer := &_timer{
		duration: time.Duration(delay) * time.Millisecond,
		call:     call,
		interval: interval,
	}
	registry[timer] = false

	timer.timer = time.AfterFunc(timer.duration, func() {
		ready <- timer
	})

	value := m.Runtime.ToValue(timer)
	return timer, value
}

func (m *TimerModule) SetTimeout(call goja.FunctionCall) goja.Value {
	_, value := m.newTimer(call, false)
	return value
}

func (m *TimerModule) SetInterval(call goja.FunctionCall) goja.Value {
	_, value := m.newTimer(call, true)
	return value
}

func (m *TimerModule) ClearTimeout(call goja.FunctionCall) goja.Value {
	timer := call.Argument(0).Export()
	if timer, ok := timer.(*_timer); ok {
		timer.timer.Stop()
		delete(registry, timer)
	}

	return goja.Undefined()
}
