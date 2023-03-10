import { createApp } from 'vue'
import VueGtag from 'vue-gtag'
import App from './App.vue'
import { setupI18n } from './i18n'
import router from './router'

import './assets/main.css'

const production = process.env.NODE_ENV === 'production';

if(!production){
	console.debug('mode:', process.env.NODE_ENV)
}

(async function(app){

	app.config.globalProperties.$apiPrefix = '/dev'

	app.use(await setupI18n())

	app.use(router)

	if(production){
		app.use(VueGtag, {
			config: {
				id: 'G-B34TLWC63Q',
			},
		}, router)
	}

	app.mount('#app')
})(createApp(App))
