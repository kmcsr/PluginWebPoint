<script setup>
import { ref } from 'vue'
import ClipboardText from 'vue-material-design-icons/ClipboardText.vue'

const props = defineProps({
	'text': String,
	'fillColor': String,
})

const showHint = ref(false)

async function copy(){
	try{
		await navigator.clipboard.writeText(props.text)
		showHint.value = true
		setTimeout(()=>{ showHint.value = false }, 1200)
	}catch(err){
		console.error('Cannot write to clipboard:', err)
	}
}

</script>

<template>
	<div>
		<button class="clipboard" @click="copy" :title="$t('message.click_to_copy')">
			<ClipboardText :fill="fillColor" size="1.3rem"/>
		</button>
		<div class="hint">
			<Transition name="hint">
				<span v-if="showHint">
					{{ $t('message.copy_successed') }}
				</span>
			</Transition>
		</div>
	</div>
</template>

<style scoped>

.clipboard {
	width: 1.5rem;
	height: 1.5rem;
	padding: 0.1rem;
	border: none;
	border-radius: 0.2rem;
	background-color: #d8d8d8;
	opacity: 0;
}

.clipboard:hover {
	background-color: #cacaca;
	cursor: pointer;
}

*:hover>div>.clipboard {
	opacity: 1;
}

.hint {
	display: flex;
	justify-content: center;
	position: absolute;
	top: -100%;
	left: 50%;
	width: 0;
	height: 1.5rem;
}

.hint>span {
	height: 100%;
	padding: 0.1rem 0.5rem;
	border-radius: 0.7rem;
	background-color: #000e;
	white-space: nowrap;
	color: #e0e0e0;
	font-size: 0.8rem;
	user-select: none;
}

.hint-enter-active,
.hint-leave-active {
	transition: opacity 0.2s ease;
}

.hint-enter-from,
.hint-leave-to {
	opacity: 0;
}

</style>
