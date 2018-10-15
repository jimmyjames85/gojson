package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gojson"
	"github.com/jimmyjames85/gotools/must"
	"github.com/pkg/profile"
)

func traverse(byts []byte) {
	j, err := gojson.ParseJSON(byts)
	must.BeNil(err)

	j.Walk(func(path string, val gojson.Value) {
		fmt.Printf("%s: %s\n", path, val.Type())
	})
}

func main() {

	example := must.ReadFile("example.json")

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		traverse(example)
		return
	}

	defer profile.Start().Stop()

	silent := len(os.Args) > 1 &&
		os.Args[1] == "-s"

	iter := 1
	if silent {
		iter = 20000
	}
	for i := 0; i < iter; i++ {

		val, err := gojson.ParseJSON(example)
		must.BeNil(err)
		if !silent {
			fmt.Printf("%s\n", val.B())
		}

	}

}
