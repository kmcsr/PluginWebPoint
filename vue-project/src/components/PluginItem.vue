<script setup>
import { defineProps } from 'vue'
import { RouterLink } from 'vue-router'
import BriefcaseDownload from 'vue-material-design-icons/BriefcaseDownload.vue'
import UpdateSvg from 'vue-material-design-icons/Update.vue'
import LabelIcon from './LabelIcon.vue'

defineProps({
	'data': Object,
})

function fmtTimestamp(ts, n){
	if(n === undefined){
		n = 2;
	}
	let unit = 'ms'
	if(ts > 1000){
		ts /= 1000
		unit = 's'
	}
	if(ts > 60){
		ts /= 60
		unit = 'min'
	}
	if(ts > 60){
		ts /= 60
		unit = 'h'
	}
	if(ts > 24){
		ts /= 24
		unit = 'd'
	}
	return (+ts.toFixed(n)) + unit
}

function sinceDate(date){
	return new Date() - new Date(date)
}

</script>

<template>
	<article class="plugin-item">
		<div class="plugin-body">
			<div class="name">
				<RouterLink :to="'/plugin/' + data.id">
					{{data.name}}
				</RouterLink>
			</div>
			<div class="authors">
				by
				<span v-for="(author, i) in data.authors">
					<span v-if="i">,</span>
					{{author}}
				</span>
			</div>
			<p class="description">
				<div v-if="desc">{{data.desc}}</div>
				<div v-else><i>No description</i></div>
			</p>
		</div>
		<div class="plugin-extra">
			<div>
				<BriefcaseDownload class="flex-box" size="1.5rem" style="margin-right:0.2rem;"/>
				<b class="plugin-downloads">{{data.downloads}}</b> downloads
			</div>
			<div>
				<UpdateSvg class="flex-box" size="1.5rem" style="margin-right:0.2rem;"/>
				Updated {{fmtTimestamp(sinceDate(data.lastUpdate), 1)}} ago
			</div>
		</div>
		<div class="labels">
			<LabelIcon class="label-item"
				v-for="label in Object.entries(data.labels).filter(([k, ok])=>ok).map(([k, _])=>k).sort()"
				:label="label" size="1rem"/>
		</div>
	</article>
</template>

<style scoped>

.plugin-item {
	display: flex;
	flex-direction: row;
	width: 100%;
	height: 7.5rem;
	margin: 0.2rem 0;
	padding: 0.5rem;
	border-radius: 1rem;
	background-color: #fafafa;
	box-shadow: #888 0 0 0.2rem;
}

.plugin-body {
	width: 77%;
	height: 100%;
}

.plugin-extra {
	display: flex;
	flex-direction: column;
	align-items: right;
	width: 23%;
	margin-left: 1rem;
}

.plugin-extra>div {
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: right;
	margin-bottom: 0.5rem;
}

.plugin-downloads {
	font-size: 1.25rem;
	font-weight: 700;
	margin-right: 0.3rem;
}

.name {
	display: inline-block;
	font-size: 1.3rem;
	font-weight: 600;
}

.authors {
	display: inline-block;
	margin-left: 0.4rem;
	font-weight: 250;
}

.description {
	text-indent: 1rem;
	font-family: monospace;
	font-weight: 100;
}

.description>div {
	display: block;
	height: 3rem;
	overflow: hidden;
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 2;
}

.labels {
	display: flex;
	flex-direction: row;
	position: absolute;
	left: 0.6rem;
	bottom: 0.5rem;
}

.label-item {
	margin-left: 0.2rem;
	padding-left: 0.2rem;
	border-left: 0.08rem solid #999;
}

.label-item:first-child {
	margin-left: 0;
	padding-left: 0;
	border-left: none;
}

</style>
