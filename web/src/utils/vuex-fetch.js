import fetch from './fetch'

module.exports = (opts = {}) => {
  const commit = (status, arg) => {
    const hasMutation = opts.commit && opts.mutation
    if (hasMutation) {
      return opts.commit(`${opts.mutation}.${status}`, arg)
    }
  }

  return new Promise((resolve, reject) => {
    commit('start')
    fetch(opts)
      .then((res) => {
        commit('success', res)
        commit('end', res)
        resolve(res)
      })
      .catch((err) => {
        commit('fail', err)
        commit('end', err)
        reject(err)
      })
  })
}
