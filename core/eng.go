package core

import (
	_ "embed"
	"fmt"
	"git.nonamestudio.me/gjs/engine/core/globals"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"os"
)

type ModuleRegistrar func(engine *Engine) goja.Value

type Engine struct {
	Runtime       *goja.Runtime
	nativeModules map[string]ModuleRegistrar
}

func CreateGJSEngine() *Engine {
	vm := goja.New()

	defer func() {
		if err := recover(); err != nil {
			_, _ = os.Stderr.WriteString(err.(string))
		}
	}()

	eng := Engine{
		Runtime:       vm,
		nativeModules: map[string]ModuleRegistrar{},
	}

	coreModule := &Module{Runtime: vm, Engine: &eng}
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

	return &eng
}

func (e Engine) ExecuteFromString(script string) {
	_, err := e.Runtime.RunString(script)
	if err != nil {
		if jse, ok := err.(*goja.Exception); ok {
			_, _ = os.Stderr.WriteString(jse.String())
		} else {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Error: %s", err.Error()))
		}
	}

}

func (e Engine) RegisterModule(module string, registerFunction ModuleRegistrar) {
	e.nativeModules[module] = registerFunction
}

func (e Engine) Start() {
	loop.Loop()
}
