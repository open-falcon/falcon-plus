const aggregator = {
  state: {
    rows: [],
    currentHostGroupName: ''
  },
  actions: require('./actions'),
  mutations: require('./mutations')
}
module.exports = aggregator
