import 'material-design-icons-iconfont/dist/material-design-icons.css'

import { createApp } from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import router from './router'
import VueSocketIO from 'vue-3-socket.io';

loadFonts()

let vue = createApp(App);
vue.use(router);
vue.use(vuetify);
vue.use(new VueSocketIO({
    connection: `//${window.location.host}`,
    options: { path: "/ws" }
  }));
vue.mount('#app');
