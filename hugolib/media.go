// Copyright © 2013 Steve Francia <spf@spf13.com>.
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

package hugolib

// An Image contains metadata for images + image sitemaps
// https://support.google.com/webmasters/answer/178636?hl=en
type Image struct {

	// The URL of the image. In some cases, the image URL may not be on the
	// same domain as your main site. This is fine, as long as both domains
	// are verified in Webmaster Tools. If, for example, you use a
	// content delivery network (CDN) to host your images, make sure that the
	// hosting site is verified in Webmaster Tools OR that you submit your
	// sitemap using robots.txt. In addition, make sure that your robots.txt
	// file doesn’t disallow the crawling of any content you want indexed.
	URL     string
	Title   string
	Caption string
	AltText string

	// The geographic location of the image. For example,
	// <image:geo_location>Limerick, Ireland</image:geo_location>.
	GeoLocation string

	// A URL to the license of the image.
	License string
}

// A Video contains metadata for videos + video sitemaps
// https://support.google.com/webmasters/answer/80471?hl=en
type Video struct {
	ThumbnailLoc         string
	Title                string
	Description          string
	ContentLoc           string
	PlayerLoc            string
	Duration             string
	ExpirationDate       string
	Rating               string
	ViewCount            string
	PublicationDate      string
	FamilyFriendly       string
	Restriction          string
	GalleryLoc           string
	Price                string
	RequiresSubscription string
	Uploader             string
	Live                 string
}
