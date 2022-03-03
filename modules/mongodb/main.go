package mongodb

import (
	"context"

	"git.nonamestudio.me/gjs/engine/core/converters"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Module struct {
	runtime *goja.Runtime
}

type NativeMongoDBHandler struct {
	runtime    *goja.Runtime
	connection *mongo.Client
}
type NativeMongoDBDatabaseHandler struct {
	runtime    *goja.Runtime
	connection *mongo.Client
	database   *mongo.Database
}

var ctx = context.TODO()

func (m NativeMongoDBDatabaseHandler) collection(call goja.FunctionCall) goja.Value {
	name := converters.String(m.runtime, call.Argument(0))
	col := m.database.Collection(name)
	handler := NativeMongoDBCollectionHandler{
		runtime:    m.runtime,
		connection: m.connection,
		database:   m.database,
		collection: col,
	}

	object := m.runtime.NewObject()
	_ = object.Set("findOne", handler.findOne)
	_ = object.Set("find", handler.find)
	_ = object.Set("insert", handler.insert)
	_ = object.Set("delete", handler.deleteOne)
	_ = object.Set("deleteMany", handler.deleteMany)
	_ = object.Set("update", handler.updateOne)
	_ = object.Set("updateMany", handler.deleteMany)
	return object
}

func (m Module) createClient(call goja.FunctionCall) goja.Value {
	optionsObj := call.Argument(0).ToObject(m.runtime)

	clientOptions := options.Client().ApplyURI(converters.StringOrDefault(m.runtime, optionsObj.Get("url"), "mongodb://localhost:27017/"))

	promise, resolve, reject := loop.NewPromise(m.runtime)
	go func() {

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			reject(m.runtime.ToValue(err))
			return
		}

		handler := NativeMongoDBHandler{
			runtime:    m.runtime,
			connection: client,
		}

		object := m.runtime.NewObject()
		_ = object.Set("db", handler.database)
		resolve(object)
	}()

	return m.runtime.ToValue(promise)
}

func CreateModule(vm *goja.Runtime) *goja.Object {
	mongodbClient := Module{runtime: vm}

	object := vm.NewObject()

	_ = object.Set("createClient", mongodbClient.createClient)
	_ = object.Set("objectId", mongodbClient.objId)

	return object
}
