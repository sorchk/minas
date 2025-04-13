<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <x-panel title="" subtitle="" :collapsed="panel !== 'dataTunnel'">
      <template #action>
        <n-button secondary strong class="toggle" size="small" @click="togglePanel('dataTunnel')">{{ panel ===
    'dataTunnel' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <n-form  inline :model="state.dataTunnel" ref="formDataTunnel"
        label-placement="left" style="padding: 4px 0 0 12px"> 
        <n-button type="primary" ></n-button>
      </n-form> 
    </x-panel>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NGrid,
  NButton,
  NSpace,
  NInput,
  NInputNumber,
  NInputGroup,
  NInputGroupLabel,
  NForm,
  NFormItem,
  NFormItemGi,
  NRadioGroup,
  NRadio,
  NSwitch,
  NAlert,
} from "naive-ui";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const state = ref({
} as any);
const panel = ref('dataTunnel')

function togglePanel(name: string) {
  if (panel.value === name) {
    panel.value = ''
  } else {
    panel.value = name
  }
}



</script>

<style scoped>
.toggle {
  width: 75px;
}
</style>