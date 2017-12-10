function login({http, path, state}) {
  const email = state.get('user.loginForm.emailAddress.value')
  const password = state.get('user.loginForm.password.value')

  return http.post('/user/login', {
      emailAddress: email,
      password: password}
      )
    .then((response) => {
      if(response.error) {
        return path.error({error: response.error})
      }
      return path.success({token: response.result.token})
    })
    .catch(path.error)
}

export default login
