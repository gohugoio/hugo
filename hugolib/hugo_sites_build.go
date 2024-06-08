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

package hugolib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/hugolib/doctree"
	"github.com/gohugoio/hugo/hugolib/pagesfromdata"
	"github.com/gohugoio/hugo/hugolib/segments"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/publisher"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/para"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/rungroup"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/siteidentities"
	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/spf13/afero"

	"github.com/fsnotify/fsnotify"
)

// Build builds all sites. If filesystem events are provided,
// this is considered to be a potential partial rebuild.
func (h *HugoSites) Build(config BuildCfg, events ...fsnotify.Event) error {
	infol := h.Log.InfoCommand("build")
	defer loggers.TimeTrackf(infol, time.Now(), nil, "")
	defer func() {
		h.buildCounter.Add(1)
	}()

	if h.Deps == nil {
		panic("must have deps")
	}

	if !config.NoBuildLock {
		unlock, err := h.BaseFs.LockBuild()
		if err != nil {
			return fmt.Errorf("failed to acquire a build lock: %w", err)
		}
		defer unlock()
	}

	defer func() {
		for _, s := range h.Sites {
			s.Deps.BuildEndListeners.Notify()
		}
	}()

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

	for _, s := range h.Sites {
		s.state = siteStateInit
	}

	if h.Metrics != nil {
		h.Metrics.Reset()
	}

	h.buildCounters = config.testCounters
	if h.buildCounters == nil {
		h.buildCounters = &buildCounters{}
	}

	// Need a pointer as this may be modified.
	conf := &config
	if conf.WhatChanged == nil {
		// Assume everything has changed
		conf.WhatChanged = &WhatChanged{needsPagesAssembly: true}
	}

	var prepareErr error

	if !config.PartialReRender {
		prepare := func() error {
			init := func(conf *BuildCfg) error {
				for _, s := range h.Sites {
					s.Deps.BuildStartListeners.Notify()
				}

				if len(events) > 0 || len(conf.WhatChanged.Changes()) > 0 {
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

			ctx := context.Background()

			if err := h.process(ctx, infol, conf, init, events...); err != nil {
				return fmt.Errorf("process: %w", err)
			}

			if err := h.assemble(ctx, infol, conf); err != nil {
				return fmt.Errorf("assemble: %w", err)
			}

			return nil
		}

		if prepareErr = prepare(); prepareErr != nil {
			h.SendError(prepareErr)
		}
	}

	for _, s := range h.Sites {
		s.state = siteStateReady
	}

	if prepareErr == nil {
		if err := h.render(infol, conf); err != nil {
			h.SendError(fmt.Errorf("render: %w", err))
		}

		if err := h.postRenderOnce(); err != nil {
			h.SendError(fmt.Errorf("postRenderOnce: %w", err))
		}

		// Make sure to write any build stats to disk first so it's available
		// to the post processors.
		if err := h.writeBuildStats(); err != nil {
			return err
		}

		if err := h.renderDeferred(infol); err != nil {
			h.SendError(fmt.Errorf("renderDeferred: %w", err))
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

	errorCount := h.Log.LoggCount(logg.LevelError) + loggers.Log().LoggCount(logg.LevelError)
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
	if !h.Configs.Base.Internal.Watch {
		return errors.New("rebuild called when not in watch mode")
	}

	h.pageTrees.treePagesResources.WalkPrefixRaw("", func(key string, n contentNodeI) bool {
		n.resetBuildState()
		return false
	})

	for _, s := range h.Sites {
		s.resetBuildState(config.WhatChanged.needsPagesAssembly)
	}

	h.reset(config)
	h.resetLogs()

	return nil
}

// process prepares the Sites' sources for a full or partial rebuild.
// This will also parse the source and create all the Page objects.
func (h *HugoSites) process(ctx context.Context, l logg.LevelLogger, config *BuildCfg, init func(config *BuildCfg) error, events ...fsnotify.Event) error {
	l = l.WithField("step", "process")
	defer loggers.TimeTrackf(l, time.Now(), nil, "")

	if _, err := h.init.layouts.Do(ctx); err != nil {
		return err
	}

	if len(events) > 0 {
		// This is a rebuild triggered from file events.
		return h.processPartialFileEvents(ctx, l, config, init, events)
	} else if len(config.WhatChanged.Changes()) > 0 {
		// Rebuild triggered from remote events.
		if err := init(config); err != nil {
			return err
		}
		return h.processPartialRebuildChanges(ctx, l, config)
	}
	return h.processFull(ctx, l, config)
}

// assemble creates missing sections, applies aggregate values (e.g. dates, cascading params),
// removes disabled pages etc.
func (h *HugoSites) assemble(ctx context.Context, l logg.LevelLogger, bcfg *BuildCfg) error {
	l = l.WithField("step", "assemble")
	defer loggers.TimeTrackf(l, time.Now(), nil, "")

	if !bcfg.WhatChanged.needsPagesAssembly {
		changes := bcfg.WhatChanged.Drain()
		if len(changes) > 0 {
			if err := h.resolveAndClearStateForIdentities(ctx, l, nil, changes); err != nil {
				return err
			}
		}
		return nil
	}

	h.translationKeyPages.Reset()
	assemblers := make([]*sitePagesAssembler, len(h.Sites))
	// Changes detected during assembly (e.g. aggregate date changes)

	for i, s := range h.Sites {
		assemblers[i] = &sitePagesAssembler{
			Site:            s,
			assembleChanges: bcfg.WhatChanged,
			ctx:             ctx,
		}
	}

	g, _ := h.workersSite.Start(ctx)
	for _, s := range assemblers {
		s := s
		g.Run(func() error {
			return s.assemblePagesStep1(ctx)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	changes := bcfg.WhatChanged.Drain()

	// Changes from the assemble step (e.g. lastMod, cascade) needs a re-calculation
	// of what needs to be re-built.
	if len(changes) > 0 {
		if err := h.resolveAndClearStateForIdentities(ctx, l, nil, changes); err != nil {
			return err
		}
	}

	for _, s := range assemblers {
		if err := s.assemblePagesStep2(); err != nil {
			return err
		}
	}

	h.renderFormats = output.Formats{}
	for _, s := range h.Sites {
		s.s.initRenderFormats()
		h.renderFormats = append(h.renderFormats, s.renderFormats...)
	}

	for _, s := range assemblers {
		if err := s.assemblePagesStepFinal(); err != nil {
			return err
		}
	}

	return nil
}

// render renders the sites.
func (h *HugoSites) render(l logg.LevelLogger, config *BuildCfg) error {
	l = l.WithField("step", "render")
	start := time.Now()
	defer func() {
		loggers.TimeTrackf(l, start, h.buildCounters.loggFields(), "")
	}()

	siteRenderContext := &siteRenderContext{cfg: config, multihost: h.Configs.IsMultihost}

	i := 0
	for _, s := range h.Sites {
		segmentFilter := s.conf.C.SegmentFilter
		if segmentFilter.ShouldExcludeCoarse(segments.SegmentMatcherFields{Lang: s.language.Lang}) {
			l.Logf("skip language %q not matching segments set in --renderSegments", s.language.Lang)
			continue
		}

		siteRenderContext.languageIdx = s.languagei
		h.currentSite = s
		for siteOutIdx, renderFormat := range s.renderFormats {
			if segmentFilter.ShouldExcludeCoarse(segments.SegmentMatcherFields{Output: renderFormat.Name, Lang: s.language.Lang}) {
				l.Logf("skip output format %q for language %q not matching segments set in --renderSegments", renderFormat.Name, s.language.Lang)
				continue
			}

			if err := func() error {
				rc := tpl.RenderingContext{Site: s, SiteOutIdx: siteOutIdx}
				h.BuildState.StartStageRender(rc)
				defer h.BuildState.StopStageRender(rc)

				siteRenderContext.outIdx = siteOutIdx
				siteRenderContext.sitesOutIdx = i
				i++

				select {
				case <-h.Done():
					return nil
				default:
					for _, s2 := range h.Sites {
						if err := s2.preparePagesForRender(s == s2, siteRenderContext.sitesOutIdx); err != nil {
							return err
						}
					}
					if !config.SkipRender {
						ll := l.WithField("substep", "pages").
							WithField("site", s.language.Lang).
							WithField("outputFormat", renderFormat.Name)

						start := time.Now()

						if config.PartialReRender {
							if err := s.renderPages(siteRenderContext); err != nil {
								return err
							}
						} else {
							if err := s.render(siteRenderContext); err != nil {
								return err
							}
						}
						loggers.TimeTrackf(ll, start, nil, "")
					}
				}
				return nil
			}(); err != nil {
				return err
			}

		}
	}

	return nil
}

func (h *HugoSites) renderDeferred(l logg.LevelLogger) error {
	l = l.WithField("step", "render deferred")
	start := time.Now()

	var deferredCount int

	for rc, de := range h.Deps.BuildState.DeferredExecutionsGroupedByRenderingContext {
		if de.FilenamesWithPostPrefix.Len() == 0 {
			continue
		}

		deferredCount += de.FilenamesWithPostPrefix.Len()

		s := rc.Site.(*Site)
		for _, s2 := range h.Sites {
			if err := s2.preparePagesForRender(s == s2, rc.SiteOutIdx); err != nil {
				return err
			}
		}
		if err := s.executeDeferredTemplates(de); err != nil {
			return herrors.ImproveRenderErr(err)
		}
	}

	loggers.TimeTrackf(l, start, logg.Fields{
		logg.Field{Name: "count", Value: deferredCount},
	}, "")

	return nil
}

func (s *Site) executeDeferredTemplates(de *deps.DeferredExecutions) error {
	handleFile := func(filename string) error {
		content, err := afero.ReadFile(s.BaseFs.PublishFs, filename)
		if err != nil {
			return err
		}

		k := 0
		changed := false

		for {
			if k >= len(content) {
				break
			}
			l := bytes.Index(content[k:], []byte(tpl.HugoDeferredTemplatePrefix))
			if l == -1 {
				break
			}
			m := bytes.Index(content[k+l:], []byte(tpl.HugoDeferredTemplateSuffix)) + len(tpl.HugoDeferredTemplateSuffix)

			low, high := k+l, k+l+m

			forward := l + m
			id := string(content[low:high])

			if err := func() error {
				deferred, found := de.Executions.Get(id)
				if !found {
					panic(fmt.Sprintf("deferred execution with id %q not found", id))
				}
				deferred.Mu.Lock()
				defer deferred.Mu.Unlock()

				if !deferred.Executed {
					tmpl := s.Deps.Tmpl()
					templ, found := tmpl.Lookup(deferred.TemplateName)
					if !found {
						panic(fmt.Sprintf("template %q not found", deferred.TemplateName))
					}

					if err := func() error {
						buf := bufferpool.GetBuffer()
						defer bufferpool.PutBuffer(buf)

						err = tmpl.ExecuteWithContext(deferred.Ctx, templ, buf, deferred.Data)
						if err != nil {
							return err
						}
						deferred.Result = buf.String()
						deferred.Executed = true

						return nil
					}(); err != nil {
						return err
					}
				}

				content = append(content[:low], append([]byte(deferred.Result), content[high:]...)...)
				changed = true

				return nil
			}(); err != nil {
				return err
			}

			k += forward
		}

		if changed {
			return afero.WriteFile(s.BaseFs.PublishFs, filename, content, 0o666)
		}

		return nil
	}

	g := rungroup.Run[string](context.Background(), rungroup.Config[string]{
		NumWorkers: s.h.numWorkers,
		Handle: func(ctx context.Context, filename string) error {
			return handleFile(filename)
		},
	})

	de.FilenamesWithPostPrefix.ForEeach(func(filename string, _ bool) {
		g.Enqueue(filename)
	})

	return g.Wait()
}

// / postRenderOnce runs some post processing that only needs to be done once, e.g. printing of unused templates.
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

// postProcess runs the post processors, e.g. writing the hugo_stats.json file.
func (h *HugoSites) postProcess(l logg.LevelLogger) error {
	l = l.WithField("step", "postProcess")
	defer loggers.TimeTrackf(l, time.Now(), nil, "")

	// This will only be set when js.Build have been triggered with
	// imports that resolves to the project or a module.
	// Write a jsconfig.json file to the project's /asset directory
	// to help JS IntelliSense in VS Code etc.
	if !h.ResourceSpec.BuildConfig().NoJSConfigInAssets {
		handleJSConfig := func(fi os.FileInfo) {
			m := fi.(hugofs.FileMetaInfo).Meta()
			if !m.IsProject {
				return
			}

			if jsConfig := h.ResourceSpec.JSConfigBuilder.Build(m.SourceRoot); jsConfig != nil {
				b, err := json.MarshalIndent(jsConfig, "", " ")
				if err != nil {
					h.Log.Warnf("Failed to create jsconfig.json: %s", err)
				} else {
					filename := filepath.Join(m.SourceRoot, "jsconfig.json")
					if h.Configs.Base.Internal.Running {
						h.skipRebuildForFilenamesMu.Lock()
						h.skipRebuildForFilenames[filename] = true
						h.skipRebuildForFilenamesMu.Unlock()
					}
					// Make sure it's  written to the OS fs as this is used by
					// editors.
					if err := afero.WriteFile(hugofs.Os, filename, b, 0o666); err != nil {
						h.Log.Warnf("Failed to write jsconfig.json: %s", err)
					}
				}
			}
		}

		fi, err := h.BaseFs.Assets.Fs.Stat("")
		if err != nil {
			if !herrors.IsNotExist(err) {
				h.Log.Warnf("Failed to resolve jsconfig.json dir: %s", err)
			}
		} else {
			handleJSConfig(fi)
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
			return afero.WriteFile(h.BaseFs.PublishFs, filename, content, 0o666)
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
	if err := afero.WriteFile(hugofs.Os, filename, js, 0o666); err != nil {
		return err
	}

	// Write to the destination as well if it's a in-memory fs.
	if !hugofs.IsOsFs(h.Fs.Source) {
		if err := afero.WriteFile(h.Fs.WorkingDirWritable, filename, js, 0o666); err != nil {
			return err
		}
	}

	// This step may be followed by a post process step that may
	// rebuild e.g. CSS, so clear any cache that's defined for the hugo_stats.json.
	h.dynacacheGCFilenameIfNotWatchedAndDrainMatching(filename)

	return nil
}

type pathChange struct {
	// The path to the changed file.
	p *paths.Path

	// If true, this is a delete operation (a delete or a rename).
	delete bool

	// If true, this is a directory.
	isDir bool
}

func (p pathChange) isStructuralChange() bool {
	return p.delete || p.isDir
}

func (h *HugoSites) processPartialRebuildChanges(ctx context.Context, l logg.LevelLogger, config *BuildCfg) error {
	if err := h.resolveAndClearStateForIdentities(ctx, l, nil, config.WhatChanged.Drain()); err != nil {
		return err
	}

	if err := h.processContentAdaptersOnRebuild(ctx, config); err != nil {
		return err
	}
	return nil
}

// processPartialFileEvents prepares the Sites' sources for a partial rebuild.
func (h *HugoSites) processPartialFileEvents(ctx context.Context, l logg.LevelLogger, config *BuildCfg, init func(config *BuildCfg) error, events []fsnotify.Event) error {
	h.Log.Trace(logg.StringFunc(func() string {
		var sb strings.Builder
		sb.WriteString("File events:\n")
		for _, ev := range events {
			sb.WriteString(ev.String())
			sb.WriteString("\n")
		}
		return sb.String()
	}))

	// For a list of events for the different OSes, see the test output in https://github.com/bep/fsnotifyeventlister/.
	events = h.fileEventsFilter(events)
	events = h.fileEventsTrim(events)
	eventInfos := h.fileEventsApplyInfo(events)

	logger := h.Log

	var (
		tmplAdded          bool
		tmplChanged        bool
		i18nChanged        bool
		needsPagesAssemble bool
	)

	changedPaths := struct {
		changedFiles []*paths.Path
		changedDirs  []*paths.Path
		deleted      []*paths.Path
	}{}

	removeDuplicatePaths := func(ps []*paths.Path) []*paths.Path {
		seen := make(map[string]bool)
		var filtered []*paths.Path
		for _, p := range ps {
			if !seen[p.Path()] {
				seen[p.Path()] = true
				filtered = append(filtered, p)
			}
		}
		return filtered
	}

	var (
		cacheBusters      []func(string) bool
		deletedDirs       []string
		addedContentPaths []*paths.Path
	)

	for _, ev := range eventInfos {
		cpss := h.BaseFs.ResolvePaths(ev.Name)
		pss := make([]*paths.Path, len(cpss))
		for i, cps := range cpss {
			p := cps.Path
			if ev.removed && !paths.HasExt(p) {
				// Assume this is a renamed/removed directory.
				// For deletes, we walk up the tree to find the container (e.g. branch bundle),
				// so we will catch this even if it is a file without extension.
				// This avoids us walking up to the home page bundle for the common case
				// of renaming root sections.
				p = p + "/_index.md"
				deletedDirs = append(deletedDirs, cps.Path)
			}

			pss[i] = h.Configs.ContentPathParser.Parse(cps.Component, p)
			if ev.added && !ev.isChangedDir && cps.Component == files.ComponentFolderContent {
				addedContentPaths = append(addedContentPaths, pss[i])
			}

			// Compile cache buster.
			np := glob.NormalizePath(path.Join(cps.Component, cps.Path))
			g, err := h.ResourceSpec.BuildConfig().MatchCacheBuster(h.Log, np)
			if err == nil && g != nil {
				cacheBusters = append(cacheBusters, g)
			}
		}

		if ev.removed {
			changedPaths.deleted = append(changedPaths.deleted, pss...)
		} else if ev.isChangedDir {
			changedPaths.changedDirs = append(changedPaths.changedDirs, pss...)
		} else {
			changedPaths.changedFiles = append(changedPaths.changedFiles, pss...)
		}
	}

	var (
		addedOrChangedContent []pathChange
		changes               []identity.Identity
	)

	// Find the most specific identity possible.
	handleChange := func(pathInfo *paths.Path, delete, isDir bool) {
		switch pathInfo.Component() {
		case files.ComponentFolderContent:
			logger.Println("Source changed", pathInfo.Path())
			isContentDataFile := pathInfo.IsContentData()
			if !isContentDataFile {
				if ids := h.pageTrees.collectAndMarkStaleIdentities(pathInfo); len(ids) > 0 {
					changes = append(changes, ids...)
				}
			} else {
				h.pageTrees.treePagesFromTemplateAdapters.DeleteAllFunc(pathInfo.Base(),
					func(s string, n *pagesfromdata.PagesFromTemplate) bool {
						changes = append(changes, n.DependencyManager)

						// Try to open the file to see if has been deleted.
						f, err := n.GoTmplFi.Meta().Open()
						if err == nil {
							f.Close()
						}
						if err != nil {
							// Remove all pages and resources below.
							prefix := pathInfo.Base() + "/"
							h.pageTrees.treePages.DeletePrefixAll(prefix)
							h.pageTrees.resourceTrees.DeletePrefixAll(prefix)
							changes = append(changes, identity.NewGlobIdentity(prefix+"*"))
						}
						return err != nil
					})
			}

			needsPagesAssemble = true

			if config.RecentlyVisited != nil {
				// Fast render mode. Adding them to the visited queue
				// avoids rerendering them on navigation.
				for _, id := range changes {
					if p, ok := id.(page.Page); ok {
						config.RecentlyVisited.Add(p.RelPermalink())
					}
				}
			}

			h.pageTrees.treeTaxonomyEntries.DeletePrefix("")

			if delete && !isContentDataFile {
				_, ok := h.pageTrees.treePages.LongestPrefixAll(pathInfo.Base())
				if ok {
					h.pageTrees.treePages.DeleteAll(pathInfo.Base())
					h.pageTrees.resourceTrees.DeleteAll(pathInfo.Base())
					if pathInfo.IsBundle() {
						// Assume directory removed.
						h.pageTrees.treePages.DeletePrefixAll(pathInfo.Base() + "/")
						h.pageTrees.resourceTrees.DeletePrefixAll(pathInfo.Base() + "/")
					}
				} else {
					h.pageTrees.resourceTrees.DeleteAll(pathInfo.Base())
				}
			}

			addedOrChangedContent = append(addedOrChangedContent, pathChange{p: pathInfo, delete: delete, isDir: isDir})

		case files.ComponentFolderLayouts:
			tmplChanged = true
			templatePath := pathInfo.Unnormalized().TrimLeadingSlash().PathNoLang()
			if !h.Tmpl().HasTemplate(templatePath) {
				tmplAdded = true
			}

			if tmplAdded {
				logger.Println("Template added", pathInfo.Path())
				// A new template may require a more coarse grained build.
				base := pathInfo.Base()
				if strings.Contains(base, "_markup") {
					// It's hard to determine the exact change set of this,
					// so be very coarse grained.
					changes = append(changes, identity.GenghisKhan)
				}
				if strings.Contains(base, "shortcodes") {
					changes = append(changes, identity.NewGlobIdentity(fmt.Sprintf("shortcodes/%s*", pathInfo.BaseNameNoIdentifier())))
				} else {
					changes = append(changes, pathInfo)
				}
			} else {
				logger.Println("Template changed", pathInfo.Path())
				if templ, found := h.Tmpl().GetIdentity(templatePath); found {
					changes = append(changes, templ)
				} else {
					changes = append(changes, pathInfo)
				}
			}
		case files.ComponentFolderAssets:
			logger.Println("Asset changed", pathInfo.Path())
			changes = append(changes, pathInfo)
		case files.ComponentFolderData:
			logger.Println("Data changed", pathInfo.Path())

			// This should cover all usage of site.Data.
			// Currently very coarse grained.
			changes = append(changes, siteidentities.Data)
			h.init.data.Reset()
		case files.ComponentFolderI18n:
			logger.Println("i18n changed", pathInfo.Path())
			i18nChanged = true
			// It's hard to determine the exact change set of this,
			// so be very coarse grained for now.
			changes = append(changes, identity.GenghisKhan)
		case files.ComponentFolderArchetypes:
			// Ignore for now.
		default:
			panic(fmt.Sprintf("unknown component: %q", pathInfo.Component()))
		}
	}

	changedPaths.deleted = removeDuplicatePaths(changedPaths.deleted)
	changedPaths.changedFiles = removeDuplicatePaths(changedPaths.changedFiles)

	h.Log.Trace(logg.StringFunc(func() string {
		var sb strings.Builder
		sb.WriteString("Resolved paths:\n")
		sb.WriteString("Deleted:\n")
		for _, p := range changedPaths.deleted {
			sb.WriteString("path: " + p.Path())
			sb.WriteString("\n")
		}
		sb.WriteString("Changed:\n")
		for _, p := range changedPaths.changedFiles {
			sb.WriteString("path: " + p.Path())
			sb.WriteString("\n")
		}
		return sb.String()
	}))

	for _, deletedDir := range deletedDirs {
		prefix := deletedDir + "/"
		predicate := func(id identity.Identity) bool {
			// This will effectively reset all pages below this dir.
			return strings.HasPrefix(paths.AddLeadingSlash(id.IdentifierBase()), prefix)
		}
		// Test in both directions.
		changes = append(changes, identity.NewPredicateIdentity(
			// Is dependent.
			predicate,
			// Is dependency.
			predicate,
		),
		)
	}

	if len(addedContentPaths) > 0 {
		// These content files are new and not in use anywhere.
		// To make sure that these gets listed in any site.RegularPages ranges or similar
		// we could invalidate everything, but first try to collect a sample set
		// from the surrounding pages.
		var surroundingIDs []identity.Identity
		for _, p := range addedContentPaths {
			if ids := h.pageTrees.collectIdentitiesSurrounding(p.Base(), 10); len(ids) > 0 {
				surroundingIDs = append(surroundingIDs, ids...)
			}
		}

		if len(surroundingIDs) > 0 {
			changes = append(changes, surroundingIDs...)
		} else {
			// No surrounding pages found, so invalidate everything.
			changes = append(changes, identity.GenghisKhan)
		}
	}

	for _, deleted := range changedPaths.deleted {
		handleChange(deleted, true, false)
	}

	for _, id := range changedPaths.changedFiles {
		handleChange(id, false, false)
	}

	for _, id := range changedPaths.changedDirs {
		handleChange(id, false, true)
	}

	resourceFiles := h.fileEventsContentPaths(addedOrChangedContent)

	changed := &WhatChanged{
		needsPagesAssembly: needsPagesAssemble,
		identitySet:        make(identity.Identities),
	}
	changed.Add(changes...)

	config.WhatChanged = changed

	if err := init(config); err != nil {
		return err
	}

	var cacheBusterOr func(string) bool
	if len(cacheBusters) > 0 {
		cacheBusterOr = func(s string) bool {
			for _, cb := range cacheBusters {
				if cb(s) {
					return true
				}
			}
			return false
		}
	}

	if err := h.resolveAndClearStateForIdentities(ctx, l, cacheBusterOr, changed.Drain()); err != nil {
		return err
	}

	if tmplChanged || i18nChanged {
		// TODO(bep) we should split this, but currently the loading of i18n and layout files are tied together. See #12048.
		h.init.layouts.Reset()
		if err := loggers.TimeTrackfn(func() (logg.LevelLogger, error) {
			// TODO(bep) this could probably be optimized to somehow
			// only load the changed templates and its dependencies, but that is non-trivial.
			ll := l.WithField("substep", "rebuild templates")
			var prototype *deps.Deps
			for i, s := range h.Sites {
				if err := s.Deps.Compile(prototype); err != nil {
					return ll, err
				}
				if i == 0 {
					prototype = s.Deps
				}
			}
			return ll, nil
		}); err != nil {
			return err
		}
	}

	if resourceFiles != nil {
		if err := h.processFiles(ctx, l, config, resourceFiles...); err != nil {
			return err
		}
	}

	if h.isRebuild() {
		if err := h.processContentAdaptersOnRebuild(ctx, config); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) LogServerAddresses() {
	if h.hugoInfo.IsMultihost() {
		for _, s := range h.Sites {
			h.Log.Printf("Web Server is available at %s (bind address %s) %s\n", s.conf.C.BaseURL, s.conf.C.ServerInterface, s.Language().Lang)
		}
	} else {
		s := h.Sites[0]
		h.Log.Printf("Web Server is available at %s (bind address %s)\n", s.conf.C.BaseURL, s.conf.C.ServerInterface)
	}
}

func (h *HugoSites) processFull(ctx context.Context, l logg.LevelLogger, config *BuildCfg) (err error) {
	if err = h.processFiles(ctx, l, config); err != nil {
		err = fmt.Errorf("readAndProcessContent: %w", err)
		return
	}
	return err
}

func (s *Site) handleContentAdapterChanges(bi pagesfromdata.BuildInfo, buildConfig *BuildCfg) {
	if !s.h.isRebuild() {
		return
	}

	if len(bi.ChangedIdentities) > 0 {
		buildConfig.WhatChanged.Add(bi.ChangedIdentities...)
		buildConfig.WhatChanged.needsPagesAssembly = true
	}

	for _, p := range bi.DeletedPaths {
		pp := path.Join(bi.Path.Base(), p)
		if v, ok := s.pageMap.treePages.Delete(pp); ok {
			buildConfig.WhatChanged.Add(v.GetIdentity())
		}
	}
}

func (h *HugoSites) processContentAdaptersOnRebuild(ctx context.Context, buildConfig *BuildCfg) error {
	// Make sure the layouts are initialized.
	if _, err := h.init.layouts.Do(context.Background()); err != nil {
		return err
	}
	g := rungroup.Run[*pagesfromdata.PagesFromTemplate](ctx, rungroup.Config[*pagesfromdata.PagesFromTemplate]{
		NumWorkers: h.numWorkers,
		Handle: func(ctx context.Context, p *pagesfromdata.PagesFromTemplate) error {
			bi, err := p.Execute(ctx)
			if err != nil {
				return err
			}
			s := p.Site.(*Site)
			s.handleContentAdapterChanges(bi, buildConfig)
			return nil
		},
	})

	h.pageTrees.treePagesFromTemplateAdapters.WalkPrefixRaw(doctree.LockTypeRead, "", func(key string, p *pagesfromdata.PagesFromTemplate) (bool, error) {
		if p.StaleVersion() > 0 {
			g.Enqueue(p)
		}
		return false, nil
	})

	return g.Wait()
}

func (s *HugoSites) processFiles(ctx context.Context, l logg.LevelLogger, buildConfig *BuildCfg, filenames ...pathChange) error {
	if s.Deps == nil {
		panic("nil deps on site")
	}

	sourceSpec := source.NewSourceSpec(s.PathSpec, buildConfig.ContentInclusionFilter, s.BaseFs.Content.Fs)

	// For inserts, we can pick an arbitrary pageMap.
	pageMap := s.Sites[0].pageMap

	c := newPagesCollector(ctx, s.h, sourceSpec, s.Log, l, pageMap, buildConfig, filenames)

	if err := c.Collect(); err != nil {
		return err
	}

	return nil
}
