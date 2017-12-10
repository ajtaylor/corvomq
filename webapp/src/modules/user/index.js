import {state} from 'cerebral/tags'
import loginRouted from './signals/loginRouted'
import loginSubmitted from './signals/loginSubmitted'

export default {
  state: {
    isLoggedIn: false,
    isLoggingIn: false,
    loginForm: {
      emailAddress: {
        value: 'antony.taylor@gmail.com',
        isRequired: true
      },
      password: {
        value: 'EasyPass1',
        isRequired: true
      }
    }
  },
  signals: {
    loginRouted: loginRouted,
    loginSubmitted: loginSubmitted
  }
}