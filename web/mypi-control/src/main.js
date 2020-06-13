import Vue from 'vue'
import VueSocketIO from 'vue-socket.io';
//import VueSocketIOExt from 'vue-socket.io-extended';
//import * as io from 'socket.io-client'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import Vuetify from './plugins/vuetify';
import VueGraph from 'vue-graph'
import router from './router'
import store from './store'

Vue.config.productionTip = false
Vue.use(Vuetify)
Vue.use(VueGraph)

Vue.use(new VueSocketIO({
  debug: true,
  connection: `//${window.location.host}`,
  vuex: {
      store,
      actionPrefix: 'SOCKET_',
      mutationPrefix: 'SOCKET_'
  },
  options: { path: "/ws" } //Optional options
}))


new Vue({
  router,
  vuetify,
  store,
  render: h => h(App)
}).$mount('#app')
