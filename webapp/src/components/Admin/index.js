import {connect} from 'cerebral/inferno'
import {state} from 'cerebral/tags'
import h from 'inferno-hyperscript'
import NavBar from '../NavBar'
import Environments from '../Environments'
import CreateEnvironment from '../CreateEnvironment'

const pages = {
  environments: {
    showNavbar: true,
    component: Environments
  },
  environments_create: {
    showNavbar: true,
    component: CreateEnvironment
  }
}

export default connect({
  currentPage: state`app.currentPage`
},
  function Admin({currentPage}) {
    const Page = pages[currentPage]

    return (
      h('div.container', [
        Page.showNavbar ?
          h(NavBar) : null,
        h('div.container', [
          h(Page.component)
        ])
      ])
    )
  }
)
