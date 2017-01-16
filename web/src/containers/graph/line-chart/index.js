import EChart from '~coms/echarts'
import s from './line-chart.scss'
import { Loading } from '@cepave/owl-ui'
import g from '~sass/global.scss'

const LineChart = {
  name: 'LineChart',
  props: {
    title: {
      type: String,
      default: ''
    },

    loading: true,

    series: {
      type: Array,
      default: () => [
        // {
        //   name: 'disk.io.await/device=sda',
        //   data: [
        //     [1482191540, 1],
        //     [1482191400, 3],
        //     [1482191700, 23],
        //   ]
        // },
      ],
    }
  },

  computed: {
    seriesExtend() {
      const { series } = this
      return series.map((s) => {
        return {
          ...s,
          type: 'line',
          areaStyle: { normal: {} },
        }
      })
    },

    echartsOptions() {
      const { seriesExtend } = this

      return {
        grid: {
          top: 20,
          left: '6%',
          right: '5%',
          bottom: 34,
        },
        tooltip: {
          trigger: 'axis',
        },
        xAxis: {
          type: 'time',
          axisLine: {
            show: false,
          },
          axisTick: {
            show: true,
          },
          splitLine: {
            show: false,
          },
        },
        yAxis: {
          type: 'value',
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              color: '#eee'
            }
          },
        },
        series: seriesExtend,
      }
    }
  },

  mounted() {

  },

  updated() {
    if (this.$refs.chart) {
      this.Chart = this.$refs.chart.Chart
    }
  },

  render(h) {
    const { title, echartsOptions, loading } = this

    return (
      <div class={[s.lineChart]}>
        <h4>{title}</h4>
        <div class={[g.ratio16x9]}>
          {
            loading
              ? <Loading class={[s.loading]} typ="bar" size="220x10" show />
              : <EChart ref="chart" size="100%x100%" options={echartsOptions} />
          }
        </div>
      </div>
    )
  }
}

module.exports = LineChart
