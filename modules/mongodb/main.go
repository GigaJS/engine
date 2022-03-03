package mongodb

import (
	"context"
	"reflect"

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

func (h NativeMongoDBCollectionHandler) extractObjectToBson(filterArgument *goja.Object) interface{} {

	var filter interface{}

	if filterArgument != nil && !goja.IsUndefined(filterArgument) && !goja.IsNull(filterArgument) {
		filter = fromValue(h.runtime, filterArgument)
	} else {
		filter = bson.D{}
	}

	return filter

}

func (h NativeMongoDBCollectionHandler) findOne(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)

	filterArgument := call.Argument(0)

	optionsArgument := call.Argument(1)
	findOptions := options.FindOne()

	filter := h.extractObjectToBson(filterArgument.ToObject(h.runtime))

	if optionsArgument != nil && !goja.IsUndefined(optionsArgument) && !goja.IsNull(optionsArgument) {
		optionsObject := optionsArgument.ToObject(h.runtime)
		if converters.IsPresentInObject(h.runtime, optionsObject, "skip") {
			skip := converters.NumberInt64(h.runtime, optionsObject.Get("skip"))
			findOptions.Skip = &skip
		}
	}

	go func() {

		singleRes := h.collection.FindOne(ctx, filter, findOptions)

		if singleRes.Err() != nil {
			reject(h.runtime.ToValue(singleRes.Err().Error()))
			return
		}

		data := bson.M{}

		err := singleRes.Decode(&data)
		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		resolve(toValue(h.runtime, data))

	}()

	return h.runtime.ToValue(promise)
}

func (h NativeMongoDBCollectionHandler) find(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)

	filterArgument := call.Argument(0)

	optionsArgument := call.Argument(1)
	findOptions := options.Find()

	filter := h.extractObjectToBson(filterArgument.ToObject(h.runtime))

	if optionsArgument != nil && !goja.IsUndefined(optionsArgument) && !goja.IsNull(optionsArgument) {
		optionsObject := optionsArgument.ToObject(h.runtime)

		if converters.IsPresentInObject(h.runtime, optionsObject, "skip") {
			skip := converters.NumberInt64(h.runtime, optionsObject.Get("skip"))
			findOptions.Skip = &skip
		}
		if converters.IsPresentInObject(h.runtime, optionsObject, "limit") {
			limit := converters.NumberInt64(h.runtime, optionsObject.Get("limit"))
			findOptions.Limit = &limit
		}
	}

	go func() {

		resultSet, err := h.collection.Find(ctx, filter, findOptions)

		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		var results []bson.M

		err = resultSet.All(ctx, &results)
		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		var gjsResults []goja.Value

		for _, res := range results {
			gjsResults = append(gjsResults, toValue(h.runtime, res))
		}

		resolve(h.runtime.ToValue(gjsResults))

	}()

	return h.runtime.ToValue(promise)
}

func (h NativeMongoDBCollectionHandler) insert(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)
	entityObject := call.Argument(0)
	entity := h.extractObjectToBson(entityObject.ToObject(h.runtime))

	go func() {

		insertResult, err := h.collection.InsertOne(ctx, entity)

		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		id := insertResult.InsertedID

		resolve(toValue(h.runtime, id))

	}()

	return h.runtime.ToValue(promise)
}

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

	if objData, ok := data.(*goja.Object); ok {
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
	_ = object.Set("findOne", handler.findOne)
	_ = object.Set("find", handler.find)
	_ = object.Set("insert", handler.insert)
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
