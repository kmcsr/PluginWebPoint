<script setup>
import { defineProps } from 'vue'
import Information from 'vue-material-design-icons/Information.vue'
import ToolBox from 'vue-material-design-icons/ToolBox.vue'
import Controller from 'vue-material-design-icons/Controller.vue'
import ApiSvg from 'vue-material-design-icons/CloudPlus.vue'

defineProps({
	'label': String,
	'text': String,
	'fillColor': String,
	'size': String,
})

const icons = {
	'information': Information,
	'tool': ToolBox,
	'management': Controller,
	'api': ApiSvg,
}
</script>

<template>
	<div class="label">
		<component :is="icons[label]" class="label-icon" :fill="fillColor" :size="size"/> {{text || label}}
	</div>
</template>

<style scoped>

.label {
	display: flex;
	flex-direction: row;
	align-items: center;
}

.label-icon {
	display: flex;
	margin-right: 0.3rem;
}

</style>
