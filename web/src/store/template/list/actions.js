import fetch from '~utils/fetch'

module.exports = {
  'getTemplates'({ commit, state }, { q }) {
    const opts = {
      url: 'template',
      params: {
        limit: 500,
        q,
      },
      commit,
      mutation: 'getTemplates',
    }

    return fetch(opts)
      .then((res) => {
        const tpls = res.data.templates.map((tpl) => {
          return [
            {
              col: tpl.template.tpl_name,
            },
            {
              col: tpl.parent_name,
            },
            {
              col: tpl.template.create_user
            },
            {
              col: [[`/#/template/${tpl.template.id}`], [tpl.template.id]]
            }
          ]
        })
        commit('getTemplates', tpls)
      })
  },
  'getSimpleTplList'({ commit, state }, { limit }) {
    const opts = {
      url: 'template_simple',
      params: {
        limit: limit || 500,
      },
      commit,
      mutation: 'getSimpleTplList',
    }

    return fetch(opts)
      .then((res) => {
        const tpls = res.data.map((tpl) => {
          return { id: tpl.id.toString(), text: tpl.tpl_name }
        })
        commit('getSimpleTplList', tpls)
      })
  },
  'createTemplate'({ commit, state, dispatch }, d) {
    const opts = {
      url: 'template',
      data: d.data,
      method: 'POST',
      commit,
      mutation: 'createTemplate',
    }

    return fetch(opts)
      .then((res) => {
        commit('createTemplate', res.data)
        dispatch('getTemplates', d.q)
      })
  },
  'deleteTemplate'({ commit, state, dispatch }, d) {
    const opts = {
      url: `/template/${d.id}`,
      method: 'DELETE',
      commit,
    }

    return fetch(opts)
      .then((res) => {
        dispatch('getTemplates', d.q)
      })
  }
}
