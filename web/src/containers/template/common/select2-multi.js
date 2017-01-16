import './select2.scss'
const { _, $ } = window

// this is only working for contact select mutiple item,
// if you want use this you need copy this as a new one and change those code below.
const Select2Muti =  {
  name: 'Select2Muti',
  props: {
    options: { type: Array, default: [] },
    value: {},
    class: { type: String, default: 'myselect' },
    setclass: { type: String, default: 'myselect_back' },
  },
  data() {
    return {
      select2Base: {
        placeholder: '请选择',
        tags: true,
        tokenSeparators: [',']
      }
    }
  },
  mounted() {
    $(this.$el)
      .select2({ data: this.options, ...this.select2Base })
  },
  methods: {
    findValueIndex(value) {
      const items = value.split(',')
      const res = _.chain(this.options).filter((m) => {
        return _.includes(items, m.text)
      }).map((m2) => {
        return m2.id
      }).value()
      return res
    }
  },
  watch: {
    value(value) {
      const itmes = this.findValueIndex(value)
      $(this.$el).val(itmes).trigger('change')
    },
    options(options) {
      // update options
      $(this.$el)
        .select2({ data: this.options, ...this.select2Base })
      $(this.$el).on('change', function (e) {
        if (this.value === 0 || this.value === '0') {
          $(`.${this.getAttribute('back')}`).val('')
        } else {
          $(`.${this.getAttribute('back')}`).val(
            _.map(this.selectedOptions, (m) => {
              return m.text
            }).join(','))
        }
      })
      const itmes = this.findValueIndex(this.value)
      $(this.$el).val(itmes).trigger('change')
    }
  },
  destroyed() {
    $(this.$el).off().select2('destroy')
  },
  render(h) {
    return (
      <select class={this.class} back={this.setclass} multiple="multiple">
        {this.options.map((option) => {
          return (<option value={option.id}> { option.text } </option>)
        })}
      </select>
    )
  }
}

module.exports = Select2Muti
