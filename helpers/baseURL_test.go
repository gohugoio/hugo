// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseURL(t *testing.T) {
	b, err := newBaseURLFromString("http://example.com")
	require.NoError(t, err)
	require.Equal(t, "http://example.com", b.String())

	p, err := b.WithProtocol("webcal://")
	require.NoError(t, err)
	require.Equal(t, "webcal://example.com", p)

	p, err = b.WithProtocol("webcal")
	require.NoError(t, err)
	require.Equal(t, "webcal://example.com", p)

	_, err = b.WithProtocol("mailto:")
	require.Error(t, err)

	b, err = newBaseURLFromString("mailto:hugo@rules.com")
	require.NoError(t, err)
	require.Equal(t, "mailto:hugo@rules.com", b.String())

	// These are pretty constructed
	p, err = b.WithProtocol("webcal")
	require.NoError(t, err)
	require.Equal(t, "webcal:hugo@rules.com", p)

	p, err = b.WithProtocol("webcal://")
	require.NoError(t, err)
	require.Equal(t, "webcal://hugo@rules.com", p)

	// Test with "non-URLs". Some people will try to use these as a way to get
	// relative URLs working etc.
	b, err = newBaseURLFromString("/")
	require.NoError(t, err)
	require.Equal(t, "/", b.String())

	b, err = newBaseURLFromString("")
	require.NoError(t, err)
	require.Equal(t, "", b.String())

}
