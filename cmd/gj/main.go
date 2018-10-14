package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gojson"
	"github.com/jimmyjames85/gotools/must"
	"github.com/pkg/profile"
)

func traverse(byts []byte) {
	j, _, err := gojson.ParseJSON(byts)
	must.BeNil(err)
	fmt.Printf("%d", j.Value.Type)

	arr, err := j.Value.Array()
	if err != nil {
		fmt.Printf("not an array: %s", err.Error())
		return
	}

	for _, elem := range arr {
		fmt.Printf("%d: %d\n", elem.Value.Type, len(elem.Value.String()))
	}

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

		b, size, err := gojson.ParseJSON(example)
		must.BeNil(err)
		if len(example) != size {
			panic(fmt.Sprintf("example size[%d] is different than parsed size[%d]", len(example), size))
		}
		if !silent {
			fmt.Printf("%s\n", b.Value.String())
		}

	}

}
