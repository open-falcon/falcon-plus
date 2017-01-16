import { Button, Grid, Icon, LightBox, Flex } from '@cepave/owl-ui'
import Link from '~coms/link'
import Select2 from '../template/common/select2'
import g from 'sass/global.scss'
import s from './portal.scss'

const templateList = {
  name: 'TemplateList',
  props: {
    closeTemplatelb: {
      type: Function
    }
  },
  data() {
    return {
      templateData: {
        heads: [
          {
            col: 'Template Name',
            width: '40%',
            sort: -1,
          },
          {
            col: 'Creator',
            width: '30%',
            sort: -1,
          },
          {
            col: 'Operation',
            width: '30%',
          },
        ],
        rows: [],
      },
    }
  },

  created() {
    this.templateData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>
          {row.tpl_name}
        </Grid.Col>,
        <Grid.Col>
          {row.create_user}
        </Grid.Col>,
        <Grid.Col>
          <ul>
            <li class={[s.operrationItem]}>
                <Link class={[s.operration]} to={`template/${row.id}`} nativeOn-click={this.closelb}>Edit</Link>
            </li>
            <li class={[s.operrationItem]}>
              <a
                class={[s.operration]}
                href
                data-group-id={row.grp_id}
                data-dir={row.dir}
                data-id={row.id}
                onClick={this.handleUnbindTemplate}
              >
                Unbind
              </a>
            </li>
          </ul>
        </Grid.Col>,
      ]
    }
  },

  methods: {
    handleBindTemplate(e) {
      const { templates, templateData } = this.$store.state.portal
      const idx = this.$refs.updateSelectedTemplate.value
      this.$store.dispatch('portal/bindOneTemplate', {
        tpl_id: templates[idx].tId,
        grp_id: templateData.hostgroup.id
      })
    },
    closelb(e) {
      this.closeTemplatelb(e)
    },
    handleUnbindTemplate(e) {
      e.preventDefault()
      const { templates, templateData } = this.$store.state.portal
      this.$store.dispatch('portal/unbindOneTemplate', {
        tpl_id: +e.target.dataset.id,
        grp_id: templateData.hostgroup.id
      })
    }
  },

  render(h) {
    const { templateData, $store } = this
    const props = {
      ...templateData,
      rows: $store.state.portal.templateData.templates
    }
    const { templates } = this.$store.state.portal

    return (
      <div class={[s.lbViewBox]}>
        <div class={[s.binding]}>
          <Flex>
            <Flex.Col size="5">
              <input type="hidden" class="templateName" ref="updateSelectedTemplate"></input>
              <Select2 class={[s.select2]} options={templates} setclass="templateName" placeholder="Please select a template" />
            </Flex.Col>
            <Flex.Col size="7">
              <Button status="primary" class={[s.buttonIcon]} nativeOn-click={this.handleBindTemplate}>
                <Icon class={[s.plus]} typ="plus" fill="#fff" size={18} />
                Bind
              </Button>
            </Flex.Col>
          </Flex>
        </div>
        <div class={[s.pluginsListWrapper]}>
          <Grid {...{ props }}></Grid>
        </div>
      </div>
    )
  }
}

module.exports = templateList
