import fetch from '~utils/fetch'

module.exports = {
  'post/newuser'(o = {}) {
    const opts = {
      method: 'POST',
      url: 'user/create',
      data: {
        ...o
      }
    }

    return fetch(opts)
  }
}
