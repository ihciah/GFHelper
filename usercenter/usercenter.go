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

package usercenter

import (
	"encoding/json"
	"github.com/asdine/storm"
	"github.com/coreos/bbolt"
	"github.com/gfhelper/GFHelper/gameact"
	"github.com/gfhelper/GFHelper/protocol"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

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

func (ud *UserData) Summary() (int64, int64, int64) {
	return ud.MP + ud.AMMO + ud.MRE + ud.PART, ud.BuildCard, ud.FastBuildCard
}

type UserCenter struct {
	DB                *storm.DB
	ServerInfoIOS     *gameact.Server
	ServerInfoAndroid *gameact.Server
	ServerInfoTw      *gameact.Server
	Logger            *log.Logger
}

func MakeUserCenter(database string, needServerInfo bool) (UserCenter, error) {
	db, err := storm.Open(database, storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	if err != nil {
		return UserCenter{}, err
	}
	logger := log.New(os.Stdout, "[UserCenter]", log.LstdFlags)
	if !needServerInfo {
		return UserCenter{db, nil, nil, nil, logger}, nil
	}
	serverIOS, err := gameact.GetServer("ios", "20310")
	if err != nil {
		log.Println("Cannot get iOS server info.")
		return UserCenter{}, err
	}
	log.Println("Got iOS server info.")
	serverAndroid, err := gameact.GetServer("android", "2020") // need update
	if err != nil {
		log.Println("Cannot get android server info.")
		return UserCenter{}, err
	}
	log.Println("Got android server info.")
	serverTw, err := gameact.GetServer("tw", "2013") // need update
	if err != nil {
		log.Println("Cannot get taiwan server info.")
		return UserCenter{}, err
	}
	log.Println("Got taiwan server info.")

	return UserCenter{db, &serverIOS, &serverAndroid, &serverTw, logger}, nil
}

func (uc *UserCenter) ReadAll() ([]UserData, error) {
	var userdata []UserData
	err := uc.DB.All(&userdata)
	return userdata, err
}

func (uc *UserCenter) ReadOne(username string) ([]UserData, error) {
	var userdata UserData
	err := uc.DB.One("ID", username, &userdata)
	return []UserData{userdata}, err
}

func (uc *UserCenter) AddUser(id, password, channelId string) error {
	ucuser := UserData{
		ID:        id,
		Password:  password,
		ChannelId: channelId,
	}
	err := uc.DB.Save(&ucuser)
	return err
}

func (uc *UserCenter) DelUser(id string) error {
	var user UserData
	err := uc.DB.One("ID", id, &user)
	if err != nil {
		return err
	}
	err = uc.DB.DeleteStruct(&user)
	return err
}

func (uc *UserCenter) ExportToFile(filepath string) error {
	users, err := uc.ReadAll()
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserCenter) Do(debugger bool) error {
	ucUsers, err := uc.ReadAll()
	if err != nil {
		uc.Logger.Println("Error when reading database.")
		return err
	}
	uc.Logger.Printf("Read %d users from database.\n", len(ucUsers))
	requestLimit := protocol.MakeRequestLimiter(battleInterval, maintainInterval)
	users := make([]gameact.User, 0, len(ucUsers))
	gunData, err := LoadGunInfo()
	if err != nil {
		uc.Logger.Println("Err when load gun data.")
		return err
	}
	uc.Logger.Printf("Load %d gun info from database.\n", len(gunData))
	for _, u := range ucUsers {
		user, err := uc.ToUser(u, &requestLimit, &gunData, debugger)
		if err == nil {
			users = append(users, user)
		} else {
			uc.Logger.Println("Error when create user", err)
		}
	}
	uc.Logger.Printf("Create %d users.\n", len(users))
	var wg sync.WaitGroup
	for u := range users {
		wg.Add(1)
		go MaxResources(&users[u], &wg)
	}
	uc.Logger.Println("Waiting go routines.")
	wg.Wait()
	return nil
}

func (uc *UserCenter) ToUser(data UserData, requestLimiter *protocol.RequestLimiter, gunData *map[int64]protocol.SingleGun, debugger bool) (gameact.User, error) {
	nbuidCallback := func(user gameact.User) error {
		var err error
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "AccessToken", user.AuthInfo.AccessToken)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "OpenID", user.AuthInfo.Openid)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "UID", user.NbUid.UID)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "Sign", user.NbUid.Sign)
		if err != nil {
			return err
		}
		uc.Logger.Println("NbUid updated.")
		return nil
	}
	updateCallback := func(user gameact.User) error {
		var err error
		userinfo := user.GetResourceInfo()
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "MP", userinfo.MP)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "AMMO", userinfo.AMMO)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "MRE", userinfo.MRE)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "PART", userinfo.PART)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "CORE", userinfo.Core)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "BuildCard", userinfo.BuildCard)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&UserData{ID: data.ID}, "FastBuildCard", userinfo.FastBuildCard)
		if err != nil {
			return err
		}
		uc.Logger.Println("Resources updated.")
		return nil
	}
	serverInfo := uc.ServerInfoIOS
	if data.ChannelId != "ios" {
		serverInfo = uc.ServerInfoAndroid
	}
	u := gameact.MakeUser(data.ID, data.Password, data.ChannelId, serverInfo, data.OpenID, data.AccessToken, data.UID, data.Sign,
		&nbuidCallback, &updateCallback, requestLimiter, gunData, debugger)
	return u, nil
}
