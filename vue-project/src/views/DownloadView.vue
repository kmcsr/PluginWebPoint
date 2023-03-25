<script setup>
import { defineProps, onMounted, onUnmounted, ref } from 'vue'
import { prefix as apiPrefix } from '../api'
import { setMetadata } from '../metadata'

const props = defineProps({
	'plugin': String,
	'tag': String,
	'filename': String,
})

const downloadURL = `${apiPrefix}/plugin/${props.plugin}/release/${props.tag}/asset/${props.filename}`

var redirectReject = null
var remain = ref(null)

function wait(n){
	return new Promise((reslove, reject)=>{
		redirectReject = reject
		setTimeout(()=>{reslove()}, n)
	})
}

function onPageShow(){
	remain.value = 3
}

const { unmount } = setMetadata({
	extras: [
		{ name: 'robots', content: 'noindex' },
	]
})

onMounted(async () => {
	window.addEventListener('pageshow', onPageShow)
	remain.value = 3;
	try{
		while(remain.value > 0){
			await wait(1000)
			remain.value--
		}
	}catch(err){
		return
	}
	redirectReject = null
	window.location.replace(downloadURL)
})

onUnmounted(() => {
	unmount()
	window.removeEventListener('pageshow', onPageShow)
	if(redirectReject){
		redirectReject()
		redirectReject = null
	}
})

</script>

<template>
	<main>
		<p class="hint">
			<center>
				<b>File will be download in {{remain}} sec...</b><br/>
				If not, click <a :href="downloadURL">here</a>
			</center>
		</p>
	</main>
</template>

<style scoped>

.hint {
	margin-top: 3rem;
}

</style>