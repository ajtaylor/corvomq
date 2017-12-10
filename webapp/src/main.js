import {render} from 'inferno'
import h from 'inferno-hyperscript'
import {Container} from 'cerebral/inferno'
import controller from './controller'
import App from './components/App'

const appContainer = document.getElementById('app')

render(h(Container, {controller: controller}, h(App)), appContainer)
