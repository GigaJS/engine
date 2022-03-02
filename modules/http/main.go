package http

import (
	"github.com/dop251/goja"
)

type Module struct {
	runtime *goja.Runtime
}

func (m *Module) createServer(_ goja.FunctionCall) goja.Value {
	object := m.runtime.NewObject()

	server := &Server{
		runtime:     m.runtime,
		releaseLock: nil,
		// uses:    make([]routing.Handler),
		groups: make(map[string]map[string]route),
	}
	_ = object.Set("use", server.Use)
	_ = object.Set("createGroup", server.CreateGroup)
	_ = object.Set("listen", server.Listen)
	_ = object.Set("close", server.Close)

	return object
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	httpModule := &Module{runtime: vm}
	object := vm.NewObject()
	_ = object.Set("createServer", httpModule.createServer)
	_ = object.Set("request", httpModule.request)
	return object
}
