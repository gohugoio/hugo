// Copyright 2026 The Hugo Authors. All rights reserved.
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

package hhtml

import (
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestScanner(t *testing.T) {
	c := qt.New(t)

	scan := func(s string) string {
		var parts []string
		sc := NewScanner(s, 0)
		for {
			t, ok := sc.Next(len(s))
			if !ok {
				break
			}
			switch t.Type {
			case StartTag:
				if t.SelfClosing {
					parts = append(parts, fmt.Sprintf("<%s/>@%d", t.Name, t.Start))
				} else {
					parts = append(parts, fmt.Sprintf("<%s>@%d", t.Name, t.Start))
				}
				if IsRawTextElement(t.Name) {
					sc.SkipRawText(len(s), t.Name)
				}
			case EndTag:
				parts = append(parts, fmt.Sprintf("</%s>@%d", t.Name, t.Start))
			case Markup:
				parts = append(parts, fmt.Sprintf("markup@%d", t.Start))
			}
		}
		return strings.Join(parts, " ")
	}

	c.Assert(scan("no tags"), qt.Equals, "")
	c.Assert(scan("<P>a</P>"), qt.Equals, "<p>@0 </p>@4")
	c.Assert(scan(`<a title="1>2" data-x='">'>x</a>`), qt.Equals, "<a>@0 </a>@28")
	c.Assert(scan("<br/><img src=x /><p>"), qt.Equals, "<br/>@0 <img/>@5 <p>@18")
	c.Assert(scan("I <3 Hugo <b>"), qt.Equals, "<b>@10")
	c.Assert(scan("<!-- <b> --><!DOCTYPE html><![CDATA[<b>]]><?xml?></ x>"), qt.Equals, "markup@0 markup@12 markup@27 markup@42 markup@49")
	c.Assert(scan("<script>if (a</b) {}</SCRIPT ><p>"), qt.Equals, "<script>@0 </script>@20 <p>@30")
	c.Assert(scan("<b"), qt.Equals, "<b>@0")
	c.Assert(scan("a <"), qt.Equals, "")
}
