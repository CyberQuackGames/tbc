package elemental

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
	. "github.com/wowsims/tbc/sim/shaman"
)

func RegisterElementalShaman() {
	core.RegisterAgentFactory(proto.PlayerOptions_ElementalShaman{}, func(sim *core.Simulation, character core.Character, options *proto.PlayerOptions) core.Agent {
		return NewElementalShaman(sim, character, options)
	})
}

func NewElementalShaman(sim *core.Simulation, character core.Character, options *proto.PlayerOptions) *Shaman {
	eleShamOptions := options.GetElementalShaman()

	selfBuffs := SelfBuffs{
		Bloodlust:    eleShamOptions.Options.Bloodlust,
		ManaSpring:   eleShamOptions.Options.ManaSpringTotem,
		TotemOfWrath: eleShamOptions.Options.TotemOfWrath,
		WrathOfAir:   eleShamOptions.Options.WrathOfAirTotem,
		WaterShield:  eleShamOptions.Options.WaterShield,
	}

	var rotation Rotation

	switch eleShamOptions.Rotation.Type {
	case proto.ElementalShaman_Rotation_Adaptive:
		rotation = NewAdaptiveRotation(sim)
	case proto.ElementalShaman_Rotation_CLOnClearcast:
		rotation = NewCLOnClearcastRotation(sim)
	case proto.ElementalShaman_Rotation_CLOnCD:
		rotation = NewCLOnCDRotation(sim)
	case proto.ElementalShaman_Rotation_FixedLBCL:
		rotation = NewFixedRotation(sim, eleShamOptions.Rotation.LbsPerCl)
	case proto.ElementalShaman_Rotation_LBOnly:
		rotation = NewLBOnlyRotation(sim)
	}

	return NewShaman(character, *eleShamOptions.Talents, selfBuffs, rotation)
}

// ################################################################
//                              LB ONLY
// ################################################################
type LBOnlyRotation struct {
}

func (agent *LBOnlyRotation) ChooseAction(shaman *Shaman, sim *core.Simulation) core.AgentAction {
	return NewLightningBolt(sim, shaman, false)
}

func (agent *LBOnlyRotation) OnActionAccepted(shaman *Shaman, sim *core.Simulation, action core.AgentAction) {
}
func (agent *LBOnlyRotation) Reset(shaman *Shaman, sim *core.Simulation) {}

func NewLBOnlyRotation(sim *core.Simulation) *LBOnlyRotation {
	return &LBOnlyRotation{}
}

// ################################################################
//                             CL ON CD
// ################################################################
type CLOnCDRotation struct {
}

func (agent *CLOnCDRotation) ChooseAction(shaman *Shaman, sim *core.Simulation) core.AgentAction {
	if shaman.IsOnCD(core.MagicIDChainLightning6, sim.CurrentTime) {
		return NewLightningBolt(sim, shaman, false)
	} else {
		return NewChainLightning(sim, shaman, false)
	}
}

func (agent *CLOnCDRotation) OnActionAccepted(shaman *Shaman, sim *core.Simulation, action core.AgentAction) {
}
func (agent *CLOnCDRotation) Reset(shaman *Shaman, sim *core.Simulation) {}

func NewCLOnCDRotation(sim *core.Simulation) *CLOnCDRotation {
	return &CLOnCDRotation{}
}

// ################################################################
//                          FIXED ROTATION
// ################################################################
type FixedRotation struct {
	numLBsPerCL       int32
	numLBsSinceLastCL int32
}

// Returns if any temporary haste buff is currently active.
// TODO: Figure out a way to make this automatic
func (agent *FixedRotation) temporaryHasteActive(shaman *Shaman) bool {
	return shaman.HasAura(core.MagicIDBloodlust) ||
		shaman.HasAura(core.MagicIDDrums) ||
		shaman.HasAura(core.MagicIDTrollBerserking) ||
		shaman.HasAura(core.MagicIDSkullGuldan) ||
		shaman.HasAura(core.MagicIDFungalFrenzy)
}

func (agent *FixedRotation) ChooseAction(shaman *Shaman, sim *core.Simulation) core.AgentAction {
	if agent.numLBsSinceLastCL < agent.numLBsPerCL {
		return NewLightningBolt(sim, shaman, false)
	}

	if !shaman.IsOnCD(core.MagicIDChainLightning6, sim.CurrentTime) {
		return NewChainLightning(sim, shaman, false)
	}

	// If we have a temporary haste effect (like bloodlust or quags eye) then
	// we should add LB casts instead of waiting
	if agent.temporaryHasteActive(shaman) {
		return NewLightningBolt(sim, shaman, false)
	}

	return core.NewWaitAction(sim, shaman.GetCharacter(), shaman.GetRemainingCD(core.MagicIDChainLightning6, sim.CurrentTime))
}

func (agent *FixedRotation) OnActionAccepted(shaman *Shaman, sim *core.Simulation, action core.AgentAction) {
	cast, isCastAction := action.(*core.DirectCastAction)
	if !isCastAction {
		return
	}

	if cast.GetActionID().SpellID == SpellIDLB12 {
		agent.numLBsSinceLastCL++
	} else if cast.GetActionID().SpellID == SpellIDCL6 {
		agent.numLBsSinceLastCL = 0
	}
}

func (agent *FixedRotation) Reset(shaman *Shaman, sim *core.Simulation) {
	agent.numLBsSinceLastCL = agent.numLBsPerCL // This lets us cast CL first
}

func NewFixedRotation(sim *core.Simulation, numLBsPerCL int32) *FixedRotation {
	return &FixedRotation{
		numLBsPerCL: numLBsPerCL,
	}
}

// ################################################################
//                          CL ON CLEARCAST
// ################################################################
type CLOnClearcastRotation struct {
	// Whether the second-to-last spell procced clearcasting
	prevPrevCastProccedCC bool
}

func (agent *CLOnClearcastRotation) ChooseAction(shaman *Shaman, sim *core.Simulation) core.AgentAction {
	if shaman.IsOnCD(core.MagicIDChainLightning6, sim.CurrentTime) || !agent.prevPrevCastProccedCC {
		return NewLightningBolt(sim, shaman, false)
	}

	return NewChainLightning(sim, shaman, false)
}

func (agent *CLOnClearcastRotation) OnActionAccepted(shaman *Shaman, sim *core.Simulation, action core.AgentAction) {
	agent.prevPrevCastProccedCC = shaman.ElementalFocusStacks == 2
}

func (agent *CLOnClearcastRotation) Reset(shaman *Shaman, sim *core.Simulation) {
	agent.prevPrevCastProccedCC = true // Lets us cast CL first
}

func NewCLOnClearcastRotation(sim *core.Simulation) *CLOnClearcastRotation {
	return &CLOnClearcastRotation{}
}

// ################################################################
//                             ADAPTIVE
// ################################################################
type AdaptiveRotation struct {
	// Circular array buffer for recent mana snapshots, within a time window
	manaSnapshots      [manaSnapshotsBufferSize]ManaSnapshot
	numSnapshots       int32
	firstSnapshotIndex int32

	baseRotation    Rotation // The agent used most of the time
	surplusRotation Rotation // The agent used when we have extra mana
}

const manaSpendingWindowNumSeconds = 60
const manaSpendingWindow = time.Second * manaSpendingWindowNumSeconds

// 2 * (# of seconds) should be plenty of slots
const manaSnapshotsBufferSize = manaSpendingWindowNumSeconds * 2

type ManaSnapshot struct {
	time      time.Duration // time this snapshot was taken
	manaSpent float64       // total amount of mana spent up to this time
}

func (agent *AdaptiveRotation) getOldestSnapshot() ManaSnapshot {
	return agent.manaSnapshots[agent.firstSnapshotIndex]
}

func (agent *AdaptiveRotation) purgeExpiredSnapshots(sim *core.Simulation) {
	expirationCutoff := sim.CurrentTime - manaSpendingWindow

	curIndex := agent.firstSnapshotIndex
	for agent.numSnapshots > 0 && agent.manaSnapshots[curIndex].time < expirationCutoff {
		curIndex = (curIndex + 1) % manaSnapshotsBufferSize
		agent.numSnapshots--
	}
	agent.firstSnapshotIndex = curIndex
}

func (agent *AdaptiveRotation) takeSnapshot(sim *core.Simulation, shaman *Shaman) {
	if agent.numSnapshots >= manaSnapshotsBufferSize {
		panic("Rotation snapshot buffer full")
	}

	snapshot := ManaSnapshot{
		time:      sim.CurrentTime,
		manaSpent: sim.GetIndividualMetrics(shaman.ID).ManaSpent,
	}

	nextIndex := (agent.firstSnapshotIndex + agent.numSnapshots) % manaSnapshotsBufferSize
	agent.manaSnapshots[nextIndex] = snapshot
	agent.numSnapshots++
}

func (agent *AdaptiveRotation) ChooseAction(shaman *Shaman, sim *core.Simulation) core.AgentAction {
	agent.purgeExpiredSnapshots(sim)
	oldestSnapshot := agent.getOldestSnapshot()

	manaSpent := sim.GetIndividualMetrics(shaman.ID).ManaSpent - oldestSnapshot.manaSpent
	timeDelta := sim.CurrentTime - oldestSnapshot.time
	if timeDelta == 0 {
		timeDelta = 1
	}

	timeRemaining := sim.Duration - sim.CurrentTime
	projectedManaCost := manaSpent * (timeRemaining.Seconds() / timeDelta.Seconds())

	if sim.Log != nil {
		manaSpendingRate := manaSpent / timeDelta.Seconds()
		sim.Log("[AI] CL Ready: Mana/s: %0.1f, Est Mana Cost: %0.1f, CurrentMana: %0.1f\n", manaSpendingRate, projectedManaCost, shaman.Stats[stats.Mana])
	}

	// If we have enough mana to burn, use the surplus agent.
	if projectedManaCost < shaman.Stats[stats.Mana] {
		return agent.surplusRotation.ChooseAction(shaman, sim)
	} else {
		return agent.baseRotation.ChooseAction(shaman, sim)
	}
}
func (agent *AdaptiveRotation) OnActionAccepted(shaman *Shaman, sim *core.Simulation, action core.AgentAction) {
	agent.takeSnapshot(sim, shaman)
	agent.baseRotation.OnActionAccepted(shaman, sim, action)
	agent.surplusRotation.OnActionAccepted(shaman, sim, action)
}

func (agent *AdaptiveRotation) Reset(shaman *Shaman, sim *core.Simulation) {
	agent.manaSnapshots = [manaSnapshotsBufferSize]ManaSnapshot{}
	agent.firstSnapshotIndex = 0
	agent.numSnapshots = 0
	agent.baseRotation.Reset(shaman, sim)
	agent.surplusRotation.Reset(shaman, sim)
}

func NewAdaptiveRotation(sim *core.Simulation) *AdaptiveRotation {
	agent := &AdaptiveRotation{}

	clearcastParams := sim.IndividualParams
	clearcastParams.Options.Debug = false
	clearcastParams.Options.Iterations = 100

	eleShamParams := *clearcastParams.PlayerOptions.GetElementalShaman()
	eleShamParams.Rotation = &proto.ElementalShaman_Rotation{Type: proto.ElementalShaman_Rotation_CLOnClearcast} // create new agent.

	// Assign new eleShamParams
	clearcastParams.PlayerOptions = &proto.PlayerOptions{
		Race: sim.IndividualParams.PlayerOptions.Race, //primitive, no pointer
		Spec: &proto.PlayerOptions_ElementalShaman{
			ElementalShaman: &eleShamParams,
		},
		// reuse pointer since this isn't mutated
		Consumes: sim.IndividualParams.PlayerOptions.Consumes,
	}

	clearcastSim := core.NewIndividualSim(clearcastParams)
	clearcastResult := clearcastSim.Run()

	if clearcastResult.Agents[0].NumOom >= 5 {
		agent.baseRotation = NewLBOnlyRotation(sim)
		agent.surplusRotation = NewCLOnClearcastRotation(sim)
	} else {
		agent.baseRotation = NewCLOnClearcastRotation(sim)
		agent.surplusRotation = NewCLOnCDRotation(sim)
	}

	return agent
}