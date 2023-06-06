
import { createApp } from './app'

export const clientRouting = true

var app, router

export async function render(pageContext){
	const { Page, pageProps } = pageContext
	if(!app){
		({ app, router } = await createApp(pageContext))
		await router.isReady()
		app.mount('#app')
	}else{
		router.push(pageContext.urlPathname)
		await router.isReady()
		console.log('page changed:', router.currentRoute)
	}
}
