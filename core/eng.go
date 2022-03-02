package core

import (
	_ "embed"
	"fmt"
	"git.nonamestudio.me/gjs/engine/core/globals"
	"github.com/dop251/goja"
	"os"
)

type Engine struct {
	vm *goja.Runtime
}

func CreateGJSEngine() *Engine {
	vm := goja.New()

	defer func() {
		if err := recover(); err != nil {
			_, _ = os.Stderr.WriteString(err.(string))
		}
	}()

	m := &Module{Runtime: vm}

	_ = vm.Set("setTimeout", m.SetTimeout)
	_ = vm.Set("setInterval", m.SetInterval)
	_ = vm.Set("clearTimeout", m.ClearTimeout)
	_ = vm.Set("clearInterval", m.ClearTimeout)

	_ = vm.Set("require", m.Require)

	RegisterCompatibility(vm)

	globals.RegisterConsole(vm)
	globals.RegisterBuffer(vm)
	globals.RegisterProcess(vm)
	globals.RegisterUrl(vm)

	eng := Engine{
		vm: vm,
	}

	return &eng
}

func (e Engine) ExecuteFromString(script string) {
	_, err := e.vm.RunString(script)
	if err != nil {
		if jse, ok := err.(*goja.Exception); ok {
			_, _ = os.Stderr.WriteString(jse.String())
		} else {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Error: %s", err.Error()))
		}
	}

}

func (e Engine) Start() {
	Loop()
}
