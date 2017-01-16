import { Input, Button } from '@cepave/owl-ui'
import Logo from '~coms/logo'
import s from './login.scss'
const mapState = window.Vuex.mapState

const LoginPage = {
  name: 'LoginPage',
  mounted() {},
  data() {
    return {
      status: false
    }
  },
  methods: {
    login() {
      this.$store.dispatch('login/login', {
        name: this.$refs.username.value,
        password: this.$refs.password.value,
        router: this.$router,
      })
    },
    handleKeyPress(e) {
      if (e.charCode === 13 && this.$refs.username.value && this.$refs.password.value) {
        this.login()
      }
    }
  },
  computed: {
    notiStyle() {
      const { notification, status } = this.$store.state.login
      const notiStyle = (!notification)
                        ? [s.noNotificaton]
                        : (notification && status)
                        ? [s.notiSuccess]
                        : [s.notiError]

      return notiStyle
    },
  },
  render(h) {
    const { notiStyle, login, handleKeyPress } = this
    const { notification } = this.$store.state.login
    return (
      <div class={[s.container]}>
        <div class={[s.wrapper]}>
          <div class={[s.inputWrapper]}>
            <Logo size={100} class={[s.openfalcon]} />
            <div class={notiStyle}>{notification}</div>
            <Input placeholder="username" class={[s.input]} ref="username" />
            <Input placeholder="password" password={true} ref="password" nativeOn-keypress={handleKeyPress} />
            <Button status="primary" class={[s.buttons]} nativeOn-click={login}>log in</Button>
            <span class={[s.signupInfo]}>
              no acoount?
              <b><router-link to="/signup" class={[s.signup]}> Sign up </router-link></b>
              now
            </span>
          </div>
        </div>
      </div>
    )
  }
}

module.exports = LoginPage
