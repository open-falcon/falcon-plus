import { isBrowser } from './is-env'

const { Cookies } = window

module.exports = (opts = {}) => {
  const hasSpin = isBrowser && opts.spin
  if (hasSpin) {
    if (typeof opts.spin === 'string') {
      opts.spin = Array.from(document.querySelectorAll(opts.spin))
    } else if (!Array.isArray(opts.spin)) {
      opts.spin = [opts.spin]
    }

    opts.spin.forEach((el) => {
      el.style.visibility = 'visible'
    })
  }

  return new Promise((resolve, reject) => {
    window.axios(_.merge({
      headers: {
        ApiToken: JSON.stringify(
          {
            name: Cookies.get('name'),
            sig: Cookies.get('sig'),
          }
        ),
      },
      baseURL: __conf.dev.apiBase,
    }, opts))
    .then((res)=> {
      if (hasSpin) {
        opts.spin.forEach((el) => {
          el.style.visibility = 'hidden'
        })
      }

      resolve(res)
    })
    .catch((err) => {
      if (hasSpin) {
        opts.spin.forEach((el) => {
          el.style.visibility = 'hidden'
        })
      }

      reject(err)
    })
  })
}
