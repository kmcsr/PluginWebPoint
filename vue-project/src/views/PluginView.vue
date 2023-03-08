<script setup>
import { defineProps, onMounted, ref } from 'vue'
import { useRequest } from 'vue-request'
import axios from 'axios'
import UpdateSvg from 'vue-material-design-icons/Update.vue'
import BriefcaseDownload from 'vue-material-design-icons/BriefcaseDownload.vue'
import SyncSvg from 'vue-material-design-icons/Sync.vue'
import Github from 'vue-material-design-icons/Github.vue'
import LinkBox from 'vue-material-design-icons/LinkBox.vue'
import LabelIcon from '../components/LabelIcon.vue'
import { fmtSize, fmtTimestamp, sinceDate } from '../utils'

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
	})]).then(([res1, res2]) => {
		if(!res1){
			return null
		}
		res1.releases = res2
		return res1
	})
})

onMounted(() => {
	freshData()
})

</script>

<template>
	<main>
		<div v-if="data" class="plugin-box">
			<header class="plugin-header">
				<RouterLink to="/">&lt;&lt;&nbsp;Back to Index</RouterLink>
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
			<div class="flex-box">
				<UpdateSvg class="flex-box" size="1.5rem" style="margin-right:0.2rem;"/>
				Last Update:&nbsp;
				<span v-if="data.lastUpdate">{{fmtTimestamp(sinceDate(data.lastUpdate), 1)}} ago</span>
				<span v-else><i>Unknown</i></span>
			</div>
			<div v-if="data.github_sync" class="flex-box">
				<SyncSvg class="flex-box" size="1.5rem" style="margin-right:0.2rem;"/>
				Synced from
				<Github class="flex-box" style="margin: 0 0.1rem;" size="1rem"/>
				Github:&nbsp;
				<span v-if="data.last_sync">{{fmtTimestamp(sinceDate(data.last_sync), 1)}} ago</span>
				<span v-else><i>Unknown</i></span>
			</div>
			<h2>Labels:</h2>
			<ul class="labels">
				<li v-for="label in Object.entries(data.labels).filter(([k, ok])=>ok).map(([k, _])=>k).sort()">
					<LabelIcon :label="label" size="1rem"/>
				</li>
			</ul>
			<h3>
				<div class="flex-box">
					<BriefcaseDownload class="flex-box" size="1.5rem"/>
					Total Download: {{data.downloads}}
				</div>
			</h3>
			<h3>
				<div class="flex-box">
					<LinkBox class="flex-box" size="1.5rem"/>
					Links:
				</div>
			</h3>
			<ul>
				<li>
					<a :href="data.repo">Repo</a>
				</li>
				<li>
					<a :href="data.link">Main page</a>
				</li>
			</ul>
			<p class="description">
				<pre style="white-space:break-spaces;">
					<div v-if="data.desc">{{data.desc}}</div>
					<div v-else><i>No description</i></div>
				</pre>
			</p>
			<h2>Releases:</h2>
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
	overflow: hidden;
}

.plugin-header {
	margin-bottom: 0.5rem;
}

.plugin-name {
	font-size: 1.5rem;
}

.plugin-version {
	font-style: italic;
	font-size: 1rem;
	font-weight: 100;
}

.plugin-authors {
	text-indent: 1rem;
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
