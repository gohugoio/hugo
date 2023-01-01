//go:build darwin && cgo

package filenotify

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"

	qt "github.com/frankban/quicktest"
)

func TestFSEventsAddRemove(t *testing.T) {
	c := qt.New(t)
	w, err := NewFSEventsWatcher()

	c.Assert(err, qt.IsNil)
	c.Assert(w.Add("foo"), qt.Not(qt.IsNil))
	c.Assert(w.Remove("foo"), qt.Not(qt.IsNil))

	f, err := ioutil.TempFile("", "asdf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(f.Name())
	c.Assert(w.Add(f.Name()), qt.IsNil)
	c.Assert(w.Remove(f.Name()), qt.IsNil)
}

func TestFSEventsEvent(t *testing.T) {
	c := qt.New(t)

	method := "fsevents"

	c.Run(fmt.Sprintf("%s, Watch dir", method), func(c *qt.C) {
		dir, w := prepareFSEventsTest(c, method)
		subdir := filepath.Join(dir, subdir1)
		c.Assert(w.Add(subdir), qt.IsNil)

		filename := filepath.Join(subdir, "file1")

		// Write to one file.
		c.Assert(ioutil.WriteFile(filename, []byte("changed"), 0600), qt.IsNil)

		var expected []fsnotify.Event

		expected = append(expected, fsnotify.Event{Name: filename, Op: fsnotify.Write})
		assertEvents(c, w, expected...)

		// Remove one file.
		filename = filepath.Join(subdir, "file2")
		c.Assert(os.Remove(filename), qt.IsNil)
		assertEvents(c, w, fsnotify.Event{Name: filename, Op: fsnotify.Remove})

		// Add one file.
		filename = filepath.Join(subdir, "file3")
		c.Assert(ioutil.WriteFile(filename, []byte("new"), 0600), qt.IsNil)
		assertEvents(c, w, fsnotify.Event{Name: filename, Op: fsnotify.Create})

		// Remove entire directory.
		subdir = filepath.Join(dir, subdir2)
		// Fsevent watcher fails if watched root is deleted
		// so parent dir is added here
		c.Assert(w.Add(dir), qt.IsNil)

		c.Assert(os.RemoveAll(subdir), qt.IsNil)

		expected = expected[:0]

		expected = append(
			expected,
			fsnotify.Event{Name: filepath.Join(subdir, "file2"), Op: fsnotify.Remove},
			fsnotify.Event{Name: filepath.Join(subdir, "file0"), Op: fsnotify.Remove},
			fsnotify.Event{Name: filepath.Join(subdir, "file1"), Op: fsnotify.Remove},
			fsnotify.Event{Name: subdir, Op: fsnotify.Remove},
		)
		assertEvents(c, w, expected...)

	})

	c.Run(fmt.Sprintf("%s, Add should not trigger event", "fsevents"), func(c *qt.C) {
		dir, w := prepareFSEventsTest(c, method)
		subdir := filepath.Join(dir, subdir1)

		w.Add(subdir)
		assertEvents(c, w)

		// Create a new sub directory and add it to the watcher.
		subdir = filepath.Join(dir, subdir1, subdir2)
		c.Assert(os.Mkdir(subdir, 0777), qt.IsNil)
		w.Add(subdir)
		// This should create only one event.
		assertEvents(c, w, fsnotify.Event{Name: subdir, Op: fsnotify.Create})
	})

}

func TestFSEventsClose(t *testing.T) {
	c := qt.New(t)
	w, err := NewFSEventsWatcher()
	c.Assert(err, qt.IsNil)
	f1, err := ioutil.TempFile("", "f1")
	c.Assert(err, qt.IsNil)
	f2, err := ioutil.TempFile("", "f2")
	c.Assert(err, qt.IsNil)
	filename1 := f1.Name()
	filename2 := f2.Name()
	f1.Close()
	f2.Close()

	c.Assert(w.Add(filename1), qt.IsNil)
	c.Assert(w.Add(filename2), qt.IsNil)
	c.Assert(w.Close(), qt.IsNil)
	c.Assert(w.Close(), qt.IsNil)
	c.Assert(ioutil.WriteFile(filename1, []byte("new"), 0600), qt.IsNil)
	c.Assert(ioutil.WriteFile(filename2, []byte("new"), 0600), qt.IsNil)
	// No more event as the watchers are closed.
	assertEvents(c, w)

	f2, err = ioutil.TempFile("", "f2")
	c.Assert(err, qt.IsNil)

	defer os.Remove(f2.Name())

	c.Assert(w.Add(f2.Name()), qt.Not(qt.IsNil))

}

func prepareFSEventsTest(c *qt.C, id string) (string, FileWatcher) {
	w, err := NewFSEventsWatcher()
	c.Assert(err, qt.IsNil)

	dir := prepareTestDirWithSomeFiles(c, id)

	c.Cleanup(func() {
		w.Close()
	})
	err = waitForInit(dir, w)
	c.Assert(err, qt.IsNil)
	drainEvents(c, w)
	return dir, w
}

// FSEvents-wathcer notification can be unstable on starting.
// It contains past events or misses new events unexpectedly.
// Though, this may be no problem for actual uses (local server), it makes tests flaky.
// Thus, we have to wait until detecting a fresh write event.
func waitForInit(dir string, w FileWatcher) (err error) {
	err = w.Add(dir)
	if err != nil {
		return err
	}
	defer func() {
		e := w.Remove(dir)
		if err == nil && e != nil {
			err = e
		}
	}()

	f, err := ioutil.TempFile(dir, "testfile")
	testfile := f.Name()
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	defer func() {
		e := os.Remove(testfile)
		if err == nil && e != nil {
			err = e
		}
	}()

	wait := func() error {
		createdOK := false
		createdSkipped := false
		skipCreatedCheck := time.After(watchWaitTime * 2)
		timeout := time.After(watchWaitTime * 3)
		for i := 0; true; i++ {
			select {
			case evt := <-w.Events():
				if evt.Name != testfile {
					continue
				}
				// first check: CREATE event of testfile
				if evt.Op == fsnotify.Create {
					createdOK = true
					err := ioutil.WriteFile(testfile, []byte(fmt.Sprint(i)), 0600)
					if err != nil {
						return err
					}
				}
				// second check: WRITE event of testfile
				if (createdOK || createdSkipped) && evt.Op == fsnotify.Write {
					return nil
				}
			case e := <-w.Errors():
				return fmt.Errorf("got unexpected error waiting for FSEvents init %v", e)
			case <-skipCreatedCheck:
				// When CREATE event of testfile is missed, skip the check
				createdSkipped = true
				err := ioutil.WriteFile(testfile, []byte(fmt.Sprint(i)), 0600)
				if err != nil {
					return err
				}
			case <-timeout:
				return fmt.Errorf("timeout during waiting for FSEvents init")
			}
		}
		return nil
	}
	return wait()
}
