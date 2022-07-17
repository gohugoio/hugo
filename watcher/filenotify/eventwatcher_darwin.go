// Copyright 2022 The Hugo Authors. All rights reserved.
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
//

//go:build darwin && cgo

package filenotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsevents"
	"github.com/fsnotify/fsnotify"
)

var (
	errFSEventsWatcherClosed              = errors.New("fsEventsWatcher is closed")
	errFSEventsWatcherStreamNotRegistered = errors.New("stream not registered")
)

type eventStream struct {
	*fsevents.EventStream
	isDir    bool
	watcher  *fsEventsWatcher
	basePath string
	removed  chan bool
}

type fsEventsWatcher struct {
	streams map[string]*eventStream
	events  chan fsnotify.Event
	errors  chan error
	mu      sync.Mutex
	done    chan bool
}

func (w *fsEventsWatcher) Events() <-chan fsnotify.Event {
	select {
	case <-w.done:
		return nil
	default:
		return w.events
	}
}

func (w *fsEventsWatcher) Errors() <-chan error {
	return w.errors
}

func (w *fsEventsWatcher) Add(path string) error {
	select {
	case <-w.done:
		return errFSEventsWatcherClosed
	default:
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	w.mu.Lock()
	_, found := w.streams[abs]
	w.mu.Unlock()

	if found {
		return fmt.Errorf("already registered: %s", abs)
	}

	if !w.hasParentEventStreamPath(abs) {
		if err := w.add(abs); err != nil {
			return err
		}
	}

	if childPaths := w.getChildEventStreamPaths(abs); len(childPaths) > 0 {
		if err := w.removePaths(childPaths); err != nil {
			return err
		}
	}
	// https://github.com/fsnotify/fsevents/issues/48
	if len(w.streams) > 4096 {
		return fmt.Errorf("too many fsevent streams: %d\n", len(w.streams))
	}

	return nil
}

func (w *fsEventsWatcher) add(path string) error {
	dev, err := fsevents.DeviceForPath(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Symlinked-path like "/temp" cannot be watched
	evaled, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	isDir := fi.IsDir()
	es := &fsevents.EventStream{
		Paths:   []string{evaled},
		Latency: 10 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot,
	}
	stream := &eventStream{
		es,
		isDir,
		w,
		path,
		make(chan bool),
	}
	w.mu.Lock()
	w.streams[path] = stream
	w.mu.Unlock()
	go func(stream *eventStream) {
		stream.Start()
		stream.Flush(true)
		for {
			select {
			case <-stream.watcher.done:
			case <-stream.removed:
				stream.Flush(true)
				stream.Stop()
				return
			case evs := <-stream.Events:
				for _, evt := range evs {
					err := stream.sendEvent(evt)
					if err != nil {
						return
					}
				}
			}
		}
	}(stream)
	return nil
}

func matchEventFlag(t, m fsevents.EventFlags) bool {
	return t&m == m
}

func (s *eventStream) convertEventPath(path string) (string, error) {
	// Symlinks-evaled path
	path = "/" + path

	evaledBasePath, err := filepath.EvalSymlinks(s.basePath)
	if err != nil {
		return "", err
	}

	rel := path[len(evaledBasePath):]

	return filepath.Join(s.basePath, rel), nil
}

func (s *eventStream) convertEvent(e fsevents.Event) (fsnotify.Event, error) {
	name, err := s.convertEventPath(e.Path)

	if err != nil {
		return fsnotify.Event{}, err
	}

	ne := fsnotify.Event{
		Name: name,
		Op:   0,
	}
	if matchEventFlag(e.Flags, fsevents.ItemCreated) {
		ne.Op = fsnotify.Create
		return ne, nil
	}
	if matchEventFlag(e.Flags, fsevents.ItemRemoved) {
		ne.Op = fsnotify.Remove
		return ne, nil
	}
	if matchEventFlag(e.Flags, fsevents.ItemRenamed) {
		ne.Op = fsnotify.Rename
		return ne, nil
	}
	if matchEventFlag(e.Flags, fsevents.ItemModified) {
		ne.Op = fsnotify.Write
		return ne, nil
	}

	return ne, nil
}

func (s *eventStream) sendEvent(e fsevents.Event) error {
	w := s.watcher
	ne, err := s.convertEvent(e)
	if err != nil {
		return err
	}
	if ne.Op == 0 {
		return nil
	}
	w.events <- ne
	return nil
}

func (w *fsEventsWatcher) sendErr(e error) {
    w.errors <- e
}

func (w *fsEventsWatcher) hasParentEventStreamPath(path string) bool {
	for p, s := range w.streams {
		if s.isDir && strings.HasPrefix(filepath.Dir(path), p) {
			return true
		}
	}
	return false
}

func (w *fsEventsWatcher) getChildEventStreamPaths(path string) (children []string) {
	for p := range w.streams {
		if strings.HasPrefix(filepath.Dir(p), path) {
			children = append(children, p)
		}
	}

	return
}

func (w *fsEventsWatcher) Remove(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return w.remove(abs)
}

func (w *fsEventsWatcher) removePaths(paths []string) error {
	for _, p := range paths {
		if err := w.remove(p); err != nil {
			return err
		}
	}
	return nil
}

func (w *fsEventsWatcher) remove(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	stream, exists := w.streams[path]
	if !exists {
		return errFSEventsWatcherStreamNotRegistered
	}
	close(stream.removed)
	delete(w.streams, path)
	return nil
}

func (w *fsEventsWatcher) Close() error {
	select {
	case <-w.done:
		return nil
	default:
	}

	close(w.done)
	for path := range w.streams {
		err := w.remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewFSEventsWatcher returns a fsevents file watcher
func NewFSEventsWatcher() (FileWatcher, error) {
	w := &fsEventsWatcher{
		streams: make(map[string]*eventStream),
		done:    make(chan bool),
		events:  make(chan fsnotify.Event),
		errors:  make(chan error),
	}
	return w, nil
}

// NewEventWatcher returns an FSEvents based file watcher on darwin
func NewEventWatcher() (FileWatcher, error) {
	return NewFSEventsWatcher()
}
