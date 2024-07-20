package hugolib

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/bep/logg"

	qt "github.com/frankban/quicktest"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/tools/txtar"
)

type TestOpt func(*IntegrationTestConfig)

// TestOptRunning will enable running in integration tests.
func TestOptRunning() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.Running = true
	}
}

// TestOptWatching will enable watching in integration tests.
func TestOptWatching() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.Watching = true
	}
}

// Enable tracing in integration tests.
// THis should only be used during development and not committed to the repo.
func TestOptTrace() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.LogLevel = logg.LevelTrace
	}
}

// TestOptDebug will enable debug logging in integration tests.
func TestOptDebug() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.LogLevel = logg.LevelDebug
	}
}

// TestOptWarn will enable warn logging in integration tests.
func TestOptWarn() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.LogLevel = logg.LevelWarn
	}
}

// TestOptWithNFDOnDarwin will normalize the Unicode filenames to NFD on Darwin.
func TestOptWithNFDOnDarwin() TestOpt {
	return func(c *IntegrationTestConfig) {
		c.NFDFormOnDarwin = true
	}
}

// TestOptWithWorkingDir allows setting any config optiona as a function al option.
func TestOptWithConfig(fn func(c *IntegrationTestConfig)) TestOpt {
	return func(c *IntegrationTestConfig) {
		fn(c)
	}
}

// Test is a convenience method to create a new IntegrationTestBuilder from some files and run a build.
func Test(t testing.TB, files string, opts ...TestOpt) *IntegrationTestBuilder {
	cfg := IntegrationTestConfig{T: t, TxtarString: files}
	for _, o := range opts {
		o(&cfg)
	}
	return NewIntegrationTestBuilder(cfg).Build()
}

// TestE is the same as Test, but returns an error instead of failing the test.
func TestE(t testing.TB, files string, opts ...TestOpt) (*IntegrationTestBuilder, error) {
	cfg := IntegrationTestConfig{T: t, TxtarString: files}
	for _, o := range opts {
		o(&cfg)
	}
	return NewIntegrationTestBuilder(cfg).BuildE()
}

// TestRunning is a convenience method to create a new IntegrationTestBuilder from some files with Running set to true and run a build.
// Deprecated: Use Test with TestOptRunning instead.
func TestRunning(t testing.TB, files string, opts ...TestOpt) *IntegrationTestBuilder {
	cfg := IntegrationTestConfig{T: t, TxtarString: files, Running: true}
	for _, o := range opts {
		o(&cfg)
	}
	return NewIntegrationTestBuilder(cfg).Build()
}

// In most cases you should not use this function directly, but the Test or TestRunning function.
func NewIntegrationTestBuilder(conf IntegrationTestConfig) *IntegrationTestBuilder {
	// Code fences.
	conf.TxtarString = strings.ReplaceAll(conf.TxtarString, "§§§", "```")
	// Multiline strings.
	conf.TxtarString = strings.ReplaceAll(conf.TxtarString, "§§", "`")

	data := txtar.Parse([]byte(conf.TxtarString))

	if conf.NFDFormOnDarwin {
		for i, f := range data.Files {
			data.Files[i].Name = norm.NFD.String(f.Name)
		}
	}

	c, ok := conf.T.(*qt.C)
	if !ok {
		c = qt.New(conf.T)
	}

	if conf.NeedsOsFS {
		if !filepath.IsAbs(conf.WorkingDir) {
			tempDir, clean, err := htesting.CreateTempDir(hugofs.Os, "hugo-integration-test")
			c.Assert(err, qt.IsNil)
			conf.WorkingDir = filepath.Join(tempDir, conf.WorkingDir)
			if !conf.PrintAndKeepTempDir {
				c.Cleanup(clean)
			} else {
				fmt.Println("\nUsing WorkingDir dir:", conf.WorkingDir)
			}
		}
	} else if conf.WorkingDir == "" {
		conf.WorkingDir = helpers.FilePathSeparator
	}

	return &IntegrationTestBuilder{
		Cfg:  conf,
		C:    c,
		data: data,
	}
}

// IntegrationTestBuilder is a (partial) rewrite of sitesBuilder.
// The main problem with the "old" one was that it was that the test data was often a little hidden,
// so it became hard to look at a test and determine what it should do, especially coming back to the
// test after a year or so.
type IntegrationTestBuilder struct {
	*qt.C

	data *txtar.Archive

	fs *hugofs.Fs
	H  *HugoSites

	Cfg IntegrationTestConfig

	changedFiles []string
	createdFiles []string
	removedFiles []string
	renamedFiles []string
	renamedDirs  []string

	buildCount   int
	GCCount      int
	counters     *buildCounters
	logBuff      lockingBuffer
	lastBuildLog string

	builderInit sync.Once
}

type lockingBuffer struct {
	sync.Mutex
	bytes.Buffer
}

func (b *lockingBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.Lock()
	n, err = b.Buffer.ReadFrom(r)
	b.Unlock()
	return
}

func (b *lockingBuffer) Write(p []byte) (n int, err error) {
	b.Lock()
	n, err = b.Buffer.Write(p)
	b.Unlock()
	return
}

// AssertLogContains asserts that the last build log contains the given strings.
// Each string can be negated with a "! " prefix.
func (s *IntegrationTestBuilder) AssertLogContains(els ...string) {
	s.Helper()
	for _, el := range els {
		var negate bool
		el, negate = s.negate(el)
		check := qt.Contains
		if negate {
			check = qt.Not(qt.Contains)
		}
		s.Assert(s.lastBuildLog, check, el)
	}
}

// AssertLogNotContains asserts that the last build log does matches the given regular expressions.
// The regular expressions can be negated with a "! " prefix.
func (s *IntegrationTestBuilder) AssertLogMatches(expression string) {
	s.Helper()
	var negate bool
	expression, negate = s.negate(expression)
	re := regexp.MustCompile(expression)
	checker := qt.IsTrue
	if negate {
		checker = qt.IsFalse
	}

	s.Assert(re.MatchString(s.lastBuildLog), checker, qt.Commentf(s.lastBuildLog))
}

func (s *IntegrationTestBuilder) AssertBuildCountData(count int) {
	s.Helper()
	s.Assert(s.H.init.data.InitCount(), qt.Equals, count)
}

func (s *IntegrationTestBuilder) AssertBuildCountGitInfo(count int) {
	s.Helper()
	s.Assert(s.H.init.gitInfo.InitCount(), qt.Equals, count)
}

func (s *IntegrationTestBuilder) AssertBuildCountLayouts(count int) {
	s.Helper()
	s.Assert(s.H.init.layouts.InitCount(), qt.Equals, count)
}

func (s *IntegrationTestBuilder) AssertFileCount(dirname string, expected int) {
	s.Helper()
	fs := s.fs.WorkingDirReadOnly
	count := 0
	afero.Walk(fs, dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	s.Assert(count, qt.Equals, expected)
}

func (s *IntegrationTestBuilder) negate(match string) (string, bool) {
	var negate bool
	if strings.HasPrefix(match, "! ") {
		negate = true
		match = strings.TrimPrefix(match, "! ")
	}
	return match, negate
}

func (s *IntegrationTestBuilder) AssertFileContent(filename string, matches ...string) {
	s.Helper()
	content := strings.TrimSpace(s.FileContent(filename))
	for _, m := range matches {
		cm := qt.Commentf("File: %s Match %s", filename, m)
		lines := strings.Split(m, "\n")
		for _, match := range lines {
			match = strings.TrimSpace(match)
			if match == "" || strings.HasPrefix(match, "#") {
				continue
			}
			var negate bool
			match, negate = s.negate(match)
			if negate {
				s.Assert(content, qt.Not(qt.Contains), match, cm)
				continue
			}
			s.Assert(content, qt.Contains, match, cm)
		}
	}
}

func (s *IntegrationTestBuilder) AssertFileContentExact(filename string, matches ...string) {
	s.Helper()
	content := s.FileContent(filename)
	for _, m := range matches {
		s.Assert(content, qt.Contains, m, qt.Commentf(m))
	}
}

func (s *IntegrationTestBuilder) AssertPublishDir(matches ...string) {
	s.AssertFs(s.fs.PublishDir, matches...)
}

func (s *IntegrationTestBuilder) AssertFs(fs afero.Fs, matches ...string) {
	s.Helper()
	var buff bytes.Buffer
	s.Assert(s.printAndCheckFs(fs, "", &buff), qt.IsNil)
	printFsLines := strings.Split(buff.String(), "\n")
	sort.Strings(printFsLines)
	content := strings.TrimSpace((strings.Join(printFsLines, "\n")))
	for _, m := range matches {
		cm := qt.Commentf("Match: %q\nIn:\n%s", m, content)
		lines := strings.Split(m, "\n")
		for _, match := range lines {
			match = strings.TrimSpace(match)
			var negate bool
			if strings.HasPrefix(match, "! ") {
				negate = true
				match = strings.TrimPrefix(match, "! ")
			}
			if negate {
				s.Assert(content, qt.Not(qt.Contains), match, cm)
				continue
			}
			s.Assert(content, qt.Contains, match, cm)
		}
	}
}

func (s *IntegrationTestBuilder) printAndCheckFs(fs afero.Fs, path string, w io.Writer) error {
	if fs == nil {
		return nil
	}

	return afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error: path %q: %s", path, err)
		}
		path = filepath.ToSlash(path)
		if path == "" {
			path = "."
		}
		if !info.IsDir() {
			f, err := fs.Open(path)
			if err != nil {
				return fmt.Errorf("error: path %q: %s", path, err)
			}
			defer f.Close()
			// This will panic if the file is a directory.
			var buf [1]byte
			io.ReadFull(f, buf[:])
		}
		fmt.Fprintln(w, path, info.IsDir())
		return nil
	})
}

func (s *IntegrationTestBuilder) AssertFileExists(filename string, b bool) {
	checker := qt.IsNil
	if !b {
		checker = qt.IsNotNil
	}
	_, err := s.fs.WorkingDirReadOnly.Stat(filename)
	if !herrors.IsNotExist(err) {
		s.Assert(err, qt.IsNil)
	}
	s.Assert(err, checker)
}

func (s *IntegrationTestBuilder) AssertIsFileError(err error) herrors.FileError {
	s.Assert(err, qt.ErrorAs, new(herrors.FileError))
	return herrors.UnwrapFileError(err)
}

func (s *IntegrationTestBuilder) AssertRenderCountContent(count int) {
	s.Helper()
	s.Assert(s.counters.contentRenderCounter.Load(), qt.Equals, uint64(count))
}

func (s *IntegrationTestBuilder) AssertRenderCountPage(count int) {
	s.Helper()
	s.Assert(s.counters.pageRenderCounter.Load(), qt.Equals, uint64(count))
}

func (s *IntegrationTestBuilder) AssertRenderCountPageBetween(from, to int) {
	s.Helper()
	i := int(s.counters.pageRenderCounter.Load())
	s.Assert(i >= from && i <= to, qt.IsTrue)
}

func (s *IntegrationTestBuilder) Build() *IntegrationTestBuilder {
	s.Helper()
	_, err := s.BuildE()
	if s.Cfg.Verbose || err != nil {
		fmt.Println(s.lastBuildLog)
		if s.H != nil && err == nil {
			for _, s := range s.H.Sites {
				m := s.pageMap
				var buff bytes.Buffer
				fmt.Fprintf(&buff, "PageMap for site %q\n\n", s.Language().Lang)
				m.debugPrint("", 999, &buff)
				fmt.Println(buff.String())
			}
		}
	} else if s.Cfg.LogLevel <= logg.LevelDebug {
		fmt.Println(s.lastBuildLog)
	}
	s.Assert(err, qt.IsNil)
	if s.Cfg.RunGC {
		s.GCCount, err = s.H.GC()
		s.Assert(err, qt.IsNil)
	}

	return s
}

func (s *IntegrationTestBuilder) LogString() string {
	return s.lastBuildLog
}

func (s *IntegrationTestBuilder) BuildE() (*IntegrationTestBuilder, error) {
	s.Helper()
	if err := s.initBuilder(); err != nil {
		return s, err
	}

	err := s.build(s.Cfg.BuildCfg)
	return s, err
}

func (s *IntegrationTestBuilder) Init() *IntegrationTestBuilder {
	if err := s.initBuilder(); err != nil {
		s.Fatalf("Failed to init builder: %s", err)
	}
	s.lastBuildLog = s.logBuff.String()
	return s
}

type IntegrationTestDebugConfig struct {
	Out io.Writer

	PrintDestinationFs bool
	PrintPagemap       bool

	PrefixDestinationFs string
	PrefixPagemap       string
}

func (s *IntegrationTestBuilder) EditFileReplaceAll(filename, old, new string) *IntegrationTestBuilder {
	return s.EditFileReplaceFunc(filename, func(s string) string {
		return strings.ReplaceAll(s, old, new)
	})
}

func (s *IntegrationTestBuilder) EditFileReplaceFunc(filename string, replacementFunc func(s string) string) *IntegrationTestBuilder {
	absFilename := s.absFilename(filename)
	b, err := afero.ReadFile(s.fs.Source, absFilename)
	s.Assert(err, qt.IsNil)
	s.changedFiles = append(s.changedFiles, absFilename)
	oldContent := string(b)
	s.writeSource(absFilename, replacementFunc(oldContent))
	return s
}

func (s *IntegrationTestBuilder) EditFiles(filenameContent ...string) *IntegrationTestBuilder {
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filepath.FromSlash(filenameContent[i]), filenameContent[i+1]
		absFilename := s.absFilename(filename)
		s.changedFiles = append(s.changedFiles, absFilename)
		s.writeSource(absFilename, content)
	}
	return s
}

func (s *IntegrationTestBuilder) AddFiles(filenameContent ...string) *IntegrationTestBuilder {
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filepath.FromSlash(filenameContent[i]), filenameContent[i+1]
		absFilename := s.absFilename(filename)
		s.createdFiles = append(s.createdFiles, absFilename)
		s.writeSource(absFilename, content)
	}
	return s
}

func (s *IntegrationTestBuilder) RemoveFiles(filenames ...string) *IntegrationTestBuilder {
	for _, filename := range filenames {
		absFilename := s.absFilename(filename)
		s.removedFiles = append(s.removedFiles, absFilename)
		s.Assert(s.fs.Source.Remove(absFilename), qt.IsNil)

	}

	return s
}

func (s *IntegrationTestBuilder) RenameFile(old, new string) *IntegrationTestBuilder {
	absOldFilename := s.absFilename(old)
	absNewFilename := s.absFilename(new)
	s.renamedFiles = append(s.renamedFiles, absOldFilename)
	s.createdFiles = append(s.createdFiles, absNewFilename)
	s.Assert(s.fs.Source.MkdirAll(filepath.Dir(absNewFilename), 0o777), qt.IsNil)
	s.Assert(s.fs.Source.Rename(absOldFilename, absNewFilename), qt.IsNil)
	return s
}

func (s *IntegrationTestBuilder) RenameDir(old, new string) *IntegrationTestBuilder {
	absOldFilename := s.absFilename(old)
	absNewFilename := s.absFilename(new)
	s.renamedDirs = append(s.renamedDirs, absOldFilename)
	s.changedFiles = append(s.changedFiles, absNewFilename)
	afero.Walk(s.fs.Source, absOldFilename, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		s.createdFiles = append(s.createdFiles, strings.Replace(path, absOldFilename, absNewFilename, 1))
		return nil
	})
	s.Assert(s.fs.Source.MkdirAll(filepath.Dir(absNewFilename), 0o777), qt.IsNil)
	s.Assert(s.fs.Source.Rename(absOldFilename, absNewFilename), qt.IsNil)
	return s
}

func (s *IntegrationTestBuilder) FileContent(filename string) string {
	s.Helper()
	return s.readWorkingDir(s, s.fs, filepath.FromSlash(filename))
}

func (s *IntegrationTestBuilder) initBuilder() error {
	var initErr error
	s.builderInit.Do(func() {
		var afs afero.Fs
		if s.Cfg.NeedsOsFS {
			afs = afero.NewOsFs()
		} else {
			afs = afero.NewMemMapFs()
		}

		if s.Cfg.LogLevel == 0 {
			s.Cfg.LogLevel = logg.LevelError
		}

		isBinaryRe := regexp.MustCompile(`^(.*)(\.png|\.jpg)$`)

		const dataSourceFilenamePrefix = "sourcefilename:"

		for _, f := range s.data.Files {
			filename := filepath.Join(s.Cfg.WorkingDir, f.Name)
			data := bytes.TrimSuffix(f.Data, []byte("\n"))
			datastr := strings.TrimSpace(string(data))
			if strings.HasPrefix(datastr, dataSourceFilenamePrefix) {
				// Read from file relative to the current dir.
				var err error
				wd, _ := os.Getwd()
				filename := filepath.Join(wd, strings.TrimSpace(strings.TrimPrefix(datastr, dataSourceFilenamePrefix)))
				data, err = os.ReadFile(filename)
				s.Assert(err, qt.IsNil)
			} else if isBinaryRe.MatchString(filename) {
				var err error
				data, err = base64.StdEncoding.DecodeString(string(data))
				s.Assert(err, qt.IsNil)

			}
			s.Assert(afs.MkdirAll(filepath.Dir(filename), 0o777), qt.IsNil)
			s.Assert(afero.WriteFile(afs, filename, data, 0o666), qt.IsNil)
		}

		configDir := "config"
		if _, err := afs.Stat(filepath.Join(s.Cfg.WorkingDir, "config")); err != nil {
			configDir = ""
		}

		var flags config.Provider
		if s.Cfg.BaseCfg != nil {
			flags = s.Cfg.BaseCfg
		} else {
			flags = config.New()
		}

		if s.Cfg.Running {
			flags.Set("internal", maps.Params{
				"running": s.Cfg.Running,
				"watch":   s.Cfg.Running,
			})
		} else if s.Cfg.Watching {
			flags.Set("internal", maps.Params{
				"watch": s.Cfg.Watching,
			})
		}

		if s.Cfg.WorkingDir != "" {
			flags.Set("workingDir", s.Cfg.WorkingDir)
		}

		var w io.Writer
		if s.Cfg.LogLevel == logg.LevelTrace {
			w = os.Stdout
		} else {
			w = &s.logBuff
		}

		logger := loggers.New(
			loggers.Options{
				Stdout:        w,
				Stderr:        w,
				Level:         s.Cfg.LogLevel,
				DistinctLevel: logg.LevelWarn,
			},
		)

		res, err := allconfig.LoadConfig(
			allconfig.ConfigSourceDescriptor{
				Flags:     flags,
				ConfigDir: configDir,
				Fs:        afs,
				Logger:    logger,
				Environ:   s.Cfg.Environ,
			},
		)
		if err != nil {
			initErr = err
			return
		}

		fs := hugofs.NewFrom(afs, res.LoadingInfo.BaseConfig)

		s.Assert(err, qt.IsNil)

		depsCfg := deps.DepsCfg{Configs: res, Fs: fs, LogLevel: logger.Level(), LogOut: logger.Out()}
		sites, err := NewHugoSites(depsCfg)
		if err != nil {
			initErr = err
			return
		}
		if sites == nil {
			initErr = errors.New("no sites")
			return
		}

		s.H = sites
		s.fs = fs

		if s.Cfg.NeedsNpmInstall {
			wd, _ := os.Getwd()
			s.Assert(os.Chdir(s.Cfg.WorkingDir), qt.IsNil)
			s.C.Cleanup(func() { os.Chdir(wd) })
			sc := security.DefaultConfig
			sc.Exec.Allow, err = security.NewWhitelist("npm")
			s.Assert(err, qt.IsNil)
			ex := hexec.New(sc, s.Cfg.WorkingDir)
			command, err := ex.New("npm", "install")
			s.Assert(err, qt.IsNil)
			s.Assert(command.Run(), qt.IsNil)

		}
	})

	return initErr
}

func (s *IntegrationTestBuilder) absFilename(filename string) string {
	filename = filepath.FromSlash(filename)
	if filepath.IsAbs(filename) {
		return filename
	}
	if s.Cfg.WorkingDir != "" && !strings.HasPrefix(filename, s.Cfg.WorkingDir) {
		filename = filepath.Join(s.Cfg.WorkingDir, filename)
	}
	return filename
}

func (s *IntegrationTestBuilder) reset() {
	s.changedFiles = nil
	s.createdFiles = nil
	s.removedFiles = nil
	s.renamedFiles = nil
}

func (s *IntegrationTestBuilder) build(cfg BuildCfg) error {
	s.Helper()
	defer func() {
		s.reset()
		s.lastBuildLog = s.logBuff.String()
		s.logBuff.Reset()
	}()

	changeEvents := s.changeEvents()
	s.counters = &buildCounters{}
	cfg.testCounters = s.counters

	if s.buildCount > 0 && (len(changeEvents) == 0) {
		return nil
	}

	s.buildCount++

	err := s.H.Build(cfg, changeEvents...)
	if err != nil {
		return err
	}

	return nil
}

// We simulate the fsnotify events.
// See the test output in https://github.com/bep/fsnotifyeventlister for what events gets produced
// by the different OSes.
func (s *IntegrationTestBuilder) changeEvents() []fsnotify.Event {
	var (
		events    []fsnotify.Event
		isLinux   = runtime.GOOS == "linux"
		isMacOs   = runtime.GOOS == "darwin"
		isWindows = runtime.GOOS == "windows"
	)

	for _, v := range s.removedFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Remove,
		})
	}
	for _, v := range s.renamedFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Rename,
		})
	}

	for _, v := range s.renamedDirs {
		events = append(events, fsnotify.Event{
			Name: v,
			// This is what we get on MacOS.
			Op: fsnotify.Remove | fsnotify.Rename,
		})
	}

	for _, v := range s.changedFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Write,
		})
		if isLinux || isWindows {
			// Duplicate write events, for some reason.
			events = append(events, fsnotify.Event{
				Name: v,
				Op:   fsnotify.Write,
			})
		}
		if isMacOs {
			events = append(events, fsnotify.Event{
				Name: v,
				Op:   fsnotify.Chmod,
			})
		}
	}
	for _, v := range s.createdFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Create,
		})
		if isLinux || isWindows {
			events = append(events, fsnotify.Event{
				Name: v,
				Op:   fsnotify.Write,
			})
		}

	}

	// Shuffle events.
	for i := range events {
		j := rand.Intn(i + 1)
		events[i], events[j] = events[j], events[i]
	}

	return events
}

func (s *IntegrationTestBuilder) readWorkingDir(t testing.TB, fs *hugofs.Fs, filename string) string {
	t.Helper()
	return s.readFileFromFs(t, fs.WorkingDirReadOnly, filename)
}

func (s *IntegrationTestBuilder) readFileFromFs(t testing.TB, fs afero.Fs, filename string) string {
	t.Helper()
	filename = filepath.Clean(filename)
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		// Print some debug info
		hadSlash := strings.HasPrefix(filename, helpers.FilePathSeparator)
		start := 0
		if hadSlash {
			start = 1
		}
		end := start + 1

		parts := strings.Split(filename, helpers.FilePathSeparator)
		if parts[start] == "work" {
			end++
		}

		s.Assert(err, qt.IsNil)

	}
	return string(b)
}

func (s *IntegrationTestBuilder) writeSource(filename, content string) {
	s.Helper()
	s.writeToFs(s.fs.Source, filename, content)
}

func (s *IntegrationTestBuilder) writeToFs(fs afero.Fs, filename, content string) {
	s.Helper()
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0o755); err != nil {
		s.Fatalf("Failed to write file: %s", err)
	}
}

type IntegrationTestConfig struct {
	T testing.TB

	// The files to use on txtar format, see
	// https://pkg.go.dev/golang.org/x/exp/cmd/txtar
	TxtarString string

	// COnfig to use as the base. We will also read the config from the txtar.
	BaseCfg config.Provider

	// Environment variables passed to the config loader.
	Environ []string

	// Whether to simulate server mode.
	Running bool

	// Watch for changes.
	// This is (currently) always set to true when Running is set.
	// Note that the CLI for the server does allow for --watch=false, but that is not used in these test.
	Watching bool

	// Will print the log buffer after the build
	Verbose bool

	// The log level to use.
	LogLevel logg.Level

	// Whether it needs the real file system (e.g. for js.Build tests).
	NeedsOsFS bool

	// Whether to run GC after each build.
	RunGC bool

	// Do not remove the temp dir after the test.
	PrintAndKeepTempDir bool

	// Whether to run npm install before Build.
	NeedsNpmInstall bool

	// Whether to normalize the Unicode filenames to NFD on Darwin.
	NFDFormOnDarwin bool

	// The working dir to use. If not absolute, a temp dir will be created.
	WorkingDir string

	// The config to pass to Build.
	BuildCfg BuildCfg
}
