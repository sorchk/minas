<template>
  <x-page-header :subtitle="model.schtask.name">
    <template #action>
      <n-button secondary size="small" @click="newHandler">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
      <n-button secondary size="small" @click="editHandler">{{ t('buttons.edit') }}</n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="16">
    <x-description cols="1 640:1" label-position="left" label-align="right" :label-width="100">
      <x-description-item :label="t('fields.id')">{{ model.schtask.id }}</x-description-item>
      <x-description-item :label="t('fields.name')">{{ model.schtask.name }}</x-description-item>
      <x-description-item :label="t('fields.type')">{{ model.schtask.type }}
        <n-button strong secondary size="small" circle type="primary" @click="copyText(model.schtask.type)">
          <template #icon>
            <n-icon>
              <copy-icon />
            </n-icon>
          </template>
        </n-button></x-description-item>

      <x-description-item :label="t('fields.cron')">{{ model.schtask.cron }}</x-description-item>
      <x-description-item :label="t('fields.status')">
        <n-tag round size="small" :type="model.schtask.is_disable == 0 ? 'primary' : 'warning'">{{
          t(model.schtask.is_disable == 0 ? 'enums.normal' : 'enums.blocked') }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        {{ model.schtask.created_at }}
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        {{ model.schtask.updated_at }}
      </x-description-item>
    </x-description>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, watch } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTime,
  useMessage,
} from "naive-ui";
import { useRoute, useRouter } from "vue-router";
import { ArrowBackCircleOutline as BackIcon, CopyOutline as CopyIcon } from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XAnchor from "@/components/Anchor.vue";
import { XDescription, XDescriptionItem } from "@/components/description";
import schtaskApi from "@/api/sch/task";
import type { SchTask } from "@/api/sch/task";
import { copyText } from "@/utils";
import { useI18n } from 'vue-i18n'

const message = useMessage()
const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const model = reactive({
  schtask: {} as SchTask,
});
const newHandler = () => {
  router.push({ name: 'schtask_list' })
}
const editHandler = () => {
  router.push({ name: 'schtask_edit', params: { id: model.schtask.id } })
}

async function fetchData() {
  let schtask = await schtaskApi.load(route.params.id as string);
  model.schtask = schtask.data as SchTask
}

watch(() => route.params.id, fetchData)

onMounted(fetchData);
</script>