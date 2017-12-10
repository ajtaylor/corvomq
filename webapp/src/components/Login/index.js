import {connect} from 'cerebral/inferno'
import {state, signal} from 'cerebral/tags'
import {Textbox, Button} from '../Forms'
import h from 'inferno-hyperscript'

export default connect({
  formSubmitted: signal`user.loginSubmitted`
},
  function Login({formSubmitted}) {
    function onSubmitted(e) {
      e.preventDefault()
      e.stopPropagation()
      formSubmitted()
    }
    return h('div.container', [
      h('div.columns', [
        h('div.column.is-4.is-offset-4.has-text-centered', [
          h('h2.title.is-size-2', ['CorvoMQ']),
          h('h4.subtitle.is-size-4', ['Account login']),
          h('form', {onSubmit: (e) => onSubmitted(e)}, [
            h(Textbox, {name: 'emailAddress', path: 'user.loginForm.emailAddress', placeholder: 'Email address'}),
            h(Textbox, {name: 'password', path: 'user.loginForm.password', placeholder: 'Password', password: true}),
            h(Button, {type: 'submit', text: 'Login', isFullWidth: true})
          ])
        ])
      ])
    ])
})
