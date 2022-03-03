package mongodb

import (
	"git.nonamestudio.me/gjs/engine/core/converters"
	"github.com/dop251/goja"
)

func (m Module) objId(call goja.FunctionCall) goja.Value {
	id := converters.String(m.runtime, call.Argument(0))

	object := m.runtime.NewObject()

	object.SetSymbol(mongodbType, "1")
	object.Set("id", id)

	return object
}
