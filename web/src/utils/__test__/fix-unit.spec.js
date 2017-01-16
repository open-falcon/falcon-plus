import fixUnit from '../fix-unit'

it('if unit not a `string` or `number`', () => {
  expect(() => {
    fixUnit({})
  }).toThrow()

  expect(() => {
    fixUnit([])
  }).toThrow()
})

it('if not a legal number string', () => {
  expect(() => {
    fixUnit('')
  }).toThrow()

  expect(() => {
    fixUnit('x123')
  }).toThrow()
})

it('If have no unit, it should be return as `px`', () => {
  expect(fixUnit('87')).toBe('87px')
  expect(fixUnit(87)).toBe('87px')
})

it('If have set unit, it should be return as unit', () => {
  expect(fixUnit('87%')).toBe('87%')
  expect(fixUnit('87em')).toBe('87em')
})
