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
	"github.com/gfhelper/GFHelper/usercenter"
	"github.com/gfhelper/GFHelper/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func DeleteUsers(uc *usercenter.UserCenter) {
	// 给定用户列表文件，从数据库中删除这些用户
	filepath := utils.GetInput("Input file path:")
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Unable to load given file.")
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sep := " "
		if strings.Contains(line, ",") {
			sep = ","
		}
		if strings.Contains(line, "\t") {
			sep = "\t"
		}
		parts := strings.Split(line, sep)
		if len(parts) > 0 {
			account := strings.Trim(parts[0], " ,\t\r\n")
			if strings.HasSuffix(account, "com") {
				err := uc.DelUser(account)
				if err != nil {
					log.Printf("Error when delete user %s: %s.", account, err)
				} else {
					log.Printf("User %s deleted.", account)
				}
			}
		}
	}
}

func AddUser(uc *usercenter.UserCenter) {
	userID := utils.GetInput("ID:")
	password := utils.GetInput("Password:")
	channel := utils.GetInput("Channel(ios, android, tw):")
	uc.AddUser(userID, password, channel)
	log.Println("Done!")
}

func ExportAllUsers(uc *usercenter.UserCenter) {
	// 导出所有用户到json文件
	err := uc.ExportToFile("export.json")
	if err != nil {
		log.Printf("Error when export: %s", err)
	}
	log.Println("Done!")
}

func writeFileChannel(users []usercenter.UserData, suffix string) error {
	outputFullPath := fmt.Sprintf("export_%s_full.txt", suffix)
	outputPath := fmt.Sprintf("export_%s.txt", suffix)
	linesFull := make([]string, 0, len(users))
	lines := make([]string, 0, len(users))
	for _, u := range users {
		lines = append(lines, fmt.Sprintf("%s,%s", u.ID, u.Password))
		linesFull = append(linesFull, fmt.Sprintf("%s %s %d %d %d %d %d %d %d", u.ID, u.Password, u.MP,
			u.AMMO, u.MRE, u.PART, u.CORE, u.BuildCard, u.FastBuildCard))
	}
	err := ioutil.WriteFile(outputPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outputFullPath, []byte(strings.Join(linesFull, "\n")), 0644)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(users []usercenter.UserData, suffix string) error {
	users_ios := make([]usercenter.UserData, 0, len(users))
	users_android := make([]usercenter.UserData, 0, len(users))
	for _, u := range users {
		if u.ChannelId == "ios" {
			users_ios = append(users_ios, u)
		} else {
			users_android = append(users_android, u)
		}
	}
	var err error
	err = writeFileChannel(users_ios, suffix+"_ios")
	if err != nil {
		return err
	}
	err = writeFileChannel(users_android, suffix+"_android")
	return err
}

func ExportFilteredUsers(uc *usercenter.UserCenter) {
	// 导出过滤后的用户
	exportLevels := [][]int{{600, 1150000}, {500, 1100000}, {400, 1050000}, {300, 800000}}

	users, err := uc.ReadAll()
	if err != nil {
		log.Printf("Error when read users: %s", err)
	}
	leveledAccounts := make(map[int][]usercenter.UserData)
	for _, level := range exportLevels {
		leveledAccounts[level[0]] = make([]usercenter.UserData, 0, len(users))
	}

	for _, u := range users {
		resources, build, fastBuild := u.Summary()
		for _, level := range exportLevels {
			if build >= int64(level[0]) && fastBuild >= int64(level[0]) && resources >= int64(level[1]) {
				leveledAccounts[level[0]] = append(leveledAccounts[level[0]], u)
				break
			}
		}
	}

	for level, usersExport := range leveledAccounts {
		writeFile(usersExport, strconv.Itoa(level))
	}
	log.Println("Done!")
}

func Interact(uc *usercenter.UserCenter) {
	choices := map[string]func(*usercenter.UserCenter){
		"1": ExportAllUsers,
		"2": ExportFilteredUsers,
		"3": DeleteUsers,
		"4": AddUser,
	}
	num := utils.GetInput("1: Export all users\n2: Export filtered users\n3: Delete given users\n4: Add user\nEnter a choice: ")
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
	uc, err := usercenter.MakeUserCenter(usercenter.MaxResourceDatabase, false)
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
		Interact(&uc)
	}
}
