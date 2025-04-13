<template>
  <x-page-header :subtitle="model.webdav.name">
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
      <x-description-item :label="t('fields.id')">{{ model.webdav.id }}</x-description-item>
      <x-description-item :label="t('fields.username')">{{ model.webdav.name }}</x-description-item>
      <x-description-item :label="t('fields.account')">{{ model.webdav.account }}
        <n-button strong secondary size="small" circle type="primary" @click="copyText(model.webdav.account)">
          <template #icon>
            <n-icon>
              <copy-icon />
            </n-icon>
          </template>
        </n-button></x-description-item>
      <x-description-item :label="t('fields.token')">{{ model.webdav.token }}
        <n-button strong secondary size="small" circle type="primary" @click="copyText(model.webdav.token)">
          <template #icon>
            <n-icon>
              <copy-icon />
            </n-icon>
          </template>
        </n-button>
      </x-description-item>
      <x-description-item :label="t('fields.root_dir')">{{ model.webdav.home }}</x-description-item>
      <x-description-item :label="t('fields.status')">
        <n-tag round size="small" :type="model.webdav.is_disable == 0 ? 'primary' : 'warning'">{{
          t(model.webdav.is_disable == 0 ? 'enums.normal' : 'enums.blocked') }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        {{ model.webdav.created_at }}
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        {{ model.webdav.updated_at }}
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
import webdavApi from "@/api/nas/webdav";
import type { WebDav } from "@/api/nas/webdav";
import { copyText } from "@/utils";
import { useI18n } from 'vue-i18n'

const message = useMessage()
const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const model = reactive({
  webdav: {} as WebDav,
});
const newHandler = () => {
  router.push({ name: 'webdav_list' })
}
const editHandler = () => {
  router.push({ name: 'webdav_edit', params: { id: model.webdav.id } })
}

async function fetchData() {
  let webdav = await webdavApi.load(route.params.id as string);
  model.webdav = webdav.data as WebDav
}

watch(() => route.params.id, fetchData)

onMounted(fetchData);
</script>