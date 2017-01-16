import { Tab, Input, Button, Grid, Icon, LightBox, Flex } from '@cepave/owl-ui'
import Link from '~coms/link'
import s from '../template.scss'
import u from '../../user/user.scss'
const { _ } = window

const TemplatePage = {
  name: 'TemplatePage',
  data() {
    const { $store, $router } = this
    return {
      gridData: {
        heads: [
          {
            col: 'Template Name',
            sort: -1,
            width: '25%',
            render(h, head) {
              return (
                <b>{head.col}</b>
              )
            }
          },
          {
            col: 'Parent',
            sort: -1,
            width: '25%',
          },
          {
            col: 'Creator',
            sort: -1,
            width: '20%',
          },
          {
            col: 'Opeartions',
            width: '30%',
          }
        ],
        rowsRender() {},
      },
      tplId: 0,
    }
  },
  created() {
    this.gridData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>{row[0].col}</Grid.Col>,
        <Grid.Col>{row[1].col}</Grid.Col>,
        <Grid.Col>{row[2].col}</Grid.Col>,
        <Grid.Col>
          <div class={[u.opeartionInline]}>
            <a class={[u.opeartions]} href={row[3].col[0]}>Edit</a>
            <span class={[u.opeartions]} tip={row[3].col[1]} onClick={(e) => this.deleteTemplateLink(e)}>Delete</span>
          </div>
        </Grid.Col>,
      ]
    }
    this.getTemplates('.+')
  },
  methods: {
    getTemplates(q) {
      const { $store, $refs } = this
      $store.dispatch('getTemplates', {
        q,
      })
    },
    openNTemplate(e) {
      this.$refs.NTemplate.open(e)
    },
    submitQuery(e) {
      if ((e.type === 'keypress' && e.charCode === 13) || e.type === 'click') {
        const queryString = this.$refs.inputRef.value || '.+'
        this.getTemplates(queryString)
      }
    },
    createTemplate(e) {
      const name = this.$refs.tplName.value
      if (name !== '') {
        const { $store, $refs } = this
        $store.dispatch('createTemplate', {
          data: {
            name,
            parent_id: 0,
          },
          q: this.$refs.inputRef.value || '.+'
        })
        this.$refs.NTemplate.close(e)
      }
    },
    deleteTemplateLink(e) {
      this.tplId = e.target.getAttribute('tip')
      this.$refs.DeleteTemplate.open(e)
    },
    deleteTemplate(id) {
      const { $store, $refs } = this
      $store.dispatch('deleteTemplate', {
        id,
        q: this.$refs.inputRef.value || '.+',
      })
    }
  },
  render(h) {
    const { $store, $refs, $slots, gridData } = this
    const props = {
      rows: $store.state.templateList.rows,
      ...gridData
    }
    //NewTemplateView
    const NewTemplateView = (
      <LightBox ref="NTemplate" closeOnClickMask closeOnESC>
        <LightBox.View>
          <h3>新增模板</h3>
          <Flex class={[s.flexWrapper]}>
            <Flex.Col size="10">
              <Input class={[s.searchTemplate]} placeholder="name" ref="tplName" />
            </Flex.Col>
            <Flex.Col size="2">
              <Button status="primary" nativeOn-click={(e) => this.createTemplate(e, this)}>
                储存
              </Button>
            </Flex.Col>
          </Flex>
        </LightBox.View>
      </LightBox>
    )

    const DeleteTemplateView = (
      <LightBox ref="DeleteTemplate" closeOnClickMask closeOnESC>
       <LightBox.View>
         <p>Delete a template</p>
         <p class={[u.deleteDes]}>You will remove this template: {this.tplId}, Are you sure？</p>
         <div class={[u.buttonGroup]}>
           <LightBox.Close class={[u.btnWrapper]}>
             <Button status="primary" class={[u.buttonBig]} nativeOn-click={(e) => this.deleteTemplate(this.tplId)}>Yes</Button>
           </LightBox.Close>
           <LightBox.Close class={[u.btnWrapper]}>
             <Button status="primaryOutline" class={[u.buttonBig]}>Cancel</Button>
           </LightBox.Close>
         </div>
       </LightBox.View>
     </LightBox>
    )

    return (
      <div class={[s.templatePage]}>
        <Tab>
          <Tab.Head slot="tabHead" name="profile" isSelected={true}>告警模板</Tab.Head>
          <Tab.Content slot="tabContent" name="template">
            <div>
              <div class={[u.contactSearchWrapper]}>
                <div style="display: flex;">
                  <Input class={s.searchTemplate} name="q" icon={['search', '#919799']} status="normal" placeholder="输入模板关键字查询" nativeOn-keypress={(e) => this.submitQuery(e)} ref="inputRef" />
                  <Button class={s.submitButton} status="primary" nativeOn-click={(e) => this.submitQuery(e)} >Submit</Button>
                  <Button status="primary" class={u.buttonIcon} nativeOn-click={(e) => this.openNTemplate(e, this)}>
                    <Icon typ="plus" size={16} />
                      新增模板
                  </Button>
                </div>
              </div>
              <div class={[u.contactWrapper]}>
                <div class={[u.gridWrapperBox]}>
                  { $slots.default }
                  <Grid {...{ props }} />
                </div>
              </div>
              {NewTemplateView}
              {DeleteTemplateView}
            </div>
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}

module.exports = TemplatePage
