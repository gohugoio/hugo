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
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/publisher"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/output"

	"errors"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/helpers"
)

func init() {
	// To avoid circular dependencies, we set this here.
	langs.DeprecationFunc = helpers.Deprecated
}

// Build builds all sites. If filesystem events are provided,
// this is considered to be a potential partial rebuild.
func (h *HugoSites) Build(config BuildCfg, events ...fsnotify.Event) error {
	if h == nil {
		return errors.New("cannot build nil *HugoSites")
	}

	if h.Deps == nil {
		return errors.New("cannot build nil *Deps")
	}

	if !config.NoBuildLock {
		unlock, err := h.BaseFs.LockBuild()
		if err != nil {
			return fmt.Errorf("failed to acquire a build lock: %w", err)
		}
		defer unlock()
	}

	infol := h.Log.InfoCommand("build")

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
						return fmt.Errorf("initRebuild: %w", err)
					}
				} else {
					if err := h.initSites(conf); err != nil {
						return fmt.Errorf("initSites: %w", err)
					}
				}

				return nil
			}

			if err := h.process(infol, conf, init, events...); err != nil {
				return fmt.Errorf("process: %w", err)
			}

			if err := h.assemble(infol, conf); err != nil {
				return fmt.Errorf("assemble: %w", err)
			}

			return nil
		}

		if prepareErr = prepare(); prepareErr != nil {
			h.SendError(prepareErr)
		}
	}

	if prepareErr == nil {
		if err := h.render(infol, conf); err != nil {
			h.SendError(fmt.Errorf("render: %w", err))
		}

		if err := h.postRenderOnce(); err != nil {
			h.SendError(fmt.Errorf("postRenderOnce: %w", err))
		}

		if err := h.postProcess(infol); err != nil {
			h.SendError(fmt.Errorf("postProcess: %w", err))
		}
	}

	if h.Metrics != nil {
		var b bytes.Buffer
		h.Metrics.WriteMetrics(&b)

		h.Log.Printf("\nTemplate Metrics:\n\n")
		h.Log.Println(b.String())
	}

	h.StopErrorCollector()

	err := <-errs
	if err != nil {
		return err
	}

	if err := h.fatalErrorHandler.getErr(); err != nil {
		return err
	}

	errorCount := h.Log.LoggCount(logg.LevelError)
	if errorCount > 0 {
		return fmt.Errorf("logged %d error(s)", errorCount)
	}

	return nil
}

// Build lifecycle methods below.
// The order listed matches the order of execution.

func (h *HugoSites) initSites(config *BuildCfg) error {
	h.reset(config)
	return nil
}

func (h *HugoSites) initRebuild(config *BuildCfg) error {
	if config.ResetState {
		return errors.New("rebuild does not support 'ResetState'")
	}

	if !h.Configs.Base.Internal.Watch {
		return errors.New("rebuild called when not in watch mode")
	}

	for _, s := range h.Sites {
		s.resetBuildState(config.whatChanged.source)
	}

	h.reset(config)
	h.resetLogs()

	return nil
}

func (h *HugoSites) process(l logg.LevelLogger, config *BuildCfg, init func(config *BuildCfg) error, events ...fsnotify.Event) error {
	defer h.timeTrack(l, time.Now(), "process")

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

func (h *HugoSites) assemble(l logg.LevelLogger, bcfg *BuildCfg) error {
	defer h.timeTrack(l, time.Now(), "assemble")

	if !bcfg.whatChanged.source {
		return nil
	}

	if err := h.getContentMaps().AssemblePages(); err != nil {
		return err
	}

	if err := h.createPageCollections(); err != nil {
		return err
	}

	return nil
}

func (h *HugoSites) timeTrack(l logg.LevelLogger, start time.Time, name string) {
	elapsed := time.Since(start)
	l.WithField("step", name).WithField("duration", elapsed).Logf("running")
}

func (h *HugoSites) render(l logg.LevelLogger, config *BuildCfg) error {
	defer h.timeTrack(l, time.Now(), "render")
	if _, err := h.init.layouts.Do(context.Background()); err != nil {
		return err
	}

	siteRenderContext := &siteRenderContext{cfg: config, multihost: h.Configs.IsMultihost}

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
		h.currentSite = s
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
		if err := h.renderCrossSitesSitemap(); err != nil {
			return err
		}
		if err := h.renderCrossSitesRobotsTXT(); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) postRenderOnce() error {
	h.postRenderInit.Do(func() {
		conf := h.Configs.Base
		if conf.PrintPathWarnings {
			// We need to do this before any post processing, as that may write to the same files twice
			// and create false positives.
			hugofs.WalkFilesystems(h.Fs.PublishDir, func(fs afero.Fs) bool {
				if dfs, ok := fs.(hugofs.DuplicatesReporter); ok {
					dupes := dfs.ReportDuplicates()
					if dupes != "" {
						h.Log.Warnln("Duplicate target paths:", dupes)
					}
				}
				return false
			})
		}

		if conf.PrintUnusedTemplates {
			unusedTemplates := h.Tmpl().(tpl.UnusedTemplatesProvider).UnusedTemplates()
			for _, unusedTemplate := range unusedTemplates {
				h.Log.Warnf("Template %s is unused, source file %s", unusedTemplate.Name(), unusedTemplate.Filename())
			}
		}

	})
	return nil
}

func (h *HugoSites) postProcess(l logg.LevelLogger) error {
	defer h.timeTrack(l, time.Now(), "postProcess")

	// Make sure to write any build stats to disk first so it's available
	// to the post processors.
	if err := h.writeBuildStats(); err != nil {
		return err
	}

	// This will only be set when js.Build have been triggered with
	// imports that resolves to the project or a module.
	// Write a jsconfig.json file to the project's /asset directory
	// to help JS IntelliSense in VS Code etc.
	if !h.ResourceSpec.BuildConfig().NoJSConfigInAssets && h.BaseFs.Assets.Dirs != nil {
		fi, err := h.BaseFs.Assets.Fs.Stat("")
		if err != nil {
			h.Log.Warnf("Failed to resolve jsconfig.json dir: %s", err)
		} else {
			m := fi.(hugofs.FileMetaInfo).Meta()
			assetsDir := m.SourceRoot
			if strings.HasPrefix(assetsDir, h.Configs.LoadingInfo.BaseConfig.WorkingDir) {
				if jsConfig := h.ResourceSpec.JSConfigBuilder.Build(assetsDir); jsConfig != nil {

					b, err := json.MarshalIndent(jsConfig, "", " ")
					if err != nil {
						h.Log.Warnf("Failed to create jsconfig.json: %s", err)
					} else {
						filename := filepath.Join(assetsDir, "jsconfig.json")
						if h.Configs.Base.Internal.Running {
							h.skipRebuildForFilenamesMu.Lock()
							h.skipRebuildForFilenames[filename] = true
							h.skipRebuildForFilenamesMu.Unlock()
						}
						// Make sure it's  written to the OS fs as this is used by
						// editors.
						if err := afero.WriteFile(hugofs.Os, filename, b, 0666); err != nil {
							h.Log.Warnf("Failed to write jsconfig.json: %s", err)
						}
					}
				}
			}

		}
	}

	var toPostProcess []postpub.PostPublishedResource
	for _, r := range h.ResourceSpec.PostProcessResources {
		toPostProcess = append(toPostProcess, r)
	}

	if len(toPostProcess) == 0 {
		// Nothing more to do.
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

	filenames := h.Deps.BuildState.GetFilenamesWithPostPrefix()
	for _, filename := range filenames {
		filename := filename
		g.Run(func() error {
			return handleFile(filename)
		})
	}

	// Prepare for a new build.
	for _, s := range h.Sites {
		s.ResourceSpec.PostProcessResources = make(map[string]postpub.PostPublishedResource)
	}

	return g.Wait()
}

type publishStats struct {
	CSSClasses string `json:"cssClasses"`
}

func (h *HugoSites) writeBuildStats() error {
	if h.ResourceSpec == nil {
		panic("h.ResourceSpec is nil")
	}
	if !h.ResourceSpec.BuildConfig().BuildStats.Enabled() {
		return nil
	}

	htmlElements := &publisher.HTMLElements{}
	for _, s := range h.Sites {
		stats := s.publisher.PublishStats()
		htmlElements.Merge(stats.HTMLElements)
	}

	htmlElements.Sort()

	stats := publisher.PublishStats{
		HTMLElements: *htmlElements,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(stats)
	if err != nil {
		return err
	}
	js := buf.Bytes()

	filename := filepath.Join(h.Configs.LoadingInfo.BaseConfig.WorkingDir, files.FilenameHugoStatsJSON)

	if existingContent, err := afero.ReadFile(hugofs.Os, filename); err == nil {
		// Check if the content has changed.
		if bytes.Equal(existingContent, js) {
			return nil
		}
	}

	// Make sure it's always written to the OS fs.
	if err := afero.WriteFile(hugofs.Os, filename, js, 0666); err != nil {
		return err
	}

	// Write to the destination as well if it's a in-memory fs.
	if !hugofs.IsOsFs(h.Fs.Source) {
		if err := afero.WriteFile(h.Fs.WorkingDirWritable, filename, js, 0666); err != nil {
			return err
		}
	}

	return nil
}
