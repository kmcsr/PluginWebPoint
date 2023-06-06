
import { createSSRApp, defineComponent, h, markRaw, reactive } from 'vue'
import VueGtag from 'vue-gtag'
import { setupI18n } from '@/i18n'
import { newRouter } from '@/router'

const production = process.env.NODE_ENV === 'production';

export async function createApp(pageContext){
	const { Page } = pageContext

	const app = createSSRApp(Page)

	// const pageContextReactive = reactive(pageContext)
	// setPageContext(app, pageContextReactive)

	if(!import.meta.env.SSR){
		app.use(await import('vue-cookies'))
	}
	app.use(await setupI18n(pageContext.i18nLocale || 'en_us'))

	const router = newRouter()
	app.use(router)

	if(production){
		app.use(VueGtag, {
			config: {
				id: 'G-B34TLWC63Q',
			},
		}, router)
	}
	return { app, router }
}
