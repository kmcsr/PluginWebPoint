<script setup>
import { RouterLink, RouterView } from 'vue-router'
import WebBox from 'vue-material-design-icons/WebBox.vue'
import { i18nLangMap } from './i18n'

</script>

<template>
	<header id="header">
		<div class="logo">
			<b>LOGO required</b>
		</div>
		<nav class="header-nav">
			<a href="/about">About</a>
		</nav>
		<div class="locale-changer">
			<label for="i18n-lang-select">
				<WebBox class="flex-box" size="2rem" />
			</label>
			<select id="i18n-lang-select" v-model="$i18n.locale">
				<option v-for="locale in $i18n.availableLocales" :key="locale"
					:value="locale">{{i18nLangMap[locale] || locale}}</option>
			</select>
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
	background: #fdfdfd;
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
	background: #fcfbf4;
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

.locale-changer {
	display: flex;
	flex-direction: row;
	align-items: center;
	margin-right: 1rem;
}

.locale-changer>select {
	width: 7rem;
	height: 1.5rem;
}

</style>
