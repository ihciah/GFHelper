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
	"strings"
)

type NbUid struct {
	UID             string `json:"uid"`
	Sign            string `json:"sign"`
	IsUsernameExist bool   `json:"is_username_exist"`
	//RealName          int    `json:"real_name"`
	AuthenticationURL string `json:"authentication_url"`
	TcOrderRetry      int    `json:"tc_order_retry"`
}

func (u *User) GetDigitalSkyNbUid() error {
	if u.ChannelId == "tw" {
		// Hard patch for tw version
		u.GetUidTianxiaQueue()
		return u.GetUidTianxiaQueue()
	}
	url := u.ConstructURL("Index", "getDigitalSkyNbUid")
	channelId := u.ChannelId
	if channelId != "ios" {
		channelId = "GWPZ"
	}
	data := map[string]interface{}{
		"openid":       u.AuthInfo.Openid,
		"access_token": u.AuthInfo.AccessToken,
		"app_id":       u.AuthInfo.AppID,
		"channelid":    channelId,
		"idfa":         "",
		"androidid":    "",
		"mac":          "02-00-00-00-00-00",
		"req_id":       u.GetReqID(),
	}
	orderedMap := MakeOrderedMap()
	orderedMap.LoadMap(data)
	formData := orderedMap.Marshal()
	ret, err := u.Post(url, []byte(formData), true)
	if err != nil {
		u.Logger.Println("[ERR]Error when doing post request.")
		return err
	}
	if strings.HasPrefix(string(ret), "{") {
		u.Logger.Println("[ERR]", string(ret))
		return ErrServerMaintain
	}
	jstr, err := cipher.AuthCodeDecodeB64(string(ret)[1:], defaultKey, true)
	if err != nil {
		u.Logger.Println("[ERR]Error when decrypting result.")
		return err
	}
	err = json.Unmarshal([]byte(jstr), &u.NbUid)
	if err != nil {
		u.Logger.Println("[ERR]Error when unmarshaling to struct.")
		return err
	}
	u.Logger.Println("SkyNbUid get.")
	return nil
}

func (u *User) GetUidTianxiaQueue() error {
	url := u.ConstructURL("Index", "getUidTianxiaQueue")
	data := map[string]interface{}{
		"openid": u.AuthInfo.Openid,
		"sid":    u.AuthInfo.AccessToken,
		"req_id": u.GetReqID(),
	}
	orderedMap := MakeOrderedMap()
	orderedMap.LoadMap(data)
	formData := orderedMap.Marshal()
	ret, err := u.Post(url, []byte(formData), true)
	if err != nil {
		u.Logger.Println("[ERR]Error when doing post request.")
		return err
	}
	if strings.HasPrefix(string(ret), "{") {
		u.Logger.Println("[ERR]", string(ret))
		return ErrServerMaintain
	}
	jstr, err := cipher.AuthCodeDecodeB64(string(ret)[1:], defaultKey, true)
	if err != nil {
		u.Logger.Println("[ERR]Error when decrypting result.")
		return err
	}
	err = json.Unmarshal([]byte(jstr), &u.NbUid)
	if err != nil {
		u.Logger.Println("[ERR]Error when unmarshaling to struct.")
		return err
	}
	u.Logger.Println("UidTianxiaQueue get.")
	return nil
}

func (u *User) ConstructURL(t, request string) (url string) {
	host := u.ServerInfo.Addr
	url = fmt.Sprintf("%s%s/%s", host, t, request)
	return
}
