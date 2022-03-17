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

//go:build !nodeploy
// +build !nodeploy

package deploy

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/compute/v1"
)

// Invalidate all of the content in a Google Cloud CDN distribution.
func InvalidateGoogleCloudCDN(ctx context.Context, origin string) error {
	parts := strings.Split(origin, "/")
	if len(parts) != 2 {
		return fmt.Errorf("origin must be <project>/<origin>")
	}
	service, err := compute.NewService(ctx)
	if err != nil {
		return err
	}
	rule := &compute.CacheInvalidationRule{Path: "/*"}
	_, err = service.UrlMaps.InvalidateCache(parts[0], parts[1], rule).Context(ctx).Do()
	return err
}
