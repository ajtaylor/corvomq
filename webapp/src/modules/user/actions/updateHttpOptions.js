function updateHttpOptions({http, props}) {
  http.updateOptions({
    headers: {
      'Authorization': `Bearer ${props.token}`
    }
  })
}

export default updateHttpOptions
