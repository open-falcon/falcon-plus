import fetch from '~utils/fetch'

module.exports = {
  'post/login'(o = {}) {
    const opts = {
      method: 'POST',
      url: 'user/login',
      params: {
        ...o
      }
    }
    return fetch(opts)
  }
}
