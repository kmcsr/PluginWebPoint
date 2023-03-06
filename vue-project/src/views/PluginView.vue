<script setup>
import { defineProps, onMounted, ref } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import LabelIcon from '../components/LabelIcon.vue'

const props = defineProps({
	'plugin': String
})

const { data, run: freshData } = useRequest(() => {
	return axios.get('/dev/plugin/' + props.plugin + '/info').then((res) => {
		if(res.data.status !== 'ok'){
			let err = new Error('Response status is not ok')
			err.response = res
			throw err
		}
		return res.data.data
	}).catch((res) => {
		console.error('Error when fetching plugin data:', res)
	})
})

onMounted(() => {
	freshData()
})

</script>

<template>
	<main>
		<div v-if="data">
			<h2>{{data.name}}</h2>
			<span>{{data.id}}</span>
			<div>v{{data.version}}</div>
			<div>
				Last Update:
				<span v-if="data.lastUpdate">{{new Date(data.lastUpdate).toJSON()}}</span>
				<span v-else><i>Unknown</i></span>
			</div>
			<h3>Authors:</h3>
			<ul>
				<li v-for="author in data.authors">
					{{author}}
				</li>
			</ul>
			<p class="description">
				<div v-if="data.desc">{{data.desc}}</div>
				<div v-else><i>No description</i></div>
			</p>
			<ul class="labels">
				<li v-for="(ok, label) in data.labels">
					<LabelIcon v-if="ok" :label="label" size="1rem"/>
				</li>
			</ul>
			<ul>
				<li>
					<a :href="data.repo">Repo</a>
				</li>
				<li>
					<a :href="data.link">Link</a>
				</li>
			</ul>
		</div>
		<div v-else>
			Loading...
		</div>
	</main>
</template>

<style scoped>

</style>
