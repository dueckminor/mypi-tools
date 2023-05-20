import 'material-design-icons-iconfont/dist/material-design-icons.css'
import 'hack-font/build/web/hack.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia';
import vuetify from './plugins/vuetify'
import { loadFonts } from './plugins/webfontloader'
import router from './router'
import App from './App.vue'

const app = createApp(App)
app.use(createPinia())

loadFonts()

app.use(router)
app.use(vuetify)

app.mount('#app')

