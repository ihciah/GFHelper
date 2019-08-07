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
	"github.com/tidwall/gjson"
	"time"
)

// 获取可以执行的后勤列表
func (u *User) GetAvailableOperations() []int64 {
	//先判断玩到第几关了
	u.Update()
	finished := make(map[int64]int)
	for _, id := range u.GetFinishedMission() {
		finished[id] = 1
	}
	// 只判断前4组后勤
	until := int64(8)
	for i := int64(10); i < 35; i += 10 {
		if finished[i] == 1 {
			until += 4
		}
	}
	// 去掉正在执行的后勤
	operating := make(map[int64]int)
	for _, id := range u.GetOperatingIDs() {
		operating[id] = 1
	}
	available := make([]int64, 0, until-4)
	for k := int64(5); k <= until; k += 1 {
		if operating[k] == 0 {
			available = append(available, k)
		}
	}
	return available
}

type Operation struct {
	ID         int64
	Level      int64
	Population int64
	Duration   int64
}

func (u *User) AbortAllOperations() error {
	u.Update()
	for _, id := range u.GetOperatingIDs() {
		err := u.AbortOperation(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// 取消所有不匹配的任务，并返回正在执行的匹配的任务=team_id:start_time
func (u *User) AbortUnmatchedOperations(plan []int64) (map[int64]int64, error) {
	operationTime := map[int64]int64{
		5:  15 * 60,
		6:  30 * 60,
		7:  60 * 60,
		8:  2 * 60 * 60,
		9:  40 * 60,
		10: 90 * 60,
		11: 4 * 60 * 60,
		12: 6 * 60 * 60,
		13: 20 * 60,
		14: 45 * 60,
		15: 90 * 60,
		16: 5 * 60 * 60,
		17: 60 * 60,
		18: 2 * 60 * 60,
		19: 6 * 60 * 60,
		20: 8 * 60 * 60,
	}

	u.Update()
	operatingIds := u.IndexInfo.Get("operation_act_info.#.operation_id").Array()
	operatingTeams := u.IndexInfo.Get("operation_act_info.#.team_id").Array()
	operatingStartTimes := u.IndexInfo.Get("operation_act_info.#.start_time").Array()
	ret := make(map[int64]int64)

	if len(operatingTeams) != len(operatingIds) || len(operatingStartTimes) != len(operatingIds) {
		return ret, ErrUnexpectedData
	}

	for i := range operatingIds {
		opId := operatingIds[i].Int()
		opTeam := operatingTeams[i].Int()
		if opTeam <= int64(len(plan)) && opTeam > 0 {
			// 符合要求的op_id直接取消或返回，不符合的不管
			if opId != plan[opTeam-1] {
				execTime, ok := operationTime[opId]
				var err error
				if ok && execTime+operatingStartTimes[i].Int()+TimePadding < time.Now().Unix() {
					u.Logger.Printf("Finish operation %d team %d account %s", opId, opTeam, u.LoginIdentify)
					err = u.FinishOperation(opId)
					u.IndexInfoDirty = true
				} else {
					u.Logger.Printf("Abort operation %d team %d account %s", opId, opTeam, u.LoginIdentify)
					err = u.AbortOperation(opId)
					u.IndexInfoDirty = true
				}
				if err != nil {
					return ret, err
				}
			} else {
				ret[opTeam] = operatingStartTimes[i].Int()
			}
		}
	}
	return ret, nil
}

// 返回team_id正在执行的任务
func (u *User) CheckTeam(teamId int64) (int64, error) {
	u.Update()
	operatingIds := u.IndexInfo.Get("operation_act_info.#.operation_id").Array()
	operatingTeams := u.IndexInfo.Get("operation_act_info.#.team_id").Array()
	if len(operatingTeams) != len(operatingIds) {
		return 0, ErrUnexpectedData
	}
	for i := range operatingIds {
		if operatingTeams[i].Int() == teamId {
			return operatingIds[i].Int(), nil
		}
	}
	return 0, nil
}

func (u *User) StartOperation(teamId, operationId, maxLevel int64) (err error) {
	err = u.User.StartOperation(teamId, operationId, maxLevel)
	u.IndexInfoDirty = true
	return
}

func (u *User) AbortOperation(operationId int64) (err error) {
	err = u.User.AbortOperation(operationId)
	u.IndexInfoDirty = true
	return
}

func (u *User) FinishOperation(operationId int64) (err error) {
	err = u.User.FinishOperation(operationId)
	u.IndexInfoDirty = true
	return
}

func (u *User) AbortMission() (gjson.Result, error) {
	u.Update()
	if u.IndexInfo.Get("mission_act_info.spot").Exists() {
		gres, err := u.User.AbortMission()
		u.IndexInfoDirty = true
		return gres, err
	}
	return gjson.Result{}, nil
}
