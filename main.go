package main

import (
	"github.com/cgentron/protoc-gen-cgentron-amzn/pkg/module"
	"github.com/cgentron/protoc-gen-cgentron/pkg/modules"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func handler() {
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV")).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}

func main() {
	runtime := modules.New()

	if err := runtime.Execute(handler); err != nil {
		panic(err)
	}
}
