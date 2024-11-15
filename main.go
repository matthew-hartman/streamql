package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"code.arista.io/lib/streamql/pkg/build"
	"github.com/itchyny/gojq"
)

func main() {

	nameQuery, err := gojq.Parse(".after.Value.name")
	if err != nil {
		panic(err)
	}
	mutname := func(a map[string]any) string {
		iter := nameQuery.Run(a)
		v, ok := iter.Next()
		if !ok {
			return ""
		}
		if _, ok := v.(error); ok {
			return ""
		}
		return v.(string)
	}

	callbackA := func(a map[string]any) {
		fmt.Printf("mhartman: %#v\n", mutname(a))
	}

	callbackB := func(a map[string]any) {
		fmt.Printf("pllong: %#v\n", mutname(a))
	}

	idQuery, err := gojq.Parse(".after.Value.id")
	if err != nil {
		panic(err)
	}
	callbackC := func(a map[string]any) {
		iter := idQuery.Run(a)
		v, ok := iter.Next()
		if !ok {
			return
		}
		if _, ok := v.(error); ok {
			return
		}
		fmt.Printf("ID: %#v\n", v)
	}

	changeNumQuery, err := gojq.Parse(".after.Value.mergedChangeNum.int")
	if err != nil {
		panic(err)
	}
	callbackD := func(a map[string]any) {
		iter := idQuery.Run(a)
		v, ok := iter.Next()
		if !ok {
			return
		}
		if _, ok := v.(error); ok {
			return
		}
		iter = changeNumQuery.Run(a)
		v1, ok := iter.Next()
		if !ok {
			return
		}
		if _, ok := v.(error); ok {
			return
		}
		fmt.Printf("Merged: %v @ %d\n", v, int(v1.(float64)))
	}

	b, err := build.NewBuilder[map[string]any]().
		AddVarExp(`.after.Value.owner.string`, "mhartman").
		AddBoolExp(`.after.Value.parent.string=="eos-trunk"`).
		Build(build.WithCallback(callbackA))
	if err != nil {
		panic(err)
	}
	b, err = b.
		AddVarExp(`.after.Value.owner.string`, "pllong").
		AddBoolExp(`.after.Value.parent.string=="eos-trunk"`).
		Build(build.WithCallback(callbackB))
	if err != nil {
		panic(err)
	}

	for i := 800000; i < 887083; i += 15 {
		b, err = b.
			AddVarExp(`.after.Value.id`, i).
			AddBoolExp(`.after.Value.parent.string=="eos-trunk"`).
			Build(build.WithCallback(callbackC))
		if err != nil {
			panic(err)
		}
	}

	b, err = b.AddBoolExp(`.before.Value.mergedChangeNum!=.after.Value.mergedChangeNum`).
		AddBoolExp(`.after.Value.mergedChangeNum.int % 13 == 0`).
		Build(build.WithCallback(callbackD))
	if err != nil {
		panic(err)
	}

	var input map[string]any
	decoder := json.NewDecoder(os.Stdin)
	total := time.Duration(0)
	count := 0
	for {
		err := decoder.Decode(&input)
		if err != nil {
			break
		}
		st := time.Now()
		_, err = b.Query(input)
		total += time.Since(st)
		count++
	}
	fmt.Printf("Processed %d Messages\n", count)
	fmt.Printf("Avg query %v\n", time.Duration(int64(total)/int64(count)))
}
