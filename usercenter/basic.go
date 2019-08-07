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
	"errors"
	"github.com/gfhelper/GFHelper/gameact"
	"github.com/gfhelper/GFHelper/protocol"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func getMissionList(user *gameact.User) []func([]int64) error {
	return []func([]int64) error{
		user.Battle_1_1,
		user.Battle_1_2,
		user.Battle_1_3,
		user.Battle_1_4,
		user.Battle_1_5,
		user.Battle_1_6,
		user.Battle_2_1,
		user.Battle_2_2,
		user.Battle_2_3,
		user.Battle_2_4,
		user.Battle_2_5,
		user.Battle_2_6,
		user.Battle_3_1,
		user.Battle_3_2,
		user.Battle_3_3,
		user.Battle_3_4,
		user.Battle_3_5,
		user.Battle_3_6,
		user.Battle_4_1,
		user.Battle_4_2,
		user.Battle_4_3,
		user.Battle_4_4,
		user.Battle_4_5,
		user.Battle_4_6,
		user.Battle_5_1,
		user.Battle_5_2,
		user.Battle_5_3,
		user.Battle_5_4,
		user.Battle_5_5,
		user.Battle_5_6,
		user.Battle_6_1,
		user.Battle_6_2,
		user.Battle_6_3,
		user.Battle_6_4,
		user.Battle_6_5,
		user.Battle_6_6,
	}
}

func getMissionEmergencyList(user *gameact.User) []func([]int64) error {
	return []func([]int64) error{
		user.Battle_4_1e,
		user.Battle_4_2e,
		user.Battle_4_3e,
		user.Battle_4_4e,
	}
}

func getSingleMissionList(user *gameact.User) []func([]int64) error {
	return []func([]int64) error{
		user.Battle_1_1,
		user.Battle_1_2,
		user.Battle_1_3,
		user.Battle_1_4,
		user.Battle_1_5,
		user.Battle_1_6,
		user.Battle_2_1,
		user.Battle_2_2,
		user.Battle_2_3,
		user.Battle_2_4,
		user.Battle_2_5,
		user.Battle_2_6,
		user.Battle_3_1,
	}
}

// 检查是否满员
func isTeamFull(user *gameact.User) bool {
	teamPopulatiuon := user.GetTeamPopulation()
	sum := int64(0)
	for _, number := range teamPopulatiuon {
		sum += number
	}
	if sum == 5*user.GetResourceInfo().MaxTeam {
		return true
	}
	return false
}

// 通关至某关战役(包括该关)
func FinishMissions(user *gameact.User, missionUntil int64, mustKeeps map[int64]protocol.SingleGun) error {
	user.Update()
	finished := user.GetFinishedMission()
	maxMission := int64(0)
	maxMissionEmergency := int64(0)
	for _, id := range finished {
		if id <= 60 && id >= 5 && id > maxMission && (id+9)%10 >= 4 {
			maxMission = id
		} else if id <= 60 && id >= 5 && id > maxMissionEmergency && (id-1)%10 < 4 {
			maxMissionEmergency = id
		}
	}
	isEmergency := (missionUntil-1)%10 < 4
	if !isEmergency && maxMission >= missionUntil ||
		isEmergency && maxMissionEmergency >= missionUntil {
		log.Printf("Mission Check Pass. MaxMission: %d, MaxMissionEnergency: %d, Require: %d", maxMission, maxMissionEmergency, missionUntil)
		return nil
	}

	user.AbortAllOperations()
	canUse := make([]int64, 0, user.GetResourceInfo().MaxTeam)
	teamLevels := user.GetTeamLevels()
	for teamId := range teamLevels {
		canUse = append(canUse, teamId)
	}
	if len(canUse) < 2 {
		return gameact.ErrUnexpectedData
	}

	var missionList []func([]int64) error
	if maxMission >= 60 {
		missionList = []func([]int64) error{}
	} else {
		mapping := []int64{0, 0, 0, 0, 0, 1, 2, 3, 4, 5}
		startId := mapping[maxMission%10] + (maxMission/10)*6
		endId := mapping[missionUntil%10] + (missionUntil/10)*6
		missionList = getMissionList(user)[startId:endId]
	}

	if isEmergency {
		// 只写了紧急4-1~紧急4-4
		if missionUntil/10 != 4 {
			return gameact.ErrUnexpectedData
		}
		startId := 0
		if maxMissionEmergency/10 == 4 {
			startId = int(maxMissionEmergency % 10)
		}
		endId := int(missionUntil % 10)
		missionList = append(missionList, getMissionEmergencyList(user)[startId:endId]...)
	}

	for _, f := range missionList {
		errCounter := 0
		for {
			if errCounter > MaxErrCount {
				log.Println("Too many failure!")
				return errors.New("too many failure")
			}
			if errCounter > 0 {
				time.Sleep(time.Duration(ErrorSleep) * time.Second)
			}

			err := user.Update()
			if err != nil {
				errCounter++
				continue
			}
			err = user.RetireOtherNormalGuns(mustKeeps, 5, true, false)
			if err != nil {
				if err == gameact.ErrNoFreeSpace {
					return err
				}
				errCounter++
				continue
			}
			err = f(canUse)
			user.IndexInfoDirty = true
			if err != nil {
				errCounter += 1

				logIds := make([]string, 0, len(canUse))
				for _, logId := range canUse {
					logIds = append(logIds, strconv.FormatInt(logId, 10))
				}
				log.Printf("ERR Team %s Start battle failed. %s: %s\n", strings.Join(logIds, ","), err, user.LoginIdentify)

				if err == gameact.ErrResourcesInsufficient {
					return err
				} else {
					abortResult, abortErr := user.AbortMission()
					if abortErr != nil {
						log.Printf("Abort error. %s, %s, %s\n", abortErr, abortResult, user.LoginIdentify)
					} else {
						log.Printf("Abort mission successfully.\n")
					}
				}
			} else {
				break
			}
		}
	}
	return nil
}

// 每个队伍升到N级
func LevelUP(user *gameact.User, levelUntil int64, mustKeeps map[int64]protocol.SingleGun, missionId int) error {
	singleMissionList := getSingleMissionList(user)
	missionList := append(getMissionList(user), getMissionEmergencyList(user)...)
	for {
		flag := true
		teamLevels := user.GetTeamLevels()
		for id, level := range teamLevels {
			// 只针对前4个队伍
			if id > 4 {
				continue
			}
			user.RetireOtherNormalGuns(mustKeeps, 5, true, false)
			if level < levelUntil {
				flag = false
				log.Printf("Team %d level %d. Start battle. %s\n", id, level, user.LoginIdentify)
				var err error
				if missionId < 0 {
					rand.Seed(time.Now().UnixNano())
					fId := rand.Intn(len(singleMissionList))
					err = singleMissionList[fId]([]int64{id})
				} else if missionId < len(missionList) && missionId >= 0 {
					err = missionList[missionId]([]int64{id, id%4 + 1, (id+1)%4 + 1})
				}

				user.IndexInfoDirty = true
				if err != nil {
					log.Printf("ERR Team %d level %d Start battle failed. %s: %s\n", id, level, err, user.LoginIdentify)
					if err == gameact.ErrResourcesInsufficient {
						return err
					} else {
						abortResult, abortErr := user.AbortMission()
						if abortErr != nil {
							log.Printf("Abort error. %s, %s, %s\n", abortErr, abortResult, user.LoginIdentify)
						} else {
							log.Printf("Abort mission successfully.\n")
						}
					}
				}
			}
		}
		if flag {
			log.Println("Level Up Done")
			break
		}
	}
	return nil
}

// 保证每个队伍都满员(需要优化)
func FillTeam(user *gameact.User) error {
	err := user.FillTeam()
	if err != nil {
		return err
	}
	if isTeamFull(user) {
		return nil
	}
	user.AbortAllOperations()
	missionList := getSingleMissionList(user)[:8] //只刷前8关
	errCounter := int64(0)
	// 循环16次，保证满
	for loop := 0; loop < 16; loop++ {
		for _, mission := range missionList {
			user.FillTeam()
			if isTeamFull(user) {
				return nil
			}
			user.RetireOtherNormalGuns(map[int64]protocol.SingleGun{}, 5, true, false)
			errBattle := mission([]int64{1})
			user.IndexInfoDirty = true
			if errBattle != nil {
				user.AbortMission()
				errCounter++
				if errCounter > 10 {
					return errBattle
				}
			} else {
				errCounter = 0
			}
		}
	}
	errFill := user.FillTeam()
	if errFill != nil {
		return errFill
	}
	if isTeamFull(user) {
		return nil
	}
	return gameact.ErrUnexpectedData
}

// 账号初始化: 过教程+全部组队+刷等级+解锁后勤
func BasicPrepare(user *gameact.User) error {
	user.Update()
	err := user.SkipGuide()
	if err != nil {
		return err
	}

	err = FillTeam(user)
	if err != nil {
		return err
	}
	mustKeeps := user.GetGuns()
	for _, missionEndPoint := range []int64{44, 60} {
		for {
			err = FinishMissions(user, missionEndPoint, mustKeeps)
			if err != nil {
				if err == gameact.ErrResourcesInsufficient {
					log.Println("ERR Insufficient resources, doing operation 5~8")
					exit := make(chan bool, 1)
					go DoSimpleOperation(user, exit)
					time.Sleep(2*time.Hour + 10*gameact.TimePadding*time.Second)
					exit <- true
				} else {
					log.Printf("ERR when finish mission! %s %s\n", err, user.LoginIdentify)
					return err
				}
			} else {
				break
			}
		}
	}

	//err = FillTeam(user)
	//if err != nil {
	//	return err
	//}
	for {
		// 每个队伍50级，使用6-6刷
		err := LevelUP(user, 50, mustKeeps, 35)
		if err == nil {
			break
		}
		if err == gameact.ErrResourcesInsufficient {
			log.Println("ERR Insufficient resources, doing operation 5~8")
			exit := make(chan bool, 1)
			go DoSimpleOperation(user, exit)
			time.Sleep(2*time.Hour + 10*gameact.TimePadding*time.Second)
			exit <- true
		} else {
			return err
		}
	}
	return nil
}
