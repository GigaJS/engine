package core

import (
	"fmt"
	gjsColors "git.nonamestudio.me/gjs/engine/modules/colors"
	gjsFs "git.nonamestudio.me/gjs/engine/modules/fs"
	gjsHttp "git.nonamestudio.me/gjs/engine/modules/http"
	gjsMongoDb "git.nonamestudio.me/gjs/engine/modules/mongodb"
	gjsPath "git.nonamestudio.me/gjs/engine/modules/path"
	gjsUrl "git.nonamestudio.me/gjs/engine/modules/url"
	"github.com/dop251/goja"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const ModuleLocationNative = 1
const ModuleLocationPackage = 2
const ModuleLocationRelative = 3

type Module struct {
	Runtime *goja.Runtime
}

var cachedModules = map[string]*goja.Object{}

func isRelativePath(str string) bool {
	return strings.HasPrefix(str, "./") || strings.HasPrefix(str, "../")
}

func moduleExists(name string) (exist, native bool, moduleLocation uint8) {
	switch name {
	case "fs", "url", "http", "https", "path", "colors", "mongodb":
		return true, true, ModuleLocationNative
	}

	if isRelativePath(name) {
		return true, false, ModuleLocationRelative
	} else {
		modulePath := path.Join("node_modules", name)

		_, err := os.Stat(modulePath)
		if os.IsNotExist(err) {
			return false, false, 0
		}

		_, err = os.Stat(path.Join(modulePath, "index.js"))
		if os.IsNotExist(err) {
			return false, false, 0
		}

		return true, false, ModuleLocationPackage
	}
}

func (m *Module) getCurrentModulePath() string {
	var buf [2]goja.StackFrame
	frames := m.Runtime.CaptureCallStack(2, buf[:0])
	if len(frames) < 2 {
		return "."
	}
	return path.Dir(frames[1].SrcName())
}

func (m *Module) importModule(originalPath, file string) goja.Value {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		m.Runtime.Interrupt(fmt.Sprintf("Cannot find module '%s'", originalPath))
		return goja.Undefined()
	}

	contentBytes, _ := ioutil.ReadFile(file)
	to := string(contentBytes)
	if filepath.Ext(file) == ".json" {
		to = fmt.Sprintf("module.exports = %s", string(contentBytes))
	}

	script, err := goja.Compile("", to, false)
	if err != nil {
		panic(err.Error())
		return goja.Undefined()
	}

	v, err := m.Runtime.RunProgram(script)

	if err != nil {
		panic(err.Error())
		return goja.Undefined()
	}

	return v
}

func (m *Module) Require(call goja.FunctionCall) goja.Value {
	moduleValue := call.Argument(0)
	if moduleValue.ExportType().Name() != "string" {
		m.Runtime.Interrupt(m.Runtime.NewTypeError("module must be a string"))
		return goja.Undefined()
	}

	moduleName := moduleValue.String()

	exist, native, _ := moduleExists(moduleName)
	if !exist {
		m.Runtime.Interrupt(fmt.Sprintf("Cannot find module '%s'", moduleName))
		return goja.Undefined()
	}

	var o *goja.Object

	if native {
		switch moduleName {
		case "fs":
			o = gjsFs.CreateModule(m.Runtime)
		case "url":
			o = gjsUrl.CreateModule(m.Runtime)
		case "http":
			o = gjsHttp.CreateModule(m.Runtime)
		case "path":
			o = gjsPath.CreateModule(m.Runtime)
		case "colors":
			o = gjsColors.CreateModule(m.Runtime)
		case "mongodb":
			o = gjsMongoDb.CreateModule(m.Runtime)
		}
	} else {
		if filepath.Ext(moduleName) == "" {
			moduleName += ".js"
		}

		op := moduleName
		moduleName = path.Clean(moduleName)

		var start string
		if path.IsAbs(op) {
			start = "/"
		} else {
			start = m.getCurrentModulePath()
		}

		moduleName = path.Join(start, moduleName)

		if cachedModules[moduleName] != nil {
			return cachedModules[moduleName]
		}

		return m.importModule(op, moduleName)
	}

	cachedModules[moduleName] = o

	return o
}
