package dyno_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/icza/dyno"
)

func ExampleGet() {
	m := map[string]interface{}{
		"a": 1,
		"b": map[interface{}]interface{}{
			3: []interface{}{1, "two", 3.3},
		},
	}

	printResults := func(v interface{}, err error) {
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			fmt.Println(v)
		}
	}

	v, err := dyno.Get(m, "a")
	printResults(v, err)

	v, err = dyno.Get(m, "b", 3, 1)
	printResults(v, err)

	v, err = dyno.Get(m, "x")
	printResults(v, err)

	sl, _ := dyno.Get(m, "b", 3) // This is: []interface{}{1, "two", 3.3}
	v, err = dyno.Get(sl, 4)
	printResults(v, err)

	// Output:
	// 1
	// two
	// ERROR: missing key: x (path element idx: 0)
	// ERROR: index out of range: 4 (path element idx: 0)
}

func ExampleSet() {
	m := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"3": []interface{}{1, "two", 3.3},
		},
	}

	printResults := func(err error) {
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			// Use JSON output so map entry order is consistent:
			json.NewEncoder(os.Stdout).Encode(m)
		}
	}

	err := dyno.Set(m, 2, "a")
	printResults(err)

	err = dyno.Set(m, "owt", "b", "3", 1)
	printResults(err)

	err = dyno.Set(m, 1, "x")
	printResults(err)

	sl, _ := dyno.Get(m, "b", "3") // This is: []interface{}{1, "owt", 3.3}
	err = dyno.Set(sl, 1, 4)
	printResults(err)

	// Output:
	// {"a":2,"b":{"3":[1,"two",3.3]}}
	// {"a":2,"b":{"3":[1,"owt",3.3]}}
	// {"a":2,"b":{"3":[1,"owt",3.3]},"x":1}
	// ERROR: index out of range: 4 (path element idx: 0)
}
