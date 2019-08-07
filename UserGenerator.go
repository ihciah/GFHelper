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
	"fmt"
	"github.com/gfhelper/GFHelper/cipher"
	"github.com/gfhelper/GFHelper/usercenter"
	"github.com/gfhelper/GFHelper/utils"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var letterNums = []rune("abcdefghjkmnprstuvwxyz01235789")
var nums = []rune("1235789")
var nonOneNums = []rune("235789")
var letters = []rune("abcdefghjkmnprstuvwxyz")

// 返回字母开头的随机串
func randomString(length int64) string {
	if length == 0 {
		return ""
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = letterNums[rand.Intn(len(letterNums))]
		b[i] = nums[rand.Intn(len(nums))]
	}
	b[0] = letters[rand.Intn(len(letters))]
	return string(b)
}

// 非1开头的数字串
func randomNumber(length int64) string {
	if length == 0 {
		return ""
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = nums[rand.Intn(len(nums))]
	}
	b[0] = nonOneNums[rand.Intn(len(nonOneNums))]
	return string(b)
}

func get163Address() (string, string) {
	return randomString(10) + "@163.com", randomString(6)
}

func getQQAddress() (string, string) {
	return randomNumber(11) + "@qq.com", randomString(6)
}

func main() {
	author := []byte("\x57\x66\xe8\xec\x62\xc4\xfd\x40\x3a\xc0\x80\x29\xaf\xe5\x60\xc2\x9a\xb3\x6c\x0a\x41\x70\x0c\x7d\xbd\xf8\x9b\xab\x42\x43\xa3\xc6\xde\xe9\x91\xa9\xdc\x48\xc4\x1f\xd1\xeb\xa9\xc3\xda\xbb\xf8\x59\xb0\x63\x33\x99\x03\x13\x6d\x38\xcf\x1c\xa2\x39\x73\x2c\xc1\xde\x27\x5e\x4e\x19\x21\x50\xe0\xb3")
	authorPlain, err := cipher.XXTEADecryptToString(author)
	if err != nil{
		return
	}
	fmt.Println(authorPlain)
	rand.Seed(time.Now().UnixNano())
	uc, err := usercenter.MakeUserCenter(usercenter.MaxResourceDatabase, false)
	if err != nil {
		log.Println(err)
		return
	}
	existed, err := uc.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}
	iosCount := 0
	androidCount := 0
	for _, u := range existed {
		if u.ChannelId == "ios" {
			iosCount++
		} else {
			androidCount++
		}
	}
	log.Println("iOS user count:", iosCount)
	log.Println("Android user count:", androidCount)

	platform := utils.GetInput("Please input platform(ios or android):")
	if platform != "ios" {
		platform = "android"
	}
	log.Printf("Platform selected: %s\n", platform)
	genCountIn := utils.GetInput("Please input user count:")
	genCount, err := strconv.Atoi(genCountIn)
	if err != nil {
		fmt.Println("Wrong input.")
		return
	}
	log.Printf("User number: %d\n", genCount)
	for i := 0; i < genCount; i++ {
		username, password := getQQAddress()
		uc.AddUser(username, password, platform)
	}
	log.Println("Done!")
}
