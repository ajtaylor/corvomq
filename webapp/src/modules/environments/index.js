import {state} from 'cerebral/tags'
import routed from './signals/routed'
import createRouted from './signals/createRouted'
import createSubmitted from './signals/createSubmitted'
import FetchedData from '../utils/fetchedData'

export default {
  signals: {
    routed: routed,
    createRouted: createRouted,
    createSubmitted: createSubmitted
  },
  state: {
    dataState: FetchedData.NOT_REQUESTED,
    list: [],
    createForm: {
      name: {
        value: '',
        isRequired: true
      },
      server: {
        value: '',
        isRequired: true
      },
      infrastructure: {
        value: '',
        isRequired: true
      },
      tls_enabled: {
        value: false,
        isRequired: true
      }
    },
    create: {
      result: ''
    }
  }
}
