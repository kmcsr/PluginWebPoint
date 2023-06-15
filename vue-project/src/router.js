
import { createRouter, createMemoryHistory, createWebHistory } from 'vue-router'
import HomeView from './views/HomeView.vue'
import PluginIndexView from './views/PluginIndexView.vue'
import PageNotFound from './views/PageNotFound.vue'

const PRODUCTION = process.env.NODE_ENV === 'production'

export function newRouter(){
	const router = createRouter({
		history: import.meta.env.SSR ?createMemoryHistory() :createWebHistory(import.meta.env.BASE_URL),
		routes: [
			{
				path: '/',
				name: 'index',
				component: HomeView,
				meta: {
					title: 'Home',
				}
			},
			{
				path: '/about',
				name: 'about',
				component: () => import('./views/AboutView.vue'),
				meta: {
					title: 'About',
				}
			},
			{
				path: '/plugins',
				name: 'plugin_index',
				component: PluginIndexView,
				meta: {
					title: 'Plugin List',
				}
			},
			{
				path: '/plugin/:plugin/',
				props: true,
				name: 'plugin',
				component: () => import('./views/PluginView.vue'),

			},
			{
				path: '/download/:plugin/:tag/:filename',
				props: true,
				name: 'download',
				component: () => import('./views/DownloadView.vue'),
			},
			{
				path: '/author/:author',
				props: true,
				name: 'author',
				component: () => import('./views/AuthorView.vue'),
			},
			{
				path: "/:pathMatch(.*)*",
				component: PageNotFound,
				name: 'not-found',
				meta: {
					is404: true,
					title: '404 Not Found',
				}
			}
		],
		scrollBehavior(to, from, savedPosition){
			console.log('calling scrollBehavior:', savedPosition)
			if(savedPosition){
				return savedPosition
			}
			return
		}
	})
	router.beforeEach((to, from) => {
		if(typeof document !== 'undefined'){
			var title = to.meta.title
			if(title || !from || from.path !== to.path){
				if(typeof title === 'function'){
					title = title(to)
				}
				if(title){
					title += ' - '
				}else{
					title = ''
				}
				document.title = title + 'MCDReforged Plugin List'
			}
		}
	})
	if(!PRODUCTION && false){
		router.beforeEach((to, from) => {
			console.log('Before each hook:', to, from)
		})
		router.afterEach((to, from) => {
			console.log('After each hook:', to, from)
		})
	}
	return router
}
