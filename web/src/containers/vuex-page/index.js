// VuexPage

const VuexPage = {
  name: 'VuexPage',

  methods: {
    enterHandler(e) {
      this.$store.dispatch('changeName', e.target.value)
    },
  },

  render(h) {
    const { state } = this.$store
    const { enterHandler } = this
    return (
      <div>
        <p>Here is Vuex Sample Page.</p>
        <input type="text" placeholder="What is your name?" on-input={enterHandler} />
        <p>{state.sample.name}</p>
      </div>
    )
  }
}

module.exports = VuexPage
