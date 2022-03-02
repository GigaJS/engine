package http

import (
	"encoding/json"
	"git.nonamestudio.me/gjs/engine/core/loop"
	"github.com/dop251/goja"
	"github.com/erikdubbelboer/fasthttp"
	routing "github.com/jackwhelpton/fasthttp-routing"
	"reflect"
)

type route struct {
	Method   string
	Handlers []routing.Handler
}

type Server struct {
	runtime     *goja.Runtime
	releaseLock func()

	uses     []routing.Handler
	groups   map[string]map[string]route
	handlers []routing.Handler
}

func (h *Server) Use(call goja.FunctionCall) goja.Value {
	callback := call.Argument(0)

	if c, ok := goja.AssertFunction(callback); ok {
		h.uses = append(h.uses, func(context *routing.Context) error {
			_, _ = c(goja.Undefined())
			return nil
		})
	} else {
		panic(h.runtime.ToValue("Callback must be passed"))
	}

	return goja.Undefined()
}

func (h *Server) CreateGroup(call goja.FunctionCall) goja.Value {
	pathValue := call.Argument(0)
	if pathValue.ExportType().Kind() != reflect.String {
		panic(h.runtime.ToValue("Group path must be a string"))
		return goja.Undefined()
	}

	groupPrefix := pathValue.Export().(string)

	h.groups[groupPrefix] = make(map[string]route)

	o := h.runtime.NewObject()
	_ = o.Set("get", func(call goja.FunctionCall) goja.Value {
		routeValue := call.Argument(0)
		if routeValue.ExportType().Kind() != reflect.String {
			return nil
		}

		routePath := routeValue.Export().(string)

		var handlers []routing.Handler

		t := call.Arguments[1:]
		for i, arg := range t {
			if cb, ok := goja.AssertFunction(arg); ok {
				if i == len(t)-1 { // Primary
					handlers = append(handlers, func(context *routing.Context) error {
						ctx := h.runtime.NewObject()
						_ = ctx.Set("reply", func(replyCall goja.FunctionCall) goja.Value {
							replyValue := replyCall.Argument(0).Export()
							var d string

							if v, ok := replyValue.(string); ok {
								d = v
							} else if v, ok := replyValue.(map[string]interface{}); ok {
								b, _ := json.Marshal(v)
								d = string(b)
							}

							_, _ = context.WriteString(d)
							return goja.Undefined()
						})

						_, _ = cb(goja.Undefined(), ctx)
						return context.Next()
					})
				} else { // Middleware
					handlers = append(handlers, func(context *routing.Context) error {
						_, _ = cb(goja.Undefined())
						return context.Next()
					})
				}
			} else {
				panic(h.runtime.ToValue("invalid middleware"))
				return nil
			}
		}

		h.groups[groupPrefix][routePath] = route{
			Method:   "get",
			Handlers: handlers,
		}
		return nil
	})
	return o
}

func (h Server) Listen(call goja.FunctionCall) goja.Value {
	addr := call.Argument(0)
	if addr.ExportType().Kind() != reflect.String {
		panic(h.runtime.ToValue("invalid address"))
		return goja.Undefined()
	}

	addStr := addr.Export().(string)

	go func() {
		router := routing.New()

		for s, m := range h.groups {
			g := router.Group(s)

			for p, r := range m {
				switch r.Method {
				case "get":
					g.Get(p, r.Handlers...)
				}
			}
		}

		err := fasthttp.ListenAndServe(addStr, router.HandleRequest)
		if err != nil {
			return
		}

		h.releaseLock = loop.NewLock(h.runtime)
	}()

	return goja.Undefined()
}

func (h Server) Close(call goja.FunctionCall) goja.Value {
	if h.releaseLock == nil {
		panic(h.runtime.ToValue("Server not launched"))
	}

	h.releaseLock()

	// TODO: Stop server

	return goja.Undefined()
}
