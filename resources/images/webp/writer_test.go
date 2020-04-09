package webp

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_Encode(t *testing.T) {
	libwebpHome := ensureLibwebpHome(t)

	givenOptions := Options{CwebpBinary: filepath.Join(libwebpHome, "bin", "cwebp")}
	givenImage := testImage()
	expectedImage := givenImage

	actualBuf := new(bytes.Buffer)
	actualErr := Encode(actualBuf, givenImage, &givenOptions)
	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}

	actualImage := decodeWebp(t, actualBuf, "buffer")

	assertYCbCrImageEquals(t, expectedImage, actualImage)
}

func Test_Encode_failsBecauseNoCwebpAvailable(t *testing.T) {
	givenImage := testImage()

	defer setEnvironmentVariable(t, LibwebpHome, "DOES_NOT_EXIST")()
	defer setEnvironmentVariable(t, "PATH", "DOES_NOT_EXIST")()

	actualBuf := new(bytes.Buffer)
	actualErr := Encode(actualBuf, givenImage, nil)
	if actualErr != ErrWebpEncodingNotSupported {
		fatalfWithEnvironment(t, "error %v expected, but got: %v", ErrWebpEncodingNotSupported, actualErr)
	}
}

func Test_encodeToPng(t *testing.T) {
	givenImage := testImage()
	expectedImage := givenImage

	actualFile, actualErr := encodeToTemp(givenImage)
	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}
	defer func() { _ = os.Remove(actualFile) }()

	if !strings.HasPrefix(actualFile, os.TempDir()) {
		fatalfWithEnvironment(t, "expected to start with temporary directory, but got: %v", actualFile)
	}
	actualImage := decodePngFrom(t, actualFile)

	assertRgbaImageEquals(t, expectedImage, actualImage)
}

func Test_lookupCwebpBinary_viaOptions(t *testing.T) {
	libwebpHome := ensureLibwebpHome(t)

	defer setEnvironmentVariable(t, LibwebpHome, "DOES_NOT_EXIST")()
	defer setEnvironmentVariable(t, "PATH", "DOES_NOT_EXIST")()

	actual, actualErr := lookupCwebpBinary(&Options{CwebpBinary: filepath.Join(libwebpHome, "bin", "cwebp")})

	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}
	expected := filepath.Join(libwebpHome, "bin", executableName())
	if actual != expected {
		fatalfWithEnvironment(t, "expected '%s', but got: %s", expected, actual)
	}
}

func Test_lookupCwebpBinary_viaHomeEnv(t *testing.T) {
	libwebpHome := ensureLibwebpHome(t)

	defer setEnvironmentVariable(t, LibwebpHome, libwebpHome)()
	defer setEnvironmentVariable(t, "PATH", "DOES_NOT_EXIST")()
	actual, actualErr := lookupCwebpBinary(nil)

	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}
	expected := filepath.Join(libwebpHome, "bin", executableName())
	if actual != expected {
		fatalfWithEnvironment(t, "expected '%s', but got: %s", expected, actual)
	}
}

func Test_lookupCwebpBinary_viaCwd(t *testing.T) {
	libwebpHome := ensureLibwebpHome(t)

	defer setEnvironmentVariable(t, LibwebpHome, "DOES_NOT_EXIST")()
	defer setEnvironmentVariable(t, "PATH", "DOES_NOT_EXIST")()
	defer setCwd(t, filepath.Join(libwebpHome, "bin"))()
	actual, actualErr := lookupCwebpBinary(nil)

	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}
	expected := filepath.Join(libwebpHome, "bin", executableName())
	if !strings.HasSuffix(actual, expected) {
		fatalfWithEnvironment(t, "expects to start with '%s', but got: %s", expected, actual)
	}
}

func Test_lookupCwebpBinary_viaPath(t *testing.T) {
	libwebpHome := ensureLibwebpHome(t)

	defer setEnvironmentVariable(t, LibwebpHome, "DOES_NOT_EXIST")()
	defer setEnvironmentVariable(t, "PATH", filepath.Join(libwebpHome, "bin"))()
	actual, actualErr := lookupCwebpBinary(nil)

	if actualErr != nil {
		fatalfWithEnvironment(t, "no error expected, but got: %v", actualErr)
	}
	expected := filepath.Join(libwebpHome, "bin", executableName())
	if actual != expected {
		fatalfWithEnvironment(t, "expected '%s', but got: %s", expected, actual)
	}
}

func executableName() string {
	switch runtime.GOOS {
	case "windows":
		return "cwebp.exe"
	default:
		return "cwebp"
	}
}
