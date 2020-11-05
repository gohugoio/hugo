// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package prettifiers contains prettifiers mapped to MIME types. This package is used
// in the publishing chain.
package prettifiers

import (
	"io"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/transform"

	"github.com/yosssi/gohtml"
)

// Client wraps a prettifier.
type Client struct {
	// Whether output prettifying is enabled (HTML in /public)
	PrettifyOutput bool

	prettifiers map[string]prettifier
}

type prettifier func(input []byte, dst io.Writer) error

// Transformer returns a func that can be used in the transformer publishing chain.
func (m Client) Transformer(mediatype media.Type) transform.Transformer {
	if !m.PrettifyOutput {
		return nil
	}
	prettifier := m.prettifiers[mediatype.Type()]
	if prettifier == nil {
		return nil
	}
	return func(ft transform.FromTo) error {
		return prettifier(ft.From().Bytes(), ft.To())
	}
}

// Prettify tries to prettify the src into dst given a MIME type.
func (m Client) Prettify(mediatype media.Type, dst io.Writer, src io.Reader) error {
	prettifier := m.prettifiers[mediatype.Type()]
	if prettifier == nil {
		// No supported prettifier. Just pass it through.
		_, err := io.Copy(dst, src)
		return err
	}

	var w = gohtml.NewWriter(dst)
	_, err := io.Copy(w, src)
	return err
}

// New creates a new Client with the provided MIME types as the mapping foundation.
// The HTML prettifier is also registered for additional HTML types (AMP etc.) in the
// provided list of output formats.
func New(mediaTypes media.Types, outputFormats output.Formats, cfg config.Provider) (Client, error) {
	conf, err := decodeConfig(cfg)

	if err != nil {
		return Client{}, err
	}

	client := Client{
		PrettifyOutput: conf.PrettifyOutput,
		prettifiers:    make(map[string]prettifier),
	}

	// We use the Type definition of the media types defined in the site if found.

	// TODO: implement other media types (see ../minifiers/minifiers.go)

	// HTML
	if !conf.DisableHTML {
		for _, of := range outputFormats {
			if of.IsHTML {
				client.prettifiers[of.MediaType.Type()] = formatHTML
			}
		}
	}

	return client, nil
}

func formatHTML(input []byte, w io.Writer) error {
	prettified := gohtml.FormatBytes(input)

	n, err := w.Write(prettified)
	if err == nil && n != len(prettified) {
		err = io.ErrShortWrite
	}

	return err
}
