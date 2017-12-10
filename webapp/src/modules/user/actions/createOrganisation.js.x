function create(firebase, organisationName, userId) {
  return firebase.push('organisations', {
    name: organisationName,
    users: [userId]
  })  
}

function createOrganisation({firebase, path, state}) {
  const uid = state.get('user.uid')
  const organisationName = state.get('user.signupForm.organisationName.value')

  return create(firebase, organisationName, uid)
    .then((result) => {
      console.log(result)
      return path.success({newOrganisationKey: result.key})})
}

export default createOrganisation
