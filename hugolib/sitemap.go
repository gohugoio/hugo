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

package hugolib

import (
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// Sitemap configures the sitemap to be generated.
type Sitemap struct {
	ChangeFreq string
	Priority   float64
	Filename   string
}

func parseSitemap(input map[string]interface{}) Sitemap {
	sitemap := Sitemap{Priority: 0.5, Filename: "sitemap.xml"}

	for key, value := range input {
		switch key {
		case "changefreq":
			if value == "always" ||
			value == "hourly" ||
			value == "daily" ||
			value == "weekly" ||
			value == "monthly" ||
			value == "yearly" ||
			value == "never" {
				sitemap.ChangeFreq = cast.ToString(value)
				break
			}
			jww.WARN.Printf("value '%s' for sitemap.changefreq is invalid, accepted values are: always, hourly, daily, weekly, monthly, yearly, never\n", value)
		case "priority":
			priority := cast.ToFloat64(value)
			if priority >= 0 &&
			priority <= 1.0 {
				sitemap.Priority = priority
				break
			}
			jww.WARN.Printf("value '%s' for sitemap.priority is invalid, value should be between 0 and 1.0 and have a maximum of 1 decimal\n", value)
		case "filename":
			sitemap.Filename = cast.ToString(value)
		default:
			jww.WARN.Printf("Unknown Sitemap field: %s\n", key)
		}
	}

	return sitemap
}