const { moment } = window

module.exports = {
  state: {
    endpointItems: [],
    counterItems: [],
    hasEndpointLoading: false,
    hasCounterLoading: false,
    endpointQ: '',
    endpointCounterQ:                                                                                                                                             '',
    charts: [],
    totalCharts: [],
    vport: 'endpoint',
    pageLimit: 6,
    pageCurrent: 1,
    viewGraphBtnDisabled: true,
    grid: 2,
    sampling: 'AVERAGE',
    _cachePages: {},
    startTime: moment().add(-24, 'hours').unix(),
    endTime: moment().unix(),
  },
  namespaced: true,

  actions: require('./actions'),
  mutations: require('./mutations'),
}
