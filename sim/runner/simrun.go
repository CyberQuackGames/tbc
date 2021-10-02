package runner

import (
	"fmt"
	"strings"

	"github.com/wowsims/tbc/items"
	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
	"github.com/wowsims/tbc/sim/druid"
	"github.com/wowsims/tbc/sim/mage"
	"github.com/wowsims/tbc/sim/paladin"
	"github.com/wowsims/tbc/sim/priest"
)

// TODO: Should we move the 'buff bots' to a subpackage so we dont have to import the full class unless we are actually wanting the whole class?

type AgentCreator interface {
	CreateAgent(player *core.Player, party *core.Party) core.Agent
}

type IndividualParams struct {
	Equip    items.EquipmentSpec
	Race     core.RaceBonusType
	Consumes core.Consumes
	Buffs    core.Buffs
	Options  core.Options

	Spec        AgentCreator // TODO: solve this better
	CustomStats stats.Stats
}

func SetupIndividualSim(params IndividualParams) *core.Simulation {
	player := core.NewPlayer(params.Equip, params.Race, params.Consumes)

	// TODO: should this be moved into the player constructor?
	for k, v := range params.CustomStats {
		player.Stats[k] += v
		player.InitialStats[k] += v
	}

	party := &core.Party{
		Players: []core.PlayerAgent{
			{Player: player},
		},
	}

	agent := params.Spec.CreateAgent(player, party)

	party.Players[0].Agent = agent
	raid := &core.Raid{Parties: []*core.Party{party}}

	buffs := params.Buffs
	options := params.Options

	sim := core.NewSim(raid, options)

	// These buffs are a one-time apply... no need to add the bots to the raid group.
	//  the constructors apply their buffs to the party.
	druid.NewBuffBot(sim, raid.Parties[0], buffs.GiftOfTheWild, buffs.Moonkin, buffs.MoonkinRavenGoddess)
	mage.NewBuffBot(sim, raid.Parties[0], buffs.ArcaneInt)

	// These apply auras on every sim reset
	priestBot := priest.NewBuffBot(sim, raid.Parties[0], buffs.Misery, float64(buffs.SpriestDPS))
	paladinBot := paladin.NewBuffBot(sim, raid.Parties[0], buffs.BlessingOfKings, buffs.ImprovedBlessingOfWisdom, buffs.ImprovedSealOfTheCrusader, buffs.JudgementOfWisdom)

	// Create a fake player and add the agent to do the buffing.
	if buffs.Misery {
		// Misery bot re-applies misery on every sim reset.
		sim.Raid.Parties[0].Players = append(sim.Raid.Parties[0].Players, core.PlayerAgent{
			Player: core.NewPlayer(items.EquipmentSpec{}, core.RaceBonusTypeNone, core.Consumes{}),
			Agent:  priestBot,
		})
	}
	if buffs.JudgementOfWisdom {
		// Judgement of wisdom is an aura that has to be reapplied on every reset.
		// create a bot that acts like a player and rebuffs us.
		sim.Raid.Parties[0].Players = append(sim.Raid.Parties[0].Players, core.PlayerAgent{
			Player: core.NewPlayer(items.EquipmentSpec{}, core.RaceBonusTypeNone, core.Consumes{}),
			Agent:  paladinBot,
		})
	}

	if len(sim.Raid.Parties[0].Players) == 1 && (buffs.TwilightOwl) {
		// Add a new player.
		sim.Raid.Parties[0].Players = append(sim.Raid.Parties[0].Players, core.PlayerAgent{
			Player: core.NewPlayer(items.EquipmentSpec{}, core.RaceBonusTypeNone, core.Consumes{}),
			Agent:  &nullAgent{}, // this player exists to pop items, no agent needed.
		})
	}

	if buffs.TwilightOwl {
		// Add neck to first bot player
		for i, item := range sim.Raid.Parties[0].Players[1].Equip {
			if item.ID == 0 { // no item in this slot.
				sim.Raid.Parties[0].Players[1].Equip[i] = items.ByID[24121]
				active := core.ActiveItemByID[24121]
				sim.Raid.Parties[0].Players[1].ActiveEquip = append(sim.Raid.Parties[0].Players[1].ActiveEquip, &active)
				break
			}
		}
	}
	if buffs.EyeOfNight {
		// Add neck to first bot player
		for i, item := range sim.Raid.Parties[0].Players[1].Equip {
			if item.ID == 0 { // no item in this slot.
				sim.Raid.Parties[0].Players[1].Equip[i] = items.ByID[24116]
				active := core.ActiveItemByID[24116]
				sim.Raid.Parties[0].Players[1].ActiveEquip = append(sim.Raid.Parties[0].Players[1].ActiveEquip, &active)
				break
			}
		}
	}

	sim.Reset()

	// Now apply all the 'final' stat improvements.
	// TODO: Figure out how to handle buffs that buff based on other buffs...
	//   for now this hardcoded buffing works...
	for _, raidParty := range sim.Raid.Parties {
		for _, pl := range raidParty.Players {
			if buffs.ImprovedDivineSpirit {
				pl.Player.InitialStats[stats.Spirit] += 50
			}
			if buffs.BlessingOfKings {
				pl.Player.InitialStats[stats.Intellect] *= 1.1
				pl.Player.InitialStats[stats.Spirit] *= 1.1
			}
			if buffs.ImprovedDivineSpirit {
				pl.Player.InitialStats[stats.SpellPower] += pl.Player.InitialStats[stats.Spirit] * 0.1
			}
			// Add SpellCrit from Int and Mana from Int
			pl.Player.InitialStats = pl.Player.InitialStats.CalculatedTotal()
			pl.Player.Stats = pl.Player.InitialStats
		}
	}

	return sim
}

// RunIndividualSim
func RunIndividualSim(sim *core.Simulation) SimResult {
	pid := 0
	for _, raidParty := range sim.Raid.Parties {
		for _, pl := range raidParty.Players {
			pl.ID = pid
			pl.AuraTracker.PID = pid
			pid++
		}
	}
	sim.AuraTracker.PID = -1 // set main sim auras to be -1 player ID.
	logsBuffer := &strings.Builder{}
	aggregator := NewMetricsAggregator()

	if sim.Options.Debug {
		sim.Log = func(s string, vals ...interface{}) {
			logsBuffer.WriteString(fmt.Sprintf("[%0.1f] "+s, append([]interface{}{sim.CurrentTime.Seconds()}, vals...)...))
		}
	}

	for i := 0; i < sim.Options.Iterations; i++ {
		metrics := sim.Run()
		aggregator.addMetrics(sim.Options, metrics)
		sim.ReturnCasts(metrics.Casts)
	}

	result := aggregator.getResult()
	result.Logs = logsBuffer.String()
	return result
}

type nullAgent struct {
}

func (a *nullAgent) ChooseAction(_ *core.Simulation, party *core.Party) core.AgentAction {
	return core.AgentAction{Wait: core.NeverExpires} // makes the bot wait forever and do nothing.
}
func (a *nullAgent) OnActionAccepted(*core.Simulation, core.AgentAction) {
}
func (a *nullAgent) BuffUp(sim *core.Simulation, party *core.Party) {
}
func (a *nullAgent) Reset(sim *core.Simulation) {
}
func (a *nullAgent) OnSpellHit(*core.Simulation, core.PlayerAgent, *core.Cast) {}
