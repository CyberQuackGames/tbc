package core

import (
	"time"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var OrcBloodFuryAuraID = NewAuraID()
var OrcBloodFuryCooldownID = NewCooldownID()

var TrollBerserkingAuraID = NewAuraID()
var TrollBerserkingCooldownID = NewCooldownID()

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceBloodElf:
		character.AddStat(stats.ArcaneResistance, 5)
		character.AddStat(stats.FireResistance, 5)
		character.AddStat(stats.FrostResistance, 5)
		character.AddStat(stats.NatureResistance, 5)
		character.AddStat(stats.ShadowResistance, 5)
		// TODO: Add major cooldown: arcane torrent
	case proto.Race_RaceDraenei:
		character.AddStat(stats.ShadowResistance, 10)
	case proto.Race_RaceDwarf:
		character.AddStat(stats.FrostResistance, 10)

		// Gun specialization (+1% ranged crit when using a gun).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun {
				character.PseudoStats.BonusRangedCritRating += 1 * MeleeCritRatingPerCritChance
			}
		}
	case proto.Race_RaceGnome:
		character.AddStat(stats.ArcaneResistance, 10)

		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Intellect,
			ModifiedStat: stats.Intellect,
			Modifier: func(intellect float64, _ float64) float64 {
				return intellect * 1.05
			},
		})
	case proto.Race_RaceHuman:
		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Spirit,
			ModifiedStat: stats.Spirit,
			Modifier: func(spirit float64, _ float64) float64 {
				return spirit * 1.1
			},
		})

		const expertiseBonus = 5 * ExpertisePerQuarterPercentReduction
		if weapon := character.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeSword || weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
				character.PseudoStats.BonusMHExpertiseRating += expertiseBonus
			}
		}
		if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeSword || weapon.WeaponType == proto.WeaponType_WeaponTypeMace {
				character.PseudoStats.BonusOHExpertiseRating += expertiseBonus
			}
		}
	case proto.Race_RaceNightElf:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStat(stats.Dodge, DodgeRatingPerDodgeChance*1)
	case proto.Race_RaceOrc:
		// Command (Pet damage +5%)
		if len(character.Pets) > 0 {
			for _, petAgent := range character.Pets {
				pet := petAgent.GetPet()
				pet.PseudoStats.DamageDealtMultiplier *= 1.05
			}
		}

		// Blood Fury
		const cd = time.Minute * 2
		const dur = time.Second * 15
		const apBonus = float64(CharacterLevel)*4 + 2
		const spBonus = float64(CharacterLevel)*2 + 3
		actionID := ActionID{SpellID: 33697}

		character.AddMajorCooldown(MajorCooldown{
			ActionID:   actionID,
			CooldownID: OrcBloodFuryCooldownID,
			Cooldown:   cd,
			Type:       CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ActivationFactory: func(sim *Simulation) CooldownActivation {
				applyStatAura := character.NewTemporaryStatsAuraApplier(OrcBloodFuryAuraID, actionID, stats.Stats{stats.AttackPower: apBonus, stats.RangedAttackPower: apBonus, stats.SpellPower: spBonus}, dur)
				return func(sim *Simulation, character *Character) {
					applyStatAura(sim)
					character.SetCD(OrcBloodFuryCooldownID, sim.CurrentTime+cd)
					character.Metrics.AddInstantCast(actionID)
				}
			},
		})

		// Axe specialization
		const expertiseBonus = 5 * ExpertisePerQuarterPercentReduction
		if weapon := character.Equip[proto.ItemSlot_ItemSlotMainHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe {
				character.PseudoStats.BonusMHExpertiseRating += expertiseBonus
			}
		}
		if weapon := character.Equip[proto.ItemSlot_ItemSlotOffHand]; weapon.ID != 0 {
			if weapon.WeaponType == proto.WeaponType_WeaponTypeAxe {
				character.PseudoStats.BonusMHExpertiseRating += expertiseBonus
			}
		}
	case proto.Race_RaceTauren:
		character.AddStat(stats.NatureResistance, 10)
		character.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Health,
			ModifiedStat: stats.Health,
			Modifier: func(health float64, _ float64) float64 {
				return health * 1.05
			},
		})
	case proto.Race_RaceTroll10, proto.Race_RaceTroll30:
		// Bow specialization (+1% ranged crit when using a bow).
		if weapon := character.Equip[proto.ItemSlot_ItemSlotRanged]; weapon.ID != 0 {
			if weapon.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow {
				character.PseudoStats.BonusRangedCritRating += 1 * MeleeCritRatingPerCritChance
			}
		}

		// Beast Slaying (+5% damage to beasts)
		character.RegisterResetEffect(func(sim *Simulation) {
			if sim.GetPrimaryTarget().MobType == proto.MobType_MobTypeBeast {
				character.PseudoStats.DamageDealtMultiplier *= 1.05
			}
		})

		// Berserking
		hasteBonus := 1.1
		if character.Race == proto.Race_RaceTroll30 {
			hasteBonus = 1.3
		}
		inverseBonus := 1 / hasteBonus
		const dur = time.Second * 10
		const cd = time.Minute * 3

		var cost ResourceCost
		var actionID ActionID
		if character.Class == proto.Class_ClassRogue {
			actionID = ActionID{SpellID: 26297, CooldownID: TrollBerserkingCooldownID}
		} else if character.Class == proto.Class_ClassWarrior {
			actionID = ActionID{SpellID: 26296, CooldownID: TrollBerserkingCooldownID}
		} else {
			actionID = ActionID{SpellID: 20554, CooldownID: TrollBerserkingCooldownID}
		}

		character.AddMajorCooldown(MajorCooldown{
			ActionID:   actionID,
			CooldownID: TrollBerserkingCooldownID,
			Cooldown:   cd,
			Type:       CooldownTypeDPS,
			CanActivate: func(sim *Simulation, character *Character) bool {
				if character.Class == proto.Class_ClassRogue {
					return character.CurrentEnergy() >= cost.Value
				} else if character.Class == proto.Class_ClassWarrior {
					return character.CurrentRage() >= cost.Value
				} else {
					return character.CurrentMana() >= cost.Value
				}
			},
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			ActivationFactory: func(sim *Simulation) CooldownActivation {
				if character.Class == proto.Class_ClassRogue {
					cost = ResourceCost{Type: stats.Energy, Value: 10}
				} else if character.Class == proto.Class_ClassWarrior {
					cost = ResourceCost{Type: stats.Rage, Value: 5}
				} else {
					cost = ResourceCost{Type: stats.Mana, Value: character.BaseMana() * 0.06}
				}

				castTemplate := SimpleCast{
					Cast: Cast{
						ActionID:  actionID,
						Character: character,
						BaseCost:  cost,
						Cost:      cost,
						Cooldown:  cd,
						OnCastComplete: func(sim *Simulation, cast *Cast) {
							character.AddAura(sim, Aura{
								ID:       TrollBerserkingAuraID,
								ActionID: actionID,
								Duration: dur,
								OnGain: func(sim *Simulation) {
									character.PseudoStats.CastSpeedMultiplier *= hasteBonus
									character.MultiplyAttackSpeed(sim, hasteBonus)
								},
								OnExpire: func(sim *Simulation) {
									character.PseudoStats.CastSpeedMultiplier /= hasteBonus
									character.MultiplyAttackSpeed(sim, inverseBonus)
								},
							})
						},
					},
				}

				return func(sim *Simulation, character *Character) {
					cast := castTemplate
					cast.Init(sim)
					cast.StartCast(sim)
				}
			},
		})
	case proto.Race_RaceUndead:
		character.AddStat(stats.ShadowResistance, 10)
	}
}
