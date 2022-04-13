package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
)

var EvocationCooldownID = core.NewCooldownID()

func (mage *Mage) registerEvocationCD() {
	cooldown := time.Minute * 8
	manaThreshold := 0.0
	actionID := core.ActionID{SpellID: 12051, CooldownID: EvocationCooldownID}

	maxTicks := int32(4)
	if ItemSetTempestRegalia.CharacterHasSetBonus(&mage.Character, 2) {
		maxTicks++
	}

	numTicks := core.MaxInt32(0, core.MinInt32(maxTicks, mage.Options.EvocationTicks))
	if numTicks == 0 {
		numTicks = maxTicks
	}

	channelTime := time.Duration(numTicks) * time.Second * 2
	manaPerTick := 0.0

	evocationSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			Cooldown: cooldown,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Target, spell *core.Spell) {
			period := spell.CurCast.ChannelTime / time.Duration(numTicks)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   period,
				NumTicks: int(numTicks),
				OnAction: func(sim *core.Simulation) {
					mage.AddMana(sim, manaPerTick, actionID, true)
				},
			})

			// All MCDs that use the GCD and have a non-zero cast time must call this.
			mage.UpdateMajorCooldowns()
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		ActionID:   actionID,
		CooldownID: EvocationCooldownID,
		Cooldown:   cooldown,
		UsesGCD:    true,
		Type:       core.CooldownTypeMana,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if character.HasActiveAuraWithTag(core.InnervateAuraTag) || character.HasActiveAuraWithTag(core.ManaTideTotemAuraTag) {
				return false
			}

			curMana := character.CurrentMana()
			if curMana > manaThreshold {
				return false
			}

			if character.HasActiveAuraWithTag(core.BloodlustAuraTag) && curMana > manaThreshold/2 {
				return false
			}

			if mage.isBlastSpamming {
				return false
			}

			return true
		},
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			manaPerTick = mage.MaxMana() * 0.15
			manaThreshold = mage.MaxMana() * 0.2

			return func(sim *core.Simulation, character *core.Character) {
				evocationSpell.Cast(sim, nil)
			}
		},
	})
}
