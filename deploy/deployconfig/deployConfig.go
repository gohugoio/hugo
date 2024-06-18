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

package deployconfig

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/config"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
)

const DeploymentConfigKey = "deployment"

// DeployConfig is the complete configuration for deployment.
type DeployConfig struct {
	Targets  []*Target
	Matchers []*Matcher
	Order    []string

	// Usually set via flags.
	// Target deployment Name; defaults to the first one.
	Target string
	// Show a confirm prompt before deploying.
	Confirm bool
	// DryRun will try the deployment without any remote changes.
	DryRun bool
	// Force will re-upload all files.
	Force bool
	// Invalidate the CDN cache listed in the deployment target.
	InvalidateCDN bool
	// MaxDeletes is the maximum number of files to delete.
	MaxDeletes int
	// Number of concurrent workers to use when uploading files.
	Workers int

	Ordering []*regexp.Regexp `json:"-"` // compiled Order
}

type Target struct {
	Name string
	URL  string

	CloudFrontDistributionID string

	// GoogleCloudCDNOrigin specifies the Google Cloud project and CDN origin to
	// invalidate when deploying this target.  It is specified as <project>/<origin>.
	GoogleCloudCDNOrigin string

	// Optional patterns of files to include/exclude for this target.
	// Parsed using github.com/gobwas/glob.
	Include string
	Exclude string

	// Parsed versions of Include/Exclude.
	IncludeGlob glob.Glob `json:"-"`
	ExcludeGlob glob.Glob `json:"-"`

	// If true, any local path matching <dir>/index.html will be mapped to the
	// remote path <dir>/. This does not affect the top-level index.html file,
	// since that would result in an empty path.
	StripIndexHTML bool
}

func (tgt *Target) ParseIncludeExclude() error {
	var err error
	if tgt.Include != "" {
		tgt.IncludeGlob, err = hglob.GetGlob(tgt.Include)
		if err != nil {
			return fmt.Errorf("invalid deployment.target.include %q: %v", tgt.Include, err)
		}
	}
	if tgt.Exclude != "" {
		tgt.ExcludeGlob, err = hglob.GetGlob(tgt.Exclude)
		if err != nil {
			return fmt.Errorf("invalid deployment.target.exclude %q: %v", tgt.Exclude, err)
		}
	}
	return nil
}

// Matcher represents configuration to be applied to files whose paths match
// a specified pattern.
type Matcher struct {
	// Pattern is the string pattern to match against paths.
	// Matching is done against paths converted to use / as the path separator.
	Pattern string

	// CacheControl specifies caching attributes to use when serving the blob.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	CacheControl string

	// ContentEncoding specifies the encoding used for the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	ContentEncoding string

	// ContentType specifies the MIME type of the blob being written.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string

	// Gzip determines whether the file should be gzipped before upload.
	// If so, the ContentEncoding field will automatically be set to "gzip".
	Gzip bool

	// Force indicates that matching files should be re-uploaded. Useful when
	// other route-determined metadata (e.g., ContentType) has changed.
	Force bool

	// Re is Pattern compiled.
	Re *regexp.Regexp `json:"-"`
}

func (m *Matcher) Matches(path string) bool {
	return m.Re.MatchString(path)
}

var DefaultConfig = DeployConfig{
	Workers:       10,
	InvalidateCDN: true,
	MaxDeletes:    256,
}

// DecodeConfig creates a config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (DeployConfig, error) {
	dcfg := DefaultConfig

	if !cfg.IsSet(DeploymentConfigKey) {
		return dcfg, nil
	}
	if err := mapstructure.WeakDecode(cfg.GetStringMap(DeploymentConfigKey), &dcfg); err != nil {
		return dcfg, err
	}

	if dcfg.Workers <= 0 {
		dcfg.Workers = 10
	}

	for _, tgt := range dcfg.Targets {
		if *tgt == (Target{}) {
			return dcfg, errors.New("empty deployment target")
		}
		if err := tgt.ParseIncludeExclude(); err != nil {
			return dcfg, err
		}
	}
	var err error
	for _, m := range dcfg.Matchers {
		if *m == (Matcher{}) {
			return dcfg, errors.New("empty deployment matcher")
		}
		m.Re, err = regexp.Compile(m.Pattern)
		if err != nil {
			return dcfg, fmt.Errorf("invalid deployment.matchers.pattern: %v", err)
		}
	}
	for _, o := range dcfg.Order {
		re, err := regexp.Compile(o)
		if err != nil {
			return dcfg, fmt.Errorf("invalid deployment.orderings.pattern: %v", err)
		}
		dcfg.Ordering = append(dcfg.Ordering, re)
	}

	return dcfg, nil
}
