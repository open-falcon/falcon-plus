const signup = {
  namespaced: true,
  state: {
    userInfo: {
      name: '',
      email: '',
      password: '',
      cnname: '',
      im: '',
      phone: '',
      qq: ''
    },
    notification: '',
    status: false
  },
  actions: require('./actions'),
  mutations: require('./mutations'),
}
module.exports = signup
