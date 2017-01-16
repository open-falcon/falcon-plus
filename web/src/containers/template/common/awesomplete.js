const { $, Awesomplete } = window

const AutoComplete = {
  name: 'AutoComplete',
  props: {
    options: { type: Array, default: [] },
    value: { type: String, default: '' },
    inputId: { type: String, default: 'myinput' },
    inputListId: { type: String, default: 'slist' },
    changeMetric: { type: Function },
  },
  data() {
    return {
      testop: ['cpu.idle', 'cpu.busy']
    }
  },
  mounted() {
    const self = this
    const input =  document.getElementById(this.inputId)
    new Awesomplete(input, { list: this.options })
    $(`#${this.inputId}`)
    .focusout(function(m) {
      self.changeMetric(m.target.value)
    })
  },
  updated() {
    const self = this
    const input = document.getElementById(this.inputId)
    new Awesomplete(input, { list: this.options })
    $(`#${this.inputId}`)
    .focusout(function(m) {
      self.changeMetric(m.target.value)
    })
  },
  methods: {
    catchMetric() {
      // this.changeMetric(this.$refs.vmetric.value)
      // return this.$refs.vmetric.value
    }
  },
  render(h) {
    return (
      <div>
        <input class="awesomplete" id="myinput" value={this.value} ref="vmetric" />
      </div>
    )
  }
}

module.exports = AutoComplete
