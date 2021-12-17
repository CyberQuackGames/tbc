import { Player } from '/tbc/core/proto/api.js';
import { Class } from '/tbc/core/proto/common.js';
import { Enchant } from '/tbc/core/proto/common.js';
import { Gem } from '/tbc/core/proto/common.js';
import { GemColor } from '/tbc/core/proto/common.js';
import { ItemCategory } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { Item } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stats } from './stats.js';
import { BalanceDruid, BalanceDruid_Rotation as BalanceDruidRotation, DruidTalents, BalanceDruid_Options as BalanceDruidOptions } from '/tbc/core/proto/druid.js';
import { ElementalShaman, ElementalShaman_Rotation as ElementalShamanRotation, ShamanTalents, ElementalShaman_Options as ElementalShamanOptions } from '/tbc/core/proto/shaman.js';
import { Hunter, Hunter_Rotation as HunterRotation, HunterTalents, Hunter_Options as HunterOptions } from '/tbc/core/proto/hunter.js';
import { Mage, Mage_Rotation as MageRotation, MageTalents, Mage_Options as MageOptions } from '/tbc/core/proto/mage.js';
import { Rogue, Rogue_Rotation as RogueRotation, RogueTalents, Rogue_Options as RogueOptions } from '/tbc/core/proto/rogue.js';
import { RetributionPaladin, RetributionPaladin_Rotation as RetributionPaladinRotation, PaladinTalents, RetributionPaladin_Options as RetributionPaladinOptions } from '/tbc/core/proto/paladin.js';
import { ShadowPriest, ShadowPriest_Rotation as ShadowPriestRotation, PriestTalents, ShadowPriest_Options as ShadowPriestOptions } from '/tbc/core/proto/priest.js';
import { Warlock, Warlock_Rotation as WarlockRotation, WarlockTalents, Warlock_Options as WarlockOptions } from '/tbc/core/proto/warlock.js';
import { Warrior, Warrior_Rotation as WarriorRotation, WarriorTalents, Warrior_Options as WarriorOptions } from '/tbc/core/proto/warrior.js';
export declare type DruidSpecs = Spec.SpecBalanceDruid;
export declare type HunterSpecs = Spec.SpecHunter;
export declare type MageSpecs = Spec.SpecMage;
export declare type RogueSpecs = Spec.SpecRogue;
export declare type PaladinSpecs = Spec.SpecRetributionPaladin;
export declare type PriestSpecs = Spec.SpecShadowPriest;
export declare type ShamanSpecs = Spec.SpecElementalShaman;
export declare type WarlockSpecs = Spec.SpecWarlock;
export declare type WarriorSpecs = Spec.SpecWarrior;
export declare const specNames: Record<Spec, string>;
export declare const classColors: Record<Class, string>;
export declare const specIconsLarge: Record<Spec, string>;
export declare const talentTreeIcons: Record<Class, Array<string>>;
export declare function getTalentTree(talentsString: string): number;
export declare function getTalentTreeIcon(spec: Spec, talentsString: string): string;
export declare type RotationUnion = BalanceDruidRotation | ElementalShamanRotation | HunterRotation | MageRotation | RogueRotation | RetributionPaladinRotation | ShadowPriestRotation | WarlockRotation | WarriorRotation;
export declare type SpecRotation<T extends Spec> = T extends Spec.SpecBalanceDruid ? BalanceDruidRotation : T extends Spec.SpecElementalShaman ? ElementalShamanRotation : T extends Spec.SpecHunter ? HunterRotation : T extends Spec.SpecMage ? MageRotation : T extends Spec.SpecRogue ? RogueRotation : T extends Spec.SpecRetributionPaladin ? RetributionPaladinRotation : T extends Spec.SpecShadowPriest ? ShadowPriestRotation : T extends Spec.SpecWarlock ? WarlockRotation : T extends Spec.SpecWarrior ? WarriorRotation : ElementalShamanRotation;
export declare type TalentsUnion = DruidTalents | HunterTalents | MageTalents | RogueTalents | PaladinTalents | PriestTalents | ShamanTalents | WarlockTalents | WarriorTalents;
export declare type SpecTalents<T extends Spec> = T extends Spec.SpecBalanceDruid ? DruidTalents : T extends Spec.SpecElementalShaman ? ShamanTalents : T extends Spec.SpecHunter ? HunterTalents : T extends Spec.SpecMage ? MageTalents : T extends Spec.SpecRogue ? RogueTalents : T extends Spec.SpecRetributionPaladin ? PaladinTalents : T extends Spec.SpecShadowPriest ? PriestTalents : T extends Spec.SpecWarlock ? WarlockTalents : T extends Spec.SpecWarrior ? WarriorTalents : ShamanTalents;
export declare type SpecOptionsUnion = BalanceDruidOptions | ElementalShamanOptions | HunterOptions | MageOptions | RogueOptions | RetributionPaladinOptions | ShadowPriestOptions | WarlockOptions | WarriorOptions;
export declare type SpecOptions<T extends Spec> = T extends Spec.SpecBalanceDruid ? BalanceDruidOptions : T extends Spec.SpecElementalShaman ? ElementalShamanOptions : T extends Spec.SpecHunter ? HunterOptions : T extends Spec.SpecMage ? MageOptions : T extends Spec.SpecRogue ? RogueOptions : T extends Spec.SpecRetributionPaladin ? RetributionPaladinOptions : T extends Spec.SpecShadowPriest ? ShadowPriestOptions : T extends Spec.SpecWarlock ? WarlockOptions : T extends Spec.SpecWarrior ? WarriorOptions : ElementalShamanOptions;
export declare type SpecProtoUnion = BalanceDruid | ElementalShaman | Hunter | Mage | Rogue | RetributionPaladin | ShadowPriest | Warlock | Warrior;
export declare type SpecProto<T extends Spec> = T extends Spec.SpecBalanceDruid ? BalanceDruid : T extends Spec.SpecElementalShaman ? ElementalShaman : T extends Spec.SpecHunter ? Hunter : T extends Spec.SpecMage ? Mage : T extends Spec.SpecRogue ? Rogue : T extends Spec.SpecRetributionPaladin ? RetributionPaladin : T extends Spec.SpecShadowPriest ? ShadowPriest : T extends Spec.SpecWarlock ? Warlock : T extends Spec.SpecWarrior ? Warrior : ElementalShaman;
export declare type SpecTypeFunctions<SpecType extends Spec> = {
    rotationCreate: () => SpecRotation<SpecType>;
    rotationEquals: (a: SpecRotation<SpecType>, b: SpecRotation<SpecType>) => boolean;
    rotationCopy: (a: SpecRotation<SpecType>) => SpecRotation<SpecType>;
    rotationToJson: (a: SpecRotation<SpecType>) => any;
    rotationFromJson: (obj: any) => SpecRotation<SpecType>;
    talentsCreate: () => SpecTalents<SpecType>;
    talentsEquals: (a: SpecTalents<SpecType>, b: SpecTalents<SpecType>) => boolean;
    talentsCopy: (a: SpecTalents<SpecType>) => SpecTalents<SpecType>;
    talentsToJson: (a: SpecTalents<SpecType>) => any;
    talentsFromJson: (obj: any) => SpecTalents<SpecType>;
    optionsCreate: () => SpecOptions<SpecType>;
    optionsEquals: (a: SpecOptions<SpecType>, b: SpecOptions<SpecType>) => boolean;
    optionsCopy: (a: SpecOptions<SpecType>) => SpecOptions<SpecType>;
    optionsToJson: (a: SpecOptions<SpecType>) => any;
    optionsFromJson: (obj: any) => SpecOptions<SpecType>;
};
export declare const specTypeFunctions: Partial<Record<Spec, SpecTypeFunctions<any>>>;
export declare enum Faction {
    Unknown = 0,
    Alliance = 1,
    Horde = 2
}
export declare const raceToFaction: Record<Race, Faction>;
export declare const specToClass: Record<Spec, Class>;
export declare const specToEligibleRaces: Record<Spec, Array<Race>>;
export declare const specToEligibleItemCategories: Record<Spec, Array<ItemCategory>>;
export declare const specToLocalStorageKey: Record<Spec, string>;
export declare function withSpecProto<SpecType extends Spec>(player: Player, rotation: SpecRotation<SpecType>, talents: SpecTalents<SpecType>, specOptions: SpecOptions<SpecType>): Player;
export declare function playerToSpec(player: Player): Spec;
export declare function getMetaGemEffectEP(spec: Spec, gem: Gem, playerStats: Stats): number;
export declare function canEquipItem(item: Item, spec: Spec): boolean;
export declare function getEligibleItemSlots(item: Item): Array<ItemSlot>;
export declare function validWeaponCombo(mainHand: Item | null | undefined, offHand: Item | null | undefined): boolean;
export declare function getEligibleEnchantSlots(enchant: Enchant): Array<ItemSlot>;
export declare function enchantAppliesToItem(enchant: Enchant, item: Item): boolean;
export declare function gemMatchesSocket(gem: Gem, socketColor: GemColor): boolean;
export declare function gemEligibleForSocket(gem: Gem, socketColor: GemColor): boolean;
export declare const NO_TARGET = -1;
