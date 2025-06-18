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
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/bep/simplecobra"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/hugolib/segments"
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
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/page/pagemeta"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"gopkg.in/yaml.v2"
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
		highlightStyle         string
		lineNumbersInlineStyle string
		lineNumbersTableStyle  string
		omitEmpty              bool
	)

	newChromaStyles := func() simplecobra.Commander {
		return &simpleCommand{
			name:  "chromastyles",
			short: "Generate CSS stylesheet for the Chroma code highlighter",
			long: `Generate CSS stylesheet for the Chroma code highlighter for a given style. This stylesheet is needed if markup.highlight.noClasses is disabled in config.

See https://xyproto.github.io/splash/docs/all.html for a preview of the available styles`,

			run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
				style = strings.ToLower(style)
				if !slices.Contains(styles.Names(), style) {
					return fmt.Errorf("invalid style: %s", style)
				}
				builder := styles.Get(style).Builder()
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

				var formatter *html.Formatter
				if omitEmpty {
					formatter = html.New(html.WithClasses(true))
				} else {
					formatter = html.New(html.WithAllClasses(true))
				}

				w := os.Stdout
				fmt.Fprintf(w, "/* Generated using: hugo %s */\n\n", strings.Join(os.Args[1:], " "))
				formatter.WriteCSS(w, style)
				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.ValidArgsFunction = cobra.NoFileCompletions
				cmd.PersistentFlags().StringVar(&style, "style", "friendly", "highlighter style (see https://xyproto.github.io/splash/docs/)")
				_ = cmd.RegisterFlagCompletionFunc("style", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&highlightStyle, "highlightStyle", "", `foreground and background colors for highlighted lines, e.g. --highlightStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("highlightStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersInlineStyle, "lineNumbersInlineStyle", "", `foreground and background colors for inline line numbers, e.g. --lineNumbersInlineStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersInlineStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().StringVar(&lineNumbersTableStyle, "lineNumbersTableStyle", "", `foreground and background colors for table line numbers, e.g. --lineNumbersTableStyle "#fff000 bg:#000fff"`)
				_ = cmd.RegisterFlagCompletionFunc("lineNumbersTableStyle", cobra.NoFileCompletions)
				cmd.PersistentFlags().BoolVar(&omitEmpty, "omitEmpty", false, `omit empty CSS rules`)
				_ = cmd.RegisterFlagCompletionFunc("omitEmpty", cobra.NoFileCompletions)
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
				yamlEnc := yaml.NewEncoder(f)
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

				camel := func(s string) string {
					if s == "" {
						return ""
					}
					return strings.ToLower(s[:1]) + s[1:]
				}

				camelType := func(t reflect.Type) string {
					name := t.Name()
					if name == "" {
						return ""
					}
					return strings.ToLower(name[:1]) + name[1:]
				}

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
				// We need to find the Hugo source directory to read Go comments
				// Use runtime to get the current file's directory as a reference to Hugo source
				_, currentFile, _, _ := runtime.Caller(0)
				hugoSourceDir := filepath.Dir(filepath.Dir(currentFile)) // Go up one level from commands/ to hugo/

				// Change to the Hugo source directory to make relative paths work
				originalWd, _ := os.Getwd()
				os.Chdir(hugoSourceDir)
				defer os.Chdir(originalWd)

				if err := rf.AddGoComments("github.com/gohugoio/hugo", "."); err != nil {
					return fmt.Errorf("failed to add Go comments: %w", err)
				}

				// Crawl documentation files to extract URLs for linking
				r.Println("Setting up documentation URL generator...")

				// Dynamically discover configuration sections from the Config struct
				configSections := discoverConfigSections()

				// Generate individual schema files for each section
				for name, fieldInfo := range configSections {

					// Declare schema variable early
					var schema *jsonschema.Schema
					var schemaErr error

					// Handle the 4 map types using additionalProperties pattern instead of wrapper structs
					switch name {
					case "taxonomies":
						// Taxonomies map[string]string - create schema with additionalProperties
						schema = &jsonschema.Schema{
							Type:        "object",
							Description: "Taxonomy configuration. Maps singular taxonomy name to plural form.",
							AdditionalProperties: &jsonschema.Schema{
								Type:        "string",
								Description: "The plural form of the taxonomy",
							},
						}
					case "languages":
						// Languages map[string]langs.LanguageConfig - create schema with additionalProperties
						// This includes both the basic language config and all localized settings
						languageConfigSchema := rf.Reflect(langs.LanguageConfig{})
						languageConfigSchema.Version = "https://json-schema.org/draft-07/schema"
						cleanupSubSchema(languageConfigSchema)

						// Add all the localized configuration settings that can be defined per language
						// These are the same configuration sections that can be localized
						if languageConfigSchema.Properties == nil {
							languageConfigSchema.Properties = jsonschema.NewProperties()
						}

						// Add references to all localizable configuration sections
						localizableConfigs := map[string]string{
							"baseURL":                    "string",
							"buildDrafts":                "boolean",
							"buildExpired":               "boolean",
							"buildFuture":                "boolean",
							"canonifyURLs":               "boolean",
							"capitalizeListTitles":       "boolean",
							"contentDir":                 "string",
							"copyright":                  "string",
							"disableAliases":             "boolean",
							"disableHugoGeneratorInject": "boolean",
							"disableKinds":               "array",
							"disableLiveReload":          "boolean",
							"disablePathToLower":         "boolean",
							"enableEmoji":                "boolean",
							"hasCJKLanguage":             "boolean",
							"mainSections":               "array",
							"pluralizeListTitles":        "boolean",
							"refLinksErrorLevel":         "string",
							"refLinksNotFoundURL":        "string",
							"relativeURLs":               "boolean",
							"removePathAccents":          "boolean",
							"renderSegments":             "array",
							"sectionPagesMenu":           "string",
							"staticDir":                  "array",
							"summaryLength":              "integer",
							"timeZone":                   "string",
							"titleCaseStyle":             "string",
						}

						// Add simple properties that have primitive types
						for propName, propType := range localizableConfigs {
							propSchema := createPrimitiveSchema(propType)
							languageConfigSchema.Properties.Set(propName, propSchema)
						}

						// Add references to complex configuration sections that can be localized
						complexConfigs := []string{
							"frontmatter", "markup", "mediatypes", "menus", "outputformats",
							"outputs", "page", "pagination", "params", "permalinks", "privacy",
							"related", "security", "services", "sitemap", "taxonomies",
						}

						for _, configName := range complexConfigs {
							languageConfigSchema.Properties.Set(configName, &jsonschema.Schema{
								Ref: fmt.Sprintf("http://gohugo.io/jsonschemas/hugo-config-%s.schema.json", configName),
							})
						}

						schema = &jsonschema.Schema{
							Type:                 "object",
							Description:          "Language configuration. Maps language code to language configuration.",
							AdditionalProperties: languageConfigSchema,
						}
					case "permalinks":
						// Permalinks map[string]map[string]string - create schema with additionalProperties
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
					case "outputs":
						// Outputs map[string][]string - create schema with additionalProperties
						schema = &jsonschema.Schema{
							Type:        "object",
							Description: "Output format configuration. Maps page kind to list of output formats.",
							AdditionalProperties: &jsonschema.Schema{
								Type:        "array",
								Description: "List of output formats for this page kind",
								Items: &jsonschema.Schema{
									Type:        "string",
									Description: "Output format name",
								},
							},
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
						if err := markupRf.AddGoComments("github.com/gohugoio/hugo", "."); err != nil {
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
							// For other maps, create an empty map
							instance = make(map[string]interface{})
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
						func() {
							defer func() {
								if r := recover(); r != nil {
									schemaErr = fmt.Errorf("panic during reflection: %v", r)
								}
							}()
							schema = rf.Reflect(instance)
						}()
					}

					if schemaErr != nil {
						continue
					}
					schema.ID = jsonschema.ID(fmt.Sprintf("https://gohugo.io/jsonschemas/hugo-config-%s.schema.json", name))
					schema.Version = "https://json-schema.org/draft-07/schema"

					// Clear all required fields
					schema.Required = nil

					// Remove required fields from nested objects by working directly with the schema
					removeRequiredFromSchema(schema)

					// Transform enums to oneOf patterns following SchemaStore recommendations
					transformEnumsToOneOf(schema)

					// Add documentation links to schema descriptions
					addDocumentationLinksToSchema(schema, name, []string{})

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
				rootConfigSchema.ID = jsonschema.ID("https://gohugo.io/jsonschemas/hugo-config.schema.json")
				rootConfigSchema.Version = "https://json-schema.org/draft-07/schema"
				rootConfigSchema.Title = "Hugo Configuration Schema"
				rootConfigSchema.Description = "JSON Schema for Hugo configuration files"
				rootConfigSchema.Type = "object"

				// Remove any required fields throughout the schema
				rootConfigSchema.Required = nil
				removeRequiredFromSchema(rootConfigSchema)

				// Add documentation links to root config schema
				addDocumentationLinksToSchema(rootConfigSchema, "rootconfig", []string{})

				// Make sure we have a properties object
				if rootConfigSchema.Properties == nil {
					rootConfigSchema.Properties = jsonschema.NewProperties()
				}

				// Add references to section schemas
				for name := range configSections {
					rootConfigSchema.Properties.Set(camel(name), &jsonschema.Schema{
						Ref: fmt.Sprintf("http://gohugo.io/jsonschemas/hugo-config-%s.schema.json", name),
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
				if err := generatePageSchemas(rf, schemaDir, r, &stats); err != nil {
					return fmt.Errorf("failed to generate page schemas: %w", err)
				}

				// Print generation statistics
				printSchemaStats(stats, r)

				return nil
			},
			withc: func(cmd *cobra.Command, r *rootCommand) {
				cmd.PersistentFlags().StringVarP(&schemaDir, "dir", "", "/tmp/hugo-schemas", "output directory for schema files")
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

// canReflectType checks if a type can be safely reflected for JSON schema generation
func canReflectType(t reflect.Type) bool {
	// Skip invalid types
	if t == nil {
		return false
	}

	// Skip interface types as they can't be properly reflected
	if t.Kind() == reflect.Interface {
		return false
	}

	// Skip function types
	if t.Kind() == reflect.Func {
		return false
	}

	// Skip channel types
	if t.Kind() == reflect.Chan {
		return false
	}

	// Skip unsafe pointer types
	if t.Kind() == reflect.UnsafePointer {
		return false
	}

	// Check for problematic type names that indicate complex generics
	typeName := t.String()
	return !strings.Contains(typeName, "ConfigNamespace")
}

// cleanupSubSchema performs standard cleanup on a schema for use as a definition
func cleanupSubSchema(schema *jsonschema.Schema) {
	schema.Required = nil
	removeRequiredFromSchema(schema)
	schema.Version = "https://json-schema.org/draft-07/schema"
	schema.ID = ""
}

// createPrimitiveSchema creates a simple schema for primitive types
func createPrimitiveSchema(typeName string) *jsonschema.Schema {
	switch typeName {
	case "string":
		return &jsonschema.Schema{Type: "string"}
	case "boolean":
		return &jsonschema.Schema{Type: "boolean"}
	case "integer":
		return &jsonschema.Schema{Type: "integer"}
	case "array":
		return &jsonschema.Schema{
			Type:  "array",
			Items: &jsonschema.Schema{Type: "string"},
		}
	default:
		return &jsonschema.Schema{Type: "string"} // fallback
	}
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

func discoverConfigSections() map[string]configFieldInfo {
	sections := make(map[string]configFieldInfo)

	// Map schema names to actual field names in the Config struct
	// These are the fields that have `mapstructure:"-"` tags and represent configuration sections
	sectionMappings := getSectionMappings()

	configType := reflect.TypeOf(allconfig.Config{})

	for schemaName, fieldName := range sectionMappings {
		if field, found := configType.FieldByName(fieldName); found {
			// Only include fields that have the mapstructure:"-" tag, indicating they are configuration sections
			if tag := field.Tag.Get("mapstructure"); tag == "-" {
				fieldType := field.Type

				// Handle ConfigNamespace types by using the specific config types directly
				if strings.Contains(fieldType.String(), "ConfigNamespace") {
					configType := getConfigNamespaceConfigType(schemaName)
					if configType != nil {
						sections[schemaName] = configFieldInfo{
							Type: configType,
							Tag:  field.Tag,
						}
					}
					continue
				}

				// For pointer types, get the underlying type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}

				// Skip types that can't be safely reflected
				if !canReflectType(fieldType) {
					continue
				}

				sections[schemaName] = configFieldInfo{
					Type: fieldType,
					Tag:  field.Tag,
				}
			}
		}
	}

	return sections
}

// getConfigNamespaceConfigType returns the specific config type for ConfigNamespace fields
// instead of trying to extract from generics which causes reflection issues
func getConfigNamespaceConfigType(schemaName string) reflect.Type {
	switch schemaName {
	case "menus":
		// For menus, use the MenuConfig struct directly instead of a map
		return reflect.TypeOf(navigation.MenuConfig{})

	case "outputformats":
		// For output formats, use the OutputFormatConfig struct directly
		return reflect.TypeOf(output.OutputFormatConfig{})

	case "mediatypes":
		// For media types, use the MediaTypeConfig struct directly
		return reflect.TypeOf(media.MediaTypeConfig{})

	case "contenttypes":
		// For content types, use the ContentTypeConfig struct directly
		return reflect.TypeOf(media.ContentTypeConfig{})

	case "imaging":
		// For imaging, use the ImagingConfig struct directly
		return reflect.TypeOf(images.ImagingConfig{})

	case "cascade":
		// For cascade, use the PageMatcherParamsConfig struct directly
		return reflect.TypeOf(page.PageMatcherParamsConfig{})

	case "segments":
		// For segments, use the SegmentConfig struct directly
		return reflect.TypeOf(segments.SegmentConfig{})
	}

	return nil
}

// generatePageSchemas creates a JSON schema for Hugo page front matter
func generatePageSchemas(rf jsonschema.Reflector, schemaDir string, r *rootCommand, stats *schemaStats) error {
	// Generate a single schema for Hugo page front matter
	schema := rf.Reflect(pagemeta.PageConfig{})

	// Set schema metadata
	schema.ID = jsonschema.ID("https://gohugo.io/jsonschemas/hugo-page.schema.json")
	schema.Version = "https://json-schema.org/draft-07/schema"
	schema.Title = "Hugo Page Front Matter Schema"
	schema.Description = "JSON Schema for Hugo page front matter structure"

	// Clear required fields
	schema.Required = nil

	// Remove required fields from nested schemas
	removeRequiredFromSchema(schema)

	// Add documentation URLs to front matter properties
	addDocumentationLinksToSchema(schema, "page-frontmatter", []string{})

	// Add flexible date definitions
	addFlexibleDateDefinitions(schema)

	// Remove the content property as it represents the file content after frontmatter, not a frontmatter field
	if schema.Properties != nil {
		schema.Properties.Delete("content")
	}

	// Write schema file
	filename := filepath.Join(schemaDir, "hugo-page.schema.json")
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

// transformEnumsToOneOf converts enum arrays to oneOf patterns following SchemaStore recommendations
// This transforms { "enum": ["foo", "bar"] } into { "oneOf": [{"const": "foo"}, {"const": "bar"}] }
// It also handles anyOf patterns that contain enums
func transformEnumsToOneOf(schema *jsonschema.Schema) {
	if schema == nil {
		return
	}

	// Transform enum to oneOf if present
	if len(schema.Enum) > 0 {
		var oneOfSchemas []*jsonschema.Schema
		for _, enumValue := range schema.Enum {
			oneOfSchemas = append(oneOfSchemas, &jsonschema.Schema{
				Const: enumValue,
			})
		}
		schema.OneOf = oneOfSchemas
		schema.Enum = nil // Clear the enum field
	}

	// Transform enums within anyOf patterns
	if schema.AnyOf != nil {
		for _, anyOfSchema := range schema.AnyOf {
			if len(anyOfSchema.Enum) > 0 {
				var oneOfSchemas []*jsonschema.Schema
				for _, enumValue := range anyOfSchema.Enum {
					oneOfSchemas = append(oneOfSchemas, &jsonschema.Schema{
						Const: enumValue,
					})
				}
				anyOfSchema.OneOf = oneOfSchemas
				anyOfSchema.Enum = nil // Clear the enum field
			}
		}
	}

	// Recursively transform properties
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			transformEnumsToOneOf(pair.Value)
		}
	}

	// Recursively transform pattern properties
	if schema.PatternProperties != nil {
		for _, propSchema := range schema.PatternProperties {
			transformEnumsToOneOf(propSchema)
		}
	}

	// Recursively transform additional properties
	if schema.AdditionalProperties != nil {
		transformEnumsToOneOf(schema.AdditionalProperties)
	}

	// Recursively transform items (for arrays)
	if schema.Items != nil {
		transformEnumsToOneOf(schema.Items)
	}

	// Recursively transform all conditional schemas
	for _, condSchema := range schema.AllOf {
		transformEnumsToOneOf(condSchema)
	}
	for _, condSchema := range schema.AnyOf {
		transformEnumsToOneOf(condSchema)
	}
	for _, condSchema := range schema.OneOf {
		transformEnumsToOneOf(condSchema)
	}

	// Transform not schema
	if schema.Not != nil {
		transformEnumsToOneOf(schema.Not)
	}

	// Transform definitions
	if schema.Definitions != nil {
		for _, defSchema := range schema.Definitions {
			transformEnumsToOneOf(defSchema)
		}
	}
}

// addDocumentationLinksToSchema enhances schema descriptions with documentation URLs
// using anchor links extracted from the built Hugo documentation
func addDocumentationLinksToSchema(schema *jsonschema.Schema, sectionName string, propertyPath []string) {
	if schema == nil {
		return
	}

	// Add documentation to properties
	if schema.Properties != nil {
		for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
			propName := pair.Key
			propSchema := pair.Value

			// Build the current property path
			currentPath := append(propertyPath, propName)

			// Generate documentation URL for this property
			docURL := generateDocumentationURL(sectionName, currentPath)
			if docURL != "" {
				if propSchema.Description != "" {
					// Append URL to existing Go comment description
					propSchema.Description = propSchema.Description + " \n" + docURL
				} else {
					// Add URL as description if no Go comment exists
					propSchema.Description = docURL
				}
			}

			// Recursively process nested schemas
			addDocumentationLinksToSchema(propSchema, sectionName, currentPath)
		}
	}

	// Process additional properties
	if schema.AdditionalProperties != nil {
		addDocumentationLinksToSchema(schema.AdditionalProperties, sectionName, propertyPath)
	}

	// Process array items
	if schema.Items != nil {
		addDocumentationLinksToSchema(schema.Items, sectionName, propertyPath)
	}

	// Process conditional schemas
	for _, condSchema := range schema.AllOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath)
	}
	for _, condSchema := range schema.AnyOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath)
	}
	for _, condSchema := range schema.OneOf {
		addDocumentationLinksToSchema(condSchema, sectionName, propertyPath)
	}

	// Process not schema
	if schema.Not != nil {
		addDocumentationLinksToSchema(schema.Not, sectionName, propertyPath)
	}

	// Process definitions with appropriate section context
	// This ensures that properties within definitions get proper documentation links
	if schema.Definitions != nil {
		for defName, defSchema := range schema.Definitions {
			// For markup definitions, use the specific subsection name for better URL generation
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
			addDocumentationLinksToSchema(defSchema, defSectionName, []string{})
		}
	}
}

var sectionAnchors map[string]map[string]string
var docAnchors map[string]string

// generateDocumentationURL generates a documentation URL based on the section and property path
func generateDocumentationURL(sectionName string, propertyPath []string) string {
	if len(propertyPath) == 0 {
		return ""
	}

	// Base URL for Hugo configuration documentation
	baseURL := "https://gohugo.io/configuration"

	// For root config properties, use the all.md page with extracted anchors
	if sectionName == "rootconfig" && len(propertyPath) > 0 {
		// Initialize anchors cache once, only for root config
		if docAnchors == nil {
			docAnchors = extractAnchorsFromHTML()
		}

		propertyName := propertyPath[len(propertyPath)-1]

		// Check if we have an anchor for this property
		if anchorID, exists := docAnchors[propertyName]; exists {
			return baseURL + "/all/#" + anchorID
		}

		// Fallback: generate anchor from property name (lowercase)
		return baseURL + "/all/#" + strings.ToLower(propertyName)
	}

	// Special case for frontmatter config: all properties are date-related and should link to #dates
	if sectionName == "frontmatter" {
		return baseURL + "/front-matter/#dates"
	}

	// Special case for page front-matter: use content management documentation with anchor extraction
	if sectionName == "page-frontmatter" {
		baseURL = "https://gohugo.io/content-management/front-matter"

		// Initialize section anchors cache for content-management front-matter
		if sectionAnchors == nil {
			sectionAnchors = make(map[string]map[string]string)
		}

		// Extract anchors from content-management front-matter if not cached
		if _, exists := sectionAnchors["content-frontmatter"]; !exists {
			// Extract anchors from content-management/front-matter page
			anchors := make(map[string]string)

			// Try different base directories
			baseDirs := []string{"docs/public", "public-docs", "."}
			for _, baseDir := range baseDirs {
				path := filepath.Join(baseDir, "content-management", "front-matter", "index.html")
				if content, err := os.ReadFile(path); err == nil {
					re := regexp.MustCompile(`id="([^"]+)"`)
					matches := re.FindAllStringSubmatch(string(content), -1)
					for _, match := range matches {
						if len(match) > 1 {
							anchorID := match[1]
							// Map anchor ID to itself for exact matches
							anchors[anchorID] = anchorID
							// Also map lowercase for case-insensitive matching
							anchors[strings.ToLower(anchorID)] = anchorID
						}
					}
					break
				}
			}
			sectionAnchors["content-frontmatter"] = anchors
		}

		// Use the same logic as standard sections for nested properties
		if len(propertyPath) > 0 {
			if anchors, exists := sectionAnchors["content-frontmatter"]; exists {
				// First try: just the property name (most specific)
				propertyName := propertyPath[len(propertyPath)-1]
				if anchorID, found := anchors[propertyName]; found {
					return baseURL + "#" + anchorID
				}

				// Second try: property name in lowercase
				if anchorID, found := anchors[strings.ToLower(propertyName)]; found {
					return baseURL + "#" + anchorID
				}

				// Third try: exact concatenated paths for known patterns
				if len(propertyPath) > 1 {
					fullPath := strings.ToLower(strings.Join(propertyPath, ""))
					if anchorID, found := anchors[fullPath]; found {
						return baseURL + "#" + anchorID
					}
				}

				// Fourth try: work backwards through the property path to find parent section anchors
				// For sitemap.changeFreq, try: sitemap
				for i := len(propertyPath) - 2; i >= 0; i-- {
					parentName := propertyPath[i]
					if anchorID, found := anchors[parentName]; found {
						return baseURL + "#" + anchorID
					}
					if anchorID, found := anchors[strings.ToLower(parentName)]; found {
						return baseURL + "#" + anchorID
					}
				}
			}
		}

		// Fallback to base URL for page front-matter
		return baseURL
	}

	// Initialize section anchors cache
	if sectionAnchors == nil {
		sectionAnchors = make(map[string]map[string]string)
	}

	// Map section names to their documentation URLs
	sectionURL := getSectionDocumentationURL(baseURL, sectionName)

	// Map some section names to their parent section for anchor extraction
	parentSection := sectionName
	if sectionName == "goldmark" || sectionName == "asciidocext" || sectionName == "highlight" || sectionName == "tableofcontents" {
		parentSection = "markup"
	}

	// For section-specific properties, try to find actual anchors from the HTML
	if len(propertyPath) > 0 {
		// Get anchors for this section if not already cached
		if _, exists := sectionAnchors[parentSection]; !exists {
			sectionAnchors[parentSection] = extractSectionAnchors(parentSection)
		}

		// Try different anchor patterns based on property path
		if anchors, exists := sectionAnchors[parentSection]; exists {
			// Special case mappings for properties that don't follow standard naming
			propertyName := propertyPath[len(propertyPath)-1]
			if parentSection == "markup" && propertyName == "defaultMarkdownHandler" {
				return sectionURL + "#default-handler"
			}

			// First try: just the property name (most specific)
			if anchorID, found := anchors[propertyName]; found {
				return sectionURL + "#" + anchorID
			}

			// Second try: property name in lowercase
			if anchorID, found := anchors[strings.ToLower(propertyName)]; found {
				return sectionURL + "#" + anchorID
			}

			// Third try: exact concatenated paths for known patterns (e.g., "rendererhardwraps" for renderer.hardWraps)
			if len(propertyPath) > 1 {
				fullPath := strings.ToLower(strings.Join(propertyPath, ""))
				if anchorID, found := anchors[fullPath]; found {
					return sectionURL + "#" + anchorID
				}
			}

			// Fourth try: work backwards through the property path to find parent section anchors
			// For extensions.typographer.disable, try: typographer, extensions
			for i := len(propertyPath) - 2; i >= 0; i-- {
				parentName := propertyPath[i]
				if anchorID, found := anchors[parentName]; found {
					return sectionURL + "#" + anchorID
				}
				if anchorID, found := anchors[strings.ToLower(parentName)]; found {
					return sectionURL + "#" + anchorID
				}
			}
		}

		// Fallback: for nested properties, use full property path concatenation
		if len(propertyPath) > 1 {
			// For nested properties like renderer.unsafe, use concatenated path: rendererunsafe
			fullPath := strings.ToLower(strings.Join(propertyPath, ""))
			return sectionURL + "#" + fullPath
		} else {
			// For top-level properties, use kebab-case
			anchor := camelToKebab(propertyPath[0])
			return sectionURL + "#" + anchor
		}
	}

	return sectionURL
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
	// Handle special cases where section name doesn't match directory name
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
	}

	// Try multiple possible base directories in order of preference
	baseDirs := []string{
		"docs/public",    // Most likely location for built docs
		"public-docs",    // Alternative location
		"../public-docs", // If running from subdirectory
		"public",         // Alternative name
	}

	var paths []string
	for _, baseDir := range baseDirs {
		path := filepath.Join(baseDir, "configuration", dirName, "index.html")
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

// getSectionMappings returns a map of schema names to Config struct field names
// These are the fields that have `mapstructure:"-"` tags and represent configuration sections
func getSectionMappings() map[string]string {
	return map[string]string{
		"build":       "Build",
		"caches":      "Caches",
		"httpcache":   "HTTPCache",
		"markup":      "Markup",
		"outputs":     "Outputs",
		"deployment":  "Deployment",
		"module":      "Module",
		"frontmatter": "Frontmatter",
		"minify":      "Minify",
		"permalinks":  "Permalinks",
		"taxonomies":  "Taxonomies",
		"sitemap":     "Sitemap",
		"related":     "Related",
		"server":      "Server",
		"pagination":  "Pagination",
		"page":        "Page",
		"privacy":     "Privacy",
		"security":    "Security",
		"services":    "Services",
		"params":      "Params",
		"languages":   "Languages",
		"uglyurls":    "UglyURLs",
		// ConfigNamespace types - now supported:
		"contenttypes":  "ContentTypes",
		"mediatypes":    "MediaTypes",
		"imaging":       "Imaging",
		"outputformats": "OutputFormats",
		"cascade":       "Cascade",
		"segments":      "Segments",
		"menus":         "Menus",
	}
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
