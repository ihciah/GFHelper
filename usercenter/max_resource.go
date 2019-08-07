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
	"log"
	"sync"
	"time"
)

func InfOperationLoop(user *gameact.User, teamId, operationId, duration int64, exit <-chan bool) {
	for {
		err := InfOperation(user, teamId, operationId, duration, exit)
		if err != nil {
			log.Println("Relogin", user.LoginIdentify)
			user.Login(false, true)
			time.Sleep(1 * time.Minute)
		}
	}
}

// 直接对单个operation用select效率较低，应采用更好的事件驱动模式
func InfOperation(user *gameact.User, teamId, operationId, duration int64, exit <-chan bool) error {
	operatingIds := user.IndexInfo.Get("operation_act_info.#.operation_id").Array()
	operatingTeams := user.IndexInfo.Get("operation_act_info.#.team_id").Array()
	operatingStartTimes := user.IndexInfo.Get("operation_act_info.#.start_time").Array()

	if len(operatingTeams) != len(operatingIds) || len(operatingStartTimes) != len(operatingIds) {
		return gameact.ErrUnexpectedData
	}

	// 清空所有任务
	for i := range operatingIds {
		// 梯队正在执行任务且不匹配OR任务正在被执行且梯队不匹配
		if operatingTeams[i].Int() == teamId && operatingIds[i].Int() != operationId ||
			(operatingTeams[i].Int() != teamId && operatingIds[i].Int() == operationId) {
			for retry := 0; true; retry++ {
				err := user.AbortOperation(operationId)
				if err == nil {
					break
				}
				if retry > 3 {
					opid, cerr := user.CheckTeam(teamId)
					if cerr == nil && opid == 0 {
						log.Println("CheckTeam abort pass", user.LoginIdentify)
						break
					}
					time.Sleep(1 * time.Minute)
				}
				if retry > 30 {
					return err
				}
			}
		}
		// 对于符合的任务，等待然后完成任务
		if operatingTeams[i].Int() == teamId && operatingIds[i].Int() == operationId {
			timeToWait := operatingStartTimes[i].Int() + duration + gameact.TimePadding - time.Now().Unix()
			log.Println("Already started, wait", timeToWait, user.LoginIdentify)
			if timeToWait > 0 {
				time.Sleep(time.Duration(timeToWait) * time.Second)
			}
			for retry := 0; true; retry++ {
				err := user.FinishOperation(operationId)
				if err == nil {
					break
				}
				log.Println("ERR finish failed", user.LoginIdentify)
				opid, cerr := user.CheckTeam(teamId)
				if cerr == nil && opid == 0 {
					log.Println("CheckTeam finish pass", user.LoginIdentify)
					break
				}
				log.Println("ERR CheckTeam finish failed", user.LoginIdentify)
				if retry > 3 {
					log.Println("ERR finish sleep 1 min, retry", retry, user.LoginIdentify)
					time.Sleep(1 * time.Minute)
				}
				if retry > 30 {
					return err
				}
			}
		}
	}

	// 获取MaxLevel
	levels := user.GetTeamLevels()
	maxLevel := int64(0)
	for _, level := range levels {
		if level > maxLevel {
			maxLevel = level
		}
	}
	// 开始任务、sleep、结束任务
	for {
		// 开始任务
		for retry := 0; true; retry++ {
			err := user.StartOperation(teamId, operationId, maxLevel)
			if err == nil {
				break
			}
			log.Println("ERR start failed", user.LoginIdentify)
			opid, cerr := user.CheckTeam(teamId)
			// 事实上已经开始了并且是该任务
			if cerr == nil && opid == operationId {
				log.Println("CheckTeam start pass", user.LoginIdentify)
				break
			} else {
				log.Println("CheckTeam start failed", user.LoginIdentify)
			}
			if retry > 3 {
				log.Println("start sleep 1 min, retry", retry, user.LoginIdentify)
				time.Sleep(1 * time.Minute)
			}
			if retry > 30 {
				return err
			}
		}

		// 等待
		log.Printf("Wait for %d seconds.\n", duration+gameact.TimePadding)
		select {
		case <-time.After(time.Duration(duration+gameact.TimePadding) * time.Second):
			// 结束任务
			for retry := 0; true; retry++ {
				err := user.FinishOperation(operationId)
				if err == nil {
					break
				}
				if retry > 3 {
					time.Sleep(1 * time.Minute)
				}
				if retry > 30 {
					return err
				}
			}
		case <-exit:
			return nil
		}
	}
}

func DoMaxMreAmmoOperatoin(user *gameact.User, exit <-chan bool) error {
	plans := []int64{17, 13, 14, 15}
	return DoOperations(user, plans, exit)
}

func DoMaxResourceOperatoin(user *gameact.User, exit <-chan bool) error {
	plans := []int64{1, 2, 29, 30}
	// 资源足够则优先刷卡
	if user.GetResourceInfo().MP > 290000 && user.GetResourceInfo().MRE > 290000 {
		plans = []int64{1, 2, 15, 30}
	}
	return DoOperations(user, plans, exit)
}

func DoSimpleOperation(user *gameact.User, exit <-chan bool) error {
	plans := []int64{5, 6, 7, 8}
	return DoOperations(user, plans, exit)
}

func Daily(user *gameact.User, exit <-chan bool) error {
	for {
		user.Attendance()
		user.QuickGetQuestsResourceInMails(7)
		deadline := time.After(time.Hour * 8)
		select {
		case <-exit:
			return nil
		case <-deadline:
			continue
		}
	}
}

func DoOperations(user *gameact.User, plans []int64, exit <-chan bool) error {
	operations := map[int64]gameact.Operation{
		1:  {1, 40, 4, 1 * 50 * 60},
		2:  {2, 45, 5, 3 * 60 * 60},
		3:  {3, 45, 5, 12 * 60 * 60},
		4:  {4, 50, 5, 24 * 60 * 60},
		5:  {5, 1, 2, 1 * 15 * 60},
		6:  {6, 3, 2, 1 * 30 * 60},
		7:  {7, 5, 3, 1 * 60 * 60},
		8:  {8, 6, 5, 2 * 60 * 60},
		9:  {9, 5, 3, 1 * 40 * 60},
		10: {10, 8, 4, 1 * 90 * 60},
		11: {11, 10, 5, 4 * 60 * 60},
		12: {12, 15, 5, 6 * 60 * 60},
		13: {13, 12, 4, 1 * 20 * 60},
		14: {14, 20, 5, 1 * 45 * 60},
		15: {15, 15, 4, 1 * 90 * 60},
		16: {16, 25, 5, 5 * 60 * 60},
		17: {17, 30, 4, 1 * 60 * 60},
		18: {18, 35, 5, 2 * 60 * 60},
		19: {19, 40, 5, 6 * 60 * 60},
		20: {20, 40, 5, 8 * 60 * 60},
		21: {21, 30, 4, 1 * 30 * 60},
		22: {22, 35, 5, 5 * 30 * 60},
		23: {23, 40, 5, 4 * 60 * 60},
		24: {24, 40, 5, 7 * 60 * 60},
		25: {25, 35, 5, 2 * 60 * 60},
		26: {26, 40, 5, 3 * 60 * 60},
		27: {27, 45, 5, 5 * 60 * 60},
		28: {28, 45, 5, 12 * 60 * 60},
		29: {29, 40, 5, 5 * 30 * 60},
		30: {30, 45, 5, 4 * 60 * 60},
		31: {31, 50, 5, 11 * 30 * 60},
		32: {32, 50, 5, 8 * 60 * 60},
	}

	for {
		_, err := user.AbortUnmatchedOperations(plans)
		if err == nil {
			break
		}
		log.Println("ERR Abort operation error, sleep 1 minute.")
		time.Sleep(1 * time.Minute)
	}

	user.Update()
	levels := user.GetTeamLevels()
	populations := user.GetTeamPopulation()

	exitSignals := make([]chan bool, 0, len(plans)+1)

	for i := int64(0); i < int64(len(plans)) && i < int64(len(levels)) && i < int64(len(populations)); i++ {
		if levels[i+1] < operations[plans[i]].Level {
			continue
		}
		if populations[i+1] < operations[plans[i]].Population {
			continue
		}
		signal := make(chan bool, 1)
		exitSignals = append(exitSignals, signal)
		go InfOperationLoop(user, i+1, plans[i], operations[plans[i]].Duration, signal)
	}
	dailyExitSignal := make(chan bool, 1)
	exitSignals = append(exitSignals, dailyExitSignal)
	go Daily(user, dailyExitSignal)
	select {
	case k := <-exit:
		for _, s := range exitSignals {
			s <- k
		}
		return nil
	}
}

func UserOperation(user *gameact.User, wg *sync.WaitGroup, needPrepare bool, f func(*gameact.User, <-chan bool) error) error {
	defer wg.Done()

	for user.NbUid.Sign == "" {
		// 只对已经登陆用户处理!
		log.Println("Not login")
		return errors.New("not login")
	}
	for {
		err := user.Login(false, true)
		if err != nil {
			// 如果未登陆则直接返回错误
			if err == gameact.ErrNotLogin {
				return err
			}
			log.Println("Error when login. Will retry after 1 minute.", err)
			time.Sleep(1 * time.Minute)
		} else {
			break
		}
	}

	for failCounter := 0; needPrepare; failCounter++ {
		err := BasicPrepare(user)
		if err == nil {
			break
		}
		log.Println("ERR Basic prepare failed, sleep 1 min", err, user.LoginIdentify)
		time.Sleep(1 * time.Minute)
		if failCounter >= MaxErrCount {
			user.Login(false, true)
			failCounter = 0
		}
	}

	signal := make(chan bool, 1)
	f(user, signal)

	return nil
}

func MaxResources(user *gameact.User, wg *sync.WaitGroup) error {
	return UserOperation(user, wg, true, DoMaxResourceOperatoin)
}

func MaxMreAmmo(user *gameact.User, wg *sync.WaitGroup) error {
	return UserOperation(user, wg, true, DoMaxMreAmmoOperatoin)
}
