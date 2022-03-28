package core

import (
	"fmt"
	"math"

	"github.com/wowsims/tbc/sim/core/stats"
)

// Callback for after a spell hits the target, before damage has been calculated.
// Use it to modify the spell damage or results.
type OnBeforeSpellHit func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellHitEffect)

// Callback for after a spell hits the target and after damage is calculated. Use it for proc effects
// or anything that comes from the final result of the spell.
type OnSpellHit func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect)

// OnPeriodicDamage is called when dots tick, after damage is calculated. Use it for proc effects
// or anything that comes from the final result of a tick.
type OnPeriodicDamage func(sim *Simulation, spellCast *SpellCast, spellEffect *SpellEffect, tickDamage float64)

// A Spell is a type of cast that can hit/miss using spell stats, and has a spell school.
type SpellCast struct {
	// Embedded Cast
	Cast

	// Results from the spell cast. Spell casts can have multiple effects (e.g.
	// Chain Lightning, Moonfire) so these are totals from all the effects.
	Hits               int32
	Misses             int32
	Crits              int32
	PartialResists_1_4 int32   // 1/4 of the spell was resisted
	PartialResists_2_4 int32   // 2/4 of the spell was resisted
	PartialResists_3_4 int32   // 3/4 of the spell was resisted
	TotalDamage        float64 // Damage done by this cast.
	TotalThreat        float64 // Threat generated by this cast.

	// Melee only stats
	Dodges  int32
	Glances int32
	Parries int32
	Blocks  int32
}

type SpellEffect struct {
	// Target of the spell.
	Target *Target

	// Bonus stats to be added to the spell.
	BonusSpellHitRating  float64
	BonusSpellPower      float64
	BonusSpellCritRating float64

	BonusHitRating        float64
	BonusAttackPower      float64
	BonusCritRating       float64
	BonusExpertiseRating  float64
	BonusArmorPenetration float64
	BonusWeaponDamage     float64

	BonusAttackPowerOnTarget float64

	// Additional multiplier that is always applied.
	DamageMultiplier float64

	// applies fixed % increases to damage at cast time.
	//  Only use multipliers that don't change for the lifetime of the sim.
	//  This should probably only be mutated in a template and not changed in auras.
	StaticDamageMultiplier float64

	// Multiplier for all threat generated by this effect.
	ThreatMultiplier float64

	// Adds a fixed amount of threat to this spell, before multipliers.
	FlatThreatBonus float64

	// Controls which effects can proc from this effect.
	ProcMask ProcMask

	// Causes the first roll for this hit to be copied from ActiveMeleeAbility.Effects[0].HitType.
	// This is only used by Shaman Stormstrike.
	ReuseMainHitRoll bool

	// Callbacks for providing additional custom behavior.
	OnSpellHit OnSpellHit

	// Results
	Outcome HitOutcome
	Damage  float64 // Damage done by this cast.
	Threat  float64

	// Certain damage multiplier, such as target debuffs and crit multipliers, do
	// not count towards the AOE cap. Store them here to they can be subtracted
	// later when calculating AOE cap.
	BeyondAOECapMultiplier float64
}

func (spellEffect *SpellEffect) Landed() bool {
	return spellEffect.Outcome.Matches(OutcomeLanded)
}

func (spellEffect *SpellEffect) TotalThreatMultiplier(spellCast *SpellCast) float64 {
	return spellEffect.ThreatMultiplier * spellCast.Character.PseudoStats.ThreatMultiplier
}

func (spellEffect *SpellEffect) calcThreat(spellCast *SpellCast) float64 {
	if spellEffect.Landed() {
		return (spellEffect.Damage + spellEffect.FlatThreatBonus) * spellEffect.TotalThreatMultiplier(spellCast)
	} else {
		return 0
	}
}

func (she *SpellHitEffect) beforeCalculations(sim *Simulation, spell *SimpleSpell) {
	se := &she.SpellEffect
	se.beforeCalculations(sim, spell, she)
}

func (spellEffect *SpellEffect) beforeCalculations(sim *Simulation, spell *SimpleSpell, she *SpellHitEffect) {
	spellEffect.BeyondAOECapMultiplier = 1
	multiplierBeforeTargetEffects := spellEffect.DamageMultiplier

	spell.Character.OnBeforeSpellHit(sim, &spell.SpellCast, she)
	spellEffect.Target.OnBeforeSpellHit(sim, &spell.SpellCast, she)

	spellEffect.BeyondAOECapMultiplier *= spellEffect.DamageMultiplier / multiplierBeforeTargetEffects

	spellEffect.determineOutcome(sim, spell, she)
}

func (hitEffect *SpellHitEffect) directCalculations(sim *Simulation, spell *SimpleSpell) {
	damage := hitEffect.calculateBaseDamage(sim, &spell.SpellCast)

	damage *= hitEffect.DamageMultiplier * hitEffect.StaticDamageMultiplier
	hitEffect.applyAttackerMultipliers(sim, &spell.SpellCast, false, &damage)
	hitEffect.applyTargetMultipliers(sim, &spell.SpellCast, false, &damage)
	hitEffect.applyResistances(sim, &spell.SpellCast, &damage)
	hitEffect.applyOutcome(sim, &spell.SpellCast, &damage)

	hitEffect.Damage = damage
}

func (hitEffect *SpellHitEffect) calculateBaseDamage(sim *Simulation, spellCast *SpellCast) float64 {
	character := spellCast.Character
	damage := 0.0

	// Weapon Damage Effects
	if hitEffect.WeaponInput.HasWeaponDamage() {
		var attackPower float64
		var bonusWeaponDamage float64
		if spellCast.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) {
			// all ranged attacks honor BonusAttackPowerOnTarget...
			attackPower = character.stats[stats.RangedAttackPower] + hitEffect.BonusAttackPower + hitEffect.BonusAttackPowerOnTarget
			bonusWeaponDamage = character.PseudoStats.BonusDamage + hitEffect.BonusWeaponDamage
		} else {
			attackPower = character.stats[stats.AttackPower] + hitEffect.BonusAttackPower
			bonusWeaponDamage = character.PseudoStats.BonusDamage + hitEffect.BonusWeaponDamage
		}

		if hitEffect.WeaponInput.CalculateDamage != nil {
			damage += hitEffect.WeaponInput.CalculateDamage(attackPower, bonusWeaponDamage)
		} else if hitEffect.WeaponInput.DamageMultiplier != 0 {
			// Bonus weapon damage applies after OH penalty: https://www.youtube.com/watch?v=bwCIU87hqTs
			// TODO not all weapon damage based attacks "scale" with +bonusWeaponDamage (e.g. Devastate, Shiv, Mutilate don't)
			// ... but for other's, BonusAttackPowerOnTarget only applies to weapon damage based attacks
			if hitEffect.WeaponInput.Normalized {
				if spellCast.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) {
					damage += character.AutoAttacks.Ranged.calculateNormalizedWeaponDamage(sim, attackPower) + bonusWeaponDamage
				} else if !hitEffect.WeaponInput.Offhand {
					damage += character.AutoAttacks.MH.calculateNormalizedWeaponDamage(sim, attackPower+hitEffect.BonusAttackPowerOnTarget) + bonusWeaponDamage
				} else {
					damage += character.AutoAttacks.OH.calculateNormalizedWeaponDamage(sim, attackPower+2*hitEffect.BonusAttackPowerOnTarget)*0.5 + bonusWeaponDamage
				}
			} else {
				if spellCast.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) {
					damage += character.AutoAttacks.Ranged.calculateWeaponDamage(sim, attackPower) + bonusWeaponDamage
				} else if !hitEffect.WeaponInput.Offhand {
					damage += character.AutoAttacks.MH.calculateWeaponDamage(sim, attackPower+hitEffect.BonusAttackPowerOnTarget) + bonusWeaponDamage
				} else {
					damage += character.AutoAttacks.OH.calculateWeaponDamage(sim, attackPower+2*hitEffect.BonusAttackPowerOnTarget)*0.5 + bonusWeaponDamage
				}
			}
			damage += hitEffect.WeaponInput.FlatDamageBonus
			damage *= hitEffect.WeaponInput.DamageMultiplier
		}

		if hitEffect.DirectInput.SpellCoefficient > 0 {
			bonus := (character.GetStat(stats.SpellPower) + character.GetStat(spellCast.SpellSchool.Stat())) * hitEffect.DirectInput.SpellCoefficient * hitEffect.WeaponInput.DamageMultiplier
			bonus += hitEffect.SpellEffect.BonusSpellPower * hitEffect.DirectInput.SpellCoefficient // does not get changed by weapon input multiplier
			damage += bonus
		}

		//if sim.Log != nil {
		//	character.Log(sim, "Melee damage calcs: AP=%0.1f, bonusWepdamage:%0.1f, damageMultiplier:%0.2f, staticMultiplier:%0.2f, result:%d, weapondamageCalc: %0.1f, critMultiplier: %0.3f, Target armor: %0.1f\n", attackPower, bonusWeaponDamage, hitEffect.DamageMultiplier, hitEffect.StaticDamageMultiplier, hitEffect.HitType, damage, spellCast.CritMultiplier, hitEffect.Target.currentArmor)
		//}
	}

	// Direct Damage Effects
	if hitEffect.DirectInput.MaxBaseDamage != 0 {
		baseDamage := hitEffect.DirectInput.MinBaseDamage + sim.RandomFloat("Base Damage Direct")*(hitEffect.DirectInput.MaxBaseDamage-hitEffect.DirectInput.MinBaseDamage)

		schoolBonus := 0.0
		// Use outcome roll to decide if it should use AP or spell school for bonus damage.
		isPhysical := spellCast.OutcomeRollCategory.Matches(OutcomeRollCategoryPhysical)
		if isPhysical {
			if spellCast.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) {
				schoolBonus = character.stats[stats.RangedAttackPower]
			} else if spellCast.SpellSchool == SpellSchoolPhysical {
				schoolBonus = character.stats[stats.AttackPower]
			}
			schoolBonus += hitEffect.BonusAttackPower
		} else {
			schoolBonus = character.GetStat(stats.SpellPower) + character.GetStat(spellCast.SpellSchool.Stat()) + hitEffect.SpellEffect.BonusSpellPower
		}
		damage += baseDamage + (schoolBonus * hitEffect.DirectInput.SpellCoefficient)
	}

	damage += hitEffect.DirectInput.FlatDamageBonus
	return damage
}

func (spellEffect *SpellEffect) determineOutcome(sim *Simulation, spell *SimpleSpell, she *SpellHitEffect) {
	if spell.OutcomeRollCategory == OutcomeRollCategoryNone || spell.SpellExtras.Matches(SpellExtrasAlwaysHits) {
		spellEffect.Outcome = OutcomeHit
		if spellEffect.critCheck(sim, &spell.SpellCast) {
			spellEffect.Outcome |= OutcomeCrit
		}
	} else if spellEffect.ReuseMainHitRoll { // TODO: can we remove this.
		spellEffect.Outcome = spell.Effects[0].Outcome
	} else if spell.OutcomeRollCategory.Matches(OutcomeRollCategoryMagic) {
		if spellEffect.hitCheck(sim, &spell.SpellCast) {
			spellEffect.Outcome = OutcomeHit
			if spellEffect.critCheck(sim, &spell.SpellCast) {
				spellEffect.Outcome |= OutcomeCrit
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
		}
	} else if spell.OutcomeRollCategory.Matches(OutcomeRollCategoryPhysical) {
		spellEffect.Outcome = spellEffect.WhiteHitTableResult(sim, spell)
		if spellEffect.Landed() && spellEffect.critCheck(sim, &spell.SpellCast) {
			spellEffect.Outcome = OutcomeCrit
		}
	}
}

// Computes an attack result using the white-hit table formula (single roll).
func (ahe *SpellEffect) WhiteHitTableResult(sim *Simulation, ability *SimpleSpell) HitOutcome {
	character := ability.Character

	roll := sim.RandomFloat("White Hit Table")

	// Miss
	missChance := ahe.Target.MissChance
	if character.AutoAttacks.IsDualWielding && ability.OutcomeRollCategory == OutcomeRollCategoryWhite {
		missChance += 0.19
	}
	hitBonus := ((character.stats[stats.MeleeHit] + ahe.BonusHitRating) / (MeleeHitRatingPerHitChance * 100)) - ahe.Target.HitSuppression
	if hitBonus > 0 {
		missChance = MaxFloat(0, missChance-hitBonus)
	}

	chance := missChance
	if roll < chance {
		return OutcomeMiss
	}

	if !ability.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) { // Ranged hits can't be dodged/glance, and are always 2-roll
		// Dodge
		if !ability.SpellExtras.Matches(SpellExtrasCannotBeDodged) {
			dodge := ahe.Target.Dodge
			expertisePercentage := MinFloat(math.Floor((character.stats[stats.Expertise]+ahe.BonusExpertiseRating)/(ExpertisePerQuarterPercentReduction))/400, dodge)
			chance += dodge - expertisePercentage
			if roll < chance {
				return OutcomeDodge
			}
		}

		// Parry (if in front)
		// If the target is a mob and defense minus weapon skill is 11 or more:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.6%

		// If the target is a mob and defense minus weapon skill is 10 or less:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.1%

		// Block (if in front)
		// If the target is a mob:
		// BlockChance = MIN(5%, 5% + (TargetLevel*5 - AttackerSkill) * 0.1%)
		// If we actually implement blocks, ranged hits can be blocked.

		// No need to crit/glance roll if we are not a white hit
		if ability.OutcomeRollCategory.Matches(OutcomeRollCategorySpecial | OutcomeRollCategoryRanged) {
			return OutcomeHit
		}

		// Glance
		chance += ahe.Target.Glance
		if roll < chance {
			return OutcomeGlance
		}

		// Crit
		critChance := ((character.stats[stats.MeleeCrit] + ahe.BonusCritRating) / (MeleeCritRatingPerCritChance * 100)) - ahe.Target.CritSuppression
		chance += critChance
		if roll < chance {
			return OutcomeCrit
		}
	}

	return OutcomeHit
}

// Calculates a hit check using the stats from this spell.
func (spellEffect *SpellEffect) hitCheck(sim *Simulation, spellCast *SpellCast) bool {
	hit := 0.83 + (spellCast.Character.GetStat(stats.SpellHit)+spellEffect.BonusSpellHitRating)/(SpellHitRatingPerHitChance*100)
	hit = MinFloat(hit, 0.99) // can't get away from the 1% miss

	return sim.RandomFloat("Magical Hit Roll") < hit
}

// Calculates a crit check using the stats from this spell.
func (spellEffect *SpellEffect) critCheck(sim *Simulation, spellCast *SpellCast) bool {
	switch spellCast.CritRollCategory {
	case CritRollCategoryMagical:
		critChance := (spellCast.Character.GetStat(stats.SpellCrit) + spellCast.BonusCritRating + spellEffect.BonusSpellCritRating) / (SpellCritRatingPerCritChance * 100)
		return sim.RandomFloat("Magical Crit Roll") < critChance
	case CritRollCategoryPhysical:
		critChance := (spellCast.Character.GetStat(stats.MeleeCrit)+spellCast.BonusCritRating+spellEffect.BonusCritRating)/(MeleeCritRatingPerCritChance*100) - spellEffect.Target.CritSuppression
		return sim.RandomFloat("Physical Crit Roll") < critChance
	default:
		return false
	}
}

func (spellEffect *SpellEffect) triggerSpellProcs(sim *Simulation, spell *SimpleSpell) {
	if spellEffect.OnSpellHit != nil {
		spellEffect.OnSpellHit(sim, &spell.SpellCast, spellEffect)
	}
	spell.Character.OnSpellHit(sim, &spell.SpellCast, spellEffect)
	spellEffect.Target.OnSpellHit(sim, &spell.SpellCast, spellEffect)
}

func (spellEffect *SpellEffect) afterCalculations(sim *Simulation, spell *SimpleSpell) {
	if sim.Log != nil && !spell.SpellExtras.Matches(SpellExtrasAlwaysHits) {
		spell.Character.Log(sim, "%s %s. (Threat: %0.3f)", spell.ActionID, spellEffect, spellEffect.calcThreat(&spell.SpellCast))
	}

	spellEffect.triggerSpellProcs(sim, spell)
}

func (spellEffect *SpellEffect) applyResultsToCast(spellCast *SpellCast) {
	if spellEffect.Outcome.Matches(OutcomeHit) {
		spellCast.Hits++
	}
	if spellEffect.Outcome.Matches(OutcomeGlance) {
		spellCast.Glances++
	}
	if spellEffect.Outcome.Matches(OutcomeCrit) {
		spellCast.Crits++
	}
	if spellEffect.Outcome.Matches(OutcomeBlock) {
		spellCast.Blocks++
	}

	if spellEffect.Landed() {
		if spellEffect.Outcome.Matches(OutcomePartial1_4) {
			spellCast.PartialResists_1_4++
		} else if spellEffect.Outcome.Matches(OutcomePartial2_4) {
			spellCast.PartialResists_2_4++
		} else if spellEffect.Outcome.Matches(OutcomePartial3_4) {
			spellCast.PartialResists_3_4++
		}
	} else {
		if spellEffect.Outcome == OutcomeMiss {
			spellCast.Misses++
		} else if spellEffect.Outcome == OutcomeDodge {
			spellCast.Dodges++
		} else if spellEffect.Outcome == OutcomeParry {
			spellCast.Parries++
		}
	}

	spellCast.TotalDamage += spellEffect.Damage
	spellCast.TotalThreat += spellEffect.calcThreat(spellCast)
}

func (spellEffect *SpellEffect) String() string {
	outcomeStr := spellEffect.Outcome.String()
	if !spellEffect.Landed() {
		return outcomeStr
	}
	return fmt.Sprintf("%s for %0.3f damage", outcomeStr, spellEffect.Damage)
}

func (hitEffect *SpellHitEffect) applyAttackerMultipliers(sim *Simulation, spellCast *SpellCast, isPeriodic bool, damage *float64) {
	attacker := spellCast.Character

	*damage *= attacker.PseudoStats.DamageDealtMultiplier
	if spellCast.SpellSchool.Matches(SpellSchoolPhysical) {
		*damage *= attacker.PseudoStats.PhysicalDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolArcane) {
		*damage *= attacker.PseudoStats.ArcaneDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolFire) {
		*damage *= attacker.PseudoStats.FireDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolFrost) {
		*damage *= attacker.PseudoStats.FrostDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolHoly) {
		*damage *= attacker.PseudoStats.HolyDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolNature) {
		*damage *= attacker.PseudoStats.NatureDamageDealtMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolShadow) {
		*damage *= attacker.PseudoStats.ShadowDamageDealtMultiplier
	}
}

func (hitEffect *SpellHitEffect) applyTargetMultipliers(sim *Simulation, spellCast *SpellCast, isPeriodic bool, damage *float64) {
	target := hitEffect.Target

	*damage *= target.PseudoStats.DamageTakenMultiplier
	if spellCast.SpellSchool.Matches(SpellSchoolPhysical) {
		*damage *= target.PseudoStats.PhysicalDamageTakenMultiplier
		if isPeriodic {
			*damage *= target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
		}
	} else if spellCast.SpellSchool.Matches(SpellSchoolArcane) {
		*damage *= target.PseudoStats.ArcaneDamageTakenMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolFire) {
		*damage *= target.PseudoStats.FireDamageTakenMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolFrost) {
		*damage *= target.PseudoStats.FrostDamageTakenMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolHoly) {
		*damage *= target.PseudoStats.HolyDamageTakenMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolNature) {
		*damage *= target.PseudoStats.NatureDamageTakenMultiplier
	} else if spellCast.SpellSchool.Matches(SpellSchoolShadow) {
		*damage *= target.PseudoStats.ShadowDamageTakenMultiplier
	}
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (hitEffect *SpellHitEffect) applyResistances(sim *Simulation, spellCast *SpellCast, damage *float64) {
	if spellCast.SpellExtras.Matches(SpellExtrasIgnoreResists) {
		return
	}

	if spellCast.SpellSchool.Matches(SpellSchoolPhysical) {
		// Physical resistance (armor).
		*damage *= 1 - hitEffect.Target.ArmorDamageReduction(spellCast.Character.stats[stats.ArmorPenetration]+hitEffect.BonusArmorPenetration)
	} else if !spellCast.SpellExtras.Matches(SpellExtrasBinary) {
		// Magical resistance.
		// https://royalgiraffe.github.io/resist-guide

		resistanceRoll := sim.RandomFloat("Partial Resist")
		if resistanceRoll > 0.18 { // 13% chance for 25% resist, 4% for 50%, 1% for 75%
			// No partial resist.
		} else if resistanceRoll > 0.05 {
			hitEffect.SpellEffect.Outcome |= OutcomePartial1_4
			*damage *= 0.75
		} else if resistanceRoll > 0.01 {
			hitEffect.SpellEffect.Outcome |= OutcomePartial2_4
			*damage *= 0.5
		} else {
			hitEffect.SpellEffect.Outcome |= OutcomePartial3_4
			*damage *= 0.25
		}
	}
}

func (hitEffect *SpellHitEffect) applyOutcome(sim *Simulation, spellCast *SpellCast, damage *float64) {
	if !hitEffect.Landed() {
		*damage = 0
	} else if hitEffect.Outcome.Matches(OutcomeCrit) {
		*damage *= spellCast.CritMultiplier
	} else if hitEffect.Outcome == OutcomeGlance {
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. 73)
		*damage *= 0.75
	}
}
