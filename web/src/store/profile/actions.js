import api from './api'
const { Cookies } = window

module.exports = {
  'getProfile'({ commit, state }) {
    api['get/profile']({
      ...state.userInfo
    })
    .then((res) => {
      commit('getUserInfo.success', {
        data: res.data,
      })
    })
    .catch((err) => {
      commit('getUserInfo.fail', {
        err
      })
    })
  },
  'updateProfile'({ commit, state, dispatch }, n) {
    api['put/user']({
      ...state.userInfo,
      name: Cookies.get('name'),
      cnname: n.cnname || state.userInfo.cnname,
      email: n.email || state.userInfo.email,
      phone: n.phone,
      im: n.im,
      qq: n.qq
    })
    .then((res) => {
      commit('updateProfile.success', {
        data: res.data,
      })
      dispatch('getProfile')
      window.setTimeout(() => {
        dispatch('hideNotification')
      }, 5000)
    })
    .catch((err) => {
      commit('updateProfile.fail', {
        err
      })
    })
  },
  'updatePwd'({ commit, state, dispatch }, n) {
    api['put/cgpwd']({
      ...n
    })
    .then((res) => {
      commit('updatePwd.success', {
        data: res.data
      })
      window.setTimeout(() => {
        dispatch('hideNotification')
      }, 5000)
    })
    .catch((err) => {
      commit('updatePwd.fail', {
        err
      })
    })
  },
  'logout'({ commit, state }, n) {
    api['get/logout']({})
    .then((res) => {
      commit('logout.success', {
        data: res.data,
        router: n.router
      })
    })
    .catch((err) => {
      commit('logout.fail', {
        err
      })
    })
  },
  'hideNotification'({ commit, state }) {
    commit('hideNotification', {})
  }
}
