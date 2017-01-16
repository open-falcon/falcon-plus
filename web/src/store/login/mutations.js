module.exports = {
  'login.success'(state, { data, router }) {
    window.Cookies.set('name', data.name)
    window.Cookies.set('sig', data.sig)
    state.notification  = 'Log in Success!'
    state.status = true
    window.setTimeout(() => {
      router.push('/graph')
    }, 2000)
  },
  'login.fail'(state, { err }) {
    state.notification = 'Username or Password error'
    state.status = false
  }
}
