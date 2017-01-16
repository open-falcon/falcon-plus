import { Tab } from '@cepave/owl-ui'
import s from './user.scss'
import UserGroup from './user-group'
import UserList from './user-list'

const User = {
  name: 'User',
  mounted() {
    this.$store.dispatch('userGroup/getUsers')
  },
  render(h) {
    return (
      <div class={[s.contactPage]}>
        <Tab>
          <Tab.Head slot="tabHead" name="userGroup" isSelected>User Group</Tab.Head>
          <Tab.Content slot="tabContent" name="userGroup">
            <UserGroup />
          </Tab.Content>
          <Tab.Head slot="tabHead" name="userList">User List</Tab.Head>
          <Tab.Content slot="tabContent" name="userList">
            <UserList />
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}

module.exports = User
