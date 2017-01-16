// Portal container
import { Tab } from '@cepave/owl-ui'
import HostGroups from './host-groups'
import s from './portal.scss'

const Portal = {
  name: 'Portal',


  render(h) {
    return (
      <div class={[s.hostGroupPage]}>
        <Tab>
          <Tab.Head slot="tabHead" isSelected={true} name="1">HostGroups</Tab.Head>
          <Tab.Content slot="tabContent" name="1">
            <HostGroups />
          </Tab.Content>
        </Tab>
      </div>
    )
  }
}

module.exports = Portal
