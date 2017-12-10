import {connect} from 'cerebral/inferno'
import {state, props, signal} from 'cerebral/tags'
import {field} from '@cerebral/forms'
import h from 'inferno-hyperscript'

export default connect({
  field: field(state`${props`path`}`),
  fieldChanged: signal`forms.fieldChanged`
  },
  function RadioList({name, options, path, field, fieldChanged, legendtext}) {
    function onClicked(e) {
      fieldChanged({
        field: path,
        value: e.target.value
      })
    }

  function displayList(name, options) {
    let list = []
    options.forEach(function(o) {
      list.push(h('li', [
        h('input', {id: name + '-' + o.name, type: 'radio', name: name, value: o.name, onClick: (e) => onClicked(e)}),
        h('label', {for: name + '-' + o.name}, [
          h('span.fa-stack', [
            h('i.fa.fa-circle-thin.fa-stack-1x'),
            h('i.fa.fa-dot-circle-o.fa-stack-1x')
          ]),
          o.label
        ])
      ]))
    });
    return list
  }

  return h('div.field', [
        h('div.control', [
          h('legend.label', [legendtext]),
          h('ul', displayList(name, options))
        ])
      ])
})

