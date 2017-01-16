import fetch from '~utils/fetch'

module.exports = {
  'get/userGroups'(o = {}) {
    const opts = {
      method: 'GET',
      url: 'team'
    }

    return fetch(opts)
  },
  'post/newUserGroup'(o = {}) {
    const opts = {
      method: 'POST',
      url: 'team',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'get/users'(o = {}) {
    const opts = {
      method: 'GET',
      url: 'user/users'
    }

    return fetch(opts)
  },
  'delete/userGroup'(o = {}) {
    const opts = {
      method: 'DELETE',
      url: `team/${o.id}`,
    }

    return fetch(opts)
  },
  'put/userGroup'(o = {}) {
    const opts = {
      method: 'PUT',
      url: 'team',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'get/singleTeam'(o = {}) {
    const opts = {
      method: 'GET',
      url: `team/${o.team_id}`
    }

    return fetch(opts)
  },
  'search/user'(o = {}) {
    const opts = {
      method: 'GET',
      url: 'user/users',
      params: {
        ...o
      }
    }

    return fetch(opts)
  },
  'search/group'(o = {}) {
    const opts = {
      method: 'GET',
      url: 'team',
      params: {
        ...o
      }
    }

    return fetch(opts)
  }
}
