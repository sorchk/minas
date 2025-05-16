<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue'
import { NTag, NSpace, NButton, useMessage, NIcon } from 'naive-ui'
import { ArrowBackCircleOutline as BackIcon, CreateOutline as EditIcon, PlayCircleOutline as RunIcon, DocumentTextOutline as LogIcon, CopyOutline as CopyIcon } from '@vicons/ionicons5'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import sflowApi, { SFlow } from '@/api/sflow'
import { t } from '@/locales'
import { getStatusTag } from '@/utils'
import { statusMapping } from '@/api/sflow'
import XPageHeader from "@/components/PageHeader.vue"
import { XDescription, XDescriptionItem } from "@/components/description"
import { copyText } from "@/utils"

// 页面状态
const route = useRoute()
const router = useRouter()
const message = useMessage()
const { t: tl } = useI18n()

// 作业流程数据
const model = reactive({
  sflow: {} as SFlow
})

// 返回列表页
const backHandler = () => {
  router.push({ name: 'sflow_list' })
}

// 编辑作业流程
const editHandler = () => {
  router.push({ name: 'sflow_edit', params: { id: model.sflow.id } })
}

// 查看作业流程日志
const logHandler = () => {
  router.push({ name: 'sflowlog_list', params: { id: model.sflow.id } })
}

// 执行作业流程
async function execHandler() {
  try {
    const res = await sflowApi.exec(model.sflow.id)
    if (res.code === 200) {
      message.success('执行成功')
      // 重新加载作业流程信息
      await fetchData()
    } else {
      message.error('执行失败')
    }
  } catch (error) {
    console.error('执行作业流程失败:', error)
    message.error('执行失败')
  }
}

// 获取作业流程详情
async function fetchData() {
  try {
    const res = await sflowApi.load(route.params.id as string)
    if (res.code === 200 && res.data) {
      model.sflow = res.data as SFlow
    } else {
      message.error('加载作业流程详情失败')
    }
  } catch (error) {
    console.error('获取作业流程详情失败:', error)
    message.error('加载作业流程详情失败')
  }
}

// 监听路由参数变化，重新加载数据
watch(() => route.params.id, fetchData)

// 组件挂载时加载数据
onMounted(fetchData)
</script>

<template>
  <x-page-header :subtitle="model.sflow.name || '未命名作业流程'">
    <template #action>
      <n-button secondary size="small" @click="backHandler">
        <template #icon>
          <n-icon>
            <BackIcon />
          </n-icon>
        </template>
        {{ tl('buttons.return') }}
      </n-button>
      <n-button secondary size="small" @click="editHandler">
        <template #icon>
          <n-icon>
            <EditIcon />
          </n-icon>
        </template>
        {{ tl('buttons.edit') }}
      </n-button>
      <n-button type="info" size="small" @click="execHandler">
        <template #icon>
          <n-icon>
            <RunIcon />
          </n-icon>
        </template>
        执行
      </n-button>
      <n-button type="primary" secondary size="small" @click="logHandler">
        <template #icon>
          <n-icon>
            <LogIcon />
          </n-icon>
        </template>
        运行记录
      </n-button>
    </template>
  </x-page-header>
  
  <n-space class="page-body" vertical :size="16">
    <x-description cols="1 640:1" label-position="left" label-align="right" :label-width="100">
      <x-description-item :label="tl('fields.name')">{{ model.sflow.name || '未命名' }}</x-description-item>
      <x-description-item label="最近状态">
        <div v-html="getStatusTag(statusMapping, model.sflow.last_status, '未执行')"></div>
      </x-description-item>
      <x-description-item label="上次执行时间">
        {{ model.sflow.last_run_time || '-' }}
      </x-description-item>
      <x-description-item label="创建人">
        {{ model.sflow.created_by_name }}
      </x-description-item>
      <x-description-item label="创建时间">
        {{ model.sflow.created_at }}
      </x-description-item>
      <x-description-item label="更新人">
        {{ model.sflow.updated_by_name }}
      </x-description-item>
      <x-description-item label="更新时间">
        {{ model.sflow.updated_at }}
      </x-description-item>
      <x-description-item label="备注">
        {{ model.sflow.remark || '-' }}
      </x-description-item>
      <x-description-item label="作业流程脚本">
        <div v-if="model.sflow.content" style="max-height: 300px; overflow-y: auto; white-space: pre-wrap;">
          {{ JSON.parse(model.sflow.content).script || '无脚本内容' }}
          <n-button v-if="model.sflow.content" strong secondary size="small" circle type="primary" 
                   @click="copyText(JSON.parse(model.sflow.content).script || '')">
            <template #icon>
              <n-icon>
                <CopyIcon />
              </n-icon>
            </template>
          </n-button>
        </div>
        <div v-else>无脚本内容</div>
      </x-description-item>
    </x-description>
  </n-space>
</template>

<style scoped>
.page-body {
  padding: 0 16px;
}
</style>