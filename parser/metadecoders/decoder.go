// Copyright 2025 The Hugo Authors. All rights reserved.
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

package metadecoders

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/niklasfasching/go-org/org"

	xml "github.com/clbanning/mxj/v2"
	yaml "github.com/goccy/go-yaml"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

// Decoder provides some configuration options for the decoders.
type Decoder struct {
	// Format specifies a specific format to decode from. If empty or
	// unspecified, it's inferred from the contents or the filename.
	Format string

	// Delimiter is the field delimiter. Used in the CSV decoder. Default is
	// ','.
	Delimiter rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored. Used in the
	// CSV decoder.
	Comment rune

	// If true, a quote may appear in an unquoted field and a non-doubled quote
	// may appear in a quoted field. Used in the CSV decoder. Default is false.
	LazyQuotes bool

	// The target data type, either slice or map. Used in the CSV decoder.
	// Default is slice.
	TargetType string
}

// OptionsKey is used in cache keys.
func (d Decoder) OptionsKey() string {
	var sb strings.Builder
	sb.WriteString(d.Format)
	sb.WriteRune(d.Delimiter)
	sb.WriteRune(d.Comment)
	sb.WriteString(strconv.FormatBool(d.LazyQuotes))
	sb.WriteString(d.TargetType)
	return sb.String()
}

// Default is a Decoder in its default configuration.
var Default = Decoder{
	Delimiter:  ',',
	TargetType: "slice",
}

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func (d Decoder) UnmarshalToMap(data []byte, f Format) (map[string]any, error) {
	m := make(map[string]any)
	if data == nil {
		return m, nil
	}

	err := d.UnmarshalTo(data, f, &m)

	if m == nil {
		// We migrated to github.com/goccy/go-yaml in v0.152.0,
		// which produces nil maps for empty YAML files (and empty map nodes), unlike gopkg.in/yaml.v2.
		//
		// To prevent crashes when trying to handle empty config files etc., we ensure we always return a non-nil map here.
		// See issue 14074.
		m = make(map[string]any)
	}

	return m, err
}

// UnmarshalFileToMap is the same as UnmarshalToMap, but reads the data from
// the given filename.
func (d Decoder) UnmarshalFileToMap(fs afero.Fs, filename string) (map[string]any, error) {
	format := FormatFromString(filename)
	if format == "" {
		return nil, fmt.Errorf("%q is not a valid configuration format", filename)
	}

	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	return d.UnmarshalToMap(data, format)
}

// UnmarshalStringTo tries to unmarshal data to a new instance of type typ.
func (d Decoder) UnmarshalStringTo(data string, typ any) (any, error) {
	data = strings.TrimSpace(data)
	// We only check for the possible types in YAML, JSON and TOML.
	switch typ.(type) {
	case string:
		return data, nil
	case map[string]any, maps.Params:
		format := d.FormatFromContentString(data)
		return d.UnmarshalToMap([]byte(data), format)
	case []any:
		// A standalone slice. Let YAML handle it.
		return d.Unmarshal([]byte(data), YAML)
	case bool:
		return cast.ToBoolE(data)
	case int:
		return cast.ToIntE(data)
	case int64:
		return cast.ToInt64E(data)
	case float64:
		return cast.ToFloat64E(data)
	default:
		return nil, fmt.Errorf("unmarshal: %T not supported", typ)
	}
}

// Unmarshal will unmarshall data in format f into an interface{}.
// This is what's needed for Hugo's /data handling.
func (d Decoder) Unmarshal(data []byte, f Format) (any, error) {
	if len(data) == 0 {
		switch f {
		case CSV:
			switch d.TargetType {
			case "map":
				return make(map[string]any), nil
			case "slice":
				return make([][]string, 0), nil
			default:
				return nil, fmt.Errorf("invalid targetType: expected either slice or map, received %s", d.TargetType)
			}
		default:
			return make(map[string]any), nil
		}
	}
	var v any
	err := d.UnmarshalTo(data, f, &v)

	return v, err
}

// UnmarshalYaml unmarshals data in YAML format into v.
func UnmarshalYaml(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return err
	}
	if err := validateAliasLimitForCollections(v, calculateCollectionAliasLimit(len(data))); err != nil {
		return err
	}

	return nil
}

// The Billion Laughs YAML example is about 500 bytes in size,
// but even halving that when converted to JSON would produce a file of about 4 MB in size,
// which, when repeated enough times, could be disruptive.
// For large data files where every row shares a common map via aliases,
// a large number of aliases could make sense.
// The primary goal here is to catch the small but malicious files.
func calculateCollectionAliasLimit(sizeInBytes int) int {
	sizeInKB := sizeInBytes / 1024
	if sizeInKB == 0 {
		sizeInKB = 1
	}
	if sizeInKB < 2 {
		// This should allow at most "thousand laughs",
		// which should be plenty of room for legitimate uses.
		return 100
	}

	// The numbers below are somewhat arbitrary, but should provide
	// a reasonable trade-off between safety and usability.
	if sizeInKB < 10 {
		return 5000
	}
	return 10000
}

// Used in benchmarks.
func unmarshalYamlNoValidation(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// See https://github.com/goccy/go-yaml/issues/461
// While it's true that yaml.Unmarshal isn't vulnerable to the Billion Laughs attack,
// we can easily get a delayed laughter when we try to render this very big structure later,
// e.g. via RenderString.
func validateAliasLimitForCollections(v any, limit int) error {
	if limit <= 0 {
		limit = 1000
	}

	collectionRefCounts := make(map[uintptr]int)

	checkCollectionRef := func(v *any) error {
		// Conversion of a Pointer to a uintptr (but not back to Pointer) is considered safe.
		// See https://pkg.go.dev/unsafe#pkg-functions
		ptr := uintptr(unsafe.Pointer(v))
		if ptr == 0 {
			return nil
		}
		collectionRefCounts[ptr]++
		if collectionRefCounts[ptr] > limit {
			return fmt.Errorf("too many YAML aliases for non-scalar nodes")
		}
		return nil
	}

	var validate func(v any) error
	validate = func(v any) error {
		switch vv := v.(type) {
		case *map[string]any:
			if err := checkCollectionRef(&v); err != nil {
				return err
			}
			for _, vvv := range *vv {
				if err := validate(vvv); err != nil {
					return err
				}
			}
		case map[string]any:
			if err := checkCollectionRef(&v); err != nil {
				return err
			}
			for _, vvv := range vv {
				if err := validate(vvv); err != nil {
					return err
				}
			}
		case []any:
			if err := checkCollectionRef(&v); err != nil {
				return err
			}
			for _, vvv := range vv {
				if err := validate(vvv); err != nil {
					return err
				}
			}
		case *any:
			return validate(*vv)
		default:
			// ok
		}
		return nil
	}

	return validate(v)
}

// UnmarshalTo unmarshals data in format f into v.
func (d Decoder) UnmarshalTo(data []byte, f Format, v any) error {
	var err error

	switch f {
	case ORG:
		err = d.unmarshalORG(data, v)
	case JSON:
		err = json.Unmarshal(data, v)
	case XML:
		var xmlRoot xml.Map
		xmlRoot, err = xml.NewMapXml(data)

		var xmlValue map[string]any
		if err == nil {
			xmlRootName, err := xmlRoot.Root()
			if err != nil {
				return toFileError(f, data, fmt.Errorf("failed to unmarshal XML: %w", err))
			}

			// Get the root value and verify it's a map
			rootValue := xmlRoot[xmlRootName]
			if rootValue == nil {
				return toFileError(f, data, fmt.Errorf("XML root element '%s' has no value", xmlRootName))
			}

			// Type check before conversion
			mapValue, ok := rootValue.(map[string]any)
			if !ok {
				return toFileError(f, data, fmt.Errorf("XML root element '%s' must be a map/object, got %T", xmlRootName, rootValue))
			}
			xmlValue = mapValue
		}

		switch v := v.(type) {
		case *map[string]any:
			*v = xmlValue
		case *any:
			*v = xmlValue
		}
	case TOML:
		err = toml.Unmarshal(data, v)
	case YAML:
		return UnmarshalYaml(data, v)
	case CSV:
		return d.unmarshalCSV(data, v)

	default:
		return fmt.Errorf("unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return toFileError(f, data, fmt.Errorf("unmarshal failed: %w", err))
}

func (d Decoder) unmarshalCSV(data []byte, v any) error {
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = d.Delimiter
	r.Comment = d.Comment
	r.LazyQuotes = d.LazyQuotes

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	switch vv := v.(type) {
	case *any:
		switch d.TargetType {
		case "map":
			if len(records) < 2 {
				return fmt.Errorf("cannot unmarshal CSV into %T: expected at least a header row and one data row", v)
			}

			seen := make(map[string]bool, len(records[0]))
			for _, fieldName := range records[0] {
				if seen[fieldName] {
					return fmt.Errorf("cannot unmarshal CSV into %T: header row contains duplicate field names", v)
				}
				seen[fieldName] = true
			}

			sm := make([]map[string]string, len(records)-1)
			for i, record := range records[1:] {
				m := make(map[string]string, len(records[0]))
				for j, col := range record {
					m[records[0][j]] = col
				}
				sm[i] = m
			}
			*vv = sm
		case "slice":
			*vv = records
		default:
			return fmt.Errorf("cannot unmarshal CSV into %T: invalid targetType: expected either slice or map, received %s", v, d.TargetType)
		}
	default:
		return fmt.Errorf("cannot unmarshal CSV into %T", v)
	}

	return nil
}

func parseORGDate(s string) string {
	r := regexp.MustCompile(`[<\[](\d{4}-\d{2}-\d{2}) .*[>\]]`)
	if m := r.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return s
}

func (d Decoder) unmarshalORG(data []byte, v any) error {
	config := org.New()
	config.Log = log.Default() // TODO(bep)
	document := config.Parse(bytes.NewReader(data), "")
	if document.Error != nil {
		return document.Error
	}
	frontMatter := make(map[string]any, len(document.BufferSettings))
	for k, v := range document.BufferSettings {
		k = strings.ToLower(k)
		if strings.HasSuffix(k, "[]") {
			frontMatter[k[:len(k)-2]] = strings.Fields(v)
		} else if strings.Contains(v, "\n") {
			frontMatter[k] = strings.Split(v, "\n")
		} else if k == "filetags" {
			trimmed := strings.TrimPrefix(v, ":")
			trimmed = strings.TrimSuffix(trimmed, ":")
			frontMatter[k] = strings.Split(trimmed, ":")
		} else if k == "date" || k == "lastmod" || k == "publishdate" || k == "expirydate" {
			frontMatter[k] = parseORGDate(v)
		} else {
			frontMatter[k] = v
		}
	}
	switch vv := v.(type) {
	case *map[string]any:
		*vv = frontMatter
	case *any:
		*vv = frontMatter
	}
	return nil
}

func toFileError(f Format, data []byte, err error) error {
	return herrors.NewFileErrorFromName(err, fmt.Sprintf("_stream.%s", f)).UpdateContent(bytes.NewReader(data), nil)
}
