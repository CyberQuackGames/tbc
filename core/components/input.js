import { Component } from './component.js';
// Shared logic for UI elements that are mapped to a value for some modifiable object.
export class Input extends Component {
    constructor(parent, cssClass, modObject, config) {
        super(parent, 'input-root');
        this.inputConfig = config;
        this.modObject = modObject;
        this.rootElem.classList.add(cssClass);
        if (config.label) {
            const label = document.createElement('span');
            label.classList.add('input-label');
            label.textContent = config.label;
            this.rootElem.appendChild(label);
            if (config.labelTooltip) {
                tippy(label, {
                    'content': config.labelTooltip,
                    'allowHTML': true,
                });
            }
        }
        config.changedEvent(this.modObject).on(() => {
            this.setInputValue(config.getValue(this.modObject));
            this.update();
        });
    }
    update() {
        const enable = !this.inputConfig.enableWhen || this.inputConfig.enableWhen(this.modObject);
        if (enable) {
            this.rootElem.classList.remove('disabled');
            this.getInputElem().removeAttribute('disabled');
        }
        else {
            this.rootElem.classList.add('disabled');
            this.getInputElem().setAttribute('disabled', '');
        }
    }
    // Can't call abstract functions in constructor, so need an init() call.
    init() {
        if (this.inputConfig.defaultValue) {
            this.setInputValue(this.inputConfig.defaultValue);
        }
        else {
            this.setInputValue(this.inputConfig.getValue(this.modObject));
        }
        this.update();
    }
    // Child classes should call this method when the value in the input element changes.
    inputChanged() {
        this.inputConfig.setValue(this.modObject, this.getInputValue());
    }
}
