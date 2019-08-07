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

package gameact

import (
	"github.com/gfhelper/GFHelper/protocol"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strconv"
	"strings"
)

type User struct {
	protocol.User
	IndexInfo        gjson.Result
	IndexInfoDirty   bool
	Logger           *log.Logger
	GunData          *map[int64]protocol.SingleGun
	GetNbUidCallBack *func(user User) error
	UpdateCallBack   *func(user User) error
}

type Server struct {
	protocol.Server
}

// 获取可用服务器
func GetServer(channelId string, checkVer string) (Server, error) {
	var server Server
	err := UpdateServer(channelId, checkVer, &server)
	return server, err
}

// 更新服务器信息
func UpdateServer(channelId string, checkVer string, server *Server) error {
	err := protocol.UpdateServer(channelId, checkVer, &server.Server)
	return err
}

// 生成已存在openid和access_token的用户实例
func MakeUser(username, password, channelId string, serverInfo *Server, openid, accessToken, UID, sign string,
	getNbuidCallback, updateCallback *func(User) error, requestLimiter *protocol.RequestLimiter, gunData *map[int64]protocol.SingleGun, debugger bool) User {
	user := protocol.MakeUser(username, password, channelId, &serverInfo.Server, requestLimiter, debugger)
	u := User{User: user, Logger: log.New(os.Stdout, "[Gameact]", log.LstdFlags), IndexInfoDirty: true, GunData: gunData,
		GetNbUidCallBack: getNbuidCallback, UpdateCallBack: updateCallback}
	u.AuthInfo.Openid = openid
	u.AuthInfo.AccessToken = accessToken
	u.AuthInfo.Result = 0
	u.AuthInfo.AppID = protocol.IOSAppId
	if channelId != "ios" {
		u.AuthInfo.AppID = protocol.AndroidAppId
	}
	u.NbUid.UID = UID
	u.NbUid.Sign = sign
	return u
}

// 登陆并获取SkyNbUid，发送资源包已下载，获取Index信息并存储
// 若有未完成任务则取消任务
func (u *User) Login(needLogin, extraInfo bool) error {
	var err error
	isLogin := false
	u.IndexInfoDirty = true

	// 尝试使用已保存的UID和Sign获取结果
	for errCount := 0; errCount < 3 && !isLogin; errCount++ {
		if u.User.NbUid.UID != "" && u.User.NbUid.Sign != "" {
			if index, err := u.Index(); err == nil && index.Get("user_info.id").Exists() {
				u.Logger.Print("UID and Sign login pass.", u.LoginIdentify)
				isLogin = true
			}
		}
	}

	// 尝试验证已有的Openid和AccessToken
	if !isLogin && u.User.AuthInfo.Openid != "" && u.User.AuthInfo.AccessToken != "" {
		if err = u.User.Auth(); err == nil {
			for errCount := 0; errCount < MaxLoginFail; errCount++ {
				err := u.User.GetDigitalSkyNbUid()
				if err == protocol.ErrServerMaintain {
					u.User.RequestLimiter.EnterMaintainMode()
				} else {
					u.User.RequestLimiter.ExitMaintainMode()
				}
				if err == nil {
					isLogin = true
					break
				}
				u.Logger.Print("Get NbUid error", err, u.LoginIdentify)
			}
			if !isLogin {
				return ErrUnexpectedData
			}
		} else {
			u.Logger.Print("Auth not pass, reset data", err, u.LoginIdentify)
			u.AuthInfo.Openid = ""
			u.AuthInfo.AccessToken = ""
			u.NbUid.Sign = ""
			if u.GetNbUidCallBack != nil {
				(*u.GetNbUidCallBack)(*u)
			}
		}
	}

	// 若验证未通过，直接使用用户名密码登陆
	if !isLogin {
		if !needLogin {
			return ErrNotLogin
		}
		err = u.User.Login()
		if err != nil {
			return err
		}
		err = u.GetDigitalSkyNbUid()
		if err != nil {
			return err
		}
	}
	if u.GetNbUidCallBack != nil {
		(*u.GetNbUidCallBack)(*u)
	}
	if !extraInfo {
		return nil
	}
	err = u.DownloadSuccess()
	if err != nil {
		u.Logger.Print("ERR download success", u.LoginIdentify)
		return err
	}
	err = u.Update()
	if err != nil {
		u.Logger.Print("ERR update", u.LoginIdentify)
		return err
	}
	_, err = u.AbortMission()
	return err
}

// 注册账号
func (u *User) Register() (err error) {
	u.IndexInfoDirty = true
	for i := 0; i < MaxRegFail; i += 1 {
		err = u.User.Register()
		if err == protocol.ErrHttpCodeError {
			continue
		} else {
			break
		}
	}
	return err
}

// 若有教程未跳过则跳过新手教程
func (u *User) SkipGuide() error {
	u.Update()
	if !strings.HasSuffix(u.IndexInfo.Get("user_info.guide_info.guide").String(), "1,1,1,1]}") {
		u.Logger.Println("Skip guide.")
		ret := u.Guide()
		u.IndexInfoDirty = true
		return ret
	}
	u.Logger.Println("No need for skipping guide.")
	return nil
}

// 更新用户的所有信息，参数为允许的最大间隔时间，秒为单位
func (u *User) Update() error {
	if !u.IndexInfoDirty {
		return nil
	}
	res, err := u.Index()
	if err != nil {
		return err
	}
	u.IndexInfo = res
	u.IndexInfoDirty = false
	u.Logger.Println("Index info updated.")
	if u.UpdateCallBack != nil {
		(*u.UpdateCallBack)(*u)
	}
	return nil
}

// 如需签到则进行签到
func (u *User) Attendance() error {
	u.Update()
	attendance_type1_time := u.IndexInfo.Get("user_info.user_record.attendance_type1_time").Int()
	tomorrow_zero, err := strconv.ParseInt(u.ServerInfo.Version.TomorrowZero, 10, 64)
	if err != nil {
		return err
	}
	if attendance_type1_time != 0 && attendance_type1_time != tomorrow_zero {
		err = u.User.Attendance()
		u.IndexInfoDirty = true
	}
	u.Logger.Println("Daily check-in.")
	return err
}

func (u *User) GetMails() error {
	homeResult, err := u.User.Home()
	if err != nil {
		return err
	}

	// 获取最大id
	ids := homeResult.Get("index_getmaillist.#.id")
	maxId := int64(0)
	for _, id := range ids.Array() {
		idi := id.Int()
		if idi > maxId {
			maxId = idi
		}
	}
	u.StartID = maxId

	// 只过滤出type in [7, 10, 21]的邮件，分别对应日常任务、主线任务、生涯任务
	types := homeResult.Get("index_getmaillist.#.type").Array()
	ids_ := homeResult.Get("index_getmaillist.#.id").Array()
	set := make(map[int64]int)
	for i, t := range types {
		ti := t.Int()
		if ti == 7 || ti == 10 || ti == 21 {
			// 一键领取
			set[ti] = 1
		}
		if ti == 2 || ti == 5 {
			// 自动领取
			u.GetResourceInMail(ids_[i].Int())
			u.IndexInfoDirty = true
			u.Logger.Println("Resource in mail get.")
		}
	}
	for t := range set {
		_, err = u.User.QuickGetQuestsResourceInMails(t)
		u.IndexInfoDirty = true
		u.Logger.Println("Quest resource in mail get.")
	}
	return err
}
