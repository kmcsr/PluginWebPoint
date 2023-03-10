<script setup>
import { RouterLink, RouterView } from 'vue-router'
import WebBox from 'vue-material-design-icons/WebBox.vue'



</script>

<template>
	<header id="header">
		<div class="locale-changer">
			<WebBox class="flex-box" size="2rem" />
			<select v-model="$i18n.locale">
				<option v-for="locale in $i18n.availableLocales" :key="locale"
					:value="locale">{{locale}}</option>
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
	justify-content: flex-end;
	width: 100%;
	height: 3.5rem;
	background: #fff;
	z-index: 1;
}

#body {
	margin-bottom: 5rem;
  padding: 0 0.8rem;
}

#footer {
	width: 100%;
	height: 5rem;
	padding: 0.4rem;
	position: fixed;
	left: 0;
	bottom: 0;
	background: #fcfbf4;
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
