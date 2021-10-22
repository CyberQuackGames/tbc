package common

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func init() {
	core.AddActiveItem(core.ItemIDTheLightningCapacitor, core.ActiveItem{BuffUp: ActivateTLC, ActivateCD: core.NeverExpires, SharedID: core.MagicIDAtkTrinket})
}

func ActivateTLC(sim *core.Simulation, agent core.Agent) {
	character := agent.GetCharacter()
	charges := 0

	const icdDur = time.Millisecond * 2500
	icd := core.NewICD()

	character.AddAura(sim, core.Aura{
		ID:      core.MagicIDTLC,
		Expires: core.NeverExpires,
		OnSpellHit: func(sim *core.Simulation, cast core.DirectCastAction, result *core.DirectCastDamageResult) {
			if icd.IsOnCD(sim) {
				return
			}

			if !result.Crit {
				return
			}

			charges++
			if charges >= 3 {
				icd = core.InternalCD(sim.CurrentTime + icdDur)
				charges = 0
				castAction := NewLightningCapacitorCast(sim, character)
				castAction.Act(sim)
			}
		},
	})
}

type LightningCapacitorCast struct {
	character *core.Character
}

func (lcc LightningCapacitorCast) GetActionID() core.ActionID {
	return core.ActionID{
		ItemID: core.ItemIDTheLightningCapacitor,
	}
}

func (lcc LightningCapacitorCast) GetName() string {
	return "Lightning Capacitor"
}

func (lcc LightningCapacitorCast) GetTag() int32 {
	return 0
}

func (lcc LightningCapacitorCast) GetCharacter() *core.Character {
	return lcc.character
}

func (lcc LightningCapacitorCast) GetBaseManaCost() float64 {
	return 0
}

func (lcc LightningCapacitorCast) GetSpellSchool() stats.Stat {
	return stats.NatureSpellPower
}

func (lcc LightningCapacitorCast) GetCooldown() time.Duration {
	return 0
}

func (lcc LightningCapacitorCast) GetCastInput(sim *core.Simulation, cast core.DirectCastAction) core.DirectCastInput {
	return core.DirectCastInput{
		IgnoreCooldowns: true,
		IgnoreManaCost: true,
		CritMultiplier: 1.5,
	}
}

func (lcc LightningCapacitorCast) GetHitInputs(sim *core.Simulation, cast core.DirectCastAction) []core.DirectCastDamageInput{
	hitInput := core.DirectCastDamageInput{
		MinBaseDamage: 694,
		MaxBaseDamage: 807,
		DamageMultiplier: 1,
	}

	return []core.DirectCastDamageInput{hitInput}
}

func (lcc LightningCapacitorCast) OnCastComplete(sim *core.Simulation, cast core.DirectCastAction) {
}
func (lcc LightningCapacitorCast) OnSpellHit(sim *core.Simulation, cast core.DirectCastAction, result *core.DirectCastDamageResult) {
}
func (lcc LightningCapacitorCast) OnSpellMiss(sim *core.Simulation, cast core.DirectCastAction) {
}

func NewLightningCapacitorCast(sim *core.Simulation, character *core.Character) core.DirectCastAction {
	return core.NewDirectCastAction(sim, LightningCapacitorCast{character: character})
}