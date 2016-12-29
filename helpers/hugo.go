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
)

// HugoVersionNumber represents the current build version.
// This should be the only one
const (
	// Major and minor version.
	HugoVersionNumber = 0.18

	// Increment this for bug releases
	HugoPatchVersion = 1
)

// HugoVersionSuffix is the suffix used in the Hugo version string.
// It will be blank for release versions.
// const HugoVersionSuffix = "-DEV" // use this when not doing a release
const HugoVersionSuffix = "" // use this line when doing a release

// HugoVersion returns the current Hugo version. It will include
// a suffix, typically '-DEV', if it's development version.
func HugoVersion() string {
	return hugoVersion(HugoVersionNumber, HugoPatchVersion, HugoVersionSuffix)
}

// HugoReleaseVersion is same as HugoVersion, but no suffix.
func HugoReleaseVersion() string {
	return hugoVersionNoSuffix(HugoVersionNumber, HugoPatchVersion)
}

// NextHugoReleaseVersion returns the next Hugo release version.
func NextHugoReleaseVersion() string {
	return hugoVersionNoSuffix(HugoVersionNumber+0.01, 0)
}

func hugoVersion(version float32, patchVersion int, suffix string) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2g.%d%s", version, patchVersion, suffix)
	}
	return fmt.Sprintf("%.2g%s", version, suffix)
}

func hugoVersionNoSuffix(version float32, patchVersion int) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2g.%d", version, patchVersion)
	}
	return fmt.Sprintf("%.2g", version)
}
