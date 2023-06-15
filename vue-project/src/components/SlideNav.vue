<script setup>
import { defineProps, defineEmits, onBeforeMount, onMounted, nextTick, ref, computed, watch } from 'vue'
import { useRouter, onBeforeRouteUpdate } from 'vue-router'

const router = useRouter()

const props = defineProps({
	'data': Array,
	'default': String,
	'active': String,
	'replace': Boolean,
	'between': String,
})

const emit = defineEmits(['update:active'])

const links = ref([])
const activeIndex = ref(0)
const sOffset = ref(0)
const sWidth = ref(0)

async function render(){
	if(links.value.length <= activeIndex.value){
		sOffset.value = 0
		sWidth.value = 0
		return
	}
	let tg, off, width
	let tgOffset = () => tg.offsetLeft - 1
	let tgWidth = () => tg.getBoundingClientRect().width + 2
	while((tg = links.value[activeIndex.value]) &&
		(sOffset.value != (off = tgOffset()) || sWidth.value != (width = tgWidth()))){
		// tg.activeClass = 'active'
		// console.log('tg:', tg, tg.$el.className)
		sOffset.value = off
		sWidth.value = width
		await nextTick()
	}
}

async function updateRouter(to, from){
	if(!from){
		from = to
	}
	let j = 0
	for(let i in props.data){
		let d = props.data[i]
		if(d.id == props.default){
			j = i
		}
		let u = new URL(d.path, 'http://example' + from.path)
		if(u.pathname !== to.path){
			continue
		}
		if(d.exactHash && u.hash !== to.hash){
			continue
		}
		let flag = false
		if(d.exactQueryNames){
			for(let n of d.exactQueryNames){
				if((to.query[n] || '') !== (u.searchParams.get(n) || '')){
					flag = true
					break
				}
			}
			if(flag){
				continue
			}
		}
		for(let [n, v] of u.searchParams.entries()){
			if(to.query[n] !== v){
				flag = true
				break
			}
		}
		if(flag){
			continue
		}
		j = i
		break
	}
	activeIndex.value = j
	emit('update:active', props.data[j].id)
	if(!import.meta.env.SSR){
		await render()
	}
}

for(let d of props.data){
	if(typeof d.text === 'function'){
		if(import.meta.env.SSR){
			d._computed = d.text()
		}else{
			d._computed = computed(d.text)
			watch(d._computed, async () => {
				await nextTick()
				await render()
			})
		}
	}else{
		d._computed = {
			value: d.text
		}
	}
}

await updateRouter(router.currentRoute.value, null)

onBeforeRouteUpdate((to, from) => {
	updateRouter(to, from)
})

function onNavItemSwitch(i, event){
	activeIndex.value = i;
	(props.replace ?router.replace :router.push)(props.data[i].path)
	emit('update:active', props.data[i].id)
	render()
}

onMounted(async () => {
	await render()
})

</script>

<template>
	<nav class="nav">
		<div class="nav-box">
			<a ref="links" v-for="(d, i) in data" :key="d.id" :href="d.path"
				:style="between && i?{'margin-left': between}:{}" :class="activeIndex == i ?'active' :''"
				@click.prevent.stop="(event) => onNavItemSwitch(i, event)" rel="nofollow">
				{{d._computed ?d._computed.value :''}}
			</a>
		</div>
		<div class="nav-after" :style="{ left: `${sOffset}px`, width: `${sWidth}px` }"></div>
	</nav>
</template>

<style scoped>
	
.nav {
	height: 2rem;
	border-bottom: 0.05rem var(--color-border-hover) solid;
	z-index: 10;
}

.nav-box {
	height: 100%;
	display: flex;
	flex-direction: row;
}

.nav-after {
	display: block;
	top: -0.1rem;
	height: 0.2rem;
	width: 0;
	border-radius: 0.1rem;
	background-color: #00d157;
	transition: 0.3s all ease-out;
}

.nav-box>a {
	display: inline-flex;
	height: 100%;
	padding: 0.1rem 0.6rem 0 0.6rem;
	border-radius: 0.2rem 0.2rem 0 0;
	transition: 0.3s all ease;
	white-space: nowrap;
}

.nav-box>a.active {
	color: var(--color-text);
	font-weight: 600;
}

.nav-box>a.active:hover {
	background-color: unset;
}

</style>
