import { EquippedItem } from '/tbc/core/proto_utils/equipped_item.js';
import { Item } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { enchantDescriptions } from '/tbc/core/proto_utils/names.js';
import { slotNames } from '/tbc/core/proto_utils/names.js';
import { getEmptyGemSocketIconUrl } from '/tbc/core/resources.js';
import { getEmptySlotIconUrl } from '/tbc/core/resources.js';
import { getIconUrl } from '/tbc/core/resources.js';
import { setWowheadHref } from '/tbc/core/resources.js';
import { setGemSocketCssClass } from '/tbc/core/css_utils.js';
import { setItemQualityCssClass } from '/tbc/core/css_utils.js';
import { Component } from './component.js';
export class GearPicker extends Component {
    constructor(parent, player) {
        super(parent, 'gear-picker-root');
        const leftSide = document.createElement('div');
        leftSide.classList.add('gear-picker-left');
        this.rootElem.appendChild(leftSide);
        const rightSide = document.createElement('div');
        rightSide.classList.add('gear-picker-right');
        this.rootElem.appendChild(rightSide);
        const selectorModal = new SelectorModal(document.body, player);
        const leftItemPickers = [
            ItemSlot.ItemSlotHead,
            ItemSlot.ItemSlotNeck,
            ItemSlot.ItemSlotShoulder,
            ItemSlot.ItemSlotBack,
            ItemSlot.ItemSlotChest,
            ItemSlot.ItemSlotWrist,
            ItemSlot.ItemSlotMainHand,
            ItemSlot.ItemSlotOffHand,
        ].map(slot => new ItemPicker(leftSide, player, slot, selectorModal));
        const rightItemPickers = [
            ItemSlot.ItemSlotHands,
            ItemSlot.ItemSlotWaist,
            ItemSlot.ItemSlotLegs,
            ItemSlot.ItemSlotFeet,
            ItemSlot.ItemSlotFinger1,
            ItemSlot.ItemSlotFinger2,
            ItemSlot.ItemSlotTrinket1,
            ItemSlot.ItemSlotTrinket2,
            ItemSlot.ItemSlotRanged,
        ].map(slot => new ItemPicker(rightSide, player, slot, selectorModal));
        this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);
    }
}
class ItemPicker extends Component {
    constructor(parent, player, slot, selectorModal) {
        super(parent, 'item-picker-root');
        // All items and enchants that are eligible for this slot
        this._items = [];
        this._enchants = [];
        this._equippedItem = null;
        this.slot = slot;
        this.player = player;
        this.rootElem.innerHTML = `
      <a class="item-picker-icon" target="_blank" data-toggle="modal" data-target="#selectorModal">
        <div class="item-picker-sockets-container">
        </div>
      </a>
      <div class="item-picker-labels-container">
        <span class="item-picker-name"></span>
        <span class="item-picker-enchant"></span>
      </div>
    `;
        this.iconElem = this.rootElem.getElementsByClassName('item-picker-icon')[0];
        this.nameElem = this.rootElem.getElementsByClassName('item-picker-name')[0];
        this.enchantElem = this.rootElem.getElementsByClassName('item-picker-enchant')[0];
        this.socketsContainerElem = this.rootElem.getElementsByClassName('item-picker-sockets-container')[0];
        this.item = null;
        player.sim.gearListEmitter.on(() => {
            this._items = this.player.getItems(this.slot);
            this._enchants = this.player.getEnchants(this.slot);
            this.iconElem.addEventListener('click', event => {
                selectorModal.setData(this.slot, this._equippedItem, this._items, this._enchants);
            });
        });
        player.gearChangeEmitter.on(() => {
            this.item = player.getEquippedItem(slot);
        });
    }
    set item(newItem) {
        // Clear everything first
        this.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
        this.iconElem.removeAttribute('data-wowhead');
        this.iconElem.removeAttribute('href');
        this.nameElem.textContent = slotNames[this.slot];
        setItemQualityCssClass(this.nameElem, null);
        this.enchantElem.textContent = '';
        //this.enchantElem.removeAttribute('data-wowhead');
        //this.enchantElem.removeAttribute('href');
        this.socketsContainerElem.innerHTML = '';
        if (newItem != null) {
            this.nameElem.textContent = newItem.item.name;
            setItemQualityCssClass(this.nameElem, newItem.item.quality);
            this.player.setWowheadData(newItem, this.iconElem);
            setWowheadHref(this.iconElem, { itemId: newItem.item.id });
            getIconUrl({ itemId: newItem.item.id }).then(url => {
                this.iconElem.style.backgroundImage = `url('${url}')`;
            });
            if (newItem.enchant) {
                this.enchantElem.textContent = enchantDescriptions.get(newItem.enchant.id) || newItem.enchant.name;
                //this.enchantElem.setAttribute('href', 'https://tbc.wowhead.com/item=' + newItem.enchant.id);
            }
            newItem.item.gemSockets.forEach((socketColor, gemIdx) => {
                const gemIconElem = document.createElement('img');
                gemIconElem.classList.add('item-picker-gem-icon');
                setGemSocketCssClass(gemIconElem, socketColor);
                if (newItem.gems[gemIdx] == null) {
                    gemIconElem.src = getEmptyGemSocketIconUrl(socketColor);
                }
                else {
                    getIconUrl({ itemId: newItem.gems[gemIdx].id }).then(url => {
                        gemIconElem.src = url;
                    });
                }
                this.socketsContainerElem.appendChild(gemIconElem);
            });
        }
        this._equippedItem = newItem;
    }
}
class SelectorModal extends Component {
    constructor(parent, player) {
        super(parent, 'selector-model-root');
        this.player = player;
        this.rootElem.innerHTML = `
    <div class="modal fade selector-modal" id="selectorModal" tabindex="-1" role="dialog" aria-labelledby="selectorModalTitle" aria-hidden="true">
      <div class="modal-dialog" role="document">
        <div class="modal-content">
          <div class="modal-body">
            <ul class="nav nav-tabs selector-modal-tabs">
            </ul>
            <div class="tab-content selector-modal-tab-content">
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary selector-modal-close-button" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
    `;
        this.tabsElem = this.rootElem.getElementsByClassName('selector-modal-tabs')[0];
        this.contentElem = this.rootElem.getElementsByClassName('selector-modal-tab-content')[0];
        this.closeButton = this.rootElem.getElementsByClassName('selector-modal-close-button')[0];
    }
    setData(slot, equippedItem, eligibleItems, eligibleEnchants) {
        this.tabsElem.innerHTML = '';
        this.contentElem.innerHTML = '';
        this.tabsElem.innerHTML = `
    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
      <span aria-hidden="true">&times;</span>
    </button>
    `;
        this.addTab('Items', slot, equippedItem, eligibleItems, item => this.player.computeItemEP(item), equippedItem => equippedItem?.item, item => {
            return {
                id: item.id,
                name: item.name,
                quality: item.quality,
                phase: item.phase,
                onEquip: item => {
                    const equippedItem = this.player.getEquippedItem(slot);
                    if (equippedItem) {
                        this.player.equipItem(slot, equippedItem.withItem(item));
                    }
                    else {
                        this.player.equipItem(slot, new EquippedItem(item));
                    }
                },
            };
        }, () => {
            this.player.equipItem(slot, null);
        });
        this.addTab('Enchants', slot, equippedItem, eligibleEnchants, enchant => this.player.computeEnchantEP(enchant), equippedItem => equippedItem?.enchant, enchant => {
            return {
                id: enchant.id,
                name: enchant.name,
                quality: enchant.quality,
                phase: 1,
                onEquip: enchant => {
                    const equippedItem = this.player.getEquippedItem(slot);
                    if (equippedItem)
                        this.player.equipItem(slot, equippedItem.withEnchant(enchant));
                },
            };
        }, () => {
            const equippedItem = this.player.getEquippedItem(slot);
            if (equippedItem)
                this.player.equipItem(slot, equippedItem.withEnchant(null));
        });
        this.addGemTabs(slot, equippedItem);
    }
    addGemTabs(slot, equippedItem) {
        equippedItem?.item.gemSockets.forEach((socketColor, socketIdx) => {
            this.addTab('Gem ' + (socketIdx + 1), slot, equippedItem, this.player.getGems(socketColor), gem => this.player.computeGemEP(gem), equippedItem => equippedItem?.gems[socketIdx], gem => {
                return {
                    id: gem.id,
                    name: gem.name,
                    quality: gem.quality,
                    phase: gem.phase,
                    onEquip: gem => {
                        const equippedItem = this.player.getEquippedItem(slot);
                        if (equippedItem)
                            this.player.equipItem(slot, equippedItem.withGem(gem, socketIdx));
                    },
                };
            }, () => {
                const equippedItem = this.player.getEquippedItem(slot);
                if (equippedItem)
                    this.player.equipItem(slot, equippedItem.withGem(null, socketIdx));
            });
        });
    }
    /**
     * Adds one of the tabs for the item selector menu.
     *
     * T is expected to be Item, Enchant, or Gem. Tab menus for all 3 looks extremely
     * similar so this function uses extra functions to do it generically.
     */
    addTab(label, slot, equippedItem, items, computeEP, equippedToItemFn, getItemData, onRemove) {
        if (items.length == 0) {
            return;
        }
        const equippedToIdFn = (equippedItem) => {
            const item = equippedToItemFn(equippedItem);
            return item ? getItemData(item).id : 0;
        };
        items.sort((itemA, itemB) => computeEP(itemB) - computeEP(itemA));
        const tabElem = document.createElement('li');
        this.tabsElem.insertBefore(tabElem, this.tabsElem.lastChild);
        const tabContentId = (label + '-tab').split(' ').join('');
        tabElem.innerHTML = `<a class="selector-modal-item-tab" data-toggle="tab" href="#${tabContentId}">${label}</a>`;
        const tabContent = document.createElement('div');
        tabContent.id = tabContentId;
        tabContent.classList.add('tab-pane', 'fade', 'selector-modal-tab-content');
        this.contentElem.appendChild(tabContent);
        tabContent.innerHTML = `
    <div class="selector-modal-tab-content-header">
      <button class="selector-modal-remove-button">Remove</button>
      <input class="selector-modal-search" type="text" placeholder="Search...">
      <select class="selector-modal-phase-select">
				<option value="1">Phase 1</option>
				<option value="2">Phase 2</option>
				<option value="3">Phase 3</option>
				<option value="4">Phase 4</option>
				<option value="5">Phase 5</option>
			</select>
    </div>
    <ul class="selector-modal-list"></ul>
    `;
        if (label == 'Items') {
            tabElem.classList.add('active', 'in');
            tabContent.classList.add('active', 'in');
        }
        const listElem = tabContent.getElementsByClassName('selector-modal-list')[0];
        const listItemElems = items.map(item => {
            const itemData = getItemData(item);
            const itemEP = computeEP(item);
            const listItemElem = document.createElement('li');
            listItemElem.classList.add('selector-modal-list-item');
            listElem.appendChild(listItemElem);
            listItemElem.dataset.id = String(itemData.id);
            listItemElem.dataset.name = itemData.name;
            listItemElem.dataset.phase = String(Math.max(itemData.phase, 1));
            listItemElem.innerHTML = `
        <a class="selector-modal-list-item-icon"></a>
        <a class="selector-modal-list-item-name">${itemData.name}</a>
        <div class="selector-modal-list-item-padding"></div>
        <div class="selector-modal-list-item-ep">
					<span class="selector-modal-list-item-ep-value">${Math.round(itemEP)}</span>
					<span class="selector-modal-list-item-ep-delta"></span>
				</div>
      `;
            setWowheadHref(listItemElem.children[0], { itemId: itemData.id });
            setWowheadHref(listItemElem.children[1], { itemId: itemData.id });
            const iconElem = listItemElem.getElementsByClassName('selector-modal-list-item-icon')[0];
            getIconUrl({ itemId: itemData.id }).then(url => {
                iconElem.style.backgroundImage = `url('${url}')`;
            });
            const nameElem = listItemElem.getElementsByClassName('selector-modal-list-item-name')[0];
            setItemQualityCssClass(nameElem, itemData.quality);
            const onclick = (event) => {
                event.preventDefault();
                itemData.onEquip(item);
                // If the item changes, the gem slots might change, so remove and recreate the gem tabs
                if (Item.is(item)) {
                    this.removeTabs('Gem');
                    this.addGemTabs(slot, this.player.getEquippedItem(slot));
                }
            };
            nameElem.addEventListener('click', onclick);
            iconElem.addEventListener('click', onclick);
            return listItemElem;
        });
        const removeButton = tabContent.getElementsByClassName('selector-modal-remove-button')[0];
        removeButton.addEventListener('click', event => {
            listItemElems.forEach(elem => elem.classList.remove('active'));
            onRemove();
        });
        const updateSelected = () => {
            const newEquippedItem = this.player.getEquippedItem(slot);
            const newItem = equippedToItemFn(newEquippedItem);
            if (!newItem)
                return;
            const newItemId = equippedToIdFn(newEquippedItem);
            const newEP = computeEP(newItem);
            listItemElems.forEach(elem => {
                const listItemId = parseInt(elem.dataset.id);
                const listItem = items.find(item => getItemData(item).id == listItemId);
                elem.classList.remove('active');
                if (listItemId == newItemId) {
                    elem.classList.add('active');
                }
                const epDeltaElem = elem.getElementsByClassName('selector-modal-list-item-ep-delta')[0];
                epDeltaElem.textContent = '';
                if (listItem) {
                    const listItemEP = computeEP(listItem);
                    const delta = Math.round(listItemEP) - Math.round(newEP);
                    if (delta > 0) {
                        epDeltaElem.textContent = '+' + delta;
                        epDeltaElem.classList.remove('negative');
                        epDeltaElem.classList.add('positive');
                    }
                    else if (delta < 0) {
                        epDeltaElem.textContent = '' + delta;
                        epDeltaElem.classList.remove('positive');
                        epDeltaElem.classList.add('negative');
                    }
                }
            });
        };
        updateSelected();
        this.player.gearChangeEmitter.on(updateSelected);
        const applyFilters = () => {
            const searchQuery = searchInput.value.toLowerCase();
            const phase = this.player.sim.getPhase();
            listItemElems.forEach(elem => {
                if (elem.dataset.name.toLowerCase().includes(searchQuery) && Number(elem.dataset.phase) <= phase) {
                    elem.style.display = 'flex';
                }
                else {
                    elem.style.display = 'none';
                }
            });
        };
        const searchInput = tabContent.getElementsByClassName('selector-modal-search')[0];
        searchInput.addEventListener('input', applyFilters);
        const phaseSelect = tabContent.getElementsByClassName('selector-modal-phase-select')[0];
        phaseSelect.value = String(this.player.sim.getPhase());
        tabContent.dataset.phase = String(this.player.sim.getPhase());
        phaseSelect.addEventListener('input', event => {
            this.player.sim.setPhase(Number(phaseSelect.value));
        });
        this.player.sim.phaseChangeEmitter.on(() => {
            tabContent.dataset.phase = String(this.player.sim.getPhase());
            phaseSelect.value = String(this.player.sim.getPhase());
            applyFilters();
        });
        applyFilters();
    }
    removeTabs(labelSubstring) {
        const tabElems = Array.prototype.slice.call(this.tabsElem.getElementsByClassName('selector-modal-item-tab')).filter(tab => tab.textContent.includes(labelSubstring));
        const contentElems = tabElems
            .map(tabElem => document.getElementById(tabElem.href.substring(1)))
            .filter(tabElem => Boolean(tabElem));
        tabElems.forEach(elem => elem.remove());
        contentElems.forEach(elem => elem.remove());
    }
}