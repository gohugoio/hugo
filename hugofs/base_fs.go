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

package hugofs

import (
	"github.com/spf13/afero"
)

// BaseFs contains the core base filesystems used by Hugo. The name "base" is used
// to underline that even if they can be composites, they all have a base path set to a specific
// resource folder, e.g "/my-project/content". So, no absolute filenames needed.
type BaseFs struct {
	// The filesystem used to capture content. This can be a composite and
	// language aware file system.
	ContentFs afero.Fs

	// The filesystem used to store resources (processed images etc.).
	// This usually maps to /my-project/resources.
	ResourcesFs afero.Fs

	// The filesystem used to publish the rendered site.
	// This usually maps to /my-project/public.
	PublishFs afero.Fs
}
