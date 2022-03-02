package mongodb

import (
	"context"
	"git.nonamestudio.me/gjs/engine/core/converters"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type NativeMongoDBCollectionHandler struct {
	runtime    *goja.Runtime
	connection *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

var ctx = context.TODO()

func (h NativeMongoDBCollectionHandler) find(call goja.FunctionCall) goja.Value {
	singleRes := h.collection.FindOne(ctx, bson.D{})

	if singleRes.Err() != nil {
		panic(singleRes.Err())
	}

	data := bson.M{}

	err := singleRes.Decode(&data)
	if err != nil {
		panic(err)
	}

	//fmt.Println(data)
	//println("OK")

	return toValue(h.runtime, data)
}

func toValue(vm *goja.Runtime, data interface{}) goja.Value {
	if val, ok := data.(bson.M); ok {
		obj := vm.NewObject()

		for key, value := range val {
			if objectId, ok := value.(primitive.ObjectID); ok {
				_ = obj.Set(key, objectId.Hex())
			} else {
				_ = obj.Set(key, vm.ToValue(value))
			}
		}

		return obj
	}
	return goja.Undefined()
}

func (m NativeMongoDBHandler) database(call goja.FunctionCall) goja.Value {
	name := converters.String(m.runtime, call.Argument(0))
	db := m.connection.Database(name)
	handler := NativeMongoDBDatabaseHandler{
		runtime:    m.runtime,
		connection: m.connection,
		database:   db,
	}

	object := m.runtime.NewObject()
	_ = object.Set("collection", handler.collection)
	return object
}

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
	_ = object.Set("find", handler.find)
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

	return object
}
