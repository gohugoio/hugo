package plugins

import (
	"errors"
	"fmt"
)

// Starlight runs the given starlight script with the given values as globals.
// The globals are expected to be associated pairs of name (string) and value,
// which will be passed to the script. It's an error if you don't pass the
// values in pairs. The value of the template function is set by calling the
// "output" function in the script with the value you wish to output.
//
// e.g. {{plugins.starlight "foo.star" "page" .Page "site" .Site }}
//
// plugins/hello.star:
//
// output("hello" + site.Title)
func (ns *Namespace) Starlight(filename interface{}, vals ...interface{}) (interface{}, error) {
	s, ok := filename.(string)
	if !ok {
		return nil, fmt.Errorf("expected first argument to be filename (string), but was %v (%T)", filename, filename)
	}

	if len(vals)%2 != 0 {
		return nil, fmt.Errorf("expected values to be pairs of <name> <value>")
	}

	var ret interface{}
	output := func(v interface{}) {
		ret = v
	}
	globals := map[string]interface{}{
		"output": output,
	}

	for i := 0; i < len(vals); i += 2 {
		name, ok := vals[i].(string)
		if !ok {
			return nil, fmt.Errorf("expected argument %d to be string label, but was %v (%T)", i, vals[i], vals[i])
		}
		if name == "output" {
			return nil, errors.New("argument label `output` is a reserved word and cannot be used")
		}
		globals[name] = vals[i+1]
	}

	_, err := ns.cache.Run(s, globals)
	return ret, err
}
