// Copyright 2019 The Hugo Authors. All rights reserved.
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

package deploy

import (
	"fmt"
	"regexp"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/config"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
	"github.com/mitchellh/mapstructure"
)

const deploymentConfigKey = "deployment"

// deployConfig is the complete configuration for deployment.
type deployConfig struct {
	Targets  []*target
	Matchers []*matcher
	Order    []string

	ordering []*regexp.Regexp // compiled Order
}

type target struct {
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
	includeGlob glob.Glob
	excludeGlob glob.Glob
}

func (tgt *target) parseIncludeExclude() error {
	var err error
	if tgt.Include != "" {
		tgt.includeGlob, err = hglob.GetGlob(tgt.Include)
		if err != nil {
			return fmt.Errorf("invalid deployment.target.include %q: %v", tgt.Include, err)
		}
	}
	if tgt.Exclude != "" {
		tgt.excludeGlob, err = hglob.GetGlob(tgt.Exclude)
		if err != nil {
			return fmt.Errorf("invalid deployment.target.exclude %q: %v", tgt.Exclude, err)
		}
	}
	return nil
}

// matcher represents configuration to be applied to files whose paths match
// a specified pattern.
type matcher struct {
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

	// re is Pattern compiled.
	re *regexp.Regexp
}

func (m *matcher) Matches(path string) bool {
	return m.re.MatchString(path)
}

// decode creates a config from a given Hugo configuration.
func decodeConfig(cfg config.Provider) (deployConfig, error) {
	var dcfg deployConfig
	if !cfg.IsSet(deploymentConfigKey) {
		return dcfg, nil
	}
	if err := mapstructure.WeakDecode(cfg.GetStringMap(deploymentConfigKey), &dcfg); err != nil {
		return dcfg, err
	}
	for _, tgt := range dcfg.Targets {
		if err := tgt.parseIncludeExclude(); err != nil {
			return dcfg, err
		}
	}
	var err error
	for _, m := range dcfg.Matchers {
		m.re, err = regexp.Compile(m.Pattern)
		if err != nil {
			return dcfg, fmt.Errorf("invalid deployment.matchers.pattern: %v", err)
		}
	}
	for _, o := range dcfg.Order {
		re, err := regexp.Compile(o)
		if err != nil {
			return dcfg, fmt.Errorf("invalid deployment.orderings.pattern: %v", err)
		}
		dcfg.ordering = append(dcfg.ordering, re)
	}
	return dcfg, nil
}
