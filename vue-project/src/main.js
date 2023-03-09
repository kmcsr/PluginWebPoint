import { createApp } from 'vue'
import VueGtag from 'vue-gtag'
import App from './App.vue'
import { setupI18n } from './i18n'
import router from './router'

import './assets/main.css'

const production = process.env.NODE_ENV === 'production';

console.debug('process-env-NODE-ENV:', process.env.NODE_ENV)
console.debug('import-meta-env-DEV:', import.meta.env.DEV)

if(!production){
	console.debug('mode:', dev)
}

(async function(app){

	app.config.globalProperties.$apiPrefix = '/dev'

	const i18n = setupI18n({
		messages: {
			'en': await import('./i18n/en_us.json'),
		}
	})

	app.use(i18n)

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
