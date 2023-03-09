import { createRouter, createWebHistory } from 'vue-router'
import IndexView from '../views/IndexView.vue'
import PluginView from '../views/PluginView.vue'
import DownloadView from '../views/DownloadView.vue'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'Plugin List',
			component: IndexView,
			meta: {
				title: 'Plugin List | PWP'
			}
		},
		{
			path: '/about',
			name: 'About',
			// route level code-splitting
			// this generates a separate chunk (About.[hash].js) for this route
			// which is lazy-loaded when the route is visited.
			component: () => import('../views/AboutView.vue'),
			meta: {
				title: 'About | PWP'
			}
		},
		{
			path: '/plugin/:plugin',
			props: true,
			name: 'Plugin',
			component: PluginView,
			meta: {
				title: (to) => `Plugin ${to.params.plugin}`,
			}
		},
		{
			path: '/download/:plugin/:tag/:filename',
			props: true,
			name: 'Download',
			component: DownloadView,
		}
	]
})

router.beforeEach((to, from) => {
	if(to.meta.title){
		if(typeof to.meta.title === 'function'){
			document.title = to.meta.title(to, from)
		}else{
			document.title = to.meta.title
		}
	}else{
		document.title = to.name
	}
})

export default router
