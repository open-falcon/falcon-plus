import bytes from '../bytes'

it('Convert from number', () => {
  expect(bytes(87)).toBe('87')
  expect(bytes(8787)).toBe('8.79K')
  expect(bytes(878787)).toBe('878.79K')
  expect(bytes(87878787)).toBe('87.88M')
  expect(bytes(8787878787)).toBe('8.79G')
  expect(bytes(8787878787878)).toBe('8.79T')
})

it('Convert from string', () => {
  expect(bytes('87')).toBe(87)
  expect(bytes('8.79K')).toBe(8790)
  expect(bytes('87K')).toBe(87000)
  expect(bytes('87M')).toBe(87000 * 1000)
  expect(bytes('87G')).toBe(87000 * 1000 * 1000)
  expect(bytes('87T')).toBe(87000 * 1000 * 1000 * 1000)
})

it('Test with `space`', () => {
  expect(bytes(8787, { space: true })).toBe('8.79 K')
})

it('Test with `decimals`', () => {
  expect(bytes(8787, { decimals: 1 })).toBe('8.8K')
})

it('Throw error if pass not a string or number', () => {
  expect(() => {
    bytes({})
  }).toThrow()
})

it('Throw error if pass a invalid string', () => {
  expect(() => {
    bytes('87A')
  }).toThrow()
})
