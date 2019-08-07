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

package main

import (
	"bufio"
	"fmt"
	"github.com/gfhelper/GFHelper/cipher"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func urlunquote(str string) string {
	str = strings.Replace(str, "%2f", "/", -1)
	str = strings.Replace(str, "%2b", "+", -1)
	str = strings.Replace(str, "%3d", "=", -1)
	str = strings.Replace(str, "%2d", "-", -1)
	str = strings.Replace(str, "%2F", "/", -1)
	str = strings.Replace(str, "%2B", "+", -1)
	str = strings.Replace(str, "%3D", "=", -1)
	str = strings.Replace(str, "%2D", "-", -1)
	return str
}

func decodeOutdatacode(outdatacode, key string) (string, error) {
	return cipher.AuthCodeDecodeB64(urlunquote(outdatacode), key, false)
}

func decodeReturn(ret, key string) (string, error) {
	return cipher.AuthCodeDecodeB64(ret[1:], key, true)
}

func decode(x, key string) (string, error) {
	if rune(x[0]) == rune('#') {
		return decodeReturn(x, key)
	} else {
		return decodeOutdatacode(x, key)
	}
}

func decode_xxtea(filepath string) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	fmt.Println(cipher.XXTEADecryptToString(content))
}

const defaultSign = "yundoudou"

func main() {
	author := []byte("\x57\x66\xe8\xec\x62\xc4\xfd\x40\x3a\xc0\x80\x29\xaf\xe5\x60\xc2\x9a\xb3\x6c\x0a\x41\x70\x0c\x7d\xbd\xf8\x9b\xab\x42\x43\xa3\xc6\xde\xe9\x91\xa9\xdc\x48\xc4\x1f\xd1\xeb\xa9\xc3\xda\xbb\xf8\x59\xb0\x63\x33\x99\x03\x13\x6d\x38\xcf\x1c\xa2\x39\x73\x2c\xc1\xde\x27\x5e\x4e\x19\x21\x50\xe0\xb3")
	authorPlain, err := cipher.XXTEADecryptToString(author)
	if err != nil{
		return
	}
	fmt.Println(authorPlain)
	sign := defaultSign
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n\r ")
		res, err := decode(text, sign)
		if err == nil {
			fmt.Println(res)
			newSign := gjson.Get(res, "sign").String()
			if len(newSign) == 32 {
				log.Printf("Key set to %s\n", newSign)
				sign = newSign
			}
		} else {
			log.Printf("Decoding errer. Key %s\n", sign)
			newres, newerr := decode(text, defaultSign)
			newSign := gjson.Get(newres, "sign").String()
			if newerr == nil && len(newSign) == 32 {
				sign = newSign
				log.Printf("Key set to %s\n", newSign)
			} else {
				log.Printf("Decoding errer again. Key %s\n", sign)
			}
		}
	}
}
