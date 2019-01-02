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

package integrity

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"html/template"
	"io"

	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

const defaultHashAlgo = "sha256"

// Client contains methods to fingerprint (cachebusting) and other integrity-related
// methods.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type fingerprintTransformation struct {
	algo string
}

func (t *fingerprintTransformation) Key() resources.ResourceTransformationKey {
	return resources.NewResourceTransformationKey("fingerprint", t.algo)
}

// Transform creates a MD5 hash of the Resource content and inserts that hash before
// the extension in the filename.
func (t *fingerprintTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	algo := t.algo

	var h hash.Hash

	switch algo {
	case "md5":
		h = md5.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return fmt.Errorf("unsupported crypto algo: %q, use either md5, sha256 or sha512", algo)
	}

	io.Copy(io.MultiWriter(h, ctx.To), ctx.From)
	d, err := digest(h)
	if err != nil {
		return err
	}

	ctx.Data["Integrity"] = integrity(algo, d)
	ctx.AddOutPathIdentifier("." + hex.EncodeToString(d[:]))
	return nil
}

// Fingerprint applies fingerprinting of the given resource and hash algorithm.
// It defaults to sha256 if none given, and the options are md5, sha256 or sha512.
// The same algo is used for both the fingerprinting part (aka cache busting) and
// the base64-encoded Subresource Integrity hash, so you will have to stay away from
// md5 if you plan to use both.
// See https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
func (c *Client) Fingerprint(res resource.Resource, algo string) (resource.Resource, error) {
	if algo == "" {
		algo = defaultHashAlgo
	}

	return c.rs.Transform(
		res,
		&fingerprintTransformation{algo: algo},
	)
}

func integrity(algo string, sum []byte) template.HTMLAttr {
	encoded := base64.StdEncoding.EncodeToString(sum)
	return template.HTMLAttr(algo + "-" + encoded)
}

func digest(h hash.Hash) ([]byte, error) {
	sum := h.Sum(nil)
	return sum, nil
}
