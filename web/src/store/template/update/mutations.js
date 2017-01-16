module.exports = {
  'getTemplate'(state, u) {
    state.name = { id: u.tpl.id, name: u.tpl.tpl_name }
    state.parent = {
      id: u.tpl.parent_id || '0',
      name: u.parentName,
    }
    state.actionId = u.caction.id
    state.uics = u.caction.uic
    state.action = {
      id: u.caction.id,
      url: u.caction.url,
      before_callback_sms: (u.caction.before_callback_sms !== 0) ? true : false,
      before_callback_mail: (u.caction.before_callback_mail !== 0) ? true : false,
      after_callback_sms: (u.caction.after_callback_sms !== 0) ? true : false,
      after_callback_mail: (u.caction.after_callback_mail !== 0) ? true : false,
    }
    state.strategys = u.strategys
  },
  'getStrategy'(state, strategy) {
    state.ustrategy = strategy
  },
  'getMetric'(state, metrics) {
    state.metrics = metrics
  },
  'getTeamList'(state, teams) {
    state.teamList = teams
  },
  'getReponse'(state, data) {
    state.apiRepose = data
  }
}
