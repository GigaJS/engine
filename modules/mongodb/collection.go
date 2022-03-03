package mongodb

import (
	"git.nonamestudio.me/gjs/engine/core/converters"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NativeMongoDBCollectionHandler struct {
	runtime    *goja.Runtime
	connection *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

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

func (h NativeMongoDBCollectionHandler) deleteOne(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)
	entityObject := call.Argument(0)
	filter := h.extractObjectToBson(entityObject.ToObject(h.runtime))

	go func() {

		res, err := h.collection.DeleteOne(ctx, filter)

		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		resolve(toValue(h.runtime, res.DeletedCount))

	}()

	return h.runtime.ToValue(promise)
}

func (h NativeMongoDBCollectionHandler) deleteMany(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)
	entityObject := call.Argument(0)
	filter := h.extractObjectToBson(entityObject.ToObject(h.runtime))

	go func() {

		res, err := h.collection.DeleteMany(ctx, filter)

		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		resolve(toValue(h.runtime, res.DeletedCount))

	}()

	return h.runtime.ToValue(promise)
}

func (h NativeMongoDBCollectionHandler) updateOne(call goja.FunctionCall) goja.Value {

	promise, resolve, reject := loop.NewPromise(h.runtime)
	entityObject := call.Argument(0)
	updateObject := call.Argument(1)

	filter := h.extractObjectToBson(entityObject.ToObject(h.runtime))
	update := h.extractObjectToBson(updateObject.ToObject(h.runtime))

	go func() {

		res, err := h.collection.UpdateOne(ctx, filter, update)

		if err != nil {
			reject(h.runtime.ToValue(err.Error()))
			return
		}

		obj := h.runtime.NewObject()

		_ = obj.Set("matched", h.runtime.ToValue(res.MatchedCount))
		_ = obj.Set("modified", h.runtime.ToValue(res.ModifiedCount))
		_ = obj.Set("upsert", h.runtime.ToValue(res.UpsertedCount))
		_ = obj.Set("upsertId", toValue(h.runtime, res.UpsertedID))

		resolve(obj)

	}()

	return h.runtime.ToValue(promise)
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
