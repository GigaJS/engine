package globals

import (
	"github.com/dop251/goja"
	"reflect"
)

type BufferModule struct {
	runtime *goja.Runtime
}

var bufferSymbol = goja.NewSymbol("buffer")

func (b *BufferModule) bufferConstructor(call goja.ConstructorCall) *goja.Object {
	input := call.Argument(0).Export()
	var data []byte

	if v, ok := input.(string); ok {
		data = []byte(v)
	} else if v, ok := input.([]uint8); ok {
		data = v
	} else if v, ok := input.(goja.ArrayBuffer); ok {
		data = v.Bytes()
	} else {
		panic(b.runtime.ToValue("invalid input type"))
		return nil
	}

	o := b.runtime.NewObject()
	_ = o.Set("data", data)
	_ = o.SetSymbol(bufferSymbol, b.runtime.ToValue(true))

	_ = o.Set("toString", func(c goja.FunctionCall) goja.Value {
		encoding := "utf8"

		encodingValue := c.Argument(0)
		if !goja.IsUndefined(encodingValue) {
			if encodingValue.ExportType().Kind() != reflect.String {
				panic(b.runtime.ToValue("Encoding must be a string"))
				return goja.Undefined()
			}

			encoding = encodingValue.Export().(string)
		}

		switch encoding {
		case "utf8", "utf-8":
			return b.runtime.ToValue(string(data))
		}

		return goja.Undefined()
	})

	return o
}

func (b *BufferModule) bufferFrom(call goja.FunctionCall) goja.Value {
	d := call.Argument(0)

	if d.ExportType().Kind() == reflect.String {
		s := d.String()
		v := b.runtime.ToValue([]byte(s))
		return v
	} else {
		panic(b.runtime.NewTypeError("invalid argument type"))
	}

	return goja.Undefined()
}

func (b *BufferModule) isBuffer(call goja.FunctionCall) goja.Value {
	d := call.Argument(0)

	if obj, ok := d.(*goja.Object); ok {
		if obj.GetSymbol(bufferSymbol) != nil {
			return b.runtime.ToValue(true)
		}
		return b.runtime.ToValue(false)
	} else if arr, ok := d.(goja.DynamicArray); ok {
		for i := 0; i <= arr.Len(); i++ {
			arrItem := arr.Get(i)
			if arrItem.ExportType().Kind() != reflect.Int64 {
				return b.runtime.ToValue(false)
			}
		}

		return b.runtime.ToValue(true)
	} else {
		return b.runtime.ToValue(false)
	}
}

func IsBuffer(value goja.Value) bool {

	if obj, ok := value.(*goja.Object); ok {
		if obj.GetSymbol(bufferSymbol) != nil {
			return true
		}
		return false
	} else if arr, ok := value.(goja.DynamicArray); ok {
		for i := 0; i <= arr.Len(); i++ {
			arrItem := arr.Get(i)
			if arrItem.ExportType().Kind() != reflect.Int64 {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func RegisterBuffer(vm *goja.Runtime) {
	b := &BufferModule{runtime: vm}
	constructor := vm.ToValue(b.bufferConstructor)

	if obj, ok := constructor.(*goja.Object); ok {
		_ = obj.Set("isBuffer", b.isBuffer)
	}

	_ = vm.GlobalObject().Set("Buffer", constructor)
	/*
		f := vm.ToValue(b.bufferFrom).(*goja.Object)
		_ = f.Set("from", b.bufferFrom)
		_ = f.Set("isBuffer", b.isBuffer)

		_ = vm.Set("Buffer", f)*/
}
