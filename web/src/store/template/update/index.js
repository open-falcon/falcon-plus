// Vuex state module

const atpl = {
  state: {
    name: { id: 0, name: 'N/A' },
    parent: { id: '0', name: 'no' },
    actionId: 0,
    uics: '',
    action: {
      id: 0,
      url: '',
      before_callback_sms: false,
      before_callback_mail: false,
      after_callback_sms: false,
      after_callback_mail: false,
    },
    metrics: [],
    strategys: [],
    ustrategy: {
      tags: '',
      run_end: '',
      run_begin: '',
      right_value: '',
      priority: 0,
      op: '',
      note: '',
      metric: '',
      max_step: 0,
      id: 0,
      func: ''
    },
    apiRepose: '',
    teamList: [],
  },
  actions: require('./actions'),
  mutations: require('./mutations'),
}

module.exports = atpl
