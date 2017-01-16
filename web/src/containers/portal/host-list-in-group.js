import { Button, Grid, Icon, LightBox, Flex } from '@cepave/owl-ui'
import g from 'sass/global.scss'
import s from './portal.scss'


const hostListInGroup = {
  name: 'HostListInGroup',

  data() {
    return {
      hostListData: {
        heads: [
          {
            col: 'Host Name',
            width: '60%',
            sort: -1,
          },
          {
            col: 'Operation',
            width: '40%',
          },
        ],
        rows: [],
      }
    }
  },

  created() {
    const state = this.$store.state.portal
    this.hostListData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>
          {row.hostname}
        </Grid.Col>,
        <Grid.Col>
          <ul>
            <li class={[s.operrationItem]}>
              <a
                class={[s.operration]}
                data-group-name={state.hostInGroupList.groupName}
                data-group-id={state.hostInGroupList.groupId}
                data-host-id={row.id}
                data-host-name={row.hostname}
                href
                onClick={(e) => this.deleteHostFromGroup(e, this)}
              >
                Delete
              </a>
            </li>
          </ul>
        </Grid.Col>,
      ]
    }
  },

  methods: {
    deleteHostFromGroup(e) {
      e.preventDefault()
      const { hostName } = e.target.dataset

      if (confirm(`Remove Host ${hostName} ?`)) {
        this.$store.dispatch('portal/deleteHostFromGroup', e.target.dataset)
      }
    }
  },

  render(h) {
    const { hostListData, $store } = this
    const props = {
      ...hostListData,
      rows: $store.state.portal.hostList.hostListItems
    }

    return (
      <div class={[s.pluginsListWrapper]}>
        <Grid {...{ props }}></Grid>

        {/* Delete LightBox */}
        <LightBox class={[g.inline]} ref="lbDeleteHostFromGroup" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>Delete this Host Group ?</p>
            <div class={[s.lbViewBox]}>
              <Flex>
                <Flex.Col size="auto">
                  <Button status="primary" class={[s.buttonBig]}>Yes</Button>
                </Flex.Col>
                <Flex.Col size="auto">
                  <Button status="primaryOutline" class={[s.buttonBig]} nativeOnClick={(e) => this.$refs.lbDeleteHostFromGroup.close(e)}>NO</Button>
                </Flex.Col>
              </Flex>
            </div>
          </LightBox.View>
        </LightBox>
      </div>
    )
  }
}

module.exports = hostListInGroup
