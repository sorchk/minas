<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="newHandler">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>
        {{ t('buttons.new') }}
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
  useMessage,
} from "naive-ui";
import {
  AddOutline as AddIcon,
  RefreshOutline as RefreshIcon,
} from "@vicons/ionicons5";
import { useRoute, useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import schtaskApi from "@/api/sch/task";
import type { SchTask } from "@/api/sch/task";
import { runStatusMapping, typeMapping } from "@/api/sch/task";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const newHandler = () => {
  router.push({ name: 'schtask_new' })
}
const args = reactive({
  name: "",
});
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    render: (row: SchTask) => renderLink({ name: 'schtask_detail', params: { id: row.id } }, row.id),
  },
  {
    title: t('fields.name'),
    key: "name",
  },
  {
    title: t('fields.type'),
    key: "type",
    render: (row: SchTask) => typeMapping[row.type + ""],
  },
  {
    title: t('fields.cron'),
    key: "cron",
  },
  // {
  //   title: "运行状态",
  //   key: "last_status",
  //   render: (row: SchTask) => renderTag(
  //     runStatusMapping[row.last_status + ""].info,
  //     runStatusMapping[row.last_status + ""].type
  //   ),
  // },
  // {
  //   title: "最近运行时间",
  //   key: "last_run_time",
  //   render: (row: SchTask) => renderTime(row.last_run_time),
  // },
  // {
  //   title: "下次运行时间",
  //   key: "next_run_time",
  //   render: (row: SchTask) => renderTime(row.next_run_time),
  // },
  {
    title: t('fields.status'),
    key: "is_disable",
    render: (row: SchTask) => renderTag(
      row.is_disable ? t('enums.blocked') : t('enums.normal'),
      row.is_disable ? "warning" : "success"
    ),
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt",
    render: (row: SchTask) => renderTime(row.updated_at),
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(row: SchTask, index: number) {
      return renderButtons([
        { type: 'success', text: t('buttons.run_now'), action: () => exec(row, index), prompt: t('prompts.run_now') },
        { type: 'info', text: t('buttons.run_logs'), action: () => router.push({ name: 'schlog_list', params: { id: row.id } }) },
        row.is_disable ?
          { type: 'success', text: t('buttons.enable'), action: () => enable(row), prompt: t('prompts.enable'), } :
          { type: 'warning', text: t('buttons.block'), action: () => disable(row), prompt: t('prompts.block'), },
        { type: 'warning', text: t('buttons.edit'), action: () => router.push({ name: 'schtask_edit', params: { id: row.id } }) },
        { type: 'error', text: t('buttons.delete'), action: () => remove(row, index), prompt: t('prompts.delete') },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(schtaskApi.search, () => {
  return { ...args, filters: route.query.filters }
})
const message = useMessage()
async function enable(u: SchTask) {
  await schtaskApi.enable(u.id);
  fetchData()
}
async function disable(u: SchTask) {
  await schtaskApi.disable(u.id);
  fetchData()
}
async function remove(u: SchTask, index: number) {
  const data = await schtaskApi.delete(u.id);
  if (data.code === 200) {
    state.data.splice(index, 1)
    message.info(t('texts.action_success'));
  } else {
    message.info(t('texts.action_error'));
  }
}
async function exec(u: SchTask, index: number) {
  const data = await schtaskApi.exec(u.id);
  if (data.code === 200) {
    message.info(t('texts.action_success'));
  } else {
    message.info(t('texts.action_error'));
  }

}

watch(() => route.query.filter, (newValue: any, oldValue: any) => {
  fetchData()
})
</script>