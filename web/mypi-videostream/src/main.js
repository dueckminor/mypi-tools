import 'material-design-icons-iconfont/dist/material-design-icons.css'

import Vue from 'vue';
import vuetify from './plugins/vuetify';
import VueRouter from 'vue-router'

import App from './App.vue';
import 'vuetify/dist/vuetify.min.css'

const router = new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    { path: '/', component: App },
  ]
})

new Vue({
  router,
  vuetify,
  el: '#app',
  render: h => h(App)
});
