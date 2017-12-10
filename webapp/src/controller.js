import {Controller} from 'cerebral'
import {state} from 'cerebral/tags'
import Devtools from 'cerebral/devtools'
import Router from '@cerebral/router'
import FormsProvider from '@cerebral/forms'
import HttpProvider from '@cerebral/http'
import AppModule from './modules/app'
import UserModule from './modules/user'
import EnvironmentsModule from './modules/environments'
import FormsModule from './modules/forms'

const router = Router({
  routes: [{
    path: '/',
    signal: 'user.loginRouted'
  }, {
  //   path: '/login',
  //   signal: 'user.loginRouted'
  // }, {
    path: '/environments',
    signal: 'environments.routed'
  }, {
    path: '/environments/create',
    signal: 'environments.createRouted'
  }]
})

const controller = Controller({
  devtools: Devtools({
    host: 'localhost:8585'
  }),
  modules: {
    app: AppModule,
    user: UserModule,
    environments: EnvironmentsModule,
    forms: FormsModule,
    router: router
  },
  providers: [
    HttpProvider({
      baseUrl: 'http://api.corvomq.com',
      headers: {'Content-Type': 'application/json; charset=UTF-8'}
    }),
    FormsProvider()]
})

export default controller
