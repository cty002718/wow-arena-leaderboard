package ctx

import (
	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()
}

func Get[T any]() T {
	var v T
	err := Container.Invoke(func(dep T) {
		v = dep
	})

	if err != nil {
		panic(err)
	}

	return v
}
