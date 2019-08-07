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
	"fmt"
	"github.com/gfhelper/GFHelper/protocol"
	"log"
	"math/rand"
	"time"
)

// 获取枪娘上限
func (u *User) GetMaxGun() (int64, error) {
	u.Update()
	count := u.IndexInfo.Get("user_info.maxgun").Int()
	// 个数小于4表示结果有问题
	if count < 100 {
		return 100, ErrUnexpectedData
	}
	return count, nil
}

// 使用仓库里的人形补全梯队
func (u *User) FillTeam() (err error) {
	// 规则：同一梯队里相同gunId的人形只能有一个√
	// 正在执行任务的成员或梯队不能编辑√
	u.Update()

	maxTeam := u.GetResourceInfo().MaxTeam

	// 获取所有枪娘id、gun_id、location信息
	ids := u.IndexInfo.Get("gun_with_user_info.#.id").Array()
	gunIds := u.IndexInfo.Get("gun_with_user_info.#.gun_id").Array()
	gunLocations := u.IndexInfo.Get("gun_with_user_info.#.location").Array()
	gunTeamId := u.IndexInfo.Get("gun_with_user_info.#.team_id").Array()
	if len(ids) != len(gunIds) || len(ids) != len(gunLocations) || len(ids) != len(gunTeamId) {
		log.Println(len(ids), len(gunIds), len(gunLocations), len(gunTeamId))
		return ErrUnexpectedData
	}

	// 梯队id: 梯队信息
	// 梯队信息为 枪娘类别id:位置
	teamMap := make(map[int64](map[int64]int64))
	for i := int64(1); i <= maxTeam; i++ {
		teamMap[i] = make(map[int64]int64)
	}
	for i, teamId := range gunTeamId {
		if teamId.Int() != 0 {
			teamMap[teamId.Int()][gunIds[i].Int()] = gunLocations[i].Int()
		}
	}

	// 梯队id: 梯队信息
	// 梯队信息为 位置:1
	teamPosMap := make(map[int64](map[int64]int))
	for i := int64(1); i <= maxTeam; i++ {
		teamPosMap[i] = make(map[int64]int)
	}
	for i, teamId := range gunTeamId {
		if teamId.Int() != 0 {
			teamPosMap[teamId.Int()][gunLocations[i].Int()] = 1
		}
	}

	busyGuns := make(map[int64]int)

	// 对所有空闲枪娘找位子
	for i, teamIdZero := range gunTeamId {
		if teamIdZero.Int() == 0 {
			for teamId, teamInfo := range teamMap {
				// 如果已处理过该枪娘则跳出对该枪娘的处理循环
				if busyGuns[ids[i].Int()] != 0 {
					break
				}
				// 如果不存在该类型的枪娘，且位子没有满
				if teamInfo[gunIds[i].Int()] == 0 && len(teamInfo) != 5 {
					// 找空闲位子
					for loc := int64(1); loc <= 5; loc += 1 {
						if teamPosMap[teamId][loc] == 0 {
							// 把枪娘放入该位置
							err = u.TeamGun(teamId, ids[i].Int(), loc)
							u.IndexInfoDirty = true
							if err != nil {
								return err
							}
							teamPosMap[teamId][loc] = 1
							teamMap[teamId][gunIds[i].Int()] = loc
							busyGuns[ids[i].Int()] = 1
							break
						}
					}
				}
			}
		}
	}
	return
}

type GunInfo struct {
	GunId    int64
	GunLife  int64
	GunLevel int64
}

func (u *User) GetTeam(teamId int64) ([]GunInfo, error) {
	u.Update()
	path := "gun_with_user_info.#[team_id==\"%d\"]#.%s"
	gunIds := u.IndexInfo.Get(fmt.Sprintf(path, teamId, "id")).Array()
	gunLife := u.IndexInfo.Get(fmt.Sprintf(path, teamId, "life")).Array()
	gunLevels := u.IndexInfo.Get(fmt.Sprintf(path, teamId, "gun_level")).Array()
	if len(gunIds) != len(gunLevels) || len(gunIds) != len(gunLife) {
		return []GunInfo{}, ErrUnexpectedData
	}
	gunInfo := make([]GunInfo, 0, len(gunIds))
	for i := range gunIds {
		gunInfo = append(gunInfo, GunInfo{gunIds[i].Int(), gunLife[i].Int(), gunLevels[i].Int()})
	}
	u.Logger.Printf("Get %d guns in team %d.\n", len(gunInfo), teamId)
	return gunInfo, nil
}

func (u *User) EatOrRetireGuns(eatList, retireList []int64) error {
	team1, err1 := u.GetTeam(1)
	team2, err2 := u.GetTeam(2)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	girls := append(team1, team2...)
	rand.Seed(time.Now().Unix())
	luckyGirlId := girls[rand.Intn(len(girls))].GunId
	luckyGirlName := (*u.GunData)[luckyGirlId].GunName

	err := u.User.EatGuns(luckyGirlId, luckyGirlName, eatList)
	if err != nil {
		return err
	}
	err = u.User.RetireGuns(retireList)
	if err != nil {
		return err
	}
	return nil
}

// 已废弃
//func (u *User) RetireExtraGuns() error {
//	u.Update()
//	guns, err := u.GetTeam(0)
//	if err != nil {
//		return err
//	}
//	if len(guns) == 0 {
//		return nil
//	}
//
//	u.IndexInfoDirty = true
//	gunIds := make([]int64, 0, len(guns))
//	for _, gun := range guns {
//		gunIds = append(gunIds, gun.GunId)
//	}
//	u.Logger.Printf("Retire %d guns.\n", len(gunIds))
//	err = u.User.RetireGuns(gunIds)
//	if err != nil {
//		return err
//	}
//	return err
//}

// 回收列表外的低星枪娘，若一定要空出来位置则按照星级从低到高回收
func (u *User) RetireOtherNormalGuns(keep map[int64]protocol.SingleGun, keepStar int64, mustFree, lazy bool) error {
	u.Update()
	guns := u.GetGuns()
	freeGuns, err := u.GetTeam(0)
	if err != nil {
		return err
	}
	retireList := make([]int64, 0, 100)
	eatList := make([]int64, 0, 100)
	for _, g := range freeGuns {
		gunId := g.GunId
		_, ok := keep[gunId]
		//若在保存列表中，或星级满足要求，或星级库中不存在，则保留
		if ok || guns[gunId].GunRank >= keepStar || guns[gunId].GunRank == 0 {
			continue
		}
		if guns[gunId].GunRank >= 3 {
			retireList = append(retireList, gunId)
		} else {
			eatList = append(eatList, gunId)
		}
	}
	// In lazy mode, if the slots are enough and the number of guns to be ate is too small, it will wait.
	if lazy && int64(len(guns)-len(retireList)+RetireGap) < u.GetResourceInfo().MaxGun && len(eatList) < MinEatCount {
		eatList = []int64{}
	}
	if len(retireList)+len(eatList) != 0 {
		u.Logger.Printf("Eat %d guns and retire %d guns.\n", len(eatList), len(retireList))
		err = u.EatOrRetireGuns(eatList, retireList)
		u.IndexInfoDirty = true
		if err != nil {
			return err
		}
	}

	u.Update()
	guns = u.GetGuns()
	freeGuns, err = u.GetTeam(0)
	if err != nil {
		return err
	}

	gunLeft := int64(len(guns))
	gunMax := u.GetResourceInfo().MaxGun
	if gunMax > gunLeft {
		return nil
	}
	if !mustFree {
		return ErrNoFreeSpace
	}
	// 按照上述规则回收后还是超量，则从低到高回收新获得的枪娘，空出一个位置
	mustRetireCount := gunLeft - gunMax + 1
	sortedRetireList := make([]int64, 0, 100)
	star2 := make([]int64, 0, 100)
	star3 := make([]int64, 0, 100)
	star4 := make([]int64, 0, 100)
	star5 := make([]int64, 0, 100)
	for _, g := range freeGuns {
		gunId := g.GunId
		_, ok := keep[gunId]
		if ok {
			continue
		}
		if guns[gunId].GunRank == 2 {
			star2 = append(star2, gunId)
		} else if guns[gunId].GunRank == 3 {
			star3 = append(star3, gunId)
		} else if guns[gunId].GunRank == 4 {
			star4 = append(star4, gunId)
		} else if guns[gunId].GunRank == 5 || guns[gunId].GunRank == 0 {
			star5 = append(star5, gunId)
		}
	}
	sortedRetireList = append(star2, star3...)
	sortedRetireList = append(sortedRetireList, star4...)
	sortedRetireList = append(sortedRetireList, star5...)

	star2 = make([]int64, 0, 100)
	star3 = make([]int64, 0, 100)
	star4 = make([]int64, 0, 100)
	star5 = make([]int64, 0, 100)
	for _, g := range freeGuns {
		gunId := g.GunId
		_, ok := keep[gunId]
		if !ok {
			continue
		}
		if guns[gunId].GunRank == 2 {
			star2 = append(star2, gunId)
		} else if guns[gunId].GunRank == 3 {
			star3 = append(star3, gunId)
		} else if guns[gunId].GunRank == 4 {
			star4 = append(star4, gunId)
		} else if guns[gunId].GunRank == 5 || guns[gunId].GunRank == 0 {
			star5 = append(star5, gunId)
		}
	}
	sortedRetireList = append(sortedRetireList, star2...)
	sortedRetireList = append(sortedRetireList, star3...)
	sortedRetireList = append(sortedRetireList, star4...)
	sortedRetireList = append(sortedRetireList, star5...)

	if int64(len(sortedRetireList)) < mustRetireCount {
		return ErrUnexpectedData
	}
	toRetire := sortedRetireList[:mustRetireCount]
	if len(star2) <= len(toRetire) {
		err = u.EatOrRetireGuns(toRetire[:len(star2)], toRetire[len(star2):])
	} else {
		err = u.EatOrRetireGuns(toRetire, []int64{})
	}
	u.IndexInfoDirty = true
	if err != nil {
		return err
	}
	return nil
}
