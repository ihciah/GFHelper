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
	"flag"
	"fmt"
	"github.com/gfhelper/GFHelper/cipher"
	"github.com/gfhelper/GFHelper/gameact"
	"github.com/gfhelper/GFHelper/protocol"
	"github.com/gfhelper/GFHelper/usercenter"
	"log"
	"sync"
)

func processUser(loginOnly, existedOnly bool, user *gameact.User, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Processing user %s\n", user.LoginIdentify)
	for user.NbUid.UID == "" && !loginOnly {
		user.Requester.RefreshProxy()
		err := user.Register()
		if err == nil {
			log.Println("User reg successful.")
			//time.Sleep(time.Second)
			break
		} else if err == protocol.ErrRegUsed {
			log.Println("User already reged.")
			//time.Sleep(time.Second)
			break
		} else {
			log.Println("User reg error: ", err)
		}
	}

	// Login for users that just reged, i.e: user.NbUid.UID == "" (if not loginOnly)
	if user.NbUid.UID != "" && !loginOnly {
		return
	}
	for errorCounter := 0; errorCounter < 3; errorCounter++ {
		if existedOnly && user.NbUid.UID == "" {
			log.Printf("Bypass user login for %s\n", user.LoginIdentify)
			break
		}
		log.Printf("Login for %s\n", user.LoginIdentify)
		err := user.Login(true, false)
		if err == nil {
			log.Println("User login successfully.")
			break
		}
		log.Println("User login error:", err)
	}
}

func main() {
	author := []byte("\x57\x66\xe8\xec\x62\xc4\xfd\x40\x3a\xc0\x80\x29\xaf\xe5\x60\xc2\x9a\xb3\x6c\x0a\x41\x70\x0c\x7d\xbd\xf8\x9b\xab\x42\x43\xa3\xc6\xde\xe9\x91\xa9\xdc\x48\xc4\x1f\xd1\xeb\xa9\xc3\xda\xbb\xf8\x59\xb0\x63\x33\x99\x03\x13\x6d\x38\xcf\x1c\xa2\x39\x73\x2c\xc1\xde\x27\x5e\x4e\x19\x21\x50\xe0\xb3")
	authorPlain, err := cipher.XXTEADecryptToString(author)
	if err != nil{
		return
	}
	fmt.Println(authorPlain)
	var (
		debugger    bool
		proxy       bool
		loginOnly   bool
		existedOnly bool
		interval    int
	)
	flag.BoolVar(&debugger, "debugger", false, "use debug http proxy or not")
	flag.BoolVar(&proxy, "proxy", false, "use random proxy or not")
	flag.BoolVar(&loginOnly, "login-only", false, "login only")
	flag.BoolVar(&existedOnly, "existed-only", false, "login existed user only")
	flag.IntVar(&interval, "interval", 0, "request interval")
	flag.Parse()

	uc, err := usercenter.MakeUserCenter(usercenter.MaxResourceDatabase, true)
	if err != nil {
		log.Println(err)
		return
	}
	users, err := uc.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	rl := protocol.MakeRequestLimiter(int64(interval), 600)
	for u := range users {
		if users[u].Sign == "" {
			user, err := uc.ToUser(users[u], &rl, &map[int64]protocol.SingleGun{}, debugger)
			if err != nil {
				log.Println("Error when convert user.")
				continue
			}
			user.Requester = protocol.MakeRequester(user.ChannelId, &rl, proxy, debugger)
			wg.Add(1)
			//go processUser(loginOnly, existedOnly, &user, &wg)
			processUser(loginOnly, existedOnly, &user, &wg)
		}
	}
	wg.Wait()
}
