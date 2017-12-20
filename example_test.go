package dyno_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/icza/dyno"
)

func Example() {
	person := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Bob",
			"last":  "Archer",
		},
		"age": 22,
		"fruits": []interface{}{
			"apple", "banana",
		},
	}

	printPerson := func(err error) {
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			// Use JSON output so map entry order is consistent:
			json.NewEncoder(os.Stdout).Encode(person)
		}
	}

	printPerson(nil)

	v, err := dyno.Get(person, "name", "first")
	fmt.Println(v, err)

	printPerson(dyno.Set(person, "Alice", "name", "first"))
	printPerson(dyno.Set(person, "Alice Archer", "name"))
	printPerson(dyno.Set(person, "lemon", "fruits", 1))
	printPerson(dyno.Append(person, "melon", "fruits"))

	// Output:
	// {"age":22,"fruits":["apple","banana"],"name":{"first":"Bob","last":"Archer"}}
	// Bob <nil>
	// {"age":22,"fruits":["apple","banana"],"name":{"first":"Alice","last":"Archer"}}
	// {"age":22,"fruits":["apple","banana"],"name":"Alice Archer"}
	// {"age":22,"fruits":["apple","lemon"],"name":"Alice Archer"}
	// {"age":22,"fruits":["apple","lemon","melon"],"name":"Alice Archer"}
}

func ExampleGet() {
	m := map[string]interface{}{
		"a": 1,
		"b": map[interface{}]interface{}{
			3: []interface{}{1, "two", 3.3},
		},
	}

	printValue := func(v interface{}, err error) {
		fmt.Printf("Value: %-5v, Error: %v\n", v, err)
	}

	printValue(dyno.Get(m, "a"))
	printValue(dyno.Get(m, "b", 3, 1))
	printValue(dyno.Get(m, "x"))

	sl, _ := dyno.Get(m, "b", 3) // This is: []interface{}{1, "two", 3.3}
	printValue(dyno.Get(sl, 4))

	// Output:
	// Value: 1    , Error: <nil>
	// Value: two  , Error: <nil>
	// Value: <nil>, Error: missing key: x (path element idx: 0)
	// Value: <nil>, Error: index out of range: 4 (path element idx: 0)
}

func ExampleSet() {
	m := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"3": []interface{}{1, "two", 3.3},
		},
	}

	printMap := func(err error) {
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			// Use JSON output so map entry order is consistent:
			json.NewEncoder(os.Stdout).Encode(m)
		}
	}

	printMap(dyno.Set(m, 2, "a"))
	printMap(dyno.Set(m, "owt", "b", "3", 1))
	printMap(dyno.Set(m, 1, "x"))

	sl, _ := dyno.Get(m, "b", "3") // This is: []interface{}{1, "owt", 3.3}
	printMap(dyno.Set(sl, 1, 4))

	// Output:
	// {"a":2,"b":{"3":[1,"two",3.3]}}
	// {"a":2,"b":{"3":[1,"owt",3.3]}}
	// {"a":2,"b":{"3":[1,"owt",3.3]},"x":1}
	// ERROR: index out of range: 4 (path element idx: 0)
}

func ExampleAppend() {
	m := map[string]interface{}{
		"a": []interface{}{
			"3", 2, []interface{}{1, "two", 3.3},
		},
	}

	printMap := func(err error) {
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			fmt.Println(m)
		}
	}

	printMap(dyno.Append(m, 4, "a"))
	printMap(dyno.Append(m, 9, "a", 2))
	printMap(dyno.Append(m, 1, "x"))

	// Output:
	// map[a:[3 2 [1 two 3.3] 4]]
	// map[a:[3 2 [1 two 3.3 9] 4]]
	// ERROR: missing key: x (path element idx: 0)
}
