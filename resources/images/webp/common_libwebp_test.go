package webp

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	// This version is required because the build agents of Travis are currently (2020-04-14) still on macOS 10.13.6
	// and later version of libwebp are not already supported.
	libwebpVersion             = "1.0.1"
	libwebpDownloadBase        = "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/"
	libwebpDownloadFileLinux   = "libwebp-{version}-linux-x86-64.tar.gz"
	libwebpDownloadFileMacOs   = "libwebp-{version}-mac-10.13.tar.gz"
	libwebpDownloadFileWindows = "libwebp-{version}-windows-x64-no-wic.zip"
)

func ensureLibwebpHome(t *testing.T) (to string) {
	from := libwebpDownloadUrl(t)
	to = toPath(from)

	if _, err := os.Stat(to); os.IsNotExist(err) {
		t.Logf("%s does not exist; downloading it...", to)
	} else if err != nil {
		fatalfWithEnvironment(t, "cannot check of existing of %s: %v", to, err)
	} else {
		return
	}

	archive := downloadLibwebpArchive(t, from)
	defer func() { _ = os.Remove(archive) }()
	unpackArchive(t, archive, from, to)

	return
}

func unpackArchive(t *testing.T, from, sourceUrl, to string) {
	fromF, err := os.Open(from)
	if err != nil {
		fatalfWithEnvironment(t, "cannot open libwebp archive %s: %v", from, err)
	}
	defer func() { _ = fromF.Close() }()

	var unpack func(from *os.File, to string) error
	if strings.HasSuffix(sourceUrl, ".tar.gz") {
		unpack = unpackTarGzArchive
	} else if strings.HasSuffix(sourceUrl, ".zip") {
		unpack = unpackZipArchive
	} else {
		t.Fatalf("unsupported libwebp archive type: %s", sourceUrl)
	}

	if err := unpack(fromF, to); err != nil {
		fatalfWithEnvironment(t, "cannot unpack libwebp archive %s: %v", from, err)
	}

	return
}

func unpackTarGzArchive(from *os.File, to string) error {
	gr, err := gzip.NewReader(from)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		p := filepath.Join(to, withoutFirstDirectory(hdr.Name))
		fi := hdr.FileInfo()
		if fi.IsDir() {
			if err := os.MkdirAll(p, fi.Mode()); err != nil {
				return fmt.Errorf("cannot create directory %s: %v", p, err)
			}
		} else {
			if err := writeFile(tr, p, fi.Mode()); err != nil {
				return err
			}
		}
	}
}

func unpackZipArchive(from *os.File, to string) error {
	fi, err := from.Stat()
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(from, fi.Size())
	if err != nil {
		return err
	}
	for _, file := range zr.File {
		p := filepath.Join(to, withoutFirstDirectory(file.Name))
		fi := file.FileInfo()
		if fi.IsDir() {
			if err := os.MkdirAll(p, fi.Mode()); err != nil {
				return fmt.Errorf("cannot create directory %s: %v", p, err)
			}
		} else if err := writeZipFile(file, to); err != nil {
			return err
		}
	}
	return nil
}

func withoutFirstDirectory(in string) (out string) {
	out = in
	index := strings.Index(out, "/")
	if index > 0 {
		out = out[index:]
	}
	return
}

func writeZipFile(file *zip.File, to string) error {
	fr, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = fr.Close() }()
	p := filepath.Join(to, withoutFirstDirectory(file.Name))
	return writeFile(fr, p, file.Mode())
}

func writeFile(r io.Reader, path string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot create parents of %s: %v", path, err)
	}
	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("cannot create output %s: %v", path, err)
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, r); err != nil {
		return fmt.Errorf("cannot write output %s: %v", path, err)
	}
	return nil
}

func downloadLibwebpArchive(t *testing.T, url string) (file string) {
	f, err := ioutil.TempFile("", "libwebp-download")
	if err != nil {
		fatalfWithEnvironment(t, "cannot download libwebp from %s, because cannot get temporary file: %v", url, err)
	}
	defer func() { _ = f.Close() }()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("cannot download libwebp from %s: %v", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		t.Fatalf("cannot download libwebp from %s: %d %s", url, resp.StatusCode, resp.Status)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		t.Fatalf("cannot download libwebp from %s: %v", url, err)
	}

	return f.Name()
}

func libwebpDownloadUrl(t *testing.T) (result string) {
	result = libwebpDownloadBase
	switch runtime.GOOS {
	case "linux":
		result += libwebpDownloadFileLinux
	case "darwin":
		result += libwebpDownloadFileMacOs
	case "windows":
		result += libwebpDownloadFileWindows
	default:
		skipBecauseLibwebpUnavailable(t)
	}

	return strings.ReplaceAll(result, "{version}", libwebpVersion)
}

func toPath(from string) (result string) {
	result = path.Base(from)
	result = strings.ReplaceAll(result, ".tar.gz", "")
	result = strings.ReplaceAll(result, ".zip", "")
	result = filepath.Join(os.TempDir(), result)
	return
}

func skipBecauseLibwebpUnavailable(t *testing.T) {
	t.Skip("There is no libwebp available on this platform." +
		" Execution of this test is not possible." +
		" It will be skipped.")
}
