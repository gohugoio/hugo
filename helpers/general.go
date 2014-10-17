// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
)

func FindAvailablePort() (*net.TCPAddr, error) {
	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer l.Close()
		addr := l.Addr()
		if a, ok := addr.(*net.TCPAddr); ok {
			return a, nil
		}
		return nil, fmt.Errorf("Unable to obtain a valid tcp port. %v", addr)
	}
	return nil, err
}

func GuessType(in string) string {
	switch strings.ToLower(in) {
	case "md", "markdown", "mdown":
		return "markdown"
	case "rst":
		return "rst"
	case "html", "htm":
		return "html"
	}

	return "unknown"
}

func ReaderToBytes(lines io.Reader) []byte {
	b := new(bytes.Buffer)
	b.ReadFrom(lines)
	return b.Bytes()
}

// sliceToLower goes through the source slice and lowers all values.
func SliceToLower(s []string) []string {
	if s == nil {
		return nil
	}

	l := make([]string, len(s))
	for i, v := range s {
		l[i] = strings.ToLower(v)
	}

	return l
}

func Md5String(f string) string {
	h := md5.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum([]byte{}))
}
