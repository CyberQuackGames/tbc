package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/wowsims/tbc/sim/core/stats"
)

const SpellCritRatingPerCritChance = 22.08

// Input needed to start casting a spell.
type DirectCastInput struct {
	ManaCost float64

	CastTime time.Duration

	// How much to multiply damage by, if this cast crits.
	CritMultiplier float64

	// If true, will force the cast to crit (if it doesnt miss).
	GuaranteedCrit bool
}

// Input needed to calculate the damage of a spell.
type DirectCastDamageInput struct {
	MinBaseDamage float64
	MaxBaseDamage float64

	// Increase in damage per point of spell power.
	SpellCoefficient float64

	// Additional multiplier that is always applied.
	DamageMultiplier float64

	BonusSpellPower float64
	BonusHit        float64 // Direct % bonus... 0.1 == 10%
	BonusCrit       float64 // Direct % bonus... 0.1 == 10%
}

type DirectCastDamageResult struct {
	Hit  bool // True = hit, False = resisted
	Crit bool // Whether this cast was a critical strike.

	PartialResist_1_4 bool // 1/4 of the spell was resisted
	PartialResist_2_4 bool // 2/4 of the spell was resisted
	PartialResist_3_4 bool // 3/4 of the spell was resisted

	Damage float64 // Damage done by this cast.
}

func (result *DirectCastDamageResult) String() string {
	if !result.Hit {
		return "Miss"
	}

	var sb strings.Builder
	sb.WriteString("")

	if result.PartialResist_1_4 {
		sb.WriteString("25% Resist ")
	} else if result.PartialResist_2_4 {
		sb.WriteString("50% Resist ")
	} else if result.PartialResist_3_4 {
		sb.WriteString("75% Resist ")
	}

	if result.Crit {
		sb.WriteString("Crit")
	} else {
		sb.WriteString("Hit")
	}

	fmt.Fprintf(&sb, " for %0.2f damage", result.Damage)
	return sb.String()
}

// Interface for direct cast spells to implement.
type DirectCastImpl interface {
	// Pass-through AgentAction methods
	GetActionID() ActionID
	GetName() string
	GetTag() int32
	GetAgent() Agent

	// This is needed because a lot of effects that 'reduce mana cost by X%' are
	// calculated from the base mana cost.
	GetBaseManaCost() float64

	// I.e. for nature spells, return stats.NatureSpellPower
	GetSpellSchool() stats.Stat

	GetCooldown() time.Duration

	GetCastInput(sim *Simulation, cast DirectCastAction) DirectCastInput
	GetHitInputs(sim *Simulation, cast DirectCastAction) []DirectCastDamageInput

	// Lifecycle callbacks for additional custom effects.
	OnCastComplete(sim *Simulation, cast DirectCastAction)
	OnSpellHit(sim *Simulation, cast DirectCastAction, result *DirectCastDamageResult)
	OnSpellMiss(sim *Simulation, cast DirectCastAction)
}

type DirectCastAction struct {
	DirectCastImpl

	// The inputs to the cast. Auras are given a reference to these to modify them
	// before the spell begins casting.
	castInput DirectCastInput
}

func (action DirectCastAction) GetManaCost() float64 {
	return action.castInput.ManaCost
}

func (action DirectCastAction) GetDuration() time.Duration {
	return action.castInput.CastTime
}

func (action DirectCastAction) Act(sim *Simulation) {
	character := action.GetAgent().GetCharacter()

	if sim.Log != nil {
		sim.Log("(%d) Casting %s (Current Mana = %0.0f, Mana Cost = %0.0f, Cast Time = %s\n",
				character.ID, action.GetName(), character.Stats[stats.Mana], action.castInput.ManaCost, action.castInput.CastTime)
	}

	if action.castInput.ManaCost > 0 {
		//fmt.Printf("Subtracting mana: %0.0f", action.castInput.ManaCost)
		character.Stats[stats.Mana] -= action.castInput.ManaCost
	}

	action.OnCastComplete(sim, action)
	for _, id := range character.ActiveAuraIDs {
		if character.Auras[id].OnCastComplete != nil {
			character.Auras[id].OnCastComplete(sim, action)
		}
	}
	for _, id := range sim.ActiveAuraIDs {
		if sim.Auras[id].OnCastComplete != nil {
			sim.Auras[id].OnCastComplete(sim, action)
		}
	}

	hitInputs := action.GetHitInputs(sim, action)

	results := make([]DirectCastDamageResult, 0, len(hitInputs))
	for _, hitInput := range hitInputs {
		result := action.calculateDirectCastDamage(sim, hitInput)
		results = append(results, result)

		if result.Hit {
			// Apply any on spell hit effects.
			action.OnSpellHit(sim, action, &result)
			for _, id := range character.ActiveAuraIDs {
				if character.Auras[id].OnSpellHit != nil {
					character.Auras[id].OnSpellHit(sim, action, &result)
				}
			}
			for _, id := range sim.ActiveAuraIDs {
				if sim.Auras[id].OnSpellHit != nil {
					sim.Auras[id].OnSpellHit(sim, action, &result)
				}
			}
		} else {
			action.OnSpellMiss(sim, action)
			for _, id := range character.ActiveAuraIDs {
				if character.Auras[id].OnSpellMiss != nil {
					character.Auras[id].OnSpellMiss(sim, action)
				}
			}
			for _, id := range sim.ActiveAuraIDs {
				if sim.Auras[id].OnSpellMiss != nil {
					sim.Auras[id].OnSpellMiss(sim, action)
				}
			}
		}

		if sim.Log != nil {
			sim.Log("(%d) %s result: %s\n", character.ID, action.GetName(), result)
		}
	}

	cooldown := action.GetCooldown()
	if cooldown > 0 {
		character.SetCD(action.GetActionID().CooldownID, sim.CurrentTime+cooldown)
	}

	sim.metricsAggregator.addCastAction(action, results)
}

func (action DirectCastAction) calculateDirectCastDamage(sim *Simulation, damageInput DirectCastDamageInput) DirectCastDamageResult {
	result := DirectCastDamageResult{}

	character := action.GetAgent().GetCharacter()

	hit := 0.83 + character.Stats[stats.SpellHit]/1260.0 + damageInput.BonusHit // 12.6 hit == 1% hit
	hit = MinFloat(hit, 0.99)                                                    // can't get away from the 1% miss

	if sim.Rando.Float64("action hit") >= hit { // Miss
		return result
	}
	result.Hit = true

	baseDamage := damageInput.MinBaseDamage + sim.Rando.Float64("action dmg")*(damageInput.MaxBaseDamage - damageInput.MinBaseDamage)
	totalSpellPower := character.Stats[stats.SpellPower] + character.Stats[action.GetSpellSchool()] + damageInput.BonusSpellPower
	damageFromSpellPower := (totalSpellPower * damageInput.SpellCoefficient)
	damage := baseDamage + damageFromSpellPower

	damage *= damageInput.DamageMultiplier

	crit := (character.Stats[stats.SpellCrit] / (SpellCritRatingPerCritChance * 100)) + damageInput.BonusCrit
	// TODO: Put guaranteed crit first, to short-circuit the random roll. Keeping it this way for now for tests.
	if sim.Rando.Float64("action crit") < crit || action.castInput.GuaranteedCrit {
		result.Crit = true
		damage *= action.castInput.CritMultiplier
	}

	// Average Resistance (AR) = (Target's Resistance / (Caster's Level * 5)) * 0.75
	// P(x) = 50% - 250%*|x - AR| <- where X is %resisted
	// Using these stats:
	//    13.6% chance of
	//  FUTURE: handle boss resists for fights/classes that are actually impacted by that.
	resVal := sim.Rando.Float64("action resist")
	if resVal < 0.18 { // 13% chance for 25% resist, 4% for 50%, 1% for 75%
		if resVal < 0.01 {
			result.PartialResist_3_4 = true
			damage *= .25
		} else if resVal < 0.05 {
			result.PartialResist_2_4 = true
			damage *= .5
		} else {
			result.PartialResist_1_4 = true
			damage *= .75
		}
	}

	result.Damage = damage

	return result
}

func NewDirectCastAction(sim *Simulation, impl DirectCastImpl) DirectCastAction {
	action := DirectCastAction{
		DirectCastImpl: impl,
	}

	castInput := impl.GetCastInput(sim, action)
	castInput.CastTime = time.Duration(float64(castInput.CastTime) / impl.GetAgent().GetCharacter().HasteBonus())

	// Apply on-cast effects.
	character := action.GetAgent().GetCharacter()
	for _, id := range character.ActiveAuraIDs {
		if character.Auras[id].OnCast != nil {
			character.Auras[id].OnCast(sim, action, &castInput)
		}
	}

	action.castInput = castInput

	return action
}