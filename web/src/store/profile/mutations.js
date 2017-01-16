module.exports = {
  'getUserInfo.success'(state, { data }) {
    state.userInfo = { ...data }
  },
  'getUserInfo.fail'(state, { err }) {},
  'updateProfile.success'(state, { data }) {
    state.notification = data.message
    state.status = true
  },
  'updateProfile.fail'(state, { err }) {
    state.notification = 'Error'
    state.status = false
  },
  'updatePwd.success'(state, { data }) {
    state.notification = data.message
    state.status = true
  },
  'updatePwd.fail'(state, { err }) {
    state.notification = 'Error'
    state.status = false
  },
  'logout.success'(state, { data, router }) {
    window.Cookies.remove('name')
    window.Cookies.remove('sig')
    router.push('/login')
  },
  'logout.fail'(state, { err }) {},
  'hideNotification'(state, { err }) {
    state.notification = ''
    state.status = false
  }
}
