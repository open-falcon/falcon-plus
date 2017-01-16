import fetch from '~utils/fetch'

module.exports = {
  'get/profile'(o = {}) {
    const opts = {
      method: 'GET',
      url: 'user/current'
    }

    return fetch(opts)
  },
  'put/user'(o = {}) {
    const opts = {
      method: 'PUT',
      url: 'user/update',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'put/cgpwd'(o = {}) {
    const opts = {
      method: 'PUT',
      url: 'user/cgpasswd',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'get/logout'() {
    const opts = {
      method: 'GET',
      url: 'user/logout'
    }

    return fetch(opts)
  }
}
