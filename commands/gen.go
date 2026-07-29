// Copyright 2024 The Hugo Authors. All rights reserved.
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

package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/bep/simplecobra"
	"github.com/goccy/go-yaml"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/hugolib/roles"
	"github.com/gohugoio/hugo/hugolib/segments"
	"github.com/gohugoio/hugo/hugolib/versions"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/asciidocext/asciidocext_config"
	"github.com/gohugoio/hugo/markup/goldmark/goldmark_config"
	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/navigation"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Schema generation statistics
type schemaStats struct {
	schemasGenerated   int
	totalProperties    int
	documentationLinks int
}

func newGenCommand() *genCommand {
	var (
		// Flags.
		gendocdir string
		genmandir string

		// Chroma flags.
		style                  string
		mode                   string
		modeSelector           bool
		highlightStyle         string
		lineNumbersInlineStyle string
		lineNumbersTableStyle  string
		omitEmpty              bool
		omitClassComments      bool
	)

	newChromaStyles := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "chromastyles",
			short: "Generate CSS stylesheet for the Chroma code highlighter",
			long: `Generate CSS stylesheet for the Chroma code highlighter for a given style. This stylesheet is needed if markup.highlight.noClasses is disabled in config.

See https://gohugo.io/quick-reference/syntax-highlighting-styles/ for a preview of the available styles.`,

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				style = strings.ToLower(style)
				if !slices.Contains(styles.Names(), style) {
					return fmt.Errorf("invalid style: %s", style)
				}
				var chromaStyle *chroma.Style
				if mode != "" {
					var chromaMode chroma.Mode
					switch mode {
					case "light":
						chromaMode = chroma.Light
					case "dark":
						chromaMode = chroma.Dark
					default:
						return fmt.Errorf("invalid mode: %s", mode)
					}

					chromaStyle = styles.GetForMode(style, chromaMode)
					if chromaStyle.Mode() != chromaMode {
						return fmt.Errorf("style %q does not have a %q mode", style, mode)
					}
				} else {
					chromaStyle = styles.Get(style)
				}
				builder := chromaStyle.Builder()
				if highlightStyle != "" {
					builder.Add(chroma.LineHighlight, highlightStyle)
				}
				if lineNumbersInlineStyle != "" {
					builder.Add(chroma.LineNumbers, lineNumbersInlineStyle)
				}
				if lineNumbersTableStyle != "" {
					builder.Add(chroma.LineNumbersTable, lineNumbersTableStyle)
				}
				style, err := builder.Build()
				if err != nil {
					return err
				}

				if omitEmpty {
					// See https://github.com/alecthomas/chroma/commit/5b2a4c5a26c503c79bc86ba3c4ae5b330028bd3d
					hugo.Deprecate("--omitEmpty", "Flag is no longer needed, empty classes are now always omitted.", "v0.149.0")
				}
				options := []html.Option{
					html.WithCSSComments(!omitClassComments),
				}
				if !modeSelector {
					// Only needed for the shared-scope overlay; --modeSelector
					// gives each mode its own scope and makes this redundant.
					options = append(options, html.WithCustomCSS(chromaCSSOverrides(style)))
				}

				formatter := html.New(options...)

				var buf bytes.Buffer
				fmt.Fprintf(&buf, "/* Generated using: hugo %s */\n\n", strings.Join(os.Args[1:], " "))
				formatter.WriteCSS(&buf, style)
				css := buf.String()
				if modeSelector {
					// Scope every selector under a top level mode class, e.g. ".dark .chroma".
					// This allows generating both light and dark stylesheets and toggling
					// them with a parent class on the page.
					// There's no upstream option for this (I think), so do string replacements for now.
					// TODO(bep) upstream option for this.
					var prefix string
					switch style.Mode() {
					case chroma.Light:
						prefix = ".light "
					case chroma.Dark:
						prefix = ".dark "
					default:
						return fmt.Errorf("style %q does not have a %q mode", style.Name, mode)
					}
					replacer := strings.NewReplacer(
						".bg {", prefix+".bg {",
						".chroma ", prefix+".chroma ",
					)
					css = replacer.Replace(css)
				}

				fmt.Print(css)
				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&style, "style", "friendly", "highlighter style")
				_ = cmd.RegisterFlagCompletionFunc("style", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&mode, "mode", "", `style mode ("light", "dark")`)
				_ = cmd.RegisterFlagCompletionFunc("mode", cobra.FixedCompletions([]string{"light", "dark"}, cobra.ShellCompDirectiveNoFileComp))
				cmd.PersistentFlags().BoolVar(&modeSelector, "modeSelector", false, `scope selectors under a top level mode class, e.g. ".dark .chroma"`)
				_ = cmd.RegisterFlagCompletionFunc("modeSelector", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&highlightStyle, "highlightStyle", "", `foreground and background colors for highlighted lines, e.g. --highlightStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("highlightStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersInlineStyle, "lineNumbersInlineStyle", "", `foreground and background colors for inline line numbers, e.g. --lineNumbersInlineStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersInlineStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersTableStyle, "lineNumbersTableStyle", "", `foreground and background colors for table line numbers, e.g. --lineNumbersTableStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersTableStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().BoolVar(&omitEmpty, "omitEmpty", false, `omit empty CSS rules (deprecated, no longer needed)`)
				_ = cmd.RegisterFlagCompletionFunc("omitEmpty", cobra.NoFileCompletions)
				cmd.PersistentFlags().BoolVar(&omitClassComments, "omitClassComments", false, `omit CSS class comment prefixes in the generated CSS`)
				_ = cmd.RegisterFlagCompletionFunc("omitClassComments", cobra.NoFileCompletions)
			},
		}
	}

	newMan := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "man",
			short: "Generate man pages for the Hugo CLI",
			long: `This command automatically generates up-to-date man pages of Hugo's
	command-line interface.  By default, it creates the man page files
	in the "man" directory under the current directory.`,

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				header := &doc.GenManHeader{
					Section: "1",
					Manual:  "Hugo Manual",
					Source:  fmt.Sprintf("Hugo %s", hugo.CurrentVersion),
				}
				if !strings.HasSuffix(genmandir, helpers.FilePathSeparator) {
					genmandir += helpers.FilePathSeparator
				}
				if found, _ := helpers.Exists(genmandir, hugofs.Os); !found {
					r.Println("Directory", genmandir, "does not exist, creating...")
					if err := hugofs.Os.MkdirAll(genmandir, 0o777); err != nil {
						return err
					}
				}
				cd.CobraCommand.Root().DisableAutoGenTag = true

				r.Println("Generating Hugo man pages in", genmandir, "...")
				doc.GenManTree(cd.CobraCommand.Root(), header, genmandir)

				r.Println("Done.")

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")
				_ = cmd.MarkFlagDirname("dir")
			},
		}
	}

	newGen := func() simplecobra.Commander {
		const gendocFrontmatterTemplate = `---
title: "%s"
slug: %s
url: %s
---
`

		return &simpleCommand{
			name:  "doc",
			short: "Generate Markdown documentation for the Hugo CLI",
			long: `Generate Markdown documentation for the Hugo CLI.
			This command is, mostly, used to create up-to-date documentation
	of Hugo's command-line interface for https://gohugo.io/.

	It creates one Markdown file per command with front matter suitable
	for rendering in Hugo.`,
			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				cd.CobraCommand.VisitParents(func(c *cobra.Command) {
					// Disable the "Auto generated by spf13/cobra on DATE"
					// as it creates a lot of diffs.
					c.DisableAutoGenTag = true
				})
				if !strings.HasSuffix(gendocdir, helpers.FilePathSeparator) {
					gendocdir += helpers.FilePathSeparator
				}
				if found, _ := helpers.Exists(gendocdir, hugofs.Os); !found {
					r.Println("Directory", gendocdir, "does not exist, creating...")
					if err := hugofs.Os.MkdirAll(gendocdir, 0o777); err != nil {
						return err
					}
				}
				prepender := func(filename string) string {
					name := filepath.Base(filename)
					base := strings.TrimSuffix(name, path.Ext(name))
					url := "/commands/" + strings.ToLower(base) + "/"
					return fmt.Sprintf(gendocFrontmatterTemplate, strings.Replace(base, "_", " ", -1), base, url)
				}

				linkHandler := func(name string) string {
					base := strings.TrimSuffix(name, path.Ext(name))
					return "/commands/" + strings.ToLower(base) + "/"
				}
				r.Println("Generating Hugo command-line documentation in", gendocdir, "...")
				doc.GenMarkdownTreeCustom(cd.CobraCommand.Root(), gendocdir, prepender, linkHandler)
				r.Println("Done.")

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&gendocdir, "dir", "/tmp/hugodoc/", "the directory to write the doc.")
				_ = cmd.MarkFlagDirname("dir")
			},
		}
	}

	var docsHelperTarget string

	newDocsHelper := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "docshelper",
			short: "Generate some data files for the Hugo docs",

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				r.Println("Generate docs data to", docsHelperTarget)

				var buf bytes.Buffer
				jsonEnc := json.NewEncoder(&buf)

				configProvider := func() docshelper.DocProvider {
					conf := hugolib.DefaultConfig()
					conf.CacheDir = "" // The default value does not make sense in the docs.
					defaultConfig := parser.NullBoolJSONMarshaller{Wrapped: parser.LowerCaseCamelJSONMarshaller{Value: conf}}
					return docshelper.DocProvider{"config": defaultConfig}
				}

				docshelper.AddDocProviderFunc(configProvider)
				if err := jsonEnc.Encode(docshelper.GetDocProvider()); err != nil {
					return err
				}

				// Decode the JSON to a map[string]interface{} and then unmarshal it again to the correct format.
				var m map[string]any
				if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
					return err
				}

				targetFile := filepath.Join(docsHelperTarget, "docs.yaml")

				f, err := os.Create(targetFile)
				if err != nil {
					return err
				}
				defer f.Close()
				yamlEnc := yaml.NewEncoder(f, yaml.UseSingleQuote(true), yaml.AutoInt())
				if err := yamlEnc.Encode(m); err != nil {
					return err
				}

				r.Println("Done!")
				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.Hidden = true
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVarP(&docsHelperTarget, "dir", "", "docs/data", "data dir")
			},
		}
	}

	newJSONSchema := func() simplecobra.Commander {
		var schemaDir string
		var schemaBaseURL string
		return &simpleCommand{
			name:  "jsonschemas",
			short: "Generate JSON Schema for Hugo config and page structures",
			long:  `Generate a JSON Schema for Hugo configuration options and page structures using reflection.`,
			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				// This function generates JSON Schema for Hugo configuration
				// It creates individual schema files for each configuration section

				// Initialize statistics tracking
				var stats schemaStats
				// and a main schema file that includes RootConfig properties and
				// references to the section schemas.
				// Go comments are used as descriptions in the schema.

				// Create the output directory if it doesn't exist
				if err := os.MkdirAll(schemaDir, 0755); err != nil {
					return fmt.Errorf("failed to create schema directory: %w", err)
				}

				// Build reference string - relative path or absolute URL
				schemaRef := func(name string) string {
					if schemaBaseURL != "" {
						return schemaBaseURL + name + ".schema.json"
					}
					return "./" + name + ".schema.json"
				}

				// Build schema ID URL
				schemaID := func(name string) string {
					if schemaBaseURL != "" {
						return schemaBaseURL + name + ".schema.json"
					}
					return "https://gohugo.io/jsonschemas/" + name + ".schema.json"
				}

				generateDocURL := makeGenerateDocumentationURL()

				// Create a reflector to use Go comments as descriptions
				rf := jsonschema.Reflector{
					KeyNamer:                   camel,
					Namer:                      camelType,
					RequiredFromJSONSchemaTags: true,   // Don't mark fields as required
					DoNotReference:             true,   // Expand references into full definitions
					ExpandedStruct:             true,   // Include all struct fields in the schema
					AllowAdditionalProperties:  true,   // Allow additional properties not in the schema
					FieldNameTag:               "json", // Use json tags for field names
				}

				// Use AddGoComments to include Go comments in the schema as descriptions
				_, currentFile, _, _ := runtime.Caller(0)
				hugoSourceDir := filepath.Dir(filepath.Dir(currentFile))

				if err := rf.AddGoComments("github.com/gohugoio/hugo", hugoSourceDir); err != nil {
					return fmt.Errorf("failed to add Go comments: %w", err)
				}

				// Crawl documentation files to extract URLs for linking
				r.Println("Setting up documentation URL generator...")

				// Dynamically discover configuration sections from the Config struct
				configSections := discoverConfigSections()

				// Generate individual schema files for each section
				for name, fieldInfo := range configSections {

					var schema *jsonschema.Schema
					var schemaErr error

					switch name {
					case "params":
						schema = &jsonschema.Schema{
							Type:        "object",
							Description: "User-defined params. \nhttps://gohugo.io/configuration/all/#params",
							AdditionalProperties: &jsonschema.Schema{
								Type: "string",
							},
						}
					case "permalinks":
						schema = &jsonschema.Schema{
							Type:        "object",
							Description: "Permalink configuration. Maps content kind (page, section, term, taxonomy) to permalink patterns.",
							AdditionalProperties: &jsonschema.Schema{
								Type:        "object",
								Description: "Permalink patterns for this content kind",
								AdditionalProperties: &jsonschema.Schema{
									Type:        "string",
									Description: "The permalink pattern for this section",
								},
							},
						}
					case "languages":
						languageConfigSchema := rf.Reflect(langs.LanguageConfig{})
						languageConfigSchema.Version = "https://json-schema.org/draft-07/schema"
						cleanupSubSchema(languageConfigSchema)

						if languageConfigSchema.Properties == nil {
							languageConfigSchema.Properties = jsonschema.NewProperties()
						}

						rootConfigType := reflect.TypeOf(allconfig.RootConfig{})
						for i := range rootConfigType.NumField() {
							f := rootConfigType.Field(i)
							if !f.IsExported() {
								continue
							}
							propName := camel(f.Name)
							if propName == "" {
								continue
							}
							if _, exists := languageConfigSchema.Properties.Get(propName); exists {
								continue
							}
							if _, isSection := configSections[propName]; isSection {
								languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{
									Ref: schemaRef("hugo-config-" + propName),
								})
								continue
							}
							switch f.Type.Kind() {
							case reflect.String:
								languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{Type: "string"})
							case reflect.Bool:
								languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{Type: "boolean"})
							case reflect.Int, reflect.Int64, reflect.Int32:
								languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{Type: "integer"})
							case reflect.Float64:
								languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{Type: "number"})
							case reflect.Slice:
								if f.Type.Elem().Kind() == reflect.String {
									languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{
										Type:  "array",
										Items: &jsonschema.Schema{Type: "string"},
									})
								}
							case reflect.Map:
								if f.Type.Elem().Kind() == reflect.String {
									languageConfigSchema.Properties.Set(propName, &jsonschema.Schema{
										Type:                 "object",
										AdditionalProperties: &jsonschema.Schema{Type: "string"},
									})
								}
							}

						}

						schema = &jsonschema.Schema{
							Type:                 "object",
							Description:          "Language configuration. Maps language code to language configuration.",
							AdditionalProperties: languageConfigSchema,
						}
					case "markup":
						// Generate markup schema with nested definitions for goldmark, highlight, tableOfContents, asciidocExt
						// Due to a bug in jsonschema reflection with the markup_config.Config struct,
						// we'll manually create the schema structure

						// Create definitions for nested structs
						markupRf := jsonschema.Reflector{
							KeyNamer:                   camel,
							Namer:                      camelType,
							RequiredFromJSONSchemaTags: true,
							DoNotReference:             true,
							ExpandedStruct:             true,
							AllowAdditionalProperties:  true,
						}

						// Add Go comments to the markup reflector
						if err := markupRf.AddGoComments("github.com/gohugoio/hugo", hugoSourceDir); err != nil {
							r.Printf("Warning: failed to add Go comments for markup reflector: %v\n", err)
						}

						// Create the main markup schema manually
						schema = &jsonschema.Schema{
							Type:        "object",
							Description: "Configuration for markup processing",
							Properties:  jsonschema.NewProperties(),
						}

						// Add properties manually based on markup_config.Config struct

						schema.Properties.Set("defaultMarkdownHandler", &jsonschema.Schema{
							Type:        "string",
							Description: "Default markdown handler for md/markdown extensions. Default is 'goldmark'.",
						})

						schema.Properties.Set("highlight", &jsonschema.Schema{
							Ref:         "#/$defs/highlight",
							Description: "Configuration for syntax highlighting.",
						})

						schema.Properties.Set("tableOfContents", &jsonschema.Schema{
							Ref:         "#/$defs/tableOfContents",
							Description: "Table of contents configuration.",
						})

						schema.Properties.Set("goldmark", &jsonschema.Schema{
							Ref:         "#/$defs/goldmark",
							Description: "Configuration for the Goldmark markdown engine.",
						})

						schema.Properties.Set("asciidocExt", &jsonschema.Schema{
							Ref:         "#/$defs/asciidocExt",
							Description: "Configuration for the Asciidoc external markdown engine.",
						})

						// Create definitions for nested structs
						schema.Definitions = make(jsonschema.Definitions)

						// Add goldmark definition
						goldmarkSchema := markupRf.Reflect(goldmark_config.Config{})
						cleanupSubSchema(goldmarkSchema)
						schema.Definitions["goldmark"] = goldmarkSchema

						// Add highlight definition
						highlightSchema := markupRf.Reflect(highlight.Config{})
						cleanupSubSchema(highlightSchema)
						schema.Definitions["highlight"] = highlightSchema

						// Add tableOfContents definition
						tocSchema := markupRf.Reflect(tableofcontents.Config{})
						cleanupSubSchema(tocSchema)
						schema.Definitions["tableOfContents"] = tocSchema

						// Add asciidocExt definition
						asciidocSchema := markupRf.Reflect(asciidocext_config.Config{})
						cleanupSubSchema(asciidocSchema)
						schema.Definitions["asciidocExt"] = asciidocSchema
					default:
						// Create a zero value instance of the type for reflection
						var instance interface{}
						// Handle other types appropriately
						switch fieldInfo.Type.Kind() {
						case reflect.Map:
							// For unnamed map types (e.g. map[string]string), the invopop library
							// crashes with nil pointer dereference when ExpandedStruct:true is set.
							// Use manual schemas for these simple cases.
							switch fieldInfo.Type.String() {
							case "map[string]string":
								schema = &jsonschema.Schema{
									Type:                 "object",
									AdditionalProperties: &jsonschema.Schema{Type: "string"},
								}
								instance = nil
							case "map[string][]string":
								schema = &jsonschema.Schema{
									Type: "object",
									AdditionalProperties: &jsonschema.Schema{
										Type:  "array",
										Items: &jsonschema.Schema{Type: "string"},
									},
								}
								instance = nil
							default:
								instance = reflect.MakeMap(fieldInfo.Type).Interface()
							}
						case reflect.Slice:
							// For slices, create an empty slice
							instance = make([]interface{}, 0)
						case reflect.Chan:
							// Skip channels as they can't be properly serialized
							continue
						case reflect.Func:
							// Skip functions as they can't be properly serialized
							continue
						case reflect.Interface:
							// Skip interfaces as they can't be properly reflected
							continue
						default:
							// For structs and other types, use reflect.New
							instance = reflect.New(fieldInfo.Type).Interface()
						}

						// Use a defer function to catch panics and continue with the next section
						if instance != nil {
							func() {
								defer func() {
									if r := recover(); r != nil {
										schemaErr = fmt.Errorf("panic during reflection: %v", r)
									}
								}()
								schema = rf.Reflect(instance)
							}()
						}
					}

					if schemaErr != nil {
						continue
					}
					schema.ID = jsonschema.ID(schemaID("hugo-config-" + name))
					schema.Version = "https://json-schema.org/draft-07/schema"

					// Clear all required fields
					schema.Required = nil

					removeRequiredFromSchema(schema)

					addDocumentationLinksToSchema(schema, name, []string{}, generateDocURL)

					filename := filepath.Join(schemaDir, fmt.Sprintf("hugo-config-%s.schema.json", name))
					f, err := os.Create(filename)
					if err != nil {
						return fmt.Errorf("failed to create schema file %s: %w", filename, err)
					}

					enc := json.NewEncoder(f)
					enc.SetIndent("", "  ")
					if err := enc.Encode(schema); err != nil {
						f.Close()
						return fmt.Errorf("failed to encode schema for %s: %w", name, err)
					}
					f.Close()
					r.Printf("Wrote %s\n", filename)

					// Update statistics
					stats.schemasGenerated++
					stats.totalProperties += countSchemaProperties(schema)
					stats.documentationLinks += countDocumentationLinks(schema)
				}

				// Generate the main hugo.schema.json that combines RootConfig and references other schemas

				// For the main schema, we'll directly reflect the RootConfig
				rootConfigSchema := rf.Reflect(allconfig.RootConfig{})

				// Set the main schema metadata directly on the schema object
				rootConfigSchema.ID = jsonschema.ID(schemaID("hugo-config"))
				rootConfigSchema.Version = "https://json-schema.org/draft-07/schema"
				rootConfigSchema.Title = "Hugo Configuration Schema"
				rootConfigSchema.Description = "JSON Schema for Hugo configuration files"
				rootConfigSchema.Type = "object"

				// Remove any required fields throughout the schema
				rootConfigSchema.Required = nil
				removeRequiredFromSchema(rootConfigSchema)

				addDocumentationLinksToSchema(rootConfigSchema, "rootconfig", []string{}, generateDocURL)

				// Make sure we have a properties object
				if rootConfigSchema.Properties == nil {
					rootConfigSchema.Properties = jsonschema.NewProperties()
				}

				// Add references to section schemas
				for name := range configSections {
					rootConfigSchema.Properties.Set(camel(name), &jsonschema.Schema{
						Ref: schemaRef("hugo-config-" + name),
					})
				}

				// Create the main hugo.schema.json
				filename := filepath.Join(schemaDir, "hugo-config.schema.json")
				f, err := os.Create(filename)
				if err != nil {
					return fmt.Errorf("failed to create main schema file: %w", err)
				}
				defer f.Close()

				enc := json.NewEncoder(f)
				enc.SetIndent("", "  ")
				if err := enc.Encode(rootConfigSchema); err != nil {
					return fmt.Errorf("failed to encode main schema: %w", err)
				}
				r.Printf("Wrote %s\n", filename)

				// Update statistics for main config schema
				stats.schemasGenerated++
				stats.totalProperties += countSchemaProperties(rootConfigSchema)
				stats.documentationLinks += countDocumentationLinks(rootConfigSchema)

				// Generate Page schemas
				if err := generatePageSchemas(rf, schemaDir, r, &stats, generateDocURL); err != nil {
					return fmt.Errorf("failed to generate page schemas: %w", err)
				}

				// Print generation statistics
				printSchemaStats(stats, r)

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.PersistentFlags().StringVarP(&schemaDir, "dir", "", filepath.Join(os.TempDir(), "hugo-schemas"), "output directory for schema files")
				cmd.PersistentFlags().StringVarP(&schemaBaseURL, "baseURL", "", "", "base URL for schema references (default: relative paths)")
			},
		}
	}

	return &genCommand{
		commands: []simplecobra.Commander{
			newChromaStyles(),
			newGen(),
			newMan(),
			newDocsHelper(),
			newJSONSchema(),
		},
	}
}

// chromaCSSOverrides re-emits leaf token colors that Chroma's minifier drops
// because they equal the style's default foreground (the .chroma color). Without
// modeSelector scoping, a paired light/dark stylesheet shares the same .chroma
// scope, and the light sheet's explicit rule (e.g. .chroma .nx) would otherwise
// leak into dark mode, since an explicit declaration beats inheritance.
func chromaCSSOverrides(style *chroma.Style) map[chroma.TokenType]string {
	bg := style.Get(chroma.Background)
	m := make(map[chroma.TokenType]string)
	for tt := range chroma.StandardTypes {
		if tt == chroma.Background || !style.Has(tt) || !chromaLeafToken(tt) {
			continue
		}
		entry := style.Get(tt)
		if !entry.Sub(bg).IsZero() || !entry.Colour.IsSet() {
			continue
		}
		if css := html.StyleEntryToCSS(chroma.StyleEntry{Colour: entry.Colour}); css != "" {
			m[tt] = css
		}
	}
	return m
}

func chromaLeafToken(tt chroma.TokenType) bool {
	for other := range chroma.StandardTypes {
		if other != tt && (other.Category() == tt || other.SubCategory() == tt) {
			return false
		}
	}
	return true
}

type genCommand struct {
	rootCmd *rootCommand

	commands []simplecobra.Commander
}

func (c *genCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *genCommand) Name() string {
	return "gen"
}

func (c *genCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return nil
}

func (c *genCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Generate documentation and syntax highlighting styles"
	cmd.Long = "Generate documentation for your project using Hugo's documentation engine, including syntax highlighting for various programming languages."

	cmd.RunE = nil
	return nil
}

func (c *genCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.rootCmd = cd.Root.Command.(*rootCommand)
	return nil
}

type configFieldInfo struct {
	Type reflect.Type
	Tag  reflect.StructTag
}

func camel(s string) string {
	if s == "" {
		return ""
	}
	if len(s) > 1 && strings.ToUpper(s[:2]) == s[:2] {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func camelType(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		return ""
	}
	return strings.ToLower(name[:1]) + name[1:]
}

func canReflectType(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Interface {
		return false
	}
	if t.Kind() == reflect.Func {
		return false
	}
	if t.Kind() == reflect.Chan {
		return false
	}
	if t.Kind() == reflect.UnsafePointer {
		return false
	}
	return true
}

func isConfigNamespaceType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// In Go 1.26, Name() for instantiated generics returns e.g. "ConfigNamespace[map[string]...]"
	return strings.HasPrefix(t.Name(), "ConfigNamespace") && t.PkgPath() == "github.com/gohugoio/hugo/config"
}

func extractConfigNamespaceType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s := t.String()
	bracket := strings.Index(s, "[")
	if bracket == -1 {
		return nil
	}
	inner := s[bracket+1 : len(s)-1]
	// inner is the type args like "map[string]navigation.MenuConfig, navigation.Menus"
	// Find the first top-level comma to split the two args
	depth := 0
	split := -1
	for i, c := range inner {
		switch c {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				split = i
				goto found
			}
		}
	}
found:
	if split == -1 {
		return nil
	}
	firstArg := strings.TrimSpace(inner[:split])
	if strings.HasPrefix(firstArg, "map[") {
		// map[string]T — extract T
		mapValStart := strings.LastIndex(firstArg, "]")
		if mapValStart == -1 {
			return nil
		}
		mapVal := strings.TrimSpace(firstArg[mapValStart+1:])
		return typeFromName(mapVal)
	}
	return typeFromName(firstArg)
}

func typeFromName(name string) reflect.Type {
	// In Go 1.26, reflect.Type.String() returns fully-qualified package paths
	// (e.g. "github.com/gohugoio/hugo/resources/images.ImagingConfig").
	// Extract just the short package.TypeName (e.g. "images.ImagingConfig").
	lastDot := strings.LastIndex(name, ".")
	if lastDot != -1 {
		pkgPath := name[:lastDot]
		lastSlash := strings.LastIndex(pkgPath, "/")
		if lastSlash != -1 {
			pkgBase := pkgPath[lastSlash+1:]
			shortName := pkgBase + name[lastDot:]
			name = shortName
		}
	}
	switch name {
	case "images.ImagingConfig":
		return reflect.TypeOf(images.ImagingConfig{})
	case "media.ContentTypeConfig":
		return reflect.TypeOf(media.ContentTypeConfig{})
	case "media.MediaTypeConfig":
		return reflect.TypeOf(media.MediaTypeConfig{})
	case "navigation.MenuConfig":
		return reflect.TypeOf(navigation.MenuConfig{})
	case "output.OutputFormatConfig":
		return reflect.TypeOf(output.OutputFormatConfig{})
	case "langs.LanguageConfig":
		return reflect.TypeOf(langs.LanguageConfig{})
	case "segments.SegmentConfig":
		return reflect.TypeOf(segments.SegmentConfig{})
	case "versions.VersionConfig":
		return reflect.TypeOf(versions.VersionConfig{})
	case "roles.RoleConfig":
		return reflect.TypeOf(roles.RoleConfig{})
	default:
		return nil
	}
}

// cleanupSubSchema performs standard cleanup on a schema for use as a definition
func cleanupSubSchema(schema *jsonschema.Schema) {
	schema.Required = nil
	removeRequiredFromSchema(schema)
	schema.Version = "https://json-schema.org/draft-07/schema"
	schema.ID = ""
}

// removeRequiredFromSchema removes "required" fields from the schema and all nested schemas
func removeRequiredFromSchema(schema *jsonschema.Schema) {
	if schema == nil {
		return
	}

	// Clear the required field
	schema.Required = nil

	// Recursively remove required from properties
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			removeRequiredFromSchema(pair.Value)
		}
	}

	// Recursively remove required from pattern properties
	if schema.PatternProperties != nil {
		for _, propSchema := range schema.PatternProperties {
			removeRequiredFromSchema(propSchema)
		}
	}

	// Recursively remove required from additional properties
	if schema.AdditionalProperties != nil {
		removeRequiredFromSchema(schema.AdditionalProperties)
	}

	// Recursively remove required from items (for arrays)
	if schema.Items != nil {
		removeRequiredFromSchema(schema.Items)
	}

	// Recursively remove required from all conditional schemas
	for _, condSchema := range schema.AllOf {
		removeRequiredFromSchema(condSchema)
	}
	for _, condSchema := range schema.AnyOf {
		removeRequiredFromSchema(condSchema)
	}
	for _, condSchema := range schema.OneOf {
		removeRequiredFromSchema(condSchema)
	}

	// Remove required from not schema
	if schema.Not != nil {
		removeRequiredFromSchema(schema.Not)
	}

	// Remove required from definitions
	if schema.Definitions != nil {
		for _, defSchema := range schema.Definitions {
			removeRequiredFromSchema(defSchema)
		}
	}
}

func walkConfigFields(t reflect.Type, fn func(reflect.StructField)) {
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			walkConfigFields(f.Type, fn)
			continue
		}
		fn(f)
	}
}

func discoverConfigSections() map[string]configFieldInfo {
	sections := make(map[string]configFieldInfo)
	configType := reflect.TypeOf(allconfig.Config{})

	walkConfigFields(configType, func(field reflect.StructField) {
		if !field.IsExported() {
			return
		}
		if field.Tag.Get("mapstructure") != "-" {
			return
		}
		if field.Tag.Get("json") == "-" {
			return
		}

		fieldType := field.Type
		schemaName := camel(field.Name)

		if isConfigNamespaceType(fieldType) {
			if t := extractConfigNamespaceType(fieldType); t != nil {
				sections[schemaName] = configFieldInfo{Type: t, Tag: field.Tag}
			}
			return
		}

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if !canReflectType(fieldType) {
			return
		}

		sections[schemaName] = configFieldInfo{Type: fieldType, Tag: field.Tag}
	})

	return sections
}

func generatePageSchemas(rf jsonschema.Reflector, schemaDir string, r *rootCommand, stats *schemaStats, generateDocURL func(string, []string) string) error {
	// Generate a single schema for Hugo page front matter
	schema := rf.Reflect(pagemeta.PageConfigLate{})

	// Set schema metadata
	schema.ID = jsonschema.ID("https://gohugo.io/jsonschemas/hugo-content.schema.json")
	schema.Version = "https://json-schema.org/draft-07/schema"
	schema.Title = "Hugo Page Front Matter Schema"
	schema.Description = "JSON Schema for Hugo page front matter structure"

	// Clear required fields
	schema.Required = nil

	// Remove required fields from nested schemas
	removeRequiredFromSchema(schema)

	addDocumentationLinksToSchema(schema, "page-frontmatter", []string{}, generateDocURL)

	// Add flexible date definitions
	addFlexibleDateDefinitions(schema)

	// Remove the content property as it represents the file content after frontmatter, not a frontmatter field
	if schema.Properties != nil {
		schema.Properties.Delete("content")
	}

	// Write schema file
	filename := filepath.Join(schemaDir, "hugo-content.schema.json")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create page schema file %s: %w", filename, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return fmt.Errorf("failed to encode page schema: %w", err)
	}

	r.Printf("Wrote page schema %s\n", filename)

	// Update statistics for page schema
	stats.schemasGenerated++
	stats.totalProperties += countSchemaProperties(schema)
	stats.documentationLinks += countDocumentationLinks(schema)

	return nil
}

// addFlexibleDateDefinitions adds flexible date definitions to the page schema
func addFlexibleDateDefinitions(schema *jsonschema.Schema) {
	if schema.Properties == nil {
		return
	}

	// Initialize $defs if not present
	if schema.Definitions == nil {
		schema.Definitions = jsonschema.Definitions{}
	}

	// Create a shared flexible date definition
	flexibleDateDef := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Format:      "date-time",
				Description: "RFC3339 date-time format (e.g., 2017-01-02T15:04:05Z07:00)",
			},
			{
				Type:        "string",
				Format:      "date",
				Description: "Date only format (e.g., 2017-01-02)",
			},
			{
				Type:        "string",
				Pattern:     "^\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}:\\d{2}$",
				Description: "Date with time format (e.g., 2017-01-02 15:04:05)",
			},
			{
				Type:        "string",
				Pattern:     "^\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}$",
				Description: "Date with hour:minute format (e.g., 2017-01-02 15:04)",
			},
		},
		Description: "Hugo date field supporting multiple formats",
	}

	// Add the definition to $defs
	schema.Definitions["flexibleDate"] = flexibleDateDef

	// Apply flexible date schema to date-related properties using $ref
	dateProperties := []string{"date", "publishDate", "expiryDate", "lastmod"}
	for _, propName := range dateProperties {
		if propSchema, exists := schema.Properties.Get(propName); exists {
			// Preserve the original description and documentation URL
			originalDesc := propSchema.Description

			// Replace with $ref to shared definition
			*propSchema = jsonschema.Schema{
				Ref:         "#/$defs/flexibleDate",
				Description: originalDesc,
			}
		}
	}
}

func addDocumentationLinksToSchema(schema *jsonschema.Schema, sectionName string, propertyPath []string, generateDocURL func(string, []string) string) {
	if schema == nil {
		return
	}

	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			propName := pair.Key
			propSchema := pair.Value
			currentPath := append(propertyPath, propName)

			docURL := generateDocURL(sectionName, currentPath)
			if docURL != "" {
				if propSchema.Description != "" {
					propSchema.Description = propSchema.Description + " \n" + docURL
				} else {
					propSchema.Description = docURL
				}
			}
			addDocumentationLinksToSchema(propSchema, sectionName, currentPath, generateDocURL)
		}
	}

	if schema.AdditionalProperties != nil {
		addDocumentationLinksToSchema(schema.AdditionalProperties, sectionName, propertyPath, generateDocURL)
	}

	if schema.Items != nil {
		addDocumentationLinksToSchema(schema.Items, sectionName, propertyPath, generateDocURL)
	}

	for _, condSchema := range schema.AllOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath, generateDocURL)
	}
	for _, condSchema := range schema.AnyOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath, generateDocURL)
	}
	for _, condSchema := range schema.OneOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath, generateDocURL)
	}

	if schema.Not != nil {
		addDocumentationLinksToSchema(schema.Not, sectionName, propertyPath, generateDocURL)
	}

	if schema.Definitions != nil {
		for defName, defSchema := range schema.Definitions {
			defSectionName := sectionName
			if sectionName == "markup" {
				switch defName {
				case "goldmark":
					defSectionName = "goldmark"
				case "highlight":
					defSectionName = "highlight"
				case "tableOfContents":
					defSectionName = "tableofcontents"
				case "asciidocExt":
					defSectionName = "asciidocext"
				}
			}
			addDocumentationLinksToSchema(defSchema, defSectionName, []string{}, generateDocURL)
		}
	}
}

func makeGenerateDocumentationURL() func(string, []string) string {
	var (
		sectionAnchors map[string]map[string]string
		docAnchors     map[string]string
		mu             sync.Mutex
	)

	getDocAnchors := func() map[string]string {
		mu.Lock()
		defer mu.Unlock()
		if docAnchors == nil {
			docAnchors = extractAnchorsFromHTML()
		}
		return docAnchors
	}

	getSectionAnchors := func(parentSection string) map[string]string {
		mu.Lock()
		defer mu.Unlock()
		if sectionAnchors == nil {
			sectionAnchors = make(map[string]map[string]string)
		}
		if _, exists := sectionAnchors[parentSection]; !exists {
			sectionAnchors[parentSection] = extractSectionAnchors(parentSection)
		}
		return sectionAnchors[parentSection]
	}

	return func(sectionName string, propertyPath []string) string {
		if len(propertyPath) == 0 {
			return ""
		}

		baseURL := "https://gohugo.io/configuration"

		if sectionName == "rootconfig" && len(propertyPath) > 0 {
			anchors := getDocAnchors()
			propertyName := propertyPath[len(propertyPath)-1]
			if anchorID, exists := anchors[propertyName]; exists {
				return baseURL + "/all/#" + anchorID
			}
			return baseURL + "/all/#" + strings.ToLower(propertyName)
		}

		if sectionName == "frontmatter" {
			return baseURL + "/front-matter/#dates"
		}

		if sectionName == "page-frontmatter" {
			ctmBaseURL := "https://gohugo.io/content-management/front-matter"
			anchors := getSectionAnchors("content-frontmatter")
			if len(propertyPath) > 0 {
				propertyName := propertyPath[len(propertyPath)-1]
				if anchorID, found := anchors[propertyName]; found {
					return ctmBaseURL + "#" + anchorID
				}
				if anchorID, found := anchors[strings.ToLower(propertyName)]; found {
					return ctmBaseURL + "#" + anchorID
				}
				if len(propertyPath) > 1 {
					fullPath := strings.ToLower(strings.Join(propertyPath, ""))
					if anchorID, found := anchors[fullPath]; found {
						return ctmBaseURL + "#" + anchorID
					}
				}
			}
			return ctmBaseURL
		}

		sectionURL := getSectionDocumentationURL(baseURL, sectionName)

		parentSection := sectionName
		if sectionName == "goldmark" || sectionName == "asciidocext" || sectionName == "highlight" || sectionName == "tableofcontents" {
			parentSection = "markup"
		}

		if len(propertyPath) > 0 {
			anchors := getSectionAnchors(parentSection)
			propertyName := propertyPath[len(propertyPath)-1]
			if parentSection == "markup" && propertyName == "defaultMarkdownHandler" {
				return sectionURL + "#default-handler"
			}
			if anchorID, found := anchors[propertyName]; found {
				return sectionURL + "#" + anchorID
			}
			if anchorID, found := anchors[strings.ToLower(propertyName)]; found {
				return sectionURL + "#" + anchorID
			}
			if len(propertyPath) > 1 {
				fullPath := strings.ToLower(strings.Join(propertyPath, ""))
				if anchorID, found := anchors[fullPath]; found {
					return sectionURL + "#" + anchorID
				}
			}
			for i := len(propertyPath) - 2; i >= 0; i-- {
				parentName := propertyPath[i]
				if anchorID, found := anchors[parentName]; found {
					return sectionURL + "#" + anchorID
				}
				if anchorID, found := anchors[strings.ToLower(parentName)]; found {
					return sectionURL + "#" + anchorID
				}
			}
			if len(propertyPath) > 1 {
				fullPath := strings.ToLower(strings.Join(propertyPath, ""))
				return sectionURL + "#" + fullPath
			}
			anchor := camelToKebab(propertyPath[0])
			return sectionURL + "#" + anchor
		}

		return sectionURL
	}
}

// camelToKebab converts camelCase strings to kebab-case
func camelToKebab(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// extractAnchorsFromHTML extracts anchor IDs from the built Hugo documentation
func extractAnchorsFromHTML() map[string]string {
	anchors := make(map[string]string)

	// Try different possible locations for the built docs using systematic approach
	possiblePaths := getRootConfigDocPaths()

	var htmlFile string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			htmlFile = path
			break
		}
	}

	if htmlFile == "" {
		fmt.Printf("Warning: Could not find built documentation HTML file\n")
		return anchors
	}

	content, err := os.ReadFile(htmlFile)
	if err != nil {
		fmt.Printf("Warning: Could not read documentation HTML file: %v\n", err)
		return anchors
	}

	// Extract anchor IDs using regex
	// Look for id="propertyname" patterns in dt elements
	re := regexp.MustCompile(`<dt[^>]+id="([^"]+)"[^>]*>([^<]+)</dt>`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) >= 3 {
			anchorID := match[1]
			propertyName := strings.TrimSpace(match[2])

			// Skip non-property anchors (like "settings")
			if anchorID != "settings" && propertyName != "" {
				// Map both the original property name and lowercase version to the anchor
				anchors[propertyName] = anchorID
				anchors[strings.ToLower(propertyName)] = anchorID
			}
		}
	}

	fmt.Printf("Found %d configuration property anchors in documentation\n", len(anchors)/2)
	return anchors
}

// getSectionPaths returns the possible HTML file paths for a given section
func getSectionPaths(sectionName string) []string {
	dirName := sectionName
	switch sectionName {
	case "outputformats":
		dirName = "output-formats"
	case "mediatypes":
		dirName = "media-types"
	case "related":
		dirName = "related-content"
	case "frontmatter":
		dirName = "front-matter"
	case "content-frontmatter":
		dirName = "front-matter"
	}

	baseDirs := []string{
		"docs/public",
		"public-docs",
		"../public-docs",
		"public",
	}

	var subDir string
	if sectionName == "content-frontmatter" {
		subDir = "content-management"
	} else {
		subDir = "configuration"
	}

	var paths []string
	for _, baseDir := range baseDirs {
		path := filepath.Join(baseDir, subDir, dirName, "index.html")
		paths = append(paths, path)
	}

	return paths
}

// extractSectionAnchors extracts anchor IDs from a specific section's documentation page
func extractSectionAnchors(sectionName string) map[string]string {
	anchors := make(map[string]string)

	// Get section paths using the getSectionPaths function
	paths := getSectionPaths(sectionName)
	if len(paths) == 0 {
		return anchors
	}

	var htmlFile string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			htmlFile = path
			break
		}
	}

	if htmlFile == "" {
		return anchors
	}

	content, err := os.ReadFile(htmlFile)
	if err != nil {
		return anchors
	}

	// Extract all anchor IDs from heading elements and definition lists
	// Look for both <h* id="..."> and <dt id="..."> patterns
	reHeadings := regexp.MustCompile(`<h[1-6][^>]+id="([^"]+)"[^>]*>([^<]+)</h[1-6]>`)
	reDefs := regexp.MustCompile(`<dt[^>]+id="([^"]+)"[^>]*>([^<]+)</dt>`)

	// Extract from headings
	matches := reHeadings.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) >= 3 {
			anchorID := match[1]
			headingText := strings.TrimSpace(match[2])

			if anchorID != "" && headingText != "" {
				anchors[headingText] = anchorID
			}
		}
	}

	// Extract from definition terms
	matches = reDefs.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) >= 3 {
			anchorID := match[1]
			termText := strings.TrimSpace(match[2])

			if anchorID != "" && termText != "" {
				anchors[termText] = anchorID
			}
		}
	}

	return anchors
}

// getSectionDocumentationURL returns the documentation URL for a given section
func getSectionDocumentationURL(baseURL, sectionName string) string {
	// Handle special cases where section name doesn't match directory name
	dirName := sectionName
	switch sectionName {
	case "outputformats":
		dirName = "output-formats"
	case "mediatypes":
		dirName = "media-types"
	case "related":
		dirName = "related-content"
	case "httpcache":
		dirName = "http-cache"
	case "frontmatter":
		dirName = "front-matter"
	case "highlight", "goldmark", "tableofcontents", "asciidocext":
		// These are all subsections of markup
		dirName = "markup"
	}

	// Simple systematic approach: baseURL already contains /configuration, just add /dirName/
	return baseURL + "/" + dirName + "/"
}

// getRootConfigDocPaths returns the possible HTML file paths for the root configuration documentation
func getRootConfigDocPaths() []string {
	// Try multiple possible base directories in order of preference
	baseDirs := []string{
		"docs/public",    // Most likely location for built docs
		"public-docs",    // Alternative location
		"../public-docs", // If running from subdirectory
		"public",         // Alternative name
	}

	var paths []string
	for _, baseDir := range baseDirs {
		path := filepath.Join(baseDir, "configuration", "all", "index.html")
		paths = append(paths, path)
	}

	return paths
}

// countSchemaProperties recursively counts all properties in a schema
func countSchemaProperties(schema *jsonschema.Schema) int {
	if schema == nil {
		return 0
	}

	count := 0

	// Count direct properties
	if schema.Properties != nil {
		count += schema.Properties.Len()

		// Recursively count nested properties
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			count += countSchemaProperties(pair.Value)
		}
	}

	// Count properties in definitions
	if schema.Definitions != nil {
		for _, defSchema := range schema.Definitions {
			count += countSchemaProperties(defSchema)
		}
	}

	// Count properties in array items
	if schema.Items != nil {
		count += countSchemaProperties(schema.Items)
	}

	// Count properties in additional properties
	if schema.AdditionalProperties != nil {
		count += countSchemaProperties(schema.AdditionalProperties)
	}

	// Count properties in conditional schemas
	for _, condSchema := range schema.AnyOf {
		count += countSchemaProperties(condSchema)
	}
	for _, condSchema := range schema.OneOf {
		count += countSchemaProperties(condSchema)
	}
	for _, condSchema := range schema.AllOf {
		count += countSchemaProperties(condSchema)
	}
	if schema.Not != nil {
		count += countSchemaProperties(schema.Not)
	}

	return count
}

// countDocumentationLinks recursively counts all documentation links in a schema
func countDocumentationLinks(schema *jsonschema.Schema) int {
	if schema == nil {
		return 0
	}

	count := 0

	// Check if this schema has a documentation link in its description
	if schema.Description != "" && strings.Contains(schema.Description, "https://gohugo.io") {
		count++
	}

	// Count links in direct properties
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			count += countDocumentationLinks(pair.Value)
		}
	}

	// Count links in definitions
	if schema.Definitions != nil {
		for _, defSchema := range schema.Definitions {
			count += countDocumentationLinks(defSchema)
		}
	}

	// Count links in array items
	if schema.Items != nil {
		count += countDocumentationLinks(schema.Items)
	}

	// Count links in additional properties
	if schema.AdditionalProperties != nil {
		count += countDocumentationLinks(schema.AdditionalProperties)
	}

	// Count links in conditional schemas
	for _, condSchema := range schema.AnyOf {
		count += countDocumentationLinks(condSchema)
	}
	for _, condSchema := range schema.OneOf {
		count += countDocumentationLinks(condSchema)
	}
	for _, condSchema := range schema.AllOf {
		count += countDocumentationLinks(condSchema)
	}
	if schema.Not != nil {
		count += countDocumentationLinks(schema.Not)
	}

	return count
}

// printSchemaStats prints a summary of schema generation statistics
func printSchemaStats(stats schemaStats, r *rootCommand) {
	r.Println()
	r.Println("Schema Generation Summary:")
	r.Printf("  Schemas generated: %d\n", stats.schemasGenerated)
	r.Printf("  Total properties: %d\n", stats.totalProperties)
	r.Printf("  Documentation links: %d\n", stats.documentationLinks)
	if stats.totalProperties > 0 {
		coverage := float64(stats.documentationLinks) / float64(stats.totalProperties) * 100
		r.Printf("  Documentation coverage: %.1f%%\n", coverage)
	}
}
