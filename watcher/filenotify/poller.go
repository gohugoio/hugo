package filenotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
)

var (
	// errPollerClosed is returned when the poller is closed
	errPollerClosed = errors.New("poller is closed")
	// errNoSuchWatch is returned when trying to remove a watch that doesn't exist
	errNoSuchWatch = errors.New("watch does not exist")
)

// filePoller is used to poll files for changes, especially in cases where fsnotify
// can't be run (e.g. when inotify handles are exhausted)
// filePoller satisfies the FileWatcher interface
type filePoller struct {
	// the duration between polls.
	interval time.Duration
	// watches is the list of files currently being polled, close the associated channel to stop the watch
	watches map[string]struct{}
	// Will be closed when done.
	done chan struct{}
	// events is the channel to listen to for watch events
	events chan fsnotify.Event
	// errors is the channel to listen to for watch errors
	errors chan error
	// mu locks the poller for modification
	mu sync.Mutex
	// closed is used to specify when the poller has already closed
	closed bool
}

// Add adds a filename to the list of watches
// once added the file is polled for changes in a separate goroutine
func (w *filePoller) Add(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return errPollerClosed
	}

	item, err := newItemToWatch(name)
	if err != nil {
		return err
	}
	if item.left.FileInfo == nil {
		return os.ErrNotExist
	}

	if w.watches == nil {
		w.watches = make(map[string]struct{})
	}
	if _, exists := w.watches[name]; exists {
		return fmt.Errorf("watch exists")
	}
	w.watches[name] = struct{}{}

	go w.watch(item)
	return nil
}

// Remove stops and removes watch with the specified name
func (w *filePoller) Remove(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.remove(name)
}

func (w *filePoller) remove(name string) error {
	if w.closed {
		return errPollerClosed
	}

	_, exists := w.watches[name]
	if !exists {
		return errNoSuchWatch
	}
	delete(w.watches, name)
	return nil
}

// Events returns the event channel
// This is used for notifications on events about watched files
func (w *filePoller) Events() <-chan fsnotify.Event {
	return w.events
}

// Errors returns the errors channel
// This is used for notifications about errors on watched files
func (w *filePoller) Errors() <-chan error {
	return w.errors
}

// Close closes the poller
// All watches are stopped, removed, and the poller cannot be added to
func (w *filePoller) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}
	w.closed = true
	close(w.done)
	for name := range w.watches {
		w.remove(name)
	}

	return nil
}

// sendEvent publishes the specified event to the events channel
func (w *filePoller) sendEvent(e fsnotify.Event) error {
	select {
	case w.events <- e:
	case <-w.done:
		return fmt.Errorf("closed")
	}
	return nil
}

// sendErr publishes the specified error to the errors channel
func (w *filePoller) sendErr(e error) error {
	select {
	case w.errors <- e:
	case <-w.done:
		return fmt.Errorf("closed")
	}
	return nil
}

// watch watches item for changes until done is closed.
func (w *filePoller) watch(item *itemToWatch) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-w.done:
			return
		}

		evs, err := item.checkForChanges()
		if err != nil {
			if err := w.sendErr(err); err != nil {
				return
			}
		}

		item.left, item.right = item.right, item.left

		for _, ev := range evs {
			if err := w.sendEvent(ev); err != nil {
				return
			}
		}

	}
}

// recording records the state of a file or a dir.
type recording struct {
	os.FileInfo

	// Set if FileInfo is a dir.
	entries map[string]os.FileInfo
}

func (r *recording) clear() {
	r.FileInfo = nil
	if r.entries != nil {
		for k := range r.entries {
			delete(r.entries, k)
		}
	}
}

func (r *recording) record(filename string) error {
	r.clear()

	fi, err := os.Stat(filename)
	if err != nil && !herrors.IsNotExist(err) {
		return err
	}

	if fi == nil {
		return nil
	}

	r.FileInfo = fi

	// If fi is a dir, we watch the files inside that directory (not recursively).
	// This matches the behavior of fsnotity.
	if fi.IsDir() {
		f, err := os.Open(filename)
		if err != nil {
			if herrors.IsNotExist(err) {
				return nil
			}
			return err
		}
		defer f.Close()

		fis, err := f.Readdir(-1)
		if err != nil {
			if herrors.IsNotExist(err) {
				return nil
			}
			return err
		}

		for _, fi := range fis {
			r.entries[fi.Name()] = fi
		}
	}

	return nil
}

// itemToWatch may be a file or a dir.
type itemToWatch struct {
	// Full path to the filename.
	filename string

	// Snapshots of the stat state of this file or dir.
	left  *recording
	right *recording
}

func newItemToWatch(filename string) (*itemToWatch, error) {
	r := &recording{
		entries: make(map[string]os.FileInfo),
	}
	err := r.record(filename)
	if err != nil {
		return nil, err
	}

	return &itemToWatch{filename: filename, left: r}, nil
}

func (item *itemToWatch) checkForChanges() ([]fsnotify.Event, error) {
	if item.right == nil {
		item.right = &recording{
			entries: make(map[string]os.FileInfo),
		}
	}

	err := item.right.record(item.filename)
	if err != nil && !herrors.IsNotExist(err) {
		return nil, err
	}

	dirOp := checkChange(item.left.FileInfo, item.right.FileInfo)

	if dirOp != 0 {
		evs := []fsnotify.Event{{Op: dirOp, Name: item.filename}}
		return evs, nil
	}

	if item.left.FileInfo == nil || !item.left.IsDir() {
		// Done.
		return nil, nil
	}

	leftIsIn := false
	left, right := item.left.entries, item.right.entries
	if len(right) > len(left) {
		left, right = right, left
		leftIsIn = true
	}

	var evs []fsnotify.Event

	for name, fi1 := range left {
		fi2 := right[name]
		fil, fir := fi1, fi2
		if leftIsIn {
			fil, fir = fir, fil
		}
		op := checkChange(fil, fir)
		if op != 0 {
			evs = append(evs, fsnotify.Event{Op: op, Name: filepath.Join(item.filename, name)})
		}

	}

	return evs, nil
}

func checkChange(fi1, fi2 os.FileInfo) fsnotify.Op {
	if fi1 == nil && fi2 != nil {
		return fsnotify.Create
	}
	if fi1 != nil && fi2 == nil {
		return fsnotify.Remove
	}
	if fi1 == nil && fi2 == nil {
		return 0
	}
	if fi1.IsDir() || fi2.IsDir() {
		return 0
	}
	if fi1.Mode() != fi2.Mode() {
		return fsnotify.Chmod
	}
	if fi1.ModTime() != fi2.ModTime() || fi1.Size() != fi2.Size() {
		return fsnotify.Write
	}

	return 0
}
