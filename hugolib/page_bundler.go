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

package hugolib

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	// Use this until errgroup gets ported to context
	// See https://github.com/golang/go/issues/19781
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type siteContentProcessor struct {
	baseDir string

	site *Site

	handleContent contentHandler

	// The input file bundles.
	fileBundlesChan chan *bundleDir

	// The input file singles.
	fileSinglesChan chan *fileInfo

	// These assets should be just copied to destination.
	fileAssetsChan chan []string

	numWorkers int

	// The output Pages
	pagesChan chan *Page

	// Used for partial rebuilds (aka. live reload)
	// Will signal replacement of pages in the site collection.
	partialBuild bool
}

func newSiteContentProcessor(baseDir string, partialBuild bool, s *Site) *siteContentProcessor {
	numWorkers := 12
	if n := runtime.NumCPU() * 3; n > numWorkers {
		numWorkers = n
	}

	numWorkers = int(math.Ceil(float64(numWorkers) / float64(len(s.owner.Sites))))

	return &siteContentProcessor{
		partialBuild:    partialBuild,
		baseDir:         baseDir,
		site:            s,
		handleContent:   newHandlerChain(s),
		fileBundlesChan: make(chan *bundleDir, numWorkers),
		fileSinglesChan: make(chan *fileInfo, numWorkers),
		fileAssetsChan:  make(chan []string, numWorkers),
		numWorkers:      numWorkers,
		pagesChan:       make(chan *Page, numWorkers),
	}
}

func (s *siteContentProcessor) closeInput() {
	close(s.fileSinglesChan)
	close(s.fileBundlesChan)
	close(s.fileAssetsChan)
}

func (s *siteContentProcessor) process(ctx context.Context) error {
	g1, ctx := errgroup.WithContext(ctx)
	g2, _ := errgroup.WithContext(ctx)

	// There can be only one of these per site.
	g1.Go(func() error {
		for p := range s.pagesChan {
			if p.s != s.site {
				panic(fmt.Sprintf("invalid page site: %v vs %v", p.s, s))
			}

			if s.partialBuild {
				s.site.replacePage(p)
			} else {
				s.site.addPage(p)
			}
		}
		return nil
	})

	for i := 0; i < s.numWorkers; i++ {
		g2.Go(func() error {
			for {
				select {
				case f, ok := <-s.fileSinglesChan:
					if !ok {
						return nil
					}
					err := s.readAndConvertContentFile(f)
					if err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})

		g2.Go(func() error {
			for {
				select {
				case filenames, ok := <-s.fileAssetsChan:
					if !ok {
						return nil
					}
					for _, filename := range filenames {
						name := strings.TrimPrefix(filename, s.baseDir)
						f, err := s.site.Fs.Source.Open(filename)
						if err != nil {
							return err
						}

						err = s.site.publish(&s.site.PathSpec.ProcessingStats.Files, name, f)
						f.Close()
						if err != nil {
							return err
						}
					}

				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})

		g2.Go(func() error {
			for {
				select {
				case bundle, ok := <-s.fileBundlesChan:
					if !ok {
						return nil
					}
					err := s.readAndConvertContentBundle(bundle)
					if err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	if err := g2.Wait(); err != nil {
		return err
	}

	close(s.pagesChan)

	if err := g1.Wait(); err != nil {
		return err
	}

	s.site.rawAllPages.Sort()

	return nil

}

func (s *siteContentProcessor) readAndConvertContentFile(file *fileInfo) error {
	ctx := &handlerContext{source: file, baseDir: s.baseDir, pages: s.pagesChan}
	return s.handleContent(ctx).err
}

func (s *siteContentProcessor) readAndConvertContentBundle(bundle *bundleDir) error {
	ctx := &handlerContext{bundle: bundle, baseDir: s.baseDir, pages: s.pagesChan}
	return s.handleContent(ctx).err
}
