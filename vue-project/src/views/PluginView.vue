<script setup>
import { defineProps, onMounted, ref } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import LabelIcon from '../components/LabelIcon.vue'

const props = defineProps({
	'plugin': String
})

const { data, run: freshData } = useRequest(() => {
	return Promise.all([axios.get(`/dev/plugin/${props.plugin}/info`).then((res) => {
		if(res.data.status !== 'ok'){
			let err = new Error('Response status is not ok')
			err.response = res
			throw err
		}
		return res.data.data
	}).catch((res) => {
		console.error('Error when fetching plugin data:', res)
	}), axios.get(`/dev/plugin/${props.plugin}/releases`).then((res) => {
		if(res.data.status !== 'ok'){
			let err = new Error('Response status is not ok')
			err.response = res
			throw err
		}
		return res.data.data
	})]).then((res) => {
		console.debug('res:', res)
		res[0].releases = res[1]
		return res[0]
	})
})

function fmtSize(size){
	let unit = 'B'
	if(size > 1024){
		size /= 1024
		unit = 'KB'
	}
	if(size > 1024){
		size /= 1024
		unit = 'MB'
	}
	if(size > 1024){
		size /= 1024
		unit = 'GB'
	}
	if(size > 1024){
		size /= 1024
		unit = 'TB'
	}
	return (+size.toFixed(2)) + unit
}

onMounted(() => {
	freshData()
})

</script>

<template>
	<main>
		<RouterLink to="/">Back to Index</RouterLink>
		<div v-if="data" class="plugin-box">
			<header class="plugin-head">
				<h1 class="plugin-name">
					{{data.name}}
					<span class="plugin-version">v{{data.version}}</span>
				</h1>
				<h2 class="plugin-authors">
					By
					<span v-for="author in data.authors">
						{{author}}
					</span>
				</h2>
			</header>
			<div>
				Last Update:
				<span v-if="data.lastUpdate">{{new Date(data.lastUpdate).toJSON()}}</span>
				<span v-else><i>Unknown</i></span>
			</div>
			<h2>Labels:</h2>
			<ul class="labels">
				<li v-for="label in Object.entries(data.labels).filter(([k, ok])=>ok).map(([k, _])=>k).sort()">
					<LabelIcon :label="label" size="1rem"/>
				</li>
			</ul>
			<h3>Total Download: {{data.downloads}}</h3>
			<h3>Links:</h3>
			<ul>
				<li>
					<a :href="data.repo">Repo: {{data.repo}}</a>
				</li>
				<li>
					<a :href="data.link">Link: {{data.link}}</a>
				</li>
			</ul>
			<p class="description">
				<pre>
					<div v-if="data.desc">{{data.desc}}</div>
					<div v-else><i>No description</i></div>
				</pre>
			</p>
			<h3>Releases:</h3>
			<table v-if="data.releases" style="border-collapse:collapse;">
				<thead>
					<th style="border: 1px solid #000;padding: 0.5rem;">File</th>
					<th style="border: 1px solid #000;padding: 0.5rem;">Size</th>
					<th style="border: 1px solid #000;padding: 0.5rem;">Downloads</th>
				</thead>
				<tbody>
					<tr v-for="r in data.releases.reverse()">
						<td style="border: 1px solid #000;padding: 0.5rem;">
							<a :href="`/download/${r.id}/${r.tag}/${r.filename}`">
								{{r.filename}}
							</a>
						</td>
						<td style="border: 1px solid #000;padding: 0.5rem;">{{fmtSize(r.size)}}</td>
						<td style="border: 1px solid #000;padding: 0.5rem;">{{r.downloads}}</td>
					</tr>
				</tbody>
			</table>
			<div v-else><i>No release</i></div>
		</div>
		<div v-else>
			Loading...
		</div>
	</main>
</template>

<style scoped>

.plugin-box {
	padding: 0.5rem;
	border: var(--color-border) 1px solid;
	border-radius: 1rem;
	background-color: var(--color-background);
}

.plugin-name {
/*	*/
}

.plugin-version {
	font-style: italic;
	font-size: 1rem;
	font-weight: 100;
}

.plugin-authors {
	font-size: 1rem;
	font-weight: 100;
}

.plugin-authors>span {
	margin-right: 0.2rem;
	font-size: 1.1rem;
	font-weight: 150;
}

.description {
	margin: 0.2rem;
}

</style>
