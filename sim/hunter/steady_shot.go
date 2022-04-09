package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

var SteadyShotActionID = core.ActionID{SpellID: 34120}

func (hunter *Hunter) registerSteadyShotSpell(sim *core.Simulation) {
	ama := core.SimpleSpell{
		SpellCast: core.SpellCast{
			Cast: core.Cast{
				ActionID:  SteadyShotActionID,
				Character: hunter.GetCharacter(),
				BaseCost: core.ResourceCost{
					Type:  stats.Mana,
					Value: 110,
				},
				Cost: core.ResourceCost{
					Type:  stats.Mana,
					Value: 110,
				},
				// Cast time is affected by ranged attack speed so set it later.
				//CastTime:     time.Millisecond * 1500,
				GCD:         core.GCDDefault + hunter.latency,
				IgnoreHaste: true, // Hunter GCD is locked at 1.5s
				SpellSchool: core.SpellSchoolPhysical,
				SpellExtras: core.SpellExtrasMeleeMetrics,
			},
		},
	}
	ama.Cost.Value *= 1 - 0.02*float64(hunter.Talents.Efficiency)

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		Template: ama,
		ModifyCast: func(sim *core.Simulation, target *core.Target, instance *core.SimpleSpell) {
			instance.CastTime = hunter.SteadyShotCastTime()
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskRangedSpecial,

			BonusCritRating:  core.TernaryFloat64(ItemSetRiftStalker.CharacterHasSetBonus(&hunter.Character, 4), 5*core.MeleeCritRatingPerCritChance, 0),
			DamageMultiplier: 1 * core.TernaryFloat64(ItemSetGronnstalker.CharacterHasSetBonus(&hunter.Character, 4), 1.1, 1),
			ThreatMultiplier: 1,

			BaseDamage: hunter.talonOfAlarDamageMod(core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (hitEffect.RangedAttackPower(spell.Character)+hitEffect.RangedAttackPowerOnTarget())*0.2 +
						hunter.AutoAttacks.Ranged.BaseDamage(sim)*2.8/hunter.AutoAttacks.Ranged.SwingSpeed +
						150
				},
				TargetSpellCoefficient: 1,
			}),
			OutcomeApplier: core.OutcomeFuncRangedHitAndCrit(hunter.critMultiplier(true, sim.GetPrimaryTarget())),

			OnSpellHit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				hunter.killCommandBlocked = false
				hunter.TryKillCommand(sim, spellEffect.Target)
				hunter.rotation(sim, false)
			},
		}),
	})
}

func (hunter *Hunter) SteadyShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*1500)/hunter.RangedSwingSpeed()) + hunter.latency
}
