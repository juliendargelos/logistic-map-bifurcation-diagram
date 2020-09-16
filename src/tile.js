const { PNG } = require('pngjs')

const iterations = 100//50000
const start = 0.25
const width = 256
const height = 256
const minimumColor = [255, 0, 255]
const maximumColor = [0, 0, 255]
const backgroundColor = [0, 0, 0]

exports.handler = async (event, context) => {
  let { x = '0', y = '0', z = '1' } = event.queryStringParameters
  z = Math.pow(parseInt(z, 10), 2)
  const range = 1 / z

  x = parseInt(x, 10) * range
  y = parseInt(y, 10) * range

  const image = new PNG({ width, height })
  const data = image.data
  const histogram = new Array(width * height)

  let i, j, k, v, f, h

  for (i = 0; i < width; i++) {
    const rate = (i / (width - 1) * range + x) * 6 + 2
    const values = new Array(height).fill(0)

    for (j = 0, v = start; j < 1000; j++) {
      v = v * rate * (1 - v)
    }

    for (j = 0, f = 0; j < iterations; j++) {
      v = v * rate * (1 - v)
      k = 1 - v

      if (k >= y && k <= y + range) {
        k = Math.round((k - y) / range * (height - 1))
        h = values[k]
        values[k] = h + 1
        if (!h) f++
      }
    }

    values.forEach((value, l) => {
      histogram[i + l * width] = Math.min(80, value) * Math.min(f, 80)
    })
  }

  const maximum = histogram.reduce((maximum, value) => (
    Math.max(value, maximum)
  ), 0)

  const smoothing = 0

  histogram.forEach((value, i) => {
    i *= 4

    data[i + 3] = 255

    if (value) {
      value /= maximum

      for (var s = 0; s < smoothing; s++) {
        value = Math.min(1, -Math.log((1 - value) * 0.9 + 0.1))
      }

      data[i    ] = Math.round(value * maximumColor[0] + (1 - value) * minimumColor[0])
      data[i + 1] = Math.round(value * maximumColor[1] + (1 - value) * minimumColor[1])
      data[i + 2] = Math.round(value * maximumColor[2] + (1 - value) * minimumColor[2])
    } else {
      data.set(backgroundColor, i)
    }
  })

  return new Promise((resolve) => {
    const chunks = []

    image.on('data', chunk => chunks.push(chunk))
    image.on('end', () => resolve({
      statusCode: 200,
      isBase64Encoded: true,
      headers: { 'content-type': 'image/png' },
      body: Buffer.concat(chunks).toString('base64')
    }))

    image.pack()
  })
}
