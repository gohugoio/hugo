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

package hexec

import (
	_ "embed"
	"encoding/base64"
	"sync"
)

//go:embed esmloader.mjs
var esmLoaderSource string

// nodeESMLoaderImportArg returns a "--import=data:..." argument that installs
// a Node.js ESM resolver hook making NODE_PATH a fallback for failed bare
// imports. See esmloader.mjs for the rationale.
var nodeESMLoaderImportArg = sync.OnceValue(func() string {
	return "--import=data:text/javascript;base64," + base64.StdEncoding.EncodeToString([]byte(esmLoaderSource))
})
