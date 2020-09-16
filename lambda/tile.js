exports.handler = async (event, context) => {
  const {x = '', y = '', z = ''} = event.queryStringParameters

  return {
    statusCode: 200,
    body: `${x},${y},${z}`
  }
}
