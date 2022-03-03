package mongodb

import (
	"reflect"

	"git.nonamestudio.me/gjs/engine/core/converters"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toValue(vm *goja.Runtime, data interface{}) goja.Value {
	// TODO: Array support
	if val, ok := data.(bson.M); ok {
		obj := vm.NewObject()

		for key, value := range val {
			_ = obj.Set(key, toValue(vm, value))
		}

		return obj
	}

	if objectId, ok := data.(primitive.ObjectID); ok {
		return vm.ToValue(objectId.Hex())
	}

	return vm.ToValue(data)
}

func fromValue(vm *goja.Runtime, data goja.Value) interface{} {
	// TODO: Array support

	if !converters.IsPresent(vm, data) {
		return nil
	}


	if objData, ok := data.(*goja.Object); ok {

		if isMongoDBObj(vm, objData) {
			return mongoDBObjExtract(vm, objData)
		}

		result := bson.M{}
		for _, key := range objData.Keys() {
			v := objData.Get(key)
			result[key] = fromValue(vm, v)
		}
		return result
	}


	switch data.ExportType().Kind() {
	case reflect.Bool:
		return data.ToBoolean()
	case reflect.String:
		return data.String()
	case reflect.Int:
		return data.ToInteger()
	case reflect.Int64:
		return data.ToInteger()
	case reflect.Float64:
		return data.ToFloat()
	case reflect.Float32:
		return data.ToFloat()
	}
	return nil
}
