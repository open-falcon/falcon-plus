import api from './api'

module.exports = {
  'getAggregators'({ commit, state }, n) {
    api['get/aggregators']({
      ...n
    })
    .then((res) => {
      commit('get/aggregators.success', {
        data: res.data
      })
    })
    .catch((err) => {
      commit('get/aggregators.fail', {
        err
      })
    })
  },
  'newAggregator'({ commit, state, dispatch }, n) {
    api['add/aggregators']({
      ...n
    })
    .then((res) => {
      commit('add/aggregators.success', {
        data: res.data
      })
      dispatch('getAggregators', { hostgroup: n.hostgroup_id })
    })
    .catch((err) => {
      commit('add/aggregators.fail', {
        err
      })
    })
  },
  'editAggregator'({ commit, state, dispatch }, n) {
    api['edit/aggregators']({
      ...n
    })
    .then((res) => {
      commit('edit/aggregators.success', {
        data: res.data
      })
      dispatch('getAggregators', { hostgroup: n.hostgroup_id })
    })
    .catch((err) => {
      commit('edit/aggregators.fail', {
        err
      })
    })
  },
  'deleteAggregator'({ commit, state, dispatch }, n) {
    api['delete/aggregators']({
      ...n
    })
    .then((res) => {
      commit('delete/aggregators.success', {
        data: res.data
      })
      dispatch('getAggregators', { hostgroup: n.hostgroup_id })
    })
    .catch((err) => {
      commit('delete/aggregators.fail', {
        err
      })
    })
  }
}
