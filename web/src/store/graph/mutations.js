module.exports = {
  'switchViewPoint'(state, { viewpoint  }) {
    state.vport = viewpoint
  },

  'getEndpoints.start'(state) {
    state.hasEndpointLoading = true
  },

  'getEndpoints.end'(state) {
    state.hasEndpointLoading = false
  },

  'getEndpoints.success'(state, { data }) {
    state.endpointItems = data
  },

  'getEndpoints.fail'(state) {

  },

  'getCounters.start'(state) {
    state.hasCounterLoading = true
  },

  'getCounters.end'(state) {
    state.hasCounterLoading = false
  },

  'getCounters.items'(state, counters) {
    state.counterItems = counters
  },

  'viewGraph.start'(state, { totalCharts }) {
    state.totalCharts = totalCharts
    state._cachePages = {}
  },

  'viewGraph.page'(state, { charts, page }) {
    state.pageCurrent = page
    state.charts = charts
  },

  'viewGraph.success'(state, { series, idx, page }) {
    state.charts[idx].loading = false
    state.charts[idx].series = series

    if (!state._cachePages[page]) {
      state._cachePages[page] = []
    }

    state._cachePages[page] = state.charts
  },

  'switchGrid'(state, { grid }) {
    state.grid = grid
  },

  'checkViewGraphBtnStatus'(state, { disabled }) {
    state.viewGraphBtnDisabled = disabled
  },

  'syncStartTime'(state, { unix }) {
    state.startTime = unix
  },

  'syncEndTime'(state, { unix }) {
    state.endTime = unix
  }

}
