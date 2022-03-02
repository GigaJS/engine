package colors

import (
	"github.com/dop251/goja"
)

func isColorsEnabled() bool {
	return true
}

func Color(colorString string) string {
	if !isColorsEnabled() {
		return ""
	}
	return colorString
}

var black = Color("\033[1;30m")
var red = Color("\033[1;31m")
var green = Color("\033[1;32m")
var yellow = Color("\033[1;33m")
var purple = Color("\033[1;34m")
var magenta = Color("\033[1;35m")
var teal = Color("\033[1;36m")
var white = Color("\033[1;37m")

var brightRed = Color("\033[31;1m")
var brightGreen = Color("\033[32;1m")
var brightYellow = Color("\033[32;1m")
var brightBlue = Color("\033[34;1m")
var brightMagenta = Color("\033[35;1m")
var brightCyan = Color("\033[36;1m")

///
/// Formatting
///

var bold = Color("\033[1m")
var reset = Color("\033[0m")
var underline = Color("\033[4m")
var reverse = Color("\033[7m")

func ReloadColors(vm *goja.Runtime, colorObject *goja.Object) {

	black = Color("\033[1;30m")
	red = Color("\033[1;31m")
	green = Color("\033[1;32m")
	yellow = Color("\033[1;33m")
	purple = Color("\033[1;34m")
	magenta = Color("\033[1;35m")
	teal = Color("\033[1;36m")
	white = Color("\033[1;37m")
	brightRed = Color("\033[31;1m")
	brightGreen = Color("\033[32;1m")
	brightYellow = Color("\033[32;1m")
	brightBlue = Color("\033[34;1m")
	brightMagenta = Color("\033[35;1m")
	brightCyan = Color("\033[36;1m")
	bold = Color("\033[1m")
	underline = Color("\033[4m")
	reverse = Color("\033[7m")
	reset = Color("\033[0m")

	_ = colorObject.Set("black", black)
	_ = colorObject.Set("red", red)
	_ = colorObject.Set("green", green)
	_ = colorObject.Set("yellow", yellow)
	_ = colorObject.Set("purple", purple)
	_ = colorObject.Set("magenta", magenta)
	_ = colorObject.Set("teal", teal)
	_ = colorObject.Set("white", white)
	_ = colorObject.Set("brightRed", brightRed)
	_ = colorObject.Set("brightGreen", brightGreen)
	_ = colorObject.Set("brightYellow", brightYellow)
	_ = colorObject.Set("brightBlue", brightBlue)
	_ = colorObject.Set("brightMagenta", brightMagenta)
	_ = colorObject.Set("brightCyan", brightCyan)
	_ = colorObject.Set("bold", bold)
	_ = colorObject.Set("underline", underline)
	_ = colorObject.Set("reverse", reverse)
	_ = colorObject.Set("reset", reset)
}

func CreateModule(vm *goja.Runtime) *goja.Object {

	object := vm.NewObject()
	_ = object.Set("reload", func(call goja.FunctionCall) goja.Value {
		ReloadColors(vm, object)
		return goja.Undefined()
	})
	ReloadColors(vm, object)
	return object
}
