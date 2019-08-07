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
	"strconv"
	"strings"
)

type TeamGunReq struct {
	TeamID        int64 `json:"team_id"`
	GunWithUserID int64 `json:"gun_with_user_id"`
	Location      int64 `json:"location"`
}

func (u *User) TeamGun(teamId, gunWithUserId, location int64) (err error) {
	data := TeamGunReq{teamId, gunWithUserId, location}
	result, err := u.SendActionStruct(u.ConstructURL("Gun", "teamGun"), data, nil, false)
	if !strings.HasPrefix(result, "1") {
		u.Logger.Println("[ERR] TeamGun:", result)
		return ErrServerError
	}
	u.Logger.Printf("Gun %d has been put to team %d location %d.\n", gunWithUserId, teamId, location)
	return
}

type EatGunReq struct {
	GunWithUserID int64   `json:"gun_with_user_id"`
	Item9Num      int64   `json:"item9_num"`
	Food          []int64 `json:"food"`
}

func (u *User) EatGun(gunLucky int64, gunLuckyName string, foodIds []int64) (err error) {
	data := EatGunReq{GunWithUserID: gunLucky, Item9Num: 0, Food: foodIds}
	result, err := u.SendActionStruct(u.ConstructURL("Gun", "eatGun"), data, nil, false)
	if !strings.HasPrefix(result, "{\"") {
		u.Logger.Println("[ERR] EatGun:", result)
		return ErrServerError
	}
	u.Logger.Printf("Gun %s(%d) has been upgraded with %d guns.\n", gunLuckyName, gunLucky, len(foodIds))
	return
}

// 回收枪娘
func (u *User) RetireGun(gunIds []int64) error {
	gunIdsStrs := make([]string, 0, len(gunIds))
	for _, id := range gunIds {
		gunIdsStrs = append(gunIdsStrs, strconv.FormatInt(id, 10))
	}
	requestStr := "[" + strings.Join(gunIdsStrs, ",") + "]"
	result, err := u.SendAction(u.ConstructURL("Gun", "retireGun"), requestStr, nil, false)
	if !strings.HasPrefix(result, "1") {
		u.Logger.Println("[ERR] RetireGun: ", result)
		return ErrServerError
	}
	u.Logger.Printf("%d guns has been retired.\n", len(gunIds))
	return err
}

// 回收枪娘(支持40以上)
func (u *User) RetireGuns(gunIds []int64) error {
	if len(gunIds) == 0 {
		return nil
	}
	for i := 0; i < len(gunIds); i += maxRetire {
		end := i + maxRetire
		if end > len(gunIds) {
			end = len(gunIds)
		}
		err := u.RetireGun(gunIds[i:end])
		if err != nil {
			return err
		}
	}
	return nil
}

// 升级枪娘(支持40以上)
func (u *User) EatGuns(gunLucky int64, gunLuckyName string, gunIds []int64) error {
	if len(gunIds) == 0 {
		return nil
	}
	for i := 0; i < len(gunIds); i += maxRetire {
		end := i + maxRetire
		if end > len(gunIds) {
			end = len(gunIds)
		}
		err := u.EatGun(gunLucky, gunLuckyName, gunIds[i:end])
		if err != nil {
			return err
		}
	}
	return nil
}
