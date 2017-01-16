import fixUnit from '~utils/fix-unit'

const ECharts = {
  name: 'ECharts',
  props: {
    options: {
      type: Object,
    },
    size: {
      type: [Number, String],
      default: '600x375',
    },
  },

  mounted() {
    const { $el, options, _events } = this

    this.Chart = window.echarts.init($el)
    this.Chart.setOption(options)

    Object.keys(_events).forEach((ev) => {
      this.Chart.on(ev, (params) => {
        this.$emit(ev, params)
      })
    })
  },

  watch: {
    options(newOptions) {
      this.Chart.setOption(newOptions)
    }
  },

  computed: {
    style() {
      const { size } = this
      let width, height

      if (typeof size === 'number') {
        width = height = size
      } else {
        [width, height] = size.split('x')

        width = fixUnit(width)
        height = fixUnit(height)
      }

      return {
        width, height
      }
    }
  },

  render(h) {
    const { style } = this
    return (
      <div style={style} />
    )
  }
}

module.exports = ECharts
