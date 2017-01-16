const userGroup = {
  namespaced: true,
  state: {
    rows: [],
    userListRows: [],
    notification: '',
    status: false,
    users: [],
    singleTeamUsers: [],
    singleTeamUsersToSelect: [],
    getSingleTeamLoading: false,
  },
  actions: require('./actions'),
  mutations: require('./mutations'),
}
module.exports = userGroup
