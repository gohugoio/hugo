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

package config

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
)

func TestGetStringSlicePreserveString(t *testing.T) {
	c := qt.New(t)
	cfg := viper.New()

	s := "This is a string"
	sSlice := []string{"This", "is", "a", "slice"}

	cfg.Set("s1", s)
	cfg.Set("s2", sSlice)

	c.Assert(GetStringSlicePreserveString(cfg, "s1"), qt.DeepEquals, []string{s})
	c.Assert(GetStringSlicePreserveString(cfg, "s2"), qt.DeepEquals, sSlice)
	c.Assert(GetStringSlicePreserveString(cfg, "s3"), qt.IsNil)
}
