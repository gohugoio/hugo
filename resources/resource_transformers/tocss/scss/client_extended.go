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

// +build extended

package scss

import (
	"github.com/bep/golibsass/libsass"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"
)

type options struct {
	// The options we receive from the end user.
	from Options

	// The options we send to the SCSS library.
	to libsass.Options
}

func (c *Client) ToCSS(res resources.ResourceTransformer, opts Options) (resource.Resource, error) {
	internalOptions := options{
		from: opts,
	}

	// Transfer values from client.
	internalOptions.to.Precision = opts.Precision
	internalOptions.to.OutputStyle = libsass.ParseOutputStyle(opts.OutputStyle)

	if internalOptions.to.Precision == 0 {
		// bootstrap-sass requires 8 digits precision. The libsass default is 5.
		// https://github.com/twbs/bootstrap-sass/blob/master/README.md#sass-number-precision
		internalOptions.to.Precision = 8
	}

	return res.Transform(&toCSSTransformation{c: c, options: internalOptions})
}

type toCSSTransformation struct {
	c       *Client
	options options
}

func (t *toCSSTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey(transformationName, t.options.from)
}
