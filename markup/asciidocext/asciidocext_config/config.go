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

// Package asciidocext_config holds asciidoc related configuration.
package asciidocext_config

var (
	// Default holds Hugo's default asciidoc configuration.
	Default = Config{
		Backend:              "html5",
		Extensions:           []string{},
		Attributes:           map[string]string{},
		NoHeaderOrFooter:     true,
		SafeMode:             "unsafe",
		SectionNumbers:       false,
		Verbose:              false,
		Trace:                false,
		FailureLevel:         "fatal",
		WorkingFolderCurrent: false,
		PreserveTOC:          false,
	}

	// CliDefault holds Asciidoctor CLI defaults (see https://asciidoctor.org/docs/user-manual/)
	CliDefault = Config{
		Backend:      "html5",
		SafeMode:     "unsafe",
		FailureLevel: "fatal",
	}

	AllowedSafeMode = map[string]bool{
		"unsafe": true,
		"safe":   true,
		"server": true,
		"secure": true,
	}

	AllowedFailureLevel = map[string]bool{
		"fatal": true,
		"warn":  true,
	}

	AllowedBackend = map[string]bool{
		"html5":     true,
		"html5s":    true,
		"xhtml5":    true,
		"docbook5":  true,
		"docbook45": true,
		"manpage":   true,
	}

	DisallowedAttributes = map[string]bool{
		"outdir": true,
	}
)

// Config configures asciidoc.
type Config struct {
	Backend              string `jsonschema:"enum=html5,enum=html5s,enum=xhtml5,enum=docbook5,enum=docbook45,enum=manpage"`
	Extensions           []string
	Attributes           map[string]string
	NoHeaderOrFooter     bool
	SafeMode             string `jsonschema:"enum=unsafe,enum=safe,enum=server,enum=secure"`
	SectionNumbers       bool
	Verbose              bool
	Trace                bool
	FailureLevel         string `jsonschema:"enum=fatal,enum=warn"`
	WorkingFolderCurrent bool
	PreserveTOC          bool
}
