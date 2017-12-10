import {connect} from 'cerebral/inferno'
import {state, props, signal} from 'cerebral/tags'
import {field} from '@cerebral/forms'
import h from 'inferno-hyperscript'

export default connect({},
  function Button({type, text, isFullWidth}) {
    const classes = '.button.is-primary' + (isFullWidth != undefined && isFullWidth ? '.is-fullwidth' : '')
    return (
      h('div.field', [
        h('div.control' + (isFullWidth != undefined && isFullWidth ? '.is-expanded' : ''), [
          h('button' + classes, {type: type}, [text])
        ])
      ])
    )
  })
