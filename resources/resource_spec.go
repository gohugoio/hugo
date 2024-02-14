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

package resources

import (
	"path"
	"sync"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/jsconfig"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/identity"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/tpl"
)

func NewSpec(
	s *helpers.PathSpec,
	common *SpecCommon, // may be nil
	fileCaches filecache.Caches,
	memCache *dynacache.Cache,
	incr identity.Incrementer,
	logger loggers.Logger,
	errorHandler herrors.ErrorSender,
	execHelper *hexec.Exec,
) (*Spec, error) {
	conf := s.Cfg.GetConfig().(*allconfig.Config)
	imgConfig := conf.Imaging

	imaging, err := images.NewImageProcessor(imgConfig)
	if err != nil {
		return nil, err
	}

	if incr == nil {
		incr = &identity.IncrementByOne{}
	}

	if logger == nil {
		logger = loggers.NewDefault()
	}

	permalinks, err := page.NewPermalinkExpander(s.URLize, conf.Permalinks)
	if err != nil {
		return nil, err
	}

	if common == nil {
		common = &SpecCommon{
			incr:       incr,
			FileCaches: fileCaches,
			PostBuildAssets: &PostBuildAssets{
				PostProcessResources: make(map[string]postpub.PostPublishedResource),
				JSConfigBuilder:      jsconfig.NewBuilder(),
			},
		}
	}

	rs := &Spec{
		PathSpec:    s,
		Logger:      logger,
		ErrorSender: errorHandler,
		imaging:     imaging,
		ImageCache: newImageCache(
			fileCaches.ImageCache(),
			memCache,
			s,
		),
		ExecHelper: execHelper,

		Permalinks: permalinks,

		SpecCommon: common,
	}

	rs.ResourceCache = newResourceCache(rs, memCache)

	return rs, nil
}

type Spec struct {
	*helpers.PathSpec

	Logger      loggers.Logger
	ErrorSender herrors.ErrorSender

	TextTemplates tpl.TemplateParseFinder

	Permalinks page.PermalinkExpander

	ImageCache *ImageCache

	// Holds default filter settings etc.
	imaging *images.ImageProcessor

	ExecHelper *hexec.Exec

	*SpecCommon
}

// The parts of Spec that's common for all sites.
type SpecCommon struct {
	incr          identity.Incrementer
	ResourceCache *ResourceCache
	FileCaches    filecache.Caches

	// Assets used after the build is done.
	// This is shared between all sites.
	*PostBuildAssets
}

type PostBuildAssets struct {
	postProcessMu        sync.RWMutex
	PostProcessResources map[string]postpub.PostPublishedResource
	JSConfigBuilder      *jsconfig.Builder
}

// NewResource creates a new Resource from the given ResourceSourceDescriptor.
func (r *Spec) NewResource(rd ResourceSourceDescriptor) (resource.Resource, error) {
	if err := rd.init(r); err != nil {
		return nil, err
	}

	dir, name := path.Split(rd.TargetPath)
	dir = paths.ToSlashPreserveLeading(dir)
	if dir == "/" {
		dir = ""
	}
	rp := internal.ResourcePaths{
		File:            name,
		Dir:             dir,
		BaseDirTarget:   rd.BasePathTargetPath,
		BaseDirLink:     rd.BasePathRelPermalink,
		TargetBasePaths: rd.TargetBasePaths,
	}

	staler := &AtomicStaler{}
	var staleInfo resource.StaleInfo = staler
	if rd.StaleInfo != nil {
		staleInfo = resource.StaleInfoFunc(func() bool {
			return staler.IsStale() || rd.StaleInfo.IsStale()
		})
	}

	gr := &genericResource{
		StaleMarker:             staler,
		StaleInfo:               staleInfo,
		h:                       &resourceHash{},
		publishInit:             &sync.Once{},
		foreachIdentityProvider: identity.NopForEeachIdentityProvider,
		paths:                   rp,
		spec:                    r,
		sd:                      rd,
		params:                  make(map[string]any),
		name:                    rd.Name,
		title:                   rd.Name,
	}

	if r.Cfg.Watching() {
		forEachIdentity := identity.ForEeachIdentityProviderFunc(func(f func(identity.Identity) bool) {
			if f(gr.GetIdentityGroup()) {
				return
			}
			if f(rd.DependencyManager) {
				return
			}
		})
		gr.foreachIdentityProvider = identity.ChainedForEeachIdentityProvider(forEachIdentity, rd.ForEeachIdentity)
	}


	if rd.MediaType.MainType == "image" {
		imgFormat, ok := images.ImageFormatFromMediaSubType(rd.MediaType.SubType)
		if ok {
			ir := &imageResource{
				Image:        images.NewImage(imgFormat, r.imaging, nil, gr),
				baseResource: gr,
			}
			ir.root = ir
			return newResourceAdapter(gr.spec, rd.LazyPublish, ir), nil
		}

	}

	return newResourceAdapter(gr.spec, rd.LazyPublish, gr), nil
}

func (r *Spec) MediaTypes() media.Types {
	return r.Cfg.GetConfigSection("mediaTypes").(media.Types)
}

func (r *Spec) OutputFormats() output.Formats {
	return r.Cfg.GetConfigSection("outputFormats").(output.Formats)
}

func (r *Spec) BuildConfig() config.BuildConfig {
	return r.Cfg.GetConfigSection("build").(config.BuildConfig)
}

func (s *Spec) String() string {
	return "spec"
}
