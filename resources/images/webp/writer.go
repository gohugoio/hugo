package webp

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LibwebpHome = "LIBWEBP_HOME"
)

var (
	ErrWebpEncodingNotSupported = errors.New("webp encoding not supported; no valid cwebp binary available;" +
		" neither by " + LibwebpHome + " environment variable" +
		" nor in PATH" +
		" nor in current working directory" +
		" nor directly configured")
)

type Options struct {
	// Quality between 1..100
	// If not present or out of range 80 will be used as default
	Quality int

	// Defines the cwebp binary explicitly.
	// If not present it will looked in the environment.
	CwebpBinary string
}

func Encode(w io.Writer, m image.Image, o *Options) (err error) {
	if o == nil {
		o = &Options{}
	}
	if o.Quality < 1 || o.Quality > 100 {
		o.Quality = 80
	}

	input, eErr := encodeToTemp(m)
	if eErr != nil {
		err = eErr
		return
	}
	defer func() {
		if rErr := os.Remove(input); rErr != nil && err == nil {
			err = fmt.Errorf("cannot remove temporary PNG file (%s) for webp conversion: %v", input, err)
		}
	}()

	cwebpBinary, lErr := lookupCwebpBinary(o)
	if lErr != nil {
		err = lErr
		return
	}

	if eErr := executeCwebp(cwebpBinary, input, w, o); eErr != nil {
		err = eErr
		return
	}

	return
}

func executeCwebp(binary string, from string, to io.Writer, o *Options) error {
	cmd := exec.Command(binary,
		"-q", strconv.FormatInt(int64(o.Quality), 10), // Quality
		"-quiet",  // Suppress useless output, we just want to see here real errors
		"-o", "-", // Output to stdout
		"--", from, // Input from file
	)
	var stderr bytes.Buffer
	cmd.Stdout = to
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot encode webp: %v; stderr: %s", err, strings.TrimSpace(stderr.String()))
	}

	if stderr.Len() > 0 {
		return fmt.Errorf("encode of webp failed: %s", strings.TrimSpace(stderr.String()))
	}

	return nil
}

func encodeToTemp(m image.Image) (filename string, err error) {
	var f *os.File
	if f, err = ioutil.TempFile("", "hugo-webp-input-*.png"); err != nil {
		err = fmt.Errorf("cannot create temporary file for webp conversion: %v", err)
		return
	}
	defer func() {
		if rErr := f.Close(); rErr != nil {
			if err == nil {
				err = fmt.Errorf("cannot close temporary PNG file (%s) for webp conversion: %v", f.Name(), err)
			}
		}
	}()
	if err = png.Encode(f, m); err != nil {
		_ = f.Close()
		err = fmt.Errorf("cannot encode temporary PNG file (%s) for webp conversion: %v", f.Name(), err)
		return
	}
	filename = f.Name()
	return
}

func lookupCwebpBinary(o *Options) (string, error) {
	if o != nil && o.CwebpBinary != "" {
		if found, err := exec.LookPath(o.CwebpBinary); isExecNotFound(err) {
			return "", ErrWebpEncodingNotSupported
		} else {
			return found, err
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		if found, err := exec.LookPath(filepath.Join(cwd, "cwebp")); err == nil {
			return found, nil
		} else if !isExecNotFound(err) {
			return "", fmt.Errorf("cannot lookup webp executable: %v", err)
		}
	}

	if webpHome := os.Getenv(LibwebpHome); webpHome != "" {
		if found, err := exec.LookPath(filepath.Join(webpHome, "bin", "cwebp")); err == nil {
			return found, nil
		} else if !isExecNotFound(err) {
			return "", fmt.Errorf("cannot lookup webp executable: %v", err)
		}
	}

	if found, err := exec.LookPath("cwebp"); err == nil {
		return found, nil
	} else if isExecNotFound(err) {
		return "", ErrWebpEncodingNotSupported
	} else {
		return "", fmt.Errorf("cannot lookup webp executable: %v", err)
	}
}

func isExecNotFound(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if eErr, ok := err.(*exec.Error); ok {
		unwrapped := eErr.Unwrap()
		return unwrapped == exec.ErrNotFound || os.IsNotExist(unwrapped)
	}
	return false
}
