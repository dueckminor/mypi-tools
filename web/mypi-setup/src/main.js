import 'material-design-icons-iconfont/dist/material-design-icons.css'
import '@mdi/font/css/materialdesignicons.css'

import Vue from 'vue';
import vuetify from './plugins/vuetify';
import VueRouter from 'vue-router'

import App from './App.vue';
import 'vuetify/dist/vuetify.min.css'
import router from "./router";

Vue.use(VueRouter)

new Vue({
  router,
  vuetify,
  el: '#app',
  render: h => h(App)
});
