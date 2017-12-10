import {connect} from 'cerebral/inferno'
import {state, props, signal} from 'cerebral/tags'
import {field} from '@cerebral/forms'
import h from 'inferno-hyperscript'

export default connect({
  field: field(state`${props`path`}`),
  fieldChanged: signal`forms.fieldChanged`
},
  function Checkbox({name, legendtext, path, field, fieldChanged, labeltext}) {
    function onChanged(e) {
      fieldChanged({
        field: path,
        value: e.target.checked
      })
    }
    
    return h('div.field', [
            h('div.control', [
              h('legend.label', [legendtext]),
              h('input', {id: name, type: 'checkbox', name: name, onChange: (e) => onChanged(e)}),
              h('label', {for: name}, [
                h('span.fa-stack', [
                  h('i.fa.fa-square-o.fa-stack-1x'),
                  h('i.fa.fa-check-square.fa-stack-1x')
                ]),
                labeltext
              ])
            ])
          ])
})
