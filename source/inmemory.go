package source

import (
	"bytes"
	"fmt"
	"path"
)

type ByteSource struct {
	Name    string
	Content []byte
	Section string
}

func (b *ByteSource) String() string {
	return fmt.Sprintf("%s %s %s", b.Name, b.Section, string(b.Content))
}

type InMemorySource struct {
	ByteSource []ByteSource
}

func (i *InMemorySource) Files() (files []*File) {
	files = make([]*File, len(i.ByteSource))
	for i, fake := range i.ByteSource {
		files[i] = &File{
			LogicalName: fake.Name,
			Contents:    bytes.NewReader(fake.Content),
			Section:     fake.Section,
			Dir:         path.Dir(fake.Name),
		}
	}
	return
}
