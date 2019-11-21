// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Logger wraps a *loggers.Logger and some other related logging state.
type Logger struct {
	*jww.Notepad

	// The writer that represents stdout.
	// Will be ioutil.Discard when in quiet mode.
	Out io.Writer

	ErrorCounter *jww.Counter
	WarnCounter  *jww.Counter

	// This is only set in server mode.
	errors *bytes.Buffer
}

// PrintTimerIfDelayed prints a time statement to the FEEDBACK logger
// if considerable time is spent.
func (l *Logger) PrintTimerIfDelayed(start time.Time, name string) {
	elapsed := time.Since(start)
	milli := int(1000 * elapsed.Seconds())
	if milli < 500 {
		return
	}
	l.FEEDBACK.Printf("%s in %v ms", name, milli)
}

func (l *Logger) PrintTimer(start time.Time, name string) {
	elapsed := time.Since(start)
	milli := int(1000 * elapsed.Seconds())
	l.FEEDBACK.Printf("%s in %v ms", name, milli)
}

func (l *Logger) Errors() string {
	if l.errors == nil {
		return ""
	}
	return ansiColorRe.ReplaceAllString(l.errors.String(), "")
}

// Reset resets the logger's internal state.
func (l *Logger) Reset() {
	l.ErrorCounter.Reset()
	if l.errors != nil {
		l.errors.Reset()
	}
}

//  NewLogger creates a new Logger for the given thresholds
func NewLogger(stdoutThreshold, logThreshold jww.Threshold, outHandle, logHandle io.Writer, saveErrors bool) *Logger {
	return newLogger(stdoutThreshold, logThreshold, outHandle, logHandle, saveErrors)
}

// NewDebugLogger is a convenience function to create a debug logger.
func NewDebugLogger() *Logger {
	return newBasicLogger(jww.LevelDebug)
}

// NewWarningLogger is a convenience function to create a warning logger.
func NewWarningLogger() *Logger {
	return newBasicLogger(jww.LevelWarn)
}

// NewErrorLogger is a convenience function to create an error logger.
func NewErrorLogger() *Logger {
	return newBasicLogger(jww.LevelError)
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

func newLogger(stdoutThreshold, logThreshold jww.Threshold, outHandle, logHandle io.Writer, saveErrors bool) *Logger {
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

	return &Logger{
		Notepad:      jww.NewNotepad(stdoutThreshold, logThreshold, outHandle, logHandle, "", log.Ldate|log.Ltime, listeners...),
		Out:          outHandle,
		ErrorCounter: errorCounter,
		WarnCounter:  warnCounter,
		errors:       errorBuff,
	}
}

func newBasicLogger(t jww.Threshold) *Logger {
	return newLogger(t, jww.LevelError, os.Stdout, ioutil.Discard, false)
}
