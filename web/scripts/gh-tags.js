require('shelljs/global')
const { tag, repo } = require('./deploy')

exec(`git tag ${tag}`)
exec(`git push ${repo} ${tag}`)
