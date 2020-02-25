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

package hugolib

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime/trace"
	"strings"

	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/output"

	"github.com/pkg/errors"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/helpers"
)

// Build builds all sites. If filesystem events are provided,
// this is considered to be a potential partial rebuild.
func (h *HugoSites) Build(config BuildCfg, events ...fsnotify.Event) error {

	if h.running {
		// Make sure we don't trigger rebuilds in parallel.
		h.runningMu.Lock()
		defer h.runningMu.Unlock()
	}

	ctx, task := trace.NewTask(context.Background(), "Build")
	defer task.End()

	errCollector := h.StartErrorCollector()
	errs := make(chan error)

	go func(from, to chan error) {
		var errors []error
		i := 0
		for e := range from {
			i++
			if i > 50 {
				break
			}
			errors = append(errors, e)
		}
		to <- h.pickOneAndLogTheRest(errors)

		close(to)

	}(errCollector, errs)

	if h.Metrics != nil {
		h.Metrics.Reset()
	}

	h.testCounters = config.testCounters

	// Need a pointer as this may be modified.
	conf := &config

	if conf.whatChanged == nil {
		// Assume everything has changed
		conf.whatChanged = &whatChanged{source: true}
	}

	var prepareErr error

	if !config.PartialReRender {
		prepare := func() error {
			init := func(conf *BuildCfg) error {
				for _, s := range h.Sites {
					s.Deps.BuildStartListeners.Notify()
				}

				if len(events) > 0 {
					// Rebuild
					if err := h.initRebuild(conf); err != nil {
						return errors.Wrap(err, "initRebuild")
					}
				} else {
					if err := h.initSites(conf); err != nil {
						return errors.Wrap(err, "initSites")
					}
				}

				return nil
			}

			var err error

			f := func() {
				err = h.process(conf, init, events...)
			}
			trace.WithRegion(ctx, "process", f)
			if err != nil {
				return errors.Wrap(err, "process")
			}

			f = func() {
				err = h.assemble(conf)
			}
			trace.WithRegion(ctx, "assemble", f)
			if err != nil {
				return err
			}

			return nil
		}

		f := func() {
			prepareErr = prepare()
		}
		trace.WithRegion(ctx, "prepare", f)
		if prepareErr != nil {
			h.SendError(prepareErr)
		}

	}

	if prepareErr == nil {
		var err error
		f := func() {
			err = h.render(conf)
		}
		trace.WithRegion(ctx, "render", f)
		if err != nil {
			h.SendError(err)
		}
	}

	if err := h.postProcess(); err != nil {
		h.SendError(err)
	}

	if h.Metrics != nil {
		var b bytes.Buffer
		h.Metrics.WriteMetrics(&b)

		h.Log.FEEDBACK.Printf("\nTemplate Metrics:\n\n")
		h.Log.FEEDBACK.Print(b.String())
		h.Log.FEEDBACK.Println()
	}

	select {
	// Make sure the channel always gets something.
	case errCollector <- nil:
	default:
	}
	close(errCollector)

	err := <-errs
	if err != nil {
		return err
	}

	if err := h.fatalErrorHandler.getErr(); err != nil {
		return err
	}

	errorCount := h.Log.ErrorCounter.Count()
	if errorCount > 0 {
		return fmt.Errorf("logged %d error(s)", errorCount)
	}

	return nil

}

// Build lifecycle methods below.
// The order listed matches the order of execution.

func (h *HugoSites) initSites(config *BuildCfg) error {
	h.reset(config)

	if config.NewConfig != nil {
		if err := h.createSitesFromConfig(config.NewConfig); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) initRebuild(config *BuildCfg) error {
	if config.NewConfig != nil {
		return errors.New("rebuild does not support 'NewConfig'")
	}

	if config.ResetState {
		return errors.New("rebuild does not support 'ResetState'")
	}

	if !h.running {
		return errors.New("rebuild called when not in watch mode")
	}

	for _, s := range h.Sites {
		s.resetBuildState(config.whatChanged.source)
	}

	h.reset(config)
	h.resetLogs()
	helpers.InitLoggers()

	return nil
}

func (h *HugoSites) process(config *BuildCfg, init func(config *BuildCfg) error, events ...fsnotify.Event) error {
	// We should probably refactor the Site and pull up most of the logic from there to here,
	// but that seems like a daunting task.
	// So for now, if there are more than one site (language),
	// we pre-process the first one, then configure all the sites based on that.

	firstSite := h.Sites[0]

	if len(events) > 0 {
		// This is a rebuild
		return firstSite.processPartial(config, init, events)
	}

	return firstSite.process(*config)

}

func (h *HugoSites) assemble(bcfg *BuildCfg) error {

	if len(h.Sites) > 1 {
		// The first is initialized during process; initialize the rest
		for _, site := range h.Sites[1:] {
			if err := site.initializeSiteInfo(); err != nil {
				return err
			}
		}
	}

	if !bcfg.whatChanged.source {
		return nil
	}

	if err := h.content.AssemblePages(); err != nil {
		return err
	}

	if err := h.createPageCollections(); err != nil {
		return err
	}

	return nil

}

func (h *HugoSites) render(config *BuildCfg) error {
	if _, err := h.init.layouts.Do(); err != nil {
		return err
	}

	siteRenderContext := &siteRenderContext{cfg: config, multihost: h.multihost}

	if !config.PartialReRender {
		h.renderFormats = output.Formats{}
		h.withSite(func(s *Site) error {
			s.initRenderFormats()
			return nil
		})

		for _, s := range h.Sites {
			h.renderFormats = append(h.renderFormats, s.renderFormats...)
		}
	}

	i := 0
	for _, s := range h.Sites {
		for siteOutIdx, renderFormat := range s.renderFormats {
			siteRenderContext.outIdx = siteOutIdx
			siteRenderContext.sitesOutIdx = i
			i++

			select {
			case <-h.Done():
				return nil
			default:
				for _, s2 := range h.Sites {
					// We render site by site, but since the content is lazily rendered
					// and a site can "borrow" content from other sites, every site
					// needs this set.
					s2.rc = &siteRenderingContext{Format: renderFormat}

					if err := s2.preparePagesForRender(s == s2, siteRenderContext.sitesOutIdx); err != nil {
						return err
					}
				}

				if !config.SkipRender {
					if config.PartialReRender {
						if err := s.renderPages(siteRenderContext); err != nil {
							return err
						}
					} else {
						if err := s.render(siteRenderContext); err != nil {
							return err
						}
					}
				}
			}

		}

	}

	if !config.SkipRender {
		if err := h.renderCrossSitesArtifacts(); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) postProcess() error {
	var toPostProcess []resource.OriginProvider
	for _, s := range h.Sites {
		for _, v := range s.ResourceSpec.PostProcessResources {
			toPostProcess = append(toPostProcess, v)
		}
	}

	if len(toPostProcess) == 0 {
		return nil
	}

	workers := para.New(config.GetNumWorkerMultiplier())
	g, _ := workers.Start(context.Background())

	handleFile := func(filename string) error {

		content, err := afero.ReadFile(h.BaseFs.PublishFs, filename)
		if err != nil {
			return err
		}

		k := 0
		changed := false

		for {
			l := bytes.Index(content[k:], []byte(postpub.PostProcessPrefix))
			if l == -1 {
				break
			}
			m := bytes.Index(content[k+l:], []byte(postpub.PostProcessSuffix)) + len(postpub.PostProcessSuffix)

			low, high := k+l, k+l+m

			field := content[low:high]

			forward := l + m

			for i, r := range toPostProcess {
				if r == nil {
					panic(fmt.Sprintf("resource %d to post process is nil", i+1))
				}
				v, ok := r.GetFieldString(string(field))
				if ok {
					content = append(content[:low], append([]byte(v), content[high:]...)...)
					changed = true
					forward = len(v)
					break
				}
			}

			k += forward
		}

		if changed {
			return afero.WriteFile(h.BaseFs.PublishFs, filename, content, 0666)
		}

		return nil

	}

	_ = afero.Walk(h.BaseFs.PublishFs, "", func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, "html") {
			return nil
		}

		g.Run(func() error {
			return handleFile(path)
		})

		return nil
	})

	// Prepare for a new build.
	for _, s := range h.Sites {
		s.ResourceSpec.PostProcessResources = make(map[string]postpub.PostPublishedResource)
	}

	return g.Wait()

}
