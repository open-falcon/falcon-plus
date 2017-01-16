import { Input, Button, Grid, Icon, LightBox, DualList, Flex } from '@cepave/owl-ui'
import Link from '~coms/link'
import g from 'sass/global.scss'
import s from './portal.scss'

import HostListInGroup from './host-list-in-group'
import HostGroupEdit from './host-group-edit'
import PluginsList from './plugins-list'
import TemplateList from './template-list'

const hostGroups = {
  name: 'HostGroups',

  data() {
    return {
      hostGroupsData: {
        heads: [
          {
            col: 'Name',
            width: '30%',
            sort: -1
          },
          {
            col: 'Creator',
            width: '22%',
          },
          {
            col: 'Host',
            width: '22%',
          },
          {
            col: 'Operation',
            width: '48%',
          },
        ],
        rows: [],
      },
      currentHostList: 0
    }
  },

  created() {
    this.hostGroupsData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>
          {row.grp_name}
        </Grid.Col>,
        <Grid.Col>
          {row.create_user}
        </Grid.Col>,
        <Grid.Col>
          <a
            class={[s.operration]}
            href
            data-group-id={row.id}
            data-group-name={row.grp_name}
            onClick={(e) => this.openHostsLightBox(e, this)}
          >
            view hosts
          </a>
        </Grid.Col>,
        <Grid.Col>
          <ul>
            <li class={[s.operrationItem]}>
              <a
                class={[s.operration]}
                href
                data-index={index}
                data-group-id={row.id}
                data-group-name={row.grp_name}
                onClick={(e) => this.openDeleteLightBox(e, this)}
              >
                Delete
              </a>
            </li>
            <li class={[s.operrationItem]}>
              <a
                class={[s.operration]}
                href
                data-group-name={row.grp_name}
                data-group-id={row.id}
                onClick={(e) => this.openHostGroupEdit(e)}
              >
                Edit
              </a>
            </li>
            <li class={[s.operrationItem]}>
              <a class={[s.operration]} href data-group-id={row.id} onClick={(e) => this.openTemplateListLightBox(e)}>Templates</a>
            </li>
            <li class={[s.operrationItem]}>
              <a class={[s.operration]} href data-group-id={row.id} onClick={(e) => this.openPluginsListLightBox(e, this)}>Plugins</a>
            </li>
            <li class={[s.operrationItem]}>
              <Link class={[s.operration]} to={`aggregator/${row.id}`}>Aggregator</Link>
            </li>
          </ul>
        </Grid.Col>,
      ]
    }
  },

  methods: {
    getEndpoints(q) {
      const { $store } = this
      $store.dispatch('portal/getEndpoints', {
        q,
      })
    },

    searchInputHandler(e) {
      if (e.charCode === 13) {
        this.searchHostGroupNameHandler()
      }
    },

    searchHostGroupNameHandler() {
      let q = this.$refs.hostGroupSearchInput.value
      if (!q.length) {
        q = '.+'
      }
      this.$store.commit('portal/updateSearchHostGroupInput', q)

      this.$store.dispatch('portal/searchHostGroup', q)
    },

    createHostGroup() {
      const { $store, $refs } = this
      const name = $refs.newHostGroupName.value
      const hosts = Object.keys($store.state.portal.selectedHosts).map((key) => {
        return $store.state.portal.selectedHosts[key].endpoint
      })

      const data = {
        name,
        hosts,
      }

      $store.dispatch('portal/createHostGroupName', data)
    },

    newHostsListHandle(q) {
      this.$store.commit('portal/clearHosts')
      this.getEndpoints(q)
    },

    newHostsSelectHandle(hosts) {
      this.$store.commit('portal/updateNewHostGroup', hosts)
    },

    getHostGroupList() {
      this.$store.dispatch('portal/getHostGroupList')
    },

    submitDeleteHostGroupHandle(e) {
      const data = this.$store.state.portal.tempDeleteCandidate
      this.$store.dispatch('portal/deleteHostGroup', data)
      this.$refs.lbDeleteHostGroupList.close(e)
    },

    // LightBox methods
    openDeleteLightBox(e) {
      const { index, groupId, groupName } = e.target.dataset
      const data = {
        index,
        name: groupName,
        id: groupId,
      }

      this.$store.commit('portal/tempDeleteCandidate', data)
      this.$refs.lbDeleteHostGroupList.open(e)
    },

    openPluginsListLightBox(e) {
      const data = e.target.dataset

      this.$store.commit('portal/bindPluginCandidate', data)
      this.$store.dispatch('portal/getBindPluginList', data)
      this.$refs.lbPluginsList.open(e)
    },

    openHostsLightBox(e) {
      const data = e.target.dataset
      this.currentHostList = data.groupId

      this.$store.dispatch('portal/getHostsList', data)
      this.$refs.lbHostsList.open(e)
    },

    openHostGroupEdit(e) {
      const data = e.target.dataset
      this.getEndpoints('.+')
      this.$store.dispatch('portal/getHostsList', data)
      this.$refs.lbHostGroupEdit.open(e)
    },

    openTemplateListLightBox(e) {
      const data = e.target.dataset
      this.$store.dispatch('portal/getBindTemplates', data)
      this.$store.dispatch('portal/getTemplates', data)
      this.$refs.lbTemplateList.open(e)
    },

    searchHostList(e) {
      if ((e.type === 'keypress' && e.charCode === 13) || e.type === 'click') {
        this.$store.dispatch('portal/searchHostList', {
          id: this.currentHostList,
          q: this.$refs.searchHostList.value
        })
      }
    },

    closeTemplatelb(e) {
      this.$refs.lbTemplateList.close(e)
    }
  },

  mounted() {
    this.getHostGroupList()
  },

  render(h) {
    const { hostGroupsData, $store } = this
    const props = {
      ...hostGroupsData,
      rows: $store.state.portal.hostGroupListItems
    }
    const state = $store.state
    const hosts = state.portal.hosts

    return (
      <div class={[s.hostGroupsContent]}>
        <div class={[s.searchInputWrapper]}>
          <Flex split>
            <Flex.Col>
              <div class={[s.searchInput]}>
                <div class={[s.inputGroups]}>
                  <Input
                    ref="hostGroupSearchInput"
                    icon={['search', '#919799']}
                    loading={state.portal.searchHostGrooupLoading}
                    nativeOnKeypress={this.searchInputHandler}
                    val={state.portal.searchHostGroupInput}
                    placeholder="Enter Host Group Name..."
                  />
                  <span class={[s.btnAppend]}>
                    <Button status="primary" nativeOnClick={this.searchHostGroupNameHandler}>Search</Button>
                  </span>
                </div>
              </div>
            </Flex.Col>
            <Flex.Col>
              <LightBox closeOnClickMask closeOnESC width="788">
                <LightBox.Open>
                  <Button status="primary" nativeOnClick={() => this.getEndpoints('.+')}>
                    <Icon class="create-icon" typ="plus" fill="#fff" size={18} />
                    Create HostGroups
                  </Button>
                </LightBox.Open>
                <LightBox.View>
                  <div class={[s.lb]}>
                    <h2>Create HostGroup</h2>
                    <div>
                      <Flex>
                        <Flex.Col size="2">
                          <div class={[s.lbText]}>Group Name</div>
                        </Flex.Col>
                        <Flex.Col size="10">
                          <Input ref="newHostGroupName"></Input>
                        </Flex.Col>
                      </Flex>
                      <Flex class={[s.dualBox]}>
                        <Flex.Col size="2">
                          <div class={[s.lbText]}>Hosts</div>
                        </Flex.Col>
                        <Flex.Col size="10">
                          <DualList
                            apiMode
                            onInputchange={this.newHostsListHandle}
                            onChange={this.newHostsSelectHandle}
                            items={hosts}
                            displayKey="endpoint"
                            leftLoading={state.portal.hasEndpointLoading}
                          />
                        </Flex.Col>
                      </Flex>
                    </div>
                    <div>
                      <LightBox.Close class={[s.cancelBtn]}>
                        <Flex class={[s.lbViewBox]}>
                          <Flex.Col offset="8" size="2">
                            <Button status="primaryOutline">Cancel</Button>
                          </Flex.Col>
                          <Flex.Col size="auto">
                            <Button status="primary" class={[s.buttonAlignRight]} nativeOnClick={this.createHostGroup}>Save</Button>
                          </Flex.Col>
                        </Flex>
                      </LightBox.Close>
                    </div>
                  </div>
                </LightBox.View>
              </LightBox>
            </Flex.Col>
          </Flex>
        </div>

        <div class={[s.gridWrapper]}>
          <div class={[s.gridWrapperBox]}>
            <Grid {...{ props }} />
          </div>
        </div>

        {/* LightBox Delete  */}
        <LightBox class={[g.inline]} ref="lbDeleteHostGroupList" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>You will remove this host group: <b>{state.portal.tempDeleteCandidate.name}</b>.</p>
            <p> Are you sure ?</p>
            <div class={[s.lbViewBox]}>
              <Flex>
                <Flex.Col size="auto">
                  <Button status="primary" class={[s.buttonBig]} nativeOnClick={(e) => this.submitDeleteHostGroupHandle(e, this)}>Yes</Button>
                </Flex.Col>
                <Flex.Col size="auto">
                  <Button status="primaryOutline" class={[s.buttonBig]} nativeOnClick={(e) => this.$refs.lbDeleteHostGroupList.close(e)}>NO</Button>
                </Flex.Col>
              </Flex>
            </div>
          </LightBox.View>
        </LightBox>

        {/* LightBox Hosts List */}
        <LightBox class={[g.inline]} ref="lbHostsList" closeOnClickMask closeOnESC>
          <LightBox.View>
            <h2>Group name: {state.portal.hostInGroupList.groupName}</h2>
            <div class={[s.lbViewBox]}>
              <div class={[s.searchInput]}>
                <div class={[s.inputGroups]}>
                  <Input
                    icon={['search', '#919799']}
                    ref="searchHostList"
                    loading={state.portal.getHostsListLoading}
                    placeholder="search hosts"
                    nativeOnKeypress={this.searchHostList}
                  />
                  <span class={[s.btnAppend]}>
                    <Button status="primary" nativeOnClick={this.searchHostList}>Search</Button>
                  </span>
                </div>
              </div>
              <HostListInGroup />
            </div>
          </LightBox.View>
        </LightBox>

        {/* LightBox Host Group Edit */}
        <LightBox class={[g.inline]} ref="lbHostGroupEdit" width="788px" closeOnClickMask closeOnESC>
          <LightBox.View>
            <HostGroupEdit lbRef={this.$refs.lbHostGroupEdit} />
          </LightBox.View>
        </LightBox>

        {/* LightBox Template List */}
        <LightBox class={[g.inline]} ref="lbTemplateList" closeOnClickMask closeOnESC>
          <LightBox.View>
            <TemplateList closeTemplatelb={this.closeTemplatelb} />
          </LightBox.View>
        </LightBox>

        {/* LightBox Plugins List */}
        <LightBox class={[g.inline]} ref="lbPluginsList" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>Plugins List</p>
            <PluginsList />
          </LightBox.View>
        </LightBox>
      </div>
    )
  }
}

module.exports = hostGroups
