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
	"encoding/json"
	"fmt"
	"github.com/asdine/storm"
	"github.com/coreos/bbolt"
	"github.com/gfhelper/GFHelper/cipher"
	"log"
	"strings"
	"time"
)

type UserDataOld struct {
	ID            string
	Password      string
	AccessToken   string
	OpenID        string
	UID           string
	Sign          string
	IsReady       bool
	MP            int64
	AMMO          int64
	MRE           int64
	PART          int64
	CORE          int64
	BuildCard     int64
	FastBuildCard int64
}

type UserData struct {
	ID            string
	Password      string
	AccessToken   string
	OpenID        string
	UID           string
	Sign          string
	ChannelId     string
	MP            int64
	AMMO          int64
	MRE           int64
	PART          int64
	CORE          int64
	BuildCard     int64
	FastBuildCard int64
}

func main() {
	author := []byte("\x57\x66\xe8\xec\x62\xc4\xfd\x40\x3a\xc0\x80\x29\xaf\xe5\x60\xc2\x9a\xb3\x6c\x0a\x41\x70\x0c\x7d\xbd\xf8\x9b\xab\x42\x43\xa3\xc6\xde\xe9\x91\xa9\xdc\x48\xc4\x1f\xd1\xeb\xa9\xc3\xda\xbb\xf8\x59\xb0\x63\x33\x99\x03\x13\x6d\x38\xcf\x1c\xa2\x39\x73\x2c\xc1\xde\x27\x5e\x4e\x19\x21\x50\xe0\xb3")
	authorPlain, err := cipher.XXTEADecryptToString(author)
	if err != nil{
		return
	}
	fmt.Println(authorPlain)
	db, err := storm.Open("database.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	if err != nil {
		log.Fatal("Cannot open database")
	}

	oldUsers := make([]UserDataOld, 0, 0)
	db.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UserData"))
		b.ForEach(func(k, v []byte) error {
			if strings.HasPrefix(string(k), "_") {
				return nil
			}
			var user UserDataOld
			if err := json.Unmarshal(v, &user); err == nil {
				oldUsers = append(oldUsers, user)
			}
			return nil
		})
		return nil
	})
	log.Printf("Read %d users from old database.", len(oldUsers))

	dbMig, err := storm.Open("database_migrated.db", storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	if err != nil {
		log.Fatal("Cannot open database")
	}
	for _, user := range oldUsers {
		newUser := UserData{
			user.ID,
			user.Password,
			user.AccessToken,
			user.OpenID,
			user.UID,
			user.Sign,
			"ios",
			user.MP,
			user.AMMO,
			user.MRE,
			user.PART,
			user.CORE,
			user.BuildCard,
			user.FastBuildCard,
		}
		dbMig.Save(&newUser)
	}
	log.Println("Done.")
}
