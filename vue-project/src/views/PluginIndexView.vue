<script setup>
import { onMounted, onBeforeUnmount, nextTick, ref, computed, watch } from 'vue'
import { useRouter, onBeforeRouteUpdate, RouterLink } from 'vue-router';
import { useRequest, usePagination } from 'vue-request'
import TextSearch from 'vue-material-design-icons/TextSearch.vue'
import Filter from 'vue-material-design-icons/Filter.vue'
import SortAscending from 'vue-material-design-icons/SortAscending.vue'
import SortDescending from 'vue-material-design-icons/SortDescending.vue'
import axios from 'axios'
import { prefix as apiPrefix } from '../api'
import LabelIcon from '../components/LabelIcon.vue'
import PluginItem from '../components/PluginItem.vue'
import Pagination from '../components/Pagination.vue'

const router = useRouter()

const pluginList = ref(null)
const pluginListHead = ref(null)
const pinHead = ref(false)
const errorText = ref(null)

const pageSlot = ref(5)

const defaultSortBy = 'downloads'
const defaultPageSize = 10

const {
	data,
	searching,
	totalPage,
	textFilter,
	tagFilters,
	sortBy,
	reverseSort,
	listCurrentPage,
	listPageSize,
} = (function(){
	console.debug('current query:', router.currentRoute.value.query)
	let q = router.currentRoute.value.query || {}
	return {
		data: ref(null),
		searching: ref(true),
		totalPage: ref(0),
		textFilter: ref(q.q || ''),
		tagFilters: ref(q.t ?q.t.split(',') :[]),
		sortBy: ref(q.s || defaultSortBy),
		reverseSort: ref(q.reversed === 'true'),
		listCurrentPage: ref(Number.parseInt(q.pg) || 1),
		listPageSize: ref(Number.parseInt(q.ps) || defaultPageSize),
	}
})()

async function getPluginCounts(){
	let res = await axios.get(`${apiPrefix}/plugins/count`, {
		params: {
			filterBy: textFilter.value,
			tags: tagFilters.value.join(','),
		}
	})
	return res.data.data
}

async function getPluginList(){
	let counts = (await getPluginCounts()).total
	totalPage.value = Math.ceil(counts / listPageSize.value) || 1
	if(listCurrentPage.value > totalPage.value){
		listCurrentPage.value = totalPage.value
	}
	if(!counts){
		return []
	}
	let res = await axios.get(`${apiPrefix}/plugins`, {
		params: {
			filterBy: textFilter.value,
			tags: tagFilters.value.sort().join(','),
			sortBy: sortBy.value,
			reversed: reverseSort.value,
			offset: (listCurrentPage.value - 1) * listPageSize.value,
			limit: listPageSize.value,
		}
	})
	res = res.data.data
	res.total = counts
	return res
}

async function refreshData(){
	try{
		searching.value = true
		data.value = await getPluginList()
		return data.value
	}catch(err){
		if(err.response && typeof err.response.data === 'object'){
			errorText.value = err.response.data.err + ': ' + err.response.data.message
		}else{
			errorText.value = err.code + ': ' + err.message
		}
		throw err
	}finally{
		searching.value = false
	}
}

async function refreshNoDelay(){
	errorText.value = null
	const oldQuery = router.currentRoute.value.query
	let query = {}
	let changed = false
	if(textFilter.value !== (oldQuery.q || '')){
		changed = 'textFilter'
		if(textFilter.value){
			query.q = textFilter.value
		}
	}else{
		query.q = oldQuery.q
	}
	let tags = tagFilters.value.sort().join(',')
	if(tags !== (oldQuery.t || '')){
		changed = 'tagFilters'
		if(tags){
			query.t = tags
		}
	}else{
		query.t = oldQuery.t
	}
	if(sortBy.value !== (oldQuery.s || defaultSortBy)){
		changed = 'sortBy'
		if(sortBy.value){
			query.s = sortBy.value
		}
	}else{
		query.s = oldQuery.s
	}
	let reversed = reverseSort.value ?'true' :''
	if(reversed !== (oldQuery.reversed || '')){
		changed = 'reverseSort'
		if(reverseSort.value){
			query.reversed = reversed
		}
	}else{
		query.reversed = oldQuery.reversed
	}
	if(listCurrentPage.value != (oldQuery.pg || 1)){
		changed = 'listCurrentPage'
		if(listCurrentPage.value > 1){
			query.pg = listCurrentPage.value.toString()
		}
	}else{
		query.pg = oldQuery.pg
	}
	if(listPageSize.value != (oldQuery.ps || defaultPageSize)){
		changed = 'listPageSize'
		if(listPageSize.value !== defaultPageSize){
			query.ps = listPageSize.value.toString()
		}
	}else{
		query.ps = oldQuery.ps
	}
	if(changed){
		console.debug('Search params changed:', changed, query)
		await router.push({ query: query })
	}
	return await refreshData()
}

if(import.meta.env.SSR){
	await refreshData()
}

const refreshDelayed = (function(){
	function _wait(t){
		return new Promise((re) => {
			setTimeout(()=>re(), t)
		})
	}
	var pending = null
	async function refreshDelayed(){
		if(pending){
			pending(true)
			return
		}
		while(await Promise.race([new Promise((re)=>{ pending = re }), _wait(600)]));
		pending = null
		return await refreshNoDelay()
	}
	return refreshDelayed
})()

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

function onQueryChange(value){
	let q = value.query
	if(q){
		if((q.q || '') !== textFilter.value){
			textFilter.value = q.q || ''
		}
		if((q.t || '') !== tagFilters.value.sort().join(',')){
			tagFilters.value = q.t ?q.t.split(',') :[]
		}
		if((q.s || defaultSortBy) !== sortBy.value){
			sortBy.value = q.s || defaultSortBy
		}
		if((q.reversed === 'true') !== reverseSort.value){
			reverseSort.value = q.reversed === 'true'
		}
		if((q.pg || 1) != listCurrentPage.value){
			listCurrentPage.value = Number.parseInt(q.pg) || 1
		}
		if((q.ps || defaultPageSize) != listPageSize.value){
			listPageSize.value = Number.parseInt(q.ps) || defaultPageSize
		}
	}
}

var mounting = false

onBeforeRouteUpdate((to, from) => {
	onQueryChange(to)
})

onMounted(() => {
	mounting = true
	window.addEventListener('scroll', onScroll)
	let { current, pageSize } = usePagination(({ page, limit }) => {
		return refreshNoDelay()
	}, {
		errorRetryCount: 10,
		pagination: {
			currentKey: 'page',
			pageSizeKey: 'limit',
		},
		defaultParams: [
			{
				page: 1,
				limit: defaultPageSize,
			}
		]
	})
	watch(listCurrentPage, (v) => (current.value = v))
	watch(listPageSize, (v) => (pageSize.value = v))
	watch(textFilter, () => {
		errorText.value = null
		searching.value = true
		refreshDelayed()
	})
	watch(tagFilters, refreshNoDelay)
	watch(sortBy, refreshNoDelay)
	watch(reverseSort, refreshNoDelay)
})

onBeforeUnmount(() => {
	mounting = false
	window.removeEventListener('scroll', onScroll)
})

</script>

<template>
	<div class="plugin-list" ref="pluginList">
		<KeepAlive>
			<div>
				<div id="plugin-filter-teleport-slot"></div>
				<div class="plugin-top-pages">
					<Pagination
						v-model:page="listCurrentPage"
						v-model:page-size="listPageSize"
						:page-count="totalPage"
						:page-slot="pageSlot"
						:page-sizes="[10, 15, 50, 100]"
					/>
				</div>
				<TransitionGroup name="pbody" tag="div">
					<div v-if="searching" class="searching-hint">
						<b>{{ $t('message.searching') }}</b>
					</div>
					<div v-if="errorText" class="error-box">
						{{ $t('message.error', { err: errorText }) }}
					</div>
					<TransitionGroup v-if="!errorText && list.length" class="plugin-list-body" name="plist" tag="div">
						<PluginItem  v-for="data in list" :key="data.id" :data="data"/>
					</TransitionGroup>
					<div v-if="!searching && !errorText && !(list.length)" class="searching-hint">
						<b>{{ $t('message.no_plugin') }}</b>
					</div>
				</TransitionGroup>
				<div class="plugin-bottom-pages">
					<Pagination
						v-model:page="listCurrentPage"
						v-model:page-size="listPageSize"
						:page-count="totalPage"
						:page-slot="pageSlot"
						:page-sizes="[10, 15, 50, 100]"
					/>
				</div>
			</div>
		</KeepAlive>
		<div class="plugin-list-head" ref="pluginListHead" :class="pinHead?'plugin-list-head-pin':''">
			<div class="plugin-head-up">
				<div class="plugin-list-searchbox">
					<TextSearch class="flex-box plugin-list-search-icon" size="1.5rem"/>
					<input class="plugin-search-input" type="search" v-model="textFilter"
						:placeholder="$t('message.search_plugins')" />
				</div>
				<div class="plugin-list-filter-box">
					<div class="plugin-filters-button">
						<button @click="showFilters=!showFilters">
							<Filter class="flex-box" size="1.5rem"/>
							{{ $t('word.filters') }}
						</button>
					</div>
					<div class="plugin-sorts">
						<label for="plugin-sorts-options">{{ $t('word.sortBy') }}</label>
						<component :is="reverseSort?SortDescending:SortAscending"
							class="flex-box plugin-sorts-icon" size="1.5rem" @click="reverseSort=!reverseSort"/>
						<select id="plugin-sorts-options" v-model="sortBy">
							<option value="downloads">{{ $t('word.downloads') }}</option>
							<option value="lastRelease">{{ $t('word.lastRelease') }}</option>
							<option value="name">{{ $t('word.name') }}</option>
							<option value="id">{{ $t('word.id') }}</option>
							<option value="authors">{{ $t('word.authors') }}</option>
						</select>
					</div>
				</div>
			</div>
			<Teleport to="#plugin-filter-teleport-slot" :disabled="pinHead" v-if="showFilters">
				<div class="plugin-filters">
					<div>
						<input type="checkbox" id="plugin-filters-information" name="scales" value="information" v-model="tagFilters">
						<label for="plugin-filters-information">
							<LabelIcon label="information" :text="$t(`label.information`)" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-tool" name="scales" value="tool" v-model="tagFilters">
						<label for="plugin-filters-tool">
							<LabelIcon label="tool" :text="$t(`label.tool`)" class="flex-box" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-management" name="scales" value="management" v-model="tagFilters">
						<label for="plugin-filters-management">
							<LabelIcon label="management" :text="$t(`label.management`)" class="flex-box" size="1rem"/>
						</label>
					</div>
					<div>
						<input type="checkbox" id="plugin-filters-api" name="scales" value="api" v-model="tagFilters">
						<label for="plugin-filters-api">
							<LabelIcon label="api" :text="$t(`label.api`)" class="flex-box" size="1rem"/>
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
	margin: 1.6rem;
	margin-top: 0;
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
	background-color: var(--color-background);
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
	background-color: var(--color-background-3);
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
	background-color: var(--color-background-3);
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
	background-color: var(--color-background-mute);
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
	background-color: var(--color-background-3);
	cursor: pointer;
}

#plugin-sorts-options {
	width: 6rem;
	height: 2rem;
	margin-left: 0.3rem;
	padding-left: 0.3rem;
	border-radius: 0.43rem;
	border: none;
	background-color: var(--color-background-3);
}

.plugin-list-pages {
	margin-left: 0.3rem;
}

.searching-hint {
	width: 100%;
	height: 6rem;
	max-height: 6rem;
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
}

.error-box {
	display: flex;
	align-items: center;
	text-indent: 1rem;
	width: 100%;
	min-height: 6rem;
	padding: 1rem;
	margin: 0 1rem;
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
	max-height: 0rem;
	opacity: 0;
	transform: translateY(-10px);
}

.plist-enter-from,
.plist-leave-to {
	opacity: 0;
	transform: translateX(30px);
}

.plist-leave-active {
	position: absolute;
}

@media (max-width: 54.2rem){
	.plugin-list {
		margin-left: 0;
		margin-right: 0;
	}
	.plugin-list>*:first-child {
		margin-top: 5.45rem;
		min-height: 12rem;
	}
	.plugin-list-head-pin {
		left: 0;
	}
	.plugin-head-up {
		flex-direction: column;
		padding-top: 0.5rem;
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
