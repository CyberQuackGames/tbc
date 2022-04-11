package mage

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

// The numbers in this file are VERY rough approximations based on logs.

var SummonWaterElementalCooldownID = core.NewCooldownID()

func (mage *Mage) registerSummonWaterElementalCD() {
	if !mage.Talents.SummonWaterElemental {
		return
	}

	manaCost := 0.0
	actionID := core.ActionID{SpellID: 31687, CooldownID: SummonWaterElementalCooldownID}

	mage.AddMajorCooldown(core.MajorCooldown{
		ActionID:   actionID,
		CooldownID: SummonWaterElementalCooldownID,
		Cooldown:   time.Minute * 3,
		UsesGCD:    true,
		Priority:   core.CooldownPriorityDrums + 1, // Always prefer to cast before drums or lust so the ele gets their benefits.
		Type:       core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			if mage.waterElemental.IsEnabled() {
				return false
			}
			if character.CurrentMana() < manaCost {
				return false
			}
			return true
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
		ActivationFactory: func(sim *core.Simulation) core.CooldownActivation {
			baseManaCost := mage.BaseMana() * 0.16
			castTemplate := core.SimpleCast{
				Cast: core.Cast{
					ActionID:  actionID,
					Character: mage.GetCharacter(),
					BaseCost: core.ResourceCost{
						Type:  stats.Mana,
						Value: baseManaCost,
					},
					Cost: core.ResourceCost{
						Type:  stats.Mana,
						Value: baseManaCost,
					},
					GCD:      core.GCDDefault,
					Cooldown: time.Minute * 3,
					OnCastComplete: func(sim *core.Simulation, cast *core.Cast) {
						mage.waterElemental.EnableWithTimeout(sim, mage.waterElemental, time.Second*45)

						// All MCDs that use the GCD and have a non-zero cast time must call this.
						mage.UpdateMajorCooldowns()
					},
				},
			}
			castTemplate.Cost.Value -= castTemplate.BaseCost.Value * float64(mage.Talents.FrostChanneling) * 0.05
			castTemplate.Cost.Value *= 1 - float64(mage.Talents.ElementalPrecision)*0.01
			manaCost = castTemplate.Cost.Value

			return func(sim *core.Simulation, character *core.Character) {
				cast := castTemplate
				cast.Init(sim)
				cast.StartCast(sim)
			}
		},
	})
}

type WaterElemental struct {
	core.Pet

	// Water Ele almost never just stands still and spams like we want, it sometimes
	// does its own thing. This controls how much it does that.
	disobeyChance float64

	Waterbolt *core.Spell
}

func (mage *Mage) NewWaterElemental(disobeyChance float64) *WaterElemental {
	waterElemental := &WaterElemental{
		Pet: core.NewPet(
			"Water Elemental",
			&mage.Character,
			waterElementalBaseStats,
			waterElementalStatInheritance,
			false,
		),
		disobeyChance: disobeyChance,
	}
	waterElemental.EnableManaBar()

	mage.AddPet(waterElemental)

	return waterElemental
}

func (we *WaterElemental) GetPet() *core.Pet {
	return &we.Pet
}

func (we *WaterElemental) Init(sim *core.Simulation) {
	we.registerWaterboltSpell(sim)
}

func (we *WaterElemental) Reset(sim *core.Simulation) {
}

func (we *WaterElemental) OnGCDReady(sim *core.Simulation) {
	spell := we.Waterbolt

	if sim.RandomFloat("Water Elemental Disobey") < we.disobeyChance {
		// Water ele has decided not to cooperate, so just wait for the cast time
		// instead of casting.
		we.WaitUntil(sim, sim.CurrentTime+spell.Template.CastTime)
		return
	}

	if success := spell.Cast(sim, sim.GetPrimaryTarget()); !success {
		// If water ele has gone OOM then there won't be enough time left for meaningful
		// regen to occur before the ele expires. So just murder itself.
		we.Disable(sim)
	}
}

// These numbers are just rough guesses based on looking at some logs.
var waterElementalBaseStats = stats.Stats{
	stats.Intellect:  100,
	stats.SpellPower: 300,
	stats.Mana:       2000,
	stats.SpellHit:   3 * core.SpellHitRatingPerHitChance,
	stats.SpellCrit:  8 * core.SpellCritRatingPerCritChance,
}

var waterElementalStatInheritance = func(ownerStats stats.Stats) stats.Stats {
	// These numbers are just rough guesses based on looking at some logs.
	return ownerStats.DotProduct(stats.Stats{
		// Computed based on my lvl 65 mage, need to ask someone with a 70 to check these
		stats.Stamina:   0.2238,
		stats.Intellect: 0.01,

		stats.SpellPower:      0.333,
		stats.FrostSpellPower: 0.333,
		stats.SpellHit:        0.01,
		stats.SpellCrit:       0.01,
	})
}

const SpellIDWaterbolt int32 = 31707

var WaterboltActionID = core.ActionID{SpellID: SpellIDWaterbolt}

func (we *WaterElemental) registerWaterboltSpell(sim *core.Simulation) {
	baseCost := we.BaseMana() * 0.1

	we.Waterbolt = we.RegisterSpell(core.SpellConfig{
		ActionID:    WaterboltActionID,
		SpellSchool: core.SpellSchoolFrost,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.NewCast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Second * 3,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagic(256, 328, 3.0/3.5),
			OutcomeApplier:   core.OutcomeFuncMagicHitAndCrit(we.DefaultSpellCritMultiplier()),
		}),
	})
}
