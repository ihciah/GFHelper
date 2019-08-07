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
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

// 获取用户信息
type IndexReq struct {
	Time int64 `json:"time"`
}

func (u *User) Index() (gjson.Result, error) {
	data := IndexReq{Time: time.Now().Unix()}
	res, err := u.SendActionStruct(u.ConstructURL("Index", "index"), data, nil, false)
	result := gjson.Parse(res)
	if err != nil {
		return result, err
	}
	if result.Get("user_info").Exists() {
		return result, nil
	}
	return result, ErrServerMaintain
}

// 资源下载/更新完毕?
func (u *User) DownloadSuccess() (err error) {
	_, err = u.SendAction(u.ConstructURL("Index", "downloadSuccess"), "", nil, false)
	return
}

// home: 游戏主页面，获取一些更新
type HomeReq struct {
	DataVersion string `json:"data_version"`
	AbVersion   string `json:"ab_version"`
	StartID     int64  `json:"start_id"`
	IgnoreTime  int    `json:"ignore_time"`
}

func (u *User) Home() (gjson.Result, error) {
	data := HomeReq{DataVersion: u.ServerInfo.Version.DataVersion,
		AbVersion:  u.ServerInfo.Version.AbVersion,
		StartID:    u.StartID,
		IgnoreTime: 1}
	res, err := u.SendActionStruct(u.ConstructURL("Index", "home"), data, nil, false)
	return gjson.Parse(res), err
}

// 获取邮件里的奖励
func (u *User) QuickGetResourceInMails() (gjson.Result, error) {
	res, err := u.SendAction(u.ConstructURL("Index", "quickGetResourceInMails"), "", nil, false)
	return gjson.Parse(res), err
}

// 获取所有指定类型的任务奖励
type QuickGetQuestsResourceInMailsReq struct {
	Type int64 `json:"type"`
}

func (u *User) QuickGetQuestsResourceInMails(t int64) (gjson.Result, error) {
	data := QuickGetQuestsResourceInMailsReq{t}
	res, err := u.SendActionStruct(u.ConstructURL("Index", "QuickGetQuestsResourceInMails"), data, nil, false)
	if strings.HasPrefix(res, "error") {
		return gjson.Parse(res), ErrServerError
	}
	return gjson.Parse(res), err
}

// 获取单个邮件资源
type GetResourceInMailReq struct {
	MailWithUserID int64 `json:"mail_with_user_id"`
}

func (u *User) GetResourceInMail(mailWithUserId int64) error {
	data := GetResourceInMailReq{mailWithUserId}
	_, err := u.SendActionStruct(u.ConstructURL("Index", "getResourceInMail"), data, nil, false)
	return err
}

// 跳过所有新手教学
type GuideReq struct {
	Guide string `json:"guide"`
}

func (u *User) Guide() (err error) {
	data := GuideReq{"{\"course\":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]}"}
	_, err = u.SendActionStruct(u.ConstructURL("Index", "guide"), data, nil, false)
	return
}

// 设置用户名
type SetUserNameReq struct {
	Name string `json:"name"`
}

func (u *User) SetUserName(username string) (err error) {
	data := SetUserNameReq{username}
	_, err = u.SendActionStruct(u.ConstructURL("Index", "setUserName"), data, nil, false)
	return
}

// 每日签到
func (u *User) Attendance() error {
	result, err := u.SendAction(u.ConstructURL("Index", "attendance"), "", nil, false)
	if !strings.HasPrefix(result, "{") {
		return ErrServerError
	}
	return err
}
