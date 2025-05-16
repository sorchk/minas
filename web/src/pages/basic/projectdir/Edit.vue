<template>
  <x-page-header :subtitle="isEdit ? ` ID:${projectDir.id || ''}` : ''">
    <template #action>
      <n-button secondary size="small" @click="listHandler">
        <template #icon>
          <n-icon>
            <arrow-back-icon />
          </n-icon>
        </template>
        {{ t('buttons.back') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-form ref="form" :model="projectDir" :rules="rules" style="margin:auto; width: 680px;" label-placement="top">
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="24" :label="t('fields.name')" path="name">
          <n-input v-model:value="projectDir.name" :placeholder="t('fields.name')" />
        </n-form-item-gi>

        <n-form-item-gi :span="24" :label="t('fields.parent')" path="parent_id">
          <n-tree-select v-model:value="projectDir.parent_id" :options="treeData"
            :placeholder="t('fields.parent_placeholder')" clearable :render-label="renderTreeLabel" />
        </n-form-item-gi>

        <n-form-item-gi :span="24" :label="t('fields.remark')" path="remark">
          <n-input v-model:value="projectDir.remark" type="textarea" :autosize="{ minRows: 2 }"
            :placeholder="t('fields.remark')" />
        </n-form-item-gi>

        <n-gi :span="2">
          <n-button type="primary" @click="submit" :loading="submiting">
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
import { reactive, ref, computed, onMounted, h } from "vue";
import {
  NSpace,
  NButton,
  NIcon,
  NCard,
  NForm,
  NFormItemGi,
  NInput,
  NGrid,
  NGi,
  NTreeSelect,
  SelectOption,
  TreeSelectOption,
  TreeOption,
  useMessage,
} from "naive-ui";
import {
  ArrowBackCircleOutline as ArrowBackIcon,
  SaveOutline as SaveIcon,
  Folder as FolderIcon
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import { useRoute, useRouter } from "vue-router";
import projectDirApi from "@/api/basic/projectdir";
import type { ProjectDirItem } from "@/api/basic/projectdir";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n';
import { deepClone } from "@/utils";

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const message = useMessage();

const listHandler = () => {
  router.push({ name: 'projectdir_list' });
};

// 是否编辑模式
const isEdit = computed(() => !!route.params.id);

// 表单数据
const projectDir = ref<any>({
  is_disable: 0,
  parent_id: 0,
  name: "",
  remark: "",
});

// 表单验证规则
const rules = {
  name: requiredRule()
};

// 表单ref
const form = ref();

// 树形选择器数据
const treeData = ref<TreeOption[]>([]);

// 加载目录树
const loadProjectDirTree = async () => {
  try {
    const res = await projectDirApi.getTree();
    if (res.data && Array.isArray(res.data)) {
      // 如果是编辑模式，需要过滤掉当前节点及其子节点，防止自引用
      const currentId = isEdit.value ? parseInt(route.params.id as string) : -1;
      treeData.value = convertToTreeOptions(
        res.data,
        currentId
      );
    }
  } catch (error) {
    console.error("加载项目目录树失败", error);
  }
};

// 转换项目目录数据为树形选项，排除自身及子节点
const convertToTreeOptions = (dirs: ProjectDirItem[], excludeId: number = -1): TreeOption[] => {
  return dirs
    .filter(dir => dir.id !== excludeId) // 过滤掉当前编辑的节点
    .map(dir => ({
      key: dir.id,
      value: dir.id,
      label: dir.name,
      children: dir.children && dir.children.length > 0
        ? convertToTreeOptions(dir.children, excludeId)
        : undefined
    }));
};

// 树节点渲染
const renderTreeLabel = (info: { option: TreeOption }) => {
  return h(
    "div",
    {
      style: "display: flex; align-items: center;"
    },
    [
      h(
        NIcon,
        {
          style: "margin-right: 4px"
        },
        {
          default: () => h(FolderIcon)
        }
      ),
      info.option.label as string
    ]
  );
};

// 提交表单
const { submit, submiting } = useForm(form, () => {
  const params = deepClone(projectDir.value);
  return projectDirApi.save(params as ProjectDirItem);
}, () => {
  message.success(t('texts.action_success'));
  router.push({ name: 'projectdir_list' });
});

// 加载数据
const fetchData = async () => {
  // 加载目录树
  await loadProjectDirTree();

  if (isEdit.value) {
    // 编辑模式下加载当前项目目录数据
    const id = parseInt(route.params.id as string);
    const result = await projectDirApi.find(id);
    projectDir.value = result.data;
  } else {
    // 新建模式下，检查是否有parent_id参数
    const parentId = route.query.parent_id;
    if (parentId) {
      projectDir.value.parent_id = parseInt(parentId as string);
    }
  }
};

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.page-body {
  padding: 10px;
}
</style>