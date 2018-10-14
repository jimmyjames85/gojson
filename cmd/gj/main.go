package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gojson"
	"github.com/jimmyjames85/gotools/must"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start().Stop()

	example := must.ReadFile("example.json")
	silent := len(os.Args) > 1 &&
		os.Args[1] == "-s"

	iter := 20000
	for i := 0; i < iter; i++ {

		b, size, err := gojson.ParseJSON(example)
		must.BeNil(err)
		if len(example) != size {
			panic(fmt.Sprintf("example size[%d] is different than parsed size[%d]", len(example), size))
		}
		if !silent {
			fmt.Printf("%s\n", string(b))
		}

	}

}
