<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="listHandler">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
      <n-popconfirm @positive-click="exec">
        <template #trigger>
          <n-button secondary size="small">
            <template #icon>
              <n-icon>
                <chevron-forward-circle-outline />
              </n-icon>
            </template>
            {{ t('buttons.execute') }}
          </n-button>
        </template>
        {{ t('prompts.execute') }}
      </n-popconfirm>

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
    <n-layout has-sider>
      <n-layout-sider bordered collapse-mode="width" class="autoHeight" :collapsed-width="64" :width="240">
        <n-menu :options="menuOptions" :render-icon="renderMenuIcon" style="width: 220px" v-model:value="state.selectId"
          @click="Load()" />
      </n-layout-sider>
      <n-layout style="margin-left: 10px;">
        <n-descriptions label-placement="top" v-if="state.selectId" :title="t('fields.run_details')">
          <n-descriptions-item :label="t('fields.start_time')">
            {{ state.data.start_time }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.end_time')">
            {{ state.data.end_time }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.status')">
            <n-tag :type="statusMapping[state.data.status]?.type || 'info'">
              {{ statusMapping[state.data.status]?.info || t('fields.unknown') }}</n-tag>
          </n-descriptions-item>
        </n-descriptions>
        <n-log :rows="50" :log="state.data.log_text" />
      </n-layout>
    </n-layout>
  </n-space>
</template>

<script setup lang="ts">
import { h, onBeforeUnmount, reactive, ref, watch } from "vue";
import hljs from 'highlight.js/lib/core'
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
  useMessage,
  NSplit,
  lightTheme,
  NLayout,
  NLayoutSider,
  NLayoutContent,
  NLayoutFooter,
  NLayoutHeader,
  NPopconfirm,
  NLog,
  NTag,
  NEllipsis,
  NMenu,
  NDescriptions,
  NDescriptionsItem,
} from "naive-ui";
import {
  ChevronForwardCircleOutline,
  RefreshOutline as RefreshIcon,
  BookmarkOutline, CaretDownOutline,
  ArrowBackCircleOutline as BackIcon, CopyOutline as CopyIcon
} from "@vicons/ionicons5";
import type { MenuOption } from 'naive-ui'
import { useRoute, useRouter } from "vue-router";
import XPageHeader from "@/components/PageHeader.vue";
import schtaskApi from "@/api/sch/task";
import schlogApi from "@/api/sch/logs";
import type { SchLog } from "@/api/sch/logs";
import { statusMapping } from "@/api/sch/logs";
import { renderLink, renderTag, renderTime } from "@/utils/render";
import { useI18n } from 'vue-i18n'
import { onMounted } from "vue";
import { stat } from "fs";
import { xxtea } from "@/utils/xxtea";
const split = ref<number>(0.3)
const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const intervalId = ref()
const intervalId2 = ref()
const message = useMessage()
const state = reactive({
  loading: false,
  name: "",
  selectId: "",
  data: {} as any
});
hljs.registerLanguage('naive-log', () => ({
  keywords: '是 否',
  literal: ['false', 'true', 'null', 'undefined', 'NaN', 'Infinity'],
  contains: [
    {
      scope: 'string',
    }
  ]
}))
const renderMenuIcon = (option: MenuOption) => {
  // 渲染图标占位符以保持缩进
  if (option.status == "1") {
    return h(NIcon, { color: '#18a058' }, { default: () => h(BookmarkOutline) })
  } else if (option.status == "0") {
    return h(NIcon, { color: 'blue' }, { default: () => h(BookmarkOutline) })
  } else if (option.status == "-2") {
    return h(NIcon, { color: '#f0a020' }, { default: () => h(BookmarkOutline) })
  } else {
    return h(NIcon, { color: 'red' }, { default: () => h(BookmarkOutline) })
  }
}
const menuOptions = ref(new Array())
async function exec() {
  const data = await schtaskApi.exec(route.params.id);
  if (data.code === 200) {
    await fetchData(true)
    console.log("data:", data)

    message.info(t('texts.action_success'));
  } else {
    message.info(t('texts.action_error'));
  }

}
const fetchData = async (selectFirst = false) => {
  const args = { page: 1, size: 50 } as any;
  args.columns = xxtea.encryptAuto(JSON.stringify("id,task_id,status,start_time,end_time".split(",")), "columns");
  args.sorts = xxtea.encryptAuto(JSON.stringify(["-start_time"]), "sorts");
  args.filters = xxtea.encryptAuto(JSON.stringify([{ "Column": "task_id", "Operator": "=?", "Value": route.params.id as string }]), "filters");
  let list = await schlogApi.search(args);
  if (list.code == 200 && list.data) {
    menuOptions.value = new Array()
    list.data.forEach((item: SchLog, index: number) => {
      if (index == 0 && (selectFirst || !state.selectId)) {
        state.selectId = item.id
        Load()
      }
      menuOptions.value.push({
        label: () =>
          h(NEllipsis, null, { default: () => item.start_time }),
        key: item.id,
        status: item.status
      })
    });
    console.log(menuOptions)
  }

}
const listHandler = () => {
  router.push({ name: 'schtask_list' })
}
const Load = () => {
  schlogApi.load(state.selectId).then(data => {
    if (data.code == 200) {
      state.data = data.data
    } else {
      state.data = {};
    }
  })
}
onMounted(() => {
  fetchData()
  intervalId.value = setInterval(Load, 3000);
  intervalId2.value = setInterval(fetchData, 5000);
})
onBeforeUnmount(() => {
  clearInterval(intervalId.value);
  clearInterval(intervalId2.value);
})

</script>
<style scoped lang="scss"></style>