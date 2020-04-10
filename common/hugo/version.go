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

package hugo

import (
	"fmt"
	"io"

	"runtime"
	"strings"

	"github.com/gohugoio/hugo/compare"
	"github.com/spf13/cast"
)

// Version represents the Hugo build version.
type Version struct {
	// Major and minor version.
	Number float32

	// Increment this for bug releases
	PatchLevel int

	// HugoVersionSuffix is the suffix used in the Hugo version string.
	// It will be blank for release versions.
	Suffix string
}

var (
	_ compare.Eqer     = (*VersionString)(nil)
	_ compare.Comparer = (*VersionString)(nil)
)

func (v Version) String() string {
	return version(v.Number, v.PatchLevel, v.Suffix)
}

// Version returns the Hugo version.
func (v Version) Version() VersionString {
	return VersionString(v.String())
}

// VersionString represents a Hugo version string.
type VersionString string

func (h VersionString) String() string {
	return string(h)
}

// Compare implements the compare.Comparer interface.
func (h VersionString) Compare(other interface{}) int {
	v := MustParseVersion(h.String())
	return compareVersionsWithSuffix(v.Number, v.PatchLevel, v.Suffix, other)
}

// Eq implements the compare.Eqer interface.
func (h VersionString) Eq(other interface{}) bool {
	s, err := cast.ToStringE(other)
	if err != nil {
		return false
	}
	return s == h.String()
}

var versionSuffixes = []string{"-test", "-DEV"}

// ParseVersion parses a version string.
func ParseVersion(s string) (Version, error) {
	var vv Version
	for _, suffix := range versionSuffixes {
		if strings.HasSuffix(s, suffix) {
			vv.Suffix = suffix
			s = strings.TrimSuffix(s, suffix)
		}
	}

	v, p := parseVersion(s)

	vv.Number = v
	vv.PatchLevel = p

	return vv, nil
}

// MustParseVersion parses a version string
// and panics if any error occurs.
func MustParseVersion(s string) Version {
	vv, err := ParseVersion(s)
	if err != nil {
		panic(err)
	}
	return vv
}

// ReleaseVersion represents the release version.
func (v Version) ReleaseVersion() Version {
	v.Suffix = ""
	return v
}

// Next returns the next Hugo release version.
func (v Version) Next() Version {
	return Version{Number: v.Number + 0.01}
}

// Prev returns the previous Hugo release version.
func (v Version) Prev() Version {
	return Version{Number: v.Number - 0.01}
}

// NextPatchLevel returns the next patch/bugfix Hugo version.
// This will be a patch increment on the previous Hugo version.
func (v Version) NextPatchLevel(level int) Version {
	return Version{Number: v.Number - 0.01, PatchLevel: level}
}

// BuildVersionString creates a version string. This is what you see when
// running "hugo version".
func BuildVersionString() string {
	program := "Hugo Static Site Generator"

	version := "v" + CurrentVersion.String()
	if commitHash != "" {
		version += "-" + strings.ToUpper(commitHash)
	}
	if IsExtended {
		version += "/extended"
	}

	osArch := runtime.GOOS + "/" + runtime.GOARCH

	date := buildDate
	if date == "" {
		date = "unknown"
	}

	return fmt.Sprintf("%s %s %s BuildDate: %s", program, version, osArch, date)

}

func version(version float32, patchVersion int, suffix string) string {
	if patchVersion > 0 || version > 0.53 {
		return fmt.Sprintf("%.2f.%d%s", version, patchVersion, suffix)
	}
	return fmt.Sprintf("%.2f%s", version, suffix)
}

// CompareVersion compares the given version string or number against the
// running Hugo version.
// It returns -1 if the given version is less than, 0 if equal and 1 if greater than
// the running version.
func CompareVersion(version interface{}) int {
	return compareVersionsWithSuffix(CurrentVersion.Number, CurrentVersion.PatchLevel, CurrentVersion.Suffix, version)
}

func compareVersions(inVersion float32, inPatchVersion int, in interface{}) int {
	return compareVersionsWithSuffix(inVersion, inPatchVersion, "", in)
}

func compareVersionsWithSuffix(inVersion float32, inPatchVersion int, suffix string, in interface{}) int {
	var c int
	switch d := in.(type) {
	case float64:
		c = compareFloatVersions(inVersion, float32(d))
	case float32:
		c = compareFloatVersions(inVersion, d)
	case int:
		c = compareFloatVersions(inVersion, float32(d))
	case int32:
		c = compareFloatVersions(inVersion, float32(d))
	case int64:
		c = compareFloatVersions(inVersion, float32(d))
	default:
		s, err := cast.ToStringE(in)
		if err != nil {
			return -1
		}

		v, err := ParseVersion(s)
		if err != nil {
			return -1
		}

		if v.Number == inVersion && v.PatchLevel == inPatchVersion {
			return strings.Compare(suffix, v.Suffix)
		}

		if v.Number < inVersion || (v.Number == inVersion && v.PatchLevel < inPatchVersion) {
			return -1
		}

		return 1
	}

	if c == 0 && suffix != "" {
		return 1
	}

	return c
}

func parseVersion(s string) (float32, int) {
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

	return v, p
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

func GoMinorVersion() int {
	return goMinorVersion(runtime.Version())
}

func goMinorVersion(version string) int {
	if strings.HasPrefix(version, "devel") {
		return 9999 // magic
	}
	var major, minor int
	var trailing string
	n, err := fmt.Sscanf(version, "go%d.%d%s", &major, &minor, &trailing)
	if n == 2 && err == io.EOF {
		// Means there were no trailing characters (i.e., not an alpha/beta)
		err = nil
	}
	if err != nil {
		return 0
	}
	return minor
}
