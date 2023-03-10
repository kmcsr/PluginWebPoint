<script setup>
import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { NPagination } from 'naive-ui'
import { usePagination } from 'vue-request'
import TextSearch from 'vue-material-design-icons/TextSearch.vue'
import Filter from 'vue-material-design-icons/Filter.vue'
import SortAscending from 'vue-material-design-icons/SortAscending.vue'
import SortDescending from 'vue-material-design-icons/SortDescending.vue'
import axios from 'axios'
import router from '../router'
import LabelIcon from '../components/LabelIcon.vue'
import PluginItem from '../components/PluginItem.vue'

const pluginList = ref(null)
const pluginListHead = ref(null)
const pinHead = ref(false)
const errorText = ref(null)

const textFilter = ref('')
const tagFilters = ref([])
const sortBy = ref('downloads')
const reverseSort = ref(false)

const {
	data,
	loading: searching,
	run: refreshPluginList,
	current: listCurrentPage,
	totalPage,
	pageSize: listPageSize,
} = usePagination(() => {
	return Promise.all([
		getPluginList(), getPluginCounts()
	]).then(([res1, res2]) => {
		if(res1 && res2){
			res1.total = res2.total
			return res1
		}
	})
}, {
	errorRetryCount: 10,
	debounceInterval: 700,
	manual: true,
	pagination: {
		currentKey: 'page',
		pageSizeKey: 'limit',
		totalKey: 'total',
	},
})

function getPluginList(){
	let query = {}
	if(textFilter.value.length){
		query.q = textFilter.value
	}
	if(tagFilters.value.length){
		query.t = tagFilters.value.join(',')
	}
	if(sortBy.value.length){
		query.s = sortBy.value
	}
	if(reverseSort.value){
		query.reversed = 'true'
	}
	if(listCurrentPage.value > 1){
		query.pg = listCurrentPage.value
		query.ps = listPageSize.value
	}
	if(router.currentRoute.value.query !== query){
		router.push({ query: query })
	}
	return axios.get('/dev/plugins', {
		params: {
			filterBy: textFilter.value,
			tags: tagFilters.value.join(','),
			sortBy: sortBy.value,
			reversed: reverseSort.value,
			offset: (listCurrentPage.value - 1) * listPageSize.value,
			limit: listPageSize.value,
		}
	}).then((res) => {
		console.debug('response for /plugins:', res)
		if(res.data.status !== 'ok'){
			let err = new Error('Response status is not ok')
			err.response = res
			throw err
		}
		return res.data.data
	}).catch((error) => {
		console.error('Error when getting plugins:', error)
		if(error.response && error.response.data){
			errorText.value = error.response.data.error + ': ' + error.response.data.message
		}else{
			errorText.value = error.code + ': ' + error.message
		}
		throw error // keep pop error
	})
}

function getPluginCounts(){
	return axios.get('/dev/plugins/count', {
		params: {
			filterBy: textFilter.value,
			tags: tagFilters.value.join(','),
		}
	}).then((res) => {
		return res.data.data
	}).catch((error) => {
		console.error('Error when getting plugin count:', error)
		if(!errorText.value){
			if(error.response && error.response.data){
				errorText.value = error.response.data.error + ': ' + error.response.data.message
			}else{
				errorText.value = error.code + ': ' + error.message
			}
		}
		throw error // keep pop error
	})
}

async function refreshFunc(){
	errorText.value = null
	searching.value = true
	refreshPluginList()
}

watch(textFilter, refreshFunc)
watch(tagFilters, refreshFunc)
watch(sortBy, refreshFunc)
watch(reverseSort, refreshFunc)

const list = computed(() =>  (data.value) || [])

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
	console.debug('currentRoute:', router.currentRoute)
	let q = router.currentRoute.value.query
	if(q){
		if(q.q){
			textFilter.value = q.q
		}
		if(q.t){
			tagFilters.value = q.t.split(',')
		}
		if(q.s){
			sortBy.value = q.s
		}
		if(q.reversed){
			reverseSort.value = Boolean(q.reversed)
		}
		if(q.pg){
			listCurrentPage.value = Number.parseInt(q.pg) || 1
		}
		if(q.ps){
			listPageSize.value = Number.parseInt(q.ps) || 5
		}
	}
	refreshFunc()
})

onUnmounted(() => {
	window.removeEventListener('scroll', onScroll)
})

</script>

<template>
	<div class="plugin-list" ref="pluginList">
		<KeepAlive>
			<div>
				<div class="plugin-top-pages">
					<NPagination
						v-model:page="listCurrentPage"
						v-model:page-size="listPageSize"
						:page-count="totalPage"
						:page-slot="6"
						:page-sizes="[5, 15, 50, 100]"
						show-size-picker
					/>
				</div>
				<TransitionGroup name="pbody" tag="div">
					<div id="plugin-filter-teleport-slot"></div>
					<div v-if="searching" style="width:100%;min-height:6rem;display:flex;flex-direction:row;justify-content:center;align-items:center;">
						<b>Searching...</b>
					</div>
					<div v-else-if="errorText" class="error-box">
						{{errorText}}
					</div>
					<TransitionGroup v-else-if="list.length" class="plugin-list-body" name="plist" tag="div">
						<PluginItem  v-for="data in list" :key="data.id" :data="data"/>
					</TransitionGroup>
					<div v-else style="width:100%;min-height:6rem;display:flex;flex-direction:row;justify-content:center;align-items:center;">
						<b>No plugin was found</b>
					</div>
				</TransitionGroup>
				<div class="plugin-bottom-pages">
					<NPagination
						v-model:page="listCurrentPage"
						v-model:page-size="listPageSize"
						:page-count="totalPage"
						:page-slot="6"
						:page-sizes="[5, 15, 50, 100]"
						show-size-picker
					/>
				</div>
			</div>
		</KeepAlive>
		<div class="plugin-list-head" ref="pluginListHead" :class="pinHead?'plugin-list-head-pin':''">
			<div class="plugin-head-up">
				<div class="plugin-list-searchbox">
					<TextSearch class="flex-box plugin-list-search-icon" size="1.5rem"/>
					<input class="plugin-search-input" type="search" v-model="textFilter"
						placeholder="Search plugins..." />
				</div>
				<div class="plugin-list-filter-box">
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
							<option value="downloads">Downloads</option>
							<option value="lastUpdate">Last Update</option>
							<option value="name">Name</option>
							<option value="id">ID</option>
							<option value="authors">Authors</option>
						</select>
					</div>
				</div>
			</div>
			<Teleport to="#plugin-filter-teleport-slot" :disabled="pinHead"  v-if="showFilters">
				<div class="plugin-filters">
					<div>
						<input type="checkbox" id="plugin-filters-information" name="scales" value="information" v-model="tagFilters">
						<label for="plugin-filters-information">
							<LabelIcon label="information" text="Information" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-tool" name="scales" value="tool" v-model="tagFilters">
						<label for="plugin-filters-tool">
							<LabelIcon label="tool" text="Tool" class="flex-box" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-management" name="scales" value="management" v-model="tagFilters">
						<label for="plugin-filters-management">
							<LabelIcon label="management" text="Management" class="flex-box" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-api" name="scales" value="api" v-model="tagFilters">
						<label for="plugin-filters-api">
							<LabelIcon label="api" text="API" class="flex-box" size="1rem"/>
						</label>
					</div>
				</div>
			</Teleport>
		</div>
	</div>
</template>

<style scoped>

.plugin-list {
	display: flex;
	border-collapse: collapse;
	flex-direction: column;
	max-width: 100%;
	width: 53rem;
	margin: 0.6rem;
	padding-top: 0;
	border-radius: 0.2rem;
	background-color: var(--color-background);
	box-shadow: #0005 0 0 0.2rem;
	overflow: hidden scroll;
}

.plugin-list>*:first-child {
	margin-top: 3.55rem;
	min-height: 6rem;
}

.plugin-list-head {
	display: flex;
	flex-direction: column;
	width: 100%;
	z-index: 1;
	position: absolute;
	top: 0;
}

.plugin-list-head-pin {
	position: fixed;
}

.plugin-head-up {
	display: flex;
	flex-direction: row;
	align-items: center;
	min-height: 3.5rem;
	max-width: 53rem;
	width: 100%;
	padding: 0.3rem;
	border-radius: 0 0 0.7rem 0.7rem;
	box-shadow: #0007 0 0 0.1rem;
	background-color: #fff;
}

.plugin-list-searchbox {
	display: flex;
	flex-direction: row;
	align-items: center;
}

.plugin-list-search-icon {
	position: absolute;
	z-index: 1;
	left: 1.3rem;
}

.plugin-search-input {
	height: 2rem;
	width: 17rem;
	margin-left: 0.8rem;
	padding-left: 2.5rem;
	padding-right: 0.5rem;
	border-radius: 1rem;
	border: none;
	background-color: #e5e7eb;
}

.plugin-list-filter-box {
	display: flex;
	flex-direction: row;
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
	max-width: 53rem;
	width: 100%;
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
	width: 6rem;
	height: 2rem;
	margin-left: 0.3rem;
	padding-left: 0.3rem;
	border-radius: 0.43rem;
	border: none;
	background-color: #e5e7eb;
}

.plugin-list-pages {
	margin-left: 0.3rem;
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
	margin: 0.3rem 0;
	padding: 0 0.8rem;
	overflow: hidden;
}

.plugin-top-pages {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-top: 0.5rem;
}

.plugin-bottom-pages {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin: 1rem 0;
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

@media (max-width: 54.2rem){
	.plugin-list {
		margin-left: 0;
		margin-right: 0;
	}
	.plugin-list>*:first-child {
		margin-top: 5.2rem;
		min-height: 12rem;
	}
	.plugin-list-head-pin {
		left: 0;
	}
	.plugin-head-up {
		flex-direction: column;
	}
	.plugin-list-searchbox {
		width: 100%;
		margin-right: 0.8rem;
		margin-bottom: 0.3rem;
	}
	.plugin-search-input {
		width: 100%;
	}
	.plugin-list-filter-box {
		margin-bottom: 0.3rem;
	}
}

</style>
