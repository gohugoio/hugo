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

package tpl

import (
	"github.com/gohugoio/hugo/identity"
)

// Increments on breaking changes.
const TemplateVersion = 2

type Info interface {
	ParseInfo() ParseInfo

	// Identifies this template and its dependencies.
	identity.Provider
}

type InfoManager interface {
	ParseInfo() ParseInfo

	// Identifies and manages this template and its dependencies.
	identity.Manager
}

type defaultInfo struct {
	identity.Manager
	parseInfo ParseInfo
}

func NewInfo(id identity.Manager, parseInfo ParseInfo) Info {
	return &defaultInfo{
		Manager:   id,
		parseInfo: parseInfo,
	}
}

func (info *defaultInfo) ParseInfo() ParseInfo {
	return info.parseInfo
}

type ParseInfo struct {
	// Set for shortcode templates with any {{ .Inner }}
	IsInner bool

	// Set for partials with a return statement.
	HasReturn bool

	// Config extracted from template.
	Config ParseConfig
}

func (info ParseInfo) IsZero() bool {
	return info.Config.Version == 0
}

// Info holds some info extracted from a parsed template.
type Info1 struct {
}

type ParseConfig struct {
	Version int
}

var DefaultParseConfig = ParseConfig{
	Version: TemplateVersion,
}

var DefaultParseInfo = ParseInfo{
	Config: DefaultParseConfig,
}
