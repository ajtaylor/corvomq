import {set} from 'cerebral/operators'
import {state} from 'cerebral/tags'
import getList from '../actions/getList'

export default [
  set(state`app.currentPage`, 'environments'),
  getList, {
    success: [],
    error: []
  }
]
