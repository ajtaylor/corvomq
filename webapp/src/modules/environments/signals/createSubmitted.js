import {set} from 'cerebral/operators'
import {state} from 'cerebral/tags'
import create from '../actions/create'

export default [
  set(state`environments.creating`, true),
  set(state`environments.create.result`, 'waiting'),
  create, {
    success: [
      set(state`environments.creating`, false),
      set(state`environments.create.result`, 'success')
    ],
    error: [
      set(state`environments.creating`, false),
      set(state`environments.create.result`, 'failed')
    ]
  }
]
