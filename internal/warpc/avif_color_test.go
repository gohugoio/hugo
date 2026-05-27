package warpc_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestAvifColorPropertyPreservation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	// dock-75-hdr.avif is Lightroom-style SDR+gain-map HDR; the decoder bakes the
	// gain map into a single BT.2020/PQ HDR AVIF.
	files := `
-- hugo.toml --
-- assets/dock.avif --
sourcefilename: ../../resources/testdata/bep/dock-75-hdr.avif
-- layouts/home.html --
{{ $dock := resources.Get "dock.avif" }}
Dock: {{ $dock.Width }}x{{ $dock.Height }}
{{ $dockReencoded := $dock.Process "avif q75" }}
DockReencoded: {{ $dockReencoded.RelPermalink }}
`

	b, err := hugolib.TestE(t, files, hugolib.TestOptWithOSFs())
	b.Assert(err, qt.IsNil)
	b.AssertFileContent("public/index.html", "Dock: 1024x683")

	// Check color properties of re-encoded files using exiftool
	checkColorProps := func(t *testing.T, name, filepath string, expectedPrimaries, expectedTransfer string) {
		cmd := exec.Command("exiftool", "-ColorPrimaries", "-TransferCharacteristics", filepath)
		output, err := cmd.Output()
		if err != nil {
			t.Skipf("exiftool not available: %v", err)
		}
		outputStr := string(output)
		t.Logf("Color properties for %s:\n%s", name, outputStr)

		if !strings.Contains(outputStr, expectedPrimaries) {
			t.Errorf("%s: Expected color primaries %q, got:\n%s", name, expectedPrimaries, outputStr)
		}
		if !strings.Contains(outputStr, expectedTransfer) {
			t.Errorf("%s: Expected transfer characteristics %q, got:\n%s", name, expectedTransfer, outputStr)
		}
	}

	// Find the generated files - use Cfg.WorkingDir
	publicDir := filepath.Join(b.Cfg.WorkingDir, "public")
	t.Logf("Looking in: %s", publicDir)

	// dock-75-hdr.avif has a gain map, so the output should be BT.2020/PQ.
	dockMatches, _ := filepath.Glob(filepath.Join(publicDir, "dock_hu*.avif"))
	if len(dockMatches) > 0 {
		checkColorProps(t, "dock", dockMatches[0], "BT.2020", "PQ")
	} else {
		t.Error("No dock AVIF file found in output")
	}
}
