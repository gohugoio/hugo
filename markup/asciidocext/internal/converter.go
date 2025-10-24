package internal

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/markup/asciidocext/asciidocext_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/spf13/cast"
	"golang.org/x/net/html"
)

type AsciiDocConverter struct {
	Ctx converter.DocumentContext
	Cfg converter.ProviderConfig
}

type AsciiDocResult struct {
	converter.ResultRender
	toc *tableofcontents.Fragments
}

type pageSubset interface {
	IsPage() bool
	RelPermalink() string
	Section() string
}

func (r AsciiDocResult) TableOfContents() *tableofcontents.Fragments {
	return r.toc
}

func (a *AsciiDocConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	b, err := a.GetAsciiDocContent(ctx.Src, a.Ctx)
	if err != nil {
		return nil, err
	}
	content, toc, err := a.extractTOC(b)
	if err != nil {
		return nil, err
	}
	return AsciiDocResult{
		ResultRender: converter.Bytes(content),
		toc:          toc,
	}, nil
}

func (a *AsciiDocConverter) Supports(_ identity.Identity) bool {
	return false
}

// GetAsciiDocContent calls asciidoctor as an external helper
// to convert AsciiDoc content to HTML.
func (a *AsciiDocConverter) GetAsciiDocContent(src []byte, ctx converter.DocumentContext) ([]byte, error) {
	if ok, err := HasAsciiDoc(); !ok {
		a.Cfg.Logger.Errorf("%s: %s", err.Error(), "leaving AsciiDoc content unrendered")
		return src, nil
	}

	args, err := a.ParseArgs(ctx)
	if err != nil {
		return nil, err
	}
	args = append(args, "-")

	a.Cfg.Logger.Infoln("Rendering", ctx.DocumentName, " using asciidoctor args", args, "...")

	return internal.ExternallyRenderContent(a.Cfg, ctx, src, AsciiDocBinaryName, args)
}

func (a *AsciiDocConverter) ParseArgs(ctx converter.DocumentContext) ([]string, error) {
	cfg := a.Cfg.MarkupConfig().AsciiDocExt
	args := []string{}

	args = a.AppendArg(args, "-b", cfg.Backend, asciidocext_config.CliDefault.Backend, asciidocext_config.AllowedBackend)

	for _, extension := range cfg.Extensions {
		if strings.LastIndexAny(extension, `\/.`) > -1 {
			a.Cfg.Logger.Errorln("Unsupported asciidoctor extension was passed in. Extension `" + extension + "` ignored. Only installed asciidoctor extensions are allowed.")
			continue
		}
		args = append(args, "-r", extension)
	}

	for attributeKey, attributeValue := range cfg.Attributes {
		if asciidocext_config.DisallowedAttributes[attributeKey] {
			a.Cfg.Logger.Errorln("Unsupported asciidoctor attribute was passed in. Attribute `" + attributeKey + "` ignored.")
			continue
		}

		// To set a document attribute to true: -a attributeKey
		// To set a document attribute to false: -a '!attributeKey'
		// For other types: -a attributeKey=attributeValue
		if b, ok := attributeValue.(bool); ok {
			arg := attributeKey
			if !b {
				arg = "'!" + attributeKey + "'"
			}
			args = append(args, "-a", arg)
		} else {
			args = append(args, "-a", attributeKey+"="+cast.ToString(attributeValue))
		}
	}

	if cfg.WorkingFolderCurrent {
		page, ok := ctx.Document.(pageSubset)
		if !ok {
			return nil, fmt.Errorf("expected pageSubset, got %T", ctx.Document)
		}

		// Derive the outdir document attribute from the relative permalink.
		relPath := strings.TrimPrefix(page.RelPermalink(), a.Cfg.Conf.BaseURL().BasePathNoTrailingSlash)
		relPath, err := url.PathUnescape(relPath)
		if err != nil {
			return nil, err
		}

		if a.Cfg.Conf.IsMultihost() {
			// In a multi-host configuration, neither absolute nor relative
			// permalinks include the language key; prepend it.
			language, ok := a.Cfg.Conf.Language().(*langs.Language)
			if !ok {
				return nil, fmt.Errorf("expected *langs.Language, got %T", a.Cfg.Conf.Language())
			}
			relPath = filepath.Join(language.Lang, relPath)
		}

		if a.Cfg.Conf.IsUglyURLs(page.Section()) {
			if page.IsPage() {
				// Remove the extension.
				relPath = strings.TrimSuffix(relPath, filepath.Ext(relPath))
			} else {
				// Remove the file name.
				relPath = filepath.Dir(relPath)
			}

			// Set imagesoutdir and imagesdir attributes.
			imagesoutdir, err := filepath.Abs(filepath.Join(a.Cfg.Conf.BaseConfig().PublishDir, relPath))
			if err != nil {
				return nil, err
			}
			imagesdir := filepath.Base(imagesoutdir)

			if page.IsPage() {
				args = append(args, "-a", "imagesoutdir="+imagesoutdir, "-a", "imagesdir@="+imagesdir)
			} else {
				args = append(args, "-a", "imagesoutdir="+imagesoutdir)
			}
		}
		// Prepend the publishDir.
		outDir, err := filepath.Abs(filepath.Join(a.Cfg.Conf.BaseConfig().PublishDir, relPath))
		if err != nil {
			return nil, err
		}

		args = append(args, "--base-dir", filepath.Dir(ctx.Filename), "-a", "outdir="+outDir)
	}

	if cfg.NoHeaderOrFooter {
		args = append(args, "--no-header-footer")
	} else {
		a.Cfg.Logger.Warnln("asciidoctor parameter NoHeaderOrFooter is expected for correct html rendering")
	}

	if cfg.SectionNumbers {
		args = append(args, "--section-numbers")
	}

	if cfg.Verbose {
		args = append(args, "--verbose")
	}

	if cfg.Trace {
		args = append(args, "--trace")
	}

	args = a.AppendArg(args, "--failure-level", cfg.FailureLevel, asciidocext_config.CliDefault.FailureLevel, asciidocext_config.AllowedFailureLevel)

	args = a.AppendArg(args, "--safe-mode", cfg.SafeMode, asciidocext_config.CliDefault.SafeMode, asciidocext_config.AllowedSafeMode)

	return args, nil
}

func (a *AsciiDocConverter) AppendArg(args []string, option, value, defaultValue string, allowedValues map[string]bool) []string {
	if value != defaultValue {
		if allowedValues[value] {
			args = append(args, option, value)
		} else {
			a.Cfg.Logger.Errorln("Unsupported asciidoctor value `" + value + "` for option " + option + " was passed in and will be ignored.")
		}
	}
	return args
}

const (
	// AsciiDocBinaryName is the executable name for the AsciiDoc converter.
	AsciiDocBinaryName = "asciidoctor"
)

// HasAsciiDoc reports whether the AsciiDoc converter is installed.
func HasAsciiDoc() (bool, error) {
	if !hexec.InPath(AsciiDocBinaryName) {
		return false, fmt.Errorf("the AsciiDoc converter (%s) is not installed", AsciiDocBinaryName)
	}
	return true, nil
}

// CanRenderDitaaDiagrams reports whether the AsciiDoc converter can render
// Ditaa diagrams. Only used in tests.
func CanRenderDitaaDiagrams() (bool, error) {
	const (
		// gemBinaryName is the executable name for the RubyGems CLI.
		gemBinaryName = "gem"
		// javaBinaryName is the executable name for the Java Runtime Environment CLI.
		javaBinaryName = "java"
	)

	// Verify that the AsciiDoc converter is installed.
	if ok, err := HasAsciiDoc(); !ok {
		return false, err
	}

	// Verify that the RubyGems CLI is installed.
	if !hexec.InPath(gemBinaryName) {
		return false, fmt.Errorf("the RubyGems CLI (%s) is not installed", gemBinaryName)
	}

	// Verify that the required AsciiDoc converter extensions are installed.
	extensions := []string{"asciidoctor-diagram", "asciidoctor-diagram-ditaamini"}

	sc := security.DefaultConfig
	sc.Exec.Allow = security.MustNewWhitelist(gemBinaryName)
	ex := hexec.New(sc, "", loggers.NewDefault())

	for _, extension := range extensions {
		args := []any{"list", extension, "--installed"}
		cmd, err := ex.New(gemBinaryName, args...)
		if err != nil {
			return false, err
		}
		err = cmd.Run()
		if err != nil {
			return false, fmt.Errorf("the %s gem is not installed", extension)
		}
	}

	// Verify that the Java Runtime Environment CLI is installed.
	if !hexec.InPath(javaBinaryName) {
		return false, fmt.Errorf("the Java Runtime Environment CLI (%s) is not installed", javaBinaryName)
	}

	return true, nil
}

// extractTOC extracts the toc from the given src html.
// It returns the html without the TOC, and the TOC data
func (a *AsciiDocConverter) extractTOC(src []byte) ([]byte, *tableofcontents.Fragments, error) {
	var buf bytes.Buffer
	buf.Write(src)
	node, err := html.Parse(&buf)
	if err != nil {
		return nil, nil, err
	}
	var (
		f       func(*html.Node) bool
		toc     *tableofcontents.Fragments
		toVisit []*html.Node
	)
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" && attr(n, "id") == "toc" {
			toc = parseTOC(n)
			if !a.Cfg.MarkupConfig().AsciiDocExt.PreserveTOC {
				n.Parent.RemoveChild(n)
			}
			return true
		}
		if n.FirstChild != nil {
			toVisit = append(toVisit, n.FirstChild)
		}
		if n.NextSibling != nil && f(n.NextSibling) {
			return true
		}
		for len(toVisit) > 0 {
			nv := toVisit[0]
			toVisit = toVisit[1:]
			if f(nv) {
				return true
			}
		}
		return false
	}
	f(node)
	if err != nil {
		return nil, nil, err
	}
	buf.Reset()
	err = html.Render(&buf, node)
	if err != nil {
		return nil, nil, err
	}
	// ltrim <html><head></head><body> and rtrim </body></html> which are added by html.Render
	res := buf.Bytes()[25:]
	res = res[:len(res)-14]
	return res, toc, nil
}

// parseTOC returns a TOC root from the given toc Node
func parseTOC(doc *html.Node) *tableofcontents.Fragments {
	var (
		toc tableofcontents.Builder
		f   func(*html.Node, int, int)
	)
	f = func(n *html.Node, row, level int) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "ul":
				if level == 0 {
					row++
				}
				level++
				f(n.FirstChild, row, level)
			case "li":
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.ElementNode || c.Data != "a" {
						continue
					}
					href := attr(c, "href")[1:]
					toc.AddAt(&tableofcontents.Heading{
						Title: nodeContent(c),
						ID:    href,
						Level: level + 1,
					}, row, level)
				}
				f(n.FirstChild, row, level)
			}
		}
		if n.NextSibling != nil {
			f(n.NextSibling, row, level)
		}
	}
	f(doc.FirstChild, -1, 0)
	return toc.Build()
}

func attr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func nodeContent(node *html.Node) string {
	var buf bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}
