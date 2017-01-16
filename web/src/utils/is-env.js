const isBrowser = typeof window !== 'undefined' && window.document && document.createElement
const isNode = !isBrowser && typeof global !== 'undefined'

module.exports = {
  isBrowser, isNode
}
