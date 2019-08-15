// Copyright 2017 The Hugo Authors. All rights reserved.
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

package commands

import (
	"os"
	"path/filepath"

	"github.com/gohugoio/hugo/hugolib/filesystems"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/fsync"
)

type staticSyncer struct {
	c *commandeer
}

func newStaticSyncer(c *commandeer) (*staticSyncer, error) {
	return &staticSyncer{c: c}, nil
}

func (s *staticSyncer) isStatic(filename string) bool {
	return s.c.hugo().BaseFs.SourceFilesystems.IsStatic(filename)
}

func (s *staticSyncer) syncsStaticEvents(staticEvents []fsnotify.Event) error {
	c := s.c

	syncFn := func(sourceFs *filesystems.SourceFilesystem) (uint64, error) {
		publishDir := c.hugo().PathSpec.PublishDir
		// If root, remove the second '/'
		if publishDir == "//" {
			publishDir = helpers.FilePathSeparator
		}

		if sourceFs.PublishFolder != "" {
			publishDir = filepath.Join(publishDir, sourceFs.PublishFolder)
		}

		syncer := fsync.NewSyncer()
		syncer.NoTimes = c.Cfg.GetBool("noTimes")
		syncer.NoChmod = c.Cfg.GetBool("noChmod")
		syncer.ChmodFilter = chmodFilter
		syncer.SrcFs = sourceFs.Fs
		syncer.DestFs = c.Fs.Destination

		// prevent spamming the log on changes
		logger := helpers.NewDistinctFeedbackLogger()

		for _, ev := range staticEvents {
			// Due to our approach of layering both directories and the content's rendered output
			// into one we can't accurately remove a file not in one of the source directories.
			// If a file is in the local static dir and also in the theme static dir and we remove
			// it from one of those locations we expect it to still exist in the destination
			//
			// If Hugo generates a file (from the content dir) over a static file
			// the content generated file should take precedence.
			//
			// Because we are now watching and handling individual events it is possible that a static
			// event that occupies the same path as a content generated file will take precedence
			// until a regeneration of the content takes places.
			//
			// Hugo assumes that these cases are very rare and will permit this bad behavior
			// The alternative is to track every single file and which pipeline rendered it
			// and then to handle conflict resolution on every event.

			fromPath := ev.Name

			relPath := sourceFs.MakePathRelative(fromPath)

			if relPath == "" {
				// Not member of this virtual host.
				continue
			}

			// Remove || rename is harder and will require an assumption.
			// Hugo takes the following approach:
			// If the static file exists in any of the static source directories after this event
			// Hugo will re-sync it.
			// If it does not exist in all of the static directories Hugo will remove it.
			//
			// This assumes that Hugo has not generated content on top of a static file and then removed
			// the source of that static file. In this case Hugo will incorrectly remove that file
			// from the published directory.
			if ev.Op&fsnotify.Rename == fsnotify.Rename || ev.Op&fsnotify.Remove == fsnotify.Remove {
				if _, err := sourceFs.Fs.Stat(relPath); os.IsNotExist(err) {
					// If file doesn't exist in any static dir, remove it
					toRemove := filepath.Join(publishDir, relPath)

					logger.Println("File no longer exists in static dir, removing", toRemove)
					_ = c.Fs.Destination.RemoveAll(toRemove)
				} else if err == nil {
					// If file still exists, sync it
					logger.Println("Syncing", relPath, "to", publishDir)

					if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
						c.logger.ERROR.Println(err)
					}
				} else {
					c.logger.ERROR.Println(err)
				}

				continue
			}

			// For all other event operations Hugo will sync static.
			logger.Println("Syncing", relPath, "to", publishDir)
			if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
				c.logger.ERROR.Println(err)
			}
		}

		return 0, nil
	}

	_, err := c.doWithPublishDirs(syncFn)
	return err

}
