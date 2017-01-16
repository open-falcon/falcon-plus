const { repository, version } = require('../package.json')

const { GH_TOKEN } = process.env
const tokenRepo = repository.replace(/(github.com)/, `${GH_TOKEN}@$1`)

module.exports = {
  ver: version,
  tag: `v${version}`,
  repo: tokenRepo,
}
