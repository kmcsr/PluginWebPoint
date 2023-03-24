import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'index',
			redirect: '/plugins',
		},
		{
			path: '/about',
			name: 'about',
			// route level code-splitting
			// this generates a separate chunk (About.[hash].js) for this route
			// which is lazy-loaded when the route is visited.
			component: () => import('../views/AboutView.vue'),
		},
		{
			path: '/plugins',
			name: 'plugin_index',
			component: () => import('../views/PluginIndexView.vue'),
		},
		{
			path: '/plugin/:plugin',
			props: true,
			name: 'plugin',
			component: () => import('../views/PluginView.vue'),
		},
		{
			path: '/download/:plugin/:tag/:filename',
			props: true,
			name: 'download',
			component: () => import('../views/DownloadView.vue'),
		},
		{
			path: '/author/:author',
			props: true,
			name: 'author',
			component: () => import('../views/AuthorView.vue'),
		}
	]
})

export default router
