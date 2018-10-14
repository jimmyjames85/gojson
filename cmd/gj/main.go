package main

import (
	"fmt"
	"os"

	"github.com/jimmyjames85/gojson"
	"github.com/jimmyjames85/gotools/must"
	"github.com/pkg/profile"
)

// todo move this into gojson package...
func Walk(prefix string, v gojson.Value) {

	switch v.Type() {
	case gojson.ObjectType:
		val := v.Object()
		for k, v := range val {
			prfx := fmt.Sprintf("%s:%s", prefix, k)
			fmt.Printf("%s\n", prfx)
			Walk(prfx, v)
		}
	case gojson.ArrayType:
		val := v.Array()
		for _, v := range val {
			// todo decend into object if it is one...
			prfx := fmt.Sprintf("%s:%s", prefix, v.B())
			fmt.Printf("%s\n", prfx)
			Walk(prfx, v)
		}
	case gojson.StringType:
	case gojson.NumberType:
	case gojson.BooleanType:
	case gojson.NullType:
	}
}

func traverse(byts []byte) {
	j, err := gojson.ParseJSON(byts)
	must.BeNil(err)
	fmt.Printf("%s: %d\n", j.Type(), len(j.B()))

	arr := j.Array()
	if arr == nil {
		fmt.Printf("not an array maybe? I didn't check tht ey  type ... shame: %s\n", j.Type())
		return
	}

	for _, val := range arr {
		fmt.Printf("%s: %d\n", val.Type(), len(val.B()))
		Walk("", val)
		fmt.Printf("\n\n")
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

		val, err := gojson.ParseJSON(example)
		must.BeNil(err)
		if !silent {
			fmt.Printf("%s\n", val.B())
		}

	}

}
