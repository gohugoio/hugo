package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

const (
	HTML_LEAD       = "<"
	YAML_LEAD       = "-"
	YAML_DELIM_UNIX = "---\n"
	YAML_DELIM_DOS  = "---\r\n"
	TOML_LEAD       = "+"
	TOML_DELIM_UNIX = "+++\n"
	TOML_DELIM_DOS  = "+++\r\n"
	JSON_LEAD       = "{"
)

var (
	delims = [][]byte{
		[]byte(YAML_DELIM_UNIX),
		[]byte(YAML_DELIM_DOS),
		[]byte(TOML_DELIM_UNIX),
		[]byte(TOML_DELIM_DOS),
		[]byte(JSON_LEAD),
	}

	unixEnding = []byte("\n")
	dosEnding  = []byte("\r\n")
)

type FrontMatter []byte
type Content []byte

type Page interface {
	FrontMatter() FrontMatter
	Content() Content
	IsRenderable() bool
}

type page struct {
	render      bool
	frontmatter FrontMatter
	content     Content
}

func (p *page) Content() Content {
	return p.content
}

func (p *page) FrontMatter() FrontMatter {
	return p.frontmatter
}

func (p *page) IsRenderable() bool {
	return p.render
}

// ReadFrom reads the content from an io.Reader and constructs a page.
func ReadFrom(r io.Reader) (p Page, err error) {
	reader := bufio.NewReader(r)

	if err = chompWhitespace(reader); err != nil && err != io.EOF {
		return
	}

	firstLine, err := peekLine(reader)
	if err != nil && err != io.EOF {
		return
	}

	newp := new(page)
	newp.render = shouldRender(firstLine)

	if newp.render && isFrontMatterDelim(firstLine) {
		left, right := determineDelims(firstLine)
		fm, err := extractFrontMatterDelims(reader, left, right)
		if err != nil {
			return nil, err
		}
		newp.frontmatter = fm
	}

	content, err := extractContent(reader)
	if err != nil {
		return nil, err
	}

	newp.content = content

	return newp, nil
}

func chompWhitespace(r io.RuneScanner) (err error) {
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(c) {
			r.UnreadRune()
			return nil
		}
	}
}

func peekLine(r *bufio.Reader) (line []byte, err error) {
	firstFive, err := r.Peek(5)
	if err != nil {
		return
	}
	idx := bytes.IndexByte(firstFive, '\n')
	if idx == -1 {
		return firstFive, nil
	}
	idx += 1 // include newline.
	return firstFive[:idx], nil
}

func shouldRender(lead []byte) (frontmatter bool) {
	if len(lead) <= 0 {
		return
	}

	if bytes.Equal(lead[:1], []byte(HTML_LEAD)) {
		return
	}
	return true
}

func isFrontMatterDelim(data []byte) bool {
	for _, d := range delims {
		if bytes.HasPrefix(data, d) {
			return true
		}
	}

	return false
}

func determineDelims(firstLine []byte) (left, right []byte) {
	switch len(firstLine) {
	case 4:
		if firstLine[0] == YAML_LEAD[0] {
			return []byte(YAML_DELIM_UNIX), []byte(YAML_DELIM_UNIX)
		}
		return []byte(TOML_DELIM_UNIX), []byte(TOML_DELIM_UNIX)

	case 5:
		if firstLine[0] == YAML_LEAD[0] {
			return []byte(YAML_DELIM_DOS), []byte(YAML_DELIM_DOS)
		}
		return []byte(TOML_DELIM_DOS), []byte(TOML_DELIM_DOS)
	case 3:
		fallthrough
	case 2:
		fallthrough
	case 1:
		return []byte(JSON_LEAD), []byte("}")
	default:
		panic(fmt.Sprintf("Unable to determine delims from %q", firstLine))
	}
}

func extractFrontMatterDelims(r *bufio.Reader, left, right []byte) (fm FrontMatter, err error) {
	var (
		c         byte
		level     int = 0
		bytesRead int = 0
		sameDelim     = bytes.Equal(left, right)
	)

	wr := new(bytes.Buffer)
	for {
		if c, err = r.ReadByte(); err != nil {
			return nil, fmt.Errorf("Unable to read frontmatter at filepos %d: %s", bytesRead, err)
		}
		bytesRead += 1

		switch c {
		case left[0]:
			var (
				buf       []byte = []byte{c}
				remaining []byte
			)

			if remaining, err = r.Peek(len(left) - 1); err != nil {
				return nil, err
			}

			buf = append(buf, remaining...)

			if bytes.Equal(buf, left) {
				if sameDelim {
					if level == 0 {
						level = 1
					} else {
						level = 0
					}
				} else {
					level += 1
				}
			}

			if _, err = wr.Write([]byte{c}); err != nil {
				return nil, err
			}

			if level == 0 {
				if _, err = r.Read(remaining); err != nil {
					return nil, err
				}
				if _, err = wr.Write(remaining); err != nil {
					return nil, err
				}
			}
		case right[0]:
			match, err := matches(r, wr, []byte{c}, right)
			if err != nil {
				return nil, err
			}
			if match {
				level -= 1
			}
		default:
			if err = wr.WriteByte(c); err != nil {
				return nil, err
			}
		}

		if level == 0 && !unicode.IsSpace(rune(c)) {
			if err = chompWhitespace(r); err != nil {
				if err != io.EOF {
					return nil, err
				}
			}
			return wr.Bytes(), nil
		}
	}
}

func matches_quick(buf, expected []byte) (ok bool, err error) {
	return bytes.Equal(expected, buf), nil
}

func matches(r *bufio.Reader, wr io.Writer, c, expected []byte) (ok bool, err error) {
	if len(expected) == 1 {
		if _, err = wr.Write(c); err != nil {
			return
		}
		return bytes.Equal(c, expected), nil
	}

	buf := make([]byte, len(expected)-1)
	if buf, err = r.Peek(len(expected) - 1); err != nil {
		return
	}

	buf = append(c, buf...)
	return bytes.Equal(expected, buf), nil
}

func extractContent(r io.Reader) (content Content, err error) {
	wr := new(bytes.Buffer)
	if _, err = wr.ReadFrom(r); err != nil {
		return
	}
	return wr.Bytes(), nil
}
