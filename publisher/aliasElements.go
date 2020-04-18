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
  "fmt"
)

type AliasElement struct {
  source string
  destination string
  redirectionType string
}

type ServerAliases struct {
  // array of aliases
  AliasElements []AliasElement
}

// Convert the slice of AliasElement stored as current ServerAliases
// as valide http server configuration.
//
// First arg is the target http server type
func (p *ServerAliases) Format(format string) []byte {
  if (format == "nginx") {
    var result string
    for _, alias := range p.AliasElements {

      // Geneate redirection type
      var redirectionType string = "permanent"
      if alias.redirectionType == "302" {
        redirectionType = "redirect"
      }

      // Compute nginx configuration line and append to result
      result = result + fmt.Sprintf("rewrite ^%s$ %s %s;\n", alias.source, alias.destination, redirectionType)
    }
    return []byte(result)
  }
  return []byte{}
}

