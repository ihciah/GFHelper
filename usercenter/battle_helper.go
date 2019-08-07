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
	"github.com/gfhelper/GFHelper/gameact"
	"github.com/gfhelper/GFHelper/protocol"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 返回GunID和{星级，名字}
func LoadGunInfo() (map[int64]protocol.SingleGun, error) {
	result := make(map[int64]protocol.SingleGun)
	gresult := gjson.Parse(gunInfoJson)
	ids := gresult.Get("guns.#.id").Array()
	ranks := gresult.Get("guns.#.rank").Array()
	names := gresult.Get("guns.#.cn_name").Array()
	if len(ids) != len(ranks) {
		return result, gameact.ErrUnexpectedData
	}
	for i := range ids {
		result[ids[i].Int()] = protocol.SingleGun{ranks[i].Int(), names[i].String()}
	}
	return result, nil
}

type EmpUserData struct {
	ID          string
	Password    string
	ChannelId   string
	AccessToken string
	OpenID      string
	UID         string
	Sign        string
	Level       string
	Count       int
}

type EmpUserCenter struct {
	UserCenter
}

func MakeEmpUserCenter(database string, needServerInfo bool) (EmpUserCenter, error) {
	userCenter, err := MakeUserCenter(database, needServerInfo)
	if err != nil {
		return EmpUserCenter{}, err
	}
	return EmpUserCenter{userCenter}, nil
}

func (uc *EmpUserCenter) ReadAll() ([]EmpUserData, error) {
	var userdata []EmpUserData
	err := uc.DB.All(&userdata)
	return userdata, err
}

func (uc *EmpUserCenter) AddUser(id, password, openid, token, channelId, level string, count int) error {
	ucuser := EmpUserData{
		ID:          id,
		Password:    password,
		OpenID:      openid,
		AccessToken: token,
		ChannelId:   channelId,
		Level:       level,
		Count:       count,
	}
	err := uc.DB.Save(&ucuser)
	return err
}

func (uc *EmpUserCenter) DelUser(id string) error {
	var user EmpUserData
	err := uc.DB.One("ID", id, &user)
	if err != nil {
		return err
	}
	err = uc.DB.DeleteStruct(&user)
	return err
}

func (uc *EmpUserCenter) MaxMreAmmo(debugger bool) error {
	ucusers, err := uc.ReadAll()
	if err != nil {
		uc.Logger.Println("Error when reading database.")
		return err
	}
	uc.Logger.Printf("Read %d users from database.\n", len(ucusers))
	requestLimit := protocol.MakeRequestLimiter(battleInterval, maintainInterval)
	gunData, err := LoadGunInfo()
	if err != nil {
		uc.Logger.Println("Err when load gun data.")
		return err
	}
	uc.Logger.Printf("Load %d gun info from database.\n", len(gunData))
	var wg sync.WaitGroup
	if err != nil {
		return err
	}
	for _, u := range ucusers {
		user, err := uc.ToUser(u, &requestLimit, &gunData, debugger)
		if err != nil {
			uc.Logger.Println("Error when create user", err)
			continue
		}
		wg.Add(1)
		uc.Logger.Printf("Create user %s.\n", user.LoginIdentify)
		uc.Logger.Printf("Doing MaxMreAmmo %s.\n", user.LoginIdentify)
		MaxMreAmmo(&user, &wg)
	}
	uc.Logger.Println("Waiting go routines.")
	wg.Wait()
	return nil
}

func (uc *EmpUserCenter) Do(debugger bool) error {
	ucusers, err := uc.ReadAll()
	if err != nil {
		uc.Logger.Println("Error when reading database.")
		return err
	}
	uc.Logger.Printf("Read %d users from database.\n", len(ucusers))
	requestLimit := protocol.MakeRequestLimiter(battleInterval, maintainInterval)
	gunData, err := LoadGunInfo()
	if err != nil {
		uc.Logger.Println("Err when load gun data.")
		return err
	}
	uc.Logger.Printf("Load %d gun info from database.\n", len(gunData))
	var wg sync.WaitGroup
	exitSignals := make([]chan bool, 0, len(ucusers))
	for _, u := range ucusers {
		user, err := uc.ToUser(u, &requestLimit, &gunData, debugger)
		if err != nil {
			uc.Logger.Println("Error when create user", err)
			continue
		}
		wg.Add(1)
		uc.Logger.Printf("Create user %s.\n", user.LoginIdentify)
		exit := make(chan bool, 1)
		exitSignals = append(exitSignals, exit)
		go BattleLoop(&user, u.Level, u.Count, &wg, exit)
	}
	uc.Logger.Println("Waiting go routines or sigterm.")
	dontWait := make(chan bool, 2)
	go func(exitSignals []chan bool, dontWait <-chan bool) {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		select {
		case <-c:
			uc.Logger.Printf("\n=======\nExit signal got!\nPlease wait for missions done.\n=======\n")
			for _, s := range exitSignals {
				s <- true
			}
		case <-dontWait:
			return
		}
	}(exitSignals, dontWait)
	wg.Wait()
	dontWait <- true
	return nil
}

func (uc *EmpUserCenter) ToUser(data EmpUserData, requestLimiter *protocol.RequestLimiter, gunData *map[int64]protocol.SingleGun, debugger bool) (gameact.User, error) {
	nbuidCallback := func(user gameact.User) error {
		var err error
		err = uc.DB.UpdateField(&EmpUserData{ID: data.ID}, "AccessToken", user.AuthInfo.AccessToken)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&EmpUserData{ID: data.ID}, "OpenID", user.AuthInfo.Openid)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&EmpUserData{ID: data.ID}, "UID", user.NbUid.UID)
		if err != nil {
			return err
		}
		err = uc.DB.UpdateField(&EmpUserData{ID: data.ID}, "Sign", user.NbUid.Sign)
		if err != nil {
			return err
		}
		uc.Logger.Println("NbUid updated.")
		return nil
	}
	serverInfo := uc.ServerInfoIOS
	if data.ChannelId == "android" {
		serverInfo = uc.ServerInfoAndroid
	} else if data.ChannelId == "tw" {
		serverInfo = uc.ServerInfoTw
	}
	u := gameact.MakeUser(data.ID, data.Password, data.ChannelId, serverInfo, data.OpenID, data.AccessToken, data.UID, data.Sign,
		&nbuidCallback, nil, requestLimiter, gunData, debugger)
	return u, nil
}

func BattleLoop(user *gameact.User, level string, count int, wg *sync.WaitGroup, exit <-chan bool) {
	defer wg.Done()

	for {
		err := user.Login(true, true)
		if err != nil {
			// 如果未登陆则直接返回错误
			if err == gameact.ErrNotLogin {
				return
			}
			log.Println("Error when login. Will retry after 1 minute.", err)
			awake := time.After(1 * time.Minute)
			select {
			case <-awake:
				continue
			case <-exit:
				return
			}
		} else {
			break
		}
	}
	// 0. 登陆并读取用户信息
	// 1. 停止所有后勤，整理多余人形
	// 2. 刷新
	// 3. 检测资源
	// 4. 完成任务

	battle2func := map[string]func([]int64) error{
		"1-6": user.Battle_1_6,
		"2-6": user.Battle_2_6,
		"3-6": user.Battle_3_6,
		"4-6": user.Battle_4_6,
		"5-6": user.Battle_5_6,
		"6-6": user.Battle_6_6sp,
	}
	battle2mission := map[string]int64{
		"1-6": 10,
		"2-6": 20,
		"3-6": 30,
		"4-6": 40,
		"5-6": 50,
		"6-6": 60,
	}
	levelFunc, ok := battle2func[level]
	if !ok {
		log.Println("Level not found.", level)
		return
	}

	mustKeeps := user.GetGuns()
	user.AbortAllOperations()
	user.AbortMission()

	for {
		err := FinishMissions(user, battle2mission[level], mustKeeps)
		user.IndexInfoDirty = true
		if err == gameact.ErrResourcesInsufficient {
			log.Println("Insufficient resources.", user.LoginIdentify)
			return
		}
		if err != nil {
			log.Println("ERR", err)
			user.AbortMission()
		} else {
			log.Println("Mission check pass.")
			break
		}
		select {
		case <-exit:
			return
		default:
			continue
		}
	}

	for missionCounter := 1; count <= 0 || missionCounter < count; missionCounter++ {
		// Set lazy to true, means it only retire guns when there's no free slots.
		// Not a good idea for new account.
		// However, when upgrading, it will cause
		user.RetireOtherNormalGuns(mustKeeps, 5, true, true)
		user.Update()
		err := levelFunc([]int64{1, 2, 3, 4})
		user.IndexInfoDirty = true
		if err == gameact.ErrResourcesInsufficient {
			log.Println("Insufficient resources.", user.LoginIdentify)
			log.Printf("Finish mission for %d times. user: %s", missionCounter, user.LoginIdentify)
			return
		}
		if err != nil {
			log.Println("ERR", err)
			user.AbortMission()
		} else {
			log.Printf("Mission pass. User: %s, Mission count: %d", user.LoginIdentify, missionCounter)
		}
		select {
		case <-exit:
			log.Printf("Mission will stop. Retiring extra guns.")
			user.RetireOtherNormalGuns(mustKeeps, 5, true, false)
			log.Printf("Mission stop. User: %s, Mission count: %d", user.LoginIdentify, missionCounter)
			return
		default:
			continue
		}
	}
}
