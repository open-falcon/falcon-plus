import qs from './utils/qs'

const { Cookies } = window
const q = qs()
const isDev = process.env.NODE_ENV === 'development' || !process.env.NODE_ENV
const isSneak = '--sneak' in q

if (isDev && isSneak) {
  Cookies.set('name', 'root')
  Cookies.set('sig', '427d6803b78311e68afd0242ac130006')
}
