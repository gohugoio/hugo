package webp

import (
	"fmt"
	"golang.org/x/image/webp"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

const (
	colorTolerance = 20
)

func setEnvironmentVariable(t *testing.T, name, to string) func() {
	oldHome := os.Getenv(name)
	if err := os.Setenv(name, to); err != nil {
		fatalfWithEnvironment(t, "cannot set environment %s variable: %v", name, err)
	}
	return func() {
		if oldHome == "" {
			if err := os.Unsetenv(name); err != nil {
				fatalfWithEnvironment(t, "cannot unset %s environment variable: %v", name, err)
			}
		} else {
			if err := os.Setenv(name, oldHome); err != nil {
				fatalfWithEnvironment(t, "cannot reset %s environment variable: %v", name, err)
			}
		}
	}
}

func setCwd(t *testing.T, to string) func() {
	oldHome, err := os.Getwd()
	if err != nil {
		fatalfWithEnvironment(t, "cannot get cwd: %v", err)
	}
	if err := os.Chdir(to); err != nil {
		fatalfWithEnvironment(t, "cannot set cwd: %v", err)
	}
	return func() {
		if err := os.Chdir(oldHome); err != nil {
			fatalfWithEnvironment(t, "cannot reset cwd: %v", err)
		}
	}
}

func decodePngFrom(t *testing.T, file string) image.Image {
	fp, err := os.Open(file)
	if err != nil {
		t.Fatalf("cannot open PNG file %s: %v", file, err)
	}
	defer func() { _ = fp.Close() }()
	return decodePng(t, fp, file)
}

func decodePng(t *testing.T, r io.Reader, name string) image.Image {
	i, err := png.Decode(r)
	if err != nil {
		t.Fatalf("cannot decode PNG file %s: %v", name, err)
	}
	return i
}

func decodeWebp(t *testing.T, r io.Reader, name string) image.Image {
	i, err := webp.Decode(r)
	if err != nil {
		t.Fatalf("cannot decode WEBP file %s: %v", name, err)
	}
	return i
}

func testImage() *image.RGBA {
	i := image.NewRGBA(image.Rect(0, 0, 100, 100))
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
	draw.Draw(i, i.Bounds(), image.NewUniform(blue), image.Point{}, draw.Src)
	return i
}

func assertRgbaImageEquals(t *testing.T, expected *image.RGBA, actual image.Image) {
	actualRgba, actualIsRgba := actual.(*image.RGBA)
	if !actualIsRgba {
		t.Fatalf("image is expected to be of type *image.RGBA, but got: %v", reflect.TypeOf(actual))
	}
	if expected.Rect != actualRgba.Rect {
		t.Fatalf("image is expected to have bounds %v, but got: %v", expected.Rect, actualRgba.Rect)
	}
	if !reflect.DeepEqual(expected, actualRgba) {
		t.Fatalf("image is not as expected")
	}
}

// This test is not perfect, but we do not test if webp itself works, we just test if webp produces a valid
// image with the same rect and at least with "similar" pixels. We assume that webp works.
func assertYCbCrImageEquals(t *testing.T, expected *image.RGBA, actual image.Image) {
	actualYCbCr, actualIsYCbCr := actual.(*image.YCbCr)
	if !actualIsYCbCr {
		t.Fatalf("image is expected to be of type *image.YCbCr, but got: %v", reflect.TypeOf(actual))
	}
	if expected.Rect != actualYCbCr.Rect {
		t.Fatalf("image is expected to have bounds %v, but got: %v", expected.Rect, actualYCbCr.Rect)
	}

	for y := actualYCbCr.Rect.Min.Y; y < actualYCbCr.Rect.Max.Y; y++ {
		for x := actualYCbCr.Rect.Min.X; x < actualYCbCr.Rect.Max.X; x++ {
			expectedColor := expected.RGBAAt(x, y)
			actualYCbCrColor := actualYCbCr.YCbCrAt(x, y)
			r, g, b := color.YCbCrToRGB(actualYCbCrColor.Y, actualYCbCrColor.Cb, actualYCbCrColor.Cr)
			actualColor := color.RGBA{R: r, G: g, B: b, A: 255}

			if !isColorPartSimilarEnough(expectedColor.R, actualColor.R) {
				t.Fatalf("image's %dx%d pixel's red is expected to be %d, but got: %d", x, y, expectedColor.R, actualColor.R)
			}
			if !isColorPartSimilarEnough(expectedColor.G, actualColor.G) {
				t.Fatalf("image's %dx%d pixel's green is expected to be %d, but got: %d", x, y, expectedColor.G, actualColor.G)
			}
			if !isColorPartSimilarEnough(expectedColor.B, actualColor.B) {
				t.Fatalf("image's %dx%d pixel's blue is expected to be %d, but got: %d", x, y, expectedColor.B, actualColor.B)
			}
			if !isColorPartSimilarEnough(expectedColor.A, actualColor.A) {
				t.Fatalf("image's %dx%d pixel's alpha is expected to be %d, but got: %d", x, y, expectedColor.A, actualColor.A)
			}
		}
	}
}

func isColorPartSimilarEnough(expected, actual uint8) bool {
	return expected-colorTolerance < actual || expected+colorTolerance > actual
}

func fatalfWithEnvironment(t *testing.T, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	t.Fatalf("%s\nPlatform: %s", message, retrievePlatformInformation())
}

func retrievePlatformInformation() string {
	switch runtime.GOOS {
	case "darwin":
		return executeAndDump(nil, "sw_vers")
	case "windows":
		return executeAndDump(regexp.MustCompile(`^OS (Name|Version):\s+(.+)`), "systeminfo")
	case "linux":
		return fmt.Sprintf("KERNEL=%q, %s",
			executeAndDump(nil, "uname", "-a"),
			executeAndDump(regexp.MustCompile(`^(NAME|VERSION)\s*=\s*"?([^"]+)"?`), "cat", "/etc/os-release"),
		)
	default:
		return unknownPlatformInformation("unknown operating system")
	}
}

func unknownPlatformInformation(reasonFormat string, args ...interface{}) string {
	return fmt.Sprintf("[%s] %s", runtime.GOOS, fmt.Sprintf(reasonFormat, args...))
}

func executeAndDump(acceptLinesOnlyIfMatches *regexp.Regexp, cmd string, args ...string) string {
	executable, err := exec.LookPath(cmd)
	if err != nil {
		return unknownPlatformInformation("cannot lookup '%s' in path: %v", cmd, err)
	}
	command := exec.Command(executable, args...)
	stdout, err := command.Output()
	if err != nil {
		return unknownPlatformInformation("cannot excute %q: %v", command, err)
	}
	result := ""
	for _, line := range strings.Split(string(stdout), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && (acceptLinesOnlyIfMatches == nil || acceptLinesOnlyIfMatches.MatchString(line)) {
			if result != "" {
				result += ", "
			}
			if acceptLinesOnlyIfMatches != nil {
				result += fmt.Sprintf("%s=%q",
					acceptLinesOnlyIfMatches.ReplaceAllString(line, "$1"),
					acceptLinesOnlyIfMatches.ReplaceAllString(line, "$2"),
				)
			} else {
				result += line
			}
		}
	}
	return result
}
