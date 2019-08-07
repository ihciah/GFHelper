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
	"fmt"
	"github.com/tidwall/gjson"
	"math/rand"
	"strings"
	"time"
)

func codeToError(res string) error {
	if strings.HasPrefix(res, "error:203") {
		return ErrGame203
	}
	if strings.HasPrefix(res, "error:3") {
		return ErrGame3
	}
	if strings.HasPrefix(res, "error:2") {
		return ErrGame2
	}
	if strings.HasPrefix(res, "error") {
		return ErrGame
	}
	return nil
}

type SpotsDeploy struct {
	SpotID int64 `json:"spot_id"`
	TeamID int64 `json:"team_id"`
}

type StartMissionReq struct {
	MissionID int64         `json:"mission_id"`
	Spots     []SpotsDeploy `json:"spots"`
	AllyID    int64         `json:"ally_id"`
}

func (u *User) StartMission(missionId int64, deploy []SpotsDeploy) (gjson.Result, error) {
	req := StartMissionReq{MissionID: missionId, Spots: deploy, AllyID: time.Now().Unix()}
	res, err := u.SendActionStruct(u.ConstructURL("Mission", "startMission"), req, nil, false)
	if e := codeToError(res); e != nil {
		u.logger.Printf("Err startmission: %s", res)
		return gjson.Result{}, e
	}
	u.Logger.Printf("Mission %d started.\n", missionId)
	return gjson.Parse(res), err
}

type SupplyTeamReq struct {
	TeamID int64 `json:"team_id"`
}

func (u *User) SupplyTeam(teamId int64) (gjson.Result, error) {
	req := SupplyTeamReq{TeamID: teamId}
	res, err := u.SendActionStruct(u.ConstructURL("Mission", "supplyTeam"), req, nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Team %d got supply.\n", teamId)
	return gjson.Parse(res), err
}

type TeamMoveReq struct {
	TeamID     int64 `json:"team_id"`
	FromSpotID int64 `json:"from_spot_id"`
	ToSpotID   int64 `json:"to_spot_id"`
	MoveType   int64 `json:"move_type"`
}

func (u *User) TeamMove(teamId, fromSpotId, toSpotId, moveType int64) (gjson.Result, error) {
	req := TeamMoveReq{TeamID: teamId, FromSpotID: fromSpotId, ToSpotID: toSpotId, MoveType: moveType}
	res, err := u.SendActionStruct(u.ConstructURL("Mission", "teamMove"), req, nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Team %d moved from location %d to %d, moving type %d.\n", teamId, fromSpotId, toSpotId, moveType)
	return gjson.Parse(res), err
}

type Gun struct {
	ID   int64 `json:"id"`
	Life int64 `json:"life"`
}

type BattleData struct {
	N1000          map[int64]int64             `json:"1000"`
	N1001          map[int64]int64             `json:"1001"`
	N1002          map[int64](map[int64]int64) `json:"1002"`
	N1003          map[int64]int64             `json:"1003"`
	BattleDamage   map[int64]int64             `json:"battle_damage"`
	BossHP         int64                       `json:"boss_hp"`
	CurrentTime    int64                       `json:"current_time"`
	Guns           []Gun                       `json:"guns"`
	UserRec        string                      `json:"user_rec"`
	IfEnemyDie     bool                        `json:"if_enemy_die"`
	LastBattleInfo string                      `json:"last_battle_info"`
	MVP            int64                       `json:"mvp"`
	SpotID         int64                       `json:"spot_id"`
}

func (u *User) BattleFinish(spotId, effect, enemyEffectClient, trueTime, lifeEnemy, clientTime, enemyCharacterTypeId,
	friendlyDamageCount, maxDamage int64, guns []Gun) (gjson.Result, error) {
	if len(guns) == 0 {
		return gjson.Result{}, ErrInputError
	}
	battleData := BattleData{
		SpotID:      spotId,
		MVP:         guns[0].ID,
		IfEnemyDie:  true,
		UserRec:     fmt.Sprintf("{\"seed\":%d,\"record\":[]}", rand.Intn(7654321)+1000000),
		Guns:        guns,
		CurrentTime: time.Now().Unix(),
	}
	battleData.N1003 = make(map[int64]int64)
	battleData.N1002 = make(map[int64]map[int64]int64)
	for _, team := range guns {
		battleData.N1002[team.ID] = make(map[int64]int64)
		battleData.N1002[team.ID][47] = 0
	}
	battleData.N1001 = make(map[int64]int64)
	battleData.N1000 = make(map[int64]int64)
	battleData.N1000[10] = effect
	battleData.N1000[11] = effect
	battleData.N1000[12] = effect
	battleData.N1000[13] = effect
	battleData.N1000[15] = enemyEffectClient
	battleData.N1000[16] = 0
	battleData.N1000[17] = trueTime
	battleData.N1000[18] = 0
	battleData.N1000[19] = 0
	battleData.N1000[20] = 0
	battleData.N1000[21] = 0
	battleData.N1000[22] = 0
	battleData.N1000[23] = 0
	battleData.N1000[24] = lifeEnemy
	battleData.N1000[25] = 0
	battleData.N1000[26] = lifeEnemy
	battleData.N1000[27] = clientTime
	battleData.N1000[33] = enemyCharacterTypeId
	battleData.N1000[34] = 0
	battleData.N1000[35] = 0
	battleData.N1000[40] = friendlyDamageCount
	battleData.N1000[41] = maxDamage
	battleData.N1000[42] = 0
	battleData.N1000[43] = 0
	battleData.N1000[44] = 0
	res, err := u.SendActionStruct(u.ConstructURL("Mission", "battleFinish"), battleData, nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Battle finished.\n")
	return gjson.Parse(res), err
}

func (u *User) EndTurn() (gjson.Result, error) {
	res, err := u.SendAction(u.ConstructURL("Mission", "endTurn"), "", nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("End turn.\n")
	return gjson.Parse(res), err
}

func (u *User) EndEnemyTurn() (gjson.Result, error) {
	res, err := u.SendAction(u.ConstructURL("Mission", "endEnemyTurn"), "", nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("End enemy turn.\n")
	return gjson.Parse(res), err
}

func (u *User) StartTurn() (gjson.Result, error) {
	res, err := u.SendAction(u.ConstructURL("Mission", "startTurn"), "", nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Start turn.\n")
	return gjson.Parse(res), err
}

type ReinforceTeamReq struct {
	SpotID int64 `json:"spot_id"`
	TeamID int64 `json:"team_id"`
}

func (u *User) ReinforceTeam(teamId, spotId int64) (gjson.Result, error) {
	res, err := u.SendActionStruct(u.ConstructURL("Mission", "reinforceTeam"), ReinforceTeamReq{spotId, teamId}, nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Reinforce team %d at %d.\n", teamId, spotId)
	return gjson.Parse(res), err
}

func (u *User) AbortMission() (gjson.Result, error) {
	res, err := u.SendAction(u.ConstructURL("Mission", "abortMission"), "", nil, false)
	if e := codeToError(res); e != nil {
		return gjson.Result{}, e
	}
	u.Logger.Printf("Abort mission.\n")
	return gjson.Parse(res), err
}
