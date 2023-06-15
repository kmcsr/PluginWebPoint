<script setup>
import { ref, watch, onMounted } from 'vue'
import IconThreeDots from './icons/IconThreeDots.vue'
import IconLeftArrowHead from './icons/IconLeftArrowHead.vue'
import IconRightArrowHead from './icons/IconRightArrowHead.vue'

const props = defineProps({
	'page': Number,
	'pageSize': Number,
	'pageCount': Number,
	'pageSlot': Number,
	'pageSizes': Array,
})

const emit = defineEmits([
	'update:page',
	'update:pageSize',
])

const pageLeft = ref([1])
const leftDots = ref(false)
const pageMiddle = ref(null)
const rightDots = ref(false)
const pageRight = ref([])

function render(){
	if(props.pageCount <= props.pageSlot){
		pageLeft.value = []
		leftDots.value = false
		pageMiddle.value = null
		rightDots.value = false
		pageRight.value = []
		for(let i = 1; i <= props.pageCount; i++){
			pageLeft.value.push(i)
		}
		return
	}

	const index = props.page
	const freeSlot = props.pageSlot - 2
	pageLeft.value = [1]
	if(index > freeSlot){ // 1 ... [4] 5 | 1 ... 4 [5] 6
		leftDots.value = true
	}else{
		leftDots.value = false
		for(let i = 2; i <= freeSlot; i++){
			pageLeft.value.push(i)
		}
	}
	pageMiddle.value = (index > freeSlot && index <= props.pageCount - freeSlot) ?index :null
	if(index <= props.pageCount - freeSlot){ // 1 [2] ... 5 | 1 ... [4] ... 7
		rightDots.value = true
		pageRight.value = props.pageCount > 1 ?[props.pageCount] :[]
	}else{
		rightDots.value = false
		pageRight.value = []
		for(let i = props.pageCount - freeSlot + 1; i <= props.pageCount; i++){
			pageRight.value.push(i)
		}
	}
	if(false){
		console.debug('pagination1: index=%d, freeSlot=%d, -freeSlot=%d', index, freeSlot, props.pageCount - freeSlot)
		console.debug('pagination2:',
			pageLeft.value.join(','),
			leftDots.value,
			pageMiddle.value,
			rightDots.value,
			pageRight.value)
	}
}

function updatePage(index){
	if(props.page === index){
		return
	}
	if(index <= 0 || index > props.pageCount){
		return
	}
	emit('update:page', index)
	render()
}

render()

onMounted(() => {
	watch(props, render)
})

</script>

<template>
	<div class="box">
		<div class="page-switch-button" :disabled="page == 1" @click="updatePage(page - 1)">
			<div>
				<IconLeftArrowHead/>
			</div>
		</div>
		<div v-for="i in pageLeft" :key="i"
			:class="['page-index', page === i?'page-index-active':'']"
			@click="updatePage(i)">
			{{i}}
		</div>
		<div v-if="leftDots">
			<IconThreeDots/>
		</div>
		<div v-if="pageMiddle"
			class="page-index page-index-active"
			@click="updatePage(pageMiddle)">
			{{pageMiddle}}
		</div>
		<div v-if="rightDots">
			<IconThreeDots/>
		</div>
		<div v-for="i in pageRight" :key="i"
			:class="['page-index', page === i?'page-index-active':'']"
			@click="updatePage(i)">
			{{i}}
		</div>
		<!--  -->
		<div class="page-switch-button" :disabled="page >= pageCount" @click="updatePage(page + 1)">
			<div>
				<IconRightArrowHead/>
			</div>
		</div>
		<select class="page-size-select"
			:value="pageSize || pageSizes[0]"
			@input="(event)=>{emit('update:pageSize', event.target.value)}">
			<option v-for="size in pageSizes">{{size}}</option>
		</select>
	</div>
</template>

<style scoped>
.box {
	display: flex;
	flex-direction: row;
	align-items: center;
	height: 1.6rem;
	margin-bottom: 0.5rem;
}

.box>* {
	display: inline-flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	width: 1.6rem;
	height: 1.6rem;
	border-radius: 4px;
	font-size: 0.8rem;
	color: var(--color-text);
	transition: color,background-color,border-color 0.5s ease;
	user-select: none;
	cursor: pointer;
}

.box>*:not(:first-child) {
	margin-left: 0.4rem;
}

.box>*:hover {
	border-color: #18a058;
}

.page-switch-button {
	border: 1px solid var(--color-text);
}

.page-switch-button[disabled=true] {
	background-color: var(--color-background-soft);
	border-color: var(--color-background-3);
	color: var(--color-background-3);
	cursor: not-allowed;
}

.page-switch-button>div {
	height: 1rem;
	width: 1rem;
}

.page-index:hover {
	color: #18a058;
}

.page-index-active {
	border: 1px solid #18a058;
	color: #18a058;
}

.page-size-select {
	border: 1px solid;
	width: auto;
	height: 100%;
}

.page-size-select::after {
	content: '/page';
}

</style>