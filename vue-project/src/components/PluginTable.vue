<script setup>
import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { useRequest } from 'vue-request'
import TextSearch from 'vue-material-design-icons/TextSearch.vue'
import Filter from 'vue-material-design-icons/Filter.vue'
import SortAscending from 'vue-material-design-icons/SortAscending.vue'
import SortDescending from 'vue-material-design-icons/SortDescending.vue'
import Information from 'vue-material-design-icons/Information.vue'
import ToolBox from 'vue-material-design-icons/ToolBox.vue'
import Controller from 'vue-material-design-icons/Controller.vue'
import ApiSvg from 'vue-material-design-icons/CloudPlus.vue'

import PluginItem from './PluginItem.vue'
import axios from 'axios'

const icons = {
	'information': Information,
	'tools': ToolBox,
	'management': Controller,
	'api': ApiSvg,
}

const pluginList = ref(null)
const pluginListHead = ref(null)
const pinHead = ref(false)
const errorText = ref(null)

const textFilter = ref('')
const tagFilters = ref([])
const sortBy = ref('name')
const reverseSort = ref(false)

const { data, loading: searching, run: refreshPluginList } = useRequest((event) => {
	if(event && event.type === 'input'){
		textFilter.value = event.target.value
	}
	errorText.value = null
	return axios.get('/dev/plugin/list', {
		params: {
			filterBy: textFilter.value,
			tags: tagFilters.value.join(','),
			sortBy: sortBy.value,
			reversed: reverseSort.value,
		}
	}).then((resp) => {
		if(resp.data.status !== 'ok'){
			console.error('response for /plugin/list:', resp)
			return null
		}
		return resp.data.data
	}).catch((error) => {
		console.error('Error when getting plugin list:', error)
		if(error.response && error.response.data){
			errorText.value = error.response.data.error + ': ' + error.response.data.message
		}else{
			errorText.value = error.code + ': ' + error.message
		}
		return null
	})
}, {
	debounceInterval: 300,
	manual: true,
})

watch(tagFilters, refreshPluginList)
watch(sortBy, refreshPluginList)
watch(reverseSort, refreshPluginList)

const list = computed(() =>  data.value || [])

const showFilters = ref(false)

function onScroll(event){
	if(pinHead.value){
		if(pluginList.value.getBoundingClientRect().y > 0){
			pinHead.value = false
		}
	}else{
		if(pluginListHead.value.getBoundingClientRect().y <= 0){
			pinHead.value = true
		}
	}
}

onMounted(() => {
	window.addEventListener('scroll', onScroll)
	searching.value = true
	refreshPluginList()
})

onUnmounted(() => {
	window.removeEventListener('scroll', onScroll)
})

</script>

<template>
	<div class="plugin-list" ref="pluginList">
		<KeepAlive>
			<TransitionGroup name="pbody" tag="div">
				<div v-if="showFilters" class="plugin-filters">
					<div>
						<input type="checkbox" id="plugin-filters-information" name="scales" value="information" v-model="tagFilters">
						<label for="plugin-filters-information">
							<Information class="flex-box" size="1rem"/>
							Information
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-tool" name="scales" value="tool" v-model="tagFilters">
						<label for="plugin-filters-tool">
							<ToolBox class="flex-box" size="1rem"/>
							Tool
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-management" name="scales" value="management" v-model="tagFilters">
						<label for="plugin-filters-management">
							<Controller class="flex-box" size="1rem"/>
							Management
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-api" name="scales" value="api" v-model="tagFilters">
						<label for="plugin-filters-api">
							<ApiSvg class="flex-box" size="1rem"/>
							API
						</label>
					</div>
				</div>
				<div v-if="searching" style="width:100%;min-height:6rem;display:flex;flex-direction:row;justify-content:center;align-items:center;">
					<b>Searching...</b>
				</div>
				<TransitionGroup v-else-if="list.length" class="plugin-list-body" name="plist" tag="div">
					<div v-for="data in list" :key="data.id">
						<PluginItem :id="data.id" :name="data.name" :authors="data.authors" :desc="data.desc" :labels="data.labels"/>
					</div>
				</TransitionGroup>
				<div v-else-if="errorText" class="error-box">
					{{errorText}}
				</div>
				<div v-else style="width:100%;min-height:6rem;display:flex;flex-direction:row;justify-content:center;align-items:center;">
					<b>No plugin was found</b>
				</div>
			</TransitionGroup>
		</KeepAlive>
		<div class="plugin-list-head" ref="pluginListHead" :style="pinHead?{position:'fixed',top:0}:{}">
			<TextSearch class="flex-box plugin-list-search-icon" size="1.5rem"/>
			<input class="plugin-list-searchbox" type="search" @input="refreshPluginList"
				placeholder="Search plugins..." />
			<div class="plugin-filters-button">
				<button @click="showFilters=!showFilters">
					<Filter class="flex-box" size="1.5rem"/>
					Filters
				</button>
			</div>
			<div class="plugin-sorts">
				<label for="plugin-sorts-options">Sort by</label>
				<component :is="reverseSort?SortDescending:SortAscending"
					class="flex-box plugin-sorts-icon" size="1.5rem" @click="reverseSort=!reverseSort"/>
				<select id="plugin-sorts-options" v-model="sortBy">
					<option value="name">Name</option>
					<option value="id">ID</option>
					<option value="authors">Authors</option>
				</select>
			</div>
			<!-- TODO: split pages -->
		</div>
	</div>
</template>

<style scoped>

.plugin-list {
	display: flex;
	border-collapse: collapse;
	flex-direction: column;
	width: 60rem;
	margin: 0.6rem;
	padding: 0.4rem;
	border-radius: 0.2rem;
	box-shadow: #0005 0 0 0.2rem;
	overflow: hidden scroll;
}

.plugin-list>*:first-child {
	margin-top: 4rem;
	min-height: 6rem;
}

.plugin-list-head {
	display: flex;
	flex-direction: row;
	align-items: center;
	height: 3.5rem;
	width: 59.2rem;
	border-radius: 0.7rem;
	box-shadow: #0007 0 0 0.1rem;
	background-color: #fff;
	z-index: 1;
	position: absolute;
	top: 0.4rem;
}

.plugin-list-search-icon {
	position: absolute;
	z-index: 1;
	left: 1.3rem;
}

.plugin-list-searchbox {
	height: 2rem;
	width: 14rem;
	margin-left: 0.8rem;
	padding-left: 2.5rem;
	padding-right: 0.5rem;
	border-radius: 1rem;
	border: none;
	background-color: #e5e7eb;
}

.plugin-filters-button>button:first-child {
	display: flex;
	flex-direction: row;
	align-items: center;
	width: 5.5rem;
	height: 2rem;
	margin-left: 0.3rem;
	padding-left: 0.3rem;
	border-radius: 0.43rem;
	border: none;
	background-color: #e5e7eb;
	box-shadow: #0008 0 0 0.1rem;
	user-select: none;
	cursor: pointer;
}

.plugin-filters-button>button:first-child:active {
	box-shadow: #0008 0 0 0.01rem;
}

.plugin-filters {
	width: 100%;
	margin-bottom: 0.5rem;
	padding: 0.4rem;
	border-radius: 0.3rem;
	background-color: #e8e8e8;
}

.plugin-filters>div {
	display: inline-flex;
	flex-direction: row;
	align-items: center;
	margin: 0.2rem 0.4rem;
}

.plugin-filters>div>label {
	display: inline-flex;
	flex-direction: row;
	align-items: center;
	margin-left: 0.1rem;
}

.plugin-sorts {
	display: flex;
	flex-direction: row;
	align-items: center;
	margin-left: 0.3rem;
}

.plugin-sorts-icon {
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
	width: 2rem;
	height: 2rem;
	margin-left: 0.3rem;
	border-radius: 0.3rem;
	background: #e5e7eb;
	cursor: pointer;
}

#plugin-sorts-options {
	width: 5.5rem;
	height: 2rem;
	margin-left: 0.3rem;
	padding-left: 0.3rem;
	border-radius: 0.43rem;
	border: none;
	background-color: #e5e7eb;
}

.error-box{
	display: flex;
	align-items: center;
	text-indent: 1rem;
	width: 100%;
	min-height: 6rem;
	padding: 1rem;
	border-radius: 1rem;
	border: #ff0000 solid 0.2rem;
	background: #ffcdc7;
	color: #ff0000;
	font-family: monospace;
	font-weight: 600;
}

.plugin-list-body {
	display: flex;
	flex-direction: column;
	width: 100%;
	padding: 0 0.4rem;
}

.pbody-enter-active,
.pbody-leave-active,
.plist-move,
.plist-enter-active,
.plist-leave-active {
	transition: all 0.5s ease;
}

.pbody-enter-from,
.pbody-leave-to {
	opacity: 0;
	transform: translateY(-10px);
}

.plist-enter-from,
.plist-leave-to {
	opacity: 0;
	transform: translateX(30px);
}

.pbody-leave-active,
.plist-leave-active {
	position: absolute;
}

.rotate-180 {
	animation: rotate-180 1s;
	transform: rotate(180deg);
}

.rotate-180-back {
	animation: rotate-180-back 1s;
	transform: rotate(0);
}

@keyframes rotate-180 {
	0% {
		transform: rotate(0);
	}
	100% {
		transform: rotate(-180deg);
	}
}

@keyframes rotate-180-back {
	0% {
		transform: rotate(-180deg);
	}
	100% {
		transform: rotate(0);
	}
}

.tb-row-1 {
	width: 20%;
}

.tb-row-2 {
	width: 20%;
}

.tb-row-3 {
	width: 45%;
}

.tb-row-4 {
	width: 15%;
}

</style>
