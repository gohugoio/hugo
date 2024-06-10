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
	"sync"

	"github.com/bep/logg"
)

func InitGlobalLogger(level logg.Level, panicOnWarnings bool) {
	logMu.Lock()
	defer logMu.Unlock()
	var logHookLast func(e *logg.Entry) error
	if panicOnWarnings {
		logHookLast = PanicOnWarningHook
	}

	log = New(
		Options{
			Level:         level,
			DistinctLevel: logg.LevelInfo,
			HandlerPost:   logHookLast,
		},
	)
}

var logMu sync.Mutex

func Log() Logger {
	logMu.Lock()
	defer logMu.Unlock()
	return log
}

// The global logger.
var log Logger

func init() {
	InitGlobalLogger(logg.LevelWarn, false)
}
