package core

import (
	"math"

	"github.com/wowsims/tbc/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64)

func OutcomeFuncAlwaysHit() OutcomeApplier {
	return func(_ *Simulation, spell *Spell, spellEffect *SpellEffect, _ *float64) {
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

// A tick always hits, but we don't count them as hits in the metrics.
func OutcomeFuncTick() OutcomeApplier {
	return func(_ *Simulation, _ *Spell, spellEffect *SpellEffect, _ *float64) {
		spellEffect.Outcome = OutcomeHit
	}
}

func OutcomeFuncMagicHitAndCrit(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		if spellEffect.magicHitCheck(sim, spell) {
			if spellEffect.magicCritCheck(sim, spell) {
				spellEffect.Outcome = OutcomeCrit
				spell.Crits++
				*damage *= critMultiplier
			} else {
				spellEffect.Outcome = OutcomeHit
				spell.Hits++
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
		}
	}
}

func OutcomeFuncMagicHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		if spellEffect.magicHitCheck(sim, spell) {
			spellEffect.Outcome = OutcomeHit
			spell.Hits++
		} else {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
		}
	}
}

func OutcomeFuncMeleeWhite(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		character := spell.Character
		roll := sim.RandomFloat("White Hit Table")

		// Miss
		missChance := spellEffect.Target.MissChance - spellEffect.PhysicalHitChance(character)
		if character.AutoAttacks.IsDualWielding {
			missChance += 0.19
		}

		chance := MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
			return
		}

		// Dodge
		chance += MaxFloat(0, spellEffect.Target.Dodge-spellEffect.ExpertisePercentage(character))
		if roll < chance {
			spellEffect.Outcome = OutcomeDodge
			spell.Dodges++
			*damage = 0
			return
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

		// Glance
		chance += spellEffect.Target.Glance
		if roll < chance {
			spellEffect.Outcome = OutcomeGlance
			spell.Glances++
			// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. 73)
			*damage *= 0.75
			return
		}

		// Crit
		chance += spellEffect.PhysicalCritChance(character, spell)
		if roll < chance {
			spellEffect.Outcome = OutcomeCrit
			spell.Crits++
			*damage *= critMultiplier
			return
		}

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func OutcomeFuncMeleeSpecialHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		character := spell.Character
		roll := sim.RandomFloat("White Hit Table")

		// Miss
		missChance := spellEffect.Target.MissChance - spellEffect.PhysicalHitChance(character)
		chance := MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
			return
		}

		// Dodge
		if !spell.SpellExtras.Matches(SpellExtrasCannotBeDodged) {
			chance += MaxFloat(0, spellEffect.Target.Dodge-spellEffect.ExpertisePercentage(character))
			if roll < chance {
				spellEffect.Outcome = OutcomeDodge
				spell.Dodges++
				*damage = 0
				return
			}
		}

		// Parry (if in front)
		// If the target is a mob and defense minus weapon skill is 11 or more:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.6%
		// If the target is a mob and defense minus weapon skill is 10 or less:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.1%

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func OutcomeFuncMeleeSpecialHitAndCrit(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		character := spell.Character
		roll := sim.RandomFloat("White Hit Table")

		// Miss
		missChance := spellEffect.Target.MissChance - spellEffect.PhysicalHitChance(character)
		chance := MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
			return
		}

		// Dodge
		if !spell.SpellExtras.Matches(SpellExtrasCannotBeDodged) {
			chance += MaxFloat(0, spellEffect.Target.Dodge-spellEffect.ExpertisePercentage(character))
			if roll < chance {
				spellEffect.Outcome = OutcomeDodge
				spell.Dodges++
				*damage = 0
				return
			}
		}

		// Parry (if in front)
		// If the target is a mob and defense minus weapon skill is 11 or more:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.6%
		// If the target is a mob and defense minus weapon skill is 10 or less:
		// ParryChance = 5% + (TargetLevel*5 - AttackerSkill) * 0.1%

		// Block (if in front). Note that critical blocks are allowed for 2-roll hits.
		// If the target is a mob:
		// BlockChance = MIN(5%, 5% + (TargetLevel*5 - AttackerSkill) * 0.1%)
		// If we actually implement blocks, ranged hits can be blocked.

		// Crit (separate roll)
		if spellEffect.physicalCritRoll(sim, spell) {
			spellEffect.Outcome = OutcomeCrit
			spell.Crits++
			*damage *= critMultiplier
			return
		}

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func OutcomeFuncMeleeSpecialCritOnly(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		if spellEffect.physicalCritRoll(sim, spell) {
			spellEffect.Outcome = OutcomeCrit
			spell.Crits++
			*damage *= critMultiplier
			return
		}

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func OutcomeFuncRangedHit() OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		character := spell.Character
		roll := sim.RandomFloat("White Hit Table")

		// Miss
		missChance := spellEffect.Target.MissChance - spellEffect.PhysicalHitChance(character)
		chance := MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
			return
		}

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func OutcomeFuncRangedHitAndCrit(critMultiplier float64) OutcomeApplier {
	return func(sim *Simulation, spell *Spell, spellEffect *SpellEffect, damage *float64) {
		character := spell.Character
		roll := sim.RandomFloat("White Hit Table")

		// Miss
		missChance := spellEffect.Target.MissChance - spellEffect.PhysicalHitChance(character)
		chance := MaxFloat(0, missChance)
		if roll < chance {
			spellEffect.Outcome = OutcomeMiss
			spell.Misses++
			*damage = 0
			return
		}

		// Block (if in front). Note that critical blocks are allowed for 2-roll hits.
		// If the target is a mob:
		// BlockChance = MIN(5%, 5% + (TargetLevel*5 - AttackerSkill) * 0.1%)
		// If we actually implement blocks, ranged hits can be blocked.

		// Crit (separate roll)
		if spellEffect.physicalCritRoll(sim, spell) {
			spellEffect.Outcome = OutcomeCrit
			spell.Crits++
			*damage *= critMultiplier
			return
		}

		// Hit
		spellEffect.Outcome = OutcomeHit
		spell.Hits++
	}
}

func (spellEffect *SpellEffect) determineOutcome(sim *Simulation, spell *Spell, isPeriodic bool) {
	if isPeriodic {
		if spellEffect.DotInput.TicksCanMissAndCrit {
			if spellEffect.magicHitCheck(sim, spell) {
				spellEffect.Outcome = OutcomeHit
				if spellEffect.critCheck(sim, spell) {
					spellEffect.Outcome = OutcomeCrit
				}
			} else {
				spellEffect.Outcome = OutcomeMiss
			}
		} else {
			spellEffect.Outcome = OutcomeHit
		}
		return
	}

	if spellEffect.OutcomeRollCategory == OutcomeRollCategoryNone || spell.SpellExtras.Matches(SpellExtrasAlwaysHits) {
		spellEffect.Outcome = OutcomeHit
		if spellEffect.critCheck(sim, spell) {
			spellEffect.Outcome = OutcomeCrit
		}
	} else if spellEffect.OutcomeRollCategory.Matches(OutcomeRollCategoryMagic) {
		if spellEffect.magicHitCheck(sim, spell) {
			spellEffect.Outcome = OutcomeHit
			if spellEffect.critCheck(sim, spell) {
				spellEffect.Outcome = OutcomeCrit
			}
		} else {
			spellEffect.Outcome = OutcomeMiss
		}
	} else if spellEffect.OutcomeRollCategory.Matches(OutcomeRollCategoryPhysical) {
		spellEffect.Outcome = spellEffect.WhiteHitTableResult(sim, spell)
		if spellEffect.Landed() && spellEffect.critCheck(sim, spell) {
			spellEffect.Outcome = OutcomeCrit
		}
	}
}

// Computes an attack result using the white-hit table formula (single roll).
func (ahe *SpellEffect) WhiteHitTableResult(sim *Simulation, spell *Spell) HitOutcome {
	character := spell.Character

	roll := sim.RandomFloat("White Hit Table")

	// Miss
	missChance := ahe.Target.MissChance - ahe.PhysicalHitChance(character)
	if character.AutoAttacks.IsDualWielding && ahe.OutcomeRollCategory == OutcomeRollCategoryWhite {
		missChance += 0.19
	}
	missChance = MaxFloat(0, missChance)

	chance := missChance
	if roll < chance {
		return OutcomeMiss
	}

	if !ahe.OutcomeRollCategory.Matches(OutcomeRollCategoryRanged) { // Ranged hits can't be dodged/glance, and are always 2-roll
		// Dodge
		if !spell.SpellExtras.Matches(SpellExtrasCannotBeDodged) {
			dodge := ahe.Target.Dodge

			expertiseRating := character.stats[stats.Expertise]
			if ahe.ProcMask.Matches(ProcMaskMeleeMH) {
				expertiseRating += character.PseudoStats.BonusMHExpertiseRating
			} else if ahe.ProcMask.Matches(ProcMaskMeleeOH) {
				expertiseRating += character.PseudoStats.BonusOHExpertiseRating
			}
			expertisePercentage := MinFloat(math.Floor(expertiseRating/ExpertisePerQuarterPercentReduction)/400, dodge)

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
		if ahe.OutcomeRollCategory.Matches(OutcomeRollCategorySpecial | OutcomeRollCategoryRanged) {
			return OutcomeHit
		}

		// Glance
		chance += ahe.Target.Glance
		if roll < chance {
			return OutcomeGlance
		}

		// Crit
		chance += ahe.PhysicalCritChance(character, spell)
		if roll < chance {
			return OutcomeCrit
		}
	}

	return OutcomeHit
}

// Calculates a hit check using the stats from this spell.
func (spellEffect *SpellEffect) magicHitCheck(sim *Simulation, spell *Spell) bool {
	hit := 0.83 + (spell.Character.GetStat(stats.SpellHit)+spellEffect.BonusSpellHitRating)/(SpellHitRatingPerHitChance*100)
	hit = MinFloat(hit, 0.99) // can't get away from the 1% miss

	return sim.RandomFloat("Magical Hit Roll") < hit
}

func (spellEffect *SpellEffect) magicCritCheck(sim *Simulation, spell *Spell) bool {
	critChance := spellEffect.SpellCritChance(spell.Character, spell)
	return sim.RandomFloat("Magical Crit Roll") < critChance
}

func (spellEffect *SpellEffect) physicalCritRoll(sim *Simulation, spell *Spell) bool {
	return sim.RandomFloat("Physical Crit Roll") < spellEffect.PhysicalCritChance(spell.Character, spell)
}

// Calculates a crit check using the stats from this spell.
func (spellEffect *SpellEffect) critCheck(sim *Simulation, spell *Spell) bool {
	switch spellEffect.CritRollCategory {
	case CritRollCategoryMagical:
		return spellEffect.magicCritCheck(sim, spell)
	case CritRollCategoryPhysical:
		return spellEffect.physicalCritRoll(sim, spell)
	default:
		return false
	}
}
