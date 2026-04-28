package unit

import (
	"battle/internal/battle/attributes"
	"battle/internal/battle/skill"
	"errors"
	"fmt"

	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/config"
)

var (
	ErrNilWorld      = errors.New("unit: world is nil")
	ErrNilPlayer     = errors.New("unit: player is nil")
	ErrNoPlayerUnits = errors.New("unit: player.Units 为空或全部为 nil")
	ErrUnknownUnit   = errors.New("unit: 未知的单位模板 id")
)

// SpawnUnitByID 根据全局 [config.Tab.UnitConfigByID] 创建实体并挂载战斗组件：
// [component.Attributes]（由属性表行展开）、可选 [component.Health]（含 hp 行时）、
// [component.SkillSet]（ability 中的技能须在 [config.Tab.SkillConfigByID] 中存在）、
// 以及初始 Buff（[buff.AddBuff] 施法者与目标均为该实体）。
// 调用方须已向 [ecs.World] 注册相关组件类型（如 [component.RegisterCombatTypesWorld]）。
func SpawnUnitByID(w *ecs.World, unitID int32) (ecs.Entity, error) {
	if w == nil {
		return 0, ErrNilWorld
	}

	unitDesc := config.GetUnitConfigByID(unitID)
	if unitDesc == nil {
		return 0, fmt.Errorf("%w: %d", ErrUnknownUnit, unitID)
	}

	e := w.CreateEntity()

	attributes.InitFromConfig(w, e, unitDesc.Stats)

	for _, sid := range unitDesc.Ability {
		if !skill.AddSkill(w, e, sid) {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 初始技能 %d 挂载失败（单位模板 %d）", sid, unitID)
		}
	}

	for _, bid := range unitDesc.BuffDefIDs {
		if !buff.AddBuff(w, e, e, bid) {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 初始 Buff %d 挂载失败（单位模板 %d）", bid, unitID)
		}
	}

	return e, nil
}

// SpawnUnitByPlayer 使用玩家存档 [Player.Units] 中的单位数据生成实体：
// 属性来自 [PlayerUnit.Stats] 行内数值（不查属性表）；技能、Buff 与 [SpawnUnitByID] 相同路径。
// 若 [Player.Units] 含多条，选用 **最小 unitKey** 对应条目（稳定、可预期）；全空或全 nil 则返回 [ErrNoPlayerUnits]。
func SpawnUnitByPlayer(w *ecs.World, player *Player) (ecs.Entity, error) {
	if w == nil {
		return 0, ErrNilWorld
	}
	if player == nil {
		return 0, ErrNilPlayer
	}
	pu, unitKey, err := pickPlayerUnit(player)
	if err != nil {
		return 0, err
	}

	e := w.CreateEntity()
	attributes.InitFromStats(w, e, pu.Stats)

	for _, sid := range pu.Ability {
		if !skill.AddSkill(w, e, sid) {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 玩家单位 %d 技能 %d 挂载失败", unitKey, sid)
		}
	}
	for _, bid := range pu.BuffDefIDs {
		if !buff.AddBuff(w, e, e, bid) {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 玩家单位 %d Buff %d 挂载失败", unitKey, bid)
		}
	}

	return e, nil
}

func pickPlayerUnit(p *Player) (*PlayerUnit, uint32, error) {
	if p == nil || len(p.Units) == 0 {
		return nil, 0, ErrNoPlayerUnits
	}
	var bestK uint32
	var best *PlayerUnit
	first := true
	for k, u := range p.Units {
		if u == nil {
			continue
		}
		if first || k < bestK {
			first = false
			bestK = k
			best = u
		}
	}
	if best == nil {
		return nil, 0, ErrNoPlayerUnits
	}
	return best, bestK, nil
}
