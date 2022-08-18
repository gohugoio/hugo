// Copyright 2020 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package npm

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

const templ = `{
        "name": "foo",
        "version": "0.1.1",
        "scripts": {},
         "dependencies": {
                "react-dom": "1.1.1",
                "tailwindcss": "1.2.0",
                "@babel/cli": "7.8.4",
                "@babel/core": "7.9.0",
                "@babel/preset-env": "7.9.5"
        },
        "devDependencies": {
                "postcss-cli": "7.1.0",
                "tailwindcss": "1.2.0",
                "@babel/cli": "7.8.4",
                "@babel/core": "7.9.0",
                "@babel/preset-env": "7.9.5"
        }
}`

func TestPackageBuilder(t *testing.T) {
	c := qt.New(t)

	b := newPackageBuilder("", strings.NewReader(templ))
	c.Assert(b.Err(), qt.IsNil)

	b.Add("mymod", strings.NewReader(`{
"dependencies": {
	 "react-dom": "9.1.1",
	 "add1": "1.1.1"
},
"devDependencies": {
	 "tailwindcss": "error",
	 "add2": "2.1.1"
}	
}`))

	b.Add("mymod", strings.NewReader(`{
"dependencies": {
	 "react-dom": "error",
	 "add1": "error",
	 "add3": "3.1.1"
},
"devDependencies": {
	 "tailwindcss": "error",
	 "add2": "error",
	 "add4": "4.1.1"
	 
}	
}`))

	c.Assert(b.Err(), qt.IsNil)

	c.Assert(b.dependencies, qt.DeepEquals, map[string]any{
		"@babel/cli":        "7.8.4",
		"add1":              "1.1.1",
		"add3":              "3.1.1",
		"@babel/core":       "7.9.0",
		"@babel/preset-env": "7.9.5",
		"react-dom":         "1.1.1",
		"tailwindcss":       "1.2.0",
	})

	c.Assert(b.devDependencies, qt.DeepEquals, map[string]any{
		"tailwindcss":       "1.2.0",
		"@babel/cli":        "7.8.4",
		"@babel/core":       "7.9.0",
		"add2":              "2.1.1",
		"add4":              "4.1.1",
		"@babel/preset-env": "7.9.5",
		"postcss-cli":       "7.1.0",
	})
}
