import 'material-design-icons-iconfont/dist/material-design-icons.css'

import Vue from 'vue';
import Vuetify from 'vuetify'
import VueRouter from 'vue-router'

import App from './App.vue';
import 'vuetify/dist/vuetify.min.css'

Vue.use(Vuetify)
Vue.use(VueRouter)

const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: App },
  ]
})

new Vue({
  router,
  el: '#app',
  render: h => h(App)
});
