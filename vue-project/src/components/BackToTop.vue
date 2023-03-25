<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import ArrowUpDropCircle from 'vue-material-design-icons/ArrowUpDropCircle.vue'

const props = defineProps({
	'position': String,
	'left': String,
	'top': String,
	'right': String,
	'bottom': String,
	'size': String,
	'background': String,
	'color': String,
})

const active = ref(false)

function onClick(){
	window.scrollTo({ top: 0, behavior: 'smooth' })
}

function onScroll(event){
	active.value = window.innerHeight < window.scrollY
}

onMounted(() => {
	window.addEventListener('scroll', onScroll)
})

onUnmounted(() => {
	window.removeEventListener('scroll', onScroll)
})

</script>

<template>
	<div class="box" :class="active?'active':''" :style="{
		'--position': position || 'fixed',
		'--left-offset': left || 'unset',
		'--top-offset': top || 'unset',
		'--right-offset': right || 'unset',
		'--bottom-offset': bottom || 'unset',
		'--size': size || '2rem',
		'--background': background || '#0000',
		'--color': color || 'currentColor',
	}" @click="onClick">
		<ArrowUpDropCircle :size="size"/>
	</div>
</template>

<style scoped>
.box {
	position: var(--position);
	left: calc(var(--left-offset) - var(--size));
	top: var(--top-offset);
	right: var(--right-offset);
	bottom: var(--bottom-offset);
	width: var(--size);
	height: var(--size);
	border-radius: calc(var(--size) / 2);
	background: var(--background);
	color: var(--color);
	opacity: 0;
	box-shadow: 0 0 0.6rem -0.2rem #000;
	z-index: 999;
	cursor: pointer;
	transition: all 0.5s ease-out;
}

.box.active {
	left: var(--left-offset);
	opacity: 1;
}
</style>