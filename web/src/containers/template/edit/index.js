import { Tab, Input, Button, Grid, Icon, Checkbox, Label, LightBox, Flex, Select } from '@cepave/owl-ui'
import Link from '~coms/link'
import s from '../template.scss'
import u from '../../user/user.scss'
import Select2 from '../common/select2'
import Select2Muti from '../common/select2-multi'

const gridBase = {
  heads: [
    {
      width: '30%',
      col: 'Metric/Tag/Note',
      render(h, head) {
        return (
          <b>{head.col}</b>
        )
      }
    },
    {
      width: '15%',
      col: 'Condition'
    },
    {
      width: '10%',
      col: 'Max'
    },
    {
      width: '10%',
      col: 'P'
    },
    {
      width: '10%',
      col: 'Run'
    },
    {
      width: '25%',
      col: 'Operation'
    }
  ],
  rowsRender() {},
}
const TemplatePage = {
  name: 'TemplatePage',
  data() {
    const { $store, $router } = this
    const action = $store.state.templateUpdate.action
    return {
      gridData: gridBase,
      metricMap: [],
      callback: action.callback || 1,
      callbackActions: {
        1: action.before_callback_sms,
        2: action.before_callback_mail,
        3: action.after_callback_sms,
        4: action.after_callback_mail,
      },
      updateStrategyError: [],
      newStrategyError: [],
      strategyId: 0,
      opOptions: [
        { value: '==', title: '==', selected: true },
        { value: '!=', title: '!=' },
        { value: '<', title: '<' },
        { value: '<=', title: '<=' },
        { value: '>', title: '>' },
        { value: '>=', title: '>=' },
      ],
    }
  },
  created() {
    this.getTemplate()
    this.getMetric()
    this.gridData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>{row[0].col}</Grid.Col>,
        <Grid.Col>{row[1].col}</Grid.Col>,
        <Grid.Col>{row[2].col}</Grid.Col>,
        <Grid.Col>{row[3].col}</Grid.Col>,
        <Grid.Col>{row[4].col}</Grid.Col>,
        <Grid.Col>
          <div class={[u.opeartionInline]}>
            <span class={[u.opeartions]} sid={row[5].col} saction="update" onClick={(e) => this.getStragtegy(e, this)}>Edit</span>
            <span class={[u.opeartions]} sid={row[5].col} saction="delete" onClick={(e) => this.deleteStrategyLink(e, this)}>Delete</span>
          </div>
        </Grid.Col>,
      ]
    }
  },
  mounted() {
    this.getSimpleTplList()
    this.getTeamList()
  },
  methods: {
    getSimpleTplList() {
      const { $store, $refs } = this
      $store.dispatch('getSimpleTplList', {})
    },
    getTeamList() {
      const { $store, $refs } = this
      $store.dispatch('getTeamList', {})
    },
    getTemplate() {
      const { $store, $refs } = this
      $store.dispatch('getTemplate', parseInt(this.$route.params.id))
    },
    getMetric() {
      const { $store, $refs } = this
      $store.dispatch('getMetric', {})
    },
    getCheckboxData(data) {
      this.callbackActions = data
    },
    getStragtegy(e) {
      e.preventDefault()
      const id = e.currentTarget.attributes.sid.value
      const { $store, $refs } = this
      $store.dispatch('getStrategy', id)
      this.getMetric()
      this.$refs.LStragtegy.open(e)
    },
    openNewStragtegy(e) {
      this.$refs.NStragtegy.open(e)
    },
    deleteStragtegy(e) {
      const id = e.currentTarget.attributes.sid.value
      this.$refs.DStragtegy.open(e)
    },
    checkFormating(data) {
      const err = []
      if (data.right_value === '') {
        err.push('right_value is empty')
      }
      if (data.op === '') {
        err.push('op is empty')
      }
      if (data.metric === '') {
        err.push('metric is empty')
      }
      if (data.max_step === 0 || data.max_step === '') {
        err.push('max_step is not set')
      }
      return err
    },
    updateMetric(e) {
      const { $store } = this
      const UpdateStrategy = {
        tags: this.$refs.updateTags.value,
        run_end: this.$refs.updateRunEnd.value,
        run_begin: this.$refs.updateRunBegin.value,
        right_value: this.$refs.updateRightValue.value,
        priority: parseInt(this.$refs.updatePriority.value) || 0,
        op: this.$refs.updateOp.value,
        note: this.$refs.updateNote.value,
        metric: this.$refs.updateMetric.value,
        max_step: parseInt(this.$refs.updateMaxStep.value) || 0,
        id: parseInt(this.$refs.updateId.value) || 0,
        func: this.$refs.updateFunc.value
      }
      const err = this.checkFormating(UpdateStrategy)
      this.updateStrategyError = err
      if (err.length === 0) {
        $store.dispatch('updateStrategy', {
          data: UpdateStrategy,
          id: `${$store.state.templateUpdate.name.id}`,
        })
        this.$refs.LStragtegy.close(e)
      }
    },
    newMetric(e) {
      const { $store } = this
      const NewStrategy = {
        tpl_id: $store.state.templateUpdate.name.id,
        tags: this.$refs.newTags.value,
        run_end: this.$refs.newRunEnd.value,
        run_begin: this.$refs.newRunEnd.value,
        right_value: this.$refs.newRightValue.value,
        priority: parseInt(this.$refs.newPriority.value) || 0,
        op: this.$refs.newOp.value,
        note: this.$refs.newNote.value,
        metric: this.$refs.newMetric.value,
        max_step:  parseInt(this.$refs.newMaxStep.value) || 0,
        func: this.$refs.newFunc.value
      }
      const err = this.checkFormating(NewStrategy)
      this.newStrategyError = err
      if (err.length === 0) {
        $store.dispatch('newStrategy', {
          data: NewStrategy,
          id: `${$store.state.templateUpdate.name.id}`,
        })
        this.getTemplate()
        this.$refs.NStragtegy.close(e)
      }
    },
    SaveTemplate(m) {
      const { $store, $refs } = this
      const parentId = $refs.updateParent.value
      const postbody = {
        tpl_id: $store.state.templateUpdate.name.id,
        parent_id: parseInt(parentId || this.parentId),
        name: $store.state.templateUpdate.name.name,
      }
      $store.dispatch('updateTemplate', {
        data: postbody,
        id: `${$store.state.templateUpdate.name.id}`,
      })
    },
    SaveAction(m) {
      const { $store, $refs } = this
      const uic = $refs.updateTeam.value
      const action = {
        url: $refs['action.url'].value,
        uic,
        id: $store.state.templateUpdate.action.id || 0,
        tpl_id: $store.state.templateUpdate.name.id,
        callback: this.callback,
        before_callback_sms: (this.callbackActions['1']) ? 1 : 0,
        before_callback_mail:  (this.callbackActions['2']) ? 1 : 0,
        after_callback_sms:  (this.callbackActions['3']) ? 1 : 0,
        after_callback_mail:  (this.callbackActions['4']) ? 1 : 0,
      }
      if (action.id === 0) {
        //create a new action
        $store.dispatch('createAction', {
          data: action,
          id: `${$store.state.templateUpdate.name.id}`,
        })
      } else {
        //update a action
        $store.dispatch('updateAction', {
          id: `${$store.state.templateUpdate.name.id}`,
          data: action,
        })
      }
    },
    deleteStrategyLink(e) {
      this.$refs.DeleteStrategy.open(e)
      this.strategyId = e.target.getAttribute('sid')
    },
    deleteStrategy(id) {
      const { $store, $refs } = this
      $store.dispatch('deleteStrategy', {
        id: this.strategyId,
        tid: $store.state.templateUpdate.name.id,
      })
    }
  },
  render(h) {
    const { $store, $refs, $slots, gridData } = this
    const props = {
      rows: $store.state.templateUpdate.strategys,
      tname: $store.state.templateUpdate.name.name,
      pname: $store.state.templateUpdate.parent.name,
      action: $store.state.templateUpdate.action,
      uics: $store.state.templateUpdate.uics,
      metric: $store.state.templateUpdate.ustrategy.metric,
      ...gridData
    }
    const merticProps = {
      options: $store.state.templateUpdate.metrics,
      value: $store.state.templateUpdate.ustrategy.metric,
      class: s.inputMetric,
      placeholder: 'Select a metric',
      accpetTag: true,
    }
    const { tags, priority, func, op, note, id, metric  } = $store.state.templateUpdate.ustrategy
    const maxStep = $store.state.templateUpdate.ustrategy.max_step
    const rightValue = $store.state.templateUpdate.ustrategy.right_value
    const runBegin = $store.state.templateUpdate.ustrategy.run_begin
    const runEnd = $store.state.templateUpdate.ustrategy.run_end
    //UpdateStrategyView
    const UpdateStrategyView = (
      <LightBox ref="LStragtegy" closeOnClickMask closeOnESC>
        <LightBox.View>
          <div class={[s.formWrapper]}>
            <Flex>
              <Flex.Col size="6">
                <p>Metric</p>
                <div>
                  <input class='newInputMetric' type="hidden" ref="newMetric"></input>
                  <Select2 {...{ props: { ...merticProps } } }  />
                </div>
              </Flex.Col>
              <Flex.Col size="6">
                <p>Tags</p>
                <Input class={[s.inputGrid6]} placeholder="tags" val={tags} ref="updateTags" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="6">
                <p>run begin</p>
                <Input class={[s.inputGrid6]} placeholder="run begin: 00:00" val={runBegin} ref="updateRunBegin" />
              </Flex.Col>
              <Flex.Col size="6">
                <p>run end</p>
                <Input class={[s.inputGrid6]} placeholder="run end: 00:00" val={runEnd} ref="updateRunEnd" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="6">
                <p>Max</p>
                <Input class={[s.inputGrid6]} placeholder="max" val={maxStep} ref="updateMaxStep" />
              </Flex.Col>
              <Flex.Col size="6">
                <p>P</p>
                <Input class={[s.inputGrid6]} placeholder="p" val={priority} ref="updatePriority" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="4">
                <p>Function</p>
                <Input placeholder="func" val={func} ref="updateFunc" />
              </Flex.Col>
              <Flex.Col size="3">
                <p>Op</p>
                {/* <Select options={this.opOptions} ref="newOp" class={[s.inputGrid4]} /> */}
                <Input class={[s.inputGrid4]} placeholder="op" val={op} ref="updateOp" />
              </Flex.Col>
              <Flex.Col size="5">
                <p>Right Value</p>
                <Input class={[s.inputGrid5]} placeholder="reght value" val={rightValue} ref="updateRightValue" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="12">
                <p>Notes</p>
                <Input class={[s.inputGrid12]} placeholder="note" val={note} ref="updateNote" />
              </Flex.Col>
            </Flex>
          </div>
          <div>
            <Flex>
              <Flex.Col size="2" offset="10">
                <Button status="primary" nativeOn-click={(e) => this.updateMetric(e, this)}>
                  Save
                </Button>
              </Flex.Col>
            </Flex>
          </div>
        </LightBox.View>
      </LightBox>
    )
    const newMerticProps = {
      options: $store.state.templateUpdate.metrics,
      value: '',
      class: s.inputMetric,
      setclass: 'newInputMetric',
      placeholder: 'Select a metric',
      accpetTag: true,
    }
    //NewStrategyView
    const NewStrategyView = (
      <LightBox ref="NStragtegy" closeOnClickMask closeOnESC>
        <LightBox.View>
          <div class={[s.formWrapper]}>
            <Flex>
              <Flex.Col size="6">
                <p>Metric</p>
                <div>
                  <input class='newInputMetric' type="hidden" ref="newMetric"></input>
                  <Select2 {...{ props: { ...newMerticProps } } }  />
                </div>
              </Flex.Col>
              <Flex.Col size="6">
                <p>Tags</p>
                <Input class={[s.inputGrid6]} placeholder="tags" ref="newTags" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="6">
                <p>run begin</p>
                <Input class={[s.inputGrid6]} placeholder="run begin: 00:00" ref="newRunBegin" />
              </Flex.Col>
              <Flex.Col size="6">
                <p>run end</p>
                <Input class={[s.inputGrid6]} placeholder="run end: 00:00" ref="newRunEnd" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="6">
                <p>Max</p>
                <Input class={[s.inputGrid6]} placeholder="max" ref="newMaxStep" />
              </Flex.Col>
              <Flex.Col size="6">
                <p>P</p>
                <Input class={[s.inputGrid6]} placeholder="p" ref="newPriority" />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="4">
                <p>Function</p>
                <Input placeholder="func" ref="newFunc" />
              </Flex.Col>
              <Flex.Col size="3">
                <p>Op</p>
                <Select options={this.opOptions} ref="newOp" class={[s.inputGrid4]} />
              </Flex.Col>
              <Flex.Col size="5">
                <p>Right Value</p>
                <Input placeholder="reght value" ref="newRightValue" class={[s.inputGrid5]} />
              </Flex.Col>
            </Flex>
            <Flex>
              <Flex.Col size="12">
                <p>Notes</p>
                <Input class={[s.inputGrid12]} placeholder="note" ref="newNote" />
              </Flex.Col>
            </Flex>
          </div>
          <div>
            <Flex>
              <Flex.Col size="2" offset="10">
                <Button status="primary" nativeOn-click={(e) => this.newMetric(e, this)}>
                  Save
                </Button>
              </Flex.Col>
            </Flex>
          </div>
        </LightBox.View>
      </LightBox>
    )
    //DeleteStrategyView
    const DeleteStrategyView = (
      <LightBox ref="DeleteStrategy" closeOnClickMask closeOnESC>
       <LightBox.View>
         <p>Delete a strategy</p>
         <p class={[u.deleteDes]}>You will remove this template: {this.strategyId}, Are you sure？</p>
         <div class={[u.buttonGroup]}>
           <LightBox.Close class={[u.btnWrapper]}>
             <Button status="primary" class={[u.buttonBig]} nativeOn-click={(e) => this.deleteStrategy(this.strategyId)}>Yes</Button>
           </LightBox.Close>
           <LightBox.Close class={[u.btnWrapper]}>
             <Button status="primaryOutline" class={[u.buttonBig]}>Cancel</Button>
           </LightBox.Close>
         </div>
       </LightBox.View>
     </LightBox>
    )

    const tplProps = {
      options: $store.state.templateList.simpleTList,
      value: $store.state.templateUpdate.parent.name,
      pid: $store.state.templateUpdate.parent.id,
      class: 'newParentSelect',
      setclass: 'newParent',
      placeholder: 'Select a template',
      accpetTag: false,
    }
    const teamProps = {
      options: $store.state.templateUpdate.teamList,
      value: $store.state.templateUpdate.uics,
      // pid: $store.state.templateUpdate.parent.id,
      class: 'newTeamSelect',
      setclass: 'newTeam',
    }

    return (
      <div class={[s.templatePage]}>
        <Tab>
          <Tab.Head slot="tabHead" name="profile" isSelected={true}>Edit Template</Tab.Head>
          <Tab.Content slot="tabContent" name="template" class={[s.templatePage]}>
            <div class={[u.contactWrapper, s.lightDivButtom]} style="top: 0px">
              <div class={[s.nameGroup]}>
                <Flex class={[s.flexWrapper]}>
                  <Flex.Col size="4">
                    <span class={[s.title]}>name:</span>
                    <Label typ="tag">{props.tname}</Label>
                  </Flex.Col>
                  <Flex.Col size="7">
                    <span class={[s.title]}>parent:</span>
                    <input type="hidden" class='newParent' placeholder="请输入模板名称" ref="updateParent" value={tplProps.pid}></input>
                    <Select2 { ...{ props: tplProps } } />
                  </Flex.Col>
                  <Flex.Col size="1">
                    <Button status="primary" nativeOn-click={(m) => this.SaveTemplate(m, this)}>
                      Save
                    </Button>
                  </Flex.Col>
                </Flex>
              </div>
              <div class={[s.templateGroup]}>
                <Flex class={[s.flexWrapper]}>
                  <Flex.Col size="11">
                    <p class={[s.templateTitle]}>报警接收组(在UIC中管理报警组，快捷入口):</p>
                    <div class={[s.questionBlock]}>
                      <input class='newTeam' type="hidden" placeholder="告警组" ref="updateTeam" value={props.uics}></input>
                      <Select2Muti { ...{ props: teamProps } } />
                    </div>
                  </Flex.Col>
                  <Flex.Col size="1">
                    <Button status="primary" nativeOn-click={(m) => this.SaveAction(m, this)}>
                      Save
                    </Button>
                  </Flex.Col>
                </Flex>
                <div>
                  <p class={[s.templateTitle]}>callback地址(只支持http get方式回调):</p>
                  <div class={[s.questionBlock]}>
                    <Input class={s.searchInput} name="q" placeholder="callback url" val={props.action.url} ref="action.url" />
                  </div>
                  <div class={[s.questionBlock]}>
                    <Checkbox.Group onChange={this.getCheckboxData}>
                      <Checkbox name="1" checked={props.action.before_callback_sms} >回调之前发提醒短信</Checkbox>
                      <Checkbox name="2" checked={props.action.before_callback_mail} >回调之前发提醒邮件</Checkbox>
                      <Checkbox name="3" checked={props.action.after_callback_sms} >回调之后发结果短信</Checkbox>
                      <Checkbox name="4" checked={props.action.after_callback_mail} >回调之后发结果邮件</Checkbox>
                    </Checkbox.Group>
                  </div>
                </div>
              </div>
            </div>
          </Tab.Content>
          <Tab.Head slot="tabHead" name="profile" >Edit Strategy</Tab.Head>
          <Tab.Content slot="tabContent" name="template">
            <div>
              <div class={[u.contactSearchWrapper]}>
                <Button class={[s.submitButton]} status="primary">
                  <Icon fill="#fff" typ="fold" size={16} />
                  <Link class={[s.link]} to="/template">Back to Template List</Link>
                </Button>
                <Button status="primary" class={u.buttonIcon} nativeOn-click={(e) => this.openNewStragtegy(e, this) }>
                  <Icon fill="#fff" typ="plus" size={16} />
                    Add new strategy
                </Button>
                {UpdateStrategyView}
                {NewStrategyView}
                {DeleteStrategyView}
              </div>
              <div class={[u.contactWrapper, s.lightDivButtom]} style="top: 0px">
                <div>
                  <Label  typ="tag">{props.tname}</Label>
                </div>
                { $slots.default }
                <div class={[u.gridWrapperBox]}>
                  <Grid {...{ props }} />
                </div>
              </div>
            </div>
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}

module.exports = TemplatePage
