import s from './user.scss'
import { Grid, Input, Button, Icon, LightBox, DualList, Loading } from '@cepave/owl-ui'

const UserGroup = {
  name: 'UserGroup',
  data() {
    return {
      heads: [
        {
          col: 'User group name',
          width: '30.7%',
          sort: -1
        },
        {
          col: 'Users',
          width: '32.3%',
          sort: -1
        },
        {
          col: 'Creator',
          width: '16.6%',
          sort: -1
        },
        {
          col: 'Opeartions',
          width: '20.4%'
        }
      ],
      userGroupData: {
        rowsRender() {},
      },
      selectedUsers: [],
      editTeamId: 0,
      editTeamName: '',
      dupliTeamId: 0,
      dupliTeamName: '',
      deleteTeamId: 0,
      deleteTeamName: '',
      editUsers: [],
    }
  },
  mounted() {
    this.$store.dispatch('userGroup/getUserGroups')
  },
  methods: {
    createUserGroup() {
      this.$store.dispatch('userGroup/newUserGroup', {
        team_name: this.$refs.newUserGroupName.value,
        resume: 'some description',
        users: Object.keys(this.selectedUsers).reduce((preVal, curVal) => {
          preVal.push(this.selectedUsers[curVal].id)
          return preVal
        }, [])
      })
      .then(() => {
        this.$refs.newUserGroupName.val = ''
      })
    },
    deleteUserGroup(id) {
      this.$store.dispatch('userGroup/deleteUserGroup', {
        id
      })
    },
    editUserGroup(id) {
      this.$store.dispatch('userGroup/editUserGroup', {
        team_name: this.$refs.editGroupName.value,
        team_id: id,
        resume: '',
        users: Object.keys(this.editUsers).reduce((preVal, curVal) => {
          preVal.push(this.editUsers[curVal].id)
          return preVal
        }, [])
      })
    },
    duplicateUserGroup(teamName) {
      this.$store.dispatch('userGroup/newUserGroup', {
        team_name: this.$refs.duplicateGroupName.value,
        resume: 'some description',
        users: Object.keys(this.editUsers).reduce((preVal, curVal) => {
          preVal.push(this.editUsers[curVal].id)
          return preVal
        }, [])
      })
    },
    searchUserGroup(e) {
      if ((e.type === 'keypress' && e.charCode === 13) || e.type === 'click') {
        this.$store.dispatch('userGroup/searchGroup', {
          q: this.$refs.searchGroup.value
        })
      }
    },
    getOneTeam(id) {
      this.$store.dispatch('userGroup/getOneTeam', {
        team_id: id,
      })
    },
    handleNewUser(data) {
      this.selectedUsers = data
    },
    handleEditUsers(data) {
      this.editUsers = data
    },
    edit(e, teamId, teamName) {
      this.editTeamId = teamId
      this.editTeamName = teamName
      this.$refs.editGroup.open(e)
      this.getOneTeam(teamId)
    },
    duplicate(e, teamId, teamName) {
      this.dupliTeamId = teamId
      this.dupliTeamName = teamName
      this.$refs.duplicateGroup.open(e)
      this.getOneTeam(teamId)
    },
    deleteTeam(e, teamId, teamName) {
      this.deleteTeamId = teamId
      this.deleteTeamName = teamName
      this.$refs.deleteGroup.open(e)
    }
  },
  created() {
    this.userGroupData.rowsRender = (h, { row, index }) => {
      return [
        <Grid.Col>{row.groupName}</Grid.Col>,
        <Grid.Col>{row.groupMember}</Grid.Col>,
        <Grid.Col>{row.creator}</Grid.Col>,
        <Grid.Col>
          <div class={[s.opeartionInline]}>
            <span class={[s.opeartions]} on-click={(e) => this.edit(e, row.id, row.groupName)}>edit</span>
            <span class={[s.opeartions]} on-click={(e) => this.duplicate(e, row.id, row.groupName)}>Duplicate</span>
            <span class={[s.opeartions]} on-click={(e) => this.deleteTeam(e, row.id, row.groupName)}>Delete</span>
          </div>
        </Grid.Col>
      ]
    }
  },
  render(h) {
    const { userGroupData, handleNewUser, handleEditUsers, createUserGroup,
            searchUserGroup, editUserGroup, deleteUserGroup, duplicateUserGroup, heads, getSingleTeamLoading } = this
    const { userGroupHeads, rows, userListRows, singleTeamUsers, singleTeamUsersToSelect } = this.$store.state.userGroup
    const props = { ...userGroupData,  heads, ...{ rows } }
    return (
      <div>
        <div class={[s.contactSearchWrapper]}>
          <div class={[s.searchInputGroup]}>
            <Input icon={['search', '#b8bdbf']} placeholder="user group name... search all:.+" ref="searchGroup" nativeOn-keypress={searchUserGroup} class={[s.contactSearch]} />
            <Button status="primary" nativeOn-click={searchUserGroup}>Search</Button>
          </div>
          <LightBox ref="createGroup" closeOnClickMask closeOnESC>
            <LightBox.Open>
              <Button status="primary" class={[s.buttonIcon]}>
                <Icon typ="plus" size={16} class={[s.plus]} />
                Add new user group
              </Button>
            </LightBox.Open>
            <LightBox.View>
              <p>Add user group</p>
              <div class={[s.inputGroup]}>
                <p>Name</p>
                <Input placeholder="user group name" class={[s.input]} ref="newUserGroupName" />
              </div>
              <div class={[s.userLists]}>
                <p>Users</p>
                <DualList class={[s.groupSelect]} items={userListRows} displayKey="name" onChange={handleNewUser} />
              </div>
              <div class={[s.lightboxFooter]}>
                <LightBox.Close>
                  <Button status="primaryOutline">Cancel</Button>
                </LightBox.Close>
                <LightBox.Close>
                  <Button status="primary" class={[s.submitBtn]} nativeOn-click={createUserGroup}>Submit</Button>
                </LightBox.Close>
              </div>
            </LightBox.View>
          </LightBox>
        </div>
        <div class={[s.contactWrapper]}>
          <div class={[s.gridWrapperBox]}>
            <Grid { ...{ props } } />
          </div>
        </div>
        <LightBox ref="editGroup" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>Edit user group</p>
            <div class={[s.inputGroup]}>
              <p>Name</p>
              <Input class={[s.input]} ref="editGroupName" val={this.editTeamName} />
            </div>
            <div class={[s.userLists]}>
              <p>Users</p>
              <DualList leftLoading={getSingleTeamLoading} class={[s.groupSelect]} items={singleTeamUsersToSelect} selectedItems={singleTeamUsers} displayKey="name" onChange={handleEditUsers} />
            </div>
            <div class={[s.lightboxFooter]}>
              <LightBox.Close>
                <Button status="primaryOutline">Cancel</Button>
              </LightBox.Close>
              <LightBox.Close>
                <Button status="primary" class={[s.submitBtn]} nativeOn-click={() => editUserGroup(this.editTeamId)}>Submit</Button>
              </LightBox.Close>
            </div>
          </LightBox.View>
        </LightBox>

        <LightBox ref="duplicateGroup" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>Edit user group</p>
            <div class={[s.inputGroup]}>
              <p>Name</p>
              <Input placeholder="user group name" ref="duplicateGroupName" class={[s.input]} val={`${this.dupliTeamName}_1`} />
            </div>
            <div class={[s.userLists]}>
              <p>Users</p>
              <DualList class={[s.groupSelect]} items={singleTeamUsersToSelect} selectedItems={singleTeamUsers} displayKey="name" onChange={handleEditUsers} />
            </div>
            <div class={[s.lightboxFooter]}>
              <LightBox.Close>
                <Button status="primaryOutline">Cancel</Button>
              </LightBox.Close>
              <LightBox.Close>
                <Button status="primary" class={[s.submitBtn]} nativeOn-click={() => duplicateUserGroup(this.dupliTeamName)}>Submit</Button>
              </LightBox.Close>
            </div>
          </LightBox.View>
        </LightBox>

        <LightBox ref="deleteGroup" closeOnClickMask closeOnESC>
          <LightBox.View>
            <p>Delete user group</p>
            <p class={[s.deleteDes]}>You will remove this user group: <b>{this.deleteTeamName}</b>. Are you sure？</p>
            <div class={[s.buttonGroup]}>
              <LightBox.Close class={[s.btnWrapper]}>
                <Button status="primary" class={[s.buttonBig]} nativeOn-click={() => deleteUserGroup(this.deleteTeamId)}>Yes</Button>
              </LightBox.Close>
              <LightBox.Close class={[s.btnWrapper]}>
                <Button status="primaryOutline" class={[s.buttonBig]}>Cancel</Button>
              </LightBox.Close>
            </div>
          </LightBox.View>
        </LightBox>
      </div>
    )
  }
}
module.exports = UserGroup
