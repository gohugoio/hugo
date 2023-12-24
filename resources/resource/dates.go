// Copyright 2019 The Hugo Authors. All rights reserved.
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

package resource

import (
	"time"

	"github.com/gohugoio/hugo/common/htime"
)

var _ Dated = Dates{}

// Dated wraps a "dated resource". These are the 4 dates that makes
// the date logic in Hugo.
type Dated interface {
	// Date returns the date of the resource.
	Date() time.Time

	// Lastmod returns the last modification date of the resource.
	Lastmod() time.Time

	// PublishDate returns the publish date of the resource.
	PublishDate() time.Time

	// ExpiryDate returns the expiration date of the resource.
	ExpiryDate() time.Time
}

// Dates holds the 4 Hugo dates.
type Dates struct {
	FDate        time.Time
	FLastmod     time.Time
	FPublishDate time.Time
	FExpiryDate  time.Time
}

func (d *Dates) IsDateOrLastModAfter(in Dated) bool {
	return d.Date().After(in.Date()) || d.Lastmod().After(in.Lastmod())
}

func (d *Dates) UpdateDateAndLastmodIfAfter(in Dated) {
	if in.Date().After(d.Date()) {
		d.FDate = in.Date()
	}
	if in.Lastmod().After(d.Lastmod()) {
		d.FLastmod = in.Lastmod()
	}
}

// IsFuture returns whether the argument represents the future.
func IsFuture(d Dated) bool {
	if d.PublishDate().IsZero() {
		return false
	}

	return d.PublishDate().After(htime.Now())
}

// IsExpired returns whether the argument is expired.
func IsExpired(d Dated) bool {
	if d.ExpiryDate().IsZero() {
		return false
	}
	return d.ExpiryDate().Before(htime.Now())
}

// IsZeroDates returns true if all of the dates are zero.
func IsZeroDates(d Dated) bool {
	return d.Date().IsZero() && d.Lastmod().IsZero() && d.ExpiryDate().IsZero() && d.PublishDate().IsZero()
}

func (p Dates) Date() time.Time {
	return p.FDate
}

func (p Dates) Lastmod() time.Time {
	return p.FLastmod
}

func (p Dates) PublishDate() time.Time {
	return p.FPublishDate
}

func (p Dates) ExpiryDate() time.Time {
	return p.FExpiryDate
}
