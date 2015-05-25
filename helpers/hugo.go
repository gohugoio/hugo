// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/parser"
)

// this should be the only one
const hugoVersionMain = 0.14
const hugoVersionSuffix = "" // blank this when doing a release

// HugoVersion returns the current Hugo version. It will include
// a suffix, typically '-DEV', if it's development version.
func HugoVersion() string {
	return hugoVersion(hugoVersionMain, hugoVersionSuffix)
}

// HugoReleaseVersion is same as HugoVersion, but no suffix.
func HugoReleaseVersion() string {
	return hugoVersionNoSuffix(hugoVersionMain)
}

// NextHugoReleaseVersion returns the next Hugo release version.
func NextHugoReleaseVersion() string {
	return hugoVersionNoSuffix(hugoVersionMain + 0.01)
}

func hugoVersion(version float32, suffix string) string {
	return fmt.Sprintf("%.2g%s", version, suffix)
}

func hugoVersionNoSuffix(version float32) string {
	return fmt.Sprintf("%.2g", version)
}

// IsThemeVsHugoVersionMismatch returns whether the current Hugo version is < theme's min_version
func IsThemeVsHugoVersionMismatch() (mismatch bool, requiredMinVersion string) {
	if !ThemeSet() {
		return
	}

	themeDir, err := getThemeDirPath("")

	if err != nil {
		return
	}

	fs := hugofs.SourceFs
	path := filepath.Join(themeDir, "theme.toml")

	exists, err := Exists(path, fs)

	if err != nil || !exists {
		return
	}

	f, err := fs.Open(path)

	if err != nil {
		return
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	if err != nil {
		return
	}

	c, err := parser.HandleTOMLMetaData(b)

	if err != nil {
		return
	}

	config := c.(map[string]interface{})

	if minVersion, ok := config["min_version"]; ok {
		switch minVersion.(type) {
		case float32:
			return hugoVersionMain < minVersion.(float32), fmt.Sprint(minVersion)
		case float64:
			return hugoVersionMain < minVersion.(float64), fmt.Sprint(minVersion)
		default:
			return
		}

	}

	return
}
