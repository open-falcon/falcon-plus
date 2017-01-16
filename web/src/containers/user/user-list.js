import { Grid, Input, Button, LightBox } from '@cepave/owl-ui'
import s from './user.scss'
const UserList = {
  name: 'UserList',
  data() {
    return {
      userListData: {
        rowsRender() {},
      },
      heads: [
        {
          col: 'Name',
          width: '15.3%',
          sort: -1
        },
        {
          col: 'Nickname',
          width: '15.3%',
          sort: -1
        },
        {
          col: 'E-mail',
          width: '23%',
          sort: -1
        },
        {
          col: 'Cellphone',
          width: '15.3%',
          sort: -1
        },
        {
          col: 'IM',
          width: '15.3%',
          sort: -1
        },
        {
          col: 'QQ',
          width: '15.3%',
          sort: -1
        }
      ]
    }
  },
  created() {
    this.userListData.rowsRender = (h, { row, index }) => {
      const result = [
        <Grid.Col>{row.name}</Grid.Col>,
        <Grid.Col>{row.cnname}</Grid.Col>,
        <Grid.Col>{row.email}</Grid.Col>,
        <Grid.Col>{row.phone}</Grid.Col>,
        <Grid.Col>{row.qq}</Grid.Col>,
        <Grid.Col>{row.im}</Grid.Col>
      ]
      return result
    }
  },
  methods: {
    searchUser(e) {
      if ((e.type === 'keypress' && e.charCode === 13) || e.type === 'click') {
        this.$store.dispatch('userGroup/searchUser', {
          q: this.$refs.searchUser.value
        })
      }
    }
  },
  render(h) {
    const { userListData, searchUser, heads } = this
    const { userListRows } = this.$store.state.userGroup
    const props = { ...userListData, heads, rows: userListRows }
    return (
      <div>
        <div class={[s.contactSearchWrapper]}>
          <div class={[s.searchInputGroup]}>
            <Input icon={['search', '#b8bdbf']} placeholder="search user... search all: .+" ref="searchUser" nativeOn-keypress={searchUser} class={[s.contactSearch]}  />
            <Button status="primary" nativeOn-click={searchUser}>Search</Button>
          </div>
        </div>
        <div class={[s.contactWrapper]}>
          <div class={[s.gridWrapperBox]}>
            <Grid { ...{ props } } />
          </div>
        </div>
      </div>
    )
  }
}
module.exports = UserList
