import './select2.scss'
const { _, $ } = window

const Select2 =  {
  name: 'Select2',
  props: {
    setMetric: {
      type: Function
    },
    options: {
      type: Array,
      default: () => {
        return []
      }
    },
    value: {
      type: [String, Number]
    },
    class: {
      type: String,
      default: 'myselect'
    },
    setclass: {
      type: String,
      default: 'myselect_back'
    },
    accpetTag: {
      type: Boolean,
      default: false
    },
    placeholder: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      select2Base: {
        tags: this.accpetTag,
        placeholder: this.placeholder,
        allowClear: true
      }
    }
  },
  mounted() {
    $(this.$el)
      .select2({ data: this.options, ...this.select2Base })
  },
  watch: {
    value(value) {
      if (!_.includes(_.keys(this.options), this.value)) {
        $(this.$el).select2({ data: [{ id: this.id, text: this.value }] })
      }
      // update value
      $(this.$el).val(this.value).trigger('change')
    },
    options(options) {
      // update options
      $(this.$el)
        .select2({ data: this.options, ...this.select2Base })
      $(this.$el).on('change', function (e) {
        if (this.value === 0 || this.value === '0') {
          $(`.${this.getAttribute('back')}`).val('')
        } else {
          $(`.${this.getAttribute('back')}`).val(this.value)
        }
      })
    }
  },
  destroyed() {
    $(this.$el).off().select2('destroy')
  },
  render(h) {
    const { placeholder } = this
    return (
      <select class={this.class} back={this.setclass}>
        <option value="0" selected="selected" disabled>{placeholder}</option>
        {
          this.options.map((option) => {
            return <option value={option.id}>{option.text}</option>
          })
        }
      </select>
    )
  }
}

module.exports = Select2
