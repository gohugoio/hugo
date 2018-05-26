package tags

import (
	"fmt"
	"html/template"

	"github.com/gohugoio/hugo/deps"
)

// New creates a pointer to a new tags.Namespace struct.
func New(d *deps.Deps) *Namespace {
	return &Namespace{d}
}

type Namespace struct {
	// how to get the source folder path info (e.g. --source YYY).
	*deps.Deps
}

// Script generates a script tag and concatenates the srcs into dest.
//
func (ns *Namespace) Script(dest string, attr map[string]interface{}, srcs ...interface{}) (template.HTML, error) {
	sourceFs := ns.Fs.WorkingDir
	destinationFs := ns.Fs.Destination
	var html template.HTML

	bytes, err := ReadBytes(sourceFs, srcs...)
	if err != nil {
		return html, err
	}

	asset := NewAsset(dest, bytes)
	url := asset.Url(true)

	// TODO (nfisher 2018-01-25): Might need to introduce a write lock to prevent partial/interleaved writes?
	err = Prepare(destinationFs, ns.PublishDir, dest)
	if err != nil {
		return html, err
	}

	err = Write(destinationFs, ns.PublishDir, url, bytes)
	if err != nil {
		return html, err
	}

	var values []interface{}

	values = append(values, url)

	// TODO (nfisher 2018-01-23): map attr into tag.
	var tagBase string
	tagBase = `<script src="%s" integrity="sha256-%s"></script>`
	values = append(values, asset.Base64Sum())

	// TODO: Not sure if values should be considered tainted input.
	return template.HTML(fmt.Sprintf(tagBase, values...)), nil
}

// Style generates a style tag and concatenates the srcs into dest.
func (ns *Namespace) Style(dest string, attr map[string]interface{}, srcs ...interface{}) (template.HTML, error) {
	sourceFs := ns.Fs.WorkingDir
	destinationFs := ns.Fs.Destination
	var html template.HTML

	bytes, err := ReadBytes(sourceFs, srcs...)
	if err != nil {
		return html, err
	}

	asset := NewAsset(dest, bytes)
	url := asset.Url(true)

	// TODO (nfisher 2018-01-25): Might need to introduce a write lock to prevent partial/interleaved writes?
	err = Prepare(destinationFs, ns.PublishDir, dest)
	if err != nil {
		return html, err
	}

	err = Write(destinationFs, ns.PublishDir, url, bytes)
	if err != nil {
		return html, err
	}

	var values []interface{}

	values = append(values, url)

	// TODO (nfisher 2018-01-23): map attr into tag.
	var tagBase string
	tagBase = `<link rel="stylesheet" href="%s" integrity="sha256-%s"/>`
	values = append(values, asset.Base64Sum())

	// TODO: Not sure if values should be considered tainted input.
	return template.HTML(fmt.Sprintf(tagBase, values...)), nil
}
