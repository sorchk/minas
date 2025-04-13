<template>
  <x-page-header :subtitle="` ID:${externalNas.id || ''}`">
    <template #action>
      <n-button secondary size="small" @click="listHandler">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>
        {{ t('buttons.return') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-form style="margin:auto; width: 480px;" :model="externalNas" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:1" :x-gap="24">
        <n-form-item-gi :label="t('fields.name')" path="name">
          <n-input :placeholder="t('fields.name')" v-model:value="externalNas.name" />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.rc_name')" path="rc_name">
          <n-input v-if="!externalNas.id" :placeholder="t('fields.rc_name')" v-model:value="externalNas.rc_name" />
          <n-tag type="info" round v-else> {{ externalNas.rc_name }} </n-tag>
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.type')" path="type">
          <n-select :disabled="!!externalNas.id" v-model:value="externalNas.type" filterable :options="providers"
            @update:value="onSelectType" clearable />
        </n-form-item-gi>
        <n-form-item-gi :label="t('fields.remark')" path="remark">
          <n-input :placeholder="t('fields.remark')" v-model:value="externalNas.remark" />
        </n-form-item-gi>
        <template v-for="(item, i) in selectTypeObj.Options">
          <template v-if="!item.Advanced">
            <n-form-item-gi :label="item.Name" :path="`config.${item.Name}`" :rule="(item.Required && item.Type != 'bool') ? {
              required: true,
              message: t('tips.cannot_be_empty'),
              trigger: ['input', 'blur'],
            } : {}">
              <n-input v-if="item.Type == 'string' && !item.IsPassword && !item.Examples" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-select v-if="item.Type == 'string' && !item.IsPassword && !!item.Examples" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" label-field="Help" value-field="Value"
                :options="item.Examples" clearable>
              </n-select>
              <n-input v-if="item.Type == 'string' && item.IsPassword" type="password" show-password-on="click"
                :placeholder="item.Help" v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr"
                clearable />

              <n-input v-if="item.Type == 'int'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-input v-else-if="item.Type == 'Duration'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-input v-else-if="item.Type == 'SizeSuffix'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-switch v-else-if="item.Type == 'bool'" :checked-value="true" :unchecked-value="false"
                v-model:value="externalNas.config[item.Name]" :default-value="item.Default" />
              <n-input v-else-if="item.Type == 'SpaceSepList'" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
            </n-form-item-gi>
          </template>
        </template>
        <n-form-item-gi :label="t('fields.adv_options')">
          <n-switch v-model:value="externalNas.is_adv" :checked-value="1" :unchecked-value="0">
            <template #checked>
              {{ t('tips.enable_adv_options') }}
            </template>
            <template #unchecked>
              {{ t('tips.disable_adv_options') }}
            </template>
          </n-switch>
        </n-form-item-gi>
        <template v-for="(item, i) in selectTypeObj.Options">
          <template v-if="item.Advanced && externalNas.is_adv == 1">
            <n-form-item-gi :label="item.Name" :path="`config.${item.Name}`" :rule="{
              required: item.Required,
              message: t('tips.cannot_be_empty'),
              trigger: ['input', 'blur'],
            }">

              <n-input v-if="item.Type == 'string' && !item.IsPassword && !item.Examples" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" :default-value="item.Default" clearable />
              <n-select v-if="item.Type == 'string' && !item.IsPassword && !!item.Examples" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" label-field="Help" value-field="Value"
                :options="item.Examples" clearable>
              </n-select>
              <n-input v-if="item.Type == 'string' && item.IsPassword" type="password" show-password-on="click"
                :placeholder="item.Help" v-model:value="externalNas.config[item.Name]" :default-value="item.Default"
                clearable />
              <n-input v-if="item.Type == 'int'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-input v-else-if="item.Type == 'Duration'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-input v-else-if="item.Type == 'SizeSuffix'" :allow-input="onlyAllowNumber"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
              <n-switch v-else-if="item.Type == 'bool'" :checked-value="true" :unchecked-value="false"
                v-model:value="externalNas.config[item.Name]" :default-value="item.Default" />
              <n-input v-else-if="item.Type == 'SpaceSepList'" :placeholder="item.Help"
                v-model:value="externalNas.config[item.Name]" :default-value="item.DefaultStr" clearable />
            </n-form-item-gi>
          </template>
        </template>
        <n-gi :span="2">
          <n-button :disabled="submiting" :loading="submiting" @click.prevent="submit" type="primary">
            <template #icon>
              <n-icon>
                <save-icon />
              </n-icon>
            </template>
            {{ t('buttons.save') }}
          </n-button>
        </n-gi>
      </n-grid>
    </n-form>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, VNodeChild, h } from "vue";
import {
  NTag,
  NTree,
  NDrawer,
  NDrawerContent,
  NButton,
  NSpace,
  NInput,
  NInputGroup,
  NIcon,
  NForm,
  NGrid,
  NGi,
  NTreeSelect,
  NFormItemGi,
  NSwitch,
  NCheckboxGroup,
  NCheckbox,
  NRadioGroup,
  NRadio,
  NDynamicInput,
  NTooltip,
  useMessage,
  TreeSelectOption,
  TreeSelectOverrideNodeClickBehavior,
  TreeOption,
  NSelect,
  NInputNumber,
  SelectOption,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
  CopyOutline as CopyIcon,
  Folder as FolderIcon,
  FolderOpenOutline as FolderOpenIcon,

} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import { useRoute, useRouter } from "vue-router";
import externalNasApi from "@/api/nas/external";
import type { ExternalNas } from "@/api/nas/external";
import { useForm, emailRule, requiredRule, customRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'
import { useClipboard } from '@vueuse/core'
import { deepClone, guid } from "@/utils";
import { log } from "console";

const { t } = useI18n()
const route = useRoute();
const router = useRouter();
const message = useMessage()
const listHandler = () => {
  router.push({ name: 'externalNas_list' })
}
const providers = ref<SelectOption[]>([])
const providersMap = ref<any>({})
const selectTypeObj = ref<any>({})
const externalNas = ref({ is_adv: 0, config: {} as any } as any)
const rules: any = {
  name: requiredRule(),
  type: requiredRule(),
  rc_name: requiredRule(),
};
const form = ref();
const onSelectType = (value: string) => {
  console.log(value)
  console.log(providersMap.value[value])
  selectTypeObj.value = providersMap.value[value]

}
const { submit, submiting } = useForm(form, () => {
  const params = deepClone(externalNas.value);
  params.config = JSON.stringify(params.config);
  return externalNasApi.save(params)
}, () => {
  message.info(t('texts.action_success'));
  router.push({ name: 'externalNas_list' })
})
const onlyAllowNumber = (value: string) => {
  return !value || /^\d+$/.test(value)
}
async function fetchData() {
  const id = route.params.id as string || ''
  const result = await externalNasApi.rcloneApi("config/providers", {}) as any
  if (result.code == 200) {
    providers.value = [];
    providersMap.value = {};
    (result.data['providers'] as Array<any>).forEach(item => {
      providers.value.push(
        { "label": item.Description, "value": item.Name }
      )
      providersMap.value[item.Name] = item;
    })
  } else {
    message.error(t('tips.get_storage_type_failed') + result.msg)
  }
  if (id) {
    const r = await externalNasApi.load(id);
    let nas = r.data as ExternalNas;
    if (nas.config) {
      nas.config = JSON.parse(nas.config)
    } else {
      nas.config = {}
    }

    externalNas.value = r.data as ExternalNas;
    selectTypeObj.value = providersMap.value[externalNas.value.type]
  }
}
onMounted(fetchData);
</script>
<style scoped>
.dir-input {
  cursor: pointer !important;
}

:deep(.dir-input .n-input__input-el) {
  cursor: pointer !important;
}
</style>
