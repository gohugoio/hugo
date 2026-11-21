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

package images

import (
	"fmt"
	"image"
	"image/draw"
	"strings"
	"time"

	"github.com/bep/imagemeta"
	"github.com/gohugoio/gift"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/types"
	"github.com/spf13/cast"
)

var _ ImageMetaProvider = (*metaFilter)(nil)

// ImageMetaProvider provides metadata to write to the encoded image.
type ImageMetaProvider interface {
	// MergeImageMeta merges its metadata into m, overwriting existing values.
	MergeImageMeta(m *imagemeta.Meta)
}

type metaFilter struct {
	meta imagemeta.Meta
}

func newMetaFilter(v any, loc *time.Location) filter {
	if loc == nil {
		loc = time.UTC
	}

	var opts map[string]any
	if !types.IsNil(v) {
		var err error
		opts, err = hmaps.ToStringMapE(v)
		if err != nil {
			panic(fmt.Sprintf("invalid meta options: %s", err))
		}
	}

	var (
		m         imagemeta.Meta
		lat, long *float64
		toString  = func(k string, v any) string {
			s, err := cast.ToStringE(v)
			if err != nil {
				panic(fmt.Sprintf("invalid %s: %s", k, err))
			}
			return s
		}
		toFloat = func(k string, v any) *float64 {
			if v == nil || v == "" {
				return nil
			}
			f, err := cast.ToFloat64E(v)
			if err != nil {
				panic(fmt.Sprintf("invalid %s: %s", k, err))
			}
			return &f
		}
	)

	for k, v := range opts {
		k = strings.ToLower(k)
		switch k {
		case "copyright":
			m.Copyright = toString(k, v)
		case "creator":
			m.Creator = toString(k, v)
		case "title":
			m.Title = toString(k, v)
		case "description":
			m.Description = toString(k, v)
		case "date":
			if v == nil || v == "" {
				continue
			}
			t, err := htime.ToTimeInDefaultLocationE(v, loc)
			if err != nil {
				panic(fmt.Sprintf("invalid date: %s", err))
			}
			m.Date = t
		case "lat":
			lat = toFloat(k, v)
		case "long":
			long = toFloat(k, v)
		default:
			panic(fmt.Sprintf("unknown meta field %q, expected one of copyright, creator, title, description, date, lat, long", k))
		}
	}

	if (lat == nil) != (long == nil) {
		panic("lat and long must be set together")
	}
	// Hugo's .Meta.Lat and .Meta.Long are 0 when the image has no GPS position.
	if lat != nil && (*lat != 0 || *long != 0) {
		m.LatLong = &imagemeta.LatLong{Lat: *lat, Long: *long}
	}
	if err := m.Validate(); err != nil {
		panic(err.Error())
	}

	var date string
	if !m.Date.IsZero() {
		date = m.Date.Format(time.RFC3339)
	}

	return filter{
		Options: newFilterOpts(m.Copyright, m.Creator, m.Title, m.Description, date, m.LatLong),
		Filter:  metaFilter{meta: m},
	}
}

func (f metaFilter) MergeImageMeta(m *imagemeta.Meta) {
	set := func(dst *string, src string) {
		if src != "" {
			*dst = src
		}
	}
	set(&m.Copyright, f.meta.Copyright)
	set(&m.Creator, f.meta.Creator)
	set(&m.Title, f.meta.Title)
	set(&m.Description, f.meta.Description)
	if !f.meta.Date.IsZero() {
		m.Date = f.meta.Date
	}
	if f.meta.LatLong != nil {
		m.LatLong = f.meta.LatLong
	}
}

func (f metaFilter) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	panic("not supported")
}

func (f metaFilter) Bounds(srcBounds image.Rectangle) image.Rectangle {
	panic("not supported")
}
