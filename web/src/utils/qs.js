module.exports = () => {
  return window.Qs.parse(location.search.replace(/^\?/, ''))
}
