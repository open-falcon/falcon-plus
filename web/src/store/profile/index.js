const profile = {
  namespaced: true,
  state: {
    userInfo: {},
    notification: '',
    status: false
  },
  actions: require('./actions'),
  mutations: require('./mutations'),
}
module.exports = profile
