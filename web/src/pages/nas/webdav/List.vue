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
      <n-input size="small" v-model:value="args.account" :placeholder="t('fields.account')" clearable />
      <n-button size="small" type="primary" @click="() => fetchData()">{{ t('buttons.search') }}</n-button>
      <n-highlight :text="`Webdav地址：${webdavUrl}`" :patterns="patterns" :highlight-style="{
        padding: '0 6px',
        borderRadius: themeVars.borderRadius,
        display: 'inline-block',
        color: themeVars.baseColor,
        background: themeVars.primaryColor,
        transition: `all .3s ${themeVars.cubicBezierEaseInOut}`,
      }" />
    </n-space>
    <n-data-table remote :row-key="(row: any) => row.id" size="small" :columns="columns" :data="state.data"
      :pagination="pagination" :loading="state.loading" @update:page="fetchData" @update-page-size="changePageSize"
      scroll-x="max-content" />
  </n-space>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
  NHighlight,
  useThemeVars,
} from "naive-ui";
import {
  AddOutline as AddIcon,
  RefreshOutline as RefreshIcon,
} from "@vicons/ionicons5";
import { useRoute, useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import webdavApi from "@/api/nas/webdav";
import type { WebDav } from "@/api/nas/webdav";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'
import { baseUrl, frontBaseUrl } from '@/config';

const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const newHandler = () => {
  router.push({ name: 'webdav_new' })
}
const themeVars = useThemeVars()
const webdavUrl = window.location.href.substring(0, window.location.href.indexOf("/", 8)) +frontBaseUrl+ "/dav";
const patterns = [webdavUrl]
const args = reactive({
  name: "",
  account: "",
});
const columns = [
  {
    title: t('fields.id'),
    key: "id",
    render: (row: WebDav) => renderLink({ name: 'webdav_detail', params: { id: row.id } }, row.id),
  },
  {
    title: t('fields.name'),
    key: "name",
  },
  {
    title: t('fields.account'),
    key: "account",
  },
  {
    title: t('fields.root_dir'),
    key: "home",
  },
  {
    title: t('fields.status'),
    key: "is_disable",
    render: (row: WebDav) => renderTag(
      row.is_disable ? t('enums.blocked') : t('enums.normal'),
      row.is_disable ? "warning" : "success"
    ),
  },
  {
    title: t('fields.updated_at'),
    key: "updatedAt",
    render: (row: WebDav) => renderTime(row.updated_at),
  },
  {
    title: t('fields.actions'),
    key: "actions",
    render(row: WebDav, index: number) {
      return renderButtons([
        row.is_disable ?
          { type: 'success', text: t('buttons.enable'), action: () => enable(row) } :
          { type: 'warning', text: t('buttons.block'), action: () => disable(row), prompt: t('prompts.block'), },
        { type: 'warning', text: t('buttons.edit'), action: () => router.push({ name: 'webdav_edit', params: { id: row.id } }) },
        { type: 'error', text: t('buttons.delete'), action: () => remove(row, index), prompt: t('prompts.delete') },
      ])
    },
  },
];
const { state, pagination, fetchData, changePageSize } = useDataTable(webdavApi.search, () => {
  return { ...args, filters: route.query.filters }
})

async function enable(u: WebDav) {
  await webdavApi.enable(u.id);
  fetchData()
}
async function disable(u: WebDav) {
  await webdavApi.disable(u.id);
  fetchData()
}
async function remove(u: WebDav, index: number) {
  await webdavApi.delete(u.id);
  state.data.splice(index, 1)
}

watch(() => route.query.filter, (newValue: any, oldValue: any) => {
  fetchData()
})
</script>