<template>
  <x-page-header :subtitle="model.externalNas.name">
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
      <x-description-item :label="t('fields.id')">{{ model.externalNas.id }}</x-description-item>
      <x-description-item :label="t('fields.username')">{{ model.externalNas.name }}</x-description-item>
      <x-description-item :label="t('fields.type')">{{ model.externalNas.type }} </x-description-item>
      <x-description-item :label="t('fields.rc_name')">{{ model.externalNas.rc_name }}
        <n-button strong secondary size="small" circle type="primary" @click="copyText(model.externalNas.rc_name)">
          <template #icon>
            <n-icon>
              <copy-icon />
            </n-icon>
          </template>
        </n-button>
      </x-description-item>
      <x-description-item :label="t('fields.status')">
        <n-tag round size="small" :type="model.externalNas.is_disable == 0 ? 'primary' : 'warning'">{{
          t(model.externalNas.is_disable == 0 ? 'enums.normal' : 'enums.need_sync') }}</n-tag>
      </x-description-item>
      <x-description-item :label="t('fields.created_at')">
        {{ model.externalNas.created_at }}
      </x-description-item>
      <x-description-item :label="t('fields.updated_at')">
        {{ model.externalNas.updated_at }}
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
import externalNasApi from "@/api/nas/external";
import type { ExternalNas } from "@/api/nas/external";
import { copyText } from "@/utils";
import { useI18n } from 'vue-i18n'

const message = useMessage()
const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const model = reactive({
  externalNas: {} as ExternalNas,
});
const newHandler = () => {
  router.push({ name: 'externalNas_list' })
}
const editHandler = () => {
  router.push({ name: 'externalNas_edit', params: { id: model.externalNas.id } })
}

async function fetchData() {
  let externalNas = await externalNasApi.load(route.params.id as string);
  model.externalNas = externalNas.data as ExternalNas
}

watch(() => route.params.id, fetchData)

onMounted(fetchData);
</script>