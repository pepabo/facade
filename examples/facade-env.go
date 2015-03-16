package main

import (
	"github.com/kentaro/facade"
)

func main() {
	f := &facade.Facade{}
	f.Environment = map[string]string{
		"FACADE_FOO": "123",
		"FACADE_BAR": "Bar Value",
	}

	f.Run()
}
