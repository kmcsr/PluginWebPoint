<script setup>
import { ref } from 'vue'
import ClipboardText from 'vue-material-design-icons/ClipboardText.vue'

const props = defineProps({
	'text': String,
	'fillColor': String,
	'hoverColor': String,
	'hoverBackgroundColor': String,
})

const copySuccess = ref(false)

function onenter(event){
	copySuccess.value = false
}

async function copy(){
	try{
		await navigator.clipboard.writeText(props.text)
		copySuccess.value = true
		setTimeout(()=>{ copySuccess.value = false }, 1200)
	}catch(err){
		console.error('Cannot write to clipboard:', err)
	}
}

</script>

<template>
	<div class="box" @click="copy">
		<div class="hint-click">
			<span :style="copySuccess?{'display':'unset'}:{}">
				{{ copySuccess ?$t('message.copy_successed') :$t('message.click_to_copy') }}
			</span>
		</div>
		<code class="textarea" @mouseenter="onenter"
			:style="{
				'--hover-text-color': hoverColor || '#474747',
				'--hover-background-color': hoverBackgroundColor || '#18eea1'
			}">
			<span>
				{{text}}
			</span>
		</code>
		<ClipboardText class="clipboard-icon" :fill="fillColor" size="1.3rem"/>
	</div>
</template>

<style scoped>
.box:hover {
	cursor: pointer;
}

.textarea {
	--hover-text-color: #474747;
	--hover-background-color: #18eea1;
	display: block;
	padding: 0.2rem 1.5rem 0.2rem 0.4rem;
	margin: 0.4rem;
	border-radius: 0.5rem;
	background-color: var(--color-background-mute);
	user-select: none;
	white-space: nowrap;
	transition: all 0.3s ease-out;
}

.box:hover>.textarea {
	background-color: var(--hover-background-color);
	color: var(--hover-text-color);
}

.textarea>span {
	display: block;
	width: 100%;
	overflow-x: scroll;
	white-space: nowrap;
}

.textarea>span::-webkit-scrollbar {
	display: none;
}

.clipboard-icon {
	position: absolute;
	top: 0.3rem;
	right: 0.6rem;
}

.hint-click {
	display: flex;
	justify-content: center;
	position: absolute;
	top: calc(-100% - 0.1rem);
	left: 50%;
	width: 0;
	height: 1.5rem;
}

.hint-click>span {
	display: none;
	height: 100%;
	padding: 0.1rem 0.5rem;
	border-radius: 0.8rem;
	background-color: #000d;
	white-space: nowrap;
	color: #e0e0e0;
	font-size: 0.8rem;
	user-select: none;
}

.box:hover>.hint-click>span {
	display: unset;
}

.hint-click>span::after {
	content: ' ';
	display: block;
	position: absolute;
	top: 100%;
	left: calc(50% - 0.5rem);
	width: 0;
	height: 0;
	border: 0.5rem solid transparent;
	border-top-color: #000d;
}

</style>