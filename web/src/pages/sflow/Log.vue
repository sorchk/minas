<script setup lang="ts">
import { h, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDescriptions, NDescriptionsItem, NEllipsis, NMenu, NSpace, useMessage } from 'naive-ui'

import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
  CopyOutline as CopyIcon,
  Folder as FolderIcon,
  FolderOpenOutline as FolderOpenIcon,
  PlayCircleOutline as PlayIcon,

} from "@vicons/ionicons5";
import { useRoute, useRouter } from 'vue-router'
import sflowLogApi, { SFlowLog, statusMapping } from '@/api/sflow/logs'
import sflowApi from '@/api/sflow'
import { t } from '@/locales'
import { xxtea } from "@/utils/xxtea";

// 页面状态
const route = useRoute()
const router = useRouter()
const message = useMessage()

// 日志数据
const state = reactive({
  data: {} as SFlowLog,
  selectId: '',
})

// 日志列表菜单选项
const menuOptions = ref<Array<{ label: () => ReturnType<typeof h>; key: string; status: any }>>([])

// 自动刷新定时器
const intervalId = ref<any>(null)
const intervalId2 = ref<any>(null)

// 获取状态样式
const getStatusClass = (status: any) => {
  if (status === -1) return 'status-error'
  if (status === 0) return 'status-processing'
  if (status === 1) return 'status-success'
  if (status === -2) return 'status-warning'
  return ''
}

// 返回列表页
const listHandler = () => {
  router.push({ name: 'sflow_list' })
}

// 执行作业流程
async function exec() {
  try {
    const res = await sflowApi.exec(route.params.id as string)
    if (res.code === 200) {
      message.success('执行成功')
      await fetchData(true)
    } else {
      message.error('执行失败')
    }
  } catch (error) {
    console.error('执行作业流程失败:', error)
    message.error('执行失败')
  }
}

// 加载日志列表
const fetchData = async (selectFirst = false) => {
  try {
    const args = { page: 1, size: 50 } as any
    // 添加需要查询的字段
    args.columns = xxtea.encryptAuto(JSON.stringify("id,sflow_id,status,start_time,end_time".split(",")), "columns")
    // 按时间倒序排序
    args.sorts = xxtea.encryptAuto(JSON.stringify(["-start_time"]), "sorts")
    // 过滤条件：仅显示当前作业流程的日志
    args.filters = xxtea.encryptAuto(JSON.stringify([
      { "Column": "sflow_id", "Operator": "=?", "Value": route.params.id as string }
    ]), "filters")
    
    // 获取日志列表
    const res = await sflowLogApi.search(args)
    if (res.code === 200 && res.data) {
      // 清空当前菜单选项
      menuOptions.value = []
      
      // 生成菜单选项
      res.data.forEach((item: SFlowLog, index: number) => {
        // 如果是第一条日志或者指定选择第一条，则自动选中
        if (index === 0 && (selectFirst || !state.selectId)) {
          state.selectId = item.id
          // 加载选中日志的详情
          Load()
        }
        
        // 添加菜单选项
        menuOptions.value.push({
          label: () => h(NEllipsis, null, { default: () => item.start_time }),
          key: item.id,
          status: item.status
        })
      })
    }
  } catch (error) {
    console.error('获取日志列表失败:', error)
  }
}

// 加载选中日志的详情
const Load = async () => {
  try {
    if (!state.selectId) return
    
    const res = await sflowLogApi.load(state.selectId)
    if (res.code === 200) {
      state.data = res.data as SFlowLog
    } else {
      state.data = {} as SFlowLog
    }
  } catch (error) {
    console.error('加载日志详情失败:', error)
  }
}

// 组件挂载时初始化
onMounted(() => {
  // 加载初始数据
  fetchData()
  
  // 设置自动刷新定时器
  intervalId.value = setInterval(Load, 3000) // 每3秒刷新当前日志详情
  intervalId2.value = setInterval(fetchData, 5000) // 每5秒刷新日志列表
})

// 组件卸载前清理定时器
onBeforeUnmount(() => {
  clearInterval(intervalId.value)
  clearInterval(intervalId2.value)
})
</script>

<template>
  <div class="sflow-log">
    <n-card title="作业流程执行日志" size="small" class="mt10">
      <template #header-extra>
        <n-space>
          <n-button @click="listHandler" secondary circle>
            <template #icon>
              <ArrowBackCircleOutline />
            </template>
          </n-button>
          <n-button @click="exec" type="info">
            <template #icon>
              <PlayCircleOutline />
            </template>
            执行作业流程
          </n-button>
        </n-space>
      </template>

      <div class="log-container">
        <!-- 左侧日志列表 -->
        <div class="log-list">
          <n-card title="历史记录" size="small">
            <n-menu
              v-model:value="state.selectId"
              :options="menuOptions"
              @update:value="Load"
              :render-label="(option) => {
                return h(
                  'div',
                  {
                    class: ['menu-item', getStatusClass(option.status)]
                  },
                  [
                    option.label()
                  ]
                )
              }"
            />
          </n-card>
        </div>

        <!-- 右侧日志详情 -->
        <div class="log-detail">
          <n-card title="日志详情" size="small" v-if="state.data.id">
            <n-descriptions bordered>
              <n-descriptions-item label="开始时间">
                {{ state.data.start_time || '-' }}
              </n-descriptions-item>
              <n-descriptions-item label="结束时间">
                {{ state.data.end_time || '执行中...' }}
              </n-descriptions-item>
              <n-descriptions-item label="状态">
                <span :class="getStatusClass(state.data.status)">
                  {{ statusMapping[state.data.status]?.info || '未知状态' }}
                </span>
              </n-descriptions-item>
              <n-descriptions-item label="日志内容" >
                <div class="log-content">{{ state.data.log_text || '无日志内容' }}</div>
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
          <n-card v-else title="日志详情" size="small">
            <div class="empty-log">
              请选择左侧的日志记录查看详情
            </div>
          </n-card>
        </div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.sflow-log {
  padding: 0 10px;
}

.log-container {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.log-list {
  width: 220px;
  flex-shrink: 0;
}

.log-detail {
  flex-grow: 1;
}

.log-content {
  white-space: pre-wrap;
  font-family: monospace;
  font-size: 14px;
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  max-height: calc(100vh - 320px);
  overflow: auto;
}

.empty-log {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}

.menu-item {
  padding: 8px 0;
}

.status-error {
  color: #f5222d;
}

.status-success {
  color: #52c41a;
}

.status-processing {
  color: #1890ff;
}

.status-warning {
  color: #faad14;
}
</style>