import {set} from 'cerebral/operators'
import {state, props, string} from 'cerebral/tags'
import {goTo} from '@cerebral/router/operators'
import {httpPost} from '@cerebral/http/operators'
import login from '../actions/login'
import updateHttpOptions from '../actions/updateHttpOptions'

// const loginWithFirebase = [
//   login, {
//     success: [
//       set(state`user.isLoggedIn`, true),
//       set(state`user.isLoggingIn`, false),
//       goTo('/environments')
//     ],
//     error: [
//       set(state`app.isLoggingIn`, false),
//       set(state`user.loginForm.password.value`, "")
//     ]
//   }
// ]

// export default loginWithFirebase

export default [
  set(state`user.isLoggingIn`, true),
  login, {
    success: [
      set(state`user.isLoggedIn`, true),
      set(state`user.isLoggingIn`, false),
      set(state`user.loginForm.password.value`, ""),
      updateHttpOptions,
      goTo('/environments')
    ],
    error: [
      set(state`user.isLoggingIn`, false),
      set(state`user.loginForm.password.value`, "")
    ]
  }
]
