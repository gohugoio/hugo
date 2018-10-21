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

package hugolib

import (
	"fmt"

	"github.com/gohugoio/hugo/common/herrors"
	errors "github.com/pkg/errors"
)

func (p *Page) errorf(err error, format string, a ...interface{}) error {
	if herrors.UnwrapErrorWithFileContext(err) != nil {
		// More isn't always better.
		return err
	}
	args := append([]interface{}{p.Lang(), p.pathOrTitle()}, a...)
	format = "[%s] page %q: " + format
	if err == nil {
		errors.Errorf(format, args...)
		return fmt.Errorf(format, args...)
	}
	return errors.Wrapf(err, format, args...)
}

func (p *Page) errWithFileContext(err error) error {

	err, _ = herrors.WithFileContextForFile(
		err,
		p.Filename(),
		p.Filename(),
		p.s.SourceSpec.Fs.Source,
		herrors.SimpleLineMatcher)

	return err
}
