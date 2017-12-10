import {connect} from 'cerebral/inferno'
import {state, props, signal} from 'cerebral/tags'
import {field} from '@cerebral/forms'
import h from 'inferno-hyperscript'

function displayHelptext(helptext) {
  return h('p.help', [helptext])
}

function displayLabeltext(name, labeltext, labelsubtext) {
  return h('label.label', {for: name}, [
           labeltext,
           labelsubtext != undefined ? h('small', [' ' + labelsubtext]) : null
         ])
}

export default connect({
  field: field(state`${props`path`}`),
  fieldChanged: signal`forms.fieldChanged`
},
  function Textbox({name, useFieldset, path, placeholder, password, field, fieldChanged, helptext, labeltext, labelsubtext, maxWidth}) {
    function onChanged(e) {
      fieldChanged({
        field: path,
        value: e.target.value
      })
    }

  function displayTextbox(labeltext, name, labelsubtext, password, helptext, field, placeholder, maxWidth) {
    var mw = (maxWidth != undefined ? '.mw' + maxWidth : '')
    return h('div.field', [
            labeltext != undefined ? displayLabeltext(name, labeltext, labelsubtext) : null,
            h('div.control', [
            h('input.input' + mw,
              {name: name, type: password ? 'password' : 'text', onChange: (e) => onChanged(e), value: field.value, placeholder: placeholder})
            ]),
            helptext != undefined ? displayHelptext(helptext) : null
          ])
  }

  if (useFieldset) {
    return h('fieldset.b--transparent', [
            displayTextbox(labeltext, name, labelsubtext, password, helptext, field, placeholder, maxWidth)
          ])
  }
  else {
    return displayTextbox(labeltext, name, labelsubtext, password, helptext, field, placeholder, maxWidth)
  }
})
