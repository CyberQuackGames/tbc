import { GemColor } from './proto/common.js';
import { ItemSlot } from './proto/common.js';
export declare const urlPathPrefix: string;
export declare function getEmptySlotIconUrl(slot: ItemSlot): string;
export declare type IconId = {
    itemId: number;
};
export declare type SpellId = {
    spellId: number;
};
export declare type ItemOrSpellId = IconId | SpellId;
export declare function getIconUrl(id: ItemOrSpellId): Promise<string>;
export declare function setWowheadHref(elem: HTMLAnchorElement, id: ItemOrSpellId): void;
export declare function getEmptyGemSocketIconUrl(color: GemColor): string;