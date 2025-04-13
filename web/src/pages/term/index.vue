<template>
    <x-page-header>
        <trem v-bind:termId="panel"></trem>
    </x-page-header>
</template>

<script lang="ts" setup>
import { computed, defineComponent, ref } from 'vue'

import { NTab, NTabs, NTabPane } from 'naive-ui'
import XPageHeader from "@/components/PageHeader.vue";
import Trem from "@/pages/term/term.vue";

const panel = ref(1) 
</script>
<style lang="css">
.summary {
    height: calc(100vh - 130px)
}
</style>