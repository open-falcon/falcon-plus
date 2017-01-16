import fetch from '~utils/fetch'

module.exports = {
  'get/aggregators'(o = {}) {
    const opts = {
      method: 'GET',
      url: `hostgroup/${o.hostgroup}/aggregators`
      // url: 'hostgroup'
    }

    return fetch(opts)
  },
  'add/aggregators'(o ={}) {
    const opts = {
      method: 'POST',
      url: '/aggregator',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'edit/aggregators'(o ={}) {
    const opts = {
      method: 'PUT',
      url: '/aggregator',
      data: {
        ...o
      }
    }

    return fetch(opts)
  },
  'delete/aggregators'(o ={}) {
    const opts = {
      method: 'DELETE',
      url: `/aggregator/${o.id}`,
    }

    return fetch(opts)
  }
}
