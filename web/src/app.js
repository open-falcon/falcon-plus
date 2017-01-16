import s from './app.scss'
import Nav from '~coms/nav'
import checkToken from './utils/check-token'

module.exports = {
  name: 'App',
  mounted() {
    const { $router } = this
    checkToken()
      .then((hasToken) => {
        if (hasToken) {
          return
        }

        $router.push('/login')
      })
  },
  computed: {
    classes() {
      const { $route } = this
      return {
        [s.isNotPage]: /^\/(login|signup)$/.test($route.path)
      }
    }
  },
  render(h) {
    const { classes } = this

    return (
      <div id="app" class={[s.app, classes]}>
        <Nav />
        <div class={[s.main]}>
          <keep-alive>
            <router-view />
          </keep-alive>
        </div>
      </div>
    )
  }
}
