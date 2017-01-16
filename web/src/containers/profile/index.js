import s from './profile.scss'
import { Tab, Input, Button } from '@cepave/owl-ui'

const ProfilePage = {
  name: 'ProfilePage',

  mounted() {
    this.$store.dispatch('profile/getProfile')
  },

  methods: {
    updateProfile() {
      const { userInfo } = this.$store.state.profile
      this.$store.dispatch('profile/updateProfile', {
        cnname: this.$refs.cnname.value,
        email: this.$refs.email.value,
        phone: this.$refs.phone.value,
        qq: this.$refs.qq.value,
        im: this.$refs.im.value
      })
    },
    updatePwd() {
      this.$store.dispatch('profile/updatePwd', {
        new_password: this.$refs.newPwd.value,
        old_password: this.$refs.curPwd.value
      })
    },
    logout() {
      this.$store.dispatch('profile/logout', {
        router: this.$router
      })
    }
  },

  computed: {
    notiStyle() {
      const { notification, status } = this.$store.state.profile
      const notiStyle = (!notification)
                        ? [s.noNotificaton]
                        : (notification && status)
                        ? [s.notiSuccess]
                        : [s.notiError]

      return notiStyle
    },
  },

  render(h) {
    const { notiStyle, updateProfile, updatePwd, logout } = this
    const { notification, userInfo } = this.$store.state.profile
    return (
      <div class={[s.profilePage]}>
        <Button status="primaryOutline" nativeOn-click={logout} class={[s.logout]}>Log out</Button>
        <Tab>
          <Tab.Head slot="tabHead" name="profile" isSelected={true}>Profile</Tab.Head>
          <Tab.Content slot="tabContent" name="profile">
            <div class={[s.box]}>
              <div>
                <div class={notiStyle}>{notification}</div>
                <p>Nickname</p>
                <Input name="cnname" val={userInfo.cnname} class={[s.input]} ref="cnname" />
                <p>E-Mail</p>
                <Input name="email" val={userInfo.email} class={[s.input]} ref="email" />
                <p>Cellphone</p>
                <Input name="cellphone" val={userInfo.phone} class={[s.input]} ref="phone" />
                <p>QQ</p>
                <Input name="qq" val={userInfo.qq} class={[s.input]} ref="qq" />
                <p>IM</p>
                <Input name="im" val={userInfo.im} class={[s.input]} ref="im" />
                <Button status="primary" class={[s.submit]} nativeOn-click={updateProfile}>Apply</Button>
              </div>
            </div>
          </Tab.Content>
          <Tab.Head slot="tabHead" name="password">Password</Tab.Head>
          <Tab.Content slot="tabContent" name="password">
            <div class={[s.box]}>
              <div>
                <div class={notiStyle}>{notification}</div>
                <p>Current Password</p>
                <Input name="currentPassword" password={true} class={[s.input]} ref="curPwd" />
                <p>New Password</p>
                <Input name="newPassword" password={true} class={[s.input]} ref="newPwd" />
                <Button status="primary" class={[s.submit]} nativeOn-click={updatePwd}>Apply</Button>
              </div>
            </div>
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}
module.exports = ProfilePage
