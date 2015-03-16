package main

import (
	"github.com/pepabo/facade"
)

func main() {
	f := &facade.Facade{}
	f.Env = map[string]string{
		"FACADE_FOO": "123",
		"FACADE_BAR": "Bar Value",
	}

	f.Run()
}
