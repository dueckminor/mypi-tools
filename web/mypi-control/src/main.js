import Vue from 'vue'
//import VueSocketio from 'vue-socket.io';
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router'

Vue.config.productionTip = false
//Vue.use(VueSocketio, `//${window.location.host}`)

new Vue({
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app')
