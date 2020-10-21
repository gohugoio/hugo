// Copyright 2020 The Hugo Authors. All rights reserved.
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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/gohugoio/hugo/common/terminal"

	jww "github.com/spf13/jwalterweatherman"
)

var (
	// Counts ERROR logs to the global jww logger.
	GlobalErrorCounter *jww.Counter
)

func init() {
	GlobalErrorCounter = &jww.Counter{}
	jww.SetLogListeners(jww.LogCounter(GlobalErrorCounter, jww.LevelError))
}

func LoggerToWriterWithPrefix(logger *log.Logger, prefix string) io.Writer {
	return prefixWriter{
		logger: logger,
		prefix: prefix,
	}
}

type prefixWriter struct {
	logger *log.Logger
	prefix string
}

func (w prefixWriter) Write(p []byte) (n int, err error) {
	w.logger.Printf("%s: %s", w.prefix, p)
	return len(p), nil
}

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	PrintTimerIfDelayed(start time.Time, name string)
	Debug() *log.Logger
	Info() *log.Logger
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Warn() *log.Logger
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Error() *log.Logger
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Errors() string

	Out() io.Writer

	Reset()

	// Used in tests.
	LogCounters() *LogCounters
}

type LogCounters struct {
	ErrorCounter *jww.Counter
	WarnCounter  *jww.Counter
}

type logger struct {
	*jww.Notepad

	// The writer that represents stdout.
	// Will be ioutil.Discard when in quiet mode.
	out io.Writer

	logCounters *LogCounters

	// This is only set in server mode.
	errors *bytes.Buffer
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.FEEDBACK.Printf(format, v...)
}

func (l *logger) Println(v ...interface{}) {
	l.FEEDBACK.Println(v...)
}

func (l *logger) Debug() *log.Logger {
	return l.DEBUG
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.INFO.Printf(format, v...)
}

func (l *logger) Infoln(v ...interface{}) {
	l.INFO.Println(v...)
}

func (l *logger) Info() *log.Logger {
	return l.INFO
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.WARN.Printf(format, v...)
}

func (l *logger) Warnln(v ...interface{}) {
	l.WARN.Println(v...)
}

func (l *logger) Warn() *log.Logger {
	return l.WARN
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.ERROR.Printf(format, v...)
}

func (l *logger) Errorln(v ...interface{}) {
	l.ERROR.Println(v...)
}

func (l *logger) Error() *log.Logger {
	return l.ERROR
}

func (l *logger) LogCounters() *LogCounters {
	return l.logCounters
}

func (l *logger) Out() io.Writer {
	return l.out
}

// PrintTimerIfDelayed prints a time statement to the FEEDBACK logger
// if considerable time is spent.
func (l *logger) PrintTimerIfDelayed(start time.Time, name string) {
	elapsed := time.Since(start)
	milli := int(1000 * elapsed.Seconds())
	if milli < 500 {
		return
	}
	l.Printf("%s in %v ms", name, milli)
}

func (l *logger) PrintTimer(start time.Time, name string) {
	elapsed := time.Since(start)
	milli := int(1000 * elapsed.Seconds())
	l.Printf("%s in %v ms", name, milli)
}

func (l *logger) Errors() string {
	if l.errors == nil {
		return ""
	}
	return ansiColorRe.ReplaceAllString(l.errors.String(), "")
}

// Reset resets the logger's internal state.
func (l *logger) Reset() {
	l.logCounters.ErrorCounter.Reset()
	if l.errors != nil {
		l.errors.Reset()
	}
}

//  NewLogger creates a new Logger for the given thresholds
func NewLogger(stdoutThreshold, logThreshold jww.Threshold, outHandle, logHandle io.Writer, saveErrors bool) Logger {
	return newLogger(stdoutThreshold, logThreshold, outHandle, logHandle, saveErrors)
}

// NewDebugLogger is a convenience function to create a debug logger.
func NewDebugLogger() Logger {
	return NewBasicLogger(jww.LevelDebug)
}

// NewWarningLogger is a convenience function to create a warning logger.
func NewWarningLogger() Logger {
	return NewBasicLogger(jww.LevelWarn)
}

// NewInfoLogger is a convenience function to create a info logger.
func NewInfoLogger() Logger {
	return NewBasicLogger(jww.LevelInfo)
}

// NewErrorLogger is a convenience function to create an error logger.
func NewErrorLogger() Logger {
	return NewBasicLogger(jww.LevelError)
}

// NewBasicLogger creates a new basic logger writing to Stdout.
func NewBasicLogger(t jww.Threshold) Logger {
	return newLogger(t, jww.LevelError, os.Stdout, ioutil.Discard, false)
}

// NewBasicLoggerForWriter creates a new basic logger writing to w.
func NewBasicLoggerForWriter(t jww.Threshold, w io.Writer) Logger {
	return newLogger(t, jww.LevelError, w, ioutil.Discard, false)
}

var (
	ansiColorRe = regexp.MustCompile("(?s)\\033\\[\\d*(;\\d*)*m")
	errorRe     = regexp.MustCompile("^(ERROR|FATAL|WARN)")
)

type ansiCleaner struct {
	w io.Writer
}

func (a ansiCleaner) Write(p []byte) (n int, err error) {
	return a.w.Write(ansiColorRe.ReplaceAll(p, []byte("")))
}

type labelColorizer struct {
	w io.Writer
}

func (a labelColorizer) Write(p []byte) (n int, err error) {
	replaced := errorRe.ReplaceAllStringFunc(string(p), func(m string) string {
		switch m {
		case "ERROR", "FATAL":
			return terminal.Error(m)
		case "WARN":
			return terminal.Warning(m)
		default:
			return m
		}
	})
	// io.MultiWriter will abort if we return a bigger write count than input
	// bytes, so we lie a little.
	_, err = a.w.Write([]byte(replaced))
	return len(p), err

}

// InitGlobalLogger initializes the global logger, used in some rare cases.
func InitGlobalLogger(stdoutThreshold, logThreshold jww.Threshold, outHandle, logHandle io.Writer) {
	outHandle, logHandle = getLogWriters(outHandle, logHandle)

	jww.SetStdoutOutput(outHandle)
	jww.SetLogOutput(logHandle)
	jww.SetLogThreshold(logThreshold)
	jww.SetStdoutThreshold(stdoutThreshold)

}

func getLogWriters(outHandle, logHandle io.Writer) (io.Writer, io.Writer) {
	isTerm := terminal.IsTerminal(os.Stdout)
	if logHandle != ioutil.Discard && isTerm {
		// Remove any Ansi coloring from log output
		logHandle = ansiCleaner{w: logHandle}
	}

	if isTerm {
		outHandle = labelColorizer{w: outHandle}
	}

	return outHandle, logHandle

}

type fatalLogWriter int

func (s fatalLogWriter) Write(p []byte) (n int, err error) {
	trace := make([]byte, 1500)
	runtime.Stack(trace, true)
	fmt.Printf("\n===========\n\n%s\n", trace)
	os.Exit(-1)

	return 0, nil
}

var fatalLogListener = func(t jww.Threshold) io.Writer {
	if t != jww.LevelError {
		// Only interested in ERROR
		return nil
	}

	return new(fatalLogWriter)
}

func newLogger(stdoutThreshold, logThreshold jww.Threshold, outHandle, logHandle io.Writer, saveErrors bool) *logger {
	errorCounter := &jww.Counter{}
	warnCounter := &jww.Counter{}
	outHandle, logHandle = getLogWriters(outHandle, logHandle)

	listeners := []jww.LogListener{jww.LogCounter(errorCounter, jww.LevelError), jww.LogCounter(warnCounter, jww.LevelWarn)}
	var errorBuff *bytes.Buffer
	if saveErrors {
		errorBuff = new(bytes.Buffer)
		errorCapture := func(t jww.Threshold) io.Writer {
			if t != jww.LevelError {
				// Only interested in ERROR
				return nil
			}
			return errorBuff
		}

		listeners = append(listeners, errorCapture)
	}

	return &logger{
		Notepad: jww.NewNotepad(stdoutThreshold, logThreshold, outHandle, logHandle, "", log.Ldate|log.Ltime, listeners...),
		out:     outHandle,
		logCounters: &LogCounters{
			ErrorCounter: errorCounter,
			WarnCounter:  warnCounter,
		},
		errors: errorBuff,
	}
}
