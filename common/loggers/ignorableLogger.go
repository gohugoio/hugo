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
	"fmt"
)

// IgnorableLogger is a logger that ignores certain log statements.
type IgnorableLogger interface {
	Logger
	Errorsf(statementID, format string, v ...any)
	Apply(logger Logger) IgnorableLogger
}

type ignorableLogger struct {
	Logger
	statements map[string]bool
}

// NewIgnorableLogger wraps the given logger and ignores the log statement IDs given.
func NewIgnorableLogger(logger Logger, statements map[string]bool) IgnorableLogger {
	if statements == nil {
		statements = make(map[string]bool)
	}
	return ignorableLogger{
		Logger:     logger,
		statements: statements,
	}
}

// Errorsf logs statementID as an ERROR if not configured as ignoreable.
func (l ignorableLogger) Errorsf(statementID, format string, v ...any) {
	if l.statements[statementID] {
		// Ignore.
		return
	}
	ignoreMsg := fmt.Sprintf(`
If you feel that this should not be logged as an ERROR, you can ignore it by adding this to your site config:
ignoreErrors = [%q]`, statementID)

	format += ignoreMsg

	l.Errorf(format, v...)
}

func (l ignorableLogger) Apply(logger Logger) IgnorableLogger {
	return ignorableLogger{
		Logger:     logger,
		statements: l.statements,
	}
}
