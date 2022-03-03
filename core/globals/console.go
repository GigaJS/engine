package globals

import (
	"encoding/hex"
	"fmt"
	"git.nonamestudio.me/gjs/engine/core/converters"
	"github.com/dop251/goja"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type ConsoleModule struct {
	runtime *goja.Runtime
}

type formatOptions struct {
	PropName     *string
	StringQuotes bool
}

func formatValue(v goja.Value, opts formatOptions) string {
	if v == nil || goja.IsUndefined(v) {
		return "undefined"
	}

	if goja.IsNull(v) {
		return "null"
	}

	switch v.ExportType().Kind() {
	case reflect.String:
		if opts.StringQuotes {
			return fmt.Sprintf("\"%s\"", v.String())
		} else {
			return v.String()
		}
	case reflect.Int64:
		return strconv.Itoa(int(v.ToInteger()))
	case reflect.Bool:
		if v.ToBoolean() {
			return "true"
		} else {
			return "false"
		}
	default:
		if _, ok := goja.AssertFunction(v); ok {
			if opts.PropName == nil {
				return "[Function null]"
			}
			return fmt.Sprintf("[Function %s]", *opts.PropName)
		} else if o, ok := v.(*goja.Object); ok {
			t := ""

			if o.ClassName() == "Array" {
				if len(o.Keys()) == 0 {
					t = "[]"
					return t
				}

				t += "[ "
				for i, k := range o.Keys() {
					t += formatValue(o.Get(k), formatOptions{StringQuotes: true})

					if i != len(o.Keys())-1 {
						t += ", "
					}
				}
				t += " ]"

				return t
			}

			if err, ok := v.Export().(error); ok {
				return fmt.Sprintf("Error: %s", err.Error())
			}

			if len(o.Keys()) == 0 {

				if converters.IsPresentInObject(o, "message") {

					message := o.Get("message").String()
					name := o.Get("name").String()
					stack := o.Get("stack").String()

					strings.IndexRune(stack, '\n')

					t = fmt.Sprintf("%s: %s\n%s", name, message, stack)
				} else {
					t = "{}"
				}
			} else {

				if IsBuffer(o) {
					bytes := o.Get("data").Export().([]uint8)
					lenBytes := len(bytes)
					if lenBytes < 1 {
						return "<Buffer>"
					} else {
						var str string
						t := make([]string, 50)

						for i := 0; i < 50; i++ {
							t[i] = hex.EncodeToString([]byte{bytes[i]})
						}

						if lenBytes > 50 {
							str = fmt.Sprintf("%s ... %d more bytes", strings.Join(t, " "), lenBytes-50)
						} else {
							str = strings.Join(t, " ")
						}

						return fmt.Sprintf("<Buffer %s>", str)
					}
				}

				t += "{ "

				for i, k := range o.Keys() {
					prop := o.Get(k)
					t += k + ": " + formatValue(prop, formatOptions{PropName: &k, StringQuotes: true})

					if i != len(o.Keys())-1 {
						t += ", "
					}
				}

				t += " }"
			}
			return t
		} else {
			return "unknown"
		}
	}
}

func logRaw(call goja.FunctionCall, file *os.File) goja.Value {
	if len(call.Arguments) == 0 {
		return goja.Undefined()
	}

	var r = ""

	for i, arg := range call.Arguments {
		r += formatValue(arg, formatOptions{StringQuotes: false})

		if i != len(call.Arguments)-1 {
			r += " "
		}
	}

	_, _ = file.WriteString(r + "\n")

	return goja.Undefined()
}

func (c *ConsoleModule) log(call goja.FunctionCall) goja.Value {
	return logRaw(call, os.Stdout)
}

func (c *ConsoleModule) error(call goja.FunctionCall) goja.Value {
	return logRaw(call, os.Stderr)
}

func RegisterConsole(vm *goja.Runtime) {
	c := &ConsoleModule{runtime: vm}

	o := vm.NewObject()
	_ = o.Set("log", c.log)
	_ = o.Set("error", c.error)

	_ = vm.Set("console", o)
}
