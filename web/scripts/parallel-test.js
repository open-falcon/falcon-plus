require('shelljs/global')
const { CIRCLE_NODE_INDEX } = process.env

const cmds = [
  'npm test && curl -s https://codecov.io/bash | bash',
  'npm run lint',
  'npm run build'
]

if (cmds[CIRCLE_NODE_INDEX]) {
  console.log(cmds[CIRCLE_NODE_INDEX])

  const run = exec(cmds[CIRCLE_NODE_INDEX])

  if (run.code) {
    exit(1)
  }
}
