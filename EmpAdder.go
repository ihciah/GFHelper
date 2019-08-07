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
	"strconv"
	"strings"
)

func DelUserEmp(uc *usercenter.EmpUserCenter) {
	username := utils.GetInput("Input username:")
	err := uc.DelUser(username)
	if err != nil {
		log.Println("Error!", err)
	} else {
		log.Println("Done!")
	}
}

func AddUserEmp(uc *usercenter.EmpUserCenter) {
	username := utils.GetInput("Input username:")
	password := utils.GetInput("Input password:")
	channelId := utils.GetInput("Input channelId(ios, android, or tw):")
	var level string
	for {
		battle := utils.GetInput("Input level(e.g. 1-6):")
		for _, b := range usercenter.BattleHelperAvailableMission {
			if battle == b {
				level = battle
				break
			}
		}
		if level != "" {
			break
		} else {
			log.Println("Invalid level! Level must be in ", strings.Join(usercenter.BattleHelperAvailableMission, ", "))
		}
	}
	var count int
	var err error
	for {
		countStr := utils.GetInput("Input count(-1=unlimited):")
		count, err = strconv.Atoi(countStr)
		if err == nil {
			break
		} else {
			log.Println("Invalid count! count be a number")
		}
	}

	openid := utils.GetInput("Input openid(Optional):")
	token := utils.GetInput("Input accessToken or sid(Optional):")
	uc.AddUser(username, password, openid, token, channelId, level, count)

	log.Println("Done!")
}

func ListUserEmp(uc *usercenter.EmpUserCenter) {
	users, err := uc.ReadAll()
	if err != nil {
		log.Println("Error!", err)
		return
	}
	for _, u := range users {
		fmt.Printf("ID: %s\tPW: %s\tChannel:%s\n", u.ID, u.Password, u.ChannelId)
		fmt.Printf("Level: %s\n", u.Level)
		fmt.Println()
	}
}

func InteractEmp(uc *usercenter.EmpUserCenter) {
	choices := map[string]func(*usercenter.EmpUserCenter){
		"1": AddUserEmp,
		"2": DelUserEmp,
		"3": ListUserEmp,
	}
	num := utils.GetInput("1: AddUser\n2: DelUser\n3: ListUser\nEnter a choice: ")
	f, ok := choices[num]
	if !ok {
		fmt.Println("Wrong choice!")
		return
	}
	f(uc)
}

func main() {
	author := []byte("\x57\x66\xe8\xec\x62\xc4\xfd\x40\x3a\xc0\x80\x29\xaf\xe5\x60\xc2\x9a\xb3\x6c\x0a\x41\x70\x0c\x7d\xbd\xf8\x9b\xab\x42\x43\xa3\xc6\xde\xe9\x91\xa9\xdc\x48\xc4\x1f\xd1\xeb\xa9\xc3\xda\xbb\xf8\x59\xb0\x63\x33\x99\x03\x13\x6d\x38\xcf\x1c\xa2\x39\x73\x2c\xc1\xde\x27\x5e\x4e\x19\x21\x50\xe0\xb3")
	authorPlain, err := cipher.XXTEADecryptToString(author)
	if err != nil{
		return
	}
	fmt.Println(authorPlain)
	uc, err := usercenter.MakeEmpUserCenter(usercenter.BattleHelperDatabase, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	existed, err := uc.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Existed user:", len(existed))
	for {
		InteractEmp(&uc)
	}
}
