<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
	<main>
		<img class="main-img" src="/assets/intro-background.png"/>
		<hr/>
		<div class="cards">
			<div class="q-a">
				<h3 class="question">What is MCDReforged?</h3>
				<p class="answer">
					<a href="https://github.com/Fallen-Breath/MCDReforged"><i><b>MCDReforged</b></i></a>
					is a Minecraft server control tool rewritten from <i><b>MCDaemon</b></i>,
					which can improve the management ability of the Minecraft server <b>without</b> <i>modification</i> of server,
					and supports <b>plugins</b> written by <i>Python</i>.
				</p>
			</div>
			<RouterLink class="searching-link" to="/plugins">
				<h3>-> Search for MCDReforged Plugins</h3>
			</RouterLink>
		</div>
	</main>
</template>

<style scoped>

.main-img {
	width: 100%;
}

.cards>div {
	display: inline-block;
}

.q-a {
	max-width: 40rem;
	height: 9rem;
	padding: 2rem;
	margin: 2rem;
	margin-bottom: 10rem;
	border-radius: 1rem;
	background-color: var(--color-background-soft);
	box-shadow: #0004 0 0 1rem;
}

.question {
	width: 80%;
	font-size: 1.5rem;
}

.answer {
	width: 90%;
	padding: 2rem;
	border-radius: 1rem;
	background-color: var(--color-background-soft);
	box-shadow: #0004 0 0 1rem;
	font-size: 1rem;
	text-indent: 2rem;
	float: left;
	top: 1rem;
	left: 2rem;
}

.searching-link {
	display: inline-block;
	padding: 1.2rem;
	padding-right: 2rem;
	margin: 2rem 0;
	border-radius: 1.3rem;
	background-color: var(--color-background-soft);
	box-shadow: #0004 0 0 1rem;
	font-size: 1.3rem;
}

.searching-link:hover {
	cursor: pointer;
}

</style>
