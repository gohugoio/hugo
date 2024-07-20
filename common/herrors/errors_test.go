// Copyright 2024 The Hugo Authors. All rights reserved.
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

package herrors

import (
	"errors"
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/afero"
)

func TestIsNotExist(t *testing.T) {
	c := qt.New(t)

	c.Assert(IsNotExist(afero.ErrFileNotFound), qt.Equals, true)
	c.Assert(IsNotExist(afero.ErrFileExists), qt.Equals, false)
	c.Assert(IsNotExist(afero.ErrDestinationExists), qt.Equals, false)
	c.Assert(IsNotExist(nil), qt.Equals, false)

	c.Assert(IsNotExist(fmt.Errorf("foo")), qt.Equals, false)

	// os.IsNotExist returns false for wrapped errors.
	c.Assert(IsNotExist(fmt.Errorf("foo: %w", afero.ErrFileNotFound)), qt.Equals, true)
}

func TestIsFeatureNotAvailableError(t *testing.T) {
	c := qt.New(t)

	c.Assert(IsFeatureNotAvailableError(ErrFeatureNotAvailable), qt.Equals, true)
	c.Assert(IsFeatureNotAvailableError(&FeatureNotAvailableError{}), qt.Equals, true)
	c.Assert(IsFeatureNotAvailableError(errors.New("asdf")), qt.Equals, false)
}
