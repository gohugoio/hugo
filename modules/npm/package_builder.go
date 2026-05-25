// Copyright 2020 The Hugo Authors. All rights reserved.
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

package npm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hmaps"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
)

const (
	dependenciesKey    = "dependencies"
	devDependenciesKey = "devDependencies"

	packageJSONName     = "package.json"
	packageMetaJSONName = "hugo_packagemeta.json"
)

var (
	workspacePackageJSON     = filepath.Join(files.FolderPackagesHugoAutoGen, packageJSONName)
	workspacePackageMetaJSON = filepath.Join(files.FolderPackagesHugoAutoGen, packageMetaJSONName)
)

func Pack(sourceFs, assetsWithDuplicatesPreservedFs afero.Fs, mods modules.Modules) error {
	b := &packageBuilder{
		devDependencies:         make(map[string]any),
		devDependenciesComments: make(map[string]any),
		dependencies:            make(map[string]any),
		dependenciesComments:    make(map[string]any),
	}

	workspacePath := filepath.ToSlash(files.FolderPackagesHugoAutoGen)

	// skip modules that shouldn't have their package files processed, either because they are the project module (handled separately)
	// or because their UsePackageJSON setting disables it.
	skipPackageJSON := buildSkipPackageJSON(mods)

	// 1. Read project deps: prefer package.hugo.json, fall back to package.json.
	var rootPkg map[string]any
	rootData, err := afero.ReadFile(sourceFs, packageJSONName)
	if err == nil {
		rootPkg = b.unmarshal(bytes.NewReader(rootData))
		if b.err != nil {
			return fmt.Errorf("npm pack: failed to parse package.json: %w", b.err)
		}
	}

	// Workspaces source: prefer package.hugo.json, fall back to package.json.
	var workspacesSource map[string]any
	hugoData, hugoErr := afero.ReadFile(sourceFs, files.FilenamePackageHugoJSON)
	if hugoErr == nil {
		hugoPkg := b.unmarshal(bytes.NewReader(hugoData))
		if b.err != nil {
			return fmt.Errorf("npm pack: failed to parse %s: %w", files.FilenamePackageHugoJSON, b.err)
		}
		b.addm("project", hugoPkg)
		workspacesSource = hugoPkg
	} else if rootPkg != nil {
		b.addm("project", rootPkg)
		workspacesSource = rootPkg
	}

	// 2. Read deps from referenced workspaces (always from package.json).
	for _, wsDir := range resolveProjectWorkspaces(sourceFs, workspacesSource, workspacePath) {
		wsFile := filepath.Join(wsDir, packageJSONName)
		wsData, err := afero.ReadFile(sourceFs, wsFile)
		if err != nil {
			continue
		}
		wsm := b.unmarshal(bytes.NewReader(wsData))
		if b.err != nil {
			return fmt.Errorf("npm pack: failed to parse %s: %w", wsFile, b.err)
		}
		b.addm("project", wsm)
	}

	// 3. Walk _jsconfig for module deps.
	// We use hugofs.Walkway (which uses ReadDir) instead of afero.Walk because
	// afero.Walk uses Readdirnames+Stat, which loses the per-module identity
	// when multiple modules mount a file (e.g. package.json) to the same virtual path.
	// Note that the order of files here is in the order of importance.
	type pkgEntry struct {
		info        hugofs.FileMetaInfo
		isHugoJSON  bool
		isRootLevel bool
	}
	var entries []pkgEntry
	modulesWithHugoJSON := make(map[string]bool)

	w := hugofs.NewWalkway(hugofs.WalkwayConfig{
		Fs:   assetsWithDuplicatesPreservedFs,
		Root: files.FolderJSConfig,
		WalkFn: func(ctx context.Context, path string, info hugofs.FileMetaInfo) error {
			if info.IsDir() {
				return nil
			}

			isPackageJSON := info.Name() == files.FilenamePackageJSON
			isHugoJSON := info.Name() == files.FilenamePackageHugoJSON

			if !isPackageJSON && !isHugoJSON {
				return nil
			}

			m := info.Meta()
			if skipPackageJSON[m.ModulePath()] {
				return nil
			}

			// package.hugo.json is only valid at module roots, not inside workspaces.
			isRootLevel := filepath.Dir(path) == files.FolderJSConfig
			if isHugoJSON && !isRootLevel {
				return nil
			}

			if isHugoJSON {
				modulesWithHugoJSON[m.ModulePath()] = true
			}

			entries = append(entries, pkgEntry{info: info, isHugoJSON: isHugoJSON, isRootLevel: isRootLevel})
			return nil
		},
	})
	if err := w.Walk(); err != nil {
		return err
	}

	// Process collected entries: for each module, prefer package.hugo.json
	// over package.json at the root level. Workspace package.json files are always processed.
	for _, e := range entries {
		m := e.info.Meta()
		// Skip root-level package.json if this module has package.hugo.json.
		if !e.isHugoJSON && e.isRootLevel && modulesWithHugoJSON[m.ModulePath()] {
			continue
		}

		f, err := m.Open()
		if err != nil {
			return fmt.Errorf("npm pack: failed to open package file: %w", err)
		}
		b.Add(m.ModulePath(), f)
		f.Close()
	}

	if b.Err() != nil {
		return fmt.Errorf("npm pack: failed to build: %w", b.Err())
	}

	// 4. Build the autogenerated workspace package.json.
	// Exclude deps already defined by the project itself — they don't
	// need to be duplicated in the workspace and it simplifies maintenance.
	moduleDeps := make(map[string]any)
	moduleDepsComments := make(map[string]any)
	for k, v := range b.dependencies {
		if b.dependenciesComments[k] != "project" {
			moduleDeps[k] = v
			moduleDepsComments[k] = b.dependenciesComments[k]
		}
	}
	moduleDevDeps := make(map[string]any)
	moduleDevDepsComments := make(map[string]any)
	for k, v := range b.devDependencies {
		if b.devDependenciesComments[k] != "project" {
			moduleDevDeps[k] = v
			moduleDevDepsComments[k] = b.devDependenciesComments[k]
		}
	}

	name := "project"
	rfi, err := sourceFs.Stat("")
	if err == nil {
		name = rfi.Name()
	}

	autoGenPkg := map[string]any{
		"name":             name,
		"version":          "0.1.0",
		dependenciesKey:    moduleDeps,
		devDependenciesKey: moduleDevDeps,
	}

	metaFile := packageMeta{
		Sum: PackageFilesSum(sourceFs, mods),
		DependencySources: dependencySources{
			Dependencies:    moduleDepsComments,
			DevDependencies: moduleDevDepsComments,
		},
	}

	if err := sourceFs.MkdirAll(files.FolderPackagesHugoAutoGen, 0o777); err != nil {
		return err
	}
	if err := writeJSON(sourceFs, workspacePackageJSON, autoGenPkg); err != nil {
		return err
	}
	if err := writeJSON(sourceFs, workspacePackageMetaJSON, metaFile); err != nil {
		return err
	}

	// 5. Ensure root package.json references the workspace.
	return ensureWorkspaceRef(sourceFs, workspacePath)
}

// ensureWorkspaceRef adds workspacePath to the "workspaces" array in root
// package.json with minimal formatting changes.
func ensureWorkspaceRef(fsys afero.Fs, workspacePath string) error {
	data, err := afero.ReadFile(fsys, packageJSONName)
	if err != nil {
		content := fmt.Sprintf("{\n  \"workspaces\": [\n    %q\n  ]\n}\n", workspacePath)
		return afero.WriteFile(fsys, packageJSONName, []byte(content), 0o666)
	}

	if runtime.GOOS == "windows" {
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	}

	// Parse to check if already present.
	var pkg map[string]any
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("npm pack: failed to parse package.json: %w", err)
	}

	indent := detectIndent(data)
	quoted := fmt.Sprintf("%q", workspacePath)

	wsVal, hasWS := pkg["workspaces"]

	switch v := wsVal.(type) {
	case []interface{}:
		// Array form: ["pkg-a", "pkg-b", ...]
		if containsString(toStringSlice(v), workspacePath) {
			if runtime.GOOS == "windows" {
				return afero.WriteFile(fsys, packageJSONName, data, 0o666)
			}
			return nil
		}
		// Fall through to byte-based insertion into the existing array below.
	case map[string]interface{}:
		// Object form: { "workspaces": { "packages": [...] } }
		packagesVal, ok := v["packages"]
		if !ok {
			return fmt.Errorf("npm pack: unsupported workspaces object; missing \"packages\" field")
		}
		packagesSlice := toStringSlice(packagesVal)
		if containsString(packagesSlice, workspacePath) {
			if runtime.GOOS == "windows" {
				return afero.WriteFile(fsys, packageJSONName, data, 0o666)
			}
			return nil
		}

		// Append the new workspace path to the packages slice.
		newPkgs := make([]interface{}, 0, len(packagesSlice)+1)
		for _, s := range packagesSlice {
			newPkgs = append(newPkgs, s)
		}
		newPkgs = append(newPkgs, workspacePath)
		v["packages"] = newPkgs

		updated, err := json.MarshalIndent(pkg, "", indent)
		if err != nil {
			return fmt.Errorf("npm pack: failed to marshal package.json with updated workspaces: %w", err)
		}
		updated = append(updated, '\n')
		return afero.WriteFile(fsys, packageJSONName, updated, 0o666)
	case nil:
		// Treat explicit null as "no workspaces"; handled below as missing key.
	default:
		return fmt.Errorf("npm pack: unsupported workspaces type %T", v)
	}

	// Try adding to existing workspaces array when present.
	if hasWS && wsVal != nil {
		if wsIdx := bytes.Index(data, []byte(`"workspaces"`)); wsIdx >= 0 {
			rest := data[wsIdx:]
			if bracketOpen := bytes.IndexByte(rest, '['); bracketOpen >= 0 {
				if bracketClose := bytes.IndexByte(rest[bracketOpen:], ']'); bracketClose >= 0 {
					pos := wsIdx + bracketOpen + bracketClose
					arrayContent := bytes.TrimSpace(data[wsIdx+bracketOpen+1 : pos])
					var insertion string
					if len(arrayContent) == 0 {
						insertion = "\n" + indent + indent + quoted + "\n" + indent
					} else {
						insertion = ",\n" + indent + indent + quoted + "\n" + indent
					}
					result := make([]byte, 0, len(data)+len(insertion))
					result = append(result, data[:pos]...)
					result = append(result, insertion...)
					result = append(result, data[pos:]...)
					return afero.WriteFile(fsys, packageJSONName, result, 0o666)
				}
			}
		}

		// We know a workspaces array exists but could not locate it reliably in the raw data.
		return fmt.Errorf("npm pack: could not locate existing workspaces array for insertion")
	}

	// No workspaces key — add before closing brace.
	lastBrace := bytes.LastIndexByte(data, '}')
	if lastBrace < 0 {
		return fmt.Errorf("npm pack: malformed package.json")
	}

	before := bytes.TrimRight(data[:lastBrace], " \t\n\r")
	needComma := len(before) > 0 && before[len(before)-1] != '{' && before[len(before)-1] != ','

	var insertion string
	if needComma {
		insertion = ",\n" + indent + `"workspaces": [` + "\n" + indent + indent + quoted + "\n" + indent + "]\n"
	} else {
		insertion = indent + `"workspaces": [` + "\n" + indent + indent + quoted + "\n" + indent + "]\n"
	}

	result := make([]byte, 0, len(data)+len(insertion))
	result = append(result, before...)
	result = append(result, insertion...)
	result = append(result, data[lastBrace:]...)
	return afero.WriteFile(fsys, packageJSONName, result, 0o666)
}

func detectIndent(data []byte) string {
	for _, line := range bytes.Split(data, []byte("\n")) {
		trimmed := bytes.TrimLeft(line, " \t")
		if len(trimmed) < len(line) && len(trimmed) > 0 && trimmed[0] == '"' {
			return string(line[:len(line)-len(trimmed)])
		}
	}
	return "  "
}

func writeJSON(fs afero.Fs, filename string, v any) error {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", strings.Repeat(" ", 2))
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("npm pack: failed to marshal JSON: %w", err)
	}
	return afero.WriteFile(fs, filename, buf.Bytes(), 0o666)
}

func toStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	if s, ok := v.([]any); ok {
		var out []string
		for _, item := range s {
			if str, ok := item.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}

func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// resolveProjectWorkspaces resolves workspace patterns from the project's
// package source, skipping the hugoautogen workspace.
func resolveProjectWorkspaces(sourceFs afero.Fs, workspacesSource map[string]any, skipPath string) []string {
	if workspacesSource == nil {
		return nil
	}
	var dirs []string
	for _, ws := range toStringSlice(workspacesSource["workspaces"]) {
		for _, wsDir := range modules.ResolveWorkspacePattern(sourceFs, "", ws) {
			if filepath.ToSlash(wsDir) != skipPath {
				dirs = append(dirs, wsDir)
			}
		}
	}
	return dirs
}

type packageBuilder struct {
	err error

	devDependencies         map[string]any
	devDependenciesComments map[string]any
	dependencies            map[string]any
	dependenciesComments    map[string]any
}

func (b *packageBuilder) Add(source string, r io.Reader) *packageBuilder {
	if b.err != nil {
		return b
	}

	m := b.unmarshal(r)
	if b.err != nil {
		return b
	}

	b.addm(source, m)

	return b
}

func (b *packageBuilder) addm(source string, m map[string]any) {
	if source == "" {
		source = "project"
	}

	// First version for a given dependency wins.
	// Packages added by order of import (project, module1, module2...),
	// so the project has control over versions.
	if devDeps, found := m[devDependenciesKey]; found {
		mm := hmaps.ToStringMapString(devDeps)
		for k, v := range mm {
			if _, added := b.devDependencies[k]; !added {
				b.devDependencies[k] = v
				b.devDependenciesComments[k] = source
			}
		}
	}

	if deps, found := m[dependenciesKey]; found {
		mm := hmaps.ToStringMapString(deps)
		for k, v := range mm {
			if _, added := b.dependencies[k]; !added {
				b.dependencies[k] = v
				b.dependenciesComments[k] = source
			}
		}
	}
}

func (b *packageBuilder) unmarshal(r io.Reader) map[string]any {
	m := make(map[string]any)
	err := json.Unmarshal(helpers.ReaderToBytes(r), &m)
	if err != nil {
		b.err = err
	}
	return m
}

func (b *packageBuilder) Err() error {
	return b.err
}

// buildSkipPackageJSON determines which modules should NOT have their package
// files processed. The project module is always skipped (handled separately).
// For other modules, the behavior is controlled by the UsePackageJSON import
// setting: "auto" (default) reads package files when a Hugo config file or
// package.hugo.json is present; "always" always reads; "never" never reads.
func buildSkipPackageJSON(mods modules.Modules) map[string]bool {
	skip := make(map[string]bool)
	for _, m := range mods {
		if m.Owner() == nil {
			skip[m.Path()] = true
			continue
		}
		if !usePackageJSON(m) {
			skip[m.Path()] = true
		}
	}
	return skip
}

// usePackageJSON checks the import config for this module and applies the
// UsePackageJSON setting. For "auto", it checks for Hugo config files or
// package.hugo.json in the module root.
func usePackageJSON(m modules.Module) bool {
	setting := findImportSetting(m)
	switch setting {
	case modules.UsePackageJSONAlways:
		return true
	case modules.UsePackageJSONNever:
		return false
	default:
		// "auto": use if Hugo config file or package.hugo.json is present.
		if len(m.ConfigFilenames()) > 0 {
			return true
		}
		if m.Dir() != "" {
			if _, err := os.Stat(filepath.Join(m.Dir(), files.FilenamePackageHugoJSON)); err == nil {
				return true
			}
		}
		return false
	}
}

func findImportSetting(m modules.Module) string {
	if m.Owner() == nil {
		return modules.UsePackageJSONAuto
	}
	for _, imp := range m.Owner().Config().Imports {
		if imp.Path == m.Path() {
			if imp.UsePackageJSON != "" {
				return imp.UsePackageJSON
			}
			return modules.UsePackageJSONAuto
		}
	}
	return modules.UsePackageJSONAuto
}

type packageMeta struct {
	// Sum is a hash of all package files that feed into the npm pack output.
	Sum string `json:"sum"`

	DependencySources dependencySources `json:"dependencySources"`
}

type dependencySources struct {
	Dependencies    map[string]any `json:"dependencies"`
	DevDependencies map[string]any `json:"devDependencies"`
}

// PackageFilesSum hashes the package files that Pack would use,
// using the same file selection logic as Pack.
func PackageFilesSum(sourceFs afero.Fs, mods modules.Modules) string {
	h := hashing.XxHasher()
	defer h.Close()

	var w io.Writer = h
	if runtime.GOOS == "windows" {
		w = &crlfReplacer{w: h}
	}

	copyFile := func(fsys afero.Fs, name string) {
		f, err := fsys.Open(name)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}
	copyOsFile := func(name string) {
		f, err := os.Open(name)
		if err != nil {
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}

	// Project level: prefer package.hugo.json, fall back to package.json.
	// We need to parse whichever file we pick to discover workspaces.
	var workspacesSource map[string]any
	if data, err := afero.ReadFile(sourceFs, files.FilenamePackageHugoJSON); err == nil {
		w.Write(data)
		var m map[string]any
		if err := json.Unmarshal(data, &m); err == nil {
			workspacesSource = m
		}
	} else if data, err := afero.ReadFile(sourceFs, packageJSONName); err == nil {
		w.Write(data)
		var m map[string]any
		if err := json.Unmarshal(data, &m); err == nil {
			workspacesSource = m
		}
	}

	// Workspace package.json files (skipping hugoautogen).
	workspacePath := filepath.ToSlash(files.FolderPackagesHugoAutoGen)
	for _, wsDir := range resolveProjectWorkspaces(sourceFs, workspacesSource, workspacePath) {
		copyFile(sourceFs, filepath.Join(wsDir, packageJSONName))
	}

	// Module package files: prefer package.hugo.json, fall back to package.json.
	skip := buildSkipPackageJSON(mods)
	modulesWithHugoJSON := make(map[string]bool)
	for _, m := range mods {
		if skip[m.Path()] || m.Dir() == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(m.Dir(), files.FilenamePackageHugoJSON)); err == nil {
			modulesWithHugoJSON[m.Path()] = true
		}
	}
	for _, m := range mods {
		if skip[m.Path()] || m.Dir() == "" {
			continue
		}
		if modulesWithHugoJSON[m.Path()] {
			copyOsFile(filepath.Join(m.Dir(), files.FilenamePackageHugoJSON))
		} else {
			copyOsFile(filepath.Join(m.Dir(), packageJSONName))
		}
	}

	return fmt.Sprintf("%x", h.Sum64())
}

// crlfReplacer wraps a writer and strips \r bytes.
type crlfReplacer struct {
	w io.Writer
}

func (c *crlfReplacer) Write(p []byte) (int, error) {
	n := len(p)
	_, err := c.w.Write(bytes.ReplaceAll(p, []byte("\r\n"), []byte("\n")))
	return n, err
}

// NpmPackNeedsUpdate checks if the npm pack output is stale by comparing
// the stored hash against the current package files.
func NpmPackNeedsUpdate(sourceFs afero.Fs, mods modules.Modules) bool {
	data, err := afero.ReadFile(sourceFs, workspacePackageMetaJSON)
	if err != nil {
		// No meta file means npm pack hasn't been run yet.
		return false
	}

	// We only need the sum for this check.
	meta := struct {
		Sum string `json:"sum"`
	}{}
	if err := json.Unmarshal(data, &meta); err != nil || meta.Sum == "" {
		return true
	}
	return meta.Sum != PackageFilesSum(sourceFs, mods)
}
