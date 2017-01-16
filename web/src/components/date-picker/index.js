import { Icon } from '@cepave/owl-ui'
import s from './date-picker.scss'
import fixUnit from '~utils/fix-unit'
const { rome, moment } = window

const DatePicker = {
  name: 'DatePicker',
  props: {
    width: {
      type: [String, Number],
      default: 168,
    },
    height: {
      type: [String, Number],
      default: 28,
    },

    readOnly: {
      type: Boolean,
      default: false,
    },
    opts: {
      type: Object,
    },
  },
  mounted() {
    const { opts } = this

    this.defaultOpts = {
      inputFormat: 'YYYY/MM/DD HH:mm',
      styles: {
        back: s.back,
        container: s.container,
        date: 'rd-date',
        dayBody: 'rd-days-body',
        dayBodyElem: s.day,
        dayConcealed: 'rd-day-concealed',
        dayDisabled: 'rd-day-disabled',
        dayHead: 'rd-days-head',
        dayHeadElem: 'rd-day-head',
        dayRow: 'rd-days-row',
        dayTable: 'rd-days',
        month: 'rd-month',
        next: s.next,
        positioned: 'rd-container-attachment',
        selectedDay: s.daySelected,
        selectedTime: s.timeSelected,
        time: 'rd-time',
        timeList: 'rd-time-list',
        timeOption: s.timeOption,
      },
    }

    this.Picker = rome(this.$refs.date, {
      ...this.defaultOpts,
      ...opts,
    })
  },

  methods: {
    showPicker() {
      return this.Picker.show()
    },

    handleChange(ev) {
      const input = ev.currentTarget
      if (!input.value) {
        return
      }

      this.value = {
        format: input.value,
        unix: moment(new Date(input.value)).unix(),
      }

      this.$emit('change', this.value)
    }
  },
  computed: {
    style() {
      const { width, height } = this
      return {
        width: fixUnit(width),
        height: fixUnit(height),
      }
    }
  },
  render(h) {
    const { style, handleChange, readOnly, showPicker } = this
    return (
      <div class={[s.picker]} style={style}>
        <input readOnly={readOnly} ref="date" onBlur={handleChange} />
        <Icon typ="date" size={18} class={[s.icon]} nativeOnMousedown={showPicker} nativeOnClick={(e) => {
          e.stopPropagation()
        }} />
      </div>
    )
  }
}

module.exports = DatePicker
