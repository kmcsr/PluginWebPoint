<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import WebBox from 'vue-material-design-icons/WebBox.vue'
import IconSun from './components/icons/IconSun.vue'
import IconMoon from './components/icons/IconMoon.vue'
import { i18nLangMap } from './i18n'

const darkTheme = ref(document.documentElement.classList.contains('dark'))

function switchTheme(){
	if(darkTheme.value){
		document.documentElement.classList.remove('dark')
	}else{
		document.documentElement.classList.add('dark')
	}
}

const observer = new MutationObserver((mutationList) => {
	for(let mutation of mutationList){
		if(mutation.type === "attributes"){
			if(mutation.attributeName === 'class'){
				darkTheme.value = document.documentElement.classList.contains('dark')
			}
		}
	}
})

onMounted(() => {
	observer.observe(document.documentElement, { attributes: true })
})

onUnmounted(() => {
	observer.disconnect()
})

</script>

<template>
	<header id="header">
		<div class="logo">
			<b>LOGO required</b>
		</div>
		<nav class="header-nav">
			<a href="/about">About</a>
		</nav>
		<div class="locale-selector">
			<label for="i18n-lang-select">
				<WebBox class="flex-box" size="2rem" />
			</label>
			<select id="i18n-lang-select" v-model="$i18n.locale">
				<option v-for="locale in $i18n.availableLocales" :key="locale"
					:value="locale">{{i18nLangMap[locale] || locale}}</option>
			</select>
		</div>
		<div class="theme-button" @click="switchTheme">
			<IconMoon v-if="darkTheme"/>
			<IconSun v-else/>
		</div>
	</header>
	<div id="body">
		<RouterView/>
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
}

#body {
  padding: 0 0.8rem;
  overflow: scroll;
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

.logo {
	/*margin-right: 1rem;*/
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
}

.locale-selector {
	display: flex;
	flex-direction: row;
	align-items: center;
	margin-right: 1rem;
}

#i18n-lang-select {
	width: 7rem;
	height: 2rem;
	border-radius: 0.3rem;
	color: var(--color-text);
	background-color: var(--color-background);
}

.theme-button {
	width: 1.7rem;
	height: 1.7rem;
	padding: 0.25rem;
	border-radius: 0.25rem;
	background: var(--color-background-3);
}

</style>
