import { setWowheadHref } from '/tbc/core/resources.js';
import { parseAuraMetrics } from './metrics_helpers.js';
import { ResultComponent } from './result_component.js';
export class BuffAuraMetrics extends ResultComponent {
    constructor(config) {
        config.rootCssClass = 'buff-aura-metrics-root';
        super(config);
        this.rootElem.innerHTML = `
		<table class="aura-metrics-table tablesorter">
			<thead class="aura-metrics-table-header">
				<tr class="aura-metrics-table-header-row">
					<th class="aura-metrics-table-header-cell"><span>Name</span></th>
					<th class="aura-metrics-table-header-cell"><span>Uptime</span></th>
				</tr>
			</thead>
			<tbody class="aura-metrics-table-body">
			</tbody>
		</table>
		`;
        this.tableElem = this.rootElem.getElementsByClassName('aura-metrics-table')[0];
        this.bodyElem = this.rootElem.getElementsByClassName('aura-metrics-table-body')[0];
        const headerElems = Array.from(this.tableElem.querySelectorAll('th'));
        // Uptime
        tippy(headerElems[1], {
            'content': 'Uptime / Encounter Duration',
            'allowHTML': true,
        });
        $(this.tableElem).tablesorter({ sortList: [[1, 1]] });
    }
    onSimResult(request, result) {
        this.bodyElem.textContent = '';
        const iterations = request.simOptions.iterations;
        const duration = request.encounter?.duration || 1;
        parseAuraMetrics(result.raidMetrics.parties[0].players[0].auras).then(auraMetrics => {
            auraMetrics.forEach(auraMetric => {
                const rowElem = document.createElement('tr');
                this.bodyElem.appendChild(rowElem);
                const nameCellElem = document.createElement('td');
                rowElem.appendChild(nameCellElem);
                nameCellElem.innerHTML = `
				<a class="aura-metrics-action-icon"></a>
				<span class="aura-metrics-action-name">${auraMetric.name}</span>
				`;
                const iconElem = nameCellElem.getElementsByClassName('aura-metrics-action-icon')[0];
                iconElem.style.backgroundImage = `url('${auraMetric.iconUrl}')`;
                if (!('otherId' in auraMetric.actionId.id)) {
                    setWowheadHref(iconElem, auraMetric.actionId.id);
                }
                const addCell = (value) => {
                    const cellElem = document.createElement('td');
                    cellElem.textContent = String(value);
                    rowElem.appendChild(cellElem);
                    return cellElem;
                };
                addCell((auraMetric.uptimeSecondsAvg / duration * 100).toFixed(2) + '%'); // Uptime
            });
            $(this.tableElem).trigger('update');
        });
    }
}
export class DebuffAuraMetrics extends ResultComponent {
    constructor(config) {
        config.rootCssClass = 'debuff-aura-metrics-root';
        super(config);
        this.rootElem.innerHTML = `
		<table class="aura-metrics-table tablesorter">
			<thead class="aura-metrics-table-header">
				<tr class="aura-metrics-table-header-row">
					<th class="aura-metrics-table-header-cell"><span>Name</span></th>
					<th class="aura-metrics-table-header-cell"><span>Uptime</span></th>
				</tr>
			</thead>
			<tbody class="aura-metrics-table-body">
			</tbody>
		</table>
		`;
        this.tableElem = this.rootElem.getElementsByClassName('aura-metrics-table')[0];
        this.bodyElem = this.rootElem.getElementsByClassName('aura-metrics-table-body')[0];
        const headerElems = Array.from(this.tableElem.querySelectorAll('th'));
        // Uptime
        tippy(headerElems[1], {
            'content': 'Uptime / Encounter Duration',
            'allowHTML': true,
        });
        $(this.tableElem).tablesorter({ sortList: [[1, 1]] });
    }
    onSimResult(request, result) {
        this.bodyElem.textContent = '';
        const iterations = request.simOptions.iterations;
        const duration = request.encounter?.duration || 1;
        parseAuraMetrics(result.encounterMetrics.targets[0].auras).then(auraMetrics => {
            auraMetrics.forEach(auraMetric => {
                const rowElem = document.createElement('tr');
                this.bodyElem.appendChild(rowElem);
                const nameCellElem = document.createElement('td');
                rowElem.appendChild(nameCellElem);
                nameCellElem.innerHTML = `
				<a class="aura-metrics-action-icon"></a>
				<span class="aura-metrics-action-name">${auraMetric.name}</span>
				`;
                const iconElem = nameCellElem.getElementsByClassName('aura-metrics-action-icon')[0];
                iconElem.style.backgroundImage = `url('${auraMetric.iconUrl}')`;
                if (!('otherId' in auraMetric.actionId.id)) {
                    setWowheadHref(iconElem, auraMetric.actionId.id);
                }
                const addCell = (value) => {
                    const cellElem = document.createElement('td');
                    cellElem.textContent = String(value);
                    rowElem.appendChild(cellElem);
                    return cellElem;
                };
                addCell((auraMetric.uptimeSecondsAvg / duration * 100).toFixed(2) + '%'); // Uptime
            });
            $(this.tableElem).trigger('update');
        });
    }
}