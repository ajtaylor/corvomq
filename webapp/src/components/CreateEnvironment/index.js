import {connect} from 'cerebral/inferno'
import {state, signal} from 'cerebral/tags'
import {Textbox, RadioList, Checkbox} from '../Forms'
import h from 'inferno-hyperscript'

export default connect({
  formSubmitted: signal`environments.createSubmitted`
},
  function CreateEnvironment({formSubmitted}) {
    function onSubmitted(e) {
      e.preventDefault()
      e.stopPropagation()
      formSubmitted()
    }
    return (
      h('div.container', [
        h('h4.is-size-4', ['Create Environment']),
        h('div.columns', [
          h('div.column.is-half', [
            h('form', {onSubmit: (e) => onSubmitted(e)}, [
              h(Textbox, {name: 'name',
                          path: 'environments.createForm.name',
                          placeholder: 'Environment name',
                          helptext: 'Required, must be unique',
                          labeltext: 'Name',
                          maxWidth: 6//,
                          // labelsubtext: 'Required'
                        }),
              // h('div.cf', [
                // h('div.fl.w-50', [
                  h(RadioList, {legendtext: 'Infrastructure',
                                name: 'infrastructure',
                                path: 'environments.createForm.infrastructure',
                                options: [ { name: 'cloud', label: 'Cloud' },
                                            { name: 'standalone_failover', label: 'Standalone with failover' },
                                            { name: 'standalone', label: 'Standalone' } ]
                              })
                // ])
                ,
                // h('div.fl.w-50', [
                  h(RadioList, {legendtext: 'Server',
                                name: 'server',
                                path: 'environments.createForm.server',
                                options: [ { name: 'message', label: 'Message' },
                                            { name: 'queue', label: 'Queue' } ]
                              })
                // ])
              // ])
              ,
              h(Checkbox, {
                name: 'tls_enabled',
                legendtext: 'Encryption',
                path: 'environments.createForm.tls_enabled',
                labeltext: 'Enable TLS'
              }),
              h('div.field', [
                h('div.control', [
                  h('button.button.is-success', {type: 'submit'}, ['Create environment'])
                ])
              ])
            ])
          ]),
          h('div.column.is-half', [
            h('div.box', [
              h('strong.is-size-6', ['Name']),
              h('p', ['You must provide a name for each environment you create. You cannot have two environments with the same name.']),
              h('strong.is-size-6', ['Infrastructure']),
              h('ul', [
                h('li', ['Cloud']),
                h('li', ['Standalone with failover']),
                h('li', ['Standalone'])
              ]),
              h('strong.is-size-6', ['Server']),
              h('ul', [
                h('li', ['Message']),
                h('li', ['Queue'])
              ]),
              h('strong.is-size-6', ['Encryption']),
              h('p', ['TLS info here.'])
            ])
          ])
        ])
      ])
    )
  }
)
