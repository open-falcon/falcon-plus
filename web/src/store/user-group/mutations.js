module.exports = {
  'userGroups.success'(state, { data }) {
    state.rows = data.reduce((preVal, curVal) => {
      preVal.push({
        id: curVal.Team.id,
        groupName: curVal.Team.name,
        groupMember: curVal.Useres.reduce((preVal, curVal) => {
          return `${preVal} ${curVal.cnname}`
        }, ''),
        creator: curVal.creator_name
      })
      return preVal
    }, [])
    // state.notification = 'Sign up Success!'
    // state.status = true
  },
  'userGroups.fail'(state, { err }) {
    // state.notification = 'Error!'
    // state.status = false
  },
  'newUserGroup.success'(state, { data }) {
    state.notification = 'created!'
    state.status = true
  },
  'newUserGroup.fail'(state, { err }) {},
  'users.success'(state, { data }) {
    state.users = data
    state.userListRows = data
  },
  'users.fail'(state, { err }) {},
  'searchUsers.success'(state, { data }) {
    state.userListRows = data
  },
  'searchUsers.fail'(state, { err }) {},
  'editUserGroup.success'(state, { data }) {},
  'editUserGroup.fail'(state, { err }) {},
  'deleteUserGroup.success'(state, { data }) {},
  'deleteUserGroup.fail'(state, { err }) {},
  'singleTeam.start'(state) {
    state.getSingleTeamLoading = true
  },
  'singleTeam.end'(state) {
  },
  'singleTeam.success'(state, { data }) {
    state.singleTeamUsers = data.users.reduce((preVal, curVal) => {
      preVal.push({
        cnname: curVal.cnname,
        role: curVal.role,
        id: curVal.id,
        name: curVal.name,
      })
      return preVal
    }, [])
    const userIds = data.users.reduce((preVal, curVal) => {
      preVal.push(curVal.id)
      return preVal
    }, [])
    state.singleTeamUsersToSelect = state.users.reduce((preVal, curVal) => {
      if (userIds.indexOf(curVal.id) < 0) {
        preVal.push(curVal)
      }
      return preVal
    }, [])
    state.getSingleTeamLoading = false
  },
  'singleTeam.fail'(state, { err }) {},
  'searchUserGroup.success'(state, { data }) {
    if (!data) {
      state.rows = []
      return
    }
    state.rows = data.reduce((preVal, curVal) => {
      preVal.push({
        id: curVal.Team.id,
        groupName: curVal.Team.name,
        groupMember: curVal.Useres.reduce((preVal, curVal) => {
          return `${preVal} ${curVal.cnname}`
        }, ''),
        creator: curVal.Team.creator
      })
      return preVal
    }, [])
  },
  'searchUserGroup.fail'(state, { err }) {}
}
