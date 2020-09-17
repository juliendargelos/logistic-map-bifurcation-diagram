// const { PNG } = require('pngjs')

// const iterations = 20
// const start = 0.25
// const width = 256
// const height = 256
// const minimumColor = [255, 100, 255]
// const maximumColor = [0, 0, 255]
// const backgroundColor = [0, 0, 0]

// exports.handler = async ({ queryStringParameters }) => {
//   let { x = '0', y = '0', z = '0' } = queryStringParameters
//   const range = 1 / Math.pow(2, parseInt(z, 10))
//   const scaledIterations = Math.min(50000, Math.max(500, iterations / Math.log(range + 1)))

//   x = parseInt(x, 10) * range
//   y = parseInt(y, 10) * range

//   const image = new PNG({ width, height })
//   const data = image.data
//   const histogram = new Array(width * height).fill(0)

//   if (x >= 0 && x <= 1 && y >= 0 && y <= 1) {
//     let i, j, k, v, f, h

//     for (i = 0; i < width; i++) {
//       const rate = (i / (width - 1) * range + x) * 3 + 1
//       const values = new Array(height).fill(0)

//       for (j = 0, v = start; j < 1000; j++) {
//         v = v * rate * (1 - v)
//       }

//       for (j = 0, f = 0; j < scaledIterations; j++) {
//         v = v * rate * (1 - v)
//         k = 1 - v

//         if (k >= y && k <= y + range) {
//           k = Math.round((k - y) / range * (height - 1))
//           h = values[k]
//           values[k] = h + 1
//           if (!h) f++
//         }
//       }

//       values.forEach((value, l) => {
//         histogram[i + l * width] = value * f
//       })
//     }
//   }

//   const maximum = histogram.reduce((maximum, value) => (
//     Math.max(value, maximum)
//   ), 0)

//   const smoothing = 0

//   histogram.forEach((value, i) => {
//     i *= 4

//     data[i + 3] = 255

//     if (value) {
//       // value = Math.log(value / height * 0.9 + 0.1) / 2 + 0.5
//       // value = Math.log(value + 1)
//       //

//       value /= height

//       for (var s = 0; s < smoothing; s++) {
//         value = Math.min(1, -Math.log((1 - value) * 0.9 + 0.1))
//       }

//       data[i    ] = Math.round(value * maximumColor[0] + (1 - value) * minimumColor[0])
//       data[i + 1] = Math.round(value * maximumColor[1] + (1 - value) * minimumColor[1])
//       data[i + 2] = Math.round(value * maximumColor[2] + (1 - value) * minimumColor[2])
//     } else {
//       data.set(backgroundColor, i)
//     }
//   })

//   return new Promise((resolve) => {
//     const chunks = []

//     image.on('data', chunk => chunks.push(chunk))
//     image.on('end', () => resolve({
//       statusCode: 200,
//       isBase64Encoded: true,
//       headers: {
//         'content-type': 'image/png',
//         'cache-control': 'max-age=31536000'
//       },
//       body: Buffer.concat(chunks).toString('base64')
//     }))

//     image.pack()
//   })
// }

package main

import (
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "strconv"
  "math"
  "image"
  "image/png"
  "image/color"
  "bytes"
  b64 "encoding/base64"
)

const iterations float64 = 100
const start float64 = 0.25
const width int = 1024
const height int = 1024

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
  minimumColor := [...]float64{255, 100, 255}
  maximumColor := [...]float64{0, 0, 255}
  backgroundColor := [...]uint8{0, 0, 0}

  xP := request.QueryStringParameters["x"]
  yP := request.QueryStringParameters["y"]
  zP := request.QueryStringParameters["z"]

  a, _ := strconv.ParseInt(zP, 10, 32)
  amplitude := 1 / math.Pow(2, float64(a))

  scaledIterations := int(math.Min(100000, math.Max(500, iterations / math.Log(amplitude + 1))))

  xI, _ := strconv.ParseInt(xP, 10, 32)
  yI, _ := strconv.ParseInt(yP, 10, 32)

  x := float64(xI) * amplitude
  y := float64(yI) * amplitude

  img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
  histogram := [width * height]float64{}

  if x >= 0 && x <= 1 && y >= 0 && y <= 1 {
    var v float64
    var k float64
    var ki int
    var f int64
    var h int64

    for i := 0; i < width; i++ {
      var values [height]int64
      rate := (float64(i) / (float64(width) - 1) * amplitude + x) * 3 + 1
      v = start
      f = 0

      for j := 0; j < 1000; j++ {
        v = v * rate * (1 - v)
      }

      for j := 0; j < scaledIterations; j++ {
        v = v * rate * (1 - v)
        k = 1 - v

        if (k >= y && k <= y + amplitude) {
          ki = int(math.Round((k - y) / amplitude * (float64(height) - 1)))
          h = values[ki]
          values[ki] = h + 1

          if h == 0 {
            f++
          }
        }
      }

      for l, value := range values {
        histogram[i + l * width] = float64(value * f)
      }
    }

    for i, value := range histogram {
      py := int(math.Floor(float64(i) / float64(width)))
      px := i - py * width

      if value != 0 {
        value = value / float64(height)

        img.Set(px, py, color.RGBA{
          uint8(math.Round(value * maximumColor[0] + (1 - value) * minimumColor[0])),
          uint8(math.Round(value * maximumColor[1] + (1 - value) * minimumColor[1])),
          uint8(math.Round(value * maximumColor[2] + (1 - value) * minimumColor[2])),
          0xff,
        })
      } else {
        img.Set(px, py, color.RGBA{
          backgroundColor[0],
          backgroundColor[1],
          backgroundColor[2],
          0xff,
        })
      }
    }
  }

  buffer := new(bytes.Buffer)
  png.Encode(buffer, img)

  return &events.APIGatewayProxyResponse{
    IsBase64Encoded: true,
    StatusCode: 200,
    Headers: map[string]string{
      "content-type": "image.png",
    },
    Body: b64.StdEncoding.EncodeToString(buffer.Bytes()),
  }, nil
}

func main() {
  lambda.Start(handler)
}
