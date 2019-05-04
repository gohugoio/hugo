package commands

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func captureStdout(f func() (*cobra.Command, error)) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_, err := f()

	if err != nil {
		return "", err
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), nil
}

func TestListAll(t *testing.T) {
	assert := require.New(t)
	dir, err := createSimpleTestSite(t, testSiteConfig{})

	assert.NoError(err)

	hugoCmd := newCommandsBuilder().addAll().build()
	cmd := hugoCmd.getCommand()

	defer func() {
		os.RemoveAll(dir)
	}()

	cmd.SetArgs([]string{"-s=" + dir, "list", "all"})

	out, err := captureStdout(cmd.ExecuteC)
	assert.NoError(err)

	r := csv.NewReader(strings.NewReader(out))

	header, err := r.Read()
	assert.NoError(err)

	assert.Equal([]string{
		"path", "slug", "title",
		"date", "expiryDate", "publishDate",
		"draft", "permalink",
	}, header)

	record, err := r.Read()
	assert.Equal([]string{
		filepath.Join("content", "p1.md"), "", "P1",
		"0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z",
		"false", "https://example.org/p1/",
	}, record)
}
