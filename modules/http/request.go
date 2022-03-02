package http

import (
	"git.nonamestudio.me/gjs/engine/core/converters"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{}

func (f *Module) request(call goja.FunctionCall) goja.Value {
	options := call.Argument(0).ToObject(f.runtime)

	promise, resolve, reject := loop.NewPromise(f.runtime)

	method := converters.StringOrDefault(f.runtime, options.Get("method"), "GET")
	uri := converters.String(f.runtime, options.Get("url"))

	requestHeaders := make(map[string]string)

	headers := options.Get("headers")
	if headers != nil && !goja.IsUndefined(headers) {
		object := headers.ToObject(f.runtime)
		for _, key := range object.Keys() {
			val := object.Get(key)
			requestHeaders[key] = converters.String(f.runtime, val)
		}
	}

	go func() {

		request, err := http.NewRequest(method, uri, nil)
		if err != nil {
			reject(err)
			return
		}

		for key, val := range requestHeaders {
			request.Header.Set(key, val)
		}

		resp, err := client.Do(request)

		if err != nil {
			reject(err)
			return
		}

		resolve(processResponse(f.runtime, resp))
	}()

	return f.runtime.ToValue(promise)
}

func processResponse(vm *goja.Runtime, response *http.Response) goja.Value {
	obj := vm.NewObject()

	_ = obj.Set("status", response.StatusCode)
	_ = obj.Set("uncompressed", response.Uncompressed)
	_ = obj.Set("protocol", response.Proto)

	bodyBytes, _ := ioutil.ReadAll(response.Body)

	cb, ok := goja.AssertFunction(vm.Get("Buffer"))
	if !ok {
		panic(vm.ToValue("Buffer not defined"))
	}
	value, err := cb(goja.Undefined(), vm.ToValue(vm.NewArrayBuffer(bodyBytes)))
	if err != nil {
		panic(vm.ToValue(err))
	}
	_ = obj.Set("body", value)

	headers := vm.NewObject()

	//var t map[string][]string
	for name, values := range response.Header {
		if len(values) > 1 {
			headerValues := vm.NewArray(values)
			_ = headers.Set(name, headerValues)
		} else {
			_ = headers.Set(name, values[0])
		}
	}
	_ = obj.Set("headers", headers)

	return obj
}
