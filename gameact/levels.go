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

func checkLength(team_id []int64, l int64) error {
	if int64(len(team_id)) < l {
		return ErrUnexpectedData
	}
	return nil
}

func (u *User) Battle_1_1(team_id []int64) error {
	if err := checkLength(team_id, 1); err != nil {
		return err
	}
	route := [][]RoutePair{
		{
			{0, 88},
		},
		{
			{0, -1},
			{0, 89},
			{0, 90},
		},
		{},
	}
	level_config := LevelConfig{MissionId: 5, Effect: 716, EnemyEffectClient: 82, TrueTime: 145, LifeEnemy: 50, ClientTime: 27,
		EnemyCharacterTypeId: 0, FriendlyDamageCount: 3, MaxDamage: 16, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_1_2(team_id []int64) error {
	if err := checkLength(team_id, 1); err != nil {
		return err
	}
	route := [][]RoutePair{
		{
			{0, 92},
		},
		{
			{0, -1},
			{0, 93},
			{0, 94},
		},
		{
			{0, -1},
			{0, 95},
			{0, 96},
		},
	}
	level_config := LevelConfig{MissionId: 6, Effect: 830, EnemyEffectClient: 264, TrueTime: 253, LifeEnemy: 388, ClientTime: 9,
		EnemyCharacterTypeId: 20001, FriendlyDamageCount: 31, MaxDamage: 12, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_1_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 99},
		},
		{
			{0, -1},
			{0, 100},
			{0, 101},
		},
		{
			{0, 106},
			{0, 107},
		},
	}
	level_config := LevelConfig{MissionId: 7, Effect: 1023, EnemyEffectClient: 480, TrueTime: 265, LifeEnemy: 415, ClientTime: 9,
		EnemyCharacterTypeId: 20002, FriendlyDamageCount: 28, MaxDamage: 14, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_1_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 108},
		},
		{
			{0, -1},
			{0, 110},
			{0, 113},
		},
		{
			{0, 114},
			{0, 116},
		},
		{
			{0, 118},
		},
	}
	level_config := LevelConfig{MissionId: 8, Effect: 2656, EnemyEffectClient: 330, TrueTime: 143, LifeEnemy: 388, ClientTime: 5,
		EnemyCharacterTypeId: 20002, FriendlyDamageCount: 16, MaxDamage: 24, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_1_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 119},
		},
		{
			{0, -1},
			{0, 121},
			{0, 124},
		},
		{
			{0, -1},
			{0, 128},
			{0, 131},
			{0, 132},
		},
	}
	level_config := LevelConfig{MissionId: 9, Effect: 2656, EnemyEffectClient: 286, TrueTime: 191, LifeEnemy: 606, ClientTime: 7,
		EnemyCharacterTypeId: 20004, FriendlyDamageCount: 30, MaxDamage: 20, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_1_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 133},
		},
		{
			{0, -1},
			{0, 134},
			{0, 135},
		},
		{
			{0, 139},
			{0, 136},
		},
		{
			{0, -1},
			{0, 144},
			{0, 148},
			{0, 149},
		},
		{
			{0, 146},
			{0, 147},
		},
	}
	level_config := LevelConfig{MissionId: 10, Effect: 2661, EnemyEffectClient: 435, TrueTime: 200, LifeEnemy: 519, ClientTime: 7,
		EnemyCharacterTypeId: 90002, FriendlyDamageCount: 25, MaxDamage: 20, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_1(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 202},
		},
		{
			{0, -1},
			{0, 205},
			{0, 206},
		},
		{
			{0, -1},
			{0, 209},
			{0, 347},
		},
	}
	level_config := LevelConfig{MissionId: 15, Effect: 2661, EnemyEffectClient: 385, TrueTime: 174, LifeEnemy: 185, ClientTime: 6,
		EnemyCharacterTypeId: 20003, FriendlyDamageCount: 19, MaxDamage: 9, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_2(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 210},
		},
		{
			{0, -1},
			{0, 213},
			{0, 212},
		},
		{
			{0, 215},
			{0, 214},
		},
		{
			{0, 217},
			{0, 219},
		},
	}
	level_config := LevelConfig{MissionId: 16, Effect: 3367, EnemyEffectClient: 645, TrueTime: 239, LifeEnemy: 1646, ClientTime: 9,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 51, MaxDamage: 32, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 220},
		},
		{
			{0, -1},
			{0, 222},
			{0, 226},
			{0, 228},
		},
		{
			{0, 224},
		},
	}
	level_config := LevelConfig{MissionId: 17, Effect: 4096, EnemyEffectClient: 781, TrueTime: 280, LifeEnemy: 1178, ClientTime: 9,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 43, MaxDamage: 27, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 233},
		},
		{
			{0, -1},
			{0, 234},
			{0, 236},
		},
		{
			{0, 239},
			{0, 241},
		},
		{
			{0, 243},
		},
	}
	level_config := LevelConfig{MissionId: 18, Effect: 4096, EnemyEffectClient: 689, TrueTime: 377, LifeEnemy: 1831, ClientTime: 14,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 48, MaxDamage: 38, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 245},
		},
		{
			{0, -1},
			{0, 247},
			{0, 252},
		},
		{
			{0, -1},
			{0, 251},
			{0, 259},
		},
	}
	level_config := LevelConfig{MissionId: 19, Effect: 4096, EnemyEffectClient: 1074, TrueTime: 517, LifeEnemy: 2504, ClientTime: 18,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 72, MaxDamage: 34, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_2_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 260},
		},
		{
			{0, -1},
			{0, 261},
			{0, 263},
		},
		{
			{0, 267},
			{0, 271},
		},
	}
	level_config := LevelConfig{MissionId: 20, Effect: 5807, EnemyEffectClient: 1008, TrueTime: 254, LifeEnemy: 2860, ClientTime: 9,
		EnemyCharacterTypeId: 90004, FriendlyDamageCount: 92, MaxDamage: 31, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_1(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 348},
		},
		{
			{0, -1},
			{0, 350},
			{0, 353},
		},
		{
			{0, 356},
		},
	}
	level_config := LevelConfig{MissionId: 25, Effect: 5807, EnemyEffectClient: 788, TrueTime: 200, LifeEnemy: 2205, ClientTime: 7,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 68, MaxDamage: 32, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_2(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 358},
		},
		{
			{0, -1},
			{0, 362},
			{0, 363},
			{1, 358},
			{1, -1},
		},
		{
			{0, -1},
			{1, -1},
			{0, 362},
			{0, 365},
			{0, 366},
			{0, 368},
			{0, 369},
		},
		{
			{0, 364},
			{0, 367},
		},
	}
	level_config := LevelConfig{MissionId: 26, Effect: 5807, EnemyEffectClient: 1350, TrueTime: 282, LifeEnemy: 3040, ClientTime: 10,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 103, MaxDamage: 29, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 382},
		},
		{
			{0, -1},
			{0, 384},
			{1, 382},
		},
		{
			{0, -1},
			{1, -1},
			{0, 383},
			{0, 378},
			{0, 373},
			{0, 372},
		},
		{},
	}
	level_config := LevelConfig{MissionId: 27, Effect: 5807, EnemyEffectClient: 1500, TrueTime: 340, LifeEnemy: 5080, ClientTime: 12,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 162, MaxDamage: 31, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 385},
		},
		{
			{0, -1},
			{0, 388},
			{1, 385},
		},
		{
			{1, -1},
			{0, 392},
			{0, 395},
		},
	}
	level_config := LevelConfig{MissionId: 28, Effect: 5807, EnemyEffectClient: 1661, TrueTime: 368, LifeEnemy: 4804, ClientTime: 12,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 150, MaxDamage: 32, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 412},
			{1, 398},
		},
		{
			{0, -1},
			{1, -1},
			{0, 411},
			{0, 414},
			{0, 413},
		},
	}
	level_config := LevelConfig{MissionId: 29, Effect: 5807, EnemyEffectClient: 2495, TrueTime: 704, LifeEnemy: 8750, ClientTime: 24,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 263, MaxDamage: 33, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_3_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 415},
			{1, 431},
		},
		{
			{0, -1},
			{1, -1},
			{0, 417},
			{0, 418},
			{0, 419},
			{0, 425},
		},
	}
	level_config := LevelConfig{MissionId: 30, Effect: 6807, EnemyEffectClient: 1644, TrueTime: 454, LifeEnemy: 4200, ClientTime: 15,
		EnemyCharacterTypeId: 90006, FriendlyDamageCount: 194, MaxDamage: 30, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_1(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 513},
		},
		{
			{0, -1},
			{0, 512},
			{0, 511},
			{0, 509},
		},
	}
	level_config := LevelConfig{MissionId: 35, Effect: 10021, EnemyEffectClient: 1680, TrueTime: 114, LifeEnemy: 4341, ClientTime: 5,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 47, MaxDamage: 92, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_2(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 522},
		},
		{
			{0, -1},
			{0, 521},
			{0, 517},
		},
		{
			{0, 518},
			{0, 519},
		},
	}
	level_config := LevelConfig{MissionId: 36, Effect: 10021, EnemyEffectClient: 1855, TrueTime: 82, LifeEnemy: 5720, ClientTime: 5,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 58, MaxDamage: 98, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 528},
		},
		{
			{0, -1},
			{0, 529},
			{0, 533},
			{0, 529},
		},
		{
			{0, 533},
			{0, 537},
			{0, 540},
		},
		{
			{0, 542},
			{0, 544},
		},
	}
	level_config := LevelConfig{MissionId: 37, Effect: 10021, EnemyEffectClient: 1456, TrueTime: 104, LifeEnemy: 4850, ClientTime: 5,
		EnemyCharacterTypeId: 20008, FriendlyDamageCount: 48, MaxDamage: 101, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 545},
			{1, 549},
		},
		{
			{0, -1},
			{1, -1},
			{0, 546},
			{0, 551},
			{0, 550},
			{0, 555},
			{0, 558},
		},
	}
	level_config := LevelConfig{MissionId: 38, Effect: 10021, EnemyEffectClient: 2523, TrueTime: 135, LifeEnemy: 6815, ClientTime: 5,
		EnemyCharacterTypeId: 20009, FriendlyDamageCount: 77, MaxDamage: 88, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 577},
			{1, 562},
		},
		{
			{0, -1},
			{1, -1},
			{0, 578},
			{0, 579},
			{0, 580},
			{0, 581},
		},
	}
	level_config := LevelConfig{MissionId: 39, Effect: 10021, EnemyEffectClient: 2048, TrueTime: 114, LifeEnemy: 4719, ClientTime: 5,
		EnemyCharacterTypeId: 20009, FriendlyDamageCount: 54, MaxDamage: 87, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 582},
			{1, 587},
		},
		{
			{0, -1},
			{1, -1},
			{0, 588},
			{0, 594},
			{0, 598},
			{0, 604},
		},
	}
	level_config := LevelConfig{MissionId: 40, Effect: 10021, EnemyEffectClient: 8846, TrueTime: 201, LifeEnemy: 16656, ClientTime: 8,
		EnemyCharacterTypeId: 90008, FriendlyDamageCount: 172, MaxDamage: 96, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_1e(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 612},
			{1, 609},
		},
		{
			{0, -1},
			{1, -1},
			{0, 613},
			{0, 617},
			{1, 610},
		},
		{
			{0, -1},
			{0, 620},
			{0, 622},
			{0, 624},
			{0, 625},
		},
	}
	level_config := LevelConfig{MissionId: 41, Effect: 32514, EnemyEffectClient: 3015, TrueTime: 98, LifeEnemy: 9465, ClientTime: 6,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 108, MaxDamage: 87, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_2e(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 642},
			{1, 630},
		},
		{
			{0, -1},
			{1, -1},
			{0, 641},
			{0, 637},
			{0, 640},
			{0, 639},
		},
	}
	level_config := LevelConfig{MissionId: 42, Effect: 32514, EnemyEffectClient: 2570, TrueTime: 109, LifeEnemy: 6865, ClientTime: 6,
		EnemyCharacterTypeId: 20008, FriendlyDamageCount: 70, MaxDamage: 98, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_3e(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 658},
			{1, 643},
		},
		{
			{0, -1},
			{1, -1},
			{0, 659},
			{0, 660},
			{0, 661},
			{0, 662},
		},
	}
	level_config := LevelConfig{MissionId: 43, Effect: 32514, EnemyEffectClient: 1597, TrueTime: 108, LifeEnemy: 3567, ClientTime: 4,
		EnemyCharacterTypeId: 20009, FriendlyDamageCount: 43, MaxDamage: 82, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_4_4e(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 664},
			{1, 669},
		},
		{
			{0, -1},
			{1, -1},
			{0, 670},
			{0, 676},
			{0, 680},
			{0, 686},
		},
	}
	level_config := LevelConfig{MissionId: 44, Effect: 32514, EnemyEffectClient: 2303, TrueTime: 154, LifeEnemy: 3881, ClientTime: 5,
		EnemyCharacterTypeId: 20009, FriendlyDamageCount: 44, MaxDamage: 88, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_1(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 690},
		},
		{
			{0, -1},
			{0, 691},
			{1, 690},
			{1, -1},
		},
		{
			{1, -1},
			{0, 694},
			{1, 691},
			{1, 692},
		},
		{
			{1, 693},
			{1, 696},
		},
		{
			{1, -1},
			{1, 698},
			{1, 700},
			{1, 703},
			{1, 702},
		},
		{
			{1, 699},
		},
	}
	level_config := LevelConfig{MissionId: 45, Effect: 10021, EnemyEffectClient: 2378, TrueTime: 229, LifeEnemy: 6031, ClientTime: 8,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 64, MaxDamage: 94, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_2(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 716},
		},
		{
			{0, -1},
			{0, 720},
			{1, 716},
			{1, -1},
		},
		{
			{0, 721},
			{0, 725},
			{0, 718},
		},
	}
	level_config := LevelConfig{MissionId: 46, Effect: 10021, EnemyEffectClient: 1908, TrueTime: 124, LifeEnemy: 5720, ClientTime: 4,
		EnemyCharacterTypeId: 20006, FriendlyDamageCount: 61, MaxDamage: 93, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 726},
			{1, 744},
		},
		{
			{0, -1},
			{1, -1},
			{1, 740},
			{1, 741},
			{1, 745},
			{1, 747},
		},
		{
			{0, -1},
			{1, -1},
			{1, 748},
			{1, 746},
			{1, 748},
			{0, 727},
			{0, 731},
		},
		{
			{1, -1},
			{0, 735},
			{1, 746},
		},
		{
			{1, 743},
			{1, 738},
			{0, 737},
		},
		{
			{0, -1},
			{1, -1},
			{1, 739},
			{1, 736},
			{1, 733},
			{1, 729},
			{1, 728},
			{1, 732},
		},
	}
	level_config := LevelConfig{MissionId: 47, Effect: 20021, EnemyEffectClient: 2644, TrueTime: 297, LifeEnemy: 8619, ClientTime: 11,
		EnemyCharacterTypeId: 20004, FriendlyDamageCount: 104, MaxDamage: 82, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 771},
		},
		{
			{0, 772},
			{0, 767},
			{0, 768},
			{0, 763},
		},
	}
	level_config := LevelConfig{MissionId: 48, Effect: 20021, EnemyEffectClient: 3044, TrueTime: 84, LifeEnemy: 8640, ClientTime: 4,
		EnemyCharacterTypeId: 20004, FriendlyDamageCount: 94, MaxDamage: 91, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 801},
			{1, 799},
		},
		{
			{0, -1},
			{1, -1},
			{1, 797},
			{0, 802},
			{0, 798},
		},
		{
			{0, -1},
			{0, 796},
			{0, 792},
			{0, 789},
			{0, 783},
			{0, 780},
		},
		{
			{0, 778},
			{0, 777},
		},
	}
	level_config := LevelConfig{MissionId: 49, Effect: 32058, EnemyEffectClient: 3493, TrueTime: 120, LifeEnemy: 8464, ClientTime: 4,
		EnemyCharacterTypeId: 20005, FriendlyDamageCount: 106, MaxDamage: 79, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_5_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 803},
			{1, 807},
		},
		{
			{0, -1},
			{0, 808},
			{0, 813},
			{0, 820},
			{0, 826},
		},
	}
	level_config := LevelConfig{MissionId: 50, Effect: 32058, EnemyEffectClient: 11633, TrueTime: 187, LifeEnemy: 19416, ClientTime: 6,
		EnemyCharacterTypeId: 900017, FriendlyDamageCount: 199, MaxDamage: 97, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_1(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1511},
		},
		{
			{0, -1},
			{0, 1512},
			{1, 1511},
			{1, -1},
		},
		{
			{1, -1},
			{0, 1518},
			{1, 1512},
			{1, 1513},
		},
		{
			{0, 1517},
			{0, 1519},
			{0, 1516},
		},
		{
			{0, 1523},
		},
	}
	level_config := LevelConfig{MissionId: 55, Effect: 32058, EnemyEffectClient: 5636, TrueTime: 105, LifeEnemy: 9648, ClientTime: 4,
		EnemyCharacterTypeId: 10004, FriendlyDamageCount: 107, MaxDamage: 90, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_2(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1524},
		},
		{
			{0, -1},
			{0, 1529},
			{0, 1531},
		},
		{
			{0, -1},
			{0, 1533},
			{0, 1535},
			{0, 1537},
		},
	}
	level_config := LevelConfig{MissionId: 56, Effect: 32058, EnemyEffectClient: 5580, TrueTime: 83, LifeEnemy: 4660, ClientTime: 3,
		EnemyCharacterTypeId: 10005, FriendlyDamageCount: 58, MaxDamage: 80, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_3(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1543},
		},
		{
			{0, -1},
			{0, 1538},
			{0, 1539},
		},
		{
			{0, -1},
			{0, 1541},
			{0, 1542},
		},
	}
	level_config := LevelConfig{MissionId: 57, Effect: 32058, EnemyEffectClient: 6930, TrueTime: 147, LifeEnemy: 14628, ClientTime: 5,
		EnemyCharacterTypeId: 10006, FriendlyDamageCount: 157, MaxDamage: 93, Route: route}
	bd, err := u.MakeBattle(team_id[0:1], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_4(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1578},
		},
		{
			{0, -1},
			{0, 1750},
			{0, 1574},
		},
		{
			{1, 1578},
			{1, -1},
			{0, 1572},
		},
		{
			{0, -1},
			{1, -1},
			{0, 1574},
			{0, 1752},
			{0, 1575},
			{0, 1576},
		},
		{
			{1, -1},
			{0, 1569},
		},
	}
	level_config := LevelConfig{MissionId: 58, Effect: 32058, EnemyEffectClient: 9491, TrueTime: 156, LifeEnemy: 14603, ClientTime: 6,
		EnemyCharacterTypeId: 10004, FriendlyDamageCount: 143, MaxDamage: 103, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_5(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1592},
		},
		{
			{0, -1},
			{0, 1589},
			{1, 1592},
			{1, -1},
		},
		{
			{1, -1},
			{0, 1596},
			{1, 1601},
			{1, 1600},
		},
		{
			{0, -1},
			{1, -1},
			{0, 1586},
			{0, 1583},
			{0, 1581},
			{0, 1579},
		},
	}
	level_config := LevelConfig{MissionId: 59, Effect: 32058, EnemyEffectClient: 7496, TrueTime: 138, LifeEnemy: 18780, ClientTime: 5,
		EnemyCharacterTypeId: 10005, FriendlyDamageCount: 202, MaxDamage: 92, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

func (u *User) Battle_6_6(team_id []int64) error {
	route := [][]RoutePair{
		{
			{1, 1616},
			{0, 1618},
		},
		{
			{0, -1},
			{1, -1},
			{0, 1619},
			{0, 1623},
			{0, 1622},
			{1, 1634},
		},
		{
			{0, -1},
			{0, 1621},
			{0, 1636},
			{0, 1632},
			{0, 1633},
		},
	}
	level_config := LevelConfig{MissionId: 60, Effect: 32058, EnemyEffectClient: 18986, TrueTime: 210, LifeEnemy: 37505, ClientTime: 8,
		EnemyCharacterTypeId: 900033, FriendlyDamageCount: 320, MaxDamage: 117, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}

// More battle!
func (u *User) Battle_6_6sp(team_id []int64) error {
	route := [][]RoutePair{
		{
			{0, 1616},
		},
		{
			{0, -1},
			{0, 1634},
			{1, 1616},
			{1, -1},
			{1, 1617},
		},
		{
			{0, 1635},
			{1, 1620},
		},
		{
			{0, 1622},
			{1, 1621},
		},
		{
			{0, -1},
			{0, 1635},
			{0, 1620},
			{1, 1622},
			{1, -1},
			{1, 1623},
		},
		{
			{1, 1622},
			{1, -1},
			{1, 1621},
			{0, 1626},
		},
		{
			{0, 1627},
			{1, 1636},
			{1, 1632},
			{1, 1628},
		},
		{
			{0, -1},
			{0, 1631},
			{0, 1632},
			{0, 1633},
		},
	}
	level_config := LevelConfig{MissionId: 60, Effect: 32058, EnemyEffectClient: 18986, TrueTime: 210, LifeEnemy: 37505, ClientTime: 8,
		EnemyCharacterTypeId: 900033, FriendlyDamageCount: 320, MaxDamage: 117, Route: route}
	bd, err := u.MakeBattle(team_id[0:2], level_config)
	if err != nil {
		return err
	}
	err = bd.WinBattle()
	return err
}
