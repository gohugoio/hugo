// Copyright 2024 The Hugo Authors. All rights reserved.
// Some functions in this file (see comments) is based on the Go source code,
// copyright The Go Authors and  governed by a BSD-style license.
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

package loggers

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/multi"
	"github.com/gohugoio/hugo/common/terminal"
)

var (
	reservedFieldNamePrefix = "__h_field_"
	// FieldNameCmd is the name of the field that holds the command name.
	FieldNameCmd = reservedFieldNamePrefix + "_cmd"
	// Used to suppress statements.
	FieldNameStatementID = reservedFieldNamePrefix + "__h_field_statement_id"
)

// Options defines options for the logger.
type Options struct {
	Level              logg.Level
	Stdout             io.Writer
	Stderr             io.Writer
	DistinctLevel      logg.Level
	StoreErrors        bool
	HandlerPost        func(e *logg.Entry) error
	SuppressStatements map[string]bool
}

// New creates a new logger with the given options.
func New(opts Options) Logger {
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stdout
	}
	if opts.Level == 0 {
		opts.Level = logg.LevelWarn
	}

	var logHandler logg.Handler
	if terminal.PrintANSIColors(os.Stdout) {
		logHandler = newDefaultHandler(opts.Stdout, opts.Stderr)
	} else {
		logHandler = newNoColoursHandler(opts.Stdout, opts.Stderr, false, nil)
	}

	errorsw := &strings.Builder{}
	logCounters := newLogLevelCounter()
	handlers := []logg.Handler{
		logCounters,
	}

	if opts.Level == logg.LevelTrace {
		// Trace is used during development only, and it's useful to
		// only see the trace messages.
		handlers = append(handlers,
			logg.HandlerFunc(func(e *logg.Entry) error {
				if e.Level != logg.LevelTrace {
					return logg.ErrStopLogEntry
				}
				return nil
			}),
		)
	}

	handlers = append(handlers, whiteSpaceTrimmer(), logHandler)

	if opts.HandlerPost != nil {
		var hookHandler logg.HandlerFunc = func(e *logg.Entry) error {
			opts.HandlerPost(e)
			return nil
		}
		handlers = append(handlers, hookHandler)
	}

	if opts.StoreErrors {
		h := newNoColoursHandler(io.Discard, errorsw, true, func(e *logg.Entry) bool {
			return e.Level >= logg.LevelError
		})

		handlers = append(handlers, h)
	}

	logHandler = multi.New(handlers...)

	var logOnce *logOnceHandler
	if opts.DistinctLevel != 0 {
		logOnce = newLogOnceHandler(opts.DistinctLevel)
		logHandler = newStopHandler(logOnce, logHandler)
	}

	if len(opts.SuppressStatements) > 0 {
		logHandler = newStopHandler(newSuppressStatementsHandler(opts.SuppressStatements), logHandler)
	}

	logger := logg.New(
		logg.Options{
			Level:   opts.Level,
			Handler: logHandler,
		},
	)

	l := logger.WithLevel(opts.Level)

	reset := func() {
		logCounters.mu.Lock()
		defer logCounters.mu.Unlock()
		logCounters.counters = make(map[logg.Level]int)
		errorsw.Reset()
		if logOnce != nil {
			logOnce.reset()
		}
	}

	return &logAdapter{
		logCounters: logCounters,
		errors:      errorsw,
		reset:       reset,
		out:         opts.Stdout,
		level:       opts.Level,
		logger:      logger,
		tracel:      l.WithLevel(logg.LevelTrace),
		debugl:      l.WithLevel(logg.LevelDebug),
		infol:       l.WithLevel(logg.LevelInfo),
		warnl:       l.WithLevel(logg.LevelWarn),
		errorl:      l.WithLevel(logg.LevelError),
	}
}

// NewDefault creates a new logger with the default options.
func NewDefault() Logger {
	opts := Options{
		DistinctLevel: logg.LevelWarn,
		Level:         logg.LevelWarn,
		Stdout:        os.Stdout,
		Stderr:        os.Stdout,
	}
	return New(opts)
}

func NewTrace() Logger {
	opts := Options{
		DistinctLevel: logg.LevelWarn,
		Level:         logg.LevelTrace,
		Stdout:        os.Stdout,
		Stderr:        os.Stdout,
	}
	return New(opts)
}

func LevelLoggerToWriter(l logg.LevelLogger) io.Writer {
	return logWriter{l: l}
}

type Logger interface {
	Debug() logg.LevelLogger
	Debugf(format string, v ...any)
	Debugln(v ...any)
	Error() logg.LevelLogger
	Errorf(format string, v ...any)
	Erroridf(id, format string, v ...any)
	Errorln(v ...any)
	Errors() string
	Info() logg.LevelLogger
	InfoCommand(command string) logg.LevelLogger
	Infof(format string, v ...any)
	Infoln(v ...any)
	Level() logg.Level
	LoggCount(logg.Level) int
	Logger() logg.Logger
	Out() io.Writer
	Printf(format string, v ...any)
	Println(v ...any)
	PrintTimerIfDelayed(start time.Time, name string)
	Reset()
	Warn() logg.LevelLogger
	WarnCommand(command string) logg.LevelLogger
	Warnf(format string, v ...any)
	Warnidf(id, format string, v ...any)
	Warnln(v ...any)
	Deprecatef(fail bool, format string, v ...any)
	Trace(s logg.StringFunc)
}

type logAdapter struct {
	logCounters *logLevelCounter
	errors      *strings.Builder
	reset       func()
	out         io.Writer
	level       logg.Level
	logger      logg.Logger
	tracel      logg.LevelLogger
	debugl      logg.LevelLogger
	infol       logg.LevelLogger
	warnl       logg.LevelLogger
	errorl      logg.LevelLogger
}

func (l *logAdapter) Debug() logg.LevelLogger {
	return l.debugl
}

func (l *logAdapter) Debugf(format string, v ...any) {
	l.debugl.Logf(format, v...)
}

func (l *logAdapter) Debugln(v ...any) {
	l.debugl.Logf(l.sprint(v...))
}

func (l *logAdapter) Info() logg.LevelLogger {
	return l.infol
}

func (l *logAdapter) InfoCommand(command string) logg.LevelLogger {
	return l.infol.WithField(FieldNameCmd, command)
}

func (l *logAdapter) Infof(format string, v ...any) {
	l.infol.Logf(format, v...)
}

func (l *logAdapter) Infoln(v ...any) {
	l.infol.Logf(l.sprint(v...))
}

func (l *logAdapter) Level() logg.Level {
	return l.level
}

func (l *logAdapter) LoggCount(level logg.Level) int {
	l.logCounters.mu.RLock()
	defer l.logCounters.mu.RUnlock()
	return l.logCounters.counters[level]
}

func (l *logAdapter) Logger() logg.Logger {
	return l.logger
}

func (l *logAdapter) Out() io.Writer {
	return l.out
}

// PrintTimerIfDelayed prints a time statement to the FEEDBACK logger
// if considerable time is spent.
func (l *logAdapter) PrintTimerIfDelayed(start time.Time, name string) {
	elapsed := time.Since(start)
	milli := int(1000 * elapsed.Seconds())
	if milli < 500 {
		return
	}
	l.Printf("%s in %v ms", name, milli)
}

func (l *logAdapter) Printf(format string, v ...any) {
	// Add trailing newline if not present.
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(l.out, format, v...)
}

func (l *logAdapter) Println(v ...any) {
	fmt.Fprintln(l.out, v...)
}

func (l *logAdapter) Reset() {
	l.reset()
}

func (l *logAdapter) Warn() logg.LevelLogger {
	return l.warnl
}

func (l *logAdapter) Warnf(format string, v ...any) {
	l.warnl.Logf(format, v...)
}

func (l *logAdapter) WarnCommand(command string) logg.LevelLogger {
	return l.warnl.WithField(FieldNameCmd, command)
}

func (l *logAdapter) Warnln(v ...any) {
	l.warnl.Logf(l.sprint(v...))
}

func (l *logAdapter) Error() logg.LevelLogger {
	return l.errorl
}

func (l *logAdapter) Errorf(format string, v ...any) {
	l.errorl.Logf(format, v...)
}

func (l *logAdapter) Errorln(v ...any) {
	l.errorl.Logf(l.sprint(v...))
}

func (l *logAdapter) Errors() string {
	return l.errors.String()
}

func (l *logAdapter) Erroridf(id, format string, v ...any) {
	id = strings.ToLower(id)
	format += l.idfInfoStatement("error", id, format)
	l.errorl.WithField(FieldNameStatementID, id).Logf(format, v...)
}

func (l *logAdapter) Warnidf(id, format string, v ...any) {
	id = strings.ToLower(id)
	format += l.idfInfoStatement("warning", id, format)
	l.warnl.WithField(FieldNameStatementID, id).Logf(format, v...)
}

func (l *logAdapter) idfInfoStatement(what, id, format string) string {
	return fmt.Sprintf("\nYou can suppress this %s by adding the following to your site configuration:\nignoreLogs = ['%s']", what, id)
}

func (l *logAdapter) Trace(s logg.StringFunc) {
	l.tracel.Log(s)
}

func (l *logAdapter) sprint(v ...any) string {
	return strings.TrimRight(fmt.Sprintln(v...), "\n")
}

func (l *logAdapter) Deprecatef(fail bool, format string, v ...any) {
	format = "DEPRECATED: " + format
	if fail {
		l.errorl.Logf(format, v...)
	} else {
		l.warnl.Logf(format, v...)
	}
}

type logWriter struct {
	l logg.LevelLogger
}

func (w logWriter) Write(p []byte) (n int, err error) {
	w.l.Log(logg.String(string(p)))
	return len(p), nil
}

func TimeTrackf(l logg.LevelLogger, start time.Time, fields logg.Fields, format string, a ...any) {
	elapsed := time.Since(start)
	if fields != nil {
		l = l.WithFields(fields)
	}
	l.WithField("duration", elapsed).Logf(format, a...)
}

func TimeTrackfn(fn func() (logg.LevelLogger, error)) error {
	start := time.Now()
	l, err := fn()
	elapsed := time.Since(start)
	l.WithField("duration", elapsed).Logf("")
	return err
}
