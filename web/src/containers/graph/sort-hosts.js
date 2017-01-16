module.exports = ({ endpoints, counters }) => {
  const res = {
    endpoint: {},
    counter: {},
  }

  endpoints.forEach((endpoint) => {
    res.endpoint[endpoint] = { counters, endpoints: [endpoint] }
  })

  counters.forEach((counter) => {
    res.counter[counter] = { endpoints, counters: [counter] }
  })

  return res
}
