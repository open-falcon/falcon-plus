import './sneak'
import App from './app'
import router from './router'
import store from './store'

new window.Vue({
  router,
  store,
  ...App
}).$mount('#app')
