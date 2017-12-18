package dyno

import "testing"

func TestGeneral(t *testing.T) {
	obj := map[string]interface{}{
		"a": 1,
		"p": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
		"ns": []interface{}{1.1, 2.2, 3.3},
		"b":  2,
	}
	sl := []interface{}{
		1, "a", 3.3, []interface{}{"inner", "inner2"},
	}

	cases := []struct {
		title string        // Human readable title of the test case
		v     interface{}   // Input dynamic object
		path  []interface{} // path whose value to get
		value interface{}   // Expected value
		isErr bool          // Tells if error is expected
	}{
		// Testing successes:

		{
			title: "nil path on map",
			v:     obj,
			path:  nil,
			value: nil,
			isErr: false,
		},
		{
			title: "nil path on slice",
			v:     sl,
			path:  nil,
			value: nil,
			isErr: false,
		},
		{
			title: "empty path on map",
			v:     obj,
			path:  []interface{}{},
			value: nil,
			isErr: false,
		},
		{
			title: "empty path on slice",
			v:     sl,
			path:  []interface{}{},
			value: nil,
			isErr: false,
		},
		{
			title: "simple map element",
			v:     obj,
			path:  []interface{}{"a"},
			value: 1,
			isErr: false,
		},
		{
			title: "nested map element",
			v:     obj,
			path:  []interface{}{"p", "x"},
			value: 1,
			isErr: false,
		},
		{
			title: "nested slice element",
			v:     sl,
			path:  []interface{}{3, 1},
			value: "inner2",
			isErr: false,
		},
		{
			title: "map element and slice element",
			v:     obj,
			path:  []interface{}{"ns", 1},
			value: 2.2,
			isErr: false,
		},

		// Testing errors:

		{
			title: "invalid node type error",
			v:     1,
			path:  []interface{}{"x"},
			value: nil,
			isErr: true,
		},
		{
			title: "missing key error",
			v:     obj,
			path:  []interface{}{"x"},
			value: nil,
			isErr: true,
		},
		{
			title: "element is not string error",
			v:     obj,
			path:  []interface{}{1},
			value: nil,
			isErr: true,
		},
		{
			title: "element is not int error",
			v:     obj,
			path:  []interface{}{"ns", "x"},
			value: nil,
			isErr: true,
		},
		{
			title: "index out of range error (negative)",
			v:     obj,
			path:  []interface{}{"ns", -1},
			value: nil,
			isErr: true,
		},
		{
			title: "index out of range error (too big)",
			v:     obj,
			path:  []interface{}{"ns", 11},
			value: nil,
			isErr: true,
		},
	}

	for _, c := range cases {
		value, err := Value(c.v, c.path...)
		if value != c.value {
			t.Errorf("[title: %s] Expected value: %v, got: %v", c.title, c.value, value)
		}
		if c.isErr != (err != nil) {
			t.Errorf("[title: %s] Expected error: %v, got: %v, err value: %v", c.title, c.isErr, err != nil, err)
		}
	}
}
