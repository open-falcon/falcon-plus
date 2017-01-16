import vfetch from '~utils/vuex-fetch'
const { _, moment } = window

module.exports = {
  'getEndpoints'({ commit, state }, { q }) {
    const opts = {
      url: 'graph/endpoint',
      params: {
        limit: 50,
        q,
      },
      commit,
      mutation: 'getEndpoints',
    }

    return vfetch(opts)
  },

  'getCounters'({ commit, state }, { eid = '6,7', metricQuery = '.+' }) {
    const opts = {
      url: 'graph/endpoint_counter',
      params: {
        limit: 50,
        eid,
        metricQuery,
      },
      commit,
      mutation: 'getCounters',
    }

    if (Array.isArray(opts.params.eid)) {
      opts.params.eid = opts.params.eid.join(',')
    }

    return vfetch(opts)
      .then((res) => {
        const items = res.data.map((counter) => {
          return { counter }
        })
        commit('getCounters.items', items)
      })
  },

  'viewGraph'({ commit, state }, { start, end, counters, endpoints, idx, vport, page, sampling }) {
    const opts = {
      method: 'POST',
      url: 'graph/history',
      data: {
        start_time: start || moment().add(-24, 'hours').unix(),
        end_time: end || moment().unix(),
        hostnames: endpoints,
        counters,
        step: 60,
        consol_fun: sampling || 'AVERAGE', // MAX, MIN
      },
    }

    return vfetch(opts)
      .then((res) => {
        const series = res.data.map((line) => {
          if (!line) {
            return {}
          }
          return {
            name: vport === 'combo' ? `${line.endpoint}${line.counter}` : (line[vport === 'endpoint' ? 'counter' : 'endpoint']),
            data: line.Values.map((Value) => {
              return [Value.timestamp, Value.value]
            }),
          }
        })
        commit('viewGraph.success', { idx, series, page })
      })
  },
}
