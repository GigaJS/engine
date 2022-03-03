package converters

import (
	"github.com/dop251/goja"
	"reflect"
	"strings"
)

func DynamicArrayToBytes(a goja.DynamicArray) []byte {
	r := make([]byte, a.Len())

	for i := 0; i <= a.Len(); i++ {
		item := a.Get(i)
		r[i] = byte(item.ToInteger())
	}

	return r
}

func replaceFirst(str string, replacement byte) string {
	out := []byte(str)
	out[0] = replacement
	return string(out)
}

func IsBuffer(d interface{}) bool {
	if arr, ok := d.(goja.DynamicArray); ok {
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

func InterfaceToObject(vm *goja.Runtime, v interface{}) *goja.Object {
	var r = vm.NewObject()

	t := reflect.Indirect(reflect.ValueOf(v))
	ty := reflect.TypeOf(v)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fn := ty.Field(i).Name
		fn = replaceFirst(fn, []byte(strings.ToLower(string(fn[0])))[0])

		switch f.Kind() {
		case reflect.Int64:
			_ = r.Set(fn, f.Int())
		case reflect.String:
			_ = r.Set(fn, f.String())
		case reflect.Bool:
			_ = r.Set(fn, f.Bool())
		case reflect.Func:
			if f.IsNil() {
				break
			}

			_ = r.Set(fn, f.Interface())
		}
	}

	return r
}

func StringOrDefault(vm *goja.Runtime, value goja.Value, def string) string {
	if value == nil {
		return def
	}

	return String(vm, value)
}

func String(vm *goja.Runtime, value goja.Value) string {
	if value == nil {
		panic(vm.ToValue("String must be string"))
	}

	if value.ExportType().Kind() == reflect.String {
		return value.Export().(string)
	}

	panic(vm.ToValue("String must be string"))
}

func NumberInt64(vm *goja.Runtime, value goja.Value) int64 {
	if !IsPresent(vm, value) {
		panic(vm.ToValue("Must be a number"))
	}

	switch value.ExportType().Kind() {
	case reflect.Int, reflect.Int64:
		return value.ToInteger()
	case reflect.Float64, reflect.Float32:
		return int64(value.ToFloat())
	}

	panic(vm.ToValue("Number type unsupported"))
}

func NumberInt(vm *goja.Runtime, value goja.Value) int32 {
	if !IsPresent(vm, value) {
		panic(vm.ToValue("Must be a number"))
	}

	switch value.ExportType().Kind() {
	case reflect.Int, reflect.Int64:
		return int32(value.ToInteger())
	case reflect.Float64, reflect.Float32:
		return int32(value.ToFloat())
	}

	panic(vm.ToValue("Number type unsupported"))
}

func IsPresent(vm *goja.Runtime, val goja.Value) bool {
	return val != nil && !goja.IsNull(val) && !goja.IsUndefined(val) && !goja.IsNaN(val)
}

func IsPresentInObject(vm *goja.Runtime, object *goja.Object, key string) bool {
	val := object.Get(key)
	return val != nil && !goja.IsNull(val) && !goja.IsUndefined(val) && !goja.IsNaN(val)
}
