<script setup>
import { defineProps } from 'vue'
import { RouterLink } from 'vue-router'
import BriefcaseDownload from 'vue-material-design-icons/BriefcaseDownload.vue'
import UpdateSvg from 'vue-material-design-icons/Update.vue'
import LabelIcon from './LabelIcon.vue'
import { fmtTimestamp, sinceDate } from '../utils'

defineProps({
	'data': Object,
})

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
			<div class="labels">
				<div class="label-item" v-for="label in Object.entries(data.labels).filter(([k, ok])=>ok).map(([k, _])=>k).sort()">
					<LabelIcon :label="label" size="1rem"/>
				</div>
			</div>
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
	align-items: flex-end;
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
	padding-left: 0.8rem;
	width: 100%;
}

.label-item {
	display: inline-block;
	margin-right: 0.2rem;
	padding-right: 0.2rem;
	border-right: 0.08rem solid #999;
}

.label-item:last-child {
	margin-right: 0;
	padding-right: 0;
	border-right: none;
}

@media (max-width: 54.2rem){
	.plugin-item {
		flex-direction: column;
		min-height: 12rem;
		height: auto;
	}
	.plugin-body {
		width: 100%;
	}
	.plugin-extra {
		align-items: flex-start;
		width: 100%;
		margin-left: 0;
	}
	.label-item {
		border-right: none;
	}
}

</style>
