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
	"strconv"
)

// 获取当前正在执行任务的梯队列表
func (u *User) GetOperationTeams() []int64 {
	u.Update()
	teamIds := make([]int64, 0)
	for _, teamId := range u.IndexInfo.Get("operation_act_info.#.team_id").Array() {
		ti := teamId.Int()
		if ti != 0 {
			teamIds = append(teamIds, ti)
		}
	}
	u.Logger.Printf("%d teams are doing operation\n", len(teamIds))
	return teamIds
}

// 获取可用梯队（不一定是有人的梯队
func (u *User) GetAvailableTeams() []int64 {
	u.Update()
	unavailableTeams := make(map[int64]int)
	for _, team := range u.GetOperationTeams() {
		unavailableTeams[team] = 1
	}
	maxTeam := u.GetResourceInfo().MaxTeam
	availableTeams := make([]int64, 0, maxTeam)
	var i int64
	for i = 1; i <= maxTeam; i++ {
		if unavailableTeams[i] == 0 {
			availableTeams = append(availableTeams, i)
		}
	}
	u.Logger.Printf("%d teams are available\n", len(availableTeams))
	return availableTeams
}

// 获取已经完成的任务列表
func (u *User) GetFinishedMission() []int64 {
	u.Update()
	missionIds := make([]int64, 0)
	for _, missionId := range u.IndexInfo.Get("mission_with_user_info.#[win_counter!=0]#.mission_id").Array() {
		ti := missionId.Int()
		if ti != 0 {
			missionIds = append(missionIds, ti)
		}
	}
	u.Logger.Printf("%d operations already finished.\n", len(missionIds))
	return missionIds
}

// 获取正在执行的后勤id
func (u *User) GetOperatingIDs() []int64 {
	u.Update()
	operatingIds := make([]int64, 0)
	for _, operatingId := range u.IndexInfo.Get("operation_act_info.#.operation_id").Array() {
		ti := operatingId.Int()
		if ti != 0 {
			operatingIds = append(operatingIds, ti)
		}
	}
	u.Logger.Printf("%d teams are doing operation\n", len(operatingIds))
	return operatingIds
}

//获取梯队(队长)等级
func (u *User) GetTeamLevels() map[int64]int64 {
	u.Update()
	leaderTeams := u.IndexInfo.Get("gun_with_user_info.#[location=1]#.team_id").Array()
	leaderLevels := u.IndexInfo.Get("gun_with_user_info.#[location=1]#.gun_level").Array()
	result := make(map[int64]int64)
	if len(leaderLevels) != len(leaderTeams) {
		return map[int64]int64{}
	}
	for i := range leaderTeams {
		result[leaderTeams[i].Int()] = leaderLevels[i].Int()
	}
	return result
}

//获取梯队人数
func (u *User) GetTeamPopulation() map[int64]int64 {
	u.Update()
	result := make(map[int64]int64)
	maxTeam := u.GetResourceInfo().MaxTeam
	for i := int64(1); i <= maxTeam; i++ {
		result[i] = int64(len(u.IndexInfo.Get("gun_with_user_info.#[team_id=" + strconv.FormatInt(i, 10) + "]#.id").Array()))
	}
	return result
}

// 获取所有人形和其对应{星级，名字}
func (u *User) GetGuns() map[int64]protocol.SingleGun {
	u.Update()
	result := make(map[int64]protocol.SingleGun)
	guns := u.IndexInfo.Get("gun_with_user_info.#.id").Array()
	gunsTypes := u.IndexInfo.Get("gun_with_user_info.#.gun_id").Array()
	for i := range guns {
		var gun_type int64 = -1
		if i < len(gunsTypes) {
			gun_type = gunsTypes[i].Int()
		}
		singleData, ok := (*u.GunData)[gun_type]
		if !ok || gun_type < 0 {
			// 防止误删，库中未找到则设置为5星
			u.Logger.Printf("Gun %d not found in database.", guns[i].Int())
			singleData = protocol.SingleGun{GunRank: 5}
		} else {
			result[guns[i].Int()] = singleData
		}
	}
	return result
}

type UserInfo struct {
	Core          int64 // 核心
	Gem           int64 // 钻石
	Lv            int64 // 等级
	MaxGun        int64 // 仓库容量
	MaxTeam       int64 // 最大梯队数
	MP            int64 // 人力
	AMMO          int64 // 弹药
	MRE           int64 // 口粮
	PART          int64 // 零件
	BuildCard     int64 // 快速建造
	FastBuildCard int64 // 制造契约
}

// 刷新用户资源信息
func (u *User) GetResourceInfo() UserInfo {
	userInfo := UserInfo{}
	userInfo.Core = u.IndexInfo.Get("user_info.core").Int()
	userInfo.Gem = u.IndexInfo.Get("user_info.gem").Int()
	userInfo.Lv = u.IndexInfo.Get("user_info.lv").Int()
	userInfo.MaxGun = u.IndexInfo.Get("user_info.maxgun").Int()
	userInfo.MaxTeam = u.IndexInfo.Get("user_info.maxteam").Int()
	userInfo.MP = u.IndexInfo.Get("user_info.mp").Int()
	userInfo.AMMO = u.IndexInfo.Get("user_info.ammo").Int()
	userInfo.MRE = u.IndexInfo.Get("user_info.mre").Int()
	userInfo.PART = u.IndexInfo.Get("user_info.part").Int()
	userInfo.BuildCard = u.IndexInfo.Get("item_with_user_info.#[item_id==1].number").Int()
	userInfo.FastBuildCard = u.IndexInfo.Get("item_with_user_info.#[item_id==3].number").Int()
	return userInfo
}
