const u = 1000 // or 1024

const K = u
const M = K * u
const G = M * u
const T = G * u

const uMap = {
  T, G, M, K
}

module.exports = (val, opts = {}) => {
  opts = {
    decimals: 2,
    space: false,
    ...opts
  }

  if (typeof val === 'number') {
    let unit
    if (val >= T) {
      unit = 'T'
    } else if (val >= G) {
      unit = 'G'
    } else if (val >= M) {
      unit = 'M'
    } else if (val >= K) {
      unit = 'K'
    }

    if (unit) {
      const space = opts.space ? ' ' : ''
      val = (val / uMap[unit]).toFixed(opts.decimals)
      return `${val}${space}${unit}`
    }

    return `${val}`
  } else if (typeof val === 'string') {
    const unit = val.substr(-1)

    if (uMap[unit]) {
      return parseFloat(val) * uMap[unit]
    } else {
      if (/^\d$/.test(unit)) {
        return +val
      }
      throw new Error('The units are not: K, M, G and T')
    }
  } else {
    throw new Error(`${val} is not a number or string`)
  }
}
