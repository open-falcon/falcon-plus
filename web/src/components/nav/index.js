import s from './nav.scss'
import { Icon } from '@cepave/owl-ui'
import Link from '~coms/link'

const Nav = {
  name: 'Nav',
  render(h) {
    return (
      <nav class={[s.nav]}>
        <div class={s.brand}>
          <Icon typ="brand-circle" size={40} />
        </div>

        {/* <div class={[s.linkIcon]}>
          <Link to="/alarm">
            <Icon typ="alarm-1" />
          </Link>
        </div> */}

        <div class={[s.linkIcon]}>
          <Link to="/graph">
            <Icon typ="linechart-1" />
          </Link>
        </div>

        <div class={[s.linkIcon]}>
          <Link to="/portal">
            <Icon typ="host-1" />
          </Link>
        </div>

        <div class={[s.linkIcon]}>
          <Link to="/user">
            <Icon typ="user-group-1" />
          </Link>
        </div>

        <div class={[s.linkIcon]}>
          <Link to="/template">
            <Icon typ="template" />
          </Link>
        </div>

        <div class={[s.navBottom]}>

          <div class={[s.linkIcon]}>
            <Link to="/profile">
              <Icon typ="person" />
            </Link>
          </div>

          <div class={[s.linkIcon, s.isBottom]}>
            <Icon typ="gear-1" />
          </div>

        </div>
      </nav>
    )
  }
}

module.exports = Nav
