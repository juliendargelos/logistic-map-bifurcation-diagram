import nodeResolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'

export default {
  input: 'src/lambda/tile.js',
  output: { dir: 'lambda', format: 'cjs' },
  plugins: [
    nodeResolve(),
    commonjs()
  ]
}
