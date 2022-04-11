package priest

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

const SpellIDMindBlast int32 = 25375

var MBCooldownID = core.NewCooldownID()
var MindBlastActionID = core.ActionID{SpellID: SpellIDMindBlast, CooldownID: MBCooldownID}

func (priest *Priest) registerMindBlastSpell(sim *core.Simulation) {
	baseCost := 450.0

	priest.MindBlast = priest.RegisterSpell(core.SpellConfig{
		ActionID:    MindBlastActionID,
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				Cost:     baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			Cooldown: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{

			BonusSpellHitRating: 0 +
				float64(priest.Talents.ShadowFocus)*2*core.SpellHitRatingPerHitChance +
				float64(priest.Talents.FocusedPower)*2*core.SpellHitRatingPerHitChance,

			BonusSpellCritRating: float64(priest.Talents.ShadowPower) * 3 * core.SpellCritRatingPerCritChance,

			DamageMultiplier: 1 *
				(1 + float64(priest.Talents.Darkness)*0.02) *
				core.TernaryFloat64(priest.Talents.Shadowform, 1.15, 1) *
				core.TernaryFloat64(ItemSetAbsolution.CharacterHasSetBonus(&priest.Character, 4), 1.1, 1),

			ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

			BaseDamage:     core.BaseDamageConfigMagic(711, 752, 0.429),
			OutcomeApplier: core.OutcomeFuncMagicHitAndCrit(priest.DefaultSpellCritMultiplier()),
		}),
	})
}
