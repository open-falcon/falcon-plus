import { Input, Button } from '@cepave/owl-ui'
import s from './signup.scss'

const SignupPage = {
  name: 'SignupPage',
  methods: {
    signup() {
      if (this.$refs.pwd.value !== this.$refs.confirmPwd.value) {
        this.$store.dispatch('signup/pwdNotMatch')
        return
      } else if (this.$refs.username.value && this.$refs.email.value && this.$refs.pwd.value && this.$refs.nickname.value) {
        this.$store.dispatch('signup/signup', {
          name: this.$refs.username.value,
          email: this.$refs.email.value,
          password: this.$refs.pwd.value,
          cnname: this.$refs.nickname.value,
          router: this.$router
        })
      }
    }
  },
  computed: {
    notiStyle() {
      const { notification, status } = this.$store.state.signup
      const notiStyle = (!notification)
                        ? [s.noNotificaton]
                        : (notification && status)
                        ? [s.notiSuccess]
                        : [s.notiError]

      return notiStyle
    }
  },
  render(h) {
    const { signup, notiStyle } = this
    const { notification } = this.$store.state.signup
    return (
      <div class={[s.container]}>
        <div class={[s.wrapper]}>
          <div class={[s.inputWrapper]}>
            <h2 class={[s.title]}>Sign Up</h2>
            <div class={notiStyle}>{notification}</div>
            <Input placeholder="username" class={[s.input]} ref="username" />
            <Input placeholder="nickname" class={[s.input]} ref="nickname" />
            <Input placeholder="password" password={true} class={[s.input]} ref="pwd" />
            <Input placeholder="confirm password" password={true} class={[s.input]} ref="confirmPwd" />
            <Input placeholder="email" class={[s.input]} ref="email" />
            <div class={[s.buttons]}>
              <Button status="outline">not now</Button>
              <Button status="primary" nativeOn-click={signup}>sign up</Button>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

module.exports = SignupPage
