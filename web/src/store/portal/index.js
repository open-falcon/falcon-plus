import fetch from '~utils/fetch'
import vfetch from '~utils/vuex-fetch'

module.exports = {
  namespaced: true,
  state: {
    hostGroups: [],
    hosts: [],
    hasEndpointLoading: true,
    searchHostGrooupLoading: false,
    getHostsListLoading: false,
    getBindPluginListLoading: false,
    searchHostGroupInput: '.+',
    hostGroupListItems: [],
    hostList: {
      hostGroup: '',
      hostListItems: [],
    },
    tempDeleteCandidate: {
      name: '',
      id: 0,
    },
    hostInGroupList: {},
    bindPluginCandidate: {
      groupId: 0,
    },
    pluginsList: [],
    templateData: {
      templates: [],
      hostgroup:{}
    },
    templates: [
      { value: '', title: '' }
    ]
  },
  actions: {
    'getEndpoints'({ commit, state }, { q }) {
      commit('getEndpoints.start')

      const opts = {
        url: 'graph/endpoint',
        params: {
          limit: 500,
          q,
        },
        commit,
        mutation: 'getEndpoints',
      }

      return vfetch(opts)
    },

    'createHostGroupName'({ commit, state, dispatch }, { name, hosts }) {
      const opts = {
        method: 'post',
        url: 'hostgroup',
        data: {
          name,
        },
      }

      fetch(opts)
        .then((res) => {
          const data = {
            id: res.data.id,
            hosts,
          }
          dispatch('addHostsIntoNewHostGroup', data)
        })
        .catch((err) => {
          console.error(err)
          reject(err)
        })
    },

    'addHostsIntoNewHostGroup'({ commit, state, dispatch }, { id, hosts }) {
      const opts = {
        method: 'post',
        url: 'hostgroup/host',
        data: {
          hostgroup_id: id,
          hosts,
        },
        commit,
        mutation: 'addHostsIntoNewHostGroup',
      }

      return vfetch(opts)
        .then((res) => {
          dispatch('getHostGroupList')
        })
    },

    'getHostGroupList'({ commit, state }) {
      const opts = {
        method: 'get',
        url: 'hostgroup',
        mutation : 'getHostGroupList',
        commit,
      }

      return vfetch(opts)
    },

    'searchHostGroup'({ commit, state }, q = '') {
      if (!q.length) {
        q = '.+'
      }

      const opts = {
        method: 'get',
        url: 'hostgroup',
        params: {
          q,
        },
        mutation : 'searchHostGroup',
        commit,
      }

      return vfetch(opts)
    },

    'getHostsList'({ commit, state }, data) {
      const opts = {
        method: 'get',
        url: `hostgroup/${data.groupId}`,
        mutation : 'getHostsList',
        commit,
      }

      return vfetch(opts)
        .then((res) => {
          commit('hostInGroupList', data)
        })
    },

    'deleteHostGroup'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'delete',
        url: `hostgroup/${data.id}`,
        mutation : 'deleteHostGroup',
        commit,
      }

      return vfetch(opts)
        .then((res) => {
          dispatch('getHostGroupList')
        })
    },

    'deleteHostFromGroup'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'put',
        url: 'hostgroup/host',
        mutation : 'deleteHostFromGroup',
        data: {
          hostgroup_id: +data.groupId,
          host_id: +data.hostId,
        },
        commit,
      }

      return vfetch(opts)
        .then((res) => {
          dispatch('getHostsList', data)
        })
    },

    // List Plugins
    'getBindPluginList'({ commit, state }, data) {
      const opts = {
        method: 'get',
        url: `hostgroup/${data.groupId}/plugins`,
        mutation : 'getBindPluginList',
        commit,
      }

      return vfetch(opts)
    },

    // Plugin binding
    'bindPluginToHostGroup'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'post',
        url: 'plugin',
        mutation : 'bindPluginToHostGroup',
        data: {
          hostgroup_id: data.groupId,
          dir_path: data.pluginDir,
        },
        commit,
      }

      return vfetch(opts)
        .then((res) => {
          dispatch('getBindPluginList', data)
        })
    },

    'unbindPluginFromGroup'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'delete',
        url: `plugin/${data.id}`,
        mutation : 'unbindPluginFromGroup',
        commit,
      }

      return vfetch(opts)
        .then((res) => {
          dispatch('getBindPluginList', data)
        })
    },

    'searchHostList'({ commit, state }, data) {
      const opts = {
        method: 'GET',
        url: `/hostgroup/${data.id}?q=${data.q}`
      }

      return fetch(opts)
              .then((res) => {
                commit('searchHostList.success', {
                  data: res.data
                })
              })
              .catch((err) => {
                commit('searchHostList.fail', {
                  err
                })
              })
    },

    'getBindTemplates'({ commit, state }, data) {
      const opts = {
        method: 'GET',
        url: `hostgroup/${data.groupId}/template`
      }

      return fetch(opts)
              .then((res) => {
                commit('getBindTemplates.success', {
                  data: res.data
                })
              })
              .catch((err) => {
                commit('getBindTemplates.fail', {
                  err
                })
              })
    },

    'getTemplates'({ commit, state }, data) {
      const opts = {
        method: 'GET',
        url: 'template'
      }

      return fetch(opts)
              .then((res) => {
                commit('getTemplates.success', {
                  data: res.data
                })
              })
              .catch((err) => {
                commit('getTemplates.fail', {
                  err
                })
              })
    },
    'bindOneTemplate'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'POST',
        url: 'hostgroup/template',
        data: {
          ...data
        }
      }

      return fetch(opts)
              .then((res) => {
                commit('bindOneTemplate.success', {
                  data: res.data
                })
                dispatch('getBindTemplates', {
                  groupId: data.grp_id
                })
              })
              .catch((err) => {
                commit('bindOneTemplate.fail', {
                  err
                })
              })
    },

    'unbindOneTemplate'({ commit, state, dispatch }, data) {
      const opts = {
        method: 'PUT',
        url: 'hostgroup/template',
        data: {
          ...data
        }
      }

      return fetch(opts)
              .then((res) => {
                commit('unbindOneTemplate.success', {
                  data: res.data
                })
                dispatch('getBindTemplates', {
                  groupId: data.grp_id
                })
              })
              .catch((err) => {
                commit('unbindOneTemplate.fail', {
                  err
                })
              })
    },
  },
  mutations: require('./mutations'),
}
