import {connect} from 'cerebral/inferno'
import {state} from 'cerebral/tags'
import {string} from 'cerebral/operators'
import h from 'inferno-hyperscript'
import FetchedData from '../../modules/utils/fetchedData'
import Css from '../Css'

export default connect({
  dataState: state`environments.dataState`,
  environmentList: state`environments.list`
},
  function Environments({dataState, environmentList}) {
    return (
      h('div.container', [
        h('div.level', [
          h('div.level-left', [
            h('div.level-item', [
              h('h4.is-size-4', ['Environments'])
            ])
          ]),
          h('div.level-right', [
            h('div.level-item', [
              h('a.button.is-success.is-pulled-right', {href: '/environments/create'}, ['Create new environment']),
            ])
          ])
        ]),
        display(dataState, environmentList)
      ])
    )
})

function display(dataState, environmentList) {
  console.log(environmentList)
  switch (dataState) {
    case FetchedData.LOADING:
      return h('p', ['Loading...'])
    case FetchedData.LOADED:
      // return h('table.w-100', {cellspacing: 0}, [
      return h('table.table', [
        h('thead', [
          h('th', ['Name']),
          h('th', ['Server']),
          h('th', ['Infrastructure']),
          h('th', ['URL']),
          h('th', [''])
        ]),
        h('tbody', [
          environmentList.map(
            (item) => {
              return h('tr', [
                h('td', [item.name]),
                h('td', [item.server]),
                h('td', [item.infrastructure]),
                h('td', [item.url]),
                h('td.has-text-success', ['Running'])
              ])
              })
        ])
      ])
    case FetchedData.FAILED:
      return h('p', ['There was a problem loading your list of environments.'])
    default:
      return h('p', ['Default'])
  }
}
