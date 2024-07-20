package internal

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/asciidocext/asciidocext_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"golang.org/x/net/html"
)

type AsciidocConverter struct {
	Ctx converter.DocumentContext
	Cfg converter.ProviderConfig
}

type AsciidocResult struct {
	converter.ResultRender
	toc *tableofcontents.Fragments
}

/* ToDo: RelPermalink patch for svg posts not working*/
type pageSubset interface {
	RelPermalink() string
}

func (r AsciidocResult) TableOfContents() *tableofcontents.Fragments {
	return r.toc
}

func (a *AsciidocConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	b, err := a.GetAsciidocContent(ctx.Src, a.Ctx)
	if err != nil {
		return nil, err
	}
	content, toc, err := a.extractTOC(b)
	if err != nil {
		return nil, err
	}
	return AsciidocResult{
		ResultRender: converter.Bytes(content),
		toc:          toc,
	}, nil
}

func (a *AsciidocConverter) Supports(_ identity.Identity) bool {
	return false
}

// GetAsciidocContent calls asciidoctor as an external helper
// to convert AsciiDoc content to HTML.
func (a *AsciidocConverter) GetAsciidocContent(src []byte, ctx converter.DocumentContext) ([]byte, error) {
	if !HasAsciiDoc() {
		a.Cfg.Logger.Errorln("asciidoctor not found in $PATH: Please install.\n",
			"                 Leaving AsciiDoc content unrendered.")
		return src, nil
	}

	args := a.ParseArgs(ctx)
	args = append(args, "-")

	a.Cfg.Logger.Infoln("Rendering", ctx.DocumentName, " using asciidoctor args", args, "...")

	return internal.ExternallyRenderContent(a.Cfg, ctx, src, asciiDocBinaryName, args)
}

func (a *AsciidocConverter) ParseArgs(ctx converter.DocumentContext) []string {
	cfg := a.Cfg.MarkupConfig().AsciidocExt
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

		args = append(args, "-a", attributeKey+"="+attributeValue)
	}

	if cfg.WorkingFolderCurrent {
		contentDir := filepath.Dir(ctx.Filename)
		destinationDir := a.Cfg.Conf.BaseConfig().PublishDir

		if destinationDir == "" {
			a.Cfg.Logger.Errorln("markup.asciidocext.workingFolderCurrent requires hugo command option --destination to be set")
		}

		var outDir string
		var err error

		file := filepath.Base(ctx.Filename)
		if a.Cfg.Conf.IsUglyURLs("") || file == "_index.adoc" || file == "index.adoc" {
			outDir, err = filepath.Abs(filepath.Dir(filepath.Join(destinationDir, ctx.DocumentName)))
		} else {
			postDir := ""
			page, ok := ctx.Document.(pageSubset)
			if ok {
				postDir = filepath.Base(page.RelPermalink())
			} else {
				a.Cfg.Logger.Errorln("unable to cast interface to pageSubset")
			}

			outDir, err = filepath.Abs(filepath.Join(destinationDir, filepath.Dir(ctx.DocumentName), postDir))
		}

		if err != nil {
			a.Cfg.Logger.Errorln("asciidoctor outDir: ", err)
		}

		args = append(args, "--base-dir", contentDir, "-a", "outdir="+outDir)
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

	return args
}

func (a *AsciidocConverter) AppendArg(args []string, option, value, defaultValue string, allowedValues map[string]bool) []string {
	if value != defaultValue {
		if allowedValues[value] {
			args = append(args, option, value)
		} else {
			a.Cfg.Logger.Errorln("Unsupported asciidoctor value `" + value + "` for option " + option + " was passed in and will be ignored.")
		}
	}
	return args
}

const asciiDocBinaryName = "asciidoctor"

func HasAsciiDoc() bool {
	return hexec.InPath(asciiDocBinaryName)
}

// extractTOC extracts the toc from the given src html.
// It returns the html without the TOC, and the TOC data
func (a *AsciidocConverter) extractTOC(src []byte) ([]byte, *tableofcontents.Fragments, error) {
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
			if !a.Cfg.MarkupConfig().AsciidocExt.PreserveTOC {
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
