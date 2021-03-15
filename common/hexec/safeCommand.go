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

package hexec

import (
	"context"

	"os/exec"

	"github.com/cli/safeexec"
)

// SafeCommand is a wrapper around os/exec Command which uses a LookPath
// implementation that does not search in current directory before looking in PATH.
// See https://github.com/cli/safeexec and the linked issues.
func SafeCommand(name string, arg ...string) (*exec.Cmd, error) {
	bin, err := safeexec.LookPath(name)
	if err != nil {
		return nil, err
	}

	return exec.Command(bin, arg...), nil
}

// SafeCommandContext wraps CommandContext
// See SafeCommand for more context.
func SafeCommandContext(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	bin, err := safeexec.LookPath(name)
	if err != nil {
		return nil, err
	}

	return exec.CommandContext(ctx, bin, arg...), nil
}
