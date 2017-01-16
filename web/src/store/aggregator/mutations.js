module.exports = {
  'get/aggregators.success'(state, { data }) {
    state.rows = data.aggregators
    state.currentHostGroupName = data.hostgroup
  },
  'get/aggregators.fail'(state, { err }) {
  },
  'add/aggregators.success'(state, { data }) {},
  'add/aggregators.fail'(state, { err }) {
  },
  'edit/aggregators.success'(state, { data }) {
  },
  'edit/aggregators.fail'(state, { err }) {
  },
  'delete/aggregators.success'(state, { data }) {
  },
  'delete/aggregators.fail'(state, { err }) {
  }
}
