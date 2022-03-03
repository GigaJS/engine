package mongodb

import (
	"git.nonamestudio.me/gjs/engine/core/converters"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongodbType = goja.NewSymbol("mongodb-type")

func getMongoDBTag(vm *goja.Runtime, obj *goja.Object) (bool, string) {

	tag := obj.GetSymbol(mongodbType)

	if converters.IsPresent(vm, tag) {
		data := converters.String(vm, tag)
		return true, data
	}

	return false, ""
}

func isMongoDBObj(vm *goja.Runtime, obj *goja.Object) bool {
	has, _ := getMongoDBTag(vm, obj)
	return has
}

func mongoDBObjExtract(vm *goja.Runtime, obj *goja.Object) interface{} {
	_, data := getMongoDBTag(vm, obj)
	switch data {
	case "1":
		id, err := primitive.ObjectIDFromHex(converters.String(vm, obj.Get("id")))
		if err != nil {
			panic(vm.ToValue(err.Error()))
		}
		return id
	}

	return nil
}
