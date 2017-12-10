import {set} from 'cerebral/operators'
import {state, props} from 'cerebral/tags'
import {goTo} from '@cerebral/router/operators'
import createUser from '../actions/createUser'
import createOrganisation from '../actions/createOrganisation'

const signupWithFirebase = [
  createUser, {
    success: [
      // set(state`user.isLoggedIn`, true),
      // set(state`user.isLoggingIn`, false),
      // goTo('/environments')
      set(state`user.signupForm.password.value`, ""),
      set(state`user.email`, props`newUser.email`),
      set(state`user.uid`, props`newUser.uid`),
      createOrganisation, {
        success: [],
        error: []
      }
    ],
    error: [
      // set(state`app.isLoggingIn`, false),
      // set(state`user.loginForm.password.value`, '')
    ]
  }
]

export default signupWithFirebase
