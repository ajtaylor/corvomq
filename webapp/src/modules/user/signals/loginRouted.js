import {set} from 'cerebral/operators'
import {state} from 'cerebral/tags'

export default [
  set(state`user.isLoggedIn`, false),
  set(state`app.currentPage`, 'login')
]