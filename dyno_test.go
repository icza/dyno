package dyno

import (
	"reflect"
	"testing"
)

var (
	s = []interface{}{
		1, "a", 3.3, []interface{}{"inner", "inner2"},
	}
	mi = map[interface{}]interface{}{
		"x": 1,
		"y": 2,
		"z": map[interface{}]interface{}{
			3: "three",
		},
	}
	ms = map[string]interface{}{
		"a": 1,
		"p": map[string]interface{}{
			"x": 1,
			"y": 2,
		},
		"pi": mi,
		"ns": []interface{}{1.1, 2.2, 3.3},
		"b":  2,
		"sl": s,
	}
)

func TestGeneral(t *testing.T) {
	cases := []struct {
		title string        // Human readable title of the test case
		v     interface{}   // Input dynamic object
		path  []interface{} // path whose value to get
		value interface{}   // Expected value
		isErr bool          // Tells if error is expected
	}{
		// Testing successes:

		{
			title: "nil path on map (MS)",
			v:     ms,
			path:  nil,
			value: ms,
			isErr: false,
		},
		{
			title: "nil path on map (MI)",
			v:     mi,
			path:  nil,
			value: mi,
			isErr: false,
		},
		{
			title: "nil path on slice",
			v:     s,
			path:  nil,
			value: s,
			isErr: false,
		},
		{
			title: "empty path on map (MS)",
			v:     ms,
			path:  []interface{}{},
			value: ms,
			isErr: false,
		},
		{
			title: "empty path on map (MI)",
			v:     mi,
			path:  []interface{}{},
			value: mi,
			isErr: false,
		},
		{
			title: "empty path on slice",
			v:     s,
			path:  []interface{}{},
			value: s,
			isErr: false,
		},
		{
			title: "simple map element",
			v:     ms,
			path:  []interface{}{"a"},
			value: 1,
			isErr: false,
		},
		{
			title: "simple map element #2",
			v:     ms,
			path:  []interface{}{"sl"},
			value: s,
			isErr: false,
		},
		{
			title: "nested map element",
			v:     ms,
			path:  []interface{}{"p", "x"},
			value: 1,
			isErr: false,
		},
		{
			title: "nested map (MI) element",
			v:     ms,
			path:  []interface{}{"pi", "x"},
			value: 1,
			isErr: false,
		},
		{
			title: "nested map (MI) element #2",
			v:     ms,
			path:  []interface{}{"pi", "z", 3},
			value: "three",
			isErr: false,
		},
		{
			title: "nested slice element",
			v:     s,
			path:  []interface{}{3, 1},
			value: "inner2",
			isErr: false,
		},
		{
			title: "map element and slice element",
			v:     ms,
			path:  []interface{}{"ns", 1},
			value: 2.2,
			isErr: false,
		},
		{
			title: "map element and slice element #2",
			v:     ms,
			path:  []interface{}{"sl", 1},
			value: "a",
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
			title: "missing key (MS) error",
			v:     ms,
			path:  []interface{}{"x"},
			value: nil,
			isErr: true,
		},
		{
			title: "missing key (MI) error",
			v:     mi,
			path:  []interface{}{"a"},
			value: nil,
			isErr: true,
		},
		{
			title: "element is not string error",
			v:     ms,
			path:  []interface{}{1},
			value: nil,
			isErr: true,
		},
		{
			title: "element is not int error",
			v:     ms,
			path:  []interface{}{"ns", "x"},
			value: nil,
			isErr: true,
		},
		{
			title: "index out of range error (negative)",
			v:     ms,
			path:  []interface{}{"ns", -1},
			value: nil,
			isErr: true,
		},
		{
			title: "index out of range error (too big)",
			v:     ms,
			path:  []interface{}{"ns", 11},
			value: nil,
			isErr: true,
		},
	}

	for _, c := range cases {
		value, err := Get(c.v, c.path...)
		if !reflect.DeepEqual(value, c.value) {
			t.Errorf("[title: %s] Expected value: %v, got: %v", c.title, c.value, value)
		}
		if c.isErr != (err != nil) {
			t.Errorf("[title: %s] Expected error: %v, got: %v, err value: %v", c.title, c.isErr, err != nil, err)
		}
	}
}

// TestTypeGetEmptyPath tests type-specific Get methods with empty path.
func TestTypeGetEmptyPath(t *testing.T) {
	cases := []struct {
		name     string                                    // Name of the type
		receiver interface{}                               // Receiver
		get      func(...interface{}) (interface{}, error) // Get method value
	}{
		{"MS", MS(ms), MS(ms).Get},
		{"MI", MI(mi), MI(mi).Get},
		{"S", S(s), S(s).Get},
	}

	for _, c := range cases {
		v, err := c.get()
		if err != nil {
			t.Errorf("%s.Value() with empty path returned error: %v", c.name, err)
		}
		if !reflect.DeepEqual(v, c.receiver) {
			t.Errorf("%s.Value() with empty path misbehaves, expected: %v, got: %v", c.name, ms, v)
		}
	}
}
