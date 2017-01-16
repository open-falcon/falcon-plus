import { Button, LightBox, DualList, Flex } from '@cepave/owl-ui'
import g from 'sass/global.scss'
import s from './portal.scss'

const { _ } = window

const hostGroupEdit = {
  name: 'HostGroupEdit',

  props: {
    lbRef: {
      type: Object,
    }
  },

  methods: {
    closeLb(e) {
      this.lbRef.close(e)
    },

    getEndpoints(q) {
      const { $store } = this
      $store.dispatch('portal/getEndpoints', {
        q,
      })
    },

    newHostsListHandle(q) {
      this.$store.commit('portal/clearHosts')
      this.getEndpoints(q)
    },

    newHostsSelectHandle(hosts) {
      this.$store.commit('portal/updateEditHostGroup', hosts)
    },

    save(e) {
      const { state, dispatch } = this.$store
      const hosts = Object.keys(state.portal.editSelectedHosts).map((key) => {
        return state.portal.editSelectedHosts[key].endpoint
      })
      const data = {
        id: state.portal.hostList.hostGroupId,
        hosts,
      }

      dispatch('portal/addHostsIntoNewHostGroup', data)
      this.closeLb(e)
    }
  },

  render(h) {
    const { $store, newHostsListHandle, newHostsSelectHandle } = this
    const state = $store.state

    // Object key `hostname` rename to `endpoint`
    let newHostListItems = {}
    const hasBindHostListItems = state.portal.hostList.hostListItems.map((o, i) => {
      newHostListItems = {
        ...newHostListItems,
        endpoint: o.hostname
      }

      return newHostListItems
    })

    const hosts = _.differenceBy(state.portal.hosts, hasBindHostListItems, 'endpoint')

    return (
      <div class={[s.hostGroupEditWrapper, s.lb]}>
        <h2>Edit HostGroup</h2>
        <Flex>
          <Flex.Col size="2">
            <div class={[s.lbText]}>Group Name</div>
          </Flex.Col>
          <Flex.Col size="10">
            <div class={[s.lbText]}>
              <b>{state.portal.hostList.hostGroup}</b>
            </div>
          </Flex.Col>
        </Flex>
        <div class={[s.dualBox]}>
          <Flex>
            <Flex.Col size="2">
              <div class={[s.lbText]}>Hosts</div>
            </Flex.Col>
            <Flex.Col size="10">
              <DualList
                apiMode
                onInputchange={newHostsListHandle}
                onChange={newHostsSelectHandle}
                items={hosts}
                selectedItems={hasBindHostListItems}
                displayKey="endpoint"
                leftLoading={state.portal.hasEndpointLoading}
              />
            </Flex.Col>
          </Flex>
        </div>
        <div class={[s.lbViewBox]}>
          <Flex class={[s.lbViewBox]}>
            <Flex.Col offset="8" size="2">
              <Button class={[s.cancelBtn]} status="primaryOutline" nativeOnClick={this.closeLb}>Cancel</Button>
            </Flex.Col>
            <Flex.Col size="auto">
              <Button class={[s.buttonAlignRight]} status="primary" nativeOnClick={this.save}>Save</Button>
            </Flex.Col>
          </Flex>
        </div>
      </div>
    )
  }
}

module.exports = hostGroupEdit
