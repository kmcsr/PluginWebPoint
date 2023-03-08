<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
	<header id="header">
		TODO: header
	</header>
	<RouterView id="body"/>
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
/*	*/
}

#body {
	margin-bottom: 5rem;
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

</style>
