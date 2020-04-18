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

package publisher

import (
	"testing"
  "strings"
	qt "github.com/frankban/quicktest"
)

func TestClassServerAliases(t *testing.T) {
	c := qt.New((t))

  // Geneator of aliases elements
  //
  // Usage :
  //
  //   generateAliases([]string{"/src,/dest,301","/src2,/dest2,302"})
  //
  generateAliases := func(elms []string) []AliasElement {
    var aliases []AliasElement
    for _,e := range elms {
      var elm = strings.Split(e, ",")
      var alias = AliasElement{
        source: elm[0],
        destination: elm[1],
        redirectionType: elm[2],
      }
      aliases = append(aliases, alias)
    }

    return aliases
  }

  t.Run("no-alias", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    c.Assert(aliases.Format("nginx"), qt.DeepEquals, []byte(""))
  })

  t.Run("bad-server", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    c.Assert(aliases.Format("dummy"), qt.DeepEquals, []byte(""))
  })

  t.Run("simple-301-nginx", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    aliases.AliasElements = generateAliases([]string{"/src,/dest,301"})
    c.Assert(string(aliases.Format("nginx")), qt.DeepEquals, "rewrite ^/src$ /dest permanent;\n")
  })

  t.Run("simple-302-nginx", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    aliases.AliasElements = generateAliases([]string{"/src,/dest,302"})
    c.Assert(string(aliases.Format("nginx")), qt.DeepEquals, "rewrite ^/src$ /dest redirect;\n")
  })

  t.Run("many-301-nginx", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    aliases.AliasElements = generateAliases([]string{"/src,/dest,301","/src2,/dest2,301"})
    c.Assert(string(aliases.Format("nginx")), qt.DeepEquals, "rewrite ^/src$ /dest permanent;\nrewrite ^/src2$ /dest2 permanent;\n")
  })

  t.Run("many-mixed-status-nginx", func(t *testing.T) {
    var aliases ServerAliases = ServerAliases{}
    aliases.AliasElements = generateAliases([]string{"/src,/dest,302","/src2,/dest2,301"})
    c.Assert(string(aliases.Format("nginx")), qt.DeepEquals, "rewrite ^/src$ /dest redirect;\nrewrite ^/src2$ /dest2 permanent;\n")
  })
}
