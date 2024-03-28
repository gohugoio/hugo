// Copyright 2015 The Hugo Authors. All rights reserved.
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

package parser

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gohugoio/hugo/parser/metadecoders"

	toml "github.com/pelletier/go-toml/v2"

	yaml "gopkg.in/yaml.v2"

	xml "github.com/clbanning/mxj/v2"
)

const (
	yamlDelimLf = "---\n"
	tomlDelimLf = "+++\n"
)

func InterfaceToConfig(in any, format metadecoders.Format, w io.Writer) error {
	if in == nil {
		return errors.New("input was nil")
	}

	switch format {
	case metadecoders.YAML:
		b, err := yaml.Marshal(in)
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		return err

	case metadecoders.TOML:
		enc := toml.NewEncoder(w)
		enc.SetIndentTables(true)
		return enc.Encode(in)
	case metadecoders.JSON:
		b, err := json.MarshalIndent(in, "", "   ")
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'\n'})
		return err
	case metadecoders.XML:
		b, err := xml.AnyXmlIndent(in, "", "\t", "root")
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		return err
	default:
		return errors.New("unsupported Format provided")
	}
}

func InterfaceToFrontMatter(in any, format metadecoders.Format, w io.Writer) error {
	if in == nil {
		return errors.New("input was nil")
	}

	switch format {
	case metadecoders.YAML:
		_, err := w.Write([]byte(yamlDelimLf))
		if err != nil {
			return err
		}

		err = InterfaceToConfig(in, format, w)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(yamlDelimLf))
		return err

	case metadecoders.TOML:
		_, err := w.Write([]byte(tomlDelimLf))
		if err != nil {
			return err
		}

		err = InterfaceToConfig(in, format, w)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(tomlDelimLf))
		return err

	default:
		return InterfaceToConfig(in, format, w)
	}
}
