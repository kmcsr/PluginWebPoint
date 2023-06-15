<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import WebBox from 'vue-material-design-icons/WebBox.vue'
import IconSun from '@/components/icons/IconSun.vue'
import IconMoon from '@/components/icons/IconMoon.vue'
import { i18nLangMap } from '@/i18n'

var observer
var darkTheme

function switchTheme(){
	if(darkTheme.value = !darkTheme.value){
		document.documentElement.classList.add('dark')
	}else{
		document.documentElement.classList.remove('dark')
	}
	$cookies.set('theme', darkTheme.value?'dark':'light', '30d')
}

onMounted(async () => {
	const VueCookies = await import('vue-cookies')
	darkTheme = ((function(){
		var themecookie = VueCookies.get('theme')
		var isdark = document.documentElement.classList.contains('dark') || themecookie === 'dark' ||
			(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
		$cookies.set('theme', isdark?'dark':'light', '30d')
		if(isdark){
			document.documentElement.classList.add('dark')
		}else{
			document.documentElement.classList.remove('dark')
		}
		return ref(isdark)
	})())

	observer = new MutationObserver((mutationList) => {
		for(let mutation of mutationList){
			if(mutation.type === "attributes"){
				if(mutation.attributeName === 'class'){
					darkTheme.value = document.documentElement.classList.contains('dark')
				}
			}
		}
	})
	observer.observe(document.documentElement, { attributes: true })
})

onUnmounted(() => {
	observer.disconnect()
})

</script>

<template>
	<header id="header">
		<div class="logo">
			<a href="/">
				<img src="/assets/favicon.png">
			</a>
		</div>
		<nav class="header-nav">
			<a class="hidden-for-mobile" href="/">{{$t(`word.home`)}}</a>
			<a class="hidden-for-mobile" href="/plugins">{{$t(`word.plugins`)}}</a>
			<a href="/about">{{$t(`word.about`)}}</a>
		</nav>
		<div class="locale-selector">
			<label for="i18n-lang-select">
				<WebBox class="flex-box" size="2rem" />
			</label>
			<select id="i18n-lang-select" v-model="$i18n.locale">
				<option v-for="locale in Object.keys(i18nLangMap)" :key="locale"
					:value="locale">{{i18nLangMap[locale] || locale}}</option>
			</select>
		</div>
		<div class="theme-button" @click="switchTheme">
			<IconMoon v-if="darkTheme"/>
			<IconSun v-else/>
		</div>
	</header>
	<div id="body">
		<RouterView v-slot="{ Component }"> 
			<!-- https://vuejs.org/guide/built-ins/suspense.html#combining-with-other-components -->
			<template v-if="Component">
				<Suspense>
					<component :is="Component"></component>
					<template #fallback>
						Loading...
					</template>
				</Suspense>
			</template>
		</RouterView>
	</div>
	<footer id="footer">
		<p>
			Powered by <a href="https://vuejs.org/">Vue.js</a>
		</p>
		<p>
			Sync from <a href="https://github.com/MCDReforged/PluginCatalogue">PluginCatalogue</a>.
			Svg icons from <a href="https://github.com/Templarian/MaterialDesign">MDI</a>
		</p>
		<address>
			Author: <a href="mailto:zyxkad@gmail.com">zyxkad@gmail.com</a> &copy;
		</address>
	</footer>
</template>

<style scoped>

#header {
	display: flex;
	flex-direction: row;
	align-items: center;
	width: 100%;
	height: 3.5rem;
	padding: 0 1rem;
	background: var(--color-background-soft);
	box-shadow: 0 0 0.5rem #000a;
	z-index: 10;
	overflow-x: scroll; /* TODO: suitable header for mobile */
}

#body {
	padding: 0 0.8rem;
	/*overflow: hidden;*/
}

#footer {
	width: 100%;
	height: 5rem;
	padding: 0.4rem;
	position: absolute;
	left: 0;
	bottom: 0;
	background: var(--color-background-soft);
}

.logo>a {
	display: inline-block;
	height: 3rem;
}

.logo>a:hover {
	background-color: unset;
}

.logo img {
	height: 3rem;
	width: 3rem;
}

.header-nav {
	display: flex;
	flex-direction: row;
	height: 100%;
	margin: 0 0.5rem;
}

.header-nav>a {
	display: flex;
	flex-direction: row;
	align-items: center;
	height: 100%;
	padding: 0 0.5rem;
	font-weight: 700;
}

.locale-selector {
	display: flex;
	flex-direction: row;
	align-items: center;
	margin-right: 1rem;
}

#i18n-lang-select {
	width: 5rem;
	height: 2rem;
	border-radius: 0.3rem;
	color: var(--color-text);
	background-color: var(--color-background);
}

.theme-button {
	min-width: 1.7rem;
	height: 1.7rem;
	padding: 0.25rem;
	border-radius: 0.25rem;
	background: var(--color-background-3);
	cursor: pointer;
}

@media (max-width: 30rem){
	.hidden-for-mobile {
		display: none !important;
	}
}

</style>
