package mongodb

import (
	"git.nonamestudio.me/gjs/engine/core/converters"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongodbType = goja.NewSymbol("mongodb-type")

func getMongoDBTag(vm *goja.Runtime, obj *goja.Object) *string {
	var st string

	tag := obj.GetSymbol(mongodbType)

	if converters.IsPresent(vm, tag) {
		data := converters.String(vm, tag)
		return &data
	}

	return &st
}

func isMongoDBObj(vm *goja.Runtime, obj *goja.Object) bool {
	return getMongoDBTag(vm, obj) != nil
}

func mongoDBObjExtract(vm *goja.Runtime, obj *goja.Object) interface{} {
	switch *getMongoDBTag(vm, obj) {
	case "1":
		 id, err := primitive.ObjectIDFromHex(converters.String(vm, obj.Get("id")))
		 if err != nil {
			 panic(vm.ToValue(err.Error()))
		 }
		 return id
	}

	return nil
}
