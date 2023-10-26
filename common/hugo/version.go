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
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/compare"
	"github.com/spf13/cast"
)

// Version represents the Hugo build version.
type Version struct {
	Major int

	Minor int

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
	return version(v.Major, v.Minor, v.PatchLevel, v.Suffix)
}

// Version returns the Hugo version.
func (v Version) Version() VersionString {
	return VersionString(v.String())
}

// Compare implements the compare.Comparer interface.
func (h Version) Compare(other any) int {
	return compareVersions(h, other)
}

// VersionString represents a Hugo version string.
type VersionString string

func (h VersionString) String() string {
	return string(h)
}

// Compare implements the compare.Comparer interface.
func (h VersionString) Compare(other any) int {
	return compareVersions(h.Version(), other)
}

func (h VersionString) Version() Version {
	return MustParseVersion(h.String())
}

// Eq implements the compare.Eqer interface.
func (h VersionString) Eq(other any) bool {
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

	vv.Major, vv.Minor, vv.PatchLevel = parseVersion(s)

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
	return Version{Major: v.Major, Minor: v.Minor + 1}
}

// Prev returns the previous Hugo release version.
func (v Version) Prev() Version {
	return Version{Major: v.Major, Minor: v.Minor - 1}
}

// NextPatchLevel returns the next patch/bugfix Hugo version.
// This will be a patch increment on the previous Hugo version.
func (v Version) NextPatchLevel(level int) Version {
	prev := v.Prev()
	prev.PatchLevel = level
	return prev
}

// BuildVersionString creates a version string. This is what you see when
// running "hugo version".
func BuildVersionString() string {
	// program := "Hugo Static Site Generator"
	program := "hugo"

	version := "v" + CurrentVersion.String()

	bi := getBuildInfo()
	if bi == nil {
		return version
	}
	if bi.Revision != "" {
		version += "-" + bi.Revision
	}
	if IsExtended {
		version += "+extended"
	}

	osArch := bi.GoOS + "/" + bi.GoArch

	date := bi.RevisionTime
	if date == "" {
		// Accept vendor-specified build date if .git/ is unavailable.
		date = buildDate
	}
	if date == "" {
		date = "unknown"
	}

	versionString := fmt.Sprintf("%s %s %s BuildDate=%s",
		program, version, osArch, date)

	if vendorInfo != "" {
		versionString += " VendorInfo=" + vendorInfo
	}

	return versionString
}

func version(major, minor, patch int, suffix string) string {
	if patch > 0 || minor > 53 {
		return fmt.Sprintf("%d.%d.%d%s", major, minor, patch, suffix)
	}
	return fmt.Sprintf("%d.%d%s", major, minor, suffix)
}

// CompareVersion compares the given version string or number against the
// running Hugo version.
// It returns -1 if the given version is less than, 0 if equal and 1 if greater than
// the running version.
func CompareVersion(version any) int {
	return compareVersions(CurrentVersion, version)
}

func compareVersions(inVersion Version, in any) int {
	var c int
	switch d := in.(type) {
	case float64:
		c = compareFloatWithVersion(d, inVersion)
	case float32:
		c = compareFloatWithVersion(float64(d), inVersion)
	case int:
		c = compareFloatWithVersion(float64(d), inVersion)
	case int32:
		c = compareFloatWithVersion(float64(d), inVersion)
	case int64:
		c = compareFloatWithVersion(float64(d), inVersion)
	case Version:
		if d.Major == inVersion.Major && d.Minor == inVersion.Minor && d.PatchLevel == inVersion.PatchLevel {
			return strings.Compare(inVersion.Suffix, d.Suffix)
		}
		if d.Major > inVersion.Major {
			return 1
		} else if d.Major < inVersion.Major {
			return -1
		}
		if d.Minor > inVersion.Minor {
			return 1
		} else if d.Minor < inVersion.Minor {
			return -1
		}
		if d.PatchLevel > inVersion.PatchLevel {
			return 1
		} else if d.PatchLevel < inVersion.PatchLevel {
			return -1
		}
	default:
		s, err := cast.ToStringE(in)
		if err != nil {
			return -1
		}

		v, err := ParseVersion(s)
		if err != nil {
			return -1
		}
		return inVersion.Compare(v)

	}

	return c
}

func parseVersion(s string) (int, int, int) {
	var major, minor, patch int
	parts := strings.Split(s, ".")
	if len(parts) > 0 {
		major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) > 1 {
		minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		patch, _ = strconv.Atoi(parts[2])
	}

	return major, minor, patch
}

// compareFloatWithVersion compares v1 with v2.
// It returns -1 if v1 is less than v2, 0 if v1 is equal to v2 and 1 if v1 is greater than v2.
func compareFloatWithVersion(v1 float64, v2 Version) int {
	mf, minf := math.Modf(v1)
	v1maj := int(mf)
	v1min := int(minf * 100)

	if v2.Major == v1maj && v2.Minor == v1min {
		return 0
	}

	if v1maj > v2.Major {
		return 1
	}

	if v1maj < v2.Major {
		return -1
	}

	if v1min > v2.Minor {
		return 1
	}

	return -1
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
