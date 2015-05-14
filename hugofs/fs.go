// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
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

package hugofs

import "github.com/spf13/afero"

var SourceFs afero.Fs = new(afero.OsFs)
var DestinationFS afero.Fs = new(afero.OsFs)
var OsFs afero.Fs = new(afero.OsFs)

//var DestinationFS afero.Fs = new(afero.MemMapFs)
