// Vuex state module

const tpl = {
  state: {
    rows: [],
    simpleTList: [],
    createRepose: '',
  },
  actions: require('./actions'),
  mutations: require('./mutations'),
}

module.exports = tpl
