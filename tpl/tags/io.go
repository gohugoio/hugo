package tags

import (
	"errors"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

var ErrNotString = errors.New("tags source input is not a string.")

// AsStrings handles interface to string conversion.
func AsStrings(srcs []interface{}) ([]string, error) {
	var arr []string
	for _, src := range srcs {
		s, ok := src.(string)
		if !ok {
			return nil, ErrNotString
		}

		arr = append(arr, s)
	}
	return arr, nil
}

// ReadBytes reads the bytes from all of the srcs filenames from fs.
func ReadBytes(fs afero.Fs, srcs ...interface{}) ([]byte, error) {
	filenames, err := AsStrings(srcs)
	if err != nil {
		return nil, err
	}

	var bytes []byte
	for _, filename := range filenames {
		r, err := fs.Open(filename)
		if err != nil {
			return nil, err
		}

		defer func(r io.Closer) {
			// not sure how I feel about ignoring Close error.
			r.Close()
		}(r)

		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, b...)
	}

	return bytes, nil
}

// Prepare generates the required directory structure in the publishDir.
func Prepare(fs afero.Fs, publishDir, url string) error {
	destDir := path.Dir(filepath.Join(publishDir, url))
	err := fs.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

// Write outputs the bytes to the url path in the publish directory on fs.
func Write(fs afero.Fs, publishDir, url string, b []byte) error {
	w, err := fs.Create(filepath.Join(publishDir, url))
	if err != nil {
		return err
	}

	// hopefully afero provides same guarantee as Go when `n != len(b)`.
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
