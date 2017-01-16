import api from './api'

module.exports = {
  'login'({ commit, state }, n) {
    api['post/login']({
      ...state.userInfo,
      name: n.name,
      password: n.password,
    })
    .then((res) => {
      commit('login.success', {
        data: res.data,
        router: n.router
      })
    })
    .catch((err) => {
      commit('login.fail', {
        err
      })
    })
  }
}
