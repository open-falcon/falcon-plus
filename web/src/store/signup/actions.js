import api from './api'

module.exports = {
  'signup'({ commit, state }, n) {
    api['post/newuser']({
      ...state.userInfo,
      name: n.name,
      email: n.email,
      password: n.password,
      cnname: n.cnname
    })
    .then((res) => {
      commit('createNewUser.success', {
        data: res.data,
        router: n.router
      })
    })
    .catch((err) => {
      commit('createNewUser.fail', {
        err
      })
    })
  },
  'pwdNotMatch'({ commit, state }) {
    commit('pwdNotMatch', {})
  }
}
