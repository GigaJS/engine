package globals

import (
	"github.com/dop251/goja"
	"os"
)

type ProcessModule struct {
	runtime *goja.Runtime
}

func (p *ProcessModule) cwd() goja.Value {
	wd, err := os.Getwd()
	if err != nil {
		panic(p.runtime.NewGoError(err))
		return goja.Undefined()
	}

	return p.runtime.ToValue(wd)
}

func (p *ProcessModule) env() goja.Value {
	r := vm.NewObject()

	for _, e := range os.Environ() {
        pair := strings.SplitN(e, "=", 2)
        _ = r.Set(pair[0], pair[1])
    }

	return p.runtime.ToValue(r)
}

func RegisterProcess(vm *goja.Runtime) {
	p := &ProcessModule{runtime: vm}

	o := vm.NewObject()
	_ = o.Set("cwd", p.cwd)
	_ = o.Set("env", p.env)

	_ = vm.GlobalObject().Set("process", o)
}
