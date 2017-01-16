const fixUnit = (unit) => {
  if (!/^(string|number)$/.test(typeof unit)) {
    throw new Error(`${unit} is not a string or number.`)
  }

  if (typeof unit === 'string' && isNaN(parseFloat(unit))) {
    throw new Error(`${unit} is not a legal number-string.`)
  }

  unit = String(unit)

  if (/\d$/.test(unit)) {
    return `${unit}px`
  }
  return unit
}

module.exports = fixUnit
