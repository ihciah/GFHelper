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
	"github.com/gfhelper/GFHelper/cipher"
	"log"
	"strconv"
)

// 发送jstr到url， 并将返回结果解析到jstruct(指针类型)
// 如果jstr为空，则发送signcode；否则发送outdatacode
// 若jstruct为空则不解析结果；否则将结果json解析到该地址
// 返回结果为api结果和error
// 若结果为加密#开头字串则返回解密后的数据；否则直接返回
func (u *User) SendAction(url, jstr string, jstruct interface{}, limit bool) (string, error) {
	om := MakeOrderedMap()
	om.Add("uid", u.NbUid.UID)
	if jstr == "" {
		// 只发送sign
		encoded, err := cipher.AuthCodeEncodeB64(u.NbUid.Sign, u.NbUid.Sign)
		if err != nil {
			return "", err
		}
		om.Add("signcode", encoded)
	} else {
		// 发送json
		encoded, err := cipher.AuthCodeEncodeB64(jstr, u.NbUid.Sign)
		if err != nil {
			return "", err
		}
		om.Add("outdatacode", encoded)
	}
	om.Add("req_id", u.GetReqID())
	formData := om.Marshal()
	ret, err := u.Post(url, []byte(formData), limit)
	if err != nil || len(ret) == 0 {
		return "", err
	}
	var res string
	// starts with "#"
	if ret[0] == byte(35) {
		res, err = cipher.AuthCodeDecodeB64(string(ret)[1:], u.NbUid.Sign, true)
		if err != nil {
			return "", err
		}
	} else {
		res = string(ret)
	}
	if jstruct == nil {
		return res, nil
	}
	switch t := jstruct.(type) {
	case *int:
		intres, err := strconv.Atoi(res)
		if err != nil {
			return "", err
		}
		*t = intres
		return res, nil
	case *string:
		*t = res
	default:
		err := json.Unmarshal([]byte(res), t)
		if err != nil {
			log.Println(err, res)
		}
		return res, err
	}
	return res, ErrUnknownError
}

func (u *User) SendActionStruct(url string, jsend interface{}, jstruct interface{}, limit bool) (string, error) {
	jstr, err := json.Marshal(jsend)
	if err != nil {
		return "", err
	}
	return u.SendAction(url, string(jstr), jstruct, limit)
}
