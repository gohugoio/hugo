// Copyright 2015 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// HugoVersion represents the Hugo build version.
type HugoVersion struct {
	// Major and minor version.
	Number float32

	// Increment this for bug releases
	PatchLevel int

	// HugoVersionSuffix is the suffix used in the Hugo version string.
	// It will be blank for release versions.
	Suffix string
}

func (v HugoVersion) String() string {
	return hugoVersion(v.Number, v.PatchLevel, v.Suffix)
}

// ReleaseVersion represents the release version.
func (v HugoVersion) ReleaseVersion() HugoVersion {
	v.Suffix = ""
	return v
}

// Next returns the next Hugo release version.
func (v HugoVersion) Next() HugoVersion {
	return HugoVersion{Number: v.Number + 0.01}
}

// Pre returns the previous Hugo release version.
func (v HugoVersion) Prev() HugoVersion {
	return HugoVersion{Number: v.Number - 0.01}
}

// NextPatchLevel returns the next patch/bugfix Hugo version.
// This will be a patch increment on the previous Hugo version.
func (v HugoVersion) NextPatchLevel(level int) HugoVersion {
	return HugoVersion{Number: v.Number - 0.01, PatchLevel: level}
}

// CurrentHugoVersion represents the current build version.
// This should be the only one.
var CurrentHugoVersion = HugoVersion{
	Number:     0.22,
	PatchLevel: 0,
	Suffix:     "-DEV",
}

func hugoVersion(version float32, patchVersion int, suffix string) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2f.%d%s", version, patchVersion, suffix)
	}
	return fmt.Sprintf("%.2f%s", version, suffix)
}

// CompareVersion compares the given version string or number against the
// running Hugo version.
// It returns -1 if the given version is less than, 0 if equal and 1 if greater than
// the running version.
func CompareVersion(version interface{}) int {
	return compareVersions(CurrentHugoVersion.Number, CurrentHugoVersion.PatchLevel, version)
}

func compareVersions(inVersion float32, inPatchVersion int, in interface{}) int {
	switch d := in.(type) {
	case float64:
		return compareFloatVersions(inVersion, float32(d))
	case float32:
		return compareFloatVersions(inVersion, d)
	case int:
		return compareFloatVersions(inVersion, float32(d))
	case int32:
		return compareFloatVersions(inVersion, float32(d))
	case int64:
		return compareFloatVersions(inVersion, float32(d))
	default:
		s, err := cast.ToStringE(in)
		if err != nil {
			return -1
		}

		var (
			v float32
			p int
		)

		if strings.Count(s, ".") == 2 {
			li := strings.LastIndex(s, ".")
			p = cast.ToInt(s[li+1:])
			s = s[:li]
		}

		v = float32(cast.ToFloat64(s))

		if v == inVersion && p == inPatchVersion {
			return 0
		}

		if v < inVersion || (v == inVersion && p < inPatchVersion) {
			return -1
		}

		return 1
	}
}

func compareFloatVersions(version float32, v float32) int {
	if v == version {
		return 0
	}
	if v < version {
		return -1
	}
	return 1
}
