
import path from 'path'
import { fileURLToPath } from 'url'
import express from 'express'
import cookieParser from 'cookie-parser'
import compression from 'compression'
import sirv from 'sirv'
import { renderPage } from 'vite-plugin-ssr/server'
import { LANG_COOKIE, i18nLangMap } from '../src/i18n/consts.js'

const isProduction = process.env.NODE_ENV === 'production'
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const root = path.dirname(__dirname)

process.on('uncaughtException', function (err){
	console.error('Uncaught exception:', err)
});

(async function(){
	const app = express()

	app.use(compression())
	app.use(cookieParser())

	if(isProduction){
		app.use(sirv(`${root}/dist/client`))
	}else{
		const vite = await import('vite')
		const viteDevMiddleware = (
			await vite.createServer({
				root,
				server: {
					middlewareMode: true,
				},
				appType: 'ssr',
			})
		).middlewares
		app.use(viteDevMiddleware)
	}

	app.get('*', async (req, res, next) => {
		const locale = req.cookies[LANG_COOKIE] || req.acceptsLanguages(...Object.keys(i18nLangMap))
		const pageContextInit = {
			urlOriginal: req.originalUrl,
			i18nLocale: locale,
		}
		const pageContext = await renderPage(pageContextInit)
		const { httpResponse } = pageContext
		if(!httpResponse){
			return next()
		}
		if(pageContext.isVue404){
			httpResponse.statusCode = 404
		}

		res.cookie(LANG_COOKIE, locale, { maxAge: 3600 * 24 * 365, httpOnly: false })
		res.
			status(httpResponse.statusCode).
			type(httpResponse.contentType)
		httpResponse.pipe(res)
	})

	const port = Number.parseInt(process.argv[2]) || (isProduction ?80 :3080)
	app.listen(port)
	console.log('Env:', process.env.NODE_ENV)
	console.log(`Server running at http://localhost:${port}`)
})()
