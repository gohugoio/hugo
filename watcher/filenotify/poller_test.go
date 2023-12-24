package filenotify

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/htesting"
)

const (
	subdir1       = "subdir1"
	subdir2       = "subdir2"
	watchWaitTime = 200 * time.Millisecond
)

var (
	isMacOs = runtime.GOOS == "darwin"
	isCI    = htesting.IsCI()
)

func TestPollerAddRemove(t *testing.T) {
	c := qt.New(t)
	w := NewPollingWatcher(watchWaitTime)

	c.Assert(w.Add("foo"), qt.Not(qt.IsNil))
	c.Assert(w.Remove("foo"), qt.Not(qt.IsNil))

	f, err := os.CreateTemp("", "asdf")
	if err != nil {
		t.Fatal(err)
	}
	c.Cleanup(func() {
		c.Assert(w.Close(), qt.IsNil)
		os.Remove(f.Name())
	})
	c.Assert(w.Add(f.Name()), qt.IsNil)
	c.Assert(w.Remove(f.Name()), qt.IsNil)
}

func TestPollerEvent(t *testing.T) {
	t.Skip("flaky test") // TODO(bep)
	c := qt.New(t)

	for _, poll := range []bool{true, false} {
		if !(poll || isMacOs) || isCI {
			// Only run the fsnotify tests on MacOS locally.
			continue
		}
		method := "fsnotify"
		if poll {
			method = "poll"
		}

		c.Run(fmt.Sprintf("%s, Watch dir", method), func(c *qt.C) {
			dir, w := preparePollTest(c, poll)
			subdir := filepath.Join(dir, subdir1)
			c.Assert(w.Add(subdir), qt.IsNil)

			filename := filepath.Join(subdir, "file1")

			// Write to one file.
			c.Assert(os.WriteFile(filename, []byte("changed"), 0o600), qt.IsNil)

			var expected []fsnotify.Event

			if poll {
				expected = append(expected, fsnotify.Event{Name: filename, Op: fsnotify.Write})
				assertEvents(c, w, expected...)
			} else {
				// fsnotify sometimes emits Chmod before Write,
				// which is hard to test, so skip it here.
				drainEvents(c, w)
			}

			// Remove one file.
			filename = filepath.Join(subdir, "file2")
			c.Assert(os.Remove(filename), qt.IsNil)
			assertEvents(c, w, fsnotify.Event{Name: filename, Op: fsnotify.Remove})

			// Add one file.
			filename = filepath.Join(subdir, "file3")
			c.Assert(os.WriteFile(filename, []byte("new"), 0o600), qt.IsNil)
			assertEvents(c, w, fsnotify.Event{Name: filename, Op: fsnotify.Create})

			// Remove entire directory.
			subdir = filepath.Join(dir, subdir2)
			c.Assert(w.Add(subdir), qt.IsNil)

			c.Assert(os.RemoveAll(subdir), qt.IsNil)

			expected = expected[:0]

			// This looks like a bug in fsnotify on MacOS. There are
			// 3 files in this directory, yet we get Remove events
			// for one of them + the directory.
			if !poll {
				expected = append(expected, fsnotify.Event{Name: filepath.Join(subdir, "file2"), Op: fsnotify.Remove})
			}
			expected = append(expected, fsnotify.Event{Name: subdir, Op: fsnotify.Remove})
			assertEvents(c, w, expected...)
		})

		c.Run(fmt.Sprintf("%s, Add should not trigger event", method), func(c *qt.C) {
			dir, w := preparePollTest(c, poll)
			subdir := filepath.Join(dir, subdir1)
			w.Add(subdir)
			assertEvents(c, w)
			// Create a new sub directory and add it to the watcher.
			subdir = filepath.Join(dir, subdir1, subdir2)
			c.Assert(os.Mkdir(subdir, 0o777), qt.IsNil)
			w.Add(subdir)
			// This should create only one event.
			assertEvents(c, w, fsnotify.Event{Name: subdir, Op: fsnotify.Create})
		})

	}
}

func TestPollerClose(t *testing.T) {
	c := qt.New(t)
	w := NewPollingWatcher(watchWaitTime)
	f1, err := os.CreateTemp("", "f1")
	c.Assert(err, qt.IsNil)
	defer os.Remove(f1.Name())
	f2, err := os.CreateTemp("", "f2")
	c.Assert(err, qt.IsNil)
	filename1 := f1.Name()
	filename2 := f2.Name()
	f1.Close()
	f2.Close()

	c.Assert(w.Add(filename1), qt.IsNil)
	c.Assert(w.Add(filename2), qt.IsNil)
	c.Assert(w.Close(), qt.IsNil)
	c.Assert(w.Close(), qt.IsNil)
	c.Assert(os.WriteFile(filename1, []byte("new"), 0o600), qt.IsNil)
	c.Assert(os.WriteFile(filename2, []byte("new"), 0o600), qt.IsNil)
	// No more event as the watchers are closed.
	assertEvents(c, w)

	f2, err = os.CreateTemp("", "f2")
	c.Assert(err, qt.IsNil)

	defer os.Remove(f2.Name())

	c.Assert(w.Add(f2.Name()), qt.Not(qt.IsNil))
}

func TestCheckChange(t *testing.T) {
	c := qt.New(t)

	dir := prepareTestDirWithSomeFiles(c, "check-change")

	stat := func(s ...string) os.FileInfo {
		fi, err := os.Stat(filepath.Join(append([]string{dir}, s...)...))
		c.Assert(err, qt.IsNil)
		return fi
	}

	f0, f1, f2 := stat(subdir2, "file0"), stat(subdir2, "file1"), stat(subdir2, "file2")
	d1 := stat(subdir1)

	// Note that on Windows, only the 0200 bit (owner writable) of mode is used.
	c.Assert(os.Chmod(filepath.Join(filepath.Join(dir, subdir2, "file1")), 0o400), qt.IsNil)
	f1_2 := stat(subdir2, "file1")

	c.Assert(os.WriteFile(filepath.Join(filepath.Join(dir, subdir2, "file2")), []byte("changed"), 0o600), qt.IsNil)
	f2_2 := stat(subdir2, "file2")

	c.Assert(checkChange(f0, nil), qt.Equals, fsnotify.Remove)
	c.Assert(checkChange(nil, f0), qt.Equals, fsnotify.Create)
	c.Assert(checkChange(f1, f1_2), qt.Equals, fsnotify.Chmod)
	c.Assert(checkChange(f2, f2_2), qt.Equals, fsnotify.Write)
	c.Assert(checkChange(nil, nil), qt.Equals, fsnotify.Op(0))
	c.Assert(checkChange(d1, f1), qt.Equals, fsnotify.Op(0))
	c.Assert(checkChange(f1, d1), qt.Equals, fsnotify.Op(0))
}

func BenchmarkPoller(b *testing.B) {
	runBench := func(b *testing.B, item *itemToWatch) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			evs, err := item.checkForChanges()
			if err != nil {
				b.Fatal(err)
			}
			if len(evs) != 0 {
				b.Fatal("got events")
			}

		}
	}

	b.Run("Check for changes in dir", func(b *testing.B) {
		c := qt.New(b)
		dir := prepareTestDirWithSomeFiles(c, "bench-check")
		item, err := newItemToWatch(dir)
		c.Assert(err, qt.IsNil)
		runBench(b, item)
	})

	b.Run("Check for changes in file", func(b *testing.B) {
		c := qt.New(b)
		dir := prepareTestDirWithSomeFiles(c, "bench-check-file")
		filename := filepath.Join(dir, subdir1, "file1")
		item, err := newItemToWatch(filename)
		c.Assert(err, qt.IsNil)
		runBench(b, item)
	})
}

func prepareTestDirWithSomeFiles(c *qt.C, id string) string {
	dir := c.TB.TempDir()
	c.Assert(os.MkdirAll(filepath.Join(dir, subdir1), 0o777), qt.IsNil)
	c.Assert(os.MkdirAll(filepath.Join(dir, subdir2), 0o777), qt.IsNil)

	for i := 0; i < 3; i++ {
		c.Assert(os.WriteFile(filepath.Join(dir, subdir1, fmt.Sprintf("file%d", i)), []byte("hello1"), 0o600), qt.IsNil)
	}

	for i := 0; i < 3; i++ {
		c.Assert(os.WriteFile(filepath.Join(dir, subdir2, fmt.Sprintf("file%d", i)), []byte("hello2"), 0o600), qt.IsNil)
	}

	c.Cleanup(func() {
		os.RemoveAll(dir)
	})

	return dir
}

func preparePollTest(c *qt.C, poll bool) (string, FileWatcher) {
	var w FileWatcher
	if poll {
		w = NewPollingWatcher(watchWaitTime)
	} else {
		var err error
		w, err = NewEventWatcher()
		c.Assert(err, qt.IsNil)
	}

	dir := prepareTestDirWithSomeFiles(c, fmt.Sprint(poll))

	c.Cleanup(func() {
		w.Close()
	})
	return dir, w
}

func assertEvents(c *qt.C, w FileWatcher, evs ...fsnotify.Event) {
	c.Helper()
	i := 0
	check := func() error {
		for {
			select {
			case got := <-w.Events():
				if i > len(evs)-1 {
					return fmt.Errorf("got too many event(s): %q", got)
				}
				expected := evs[i]
				i++
				if expected.Name != got.Name {
					return fmt.Errorf("got wrong filename, expected %q: %v", expected.Name, got.Name)
				} else if got.Op&expected.Op != expected.Op {
					return fmt.Errorf("got wrong event type, expected %q: %v", expected.Op, got.Op)
				}
			case e := <-w.Errors():
				return fmt.Errorf("got unexpected error waiting for events %v", e)
			case <-time.After(watchWaitTime + (watchWaitTime / 2)):
				return nil
			}
		}
	}
	c.Assert(check(), qt.IsNil)
	c.Assert(i, qt.Equals, len(evs))
}

func drainEvents(c *qt.C, w FileWatcher) {
	c.Helper()
	check := func() error {
		for {
			select {
			case <-w.Events():
			case e := <-w.Errors():
				return fmt.Errorf("got unexpected error waiting for events %v", e)
			case <-time.After(watchWaitTime * 2):
				return nil
			}
		}
	}
	c.Assert(check(), qt.IsNil)
}
