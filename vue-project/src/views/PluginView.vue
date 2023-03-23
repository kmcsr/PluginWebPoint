<script setup>
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useRequest } from 'vue-request'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import UpdateSvg from 'vue-material-design-icons/Update.vue'
import BriefcaseDownload from 'vue-material-design-icons/BriefcaseDownload.vue'
import SyncSvg from 'vue-material-design-icons/Sync.vue'
import Github from 'vue-material-design-icons/Github.vue'
import LinkBox from 'vue-material-design-icons/LinkBox.vue'
import DownloadBox from 'vue-material-design-icons/DownloadBox.vue'
import LabelIcon from '../components/LabelIcon.vue'
import SlideNav from '../components/SlideNav.vue'
import CopyableText from '../components/CopyableText.vue'
import { setMetadata } from '../metadata.js'
import { prefix as apiPrefix } from '../api'
import { fmtSize, fmtTimestamp, sinceDate, fmtDateTime } from '../utils'

const props = defineProps({
	'plugin': String,
})

if(props.plugin === 'mcdreforged'){
	window.location.replace('https://github.com/Fallen-Breath/MCDReforged')
}

const errorText = ref(null)

const { t } = useI18n()

const navActive = ref('readme')
const navData = [
	{
		id: 'readme',
		path: '',
		exactQueryNames: ['i'],
		text: () => t('word.readme'),
	},
	{
		id: 'depend',
		path: '?i=depend',
		exactQueryNames: ['i'],
		text: () => t('word.dependencies'),
	},
	{
		id: 'releases',
		path: '?i=releases',
		exactQueryNames: ['i'],
		text: () => t('word.releases'),
	}
]

var unmountCall = null

const { data, run: getPluginInfo } = useRequest(async () => {
	try{
		let res = await axios.get(`${apiPrefix}/plugin/${props.plugin}/info`)
		res = res.data.data
		if(!unmountCall){
			({ unmount: unmountCall } = setMetadata({
				title: `${res.name} | PWP`,
				keywords: [res.id, res.name],
				description: {
					'': res.desc,
					'zh': res.desc_zhCN || res.desc,
				}
			}))
		}
		return res
	}catch(err){
		console.error('Error when fetching plugin data:', err)
		if(err.response && err.response.data){
			errorText.value = err.response.data.err + ': ' + err.response.data.message
		}else{
			errorText.value = err.code + ': ' + err.message
		}
		throw err
	}
})

const { data: dataReadme, run: getPluginReadme } = useRequest(async () => {
	try{
		const res = await axios.get(`${apiPrefix}/plugin/${props.plugin}/readme`, {
			params: {
				render: true,
			}
		})
		const data = res.data
		return data
	}catch(err){
		if(err.response){
			if(err.response.status === 404){
				return false
			}
		}
		console.error('Error when getting readme:', err)
		throw err
	}
})

const { data: dataReleases, run: getPluginReleases } = useRequest(async () => {
	return (await axios.get(`${apiPrefix}/plugin/${props.plugin}/releases`)).data.data
})

const requireInstallCmd = computed(() => ((data.value && data.value.requirements) ? 
	("pip3 install " + Object.entries(data.value.requirements).map(([id, cond])=>`'${id}${cond}'`).join(' '))
	: ""))

async function freshData(){
	return await Promise.all([ getPluginInfo(), getPluginReadme(), getPluginReleases() ])
}

function pluginDependUrl(id){
	if(id === 'mcdreforged'){
		return 'https://github.com/Fallen-Breath/MCDReforged'
	}
	return `/plugin/${id}`
}

onMounted(() => {
	// freshData()
})

onUnmounted(() => {
	if(unmountCall){
		unmountCall()
		unmountCall = null
	}
})

</script>

<template>
	<main>
		<div v-if="data" class="plugin-box">
			<section class="plugin-section-box">
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
					{{ $t('message.lastRelease') }}:&nbsp;
					<span v-if="data.lastRelease">{{fmtTimestamp(sinceDate(data.lastRelease), 1)}} {{ $t('word.ago') }}</span>
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
					<div v-if="data.desc">{{$i18n.locale === 'zh_cn' ?data.desc_zhCN :data.desc}}</div>
					<div v-else><i>{{ $t('message.no_description') }}</i></div>
				</p>
			</section>
			<div class="plugin-main-box">
				<SlideNav :data="navData" default="readme" v-model:active="navActive" :replace="true"/>
				<article v-if="navActive === 'readme'" class="markdown-body plugin-readme"
					v-html="dataReadme === false?'<i>No readme :&lt;</i>' :(dataReadme || '<i>Loading ...</i>')">
				</article>
				<article v-else-if="navActive === 'depend'">
					<div v-if="data.dependencies || data.requirements">
						<div v-if="data.dependencies">
							<h2>{{ $t('word.dependencies') }}</h2>
							<hr style="margin-bottom:0.5rem;"/>
							<table style="margin-bottom: 1rem">
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
						</div>
						<div v-if="data.requirements">
							<h2>{{ $t('word.requirements') }}</h2>
							<hr style="margin-bottom:0.5rem;"/>
							<h3>{{ $t('word.install_cmd') }}</h3>
							<CopyableText :text="requireInstallCmd"/>
							<table>
								<thead>
									<th>Name</th>
									<th>{{ $t('word.require') }}</th>
								</thead>
								<tbody>
									<tr v-for="[id, cond] in Object.entries(data.requirements)">
										<td>
											<a :href="`https://pypi.org/project/${id}`">
												{{id}}
											</a>
										</td>
										<td>{{cond}}</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
					<div v-else>
						No dependencies
					</div>
				</article>
				<article v-else-if="navActive === 'releases'">
					<h2>{{ $t('word.releases') }}</h2>
					<div v-if="dataReleases">
						<div class="plugin-release" v-for="r in dataReleases">
							<a :href="`/download/${r.id}/${r.tag}/${r.filename}`">
								<DownloadBox class="flex-box release-download-icon" size="2rem"/>
								<div class="release-type-box">
									<div class="release-head">
										<div class="release-filename">
											<b>{{r.filename}}</b>
										</div>
										<div class="release-size">{{fmtSize(r.size)}}</div>
									</div>
									<div>
										<div class="release-download">{{ $t('word.downloads') }} <b>{{r.downloads}}</b></div>
										<div class="release-uploaded">{{ $t('word.published_at') }}
											<b style="white-space:nowrap;">{{fmtDateTime(r.uploaded)}}</b>
										</div>
									</div>
								</div>
							</a>
						</div>
					</div>
					<div v-else><i>{{ $t('message.no_release') }}</i></div>
				</article>
			</div>
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
	display: flex;
	flex-direction: row;
}

.plugin-section-box {
	min-width: 21rem;
	width: 21rem;
	height: fit-content;
	margin-top: 1rem;
	padding: 0.5rem;
	padding-bottom: 1rem;
	border: var(--color-border) 1px solid;
	border-radius: 1rem;
	background-color: var(--color-background);
	overflow: hidden;
}

.plugin-main-box {
	max-width: calc(100% - 21rem);
	width: 52rem;
	margin-top: 1rem;
	margin-bottom: 5rem;
	padding: 1rem;
	border: var(--color-border) 1px solid;
	border-radius: 1rem;
	background-color: var(--color-background);
}

.plugin-section-box>* {
	margin-bottom: 0.2rem;
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
	margin: 0.5rem 0.2rem;
	text-indent: 1rem;
	color: #888;
	font-family: monospace;
	white-space: break-spaces;
}

th, td {
	border: 0.05rem solid #000;
	padding: 0.5rem;
}

.plugin-main-box>article {
	padding: 0.5rem;
}

.plugin-release {
	width: 100%;
	height: 4rem;
	margin-bottom: 0.1rem;
}

.plugin-release>a {
	display: flex;
	flex-direction: row;
	width: 100%;
	height: 100%;
	border-radius: 1rem;
	padding: 0.5rem;
	color: #000;
}

.plugin-release:hover>a {
	background-color: #eee;
}

.plugin-release:active {

}

.release-type-box {
	display: flex;
	flex-direction: row;
	justify-content: space-between;
	width: 100%;
	font-weight: 500;
}

.release-download-icon {
	color: #00dc6e;
	margin-right: 0.5rem;
}

.release-size {
	margin-left: 0.5rem;
}

@media (max-width: 54.2rem){
	.plugin-box {
		flex-direction: column;
	}
	.plugin-section-box {
		width: 100%;
	}
	.plugin-main-box {
		max-width: 100%;
		width: 100%;
	}
	.plugin-release {
		height: 9rem;
	}
	.release-type-box {
		flex-direction: column;
	}
}

</style>
