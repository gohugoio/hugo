package internal

import (
	"bytes"
	"strings"

	"github.com/cli/safeexec"
	"github.com/gohugoio/hugo/common/hexec"

	"github.com/gohugoio/hugo/markup/converter"
)

func ExternallyRenderContent(
	cfg converter.ProviderConfig,
	ctx converter.DocumentContext,
	content []byte, path string, args []string) []byte {

	logger := cfg.Logger
	cmd, err := hexec.SafeCommand(path, args...)
	if err != nil {
		logger.Errorf("%s rendering %s: %v", path, ctx.DocumentName, err)
		return nil
	}
	cmd.Stdin = bytes.NewReader(content)
	var out, cmderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmderr
	err = cmd.Run()
	// Most external helpers exit w/ non-zero exit code only if severe, i.e.
	// halting errors occurred. -> log stderr output regardless of state of err
	for _, item := range strings.Split(cmderr.String(), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			logger.Errorf("%s: %s", ctx.DocumentName, item)
		}
	}
	if err != nil {
		logger.Errorf("%s rendering %s: %v", path, ctx.DocumentName, err)
	}

	return normalizeExternalHelperLineFeeds(out.Bytes())
}

// Strips carriage returns from third-party / external processes (useful for Windows)
func normalizeExternalHelperLineFeeds(content []byte) []byte {
	return bytes.Replace(content, []byte("\r"), []byte(""), -1)
}

func GetPythonExecPath() string {
	path, err := safeexec.LookPath("python")
	if err != nil {
		path, err = safeexec.LookPath("python.exe")
		if err != nil {
			return ""
		}
	}
	return path
}
