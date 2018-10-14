package main

import (
	"fmt"

	"github.com/jimmyjames85/gojson"
	"github.com/jimmyjames85/gotools/must"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start().Stop()

	for i := 0; i < 10000; i++ {
		example := must.ReadFile("example.json")
		b, size, err := gojson.ParseJSON(example)
		must.BeNil(err)

		if len(example) != size {
			panic(fmt.Sprintf("example size[%d] is different than parsed size[%d]", len(example), size))
		}

		fmt.Printf("%s\n", string(b))
	}

}
