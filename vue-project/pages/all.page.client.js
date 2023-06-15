
import { nextTick } from 'vue'
import { createApp } from './app'

export const clientRouting = true

var app, router

export async function render(pageContext){
	if(!app){
		({ app, router } = await createApp(pageContext))
		await router.isReady()
		app.mount('#app')
	}else{
		if(pageContext.isBackwardNavigation !== null){ // TODO: as far as we know, isBackwardNavigation could be null or false
			console.debug('pageContext updating:', pageContext, pageContext.isBackwardNavigation)
			await router.push(pageContext.urlOriginal)
		}
	}
}
