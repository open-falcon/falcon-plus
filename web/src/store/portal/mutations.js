module.exports = {
  'getEndpoints.start'(state) {
    state.hasEndpointLoading = true
  },

  'getEndpoints.end'(state) {
    state.hasEndpointLoading = false
  },

  'getEndpoints.success'(state, { data }) {
    state.hosts = data
  },

  'getEndpoints.fail'(state) {

  },

  'clearHosts'(state) {
    state.hosts = []
  },

  hostGroupSearchInput(state, value) {
    state.hostGroupSearchInput = value
  },

  'updateNewHostGroup'(state, hosts) {
    state.selectedHosts = hosts
  },

  'updateEditHostGroup'(state, hosts) {
    state.editSelectedHosts = hosts
  },

  'addHostsIntoNewHostGroup.start'(state) {
    state.hasCreateLoading = true
  },

  'addHostsIntoNewHostGroup.end'(state) {
    state.hasCreateLoading = false
  },

  'addHostsIntoNewHostGroup.success'(state, { data }) {

  },

  // getHostGroupList
  'getHostGroupList.start'(state) {
    state.searchHostGrooupLoading = true
  },

  'getHostGroupList.end'(state) {
    state.searchHostGrooupLoading = false
  },

  'getHostGroupList.success'(state, { data }) {
    state.hostGroupListItems = data
  },

  'getHostGroupList.fail'(state) {

  },

  // searchHostGroup
  'searchHostGroup.start'(state) {
    state.searchHostGrooupLoading = true
  },

  'searchHostGroup.end'(state) {
    state.searchHostGrooupLoading = false
  },

  'searchHostGroup.success'(state, { data }) {
    state.hostGroupListItems = data
  },

  'updateSearchHostGroupInput'(state, q) {
    state.searchHostGroupInput = q
  },

  'searchHostGroup.fail'(state) {

  },

  'tempDeleteCandidate'(state, data) {
    state.tempDeleteCandidate = data
  },

  // getHostsList
  'getHostsList.start'(state) {
    state.getHostsListLoading = true
  },

  'getHostsList.end'(state) {
    state.getHostsListLoading = false
  },

  'getHostsList.success'(state, { data }) {
    state.hostList.hostListItems = data.hosts
    state.hostList.hostGroup = data.hostgroup.grp_name
    state.hostList.hostGroupId = data.hostgroup.id
  },

  'getHostsList.fail'(state, { data }) {

  },

  'searchHostList.success'(state, { data }) {
    state.hostList.hostListItems = data.hosts
  },

  'searchHostList.fail'(state, { err }) {

  },

  'hostInGroupList'(state, data) {
    state.hostInGroupList = data
  },

  // deleteHostGroup
  'deleteHostGroup.start'(state) {

  },

  'deleteHostGroup.end'(state) {

  },

  'deleteHostGroup.success'(state, data) {

  },

  'deleteHostGroup.fail'(state, err) {

  },

  'deleteHostFromGroup.start'(state) {

  },

  'deleteHostFromGroup.end'(state) {

  },

  'deleteHostFromGroup.success'(state, data) {

  },

  'deleteHostFromGroup.fail'(state, err) {

  },

  // Get PluginsList
  'getBindPluginList.start'(state) {
    state.getBindPluginListLoading = true
  },
  'getBindPluginList.end'(state) {
    state.getBindPluginListLoading = false
  },
  'getBindPluginList.success'(state, res) {
    state.pluginsList = res.data
  },
  'getBindPluginList.fail'(state, err) {
    console.error(err)
  },

  // Plugin binding
  'bindPluginCandidate'(state, data) {
    state.bindPluginCandidate.groupId = +data.groupId
  },

  'bindPluginToHostGroup.start'(state) {},
  'bindPluginToHostGroup.end'(state) {},
  'bindPluginToHostGroup.success'(state, data) {},
  'bindPluginToHostGroup.fail'(state, err) {
    console.error(err)
  },

  // Plugin unbind from hostgroup
  'unbindPluginFromGroup.start'(state) {},
  'unbindPluginFromGroup.end'(state) {},
  'unbindPluginFromGroup.success'(state, data) {},
  'unbindPluginFromGroup.fail'(state, err) {
    console.error(err)
  },

  'getBindTemplates.success'(state, { data }) {
    state.templateData = data
  },
  'getBindTemplates.fail'(state, err) {

  },
  'getTemplates.success'(state, { data }) {
    state.templates = data.templates.reduce((preVal, curVal, idx) => {
      preVal.push({
        id: idx,
        text: curVal.template.tpl_name,
        tId: curVal.template.id,
      })
      return preVal
    }, [])
  },
  'getTemplates.fail'(state, err) {

  },
  'bindOneTemplate.success'(state, { data }) {

  },
  'bindOneTemplate.fail'(state, err) {

  },
  'unbindOneTemplate.success'(state, { data }) {

  },
  'unbindOneTemplate.fail'(state, err) {

  },
}
