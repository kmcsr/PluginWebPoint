import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import ssr from 'vite-plugin-ssr/plugin'

// https://vitejs.dev/config/
export default defineConfig(async ({ command, mode }) => {
	console.log(command, mode);
	const isdev = mode === 'development';
	const minify = isdev ?'' :'esbuild';

	return {
		plugins: [
			vue(),
			vueJsx(),
			ssr()
		],
		base: '/',
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
		mode: mode,
		build: {
			minify: minify,
		},
		esbuild: {
			pure: mode === 'production' ? ['console.debug'] : [],
		}
	}
})
