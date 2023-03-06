<script setup>
import { defineProps } from 'vue'
import Information from 'vue-material-design-icons/Information.vue'
import ToolBox from 'vue-material-design-icons/ToolBox.vue'
import Controller from 'vue-material-design-icons/Controller.vue'
import ApiSvg from 'vue-material-design-icons/CloudPlus.vue'

defineProps({
	'id': String,
	'name': String,
	'desc': String,
	'authors': Array,
	'labels': Array,
})

const icons = {
	'information': Information,
	'tool': ToolBox,
	'management': Controller,
	'api': ApiSvg,
}

</script>

<template>
	<div class="plugin-item">
		<div class="name">
			{{name}}
			<span class="id">({{id}})</span>
		</div>
		<div class="authors">
			by
			<span v-for="(author, i) in authors">
				<span v-if="i">,</span>
				{{author}}
			</span>
		</div>
		<p class="description">
			<div v-if="desc">{{desc}}</div>
			<div v-else><i>No description</i></div>
		</p>
		<div class="labels">
			<div v-for="(ok, label) in labels">
				<div v-if="ok" class="label-item">
					<component :is="icons[label]" class="flex-box" size="1rem"/> {{label}}
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>

.plugin-item {
	width: 100%;
	height: 7.5rem;
	margin: 0.2rem 0;
	padding: 0.5rem;
	border-radius: 1rem;
	background-color: #fafafa;
	box-shadow: #888 0 0 0.2rem;
}

.name {
	display: inline-block;
	font-size: 1.3rem;
	font-weight: 600;
}

.id {
	font-family: monospace;
	font-size: 0.8rem;
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
	display: flex;
	flex-direction: row;
	align-items: center;
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
