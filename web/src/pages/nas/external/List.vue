<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="newHandler">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>
        {{ t('buttons.register') }}
      </n-button>
      <n-button secondary size="small" @click="fetchData()">
        <template #icon>
          <n-icon>
            <refresh-icon />
          </n-icon>
        </template>
        {{ t('buttons.refresh') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input size="small" v-model:value="args.name" :placeholder="t('fields.name')" clearable />
      <n-input size="small" v-model:value="args.rc_name" :placeholder="t('fields.rc_name')" clearable />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
    </n-space>
    <n-data-table remote :row-key="(row: any) => row.id" size="small" :columns="columns" :data="state.data"
      :pagination="pagination" :loading="state.loading" @update:page="fetchData" @update-page-size="changePageSize"
      scroll-x="max-content" />
  </n-space>
</template>

<script setup lang="ts">
import { reactive, watch } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
} from "naive-ui";
import {
  AddOutline as AddIcon,
  RefreshOutline as RefreshIcon,
} from "@vicons/ionicons5";
import { useRoute, useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import externalNasApi from "@/api/nas/external";
import type { ExternalNas } from "@/api/nas/external";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const newHandler = () => {
  router.push({ name: 'externalNas_new' })
}
const args = reactive({
  name: "",
  rc_name: "",
});
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    render: (row: ExternalNas) => renderLink({ name: 'externalNas_detail', params: { id: row.id } }, row.id),
  },
  {
    title: t('fields.name'),
    key: "name",
  },
  {
    title: t('fields.type'),
    key: "type",
  },
  {
    title: t('fields.rc_name'),
    key: "rc_name",
  },
  {
    title: t('fields.status'),
    key: "is_sync",
    render: (row: ExternalNas) => renderTag(
      row.is_sync ? t('enums.normal') : t('enums.need_sync'),
      row.is_sync ? "success" : "warning"
    ),
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt",
    render: (row: ExternalNas) => renderTime(row.updated_at),
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(row: ExternalNas, index: number) {
      return renderButtons([
        { type: 'warning', text: t('buttons.edit'), action: () => router.push({ name: 'externalNas_edit', params: { id: row.id } }) },
        { type: 'error', text: t('buttons.delete'), action: () => remove(row, index), prompt: t('prompts.delete') },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(externalNasApi.search, () => {
  return { ...args, filters: route.query.filters }
})
async function remove(u: ExternalNas, index: number) {
  await externalNasApi.delete(u.id);
  state.data.splice(index, 1)
}

watch(() => route.query.filter, () => {
  fetchData()
})
</script>