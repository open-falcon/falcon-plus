import api from './api'
import vfetch from '~utils/vuex-fetch'

module.exports = {
  'getUserGroups'({ commit, state }, n) {
    api['get/userGroups']({
      ...state.rows,
    })
    .then((res) => {
      commit('userGroups.success', {
        data: res.data,
      })
    })
    .catch((err) => {
      commit('userGroups.fail', {
        err
      })
    })
  },
  'newUserGroup'({ commit, state, dispatch }, n) {
    api['post/newUserGroup']({
      ...n
    })
    .then((res) => {
      commit('newUserGroup.success', {
        data: res.data
      })
      dispatch('getUserGroups')
    })
    .catch((err) => {
      commit('newUserGroup.fail', {
        err
      })
    })
  },
  'deleteUserGroup'({ commit, state, dispatch }, n) {
    api['delete/userGroup']({
      ...n
    })
    .then((res) => {
      commit('deleteUserGroup.success', {
        data: res.data
      })
      dispatch('getUserGroups')
    })
    .catch((err) => {
      commit('deleteUserGroup.fail', {
        err
      })
    })
  },
  'editUserGroup'({ commit, state, dispatch }, n) {
    api['put/userGroup']({
      ...n
    })
    .then((res) => {
      commit('editUserGroup.success', {
        data: res.data
      })
      dispatch('getUserGroups')
    })
    .catch((err) => {
      commit('editUserGroup.fail', {
        err
      })
    })
  },
  'getOneTeam'({ commit, state }, n) {
    const opts = {
      url: `http://113.207.30.198:8088/api/v1/team/${n.team_id}`,
      commit,
      mutation: 'singleTeam'
    }
    return vfetch(opts)
  },
  'searchGroup'({ commit, state }, n) {
    api['search/group']({
      ...n
    })
    .then((res) => {
      commit('searchUserGroup.success', {
        data: res.data
      })
    })
    .catch((err) => {
      commit('searchUserGroup.fail', {
        err
      })
    })
  },
  'getUsers'({ commit, state }, n) {
    api['get/users']({
      ...state.users,
    })
    .then((res) => {
      commit('users.success', {
        data: res.data
      })
    })
    .catch((err) => {
      commit('users.fail', {
        err
      })
    })
  },
  'searchUser'({ commit, state }, n) {
    api['search/user']({
      ...n
    })
    .then((res) => {
      commit('searchUsers.success', {
        data: res.data
      })
    })
    .catch((err) => {
      commit('searchUsers.fail', {
        err
      })
    })
  },
}
