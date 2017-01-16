import fetch from '~utils/fetch'

module.exports = {
  'getStrategy'({ commit, state }, id) {
    const opts = {
      url: `strategy/${id}`,
      method: 'GET',
      commit,
      mutation: 'getStrategy',
    }
    return fetch(opts)
      .then((res) => {
        commit('getStrategy', res.data)
      })
      .catch((err) => {
        commit('getStrategy', err)
      })
  },
  'getTemplate'({ commit, state }, id) {
    const opts = {
      url: `template/${id}`,
      method: 'GET',
      commit,
      mutation: 'getTemplate',
    }

    return fetch(opts)
      .then((res) => {
        const tpl = res.data.template
        const caction = res.data.action
        const strategys = res.data.stratges.map((stratge) => {
          let col1 = stratge.metric
          if (stratge.tags.match(/.+/g)) {
            col1 = `${col1}/${stratge.tags}`
          }
          col1 = `${col1} [${stratge.note}]`
          let col4 = 'all day'
          if (stratge.run !== '' && stratge.run_end !== '') {
            col4 = `${stratge.run_begin}~${stratge.run_end}`
          }
          return [
            {
              col: col1
            },
            {
              col: `${stratge.func}${stratge.op}${stratge.right_value}`
            },
            {
              col: stratge.max_step
            },
            {
              col: stratge.priority
            },
            {
              col: col4
            },
            {
              col: parseInt(stratge.id)
            }
          ]
        })
        const parentName = res.data.parent_name
        commit('getTemplate', { tpl, parentName, caction, strategys })
      })
  },
  'getMetric'({ commit, state }, {}) {
    const opts = {
      url: `metric/tmplist`,
      method: 'GET',
      commit,
      mutation: 'getMetric',
    }

    return fetch(opts)
      .then((res) => {
        const metricMap = res.data.map((m) =>{
          return { id: m, text: m }
        })
        commit('getMetric', metricMap)
      })
  },
  'newStrategy'({ commit, state, dispatch }, d) {
    const opts = {
      url: `strategy`,
      method: 'POST',
      commit,
      data: {
        ...d.data
      },
      mutation: 'newStrategy',
    }

    return fetch(opts)
    .then((res) => {
      dispatch('getTemplate', d.id)
      // commit('newStrategy', res.data)
    })
  },
  'updateStrategy'({ commit, state, dispatch }, d) {
    const opts = {
      url: `strategy`,
      method: 'PUT',
      commit,
      data: {
        ...d.data
      },
      mutation: 'updateStrategy',
    }

    return fetch(opts)
      .then((res) => {
        dispatch('getTemplate', d.id)
        // commit('updateStrategy', res.data)
      })
  },
  'getTeamList'({ commit, state, dispatch }, {}) {
    const opts = {
      url: `team`,
      method: 'GET',
      commit,
      mutation: 'getTeamList',
    }
    return fetch(opts)
      .then((res) => {
        const teams = res.data.map((t) => {
          return { id: t.Team.id.toString(), text: t.Team.name }
        })
        commit('getTeamList', teams)
      })
  },
  'updateTemplate'({ commit, state, dispatch }, d) {
    const opts = {
      url: `template`,
      method: 'PUT',
      data: {
        ...d.data
      },
      commit,
      mutation: 'getReponse',
    }
    return fetch(opts)
      .then((res) => {
        commit('getReponse', res.data)
        dispatch('getTemplate', d.id)
      })
  },
  'createAction'({ commit, state, dispatch }, d) {
    const opts = {
      url: `template/action`,
      method: 'POST',
      data: {
        ...d.data
      },
      commit,
      mutation: 'getReponse',
    }
    return fetch(opts)
      .then((res) => {
        commit('getReponse', res.data)
        dispatch('getTemplate', d.id)
      })
  },
  'updateAction'({ commit, state, dispatch }, d) {
    const opts = {
      url: `template/action`,
      method: 'PUT',
      data: {
        ...d.data
      },
      commit,
      mutation: 'getReponse',
    }
    return fetch(opts)
      .then((res) => {
        commit('getReponse', res.data)
        dispatch('getTemplate', d.id)
      })
  },
  'deleteStrategy'({ commit, state, dispatch }, d) {
    const opts = {
      url: `/strategy/${d.id}`,
      method: 'DELETE',
      commit,
    }

    return fetch(opts)
      .then((res) => {
        dispatch('getTemplate', d.tid)
      })
  }
}
