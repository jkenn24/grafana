import _ from 'lodash';
import { Variable, assignModelProperties, variableTypes } from './variable';

export class CustomVariable implements Variable {
  query: string;
  options: any;
  includeAll: boolean;
  multi: boolean;
  current: any;
  skipUrlSync: boolean;

  defaults = {
    type: 'custom',
    name: '',
    label: '',
    hide: 0,
    options: [],
    current: {},
    query: '',
    includeAll: false,
    multi: false,
    allValue: null,
    skipUrlSync: false,
  };

  /** @ngInject */
  constructor(private model, private variableSrv) {
    assignModelProperties(this, model, this.defaults);
  }

  setValue(option) {
    return this.variableSrv.setOptionAsCurrent(this, option);
  }

  getSaveModel() {
    assignModelProperties(this.model, this, this.defaults);
    return this.model;
  }

  updateOptions() {
    // extract options in comma separated string
    this.options = [];

    try {
      //attempt to carse options as JSON
      const arr = JSON.parse(this.query);
      for (let i = 0; i < arr.length; i++) {
        const obj = arr[i];
        let txt, val;
        for (const key in obj) {
          if (key === '__text') {
            txt = obj[key];
          } else if (key === '__value') {
            val = obj[key];
          }
        }
        if (txt && val) {
          this.options.push({ text: txt.trim(), value: val.trim() });
        }
      }
    } catch (e) {
      //If JSON parsing fails, use as comma separated string.
      this.options = [];
      this.options = _.map(this.query.split(/[,]+/), text => {
        return { text: text.trim(), value: text.trim() };
      });
    }

    if (this.includeAll) {
      this.addAllOption();
    }

    return this.variableSrv.validateVariableSelectionState(this);
  }

  addAllOption() {
    this.options.unshift({ text: 'All', value: '$__all' });
  }

  dependsOn(variable) {
    return false;
  }

  setValueFromUrl(urlValue) {
    return this.variableSrv.setOptionFromUrl(this, urlValue);
  }

  getValueForUrl() {
    if (this.current.text === 'All') {
      return 'All';
    }
    return this.current.value;
  }
}

variableTypes['custom'] = {
  name: 'Custom',
  ctor: CustomVariable,
  description: 'Define variable values manually',
  supportsMulti: true,
};
