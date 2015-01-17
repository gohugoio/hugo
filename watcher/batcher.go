package watcher

import (
	"time"

	"gopkg.in/fsnotify.v0"
)

type Batcher struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}

	Event chan []*fsnotify.FileEvent // Events are returned on this channel
}

func New(interval time.Duration) (*Batcher, error) {
	watcher, err := fsnotify.NewWatcher()

	batcher := &Batcher{}
	batcher.Watcher = watcher
	batcher.interval = interval
	batcher.done = make(chan struct{}, 1)
	batcher.Event = make(chan []*fsnotify.FileEvent, 1)

	if err == nil {
		go batcher.run()
	}

	return batcher, err
}

func (b *Batcher) run() {
	tick := time.Tick(b.interval)
	evs := make([]*fsnotify.FileEvent, 0)
OuterLoop:
	for {
		select {
		case ev := <-b.Watcher.Event:
			evs = append(evs, ev)
		case <-tick:
			if len(evs) == 0 {
				continue
			}
			b.Event <- evs
			evs = make([]*fsnotify.FileEvent, 0)
		case <-b.done:
			break OuterLoop
		}
	}
	close(b.done)
}

func (b *Batcher) Close() {
	b.done <- struct{}{}
	b.Watcher.Close()
}
