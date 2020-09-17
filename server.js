const serve = require('serve-handler')
const http = require('http')
const URL = require('url')
const querystring = require('querystring')
const { handler: tileLambda } = require('./src/lambda/tile.js')

const server = http.createServer((request, response) => {
  const url = URL.parse(request.url)

  if (url.pathname === '/.netlify/functions/tile') {
    return tileLambda({
      queryStringParameters: querystring.parse(url.query)
    }).then((result) => {
      Object.entries(result.headers).forEach(({ 0: header, 1: content}) => {
        response.setHeader(header, content)
      })

      response.writeHead(result.statusCode)

      if (result.isBase64Encoded) {
        result.body = Buffer.from(result.body, 'base64')
      }

      response.end(result.body)
    })
  }

  return serve(request, response, { public: 'public' })
})

server.listen(3000, () => {
  console.log('Running at http://localhost:3000');
});
