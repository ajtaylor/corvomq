function create({http, path, state}) {
  const name = state.get('environments.createForm.name.value')
  const server = state.get('environments.createForm.server.value')
  const infrastructure = state.get('environments.createForm.infrastructure.value')
  const tls_enabled = state.get('environments.createForm.tls_enabled.value')

  return http.post('/environments/create', {
      name: name,
      server: server,
      infrastructure: infrastructure,
      tls_enabled: tls_enabled})
    .then((response) => {
      console.log(response)
      if(response.error) {
        return path.error({error: response.error})
      }
      return path.success({result: response.result})
    })
    .catch(path.error)
}

export default create
