import {resolve} from 'path'

import {generateApi} from 'swagger-typescript-api'

const inputFile = resolve(
  process.cwd(),
  '../backend-api/api/service/openapi.yaml',
)

generateApi({
  name: 'Api.ts',
  output: resolve(process.cwd(), './src/core/api'),
  input: inputFile,
  httpClientType: 'axios',
})
