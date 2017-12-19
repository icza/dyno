package dyno

import "fmt"

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

	v, err := Get(m, "a")
	printResults(v, err)

	v, err = Get(m, "b", 3, 1)
	printResults(v, err)

	v, err = Get(m, "x")
	printResults(v, err)

	sl, _ := Get(m, "b", 3) // This is: []interface{}{1, "two", 3.3}
	v, err = Get(sl, 4)
	printResults(v, err)

	// Output:
	// 1
	// two
	// ERROR: missing key: x (path element idx: 0)
	// ERROR: index out of range: 4 (path element idx: 0)
}
