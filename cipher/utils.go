// Copyright 2019 ihciah <ihciah@gmail.com>
//
// Licensed under the GNU Affero General Public License, Version 3.0
// (the "License"); you may not use this file except in compliance with the
// License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cipher

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"io/ioutil"
	"strings"
)

// Return a md5 string(lower case) of a input string
func MD5s(s string) string {
	return MD5b([]byte(s))
}

// Return a md5 string(upper case) of a input string
func MD5S(s string) string {
	return strings.ToUpper(MD5b([]byte(s)))
}

// Return a md5 string of input bytes
func MD5b(b []byte) string {
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Return the RC4 encrypted/decrypted bytes of a input bytes and a key string
func RC4(inputs []byte, key string) (output []byte, err error) {
	var c *rc4.Cipher
	if c, err = rc4.NewCipher([]byte(key)); err != nil {
		return nil, err
	}
	output = make([]byte, len(inputs))
	c.XORKeyStream(output, inputs)
	return
}

// Return gzip compressed bytes
func GzipCompress(inputs []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(inputs); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Return gzip decompressed bytes
func GzipDecompress(inputs []byte) ([]byte, error) {
	buf := bytes.NewBuffer(inputs)
	r, err := gzip.NewReader(buf)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	plain, err := ioutil.ReadAll(r)
	return plain, err
}
