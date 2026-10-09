// Copyright 2026 The Hugo Authors. All rights reserved.
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

package terminal

import (
	"encoding/base64"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// See issue 15425.
func TestReportProgress(t *testing.T) {
	c := qt.New(t)

	report := func(state ProgressState, progress float64, msg string) string {
		var b strings.Builder
		ReportProgress(&b, state, progress, msg)
		return b.String()
	}

	c.Assert(report(ProgressIntermediate, 0, ""), qt.Equals, "\033]9;4;3;0\007\033]7501;state=working:app=hugo\033\\")
	c.Assert(report(ProgressNormal, 0.42, ""), qt.Equals, "\033]9;4;1;42\007\033]7501;state=working:progress=42:app=hugo\033\\")
	c.Assert(report(ProgressNormal, 1.5, ""), qt.Equals, "\033]9;4;1;100\007\033]7501;state=working:progress=100:app=hugo\033\\")
	c.Assert(report(ProgressHidden, 1, ""), qt.Equals, "\033]9;4;0;100\007\033]7501;state=idle:app=hugo\033\\")
	c.Assert(report(ProgressDone, 1, ""), qt.Equals, "\033]9;4;0;100\007\033]7501;state=done:app=hugo\033\\")

	msg := base64.StdEncoding.EncodeToString([]byte("failed: boom"))
	c.Assert(report(ProgressError, 1, Error("failed: boom")+"\nmore\n"), qt.Equals, "\033]9;4;2;100\007\033]7501;state=error:app=hugo:msg="+msg+"\033\\")
}

func TestStatusMsg(t *testing.T) {
	c := qt.New(t)

	c.Assert(statusMsg("  \x1b[1;31mfoo\x1b[0m\x07 bar\nbaz"), qt.Equals, "foo bar")
	c.Assert(statusMsg(strings.Repeat("a", 3000)), qt.HasLen, maxStatusMsgLen)
	c.Assert(statusMsg(strings.Repeat("a", maxStatusMsgLen-1)+"ø"), qt.Equals, strings.Repeat("a", maxStatusMsgLen-1))
}
