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

package security

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/mitchellh/mapstructure"
)

const securityConfigKey = "security"

// DefaultConfig holds the default security policy.
var DefaultConfig = Config{
	Exec: Exec{
		Allow: MustNewWhitelist(
			"^(dart-)?sass(-embedded)?$", // sass, dart-sass, dart-sass-embedded.
			"^go$",                       // for Go Modules
			"^git$",                      // For Git info
			"^node$",                     // Used as the runtime for Node tools.
			"^postcss$",
			"^tailwindcss$",
		),
		// These have been tested to work with Hugo's external programs
		// on Windows, Linux and MacOS.
		OsEnv: MustNewWhitelist(`(?i)^((HTTPS?|NO)_PROXY|PATH(EXT)?|APPDATA|TE?MP|TERM|GO\w+|(XDG_CONFIG_)?HOME|USERPROFILE|SSH_AUTH_SOCK|DISPLAY|LANG|SYSTEMDRIVE|PROGRAMDATA)$`),
	},
	Funcs: Funcs{
		Getenv: MustNewWhitelist("^HUGO_", "^CI$"),
	},
	HTTP: HTTP{
		// Allow URLs whose host starts with a letter (the typical
		// "https://example.com" shape), deny anything that looks like
		// localhost, and deny URLs with userinfo ("http://user@...") to
		// foil the obvious SSRF bypass. Public IP literals are collateral
		// blocks; users who need them can override security.http.urls.
		URLs: MustNewWhitelist(
			`(?i)^https?://[a-z]`,
			`! (?i)localhost`,
			`! @`,
		),
		Methods: MustNewWhitelist("(?i)GET|POST"),
	},
	Node: Node{
		Permissions: NodePermissions{
			Disable:           false,
			AllowRead:         []string{"."},
			AllowWrite:        []string{},              // No write access by default.
			AllowAddons:       []string{"tailwindcss"}, // tailwindcss does not work without addon permissions.
			AllowChildProcess: []string{"tailwindcss"},
			AllowWorker:       []string{"tailwindcss"}, // tailwindcss needs worker access.
		},
	},
}

// Config is the top level security config.
// <docsmeta>{"name": "security", "description": "This section holds the top level security config.", "newIn": "0.91.0" }</docsmeta>
type Config struct {
	// Restricts access to os.Exec....
	// <docsmeta>{ "newIn": "0.91.0" }</docsmeta>
	Exec Exec `json:"exec"`

	// Restricts access to certain template funcs.
	Funcs Funcs `json:"funcs"`

	// Restricts access to resources.GetRemote, getJSON, getCSV.
	HTTP HTTP `json:"http"`

	// Node holds Node.js security settings.
	Node Node `json:"node"`

	// Allow inline shortcodes
	EnableInlineShortcodes bool `json:"enableInlineShortcodes"`
}

// Exec holds os/exec policies.
type Exec struct {
	Allow Whitelist `json:"allow"`
	OsEnv Whitelist `json:"osEnv"`
}

// Funcs holds template funcs policies.
type Funcs struct {
	// OS env keys allowed to query in os.Getenv.
	Getenv Whitelist `json:"getenv"`
}

type HTTP struct {
	// URLs to allow in remote HTTP (resources.Get, getJSON, getCSV).
	URLs Whitelist `json:"urls"`

	// HTTP methods to allow.
	Methods Whitelist `json:"methods"`

	// Media types where the Content-Type in the response is used instead of resolving from the file content.
	MediaTypes Whitelist `json:"mediaTypes"`
}

// Node holds Node.js security settings.
type Node struct {
	// Permissions configures Node's --permission flag for file system access control.
	Permissions NodePermissions `json:"permissions"`
}

// NodePermissions configures the Node.js permission model (--permission).
// Paths are relative to the working directory; "." means the working directory itself.
// Use "*" to allow all paths.
type NodePermissions struct {
	// Disable turns off the Node.js permission model entirely.
	Disable           bool     `json:"disable"`
	AllowRead         []string `json:"allowRead"`
	AllowWrite        []string `json:"allowWrite"`
	AllowAddons       []string `json:"allowAddons"`
	AllowChildProcess []string `json:"allowChildProcess"`
	AllowWorker       []string `json:"allowWorker"`
}

// IsEnabled reports whether the Node.js permission model is active.
func (p NodePermissions) IsEnabled() bool {
	return !p.Disable
}

// ToTOML converts c to TOML with [security] as the root.
func (c Config) ToTOML() string {
	sec := c.ToSecurityMap()

	var b bytes.Buffer

	if err := parser.InterfaceToConfig(sec, metadecoders.TOML, &b); err != nil {
		panic(err)
	}

	return strings.TrimSpace(b.String())
}

func (c Config) CheckAllowedExec(name string) error {
	if !c.Exec.Allow.Accept(name) {
		return &AccessDeniedError{
			name:     name,
			path:     "security.exec.allow",
			policies: c.ToTOML(),
		}
	}
	return nil
}

func (c Config) CheckAllowedGetEnv(name string) error {
	if !c.Funcs.Getenv.Accept(name) {
		return &AccessDeniedError{
			name:     name,
			path:     "security.funcs.getenv",
			policies: c.ToTOML(),
		}
	}
	return nil
}

func (c Config) CheckAllowedHTTPURL(url string) error {
	if !c.HTTP.URLs.Accept(url) {
		return &AccessDeniedError{
			name:     url,
			path:     "security.http.urls",
			policies: c.ToTOML(),
		}
	}
	return nil
}

func (c Config) CheckAllowedHTTPMethod(method string) error {
	if !c.HTTP.Methods.Accept(method) {
		return &AccessDeniedError{
			name:     method,
			path:     "security.http.method",
			policies: c.ToTOML(),
		}
	}
	return nil
}

// ToSecurityMap converts c to a map with 'security' as the root key.
func (c Config) ToSecurityMap() map[string]any {
	// Take it to JSON and back to get proper casing etc.
	asJson, err := json.Marshal(c)
	herrors.Must(err)
	m := make(map[string]any)
	herrors.Must(json.Unmarshal(asJson, &m))

	// Add the root
	sec := map[string]any{
		"security": m,
	}
	return sec
}

// DecodeConfig creates a privacy Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (Config, error) {
	sc := DefaultConfig
	// Deep copy slices to prevent mapstructure from mutating DefaultConfig.
	sc.Node.Permissions.AllowRead = slices.Clone(sc.Node.Permissions.AllowRead)
	sc.Node.Permissions.AllowWrite = slices.Clone(sc.Node.Permissions.AllowWrite)
	sc.Node.Permissions.AllowAddons = slices.Clone(sc.Node.Permissions.AllowAddons)
	sc.Node.Permissions.AllowChildProcess = slices.Clone(sc.Node.Permissions.AllowChildProcess)
	sc.Node.Permissions.AllowWorker = slices.Clone(sc.Node.Permissions.AllowWorker)
	if cfg.IsSet(securityConfigKey) {
		m := cfg.GetStringMap(securityConfigKey)
		dec, err := mapstructure.NewDecoder(
			&mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           &sc,
				DecodeHook:       stringSliceToWhitelistHook(),
			},
		)
		if err != nil {
			return sc, err
		}

		if err = dec.Decode(m); err != nil {
			return sc, err
		}
	}

	if !sc.EnableInlineShortcodes {
		// Legacy
		sc.EnableInlineShortcodes = cfg.GetBool("enableInlineShortcodes")
	}

	return sc, nil
}

func stringSliceToWhitelistHook() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if t != reflect.TypeFor[Whitelist]() {
			return data, nil
		}

		wl := types.ToStringSlicePreserveString(data)

		return NewWhitelist(wl...)
	}
}

// AccessDeniedError represents a security policy conflict.
type AccessDeniedError struct {
	path     string
	name     string
	policies string
}

func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("access denied: %q is not whitelisted in policy %q; the current security configuration is:\n\n%s\n\n", e.name, e.path, e.policies)
}

// IsAccessDenied reports whether err is an AccessDeniedError
func IsAccessDenied(err error) bool {
	var notFoundErr *AccessDeniedError
	return errors.As(err, &notFoundErr)
}
