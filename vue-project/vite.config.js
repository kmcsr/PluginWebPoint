import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig(async ({ command, mode }) => {
	console.log(command, mode);
	const isdev = mode === 'development';
	const minify = isdev ?'' :'esbuild';

	return {
		plugins: isdev ?[
			await import('@vitejs/plugin-vue'),
			await import('@vitejs/plugin-vue-jsx'),
			await import('vite-plugin-ssr/plugin')
		] :[],
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
