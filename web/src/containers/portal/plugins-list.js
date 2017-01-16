import { Input, Button, Grid, Icon, LightBox } from '@cepave/owl-ui'
import g from 'sass/global.scss'
import s from './portal.scss'

const pluginsList = {
  name: 'PluginsList',
  data() {
    return {
      pluginData: {
        heads: [
          {
            col: 'Plugin Dir',
            width: '40%',
            sort: -1,
          },
          {
            col: 'Creator',
            width: '40%',
            sort: -1,
          },
          {
            col: 'Operation',
            width: '20%',
          },
        ],
        rows: [],
      }
    }
  },

  created() {
    this.pluginData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>
          {row.dir}
        </Grid.Col>,
        <Grid.Col>
          {row.create_user}
        </Grid.Col>,
        <Grid.Col>
          <ul>
            <li class={[s.operrationItem]}>
              <a
                class={[s.operration]}
                href
                data-group-id={row.grp_id}
                data-dir={row.dir}
                data-id={row.id}
                onClick={this.unbindPluginHandler}
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
    getBindPluginList(groupId) {
      this.$store.dispatch('portal/getBindPluginList', groupId)
    },

    bindInputHandler(e) {
      if (e.charCode === 13) {
        this.bindPluginHandler()
      }
    },

    bindPluginHandler() {
      const plugin = this.$refs.pluginBindInput.value || ''
      if (!plugin.length) {
        return
      }

      const data = {
        groupId: this.$store.state.portal.bindPluginCandidate.groupId,
        pluginDir: plugin,
      }

      this.$store.dispatch('portal/bindPluginToHostGroup', data)
    },

    unbindPluginHandler(e) {
      e.preventDefault()
      const { dir, id } = e.target.dataset

      if (confirm(`Unbind ${dir} ?`)) {
        this.$store.dispatch('portal/unbindPluginFromGroup', e.target.dataset)
      }
    },

  },

  render(h) {
    const { pluginData, $store } = this
    const props = {
      ...pluginData,
      rows: $store.state.portal.pluginsList
    }

    return (
      <div class={[s.lbViewBox]}>
        <div class={[s.searchInput]}>
          <div class={[s.inputGroups]}>
            <Input
              ref="pluginBindInput"
              loading={$store.state.portal.getBindPluginListLoading}
              nativeOnKeypress={(e) => this.bindInputHandler(e)}
            />
            <span class={[s.btnAppend]}>
              <Button status="primary" nativeOnClick={this.bindPluginHandler}>
                <Icon class="create-icon" typ="plus" fill="#fff" size={18} />
                Bind
              </Button>
            </span>
          </div>
        </div>
        <div class={[s.pluginsListWrapper]}>
          <Grid {...{ props }}></Grid>
        </div>
      </div>
    )
  }
}

module.exports = pluginsList
