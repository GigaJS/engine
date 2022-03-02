package core

import (
	_ "embed"
	"fmt"
	"git.nonamestudio.me/gjs/engine/core/globals"
	"git.nonamestudio.me/gjs/engine/core/loop"
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

	coreModule := &Module{Runtime: vm}
	timerModule := &loop.TimerModule{Runtime: vm}

	_ = vm.Set("setTimeout", timerModule.SetTimeout)
	_ = vm.Set("setInterval", timerModule.SetInterval)
	_ = vm.Set("clearTimeout", timerModule.ClearTimeout)
	_ = vm.Set("clearInterval", timerModule.ClearTimeout)

	_ = vm.Set("require", coreModule.Require)

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
	loop.Loop()
}
