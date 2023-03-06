import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import './assets/main.css'

const app = createApp(App)

app.config.globalProperties.$apiPrefix = '/dev'

app.use(router)

app.mount('#app')
