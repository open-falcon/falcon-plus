import fetch from './fetch'

module.exports = () => {
  return fetch({
    url: '/user/auth_session',
  })
  .then((res) => {
    return !res.data.error
  })
  .catch((err) => {
    return false
  })
}
