
import { createRouter, createMemoryHistory, createWebHistory } from 'vue-router'
import PluginIndexView from './views/PluginIndexView.vue'
import PageNotFound from './views/PageNotFound.vue'

export function newRouter(){
	return createRouter({
		history: import.meta.env.SSR ?createMemoryHistory() :createWebHistory(import.meta.env.BASE_URL),
		routes: [
			{
				path: '/',
				name: 'index',
				redirect: '/plugins',
				meta: {
					title: 'Index',
				}
			},
			{
				path: '/about',
				name: 'about',
				// route level code-splitting
				// this generates a separate chunk (About.[hash].js) for this route
				// which is lazy-loaded when the route is visited.
				component: () => import('./views/AboutView.vue'),
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
				path: '/plugin/:plugin',
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
				meta: {
					is404: true,
					title: '404 Not Found',
				}
			}
		]
	})
}
