import {connect} from 'cerebral/inferno'
import {state, signal} from 'cerebral/tags'
import {form} from '@cerebral/forms'
import Input from '../Forms/Input'
import h from 'inferno-hyperscript'

export default connect({
  form: form(state`user.signupForm`),
  formSubmitted: signal`user.signupSubmitted`
},
  function Signup({formSubmitted}) {
    function onSubmitted(e) {
      e.preventDefault()
      e.stopPropagation()
      formSubmitted()
    }
    return h('container', [
      h('div.measure-narrow.center.pt4', [
        h('h2.f2.fw4.tc.dark-blue', ['CorvoMQ']),
        h('h3.f3.fw3.tc', ['Signup']),
        h('form', {onSubmit: (e) => onSubmitted(e)}, [
          h('.mt3', [
            h(Input, {name: 'organisationName', path: 'user.signupForm.organisationName', placeholder: 'Organisation name'})
          ]),
          h('.mt3', [
            h(Input, {name: 'emailAddress', path: 'user.signupForm.emailAddress', placeholder: 'Email address'})
          ]),
          h('.mt3', [
            h(Input, {name: 'password', path: 'user.signupForm.password', placeholder: 'Password', password: true})
          ]),
          h('.mt3', [
            h('button.br1.bw0.bg-blue.pa3.dim.dib.white.w-100.pointer.ttu.tracked', {type: 'submit'}, ['Signup'])
          ])
        ])
      ])
    ])
})
