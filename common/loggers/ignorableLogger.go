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
	"strings"
)

// IgnorableLogger is a logger that ignores certain log statements.
type IgnorableLogger interface {
	Logger
	Errorsf(statementID, format string, v ...interface{})
}

type ignorableLogger struct {
	Logger
	statements map[string]bool
}

// NewIgnorableLogger wraps the given logger and ignores the log statement IDs given.
func NewIgnorableLogger(logger Logger, statements ...string) IgnorableLogger {
	statementsSet := make(map[string]bool)
	for _, s := range statements {
		statementsSet[strings.ToLower(s)] = true

	}
	return ignorableLogger{
		Logger:     logger,
		statements: statementsSet,
	}
}

// Errorsf logs statementID as an ERROR if not configured as ignoreable.
func (l ignorableLogger) Errorsf(statementID, format string, v ...interface{}) {
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
