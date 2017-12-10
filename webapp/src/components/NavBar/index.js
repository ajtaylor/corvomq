import {connect} from 'cerebral/inferno'
// import {state} from 'cerebral/tags'
import h from 'inferno-hyperscript'

export default connect({},
  function NavBar() {
    return (
      h('nav.level', [
        h('div.level-left', [
          h('p.is-size-3', ['CorvoMQ'])
        ])
      ])
    )
  }
)
