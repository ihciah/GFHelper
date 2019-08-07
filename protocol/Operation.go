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

import "strings"

// 开始后勤
type StartOperationReq struct {
	TeamID      int64 `json:"team_id"`
	OperationID int64 `json:"operation_id"`
	MaxLevel    int64 `json:"max_level"`
}

func (u *User) StartOperation(teamId, operationId, maxLevel int64) (err error) {
	data := StartOperationReq{TeamID: teamId, OperationID: operationId, MaxLevel: maxLevel}
	result, err := u.SendActionStruct(u.ConstructURL("Operation", "startOperation"), data, nil, true)
	if !strings.HasPrefix(result, "1") {
		u.Logger.Printf("Error when team %d strat operation %d: %s, user: %s", teamId, operationId, result, u.LoginIdentify)
		return ErrServerError
	}
	u.Logger.Printf("Operation %d with team %d started. User: %s", operationId, teamId, u.LoginIdentify)
	return
}

type OperationReq struct {
	OperationID int64 `json:"operation_id"`
}

// 取消后勤
func (u *User) AbortOperation(operationId int64) (err error) {
	data := OperationReq{operationId}
	result, err := u.SendActionStruct(u.ConstructURL("Operation", "abortOperation"), data, nil, true)
	if !strings.HasPrefix(result, "1") {
		u.Logger.Printf("Error when abort operation %d: %s, user: %s", operationId, result, u.LoginIdentify)
		return ErrServerError
	}
	u.Logger.Printf("Operation %d aborted. User: %s", operationId, u.LoginIdentify)
	return
}

// 完成后勤
func (u *User) FinishOperation(operationId int64) (err error) {
	data := OperationReq{operationId}
	result, err := u.SendActionStruct(u.ConstructURL("Operation", "finishOperation"), data, nil, true)
	if !strings.HasPrefix(result, "{") {
		u.Logger.Printf("Error when finish operation %d: %s, user: %s", operationId, result, u.LoginIdentify)
		return ErrServerError
	}
	u.Logger.Printf("Operation %d finished. User: %s", operationId, u.LoginIdentify)
	return
}
