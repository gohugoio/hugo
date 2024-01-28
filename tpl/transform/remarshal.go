package transform

import (
	"bytes"
	"errors"
	"strings"

	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/cast"
)

// Remarshal is used in the Hugo documentation to convert configuration
// examples from YAML to JSON, TOML (and possibly the other way around).
// The is primarily a helper for the Hugo docs site.
// It is not a general purpose YAML to TOML converter etc., and may
// change without notice if it serves a purpose in the docs.
// Format is one of json, yaml or toml.
func (ns *Namespace) Remarshal(format string, data any) (string, error) {
	var meta map[string]any

	format = strings.TrimSpace(strings.ToLower(format))

	mark, err := toFormatMark(format)
	if err != nil {
		return "", err
	}

	if m, ok := data.(map[string]any); ok {
		meta = m
	} else {
		from, err := cast.ToStringE(data)
		if err != nil {
			return "", err
		}

		from = strings.TrimSpace(from)
		if from == "" {
			return "", nil
		}

		fromFormat := metadecoders.Default.FormatFromContentString(from)
		if fromFormat == "" {
			return "", errors.New("failed to detect format from content")
		}

		meta, err = metadecoders.Default.UnmarshalToMap([]byte(from), fromFormat)
		if err != nil {
			return "", err
		}
	}

	// Make it so 1.0 float64 prints as 1 etc.
	applyMarshalTypes(meta)

	var result bytes.Buffer
	if err := parser.InterfaceToConfig(meta, mark, &result); err != nil {
		return "", err
	}

	return result.String(), nil
}

// The unmarshal/marshal dance is extremely type lossy, and we need
// to make sure that integer types prints as "43" and not "43.0" in
// all formats, hence this hack.
func applyMarshalTypes(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			applyMarshalTypes(t)
		case float64:
			i := int64(t)
			if t == float64(i) {
				m[k] = i
			}
		}
	}
}

func toFormatMark(format string) (metadecoders.Format, error) {
	if f := metadecoders.FormatFromString(format); f != "" {
		return f, nil
	}

	return "", errors.New("failed to detect target data serialization format")
}
