function createUser({firebase, path, state}) {
  const email = state.get('user.signupForm.emailAddress.value')
  const password = state.get('user.signupForm.password.value')

  return firebase.createUserWithEmailAndPassword(email, password)
    .then((result) => {
      if(result.error) {
        return path.error({error: result.error})
      }
      return path.success({newUser: result.user})
    })
    .catch(path.error)
}

export default createUser
