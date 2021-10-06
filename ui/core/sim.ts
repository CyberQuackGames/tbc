import { EquippedItem } from './api/equipped_item.js';
import { Gear } from './api/gear.js';
import { Buffs } from './proto/common.js';
import { Class } from './proto/common.js';
import { Consumes } from './proto/common.js';
import { Enchant } from './proto/common.js';
import { Encounter } from './proto/common.js';
import { EquipmentSpec } from './proto/common.js';
import { Gem } from './proto/common.js';
import { GemColor } from './proto/common.js';
import { ItemQuality } from './proto/common.js';
import { ItemSlot } from './proto/common.js';
import { ItemSpec } from './proto/common.js';
import { ItemType } from './proto/common.js';
import { Item } from './proto/common.js';
import { Race } from './proto/common.js';
import { Spec } from './proto/common.js';
import { Stat } from './proto/common.js';
import { makeComputeStatsRequest } from './api/request_helpers.js';
import { makeIndividualSimRequest } from './api/request_helpers.js';
import { Stats } from './api/stats.js';
import { SpecAgent } from './api/utils.js';
import { SpecTalents } from './api/utils.js';
import { SpecTypeFunctions } from './api/utils.js';
import { specTypeFunctions } from './api/utils.js';
import { SpecOptions } from './api/utils.js';
import { specToClass } from './api/utils.js';
import { specToEligibleRaces } from './api/utils.js';
import { getEligibleItemSlots } from './api/utils.js';
import { getEligibleEnchantSlots } from './api/utils.js';
import { gemEligibleForSocket } from './api/utils.js';
import { gemMatchesSocket } from './api/utils.js';

import { Player } from './proto/api.js';
import { PlayerOptions } from './proto/api.js';
import { ComputeStatsRequest, ComputeStatsResult } from './proto/api.js';
import { GearListRequest, GearListResult } from './proto/api.js';
import { IndividualSimRequest, IndividualSimResult } from './proto/api.js';
import { StatWeightsRequest, StatWeightsResult } from './proto/api.js';

import { Listener } from './typed_event.js';
import { TypedEvent } from './typed_event.js';
import { sum } from './utils.js';
import { wait } from './utils.js';
import { WorkerPool } from './worker_pool.js';

export interface SimConfig<SpecType extends Spec> {
  spec: Spec;
  epStats: Array<Stat>;
  epReferenceStat: Stat;
  defaults: {
		phase: number,
		epWeights: Stats,
    encounter: Encounter,
    buffs: Buffs,
    consumes: Consumes,
    agent: SpecAgent<SpecType>,
    talents: string,
    specOptions: SpecOptions<SpecType>,
  },
	metaGemEffectEP?: ((gem: Gem, sim: Sim<SpecType>) => number),
}

// Core Sim module which deals only with api types, no UI-related stuff.
export class Sim<SpecType extends Spec> extends WorkerPool {
  readonly spec: Spec;

  readonly phaseChangeEmitter = new TypedEvent<void>();
  readonly buffsChangeEmitter = new TypedEvent<void>();
  readonly consumesChangeEmitter = new TypedEvent<void>();
  readonly customStatsChangeEmitter = new TypedEvent<void>();
  readonly encounterChangeEmitter = new TypedEvent<void>();
  readonly gearChangeEmitter = new TypedEvent<void>();
  readonly raceChangeEmitter = new TypedEvent<void>();
  readonly agentChangeEmitter = new TypedEvent<void>();
  readonly talentsChangeEmitter = new TypedEvent<void>();
  // Talents dont have all fields so we need this
  readonly talentsStringChangeEmitter = new TypedEvent<void>();
  readonly specOptionsChangeEmitter = new TypedEvent<void>();

  // Emits when any of the above emitters emit.
  readonly changeEmitter = new TypedEvent<void>();

  readonly gearListEmitter = new TypedEvent<void>();
  readonly characterStatsEmitter = new TypedEvent<void>();
	private _currentStats: ComputeStatsResult;

  // Database
  private _items: Record<number, Item> = {};
  private _enchants: Record<number, Enchant> = {};
  private _gems: Record<number, Gem> = {};

  // Current values
  private _phase: number;
  private _buffs: Buffs;
  private _consumes: Consumes;
  private _customStats: Stats;
  private _gear: Gear;
  private _encounter: Encounter;
  private _race: Race;
  private _agent: SpecAgent<SpecType>;
  private _talents: SpecTalents<SpecType>;
  private _talentsString: string;
  private _specOptions: SpecOptions<SpecType>;
	private _epWeights: Stats;

  readonly specTypeFunctions: SpecTypeFunctions<SpecType>;
	private readonly _metaGemEffectEP: (gem: Gem, sim: Sim<SpecType>) => number;

  private _init = false;

  constructor(config: SimConfig<SpecType>) {
		super(3);

    this.spec = config.spec;
    this._race = specToEligibleRaces[this.spec][0];

    this.specTypeFunctions = specTypeFunctions[this.spec] as SpecTypeFunctions<SpecType>;
		this._metaGemEffectEP = config.metaGemEffectEP || (() => 0);

		this._phase = config.defaults.phase;
    this._buffs = config.defaults.buffs;
    this._consumes = config.defaults.consumes;
    this._customStats = new Stats();
    this._encounter = config.defaults.encounter;
    this._gear = new Gear({});
    this._agent = config.defaults.agent;
    this._talents = this.specTypeFunctions.talentsCreate();
    this._talentsString = config.defaults.talents;
		this._epWeights = config.defaults.epWeights;
    this._specOptions = config.defaults.specOptions;

    [
      this.buffsChangeEmitter,
      this.consumesChangeEmitter,
      this.customStatsChangeEmitter,
      this.encounterChangeEmitter,
      this.gearChangeEmitter,
      this.raceChangeEmitter,
      this.agentChangeEmitter,
      this.talentsChangeEmitter,
      this.talentsStringChangeEmitter,
      this.specOptionsChangeEmitter,
    ].forEach(emitter => emitter.on(() => this.changeEmitter.emit()));

		this._currentStats = ComputeStatsResult.create();
		this.changeEmitter.on(() => {
			this.updateCharacterStats();
		});
  }

  async init(): Promise<void> {
    if (this._init)
      return;
    this._init = true;

    const result = await this.getGearList(GearListRequest.create({
      spec: this.spec,
    }));

    result.items.forEach(item => this._items[item.id] = item);
    result.enchants.forEach(enchant => this._enchants[enchant.id] = enchant);
    result.gems.forEach(gem => this._gems[gem.id] = gem);

    this.gearListEmitter.emit();
  }

  async statWeights(request: StatWeightsRequest): Promise<StatWeightsResult> {
		const result = await super.statWeights(request);
		this._epWeights = new Stats(result.epValues);
		return result;
	}

	// This should be invoked internally whenever stats might have changed.
	private async updateCharacterStats() {
		// Sometimes a ui change triggers other changes, so waiting a bit makes sure
		// we get all of them.
		await wait(10);

		const computeStatsResult = await this.computeStats(makeComputeStatsRequest(
      this._buffs,
      this._consumes,
      this._customStats,
      this._encounter,
      this._gear,
      this._race,
      this._agent,
      this._talents,
      this._specOptions));

		this._currentStats = computeStatsResult;
		this.characterStatsEmitter.emit();
	}

	getCurrentStats(): ComputeStatsResult {
		return ComputeStatsResult.clone(this._currentStats);
	}

	getItems(slot: ItemSlot | undefined): Array<Item> {
		let items = Object.values(this._items);
		if (slot != undefined) {
			items = items.filter(item => getEligibleItemSlots(item).includes(slot));
		}
		return items;
	}

	getEnchants(slot: ItemSlot | undefined): Array<Enchant> {
		let enchants = Object.values(this._enchants);
		if (slot != undefined) {
			enchants = enchants.filter(enchant => getEligibleEnchantSlots(enchant).includes(slot));
		}
		return enchants;
	}

  getGems(socketColor: GemColor | undefined): Array<Gem> {
    let gems = Object.values(this._gems);
		if (socketColor) {
			gems = gems.filter(gem => gemEligibleForSocket(gem, socketColor));
		}
		return gems;
  }

	getMatchingGems(socketColor: GemColor): Array<Gem> {
    return Object.values(this._gems).filter(gem => gemMatchesSocket(gem, socketColor));
	}
  
  getPhase(): number {
    return this._phase;
  }
  setPhase(newPhase: number) {
    if (newPhase != this._phase) {
      this._phase = newPhase;
      this.phaseChangeEmitter.emit();
    }
  }
  
  getRace(): Race {
    return this._race;
  }
  setRace(newRace: Race) {
    if (newRace != this._race) {
      this._race = newRace;
      this.raceChangeEmitter.emit();
    }
  }

  getBuffs(): Buffs {
    // Make a defensive copy
    return Buffs.clone(this._buffs);
  }

  setBuffs(newBuffs: Buffs) {
    if (Buffs.equals(this._buffs, newBuffs))
      return;

    // Make a defensive copy
    this._buffs = Buffs.clone(newBuffs);
    this.buffsChangeEmitter.emit();
  }

  getConsumes(): Consumes {
    // Make a defensive copy
    return Consumes.clone(this._consumes);
  }

  setConsumes(newConsumes: Consumes) {
    if (Consumes.equals(this._consumes, newConsumes))
      return;

    // Make a defensive copy
    this._consumes = Consumes.clone(newConsumes);
    this.consumesChangeEmitter.emit();
  }

  getEncounter(): Encounter {
    // Make a defensive copy
    return Encounter.clone(this._encounter);
  }

  setEncounter(newEncounter: Encounter) {
    if (Encounter.equals(this._encounter, newEncounter))
      return;

    // Make a defensive copy
    this._encounter = Encounter.clone(newEncounter);
    this.encounterChangeEmitter.emit();
  }

  equipItem(slot: ItemSlot, newItem: EquippedItem | null) {
    const newGear = this._gear.withEquippedItem(slot, newItem);
    if (newGear.equals(this._gear))
      return;

    this._gear = newGear;
    this.gearChangeEmitter.emit();
  }

  getEquippedItem(slot: ItemSlot): EquippedItem | null {
    return this._gear.getEquippedItem(slot);
  }

  getGear(): Gear {
    return this._gear;
  }

  setGear(newGear: Gear) {
    if (newGear.equals(this._gear))
      return;

    this._gear = newGear;
    this.gearChangeEmitter.emit();
  }

  getCustomStats(): Stats {
    return this._customStats;
  }

  setCustomStats(newCustomStats: Stats) {
    if (newCustomStats.equals(this._customStats))
      return;

    this._customStats = newCustomStats;
    this.customStatsChangeEmitter.emit();
  }

  getAgent(): SpecAgent<SpecType> {
    return this.specTypeFunctions.agentCopy(this._agent);
  }

  setAgent(newAgent: SpecAgent<SpecType>) {
    if (this.specTypeFunctions.agentEquals(newAgent, this._agent))
      return;

    this._agent = this.specTypeFunctions.agentCopy(newAgent);
    this.agentChangeEmitter.emit();
  }

  // Commented because this should NOT be used; all external uses should be able to use getTalentsString()
  //getTalents(): SpecTalents<SpecType> {
  //  return this.specTypeFunctions.talentsCopy(this._talents);
  //}

  setTalents(newTalents: SpecTalents<SpecType>) {
    if (this.specTypeFunctions.talentsEquals(newTalents, this._talents))
      return;

    this._talents = this.specTypeFunctions.talentsCopy(newTalents);
    this.talentsChangeEmitter.emit();
  }

  getTalentsString(): string {
    return this._talentsString;
  }

  setTalentsString(newTalentsString: string) {
    if (newTalentsString == this._talentsString)
      return;

    this._talentsString = newTalentsString;
    this.talentsStringChangeEmitter.emit();
  }

  getSpecOptions(): SpecOptions<SpecType> {
    return this.specTypeFunctions.optionsCopy(this._specOptions);
  }

  setSpecOptions(newSpecOptions: SpecOptions<SpecType>) {
    if (this.specTypeFunctions.optionsEquals(newSpecOptions, this._specOptions))
      return;

    this._specOptions = this.specTypeFunctions.optionsCopy(newSpecOptions);
    this.specOptionsChangeEmitter.emit();
  }

  lookupItemSpec(itemSpec: ItemSpec): EquippedItem | null {
    const item = this._items[itemSpec.id];
    if (!item)
      return null;

    const enchant = this._enchants[itemSpec.enchant] || null;
    const gems = itemSpec.gems.map(gemId => this._gems[gemId] || null);

    return new EquippedItem(item, enchant, gems);
  }

  lookupEquipmentSpec(equipSpec: EquipmentSpec): Gear {
    // EquipmentSpec is supposed to be indexed by slot, but here we assume
    // it isn't just in case.
    const gearMap: Partial<Record<ItemSlot, EquippedItem | null>> = {};

    equipSpec.items.forEach(itemSpec => {
      const item = this.lookupItemSpec(itemSpec);
      if (!item)
        return;

      const itemSlots = getEligibleItemSlots(item.item);

      const assignedSlot = itemSlots.find(slot => !gearMap[slot]);
      if (assignedSlot == null)
        throw new Error('No slots left to equip ' + Item.toJsonString(item.item));

      gearMap[assignedSlot] = item;
    });

    return new Gear(gearMap);
  }

	computeGemEP(gem: Gem): number {
		const epFromStats = new Stats(gem.stats).computeEP(this._epWeights);
		const epFromEffect = this._metaGemEffectEP(gem, this);
		return epFromStats + epFromEffect;
	}

	computeEnchantEP(enchant: Enchant): number {
		return new Stats(enchant.stats).computeEP(this._epWeights);
	}

	computeItemEP(item: Item): number {
		if (item == null)
			return 0;

		let ep = new Stats(item.stats).computeEP(this._epWeights);

		const slot = getEligibleItemSlots(item)[0];
		const enchants = this.getEnchants(slot);
		if (enchants.length > 0) {
			ep += Math.max(...enchants.map(enchant => this.computeEnchantEP(enchant)));
		}

		// Compare whether its better to match sockets + get socket bonus, or just use best gems.
		const bestGemEPNotMatchingSockets = sum(item.gemSockets.map(socketColor => {
			const gems = this.getGems(socketColor).filter(gem => !gem.unique && gem.phase <= this.getPhase());
			if (gems.length > 0) {
				return Math.max(...gems.map(gem => this.computeGemEP(gem)));
			} else {
				return 0;
			}
		}));

		const bestGemEPMatchingSockets = sum(item.gemSockets.map(socketColor => {
			const gems = this.getGems(socketColor).filter(gem => !gem.unique && gem.phase <= this.getPhase() && gemMatchesSocket(gem, socketColor));
			if (gems.length > 0) {
				return Math.max(...gems.map(gem => this.computeGemEP(gem)));
			} else {
				return 0;
			}
		})) + new Stats(item.socketBonus).computeEP(this._epWeights);

		ep += Math.max(bestGemEPMatchingSockets, bestGemEPNotMatchingSockets);

		return ep;
	}

  makeCurrentIndividualSimRequest(iterations: number, debug: boolean): IndividualSimRequest {
    return makeIndividualSimRequest(
      this._buffs,
      this._consumes,
      this._customStats,
      this._encounter,
      this._gear,
      this._race,
      this._agent,
      this._talents,
      this._specOptions,
      iterations,
      debug);
  }

  setWowheadData(equippedItem: EquippedItem, elem: HTMLElement) {
    let parts = [];
    if (equippedItem.gems.length > 0) {
      parts.push('gems=' + equippedItem.gems.map(gem => gem ? gem.id : 0).join(':'));
    }
    if (equippedItem.enchant != null) {
      parts.push('ench=' + equippedItem.enchant.effectId);
    }
    parts.push('pcs=' + this._gear.asArray().filter(ei => ei != null).map(ei => ei!.item.id).join(':'));

    elem.setAttribute('data-wowhead', parts.join('&'));
  }

  // Returns JSON representing all the current values.
  toJson(): Object {
    return {
      'buffs': Buffs.toJson(this._buffs),
      'consumes': Consumes.toJson(this._consumes),
      'customStats': this._customStats.toJson(),
      'encounter': Encounter.toJson(this._encounter),
      'gear': EquipmentSpec.toJson(this._gear.asSpec()),
      'race': this._race,
      'agent': this.specTypeFunctions.agentToJson(this._agent),
      'talents': this._talentsString,
      'specOptions': this.specTypeFunctions.optionsToJson(this._specOptions),
    };
  }

  // Set all the current values, assumes obj is the same type returned by toJson().
  fromJson(obj: any) {
    this.setBuffs(Buffs.fromJson(obj['buffs']));
    this.setConsumes(Consumes.fromJson(obj['consumes']));
    this.setCustomStats(Stats.fromJson(obj['customStats']));
    this.setEncounter(Encounter.fromJson(obj['encounter']));
    this.setGear(this.lookupEquipmentSpec(EquipmentSpec.fromJson(obj['gear'])));
    this.setRace(obj['race']);
    this.setAgent(this.specTypeFunctions.agentFromJson(obj['agent']));
    this.setTalentsString(obj['talents']);
    this.setSpecOptions(this.specTypeFunctions.optionsFromJson(obj['specOptions']));
  }
}