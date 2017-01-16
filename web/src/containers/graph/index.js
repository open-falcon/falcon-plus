import { Tab, DualList, Loading, Button, Page, Icon, Select, Flex } from '@cepave/owl-ui'
import LineChart from './line-chart/index.js'
import sortHosts from './sort-hosts'
import g from '~sass/global.scss'
import s from './graph.scss'
import delegate from 'delegate-to'
import DatePicker from '~coms/date-picker'

const GraphView = {
  name: 'GraphView',

  data() {
    return {
      samplingOptions: [
        { title: 'AVERAGE', value: 'AVERAGE' },
        { title: 'MAX', value: 'MAX' },
        { title: 'MIN', value: 'MIN' },
      ]
    }
  },
  mounted() {

  },

  methods: {
    getEndpoints(q) {
      const { $store, $refs } = this
      $store.dispatch('graph/getEndpoints', {
        q,
      })
    },

    getCounters(metricQuery) {
      const { $store, $refs } = this

      $store.dispatch('graph/getCounters', {
        eid: Object.keys($refs.dualEndpoint.rightList).map((k)=>{
          return $refs.dualEndpoint.rightList[k].id
        }),
        metricQuery
      })
    },

    switchViewPoint(d) {
      const { $store } = this

      $store.commit('graph/switchViewPoint', {
        viewpoint: d.value
      })
    },

    viewGraph() {
      const { $store, $refs, goToPage } = this
      if ($store.state.graph.viewGraphBtnDisabled) {
        return
      }

      const endpoints = Object.keys($refs.dualEndpoint.rightList).map((k) => $refs.dualEndpoint.rightList[k].endpoint)
      const counters = Object.keys($refs.dualCounter.rightList).map((k) => $refs.dualCounter.rightList[k].counter)
      const vport = $store.state.graph.vport

      let totalCharts
      if (vport === 'combo') {
        totalCharts = [
          {
            title: 'Combo',
            endpoints,
            counters,
          }
        ]
      } else {
        const vports = sortHosts({
          endpoints,
          counters,
        })[vport]

        totalCharts = Object.keys(vports).map((k) => {
          return {
            title: k,
            endpoints: vports[k].endpoints,
            counters: vports[k].counters,
          }
        })
      }
      $store.commit('graph/viewGraph.start', { totalCharts })
      goToPage(1)
    },

    goToPage(page) {
      const { $store, $refs } = this
      const graph = $store.state.graph
      const vport = graph.vport

      const sliceStart = (page - 1) * graph.pageLimit
      const sliceEnd = page * graph.pageLimit

      let start = graph.startTime
      let end = graph.endTime

      if (start > end) {
        start = graph.endTime
        end = graph.startTime
      }

      const charts = graph.totalCharts
        .slice(sliceStart, sliceEnd)
        .map(({ title, endpoints, counters }, idx) => {
          $store.dispatch('graph/viewGraph', {
            endpoints,
            counters,
            idx, vport, page,
            sampling: $refs.sampling.value,
            start, end
          })
          return {
            title,
            loading: true,
            series: [],
          }
        })
      $store.commit('graph/viewGraph.page', { charts, page })
    },

    switchGrid: delegate('[data-grid]', function(ev) {
      const { $store, $refs } = this
      const grid = ev.delegateTarget.dataset.grid

      $store.commit('graph/switchGrid', { grid })

      this.$nextTick(() => {
        ;[...Array($store.state.graph.pageLimit)].forEach((ref, i) => {
          if ($refs[`chart${i}`].Chart) {
            $refs[`chart${i}`].Chart.resize()
          }
        })
      })
    }),

    checkBtnStatus(d) {
      const { $store, $refs } = this

      $store.commit('graph/checkViewGraphBtnStatus', {
        disabled: !(Object.keys($refs.dualEndpoint.rightList).length &&
          Object.keys($refs.dualCounter.rightList).length)
      })
    },

    onChangeStartTime({ unix }) {
      const { $store } = this

      $store.commit('graph/syncStartTime', { unix })
    },

    onChangeEndTime({ unix }) {
      const { $store } = this

      $store.commit('graph/syncEndTime', { unix })
    }
  },

  computed: {
    pageSum() {
      const { $store } = this
      const graph = $store.state.graph
      const start = graph.pageLimit * (graph.pageCurrent - 1) + 1
      const end = Math.min(graph.pageCurrent *  graph.pageLimit, graph.totalCharts.length)

      return `${start} - ${end}`
    }
  },

  render(h) {
    const { $store, $router, getEndpoints, getCounters, switchViewPoint, viewGraph, goToPage,
      pageSum, switchGrid, checkBtnStatus, samplingOptions, onChangeStartTime, onChangeEndTime } = this
    const { graph } = $store.state

    return (
      <div>
        <Tab class={[s.tab]}>
          <Tab.Head slot="tabHead" isSelected name="graph">Graph</Tab.Head>
          <Tab.Content slot="tabContent" name="graph">
            <Flex split class={[s.topCtrl]}>
              <Flex.Col>
                <Flex>
                  <Flex.Col size="auto">
                    <Flex mid>
                      Start <DatePicker ref="startTime" onChange={onChangeStartTime} readOnly height={34} class={[s.picker]} opts={{
                        max: new Date(),
                        initialValue: new Date(graph.startTime * 1000),
                      }} />
                    </Flex>
                  </Flex.Col>
                  <Flex.Col size="auto">
                    <Flex mid>
                      End <DatePicker ref="endTime" onChange={onChangeEndTime} readOnly height={34} class={[s.picker]} opts={{
                        max: new Date(),
                        initialValue: new Date(graph.endTime * 1000),
                      }} />
                    </Flex>
                  </Flex.Col>
                  <Flex.Col size="auto">
                    <Flex mid>
                      Sampling
                      <Select class={[s.sampling]} ref="sampling" options={samplingOptions} />
                    </Flex>
                  </Flex.Col>
                </Flex>
              </Flex.Col>

              <Flex.Col>
                <Button disabled={graph.viewGraphBtnDisabled} status="primary" nativeOnClick={viewGraph}>View Graph</Button>
              </Flex.Col>
            </Flex>

            <Flex class={[s.submitSection]}>
              <span class={[s.vpoint]}>
                View Point
              </span>
              <Button.Group onChange={switchViewPoint} options={[
                 { value: 'endpoint', title: 'Endpoint', selected: true },
                 { value: 'counter', title: 'Counter' },
                 { value: 'combo', title: 'Combo' },
              ]} />
            </Flex>

            <Flex class={[s.querySection]}>
              <Flex.Col size="6">
                <h4>Search Endpoints</h4>
                <DualList
                  ref="dualEndpoint"
                  leftLoading={graph.hasEndpointLoading}
                  apiMode
                  displayKey="endpoint"
                  items={graph.endpointItems}
                  onChange={checkBtnStatus}
                  onInputchange={getEndpoints}
                />
              </Flex.Col>
              <Flex.Col size="6">
                <h4>Search Counters</h4>
                <DualList
                  apiMode
                  onChange={checkBtnStatus}
                  ref="dualCounter"
                  displayKey="counter"
                  onInputchange={getCounters}
                  items={graph.counterItems}
                  leftLoading={graph.hasCounterLoading}
                />
              </Flex.Col>
            </Flex>

            {
              graph.totalCharts.length
              ? <div class={[s.pageSum, g.flexSplit]} data-grid={graph.grid}>
                  <div>
                    Total Charts: { pageSum } of { graph.totalCharts.length }
                  </div>
                  <div class={[s.chartGrid]} onClick={switchGrid}>
                    <Icon typ="grid-1" data-grid="1" />
                    <Icon typ="grid-4" data-grid="2" />
                    <Icon typ="grid-9" data-grid="3" />
                    <Icon typ="grid-16" data-grid="4" />
                  </div>
                </div>
              : null
            }
            <div class={[s.boxCharts]} data-grid={graph.grid}>
              {graph.charts.map((chart, i) => {
                return (
                  <div class={[s.chart]}>
                    <LineChart ref={`chart${i}`} title={chart.title} loading={chart.loading} series={chart.series} />
                  </div>
                )
              })}
            </div>
            { graph.totalCharts.length ? <Page class={[s.pager]} onPage={({ page }) => {
              goToPage(page)
            }} total={graph.totalCharts.length} limit={graph.pageLimit} /> : null }
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}

module.exports = GraphView
