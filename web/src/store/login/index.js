const login = {
  namespaced: true,
  state: {
    notification: '',
    status: false
  },

  actions: require('./actions'),
  mutations: require('./mutations'),
}
module.exports = login
