<script setup>
import { defineProps, onBeforeMount, onMounted, ref } from 'vue'
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

const errorText = ref(null)

const { data, run: freshData } = useRequest(() => {
	return Promise.all([axios.get(`/dev/plugin/${props.plugin}/info`).then((res) => {
		if(res.data.status !== 'ok'){
			let err = new Error('Response status is not ok')
			err.response = res
			throw err
		}
		const data = res.data.data
		document.title = `${data.name} | PWP`
		return data
	}).catch((error) => {
		console.error('Error when fetching plugin data:', error)
		if(error.response && error.response.data){
			errorText.value = error.response.data.error + ': ' + error.response.data.message
		}else{
			errorText.value = error.code + ': ' + error.message
		}
		throw res
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

function pluginDependUrl(id){
	if(id === 'mcdreforged'){
		return 'https://github.com/MCDReforged/PluginCatalogue'
	}
	return `/plugin/${id}`
}

onBeforeMount(() => {
	if(props.plugin === 'mcdreforged'){
		window.location.replace('https://github.com/MCDReforged/PluginCatalogue')
		return
	}
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
				{{ $t('message.lastUpdate') }}:&nbsp;
				<span v-if="data.lastUpdate">{{fmtTimestamp(sinceDate(data.lastUpdate), 1)}} {{ $t('word.ago') }}</span>
				<span v-else><i>{{ $t('word.unknown') }}</i></span>
			</div>
			<div v-if="data.github_sync" class="flex-box">
				<SyncSvg class="flex-box" size="1.5rem" style="margin-right:0.2rem;"/>
				{{ $t('message.synced_from_gh_1') }}
				<Github class="flex-box" style="margin: 0 0.1rem;" size="1rem"/>
				Github
				{{ $t('message.synced_from_gh_2') }}:&nbsp;
				<span v-if="data.last_sync">{{fmtTimestamp(sinceDate(data.last_sync), 1)}} {{ $t('word.ago') }}</span>
				<span v-else><i>{{ $t('word.unknown') }}</i></span>
			</div>
			<h2>{{ $t('word.labels') }}:</h2>
			<ul class="labels">
				<li v-for="label in Object.entries(data.labels).filter(([k, ok])=>ok).map(([k, _])=>k).sort()">
					<LabelIcon :label="label" :text="$t(`label.${label}`)" size="1rem"/>
				</li>
			</ul>
			<h3>
				<div class="flex-box">
					<BriefcaseDownload class="flex-box" size="1.5rem"/>
					{{ $t('message.totalDownload') }}: {{data.downloads}}
				</div>
			</h3>
			<h3>
				<div class="flex-box">
					<LinkBox class="flex-box" size="1.5rem"/>
					{{ $t('word.links') }}:
				</div>
			</h3>
			<ul>
				<li>
					<a :href="data.repo">{{ $t('message.repo') }}</a>
				</li>
				<li>
					<a :href="data.link">{{ $t('message.main_page') }}</a>
				</li>
			</ul>
			<p class="description">
				<pre style="white-space:break-spaces;">
					<div v-if="data.desc">{{$i18n.locale === 'zh_cn' ?data.desc_zhCN :data.desc}}</div>
					<div v-else><i>{{ $t('message.no_description') }}</i></div>
				</pre>
			</p>
			<h2>{{ $t('word.dependencies') }}:</h2>
			<table>
				<thead>
					<th>ID</th>
					<th>Tag</th>
				</thead>
				<tbody>
					<tr v-for="[id, cond] in Object.entries(data.dependencies)">
						<td>
							<a :href="pluginDependUrl(id)">
								{{id}}
							</a>
						</td>
						<td>{{cond}}</td>
					</tr>
				</tbody>
			</table>
			<h2>{{ $t('word.releases') }}:</h2>
			<table v-if="data.releases">
				<thead>
					<th>{{ $t('word.filename') }}</th>
					<th>{{ $t('word.size') }}</th>
					<th>{{ $t('word.downloads') }}</th>
				</thead>
				<tbody>
					<tr v-for="r in data.releases.reverse()">
						<td>
							<a :href="`/download/${r.id}/${r.tag}/${r.filename}`">
								{{r.filename}}
							</a>
						</td>
						<td>{{fmtSize(r.size)}}</td>
						<td>{{r.downloads}}</td>
					</tr>
				</tbody>
			</table>
			<div v-else><i>{{ $t('message.no_release') }}</i></div>
		</div>
		<div v-else-if="errorText" class="error-box">
			{{errorText}}
		</div>
		<div v-else>
			{{ $t('message.loading') }}
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
	font-weight: 250;
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

th, td {
	border: 1px solid #000;
	padding: 0.5rem;
}

</style>
