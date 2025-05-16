<script setup lang="ts">
import { h, onMounted, ref, watch, reactive } from 'vue'
import { 
  NButton, 
  NCard, 
  NDataTable, 
  NEllipsis, 
  NInput, 
  NPopconfirm, 
  NSpace, 
  NTree, 
  NIcon,
  NGrid,
  NGridItem,
  NSpin,
  useMessage,
  TreeOption
} from 'naive-ui'
import sflowApi, { SFlow, statusMapping } from '@/api/sflow'
import { 
  SearchOutline, 
  AddOutline, 
  PlayCircleOutline, 
  DocumentTextOutline, 
  CreateOutline, 
  TrashOutline, 
  EyeOutline,
  RefreshOutline
} from '@vicons/ionicons5'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import type { ProjectDirItem } from "@/api/basic/projectdir";
import projectDirApi from '@/api/basic/projectdir'
import { t } from '@/locales'
import { xxtea } from "@/utils/xxtea";
import { renderButtons, renderLink, renderTag, renderTime } from "@/utils/render";
import { getStatusTag } from '@/utils'
import XPageHeader from "@/components/PageHeader.vue"
import { useDataTable } from "@/utils/data-table"
const { state, pagination, fetchData, changePageSize } = useDataTable(sflowApi.search, () => {
  return { ...args, filters: route.query.filters }
})
// 页面状态
const route = useRoute()
const router = useRouter()
const message = useMessage()

// 搜索相关状态
const args = reactive({
  name: "",
  type: "job",
  project_dir_id: undefined as number | undefined
})

// 项目目录树相关状态
const treeData = ref<TreeOption[]>([]);
const treeLoading = ref(false);
const expandedKeys = ref<string[]>([]);
const selectedKeys = ref<string[]>([]);
const logLevelName = [ '无', '错误', '警告', '信息',  '调试',  '跟踪' ]
const logLevelType = [ 'default', 'error', 'warning', 'info', 'success', 'success' ]
// 表格列定义
const columns = [
{
    title: t('fields.id'),
    key: "id",
    render: (row: SFlow) => renderLink({ name: 'sflow_detail', params: { id: row.id } }, row.id),
  },
  {
    title: '名称',
    key: 'name',
  },
  {
    // NONE(0), ERROR(1), WARN(2), INFO(3), DEBUG(4), TRACE(5);
    title: '日志级别',
    key: 'log_level',
    render(row: SFlow) {
      const logl = logLevelType[row.log_level] as ("info" | "error" | "default" | "success" | "warning" | undefined)
      return renderTag(logLevelName[row.log_level], logl, "small")
    }
  },
  {
    title: '最近状态',
    key: 'last_status',
    render(row: SFlow) {
      return getStatusTag(statusMapping, row.last_status, "未执行", false)
    }
  },
  {
    title: '上次执行时间',
    key: 'last_run_time',
    render(row: SFlow) {
      return renderTime(row.last_run_time)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 300,
    render(row: SFlow, index: number) {
      return renderButtons([
        { type: 'success', text: t('buttons.run_now'), action: () => exec(row, index), prompt: t('prompts.run_now') },
        { type: 'info', text: t('buttons.run_logs'), action: () => router.push({ name: 'sflowlog_list', params: { id: row.id } }) },
        row.is_disable ?
          { type: 'success', text: t('buttons.enable'), action: () => enable(row), prompt: t('prompts.enable'), } :
          { type: 'warning', text: t('buttons.block'), action: () => disable(row), prompt: t('prompts.block'), },
          { type: 'success', text: t('buttons.design'), action: () => router.push({ name: 'sflow_design', params: { id: row.id, type: row.type } }) },
        { type: 'warning', text: t('buttons.edit'), action: () => router.push({ name: 'sflow_edit', params: { id: row.id } }) },
        { type: 'error', text: t('buttons.delete'), action: () => remove(row, index), prompt: t('prompts.delete') },
      ])
    },
  }
]
 
// 项目目录树加载
const loadProjectDirTree = async () => {
  treeLoading.value = true;
  try {
    // 获取项目目录树
    const res = await projectDirApi.getTree(); 
    if (res.data && Array.isArray(res.data)) {
      treeData.value = convertToTreeOptions(res.data);
      
      // 如果有数据，默认展开第一级
      if (treeData.value.length > 0) {
        expandedKeys.value = [treeData.value[0].key as string];
      }
    }
  } catch (error) {
    console.error("加载项目目录树失败", error);
  } finally {
    treeLoading.value = false;
  }
};

// 转换项目目录数据为树形选项
const convertToTreeOptions = (dirs: ProjectDirItem[]): TreeOption[] => {
  return dirs.map(dir => ({
    key: dir.id.toString(),
    label: dir.name,
    children: dir.children && dir.children.length > 0 
      ? convertToTreeOptions(dir.children) 
      : undefined
  }));
};

// 项目目录树节点选择处理函数
const handleSelectNode = (keys: string[]) => {
  selectedKeys.value = keys;
  // 根据选中的目录节点过滤任务列表
  if (keys.length > 0) {
    args.project_dir_id = parseInt(keys[0]);
  } else {
    args.project_dir_id = undefined;
  }
  fetchData(1);
};

// 添加新作业流程
const addHandler = () => {
  // 如果已选择了项目目录，则将目录ID传递给新建页面
  if (selectedKeys.value.length > 0) {
    router.push({ 
      name: 'sflow_new',
      query: { project_dir_id: selectedKeys.value[0] }
    });
  } else {
    router.push({ name: 'sflow_new' });
  }
}
async function enable(u: SFlow) {
  await sflowApi.enable(u.id);
  fetchData()
}
async function disable(u: SFlow) {
  await sflowApi.disable(u.id);
  fetchData()
}
async function remove(u: SFlow, index: number) {
  const data = await sflowApi.delete(u.id);
  if (data.code === 200) {
    state.data.splice(index, 1)
    message.info(t('texts.action_success'));
  } else {
    message.info(t('texts.action_error'));
  }
} 

// 执行作业流程
async function exec(row: SFlow, index: number) {
  try {
    const res = await sflowApi.exec(row.id)
    if (res.code === 200) {
      message.success(t('texts.action_success'))
      // 延迟一秒后刷新数据，以便看到最新状态
      setTimeout(() => {
        fetchData()
      }, 1000)
    } else {
      message.error(t('texts.action_error'))
    }
  } catch (error) {
    console.error('执行作业流程失败:', error)
    message.error(t('texts.action_error'))
  }
}

// 查看作业流程运行日志
const toLog = (row: SFlow) => {
  router.push({
    name: 'sflowlog_list',
    params: { id: row.id }
  })
}

// 监听路由查询参数变化
watch(() => route.query.filter, (newValue: any) => {
  fetchData()
})

// 组件挂载时加载项目目录树和初始数据
onMounted(() => {
  loadProjectDirTree()
  fetchData()
})
</script>

<template>
  <x-page-header>
    <template #action>
      <n-button secondary size="small" @click="addHandler">
        <template #icon>
          <n-icon>
            <AddOutline />
          </n-icon>
        </template>
        {{ t('buttons.new') }}
      </n-button>
      <n-button secondary size="small" @click="fetchData()">
        <template #icon>
          <n-icon>
            <RefreshOutline />
          </n-icon>
        </template>
        {{ t('buttons.refresh') }}
      </n-button>
    </template>
  </x-page-header>
  <n-space class="page-body" vertical :size="12">
    <n-grid :cols="24" :x-gap="12">
      <!-- 左侧项目目录树 -->
      <n-grid-item :span="6">
        <n-card :bordered="false" size="small" title="项目目录">
          <n-spin :show="treeLoading">
            <n-tree
              block-line
              :data="treeData"
              :default-expanded-keys="expandedKeys"
              :selected-keys="selectedKeys"
              :on-update:selected-keys="handleSelectNode"
            />
          </n-spin>
        </n-card>
      </n-grid-item>
      <!-- 右侧作业流程列表 -->
      <n-grid-item :span="18">
        <n-space :size="12">
          <n-input size="small" v-model:value="args.name" :placeholder="t('fields.name')" clearable @keydown.enter="fetchData(1)" />
          <n-button size="small" type="primary" @click="() => fetchData(1)">{{ t('buttons.search') }}</n-button>
        </n-space>

        <n-data-table remote :row-key="(row: any) => row.id" size="small" :columns="columns" :data="state.data"
          :pagination="pagination" :loading="state.loading" @update:page="fetchData" @update-page-size="changePageSize" scroll-x="max-content"  />
     
      </n-grid-item>
    </n-grid>
  </n-space>
</template>
