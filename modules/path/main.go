package path

import (
	"github.com/dop251/goja"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

type Module struct {
	runtime *goja.Runtime
}

func (p *Module) join(call goja.FunctionCall) goja.Value {
	parts := make([]string, len(call.Arguments))

	for i, v := range call.Arguments {
		if v.ExportType().Kind() != reflect.String {
			panic(p.runtime.NewTypeError("path part must be a string"))
			return goja.Undefined()
		}

		parts[i] = v.String()
	}

	return p.runtime.ToValue(path.Join(parts...))
}

func (p *Module) dirname(call goja.FunctionCall) goja.Value {
	dirValue := call.Argument(0)

	if dirValue.ExportType().Kind() != reflect.String {
		panic(p.runtime.NewTypeError("path must be a string"))
		return goja.Undefined()
	}

	dir := dirValue.String()
	return p.runtime.ToValue(filepath.Dir(dir))
}

func (p *Module) extname(call goja.FunctionCall) goja.Value {
	fileValue := call.Argument(0)

	if fileValue.ExportType().Kind() != reflect.String {
		panic(p.runtime.NewTypeError("filename must be a string"))
		return goja.Undefined()
	}

	file := fileValue.String()
	return p.runtime.ToValue(filepath.Ext(file))
}

func (p *Module) isAbsolute(call goja.FunctionCall) goja.Value {
	fileValue := call.Argument(0)

	if fileValue.ExportType().Kind() != reflect.String {
		p.runtime.Interrupt(p.runtime.NewTypeError("path  must be a string"))
		return goja.Undefined()
	}

	return p.runtime.ToValue(filepath.IsAbs(fileValue.String()))
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	p := &Module{runtime: vm}

	object := vm.NewObject()
	switch runtime.GOOS {
	case "windows":
		_ = object.Set("sep", "\\")
	case "linux", "darwin":
		_ = object.Set("sep", "/")
	}

	_ = object.Set("join", p.join)
	_ = object.Set("dirname", p.dirname)
	_ = object.Set("extname", p.extname)
	_ = object.Set("isAbsolute", p.isAbsolute)
	return object
}
