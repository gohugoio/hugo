---
title: JSON schemas
linkTitle: JSON schemas
description: Use JSON Schema files to validate Hugo configuration files and page frontmatter.
categories: []
keywords: []
weight: 70
---

Hugo provides JSON Schema files for validating Hugo configuration files and page frontmatter. These schemas enable IntelliSense, auto-completion, and validation in code editors that support JSON Schema.

## Available schemas

### Configuration schemas

These schemas validate Hugo configuration files (`hugo.yaml`, `hugo.json`, `hugo.toml`):

- [`hugo-config.schema.json`](hugo-config.schema.json) - Main configuration schema that references all other config schemas
- [`hugo-config-build.schema.json`](hugo-config-build.schema.json) - Build configuration (`build`)
- [`hugo-config-caches.schema.json`](hugo-config-caches.schema.json) - Cache configuration (`caches`)
- [`hugo-config-cascade.schema.json`](hugo-config-cascade.schema.json) - Cascade configuration (`cascade`)
- [`hugo-config-contenttypes.schema.json`](hugo-config-contenttypes.schema.json) - Content types configuration (`contenttypes`)
- [`hugo-config-deployment.schema.json`](hugo-config-deployment.schema.json) - Deployment configuration (`deployment`)
- [`hugo-config-frontmatter.schema.json`](hugo-config-frontmatter.schema.json) - Front matter configuration (`frontmatter`)
- [`hugo-config-httpcache.schema.json`](hugo-config-httpcache.schema.json) - HTTP cache configuration (`httpcache`)
- [`hugo-config-imaging.schema.json`](hugo-config-imaging.schema.json) - Image processing configuration (`imaging`)
- [`hugo-config-languages.schema.json`](hugo-config-languages.schema.json) - Language configuration (`languages`)
- [`hugo-config-markup.schema.json`](hugo-config-markup.schema.json) - Markup configuration (`markup`)
- [`hugo-config-mediatypes.schema.json`](hugo-config-mediatypes.schema.json) - Media types configuration (`mediatypes`)
- [`hugo-config-menus.schema.json`](hugo-config-menus.schema.json) - Menu configuration (`menus`)
- [`hugo-config-minify.schema.json`](hugo-config-minify.schema.json) - Minification configuration (`minify`)
- [`hugo-config-module.schema.json`](hugo-config-module.schema.json) - Module configuration (`module`)
- [`hugo-config-outputformats.schema.json`](hugo-config-outputformats.schema.json) - Output formats configuration (`outputformats`)
- [`hugo-config-outputs.schema.json`](hugo-config-outputs.schema.json) - Outputs configuration (`outputs`)
- [`hugo-config-page.schema.json`](hugo-config-page.schema.json) - Page configuration (`page`)
- [`hugo-config-pagination.schema.json`](hugo-config-pagination.schema.json) - Pagination configuration (`pagination`)
- [`hugo-config-params.schema.json`](hugo-config-params.schema.json) - Parameters configuration (`params`)
- [`hugo-config-permalinks.schema.json`](hugo-config-permalinks.schema.json) - Permalinks configuration (`permalinks`)
- [`hugo-config-privacy.schema.json`](hugo-config-privacy.schema.json) - Privacy configuration (`privacy`)
- [`hugo-config-related.schema.json`](hugo-config-related.schema.json) - Related content configuration (`related`)
- [`hugo-config-security.schema.json`](hugo-config-security.schema.json) - Security configuration (`security`)
- [`hugo-config-segments.schema.json`](hugo-config-segments.schema.json) - Segments configuration (`segments`)
- [`hugo-config-server.schema.json`](hugo-config-server.schema.json) - Server configuration (`server`)
- [`hugo-config-services.schema.json`](hugo-config-services.schema.json) - Services configuration (`services`)
- [`hugo-config-sitemap.schema.json`](hugo-config-sitemap.schema.json) - Sitemap configuration (`sitemap`)
- [`hugo-config-taxonomies.schema.json`](hugo-config-taxonomies.schema.json) - Taxonomies configuration (`taxonomies`)

### Page frontmatter schema

This schema validates Hugo page frontmatter in Markdown files:

- [`hugo-page.schema.json`](hugo-page.schema.json) - Page frontmatter schema with flexible date handling and build options

## Configuration file organization

Hugo supports two approaches for organizing configuration:

### Single configuration file
All configuration in one file: `hugo.yaml`, `hugo.json`, or `hugo.toml`

### Split configuration files
Configuration sections in separate files in the `config/` directory, organized by environment:

```text
config/
├── _default/           # Default configuration (all environments)
│   ├── hugo.yaml       # Main configuration
│   ├── build.yaml      # Build configuration
│   ├── markup.yaml     # Markup and rendering configuration  
│   ├── params.yaml     # Site parameters
│   ├── module.yaml     # Hugo Modules configuration
│   ├── deployment.yaml # Deployment settings
│   ├── languages.yaml  # Multi-language configuration
│   ├── menus.yaml      # Site navigation menus
│   └── caches.yaml     # File cache settings
├── production/         # Production environment overrides
│   ├── params.yaml
│   └── build.yaml
├── staging/           # Staging environment overrides
│   └── params.yaml
└── development/       # Development environment overrides
    └── server.yaml    # Development server settings
```

Each file can use `.yaml`, `.yml`, `.json`, or `.toml` extensions. The individual schemas validate the specific configuration section, providing targeted IntelliSense and validation for each file.

## Using the schemas

### SchemaStore integration

Hugo's JSON schemas are being integrated with [SchemaStore](https://schemastore.org/), the central repository for JSON schemas used by most editors and IDEs. Once available, schemas will be automatically detected by compatible editors without any manual configuration required.

For current editor integrations and setup instructions, visit [schemastore.org](https://schemastore.org/).

### Recommended VS Code extensions

Until Hugo schemas are available in SchemaStore, use these extensions that provide automatic JSON Schema validation:

#### YAML files
- **[YAML Language Support](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)** by Red Hat
  - Automatically uses SchemaStore for schema detection
  - Provides validation, auto-completion, and hover documentation
  - Supports Hugo configuration files when schemas are available in SchemaStore

#### JSON files  
- **Built-in JSON Language Features** (included with VS Code)
  - Automatically uses SchemaStore for JSON schema validation
  - No additional configuration needed once schemas are in SchemaStore

#### TOML files
- **[Even Better TOML](https://marketplace.visualstudio.com/items?itemName=tamasfe.even-better-toml)** by tamasfe
  - Uses SchemaStore for automatic schema detection
  - Provides syntax highlighting, validation, and IntelliSense
  - Supports Hugo's TOML configuration files

#### Markdown files
- **[YAML Language Support](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)** also validates YAML frontmatter in Markdown files

### Manual configuration (temporary)

Until Hugo schemas are available in SchemaStore, you can manually configure VS Code by adding this to your workspace settings (`.vscode/settings.json`):

```json
{
  "yaml.schemas": {
    "https://gohugo.io/jsonschemas/hugo-config.schema.json": [
      "hugo.{yaml,yml}",
      "config.{yaml,yml}",
      "config/**/*.{yaml,yml}"
    ],
    "https://gohugo.io/jsonschemas/hugo-page.schema.json": [
      "content/**/*.md",
      "archetypes/*.md"
    ]
  },
  "json.schemas": [
    {
      "fileMatch": [
        "hugo.json",
        "config.json",
        "config/**/*.json"
      ],
      "url": "https://gohugo.io/jsonschemas/hugo-config.schema.json"
    }
  ]
}
```

### Other editors

Most modern editors with JSON Schema support will automatically use schemas from SchemaStore. For manual configuration, reference these URLs:

- Main configuration: `https://gohugo.io/jsonschemas/hugo-config.schema.json`
- Page frontmatter: `https://gohugo.io/jsonschemas/hugo-page.schema.json`
- Individual config sections: `https://gohugo.io/jsonschemas/hugo-config-{section}.schema.json`

## Generating schemas

The schemas are generated using the `hugo gen jsonschemas` command, which uses reflection to create schemas directly from Hugo's Go source code. This ensures the schemas are always accurate and up-to-date.

```bash
hugo gen jsonschemas --dir docs/content/en/jsonschemas
```

For more information about configuration options, see the [configuration documentation](../configuration/).
