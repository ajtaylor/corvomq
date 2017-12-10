import {connect} from 'cerebral/inferno'
import {state} from 'cerebral/tags'
import h from 'inferno-hyperscript'
import Login from '../Login'
import Admin from '../Admin'
// import Signup from '../Signup'

// const pages = {
//   login: {
//     component: Login
//   },
//   signup: {
//     component: Signup
//   }
// }

export default connect({
  isLoggedIn: state`user.isLoggedIn`//,
  // currentPage: state`app.currentPage`
},
  function App({isLoggedIn, currentPage}) {
    // const Page = pages[currentPage]

    return (
      h('div.container', [
        isLoggedIn ?
          h(Admin) : h(Login)
          // h(Admin) : h(Page.component)
      ])
    )
  }
)
