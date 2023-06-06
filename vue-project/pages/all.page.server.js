
import { renderToNodeStream } from '@vue/server-renderer'
import { escapeInject } from 'vite-plugin-ssr/server'
import { createApp } from './app'

export const passToClient = ['pageProps', 'documentProps', 'i18nLocale']

export async function render(pageContext){
	const { app, router } = await createApp(pageContext)

	router.push(pageContext.urlPathname)
	await router.isReady()
	const is404 = router.currentRoute.value.meta && router.currentRoute.value.meta.is404

	const stream = await renderToNodeStream(app)

	var title = router.currentRoute.value.meta && router.currentRoute.value.meta.title
	if(typeof title === 'function'){
		title = title(router.currentRoute.value)
	}
	if(title){
		title += ' - '
	}else{
		title = ''
	}

	const documentHtml = escapeInject`<!DOCTYPE html>
		<html>
			<head>
				<title>${title}MCDReforged Plugin List</title>
				<meta charset="UTF-8">
				<link rel="icon" href="/assets/favicon.ico">
				<meta name="keywords" content="mcdr,mcdreforged,minecraft,plugin,api">
				<meta name="description" content="An MCDReforged plugin list access point">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<link rel="stylesheet" type="text/css" href="/assets/main.css">
			</head>
			<body>
				<div id="app">${stream}</div>
			</body>
		</html>
	`

	return {
		documentHtml,
		pageContext: {
			enableEagerStreaming: true,
			isVue404: is404, // cannot use reversed key `is404`
		},
	}
}
