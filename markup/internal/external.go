package internal

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gohugoio/hugo/markup/converter"
)

func ExternallyRenderContent(
	cfg converter.ProviderConfig,
	ctx converter.DocumentContext,
	content []byte, path string, args []string) []byte {

	logger := cfg.Logger
	cmd := exec.Command(path, args...)
	cmd.Stdin = bytes.NewReader(content)
	var out, cmderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmderr
	err := cmd.Run()
	// Most external helpers exit w/ non-zero exit code only if severe, i.e.
	// halting errors occurred. -> log stderr output regardless of state of err
	for _, item := range strings.Split(cmderr.String(), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			logger.ERROR.Printf("%s: %s", ctx.DocumentName, item)
		}
	}
	if err != nil {
		logger.ERROR.Printf("%s rendering %s: %v", path, ctx.DocumentName, err)
	}

	return normalizeExternalHelperLineFeeds(out.Bytes())
}

// Strips carriage returns from third-party / external processes (useful for Windows)
func normalizeExternalHelperLineFeeds(content []byte) []byte {
	return bytes.Replace(content, []byte("\r"), []byte(""), -1)
}

func GetPythonExecPath() string {
	path, err := exec.LookPath("python")
	if err != nil {
		path, err = exec.LookPath("python.exe")
		if err != nil {
			return ""
		}
	}
	return path
}
