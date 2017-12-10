import {when} from 'cerebral/operators'
import FetchedData from '../../utils/fetchedData'

function getList({state, http, path}) {
  switch (state.get('environments.dataState')) {
    case FetchedData.NOT_REQUESTED:
      state.set('environments.dataState', FetchedData.LOADING)
      return http.get('/environments')
        .then((response) => {
          if(response.error) {
            return path.error({error: response.error})
          }
          state.set('environments.dataState', FetchedData.LOADED)
          state.set('environments.list', response.result.environments)
          return path.success()
        })
        .catch(path.error)

      case FetchedData.FAILED:
        return path.error()

      default:
        return path.success()
  }
}

export default getList
