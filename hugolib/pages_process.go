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
	"context"
	"fmt"
	"path/filepath"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/hugofs"
)

func newPagesProcessor(h *HugoSites, sp *source.SourceSpec) *pagesProcessor {
	procs := make(map[string]pagesCollectorProcessorProvider)
	for _, s := range h.Sites {
		procs[s.Lang()] = &sitePagesProcessor{
			m:           s.pageMap,
			errorSender: s.h,
			itemChan:    make(chan interface{}, config.GetNumWorkerMultiplier()*2),
		}
	}
	return &pagesProcessor{
		procs: procs,
	}
}

type pagesCollectorProcessorProvider interface {
	Process(item interface{}) error
	Start(ctx context.Context) context.Context
	Wait() error
}

type pagesProcessor struct {
	// Per language/Site
	procs map[string]pagesCollectorProcessorProvider
}

func (proc *pagesProcessor) Process(item interface{}) error {
	switch v := item.(type) {
	// Page bundles mapped to their language.
	case pageBundles:
		for _, vv := range v {
			proc.getProcFromFi(vv.header).Process(vv)
		}
	case hugofs.FileMetaInfo:
		proc.getProcFromFi(v).Process(v)
	default:
		panic(fmt.Sprintf("unrecognized item type in Process: %T", item))

	}

	return nil
}

func (proc *pagesProcessor) Start(ctx context.Context) context.Context {
	for _, p := range proc.procs {
		ctx = p.Start(ctx)
	}
	return ctx
}

func (proc *pagesProcessor) Wait() error {
	var err error
	for _, p := range proc.procs {
		if e := p.Wait(); e != nil {
			err = e
		}
	}
	return err
}

func (proc *pagesProcessor) getProcFromFi(fi hugofs.FileMetaInfo) pagesCollectorProcessorProvider {
	if p, found := proc.procs[fi.Meta().Lang()]; found {
		return p
	}
	return defaultPageProcessor
}

type nopPageProcessor int

func (nopPageProcessor) Process(item interface{}) error {
	return nil
}

func (nopPageProcessor) Start(ctx context.Context) context.Context {
	return context.Background()
}

func (nopPageProcessor) Wait() error {
	return nil
}

var defaultPageProcessor = new(nopPageProcessor)

type sitePagesProcessor struct {
	m           *pageMap
	errorSender herrors.ErrorSender

	itemChan  chan interface{}
	itemGroup *errgroup.Group
}

func (p *sitePagesProcessor) Process(item interface{}) error {
	p.itemChan <- item
	return nil
}

func (p *sitePagesProcessor) Start(ctx context.Context) context.Context {
	p.itemGroup, ctx = errgroup.WithContext(ctx)
	p.itemGroup.Go(func() error {
		for item := range p.itemChan {
			if err := p.doProcess(item); err != nil {
				return err
			}
		}
		return nil
	})
	return ctx
}

func (p *sitePagesProcessor) Wait() error {
	close(p.itemChan)
	return p.itemGroup.Wait()
}

func (p *sitePagesProcessor) copyFile(fim hugofs.FileMetaInfo) error {
	meta := fim.Meta()
	f, err := meta.Open()
	if err != nil {
		return errors.Wrap(err, "copyFile: failed to open")
	}

	s := p.m.s

	target := filepath.Join(s.PathSpec.GetTargetLanguageBasePath(), meta.Path())

	defer f.Close()

	return s.publish(&s.PathSpec.ProcessingStats.Files, target, f)

}

func (p *sitePagesProcessor) doProcess(item interface{}) error {
	m := p.m
	switch v := item.(type) {
	case *fileinfoBundle:
		if err := m.AddFilesBundle(v.header, v.resources...); err != nil {
			return err
		}
	case hugofs.FileMetaInfo:
		if p.shouldSkip(v) {
			return nil
		}
		meta := v.Meta()

		classifier := meta.Classifier()
		switch classifier {
		case files.ContentClassContent:
			if err := m.AddFilesBundle(v); err != nil {
				return err
			}
		case files.ContentClassFile:
			if err := p.copyFile(v); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("invalid classifier: %q", classifier))
		}
	default:
		panic(fmt.Sprintf("unrecognized item type in Process: %T", item))
	}
	return nil

}

func (p *sitePagesProcessor) shouldSkip(fim hugofs.FileMetaInfo) bool {
	// TODO(ep) unify
	return p.m.s.SourceSpec.DisabledLanguages[fim.Meta().Lang()]
}
