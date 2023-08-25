// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package highlight provides code highlighting.
package highlight

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/markup/converter/hooks"

	"github.com/mitchellh/mapstructure"
)

const (
	lineanchorsKey = "lineanchors"
	lineNosKey     = "linenos"
	hlLinesKey     = "hl_lines"
	linosStartKey  = "linenostart"
	noHlKey        = "nohl"
)

var DefaultConfig = Config{
	// The highlighter style to use.
	// See https://xyproto.github.io/splash/docs/all.html
	Style:              "monokai",
	LineNoStart:        1,
	CodeFences:         true,
	NoClasses:          true,
	LineNumbersInTable: true,
	TabWidth:           4,
}

type Config struct {
	Style string

	CodeFences bool

	// Use inline CSS styles.
	NoClasses bool

	// No highlighting.
	NoHl bool

	// When set, line numbers will be printed.
	LineNos            bool
	LineNumbersInTable bool

	// When set, add links to line numbers
	AnchorLineNos bool
	LineAnchors   string

	// Start the line numbers from this value (default is 1).
	LineNoStart int

	// A space separated list of line numbers, e.g. “3-8 10-20”.
	Hl_Lines string

	// If set, the markup will not be wrapped in any container.
	Hl_inline bool

	// A parsed and ready to use list of line ranges.
	HL_lines_parsed [][2]int `json:"-"`

	// TabWidth sets the number of characters for a tab. Defaults to 4.
	TabWidth int

	GuessSyntax bool
}

func (cfg Config) toHTMLOptions() []html.Option {
	var lineAnchors string
	if cfg.LineAnchors != "" {
		lineAnchors = cfg.LineAnchors + "-"
	}
	options := []html.Option{
		html.TabWidth(cfg.TabWidth),
		html.WithLineNumbers(cfg.LineNos),
		html.BaseLineNumber(cfg.LineNoStart),
		html.LineNumbersInTable(cfg.LineNumbersInTable),
		html.WithClasses(!cfg.NoClasses),
		html.WithLinkableLineNumbers(cfg.AnchorLineNos, lineAnchors),
		html.InlineCode(cfg.Hl_inline),
	}

	if cfg.Hl_Lines != "" || cfg.HL_lines_parsed != nil {
		var ranges [][2]int
		if cfg.HL_lines_parsed != nil {
			ranges = cfg.HL_lines_parsed
		} else {
			var err error
			ranges, err = hlLinesToRanges(cfg.LineNoStart, cfg.Hl_Lines)
			if err != nil {
				ranges = nil
			}
		}

		if ranges != nil {
			options = append(options, html.HighlightLines(ranges))
		}
	}

	return options
}

func applyOptions(opts any, cfg *Config) error {
	if opts == nil {
		return nil
	}
	switch vv := opts.(type) {
	case map[string]any:
		return applyOptionsFromMap(vv, cfg)
	default:
		s, err := cast.ToStringE(opts)
		if err != nil {
			return err
		}
		return applyOptionsFromString(s, cfg)
	}
}

func applyOptionsFromString(opts string, cfg *Config) error {
	optsm, err := parseHighlightOptions(opts)
	if err != nil {
		return err
	}
	return mapstructure.WeakDecode(optsm, cfg)
}

func applyOptionsFromMap(optsm map[string]any, cfg *Config) error {
	normalizeHighlightOptions(optsm)
	return mapstructure.WeakDecode(optsm, cfg)
}

func applyOptionsFromCodeBlockContext(ctx hooks.CodeblockContext, cfg *Config) error {
	if cfg.LineAnchors == "" {
		const lineAnchorPrefix = "hl-"
		// Set it to the ordinal with a prefix.
		cfg.LineAnchors = fmt.Sprintf("%s%d", lineAnchorPrefix, ctx.Ordinal())
	}

	return nil
}

// ApplyLegacyConfig applies legacy config from back when we had
// Pygments.
func ApplyLegacyConfig(cfg config.Provider, conf *Config) error {
	if conf.Style == DefaultConfig.Style {
		if s := cfg.GetString("pygmentsStyle"); s != "" {
			conf.Style = s
		}
	}

	if conf.NoClasses == DefaultConfig.NoClasses && cfg.IsSet("pygmentsUseClasses") {
		conf.NoClasses = !cfg.GetBool("pygmentsUseClasses")
	}

	if conf.CodeFences == DefaultConfig.CodeFences && cfg.IsSet("pygmentsCodeFences") {
		conf.CodeFences = cfg.GetBool("pygmentsCodeFences")
	}

	if conf.GuessSyntax == DefaultConfig.GuessSyntax && cfg.IsSet("pygmentsCodefencesGuessSyntax") {
		conf.GuessSyntax = cfg.GetBool("pygmentsCodefencesGuessSyntax")
	}

	if cfg.IsSet("pygmentsOptions") {
		if err := applyOptionsFromString(cfg.GetString("pygmentsOptions"), conf); err != nil {
			return err
		}
	}

	return nil
}

func parseHighlightOptions(in string) (map[string]any, error) {
	in = strings.Trim(in, " ")
	opts := make(map[string]any)

	if in == "" {
		return opts, nil
	}

	for _, v := range strings.Split(in, ",") {
		keyVal := strings.Split(v, "=")
		key := strings.Trim(keyVal[0], " ")
		if len(keyVal) != 2 {
			return opts, fmt.Errorf("invalid Highlight option: %s", key)
		}
		opts[key] = keyVal[1]

	}

	normalizeHighlightOptions(opts)

	return opts, nil
}

func normalizeHighlightOptions(m map[string]any) {
	if m == nil {
		return
	}

	// lowercase all keys
	for k, v := range m {
		delete(m, k)
		m[strings.ToLower(k)] = v
	}

	baseLineNumber := 1
	if v, ok := m[linosStartKey]; ok {
		baseLineNumber = cast.ToInt(v)
	}

	for k, v := range m {
		switch k {
		case noHlKey:
			m[noHlKey] = cast.ToBool(v)
		case lineNosKey:
			if v == "table" || v == "inline" {
				m["lineNumbersInTable"] = v == "table"
			}
			if vs, ok := v.(string); ok {
				m[k] = vs != "false"
			}
		case hlLinesKey:
			if hlRanges, ok := v.([][2]int); ok {
				for i := range hlRanges {
					hlRanges[i][0] += baseLineNumber
					hlRanges[i][1] += baseLineNumber
				}
				delete(m, k)
				m[k+"_parsed"] = hlRanges
			}
		}
	}
}

// startLine compensates for https://github.com/alecthomas/chroma/issues/30
func hlLinesToRanges(startLine int, s string) ([][2]int, error) {
	var ranges [][2]int
	s = strings.TrimSpace(s)

	if s == "" {
		return ranges, nil
	}

	// Variants:
	// 1 2 3 4
	// 1-2 3-4
	// 1-2 3
	// 1 3-4
	// 1    3-4
	fields := strings.Split(s, " ")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		numbers := strings.Split(field, "-")
		var r [2]int
		first, err := strconv.Atoi(numbers[0])
		if err != nil {
			return ranges, err
		}
		first = first + startLine - 1
		r[0] = first
		if len(numbers) > 1 {
			second, err := strconv.Atoi(numbers[1])
			if err != nil {
				return ranges, err
			}
			second = second + startLine - 1
			r[1] = second
		} else {
			r[1] = first
		}

		ranges = append(ranges, r)
	}
	return ranges, nil
}
