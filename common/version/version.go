// Copyright 2025 The Hugo Authors. All rights reserved.
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

package version

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
	return CompareVersions(h, other)
}

// VersionString represents a Hugo version string.
type VersionString string

func (h VersionString) String() string {
	return string(h)
}

// Compare implements the compare.Comparer interface.
func (h VersionString) Compare(other any) int {
	return CompareVersions(h.Version(), other)
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

// ParseVersion parses a version string.
func ParseVersion(s string) (Version, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	var vv Version
	hyphen := strings.Index(s, "-")
	if hyphen > 0 {
		suffix := s[hyphen:]
		if len(suffix) > 1 {
			if suffix[0] == '-' {
				suffix = suffix[1:]
			}
			if len(suffix) > 0 {
				vv.Suffix = suffix
				s = s[:hyphen]
			}
		}
		vv.Suffix = suffix
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

func version(major, minor, patch int, suffix string) string {
	if suffix != "" {
		if suffix[0] != '-' {
			suffix = "-" + suffix
		}
	}
	if patch > 0 || minor > 53 {
		return fmt.Sprintf("%d.%d.%d%s", major, minor, patch, suffix)
	}
	return fmt.Sprintf("%d.%d%s", major, minor, suffix)
}

// CompareVersion compares v1 with v2.
// It returns -1 if the v2 is less than, 0 if equal and 1 if greater than
// v1.
func CompareVersions(v1 Version, v2 any) int {
	var c int
	switch d := v2.(type) {
	case float64:
		c = compareFloatWithVersion(d, v1)
	case float32:
		c = compareFloatWithVersion(float64(d), v1)
	case int:
		c = compareFloatWithVersion(float64(d), v1)
	case int32:
		c = compareFloatWithVersion(float64(d), v1)
	case int64:
		c = compareFloatWithVersion(float64(d), v1)
	case Version:
		if d.Major == v1.Major && d.Minor == v1.Minor && d.PatchLevel == v1.PatchLevel {
			return strings.Compare(v1.Suffix, d.Suffix)
		}
		if d.Major > v1.Major {
			return 1
		} else if d.Major < v1.Major {
			return -1
		}
		if d.Minor > v1.Minor {
			return 1
		} else if d.Minor < v1.Minor {
			return -1
		}
		if d.PatchLevel > v1.PatchLevel {
			return 1
		} else if d.PatchLevel < v1.PatchLevel {
			return -1
		}
	default:
		s, err := cast.ToStringE(v2)
		if err != nil {
			return -1
		}

		v, err := ParseVersion(s)
		if err != nil {
			return -1
		}
		return v1.Compare(v)

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
