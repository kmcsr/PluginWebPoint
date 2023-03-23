<script setup>
import { defineProps, onMounted, onUnmounted, ref } from 'vue'
import { prefix as apiPrefix } from '../api'

const props = defineProps({
	'plugin': String,
	'tag': String,
	'filename': String,
})

var redirectReject = null
var remain = ref(null)

function wait(n){
	return new Promise((reslove, reject)=>{
		redirectReject = reject
		setTimeout(()=>{reslove()}, n)
	})
}

onMounted(() => {
	remain.value = 3;
	(async function(){
		try{
			while(remain.value > 0){
				await wait(1000)
				remain.value--
				console.log(1)
			}
		}catch(_){
			return
		}
		redirectReject = null
		window.location.replace(`${apiPrefix}/plugin/${props.plugin}/release/${props.tag}/asset/${props.filename}`)
	})()
})

onUnmounted(() => {
	if(redirectReject){
		redirectReject()
		redirectReject = null
	}
})

</script>

<template>
	<main>
		<pre>
			<center>File will be download in {{remain}} sec...</center>
		</pre>
	</main>
</template>

<style scoped>

</style>