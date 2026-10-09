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
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/minifiers"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/helpers"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/transform"
	"github.com/gohugoio/hugo/transform/livereloadinject"
	"github.com/gohugoio/hugo/transform/metainject"
	"github.com/gohugoio/hugo/transform/urlreplacers"
)

// Descriptor describes the needed publishing chain for an item.
type Descriptor struct {
	// The content to publish.
	Src io.Reader

	// The OutputFormat of the this content.
	OutputFormat output.Format

	// Where to publish this content. This is a filesystem-relative path.
	TargetPath string

	// Counter for the end build summary.
	StatCounter *uint64

	// Configuration that trigger pre-processing.
	// LiveReload script will be injected if this is != nil
	LiveReloadBaseURL *url.URL

	// Enable to inject the Hugo generated tag in the header. Is currently only
	// injected on the home page for HTML type of output formats.
	AddHugoGeneratorTag bool

	// If set, will replace all relative URLs with this one.
	AbsURLPath string

	// Enable to minify the output using the OutputFormat defined above to
	// pick the correct minifier configuration.
	Minify bool

	// Whether Src contains templates.Defer placeholders.
	HasDeferred bool
}

// DestinationPublisher is the default and currently only publisher in Hugo. This
// publisher prepares and publishes an item to the defined destination, e.g. /public.
type DestinationPublisher struct {
	fs                    afero.Fs
	min                   minifiers.Client
	htmlElementsCollector *htmlElementsCollector

	// Descriptors for files with templates.Defer placeholders, keyed by target path.
	// Transformations are applied in PublishDeferred.
	deferred *hmaps.Cache[string, Descriptor]
}

// NewDestinationPublisher creates a new DestinationPublisher.
func NewDestinationPublisher(rs *resources.Spec, outputFormats output.Formats, mediaTypes media.Types) (pub DestinationPublisher, err error) {
	fs := rs.BaseFs.PublishFs
	cfg := rs.Cfg
	var classCollector *htmlElementsCollector
	if rs.BuildConfig().BuildStats.Enabled() {
		classCollector = newHTMLElementsCollector(rs.BuildConfig().BuildStats)
	}
	pub = DestinationPublisher{fs: fs, htmlElementsCollector: classCollector, deferred: hmaps.NewCache[string, Descriptor]()}
	pub.min, err = minifiers.New(mediaTypes, outputFormats, cfg)
	return
}

// Publish applies any relevant transformations and writes the file
// to its destination, e.g. /public.
func (p DestinationPublisher) Publish(d Descriptor) error {
	if d.TargetPath == "" {
		return errors.New("Publish: must provide a TargetPath")
	}

	transformers := p.createTransformerChain(d)
	collect := p.htmlElementsCollector != nil && d.OutputFormat.IsHTML

	if len(transformers) != 0 && d.HasDeferred {
		// Transform when the placeholders are replaced.
		dd := d
		dd.Src = nil
		dd.StatCounter = nil
		p.deferred.Set(filepath.Clean(d.TargetPath), dd)
		transformers = nil
	}

	return p.publish(d, transformers, collect)
}

// PublishDeferred publishes content with its templates.Defer placeholders replaced.
func (p DestinationPublisher) PublishDeferred(filename string, content []byte) error {
	d, found := p.deferred.Get(filename)
	if !found {
		return afero.WriteFile(p.fs, filename, content, 0o666)
	}
	d.Src = bytes.NewReader(content)
	return p.publish(d, p.createTransformerChain(d), false)
}

func (p DestinationPublisher) publish(d Descriptor, transformers transform.Chain, collect bool) error {
	src := d.Src

	if len(transformers) != 0 {
		b := bp.GetBuffer()
		defer bp.PutBuffer(b)

		if err := transformers.Apply(b, d.Src); err != nil {
			return fmt.Errorf("failed to process %q: %w", d.TargetPath, err)
		}

		// This is now what we write to disk.
		src = b
	}

	f, err := helpers.OpenFileForWriting(p.fs, d.TargetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var w io.Writer = f

	if collect {
		w = io.MultiWriter(w, newHTMLElementsCollectorWriter(p.htmlElementsCollector))
	}

	_, err = io.Copy(w, src)
	if err == nil && d.StatCounter != nil {
		atomic.AddUint64(d.StatCounter, uint64(1))
	}

	return err
}

func (p DestinationPublisher) PublishStats() PublishStats {
	if p.htmlElementsCollector == nil {
		return PublishStats{}
	}

	return PublishStats{
		HTMLElements: p.htmlElementsCollector.getHTMLElements(),
	}
}

type PublishStats struct {
	HTMLElements HTMLElements `json:"htmlElements"`
}

// Publisher publishes a result file.
type Publisher interface {
	Publish(d Descriptor) error
	PublishDeferred(filename string, content []byte) error
	PublishStats() PublishStats
}

// XML transformer := transform.New(urlreplacers.NewAbsURLInXMLTransformer(path))
func (p DestinationPublisher) createTransformerChain(f Descriptor) transform.Chain {
	transformers := transform.NewEmpty()

	isHTML := f.OutputFormat.IsHTML

	if f.AbsURLPath != "" {
		if isHTML {
			transformers = append(transformers, urlreplacers.NewAbsURLTransformer(f.AbsURLPath))
		} else {
			// Assume XML.
			transformers = append(transformers, urlreplacers.NewAbsURLInXMLTransformer(f.AbsURLPath))
		}
	}

	if isHTML {
		if f.LiveReloadBaseURL != nil {
			transformers = append(transformers, livereloadinject.New(f.LiveReloadBaseURL))
		}

		// This is only injected on the home page.
		if f.AddHugoGeneratorTag {
			transformers = append(transformers, metainject.HugoGenerator)
		}

	}

	if p.min.MinifyOutput {
		minifyTransformer := p.min.Transformer(f.OutputFormat.MediaType)
		if minifyTransformer != nil {
			transformers = append(transformers, minifyTransformer)
		}
	}

	return transformers
}
