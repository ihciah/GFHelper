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
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

const default_expiry int64 = 3600

// Encrypt the body string to bytes
func AuthCodeEncode(body, key string) ([]byte, error) {
	expiry := default_expiry
	hash := MD5s(key)
	x := MD5s(hash[:16])
	real_key_x := MD5s(hash[16:])
	real_key := real_key_x + MD5s(real_key_x)
	deadline := "0000000000"
	if expiry != 0 {
		deadline = strconv.FormatInt(time.Now().Unix()+expiry, 10)
	}
	mac := MD5s(body + x)[:16]
	unencrypted := deadline + mac + body
	Bytes, err := RC4([]byte(unencrypted), real_key)
	return Bytes, err
}

// Encrypt the body string to base64 string
func AuthCodeEncodeB64(body, key string) (string, error) {
	b64 := base64.StdEncoding
	b, err := AuthCodeEncode(body, key)
	if err != nil {
		return "", err
	}
	return b64.EncodeToString(b), nil
}

// Decrypt the enc bytes
func AuthCodeDecode(enc []byte, key string, withGzip bool) (string, error) {
	hash := MD5s(key)
	x := MD5s(hash[:16])
	realKeyX := MD5s(hash[16:])
	realKey := realKeyX + MD5s(realKeyX)
	decrypted, err := RC4(enc, realKey)
	if err != nil {
		return "", err
	}
	deadline, err := strconv.ParseInt(string(decrypted[:10]), 10, 64)
	if err != nil {
		return "", err
	}
	var retErr error
	if time.Now().Unix() > deadline {
		retErr = ErrDeadline
	}
	mac := string(decrypted[10:26])
	body := decrypted[26:]
	if MD5b(append(body, []byte(x)...))[:16] != mac {
		return "", ErrWrongMAC
	}
	if withGzip {
		body, err = GzipDecompress(body)
		if err != nil {
			return "", err
		}
	}
	return string(body), retErr
}

// Decrypt the enc b64 string
func AuthCodeDecodeB64(encB64, key string, withGzip bool) (string, error) {
	b64 := base64.StdEncoding
	padding := strings.Repeat("=", (len(encB64)/4*4+4-len(encB64))%4)
	enc, err := b64.DecodeString(encB64 + padding)
	if err != nil {
		return "", err
	}
	return AuthCodeDecode(enc, key, withGzip)
}
