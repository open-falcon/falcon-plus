import s from './aggregator.scss'
import { Tab, Grid, Input, Button, Icon, LightBox, Flex } from '@cepave/owl-ui'
import Link from '~coms/link'

const AggregatorPage = {
  name: 'AggregatorPage',
  data() {
    return {
      heads: [
        {
          col: 'HostGroup name',
          width: '14%',
          sort: -1
        },
        {
          col: 'Aggregator id',
          width: '14%',
          sort: -1
        },
        {
          col: 'endpoint',
          width: '14%',
          sort: -1
        },
        {
          col: 'metric',
          width: '14%',
          sort: -1
        },
        {
          col: 'tags',
          width: '14%',
          sort: -1
        },
        {
          col: 'creator',
          width: '14%',
          sort: -1
        },
        {
          col: 'operation',
          width: '16%'
        },
      ],
      rowsRender() {},
      grpToEdit: {},
      grpToRemove: 0,
      currentHostGroupName: ''
    }
  },
  mounted() {
    this.$store.dispatch('getAggregators', {
      hostgroup: this.$route.params.id
    })
  },
  created() {
    this.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>{this.currentHostGroupName}</Grid.Col>,
        <Grid.Col>{row.id}</Grid.Col>,
        <Grid.Col>{row.endpoint}</Grid.Col>,
        <Grid.Col>{row.metric}</Grid.Col>,
        <Grid.Col>{row.tags}</Grid.Col>,
        <Grid.Col>{row.creator}</Grid.Col>,
        <Grid.Col>
          <div class={[s.opeartionInline]}>
            <span class={[s.opeartions]} on-click={(e) => this.edit(e, row)}>edit</span>
            <span class={[s.opeartions]} on-click={(e) => this.duplicate(e, row)}>Duplicate</span>
            <span class={[s.opeartions]} on-click={(e) => this.delete(e, row.id)}>Delete</span>
          </div>
        </Grid.Col>
      ]
    }
  },
  methods: {
    addAggregator() {
      this.$store.dispatch('newAggregator', {
        numerator: this.$refs.numerator.value,
        denominator: this.$refs.denominator.value,
        endpoint: this.$refs.endpoint.value,
        metric: this.$refs.metric.value,
        tags: this.$refs.tags.value,
        step: +this.$refs.step.value,
        hostgroup_id: +this.$route.params.id
      })
    },
    duplicate(e, row) {
      this.grpToEdit = row
      this.$refs.duplicateAggregator.open(e)
    },
    duplicateAggregator() {
      this.$store.dispatch('newAggregator', {
        numerator: this.$refs.dupliNumerator.value,
        denominator: this.$refs.dupliDenominator.value,
        endpoint: this.$refs.dupliEndpoint.value,
        metric: this.$refs.dupliMetric.value,
        tags: this.$refs.dupliTags.value,
        step: +this.$refs.dupliStep.value,
        hostgroup_id: +this.$route.params.id
      })
    },
    edit(e, row) {
      this.grpToEdit = row
      this.$refs.editAggregator.open(e)
    },
    editAggregator() {
      const { grpToEdit } = this
      this.$store.dispatch('editAggregator', {
        id: grpToEdit.id,
        numerator: this.$refs.editNumerator.value || grpToEdit.numerator,
        denominator: this.$refs.editDenominator.value || grpToEdit.denominator,
        endpoint: this.$refs.editEndpoint.value || grpToEdit.endpoint,
        metric: this.$refs.editMetric.value || grpToEdit.metric,
        tags: this.$refs.editTags.value || grpToEdit.tags,
        step: +this.$refs.editStep.value || grpToEdit.step,
        hostgroup_id: +this.$route.params.id
      })
    },
    delete(e, grpId) {
      this.grpToRemove = grpId
      this.$refs.deleteAggregator.open(e)
    },
    deleteAggregator() {
      this.$store.dispatch('deleteAggregator', {
        id: this.grpToRemove,
        hostgroup_id: +this.$route.params.id
      })
    },
    newAggregator(e) {
      this.$refs.newAggregator.open(e)
    }
  },

  render(h) {
    const { heads, rowsRender, addAggregator, grpToEdit, editAggregator, deleteAggregator, duplicateAggregator, newAggregator } = this
    const { rows } = this.$store.state.aggregator
    this.currentHostGroupName = this.$store.state.aggregator.currentHostGroupName
    const props = { heads, rowsRender, rows }
    return (
      <div class={[s.aggregatorPage]}>
        <Tab>
          <Tab.Head slot="tabHead" name="current hostgroup aggregators">
            Current HostGroup Aggregators
          </Tab.Head>
          <Tab.Content slot="tabContent" name="current hostgroup aggregators">
            <div class={[s.head]}>
              <Flex grids={24}>
                <Flex.Col size="6">
                  <Link to="/portal" class={[s.linkToPortal]}>
                    <Button status="primary" class={[s.buttonIcon]}>
                      <Icon typ="fold" size={16} class={[s.plus]} />
                      Back to HostGroup List
                    </Button>
                </Link>
                </Flex.Col>
                <Flex.Col offset="14" size="5">
                  <Button status="primary" class={[s.buttonIcon]} nativeOnClick={(e) => newAggregator(e)}>
                    <Icon typ="plus" size={16} class={[s.plus]} />
                    Add new aggregator
                  </Button>
                </Flex.Col>
              </Flex>

              {/* create a new aggregator */}
              <LightBox ref="newAggregator" class={[s.headlb]} closeOnClickMask closeOnESC>
                <LightBox.View>
                  <p>Add a new aggregator</p>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Numerator</p>
                    <Input name="numerator" class={[s.input]} ref="numerator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Denominator</p>
                    <Input name="denominator" class={[s.input]} ref="denominator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Endpoint</p>
                    <Input name="endpoint" class={[s.input]} ref="endpoint" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Metric</p>
                    <Input name="metric" class={[s.input]} ref="metric" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Tags</p>
                    <Input name="tags" class={[s.input]} ref="tags" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Step</p>
                    <Input name="step" class={[s.input]} ref="step" />
                  </div>
                  <div class={[s.lightboxFooter]}>
                    <LightBox.Close>
                      <Button status="primaryOutline">Cancel</Button>
                    </LightBox.Close>
                    <LightBox.Close>
                      <Button status="primary" class={[s.submitBtn]} nativeOn-click={addAggregator}>Submit</Button>
                    </LightBox.Close>
                  </div>
                </LightBox.View>
              </LightBox>

              <LightBox ref="editAggregator" closeOnClickMask closeOnESC>
                <LightBox.View>
                  <p>Edit aggregator</p>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Numerator</p>
                    <Input name="numerator" val={grpToEdit.numerator} class={[s.input]} ref="editNumerator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Denominator</p>
                    <Input name="denominator" val={grpToEdit.denominator} class={[s.input]} ref="editDenominator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Endpoint</p>
                    <Input name="endpoint" val={grpToEdit.endpoint} class={[s.input]} ref="editEndpoint" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Metric</p>
                    <Input name="metric" val={grpToEdit.metric} class={[s.input]} ref="editMetric" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Tags</p>
                    <Input name="tags" val={grpToEdit.tags} class={[s.input]} ref="editTags" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Step</p>
                    <Input name="step" val={grpToEdit.step} class={[s.input]} ref="editStep" />
                  </div>
                  <div class={[s.lightboxFooter]}>
                    <LightBox.Close>
                      <Button status="primaryOutline">Cancel</Button>
                    </LightBox.Close>
                    <LightBox.Close>
                      <Button status="primary" class={[s.submitBtn]} nativeOn-click={editAggregator}>Submit</Button>
                    </LightBox.Close>
                  </div>
                </LightBox.View>
              </LightBox>

              <LightBox ref="duplicateAggregator" closeOnClickMask closeOnESC>
                <LightBox.View>
                  <p>Duplicate aggregator</p>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Numerator</p>
                    <Input name="numerator" val={grpToEdit.numerator} class={[s.input]} ref="dupliNumerator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Denominator</p>
                    <Input name="denominator" val={grpToEdit.denominator} class={[s.input]} ref="dupliDenominator" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Endpoint</p>
                    <Input name="endpoint" val={grpToEdit.endpoint} class={[s.input]} ref="dupliEndpoint" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Metric</p>
                    <Input name="metric" val={grpToEdit.metric} class={[s.input]} ref="dupliMetric" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Tags</p>
                    <Input name="tags" val={grpToEdit.tags} class={[s.input]} ref="dupliTags" />
                  </div>
                  <div class={[s.inputGroup]}>
                    <p class={[s.formTitle]}>Step</p>
                    <Input name="step" val={grpToEdit.step} class={[s.input]} ref="dupliStep" />
                  </div>
                  <div class={[s.lightboxFooter]}>
                    <LightBox.Close>
                      <Button status="primaryOutline">Cancel</Button>
                    </LightBox.Close>
                    <LightBox.Close>
                      <Button status="primary" class={[s.submitBtn]} nativeOn-click={duplicateAggregator}>Submit</Button>
                    </LightBox.Close>
                  </div>
                </LightBox.View>
              </LightBox>

              <LightBox ref="deleteAggregator" closeOnClickMask closeOnESC>
               <LightBox.View>
                 <p>Delete aggregator</p>
                 <p class={[s.deleteDes]}>You will remove this aggregator: <b>{this.grpToRemove}</b>. Are you sure？</p>
                 <div class={[s.buttonGroup]}>
                   <LightBox.Close class={[s.btnWrapper]}>
                     <Button status="primary" class={[s.buttonBig]} nativeOn-click={() => deleteAggregator(this.grpToRemove)}>Yes</Button>
                   </LightBox.Close>
                   <LightBox.Close class={[s.btnWrapper]}>
                     <Button status="primaryOutline" class={[s.buttonBig]}>Cancel</Button>
                   </LightBox.Close>
                 </div>
               </LightBox.View>
             </LightBox>

            </div>
            <Grid { ...{ props } } />
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}
module.exports = AggregatorPage
