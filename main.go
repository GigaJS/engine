package main

import (
	"git.nonamestudio.me/gjs/engine/core"
	"io/ioutil"
)

func main() {

	engine := core.CreateGJSEngine()
	data, err := ioutil.ReadFile("test.js")
	if err != nil {
		panic(err)
	}

	engine.ExecuteFromString(string(data))
	engine.Start()
}
