module.exports = {
  'createNewUser.success'(state, { data, router }) {
    state.notification = 'Sign up Success!'
    state.status = true
    window.setTimeout(() => {
      router.push('/login')
    }, 2000)
  },
  'createNewUser.fail'(state, { err }) {
    state.notification = 'Error!'
    state.status = false
  },
  'pwdNotMatch'(state, { data }) {
    state.notification = 'password not matched'
    state.status = false
  }
}
