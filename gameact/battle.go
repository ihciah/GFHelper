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
	"github.com/tidwall/gjson"
)

type RoutePair struct {
	TeamIndex int64 // start from 0
	Dst       int64
}

type LevelConfig struct {
	MissionId            int64
	Effect               int64
	EnemyEffectClient    int64
	TrueTime             int64
	LifeEnemy            int64
	ClientTime           int64
	EnemyCharacterTypeId int64
	FriendlyDamageCount  int64
	MaxDamage            int64
	Route                [][]RoutePair
}

type BattleData struct {
	U             *User
	LevelConfig   LevelConfig
	TeamId        []int64
	TeamPosition  map[int64]int64 // team index -> pos
	Battleguns    [][]protocol.Gun
	EnemyPosition map[int64]int64
	TurnStarted   bool
	//TeamPosition  map[int64]int64
}

// 生成BattleData
func (u *User) MakeBattle(teamIds []int64, levelConfig LevelConfig) (*BattleData, error) {
	battleGuns := make([][]protocol.Gun, 0, len(teamIds))
	for i := range teamIds {
		guns, err := u.GetTeam(teamIds[i])
		if err != nil {
			return nil, err
		}
		singleBattleGun := make([]protocol.Gun, 0, len(guns))
		for _, g := range guns {
			singleBattleGun = append(singleBattleGun, protocol.Gun{g.GunId, g.GunLife})
		}
		battleGuns = append(battleGuns, singleBattleGun)
	}
	teamPositions := make(map[int64]int64)
	enemyPositions := make(map[int64]int64)
	return &BattleData{U: u, TeamId: teamIds, LevelConfig: levelConfig, TeamPosition: teamPositions,
		Battleguns: battleGuns, EnemyPosition: enemyPositions, TurnStarted: true}, nil
}

// 开始任务
// 使用LevelConfig的第一个元素作为初始部署位置
func (b *BattleData) DoStartMission() error {
	if len(b.LevelConfig.Route) < 1 || len(b.LevelConfig.Route[0]) < 1 {
		return ErrUnexpectedData
	}
	deploys := make([]protocol.SpotsDeploy, 0, len(b.TeamId))
	for _, pair := range b.LevelConfig.Route[0] {
		if pair.TeamIndex >= int64(len(b.TeamId)) {
			return ErrUnexpectedData
		}
		deploys = append(deploys, protocol.SpotsDeploy{pair.Dst, b.TeamId[pair.TeamIndex]})
		b.TeamPosition[pair.TeamIndex] = pair.Dst
	}
	result, err := b.U.StartMission(b.LevelConfig.MissionId, deploys)
	if err != nil {
		return err
	}
	enemySpotIds := result.Get("spot_act_info.#[enemy_team_id!=\"0\"]#.spot_id").Array()
	for _, spotId := range enemySpotIds {
		if spotId.Exists() {
			b.EnemyPosition[spotId.Int()] = 1
		}
	}
	return nil
}

// 结束己方回合
func (b *BattleData) DoEndTurn() (bool, error) {
	result, err := b.U.EndTurn()
	if err != nil {
		return false, err
	}
	isEnd := result.Get("mission_lose_result.turn").Exists() || result.Get("mission_win_result.rank").Exists()
	if isEnd {
		b.U.Logger.Println("Mission Finished.")
		return true, nil
	}
	growEnemyPos := result.Get("grow_enemy.#.spot_id").Array()
	for _, pos := range growEnemyPos {
		b.EnemyPosition[pos.Int()] = 1
	}
	from := result.Get("enemy_move.#.from_spot_id").Array()
	to := result.Get("enemy_move.#.to_spot_id").Array()
	if len(from) != len(to) {
		return false, ErrUnexpectedData
	}
	for i := range from {
		b.EnemyPosition[from[i].Int()] = b.EnemyPosition[from[i].Int()] - 1
		b.EnemyPosition[to[i].Int()] = b.EnemyPosition[to[i].Int()] + 1
	}
	for i, teamPos := range b.TeamPosition {
		if b.EnemyPosition[teamPos] != 0 {
			_, err := b.U.BattleFinish(teamPos, b.LevelConfig.Effect, b.LevelConfig.EnemyEffectClient, b.LevelConfig.TrueTime, b.LevelConfig.LifeEnemy,
				b.LevelConfig.ClientTime, b.LevelConfig.EnemyCharacterTypeId, b.LevelConfig.FriendlyDamageCount, b.LevelConfig.MaxDamage, b.Battleguns[i])
			if err != nil {
				return false, err
			}
			b.EnemyPosition[teamPos] = 0
		}
	}

	return false, nil
}

// 结束敌方回合
func (b *BattleData) DoEndEnemyTurn() (bool, error) {
	result, err := b.U.EndEnemyTurn()
	if err != nil {
		return false, err
	}
	isEnd := result.Get("mission_lose_result.turn").Exists() || result.Get("mission_win_result.rank").Exists()
	if isEnd {
		b.U.Logger.Println("Mission Finished.")
	}
	return isEnd, nil
}

// 开始己方回合
func (b *BattleData) DoStartTurn() (bool, error) {
	result, err := b.U.StartTurn()
	if err != nil {
		return false, err
	}
	isEnd := result.Get("mission_lose_result.turn").Exists() || result.Get("mission_win_result.rank").Exists()
	if isEnd {
		b.U.Logger.Println("Mission Finished.")
	}
	return isEnd, nil
}

// 移动/部署队伍，并更新TeamPosition
func (b *BattleData) DoMove(routePair RoutePair) (err error) {
	var moveResult gjson.Result
	battleNow := false
	if _, exists := b.TeamPosition[routePair.TeamIndex]; !exists {
		err = b.DoReinforce(routePair)
	} else {
		moveResult, err = b.U.TeamMove(b.TeamId[routePair.TeamIndex], b.TeamPosition[routePair.TeamIndex], routePair.Dst, 1)
		if err == nil && moveResult.Get("enemy_team_id").Exists() {
			battleNow = true
		}
	}
	if err != nil {
		return
	}
	b.TeamPosition[routePair.TeamIndex] = routePair.Dst
	if b.EnemyPosition[routePair.Dst] != 0 || battleNow {
		_, err = b.U.BattleFinish(routePair.Dst, b.LevelConfig.Effect, b.LevelConfig.EnemyEffectClient, b.LevelConfig.TrueTime, b.LevelConfig.LifeEnemy,
			b.LevelConfig.ClientTime, b.LevelConfig.EnemyCharacterTypeId, b.LevelConfig.FriendlyDamageCount, b.LevelConfig.MaxDamage, b.Battleguns[routePair.TeamIndex])
		if err != nil {
			return err
		}
		b.EnemyPosition[routePair.Dst] = 0
	}
	return nil
}

// 部署队伍
func (b *BattleData) DoReinforce(pair RoutePair) error {
	if pair.TeamIndex >= int64(len(b.TeamId)) {
		return ErrUnexpectedData
	}
	teamId := b.TeamId[pair.TeamIndex]
	_, err := b.U.ReinforceTeam(teamId, pair.Dst)
	if err != nil {
		return err
	}
	b.TeamPosition[pair.TeamIndex] = pair.Dst
	return nil
}

// 完成整个任务
func (b *BattleData) WinBattle() error {
	// 弹药口粮检查
	b.U.Update()
	resources := b.U.GetResourceInfo()
	if resources.MRE < MinMRE || resources.AMMO < MinAMMO || resources.MP < MinMP {
		b.U.Logger.Printf("Insufficient MRE(%d) or AMMO(%d) or MP(%d)", resources.MRE, resources.AMMO, resources.MP)
		return ErrResourcesInsufficient
	}
	b.U.AbortAllOperations()
	err := b.DoStartMission()
	if err != nil {
		return err
	}
	for _, move := range b.LevelConfig.Route[1:] {
		if !b.TurnStarted {
			isEnd, err := b.DoStartTurn()
			if err != nil || isEnd {
				return err
			}
		}
		for _, rpos := range move {
			if rpos.Dst < 0 {
				if WithSupply {
					b.U.SupplyTeam(b.TeamId[rpos.TeamIndex])
				}
			} else {
				err = b.DoMove(rpos)
				if err != nil {
					return err
				}
			}
		}
		isEnd, err := b.DoEndTurn()
		if err != nil || isEnd {
			return err
		}
		b.TurnStarted = false
		isEnd, err = b.DoEndEnemyTurn()
		if err != nil || isEnd {
			return err
		}
	}
	b.U.Logger.Println("Unknown error.")
	return ErrUnexpectedData
}
