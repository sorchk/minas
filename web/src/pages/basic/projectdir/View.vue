<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="listHandler">
        <template #icon>
          <n-icon>
            <arrow-back-icon />
          </n-icon>
        </template>
        {{ t('buttons.back') }}
      </n-button>
      <n-button secondary size="small" @click="editHandler">
        <template #icon>
          <n-icon>
            <edit-icon />
          </n-icon>
        </template>
        {{ t('buttons.edit') }}
      </n-button>
      <n-button v-if="projectDir.is_disable" secondary type="success" size="small" @click="enable">
        <template #icon>
          <n-icon>
            <checkmark-icon />
          </n-icon>
        </template>
        {{ t('buttons.enable') }}
      </n-button>
      <n-button v-else secondary type="warning" size="small" @click="disable">
        <template #icon>
          <n-icon>
            <stop-icon />
          </n-icon>
        </template>
        {{ t('buttons.block') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-spin :show="loading">
      <n-card :title="t('projectdir.detail')" size="small">
        <n-descriptions :column="2" label-placement="left" bordered>
          <n-descriptions-item :label="t('fields.id')">
            {{ projectDir.id }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.name')">
            {{ projectDir.name }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.parent')">
            {{ parentName }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.status')">
            <n-tag :type="projectDir.is_disable ? 'warning' : 'success'">
              {{ projectDir.is_disable ? t('enums.blocked') : t('enums.normal') }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.has_child')">
            {{ projectDir.has_child ? t('enums.yes') : t('enums.no') }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.remark')" :span="2">
            {{ projectDir.remark || '-' }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.created_at')" :span="2">
            {{ projectDir.created_at }}
          </n-descriptions-item>
          <n-descriptions-item :label="t('fields.updated_at')" :span="2">
            {{ projectDir.updated_at }}
          </n-descriptions-item>
        </n-descriptions>
      </n-card>
      
      <n-card v-if="hasChildren" class="mt-4" :title="t('projectdir.children')" size="small">
        <n-data-table
          :row-key="(row: any) => row.id"
          size="small"
          :columns="childrenColumns"
          :data="childrenList"
          :pagination="{ pageSize: 10 }"
          scroll-x="max-content"
        />
      </n-card>
    </n-spin>
  </n-space>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import {
  NSpace,
  NButton,
  NIcon,
  NCard,
  NSpin,
  NDescriptions,
  NDescriptionsItem,
  NTag,
  NDataTable,
  useMessage,
} from "naive-ui";
import {
  ArrowBackOutline as ArrowBackIcon,
  CreateOutline as EditIcon,
  CheckmarkOutline as CheckmarkIcon,
  StopOutline as StopIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import { useRoute, useRouter } from "vue-router";
import projectDirApi from "@/api/basic/projectdir";
import type { ProjectDirItem } from "@/api/basic/projectdir";
import { renderLink, renderTag } from "@/utils/render";
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const message = useMessage();

// 页面状态
const loading = ref(true);
const projectDir = ref<ProjectDirItem>({} as ProjectDirItem);
const childrenList = ref<ProjectDirItem[]>([]);
const parentInfo = ref<ProjectDirItem | null>(null);

// 计算属性
const hasChildren = computed(() => childrenList.value && childrenList.value.length > 0);
const parentName = computed(() => parentInfo.value ? parentInfo.value.name : '-');

// 导航处理函数
const listHandler = () => {
  router.push({ name: 'projectdir_list' });
};

const editHandler = () => {
  router.push({ name: 'projectdir_edit', params: { id: projectDir.value.id } });
};

// 子目录列表列定义
const childrenColumns = [
  {
    title: t('fields.id'),
    key: "id",
    render: (row: ProjectDirItem) => renderLink({ name: 'projectdir_detail', params: { id: row.id } }, row.id),
  },
  {
    title: t('fields.name'),
    key: "name",
  },
  {
    title: t('fields.status'),
    key: "is_disable",
    render: (row: ProjectDirItem) => renderTag(
      row.is_disable ? t('enums.blocked') : t('enums.normal'),
      row.is_disable ? "warning" : "success"
    ),
  },
  {
    title: t('fields.remark'),
    key: "remark",
  },
];

// 启用目录
async function enable() {
  try {
    await projectDirApi.enable(projectDir.value.id);
    message.success(t('texts.action_success'));
    fetchData(); // 重新加载数据
  } catch (error) {
    message.error(t('errors.operation_failed'));
  }
}

// 禁用目录
async function disable() {
  try {
    await projectDirApi.disable(projectDir.value.id);
    message.success(t('texts.action_success'));
    fetchData(); // 重新加载数据
  } catch (error) {
    message.error(t('errors.operation_failed'));
  }
}

// 加载数据
async function fetchData() {
  loading.value = true;
  try {
    const id = parseInt(route.params.id as string);
    
    // 加载当前目录详情
    const result = await projectDirApi.find(id);
    if (result.data) {
      projectDir.value = result.data;
      
      // 如果有父目录ID，加载父目录信息
      if (projectDir.value.parent_id) {
        try {
          const parentResult = await projectDirApi.find(projectDir.value.parent_id);
          parentInfo.value = parentResult.data;
        } catch (error) {
          console.error("加载父目录信息失败", error);
        }
      }
      
      // 加载子目录列表
      if (projectDir.value.has_child) {
        try {
          const childrenResult = await projectDirApi.getChildren(id);
          childrenList.value = childrenResult.data || [];
        } catch (error) {
          console.error("加载子目录信息失败", error);
        }
      }
    }
  } catch (error) {
    console.error("加载项目目录信息失败", error);
    message.error(t('errors.load_failed'));
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.page-body {
  padding: 10px;
}
.mt-4 {
  margin-top: 16px;
}
</style>