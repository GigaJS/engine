package loop

import (
	"github.com/dop251/goja"
	"sync"
)

var promises = make(map[*goja.Promise]bool)
var promiseMutex = sync.Mutex{}

var locks = make(map[*ExitLock]bool)
var lockMutex = sync.Mutex{}

var callbacks = make(map[int64]goja.Callable)

type ExitLock struct{}

func NewPromise(vm *goja.Runtime) (*goja.Promise, func(result interface{}), func(reason interface{})) {
	promise, resolve, reject := vm.NewPromise()
	promises[promise] = false
	return promise, func(result interface{}) {
			resolve(result)
			promiseMutex.Lock()
			promises[promise] = true
			promiseMutex.Unlock()
		}, func(reason interface{}) {
			reject(reason)
			promiseMutex.Lock()
			promises[promise] = true
			promiseMutex.Unlock()
		}
}

func NewLock(vm *goja.Runtime) func() {
	lock := &ExitLock{}

	lockMutex.Lock()
	locks[lock] = false
	lockMutex.Unlock()

	return func() {
		lockMutex.Lock()
		delete(locks, lock)
		lockMutex.Unlock()
	}
}

func Loop() {
	for {
		select {
		case timer := <-ready:
			var arguments []goja.Value
			if len(timer.call.Arguments) > 2 {
				tmp := timer.call.Arguments[2:]
				arguments = make([]goja.Value, 2+len(tmp))
				for i, value := range tmp {
					arguments[i+2] = value
				}
			} else {
				arguments = make([]goja.Value, 1)
			}

			arguments[0] = timer.call.Arguments[0]

			if fn, ok := goja.AssertFunction(arguments[0]); ok {
				_, err := fn(nil, arguments...)
				if err != nil {
					return
				}
			}

			for timer, _ := range registry {
				timer.timer.Stop()
				delete(registry, timer)
				return
			}

			if timer.interval {
				timer.timer.Reset(timer.duration)
			} else {
				delete(registry, timer)
			}
		default:
			// Escape valve!
			// If this isn't here, we deadlock...
		}

		promiseMutex.Lock()
		for promise, ended := range promises {
			if ended {
				delete(promises, promise)
			}
		}
		promiseMutex.Unlock()

		if len(registry) == 0 && len(promises) == 0 && len(locks) == 0 {
			break
		}
	}
}
