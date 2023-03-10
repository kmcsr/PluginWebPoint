
<script setup>
import { defineProps, defineEmits, onBeforeMount, onMounted, nextTick, ref, computed, watch } from 'vue'
import { RouterLink } from 'vue-router'
import router from '../router'

const props = defineProps({
	'data': Array,
	'default': String,
	'active': String,
	'replace': Boolean,
})

const emit = defineEmits(['update:active'])

const navAfter = ref(null)
const links = ref([])
var activeIndex = 0
const sOffset = ref(0)
const sWidth = ref(0)

function updateRouter(to){
	let j = 0
	for(let i in props.data){
		let d = props.data[i]
		if(d.id == props.default){
			j = i
		}
		let u = new URL(d.path, window.location.origin + window.location.pathname)
		if(u.pathname !== to.path){
			continue
		}
		if((d.exactHash || window.location.hash) && window.location.hash !== to.hash){
			continue
		}
		let flag = true
		if(d.exactQueryNames){
			for(let n of d.exactQueryNames){
				if((to.query[n] || '') !== (u.searchParams.get(n) || '')){
					flag = false
					break
				}
			}
			if(!flag){
				continue
			}
		}
		for(let [n, v] of u.searchParams.entries()){
			if(to.query[n] !== v){
				flag = false
				break
			}
		}
		if(!flag){
			continue
		}
		j = i
		break
	}
	activeIndex = j
	const tg = links.value[j]
	sOffset.value = tg.$el.offsetLeft
	sWidth.value = tg.$el.getBoundingClientRect().width
	emit('update:active', props.data[j].id)
}

watch(router.currentRoute, updateRouter)

for(let d of props.data){
	if(typeof d.text === 'function'){
		d._computed = computed(d.text)
		watch(d._computed, async () => {
			await nextTick()
			const tg = links.value[activeIndex]
			console.log('el:', tg.$el.offsetLeft, tg.$el.getBoundingClientRect())
			sOffset.value = tg.$el.offsetLeft
			sWidth.value = tg.$el.getBoundingClientRect().width
		})
	}else{
		d._computed = {
			value: d.text
		}
	}
}

onMounted(() => {
	updateRouter(router.currentRoute.value)
})

</script>

<template>
	<nav class="nav">
		<RouterLink ref="links" v-for="(d, i) in data" :key="d.id" :to="d.path" :replace="replace">
			{{d._computed ?d._computed.value :''}}
		</RouterLink>
		<div ref="navAfter" class="nav-after" :style="{ left: `${sOffset}px`, width: `${sWidth}px` }"></div>
	</nav>
</template>

<style scoped>
	
.nav {
	height: 2rem;
	border-bottom: 0.05rem #000 solid;
}

.nav-after {
	height: 0.08rem;
	background-color: #00d157;
	top: -0.05rem;
	transition: 0.3s all ease-out;
}

.nav>a {
	display: inline-block;
	height: 100%;
	padding: 0 0.2rem;
}

.nav>a.active {
	color: #000;
	font-weight: 400;
}

</style>
