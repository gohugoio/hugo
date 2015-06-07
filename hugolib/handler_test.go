package hugolib

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
)

func TestDefaultHandler(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.DestinationFS = new(afero.MemMapFs)
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{filepath.FromSlash("sect/doc2.html"), []byte("<!doctype html><html><body>more content</body></html>")},
		{filepath.FromSlash("sect/doc3.md"), []byte("# doc3\n*some* content")},
		{filepath.FromSlash("sect/doc4.md"), []byte("---\ntitle: doc4\n---\n# doc4\n*some content*")},
		{filepath.FromSlash("sect/doc3/img1.png"), []byte("‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚")},
		{filepath.FromSlash("sect/img2.gif"), []byte("GIF89a��€��ÿÿÿ���,�������D�;")},
		{filepath.FromSlash("sect/img2.spf"), []byte("****FAKE-FILETYPE****")},
		{filepath.FromSlash("doc7.html"), []byte("<html><body>doc7 content</body></html>")},
		{filepath.FromSlash("sect/doc8.html"), []byte("---\nmarkup: md\n---\n# title\nsome *content*")},
	}

	viper.Set("DefaultExtension", "html")
	viper.Set("verbose", true)

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		Targets: targetList{Page: &target.PagePub{UglyURLs: true}},
	}

	s.initializeSiteInfo()
	// From site_test.go
	templatePrep(s)

	must(s.addTemplate("_default/single.html", "{{.Content}}"))
	must(s.addTemplate("head", "<head><script src=\"script.js\"></script></head>"))
	must(s.addTemplate("head_abs", "<head><script src=\"/script.js\"></script></head>"))

	// From site_test.go
	createAndRenderPages(t, s)

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("sect/doc1.html"), "\n\n<h1 id=\"title:5d74edbb89ef198cd37882b687940cda\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("sect/doc3.html"), "\n\n<h1 id=\"doc3:28c75a9e2162b8eccda73a1ab9ce80b4\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("sect/doc3/img1.png"), string([]byte("‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚"))},
		{filepath.FromSlash("sect/img2.gif"), string([]byte("GIF89a��€��ÿÿÿ���,�������D�;"))},
		{filepath.FromSlash("sect/img2.spf"), string([]byte("****FAKE-FILETYPE****"))},
		{filepath.FromSlash("doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("sect/doc8.html"), "\n\n<h1 id=\"title:0ae308ad73e2f37bd09874105281b5d8\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		file, err := hugofs.DestinationFS.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target.", test.doc)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}
