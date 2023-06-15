
import { renderToNodeStream } from '@vue/server-renderer'
import { escapeInject } from 'vite-plugin-ssr/server'
import { createApp } from './app'

export const passToClient = ['pageProps', 'documentProps', 'i18nLocale']

export async function render(pageContext){
	const { app, router } = await createApp(pageContext)

	await router.push(pageContext.urlOriginal)
	await router.isReady()
	const route = router.currentRoute.value
	const meta = route.meta
	const is404 = meta.is404

	const stream = await renderToNodeStream(app)

	var keywords = ['mcdr', 'mcdreforged', 'minecraft', 'plugin', 'api']
	var description = meta.description
	if(typeof description === 'function'){
		description = description(route)
	}
	if(!description){
		description = 'An MCDReforged plugin list access point'
	}
	var title = meta.title
	if(typeof title === 'function'){
		title = title(route)
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
				<meta name="keywords" content="${keywords.join(',')}">
				<meta name="description" content="${description}">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<link rel="stylesheet" type="text/css" href="/assets/main.css">
			</head>
			<body>
				<div id="app">${stream}</div>
				<script defer async id="mermaid-script" type="text/javascript" src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
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
