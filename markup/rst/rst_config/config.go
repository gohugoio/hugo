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

// Package rst_config holds reStructuredText related configuration.
package rst_config

import "fmt"

// Default holds Hugo's default reStructuredText configuration.
var Default = Config{
	SyntaxHighlight: "long",
}

// Config configures reStructuredText.
type Config struct {
	// Configures Pygments syntax highlighting. Valid values are "short", "long" (default) and "none".
	SyntaxHighlight string
}

func (c *Config) Init() error {
	if c.SyntaxHighlight != "short" && c.SyntaxHighlight != "long" && c.SyntaxHighlight != "none" {
		return fmt.Errorf("invalid value for syntaxHighlight: %q", c.SyntaxHighlight)
	}
	return nil
}
