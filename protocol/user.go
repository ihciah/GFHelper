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

package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/gfhelper/GFHelper/cipher"
	"log"
	"os"
	"strings"
)

type User struct {
	AppID         string `json:"app_id"`
	LoginIdentify string `json:"login_identify"`
	LoginPwd      string `json:"login_pwd"`
	ChannelId     string `json:"channel_id"`
	NbUid         NbUid  `json:"-"`
	Requester     `json:"-"`
	AuthInfo      UserLoginRes `json:"-"`
	ServerInfo    *Server      `json:"-"`
	StartID       int64        `json:"-"`
	Logger        *log.Logger  `json:"-"`
}

type UserLoginRes struct {
	AccessToken  string `json:"access_token"`
	Result       int    `json:"result"`
	Msg          string `json:"msg"`
	Openid       string `json:"openid"`
	RegisterMode string `json:"register_mode"`
	State        int    `json:"state"`
	AppID        string `json:"app_id"`
}

type UserAuthRes struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
	State  int    `json:"state"`
}

type UserAuthRequest struct {
	AppID       string `json:"app_id"`
	Openid      string `json:"openid"`
	AccessToken string `json:"access_token"`
	Language    string `json:"language"`
}

type UserRegRes struct {
	AppID  string `json:"app_id"`
	Email  string `json:"email"`
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

type UserRegRequest struct {
	AppID string `json:"app_id"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type UserPostInfo struct {
	AppID   string `json:"app_id"`
	Version string `json:"version"`
}

func MakeUser(username, password, channelId string, serverInfo *Server, RequestLimiter *RequestLimiter, debugger bool) User {
	requester := MakeRequester(channelId, RequestLimiter, false, debugger)
	appId := IOSAppId
	passMD5 := cipher.MD5S(password)
	if channelId != "ios" {
		appId = AndroidAppId
		passMD5 = cipher.MD5s(password)
	}
	return User{AppID: appId, LoginIdentify: username, LoginPwd: passMD5, ChannelId: channelId,
		Requester: requester, ServerInfo: serverInfo, Logger: log.New(os.Stdout, "[Protocol]", log.LstdFlags)}
}

func (u *User) sendUserRequest(dataStruct interface{}, url string, dst interface{}) error {
	var data []byte
	var err error
	switch t := dataStruct.(type) {
	case string:
		data = []byte(t)
	default:
		data, err = json.Marshal(dataStruct)
		if err != nil {
			return err
		}
	}
	data = cipher.XXTEAEncrypt(data)
	jbytes, err := json.Marshal(UserPostInfo{u.AppID, "1.0"})
	if err != nil {
		return err
	}
	body := append(jbytes, byte(0))
	body = append(body, data...)
	respBody, err := u.Post(url, body, true)
	strdst, ok := dst.(*string)
	if ok {
		*strdst = string(respBody)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, dst)
	return err
}

// 发送用户密码登陆请求，返回是否成功，并将结果写入AuthInfo
func (u *User) Login() error {
	url := "http://l.ucenter.ppgame.com/normal_login"
	if u.ChannelId != "ios" {
		url = "http://gf.ucenter.ppgame.com/normal_login"
	}
	err := u.sendUserRequest(u, url, &u.AuthInfo)
	if err != nil {
		return err
	}
	if u.AuthInfo.Result != 0 {
		u.Logger.Println("[ERR] Normal login:", u.AuthInfo.Result, u.AuthInfo.Msg)
		if u.AuthInfo.Result == 60111 {
			return ErrUserNotExist
		}
		return ErrWrongPass
	}
	u.Logger.Printf("User %s normal login pass.", u.LoginIdentify)
	return nil
}

// 发送Register请求，返回是否注册成功
func (u *User) Register() error {
	url := "http://l.ucenter.ppgame.com/email_register"
	if u.ChannelId != "ios" {
		url = "http://gf.ucenter.ppgame.com/email_register"
	}
	//data := UserRegRequest{AppID: u.AppID, Email: u.LoginIdentify, Pwd: u.LoginPwd}
	data := fmt.Sprintf("{\"email\":\"%s\",\"app_id\":\"%s\",\"pwd\":\"%s\"}", u.LoginIdentify, u.AppID, u.LoginPwd)
	var result UserRegRes
	err := u.sendUserRequest(data, url, &result)
	if err != nil {
		return err
	}
	if result.Result != 0 {
		u.Logger.Println("[ERR] Register:", result.Msg, result.Result)
		if result.Result == 60403 {
			return ErrRegUsed
		}
		return ErrRegError
	}
	u.Logger.Printf("User %s with password %s registered.", u.LoginIdentify, u.LoginPwd)
	return nil
}

// 发送Auth请求，返回是否认证成功
func (u *User) Auth() error {
	// Patch for tw version
	if u.ChannelId == "tw" {
		return nil
	}
	url := "http://l.ucenter.ppgame.com/auth"
	if u.ChannelId != "ios" {
		url = "http://gf.ucenter.ppgame.com/auth"
	}
	data := UserAuthRequest{u.AppID, u.AuthInfo.Openid, u.AuthInfo.AccessToken, "cn"}
	var result UserAuthRes
	err := u.sendUserRequest(data, url, &result)
	if err != nil {
		return err
	}
	if result.Result != 0 {
		u.Logger.Println("[ERR] Auth:", result.Result, result.Msg)
		if strings.HasPrefix(result.Msg, "access_token invalid") {
			return ErrWrongToken
		}
		return ErrServerError
	}
	u.Logger.Println("Auth pass.")
	return nil
}
