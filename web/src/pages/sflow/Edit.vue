<script setup lang="ts">
import { onMounted, ref, reactive, h } from 'vue'
import { NButton,NRadio,NRadioButton,NRadioGroup, NCard, NForm, NFormItem, NIcon, NInput, NInputNumber, NSpace, NTree, TreeOption, useMessage, NGrid, NSwitch, NFormItemGi, NGi, NTreeSelect } from 'naive-ui'
import { ArrowBackCircleOutline as BackIcon, SaveOutline as SaveIcon, FolderOutline as FolderIcon } from '@vicons/ionicons5'
import type { ProjectDirItem } from "@/api/basic/projectdir";
import { useRoute, useRouter } from 'vue-router'
import sflowApi, { SFlow } from '@/api/sflow'
import projectDirApi from "@/api/basic/projectdir";
import { t } from '@/locales'
import type { FormInst, FormRules } from 'naive-ui'
import XPageHeader from "@/components/PageHeader.vue";
import { useForm, requiredRule } from "@/utils/form";
import { deepClone } from "@/utils";

// 页面状态
const route = useRoute()
const router = useRouter()
const message = useMessage()
const form = ref<FormInst | null>(null)

// 是否为编辑模式
const isEdit = route.name === 'sflow_edit'
const logLevelName = [ '无', '错误', '警告', '信息',  '调试',  '跟踪' ]
// 表单验证规则
const rules = {
  name: requiredRule()
}

// 作业流程数据
const sflow = ref<any>({
  name: '',
  type: 'job',
  log_level: 3,  // 日志级别 
  project_dir_id: '0',
  remark: '',
})
 

// 项目目录树相关状态
const treeData = ref<TreeOption[]>([]);
const treeLoading = ref(false);
const expandedKeys = ref<string[]>([]);

// 返回列表页
const listHandler = () => {
  router.push({ name: 'sflow_list' })
}

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

// 树节点属性
const treeNodeProps = ({ option }: { option: TreeOption }) => {
  return {
    style: "cursor: pointer;"
  };
};

// 处理表单提交
const { submit, submiting } = useForm(form, () => {
  const params = deepClone(sflow.value); 
  return sflowApi.save(params);
}, () => {
  message.success(t('texts.action_success'));
  router.push({ name: 'sflow_list' });
});

// 加载作业流程详情
const loadSFlow = async () => {
  if (isEdit) {
    try {
      const res = await sflowApi.load(Number(route.params.id as string))
      if (res.code === 200 && res.data) {
        sflow.value = res.data
         
      }
    } catch (error) {
      console.error('获取作业流程详情失败:', error)
      message.error('加载作业流程详情失败')
    }
  }
}

// 组件挂载时加载数据
onMounted(async () => {
  // 检查URL查询参数是否包含项目目录ID
  const projectDirId = route.query.project_dir_id;
  
  if (!isEdit && projectDirId) {
    // 如果是新建作业流程，且URL中包含项目目录ID参数
    sflow.value.project_dir_id = parseInt(projectDirId as string);
  }
  
  await Promise.all([loadProjectDirTree(), loadSFlow()])
})
</script>

<template>
  <x-page-header :subtitle="isEdit ? ` ID:${sflow.id || ''}` : ''">
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
    <n-form style="margin:auto; width: 680px;" :model="sflow" :rules="rules" ref="form" label-placement="top">
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="24" :label="t('fields.name')" path="name">
          <n-input :placeholder="t('fields.name')" v-model:value="sflow.name" />
        </n-form-item-gi>
        
        <!-- 项目目录选择 -->
        <n-form-item-gi :span="12" :label="t('fields.projectdir')" path="project_dir_id">
          <n-tree-select
            v-model:value="sflow.project_dir_id"
            :options="treeData"
            :default-expanded-keys="expandedKeys"
            placeholder="选择项目目录"
            :node-props="treeNodeProps"
            :render-label="renderTreeLabel"
          />
        </n-form-item-gi>
        <n-form-item-gi :span="12" label="日志级别">
              <n-radio-group v-model:value="sflow.log_level" name="log_level">
                <n-radio v-for="(item,index) in logLevelName" :key="index" :value="index"
                  :label="item" />
              </n-radio-group>
          </n-form-item-gi>
        <n-form-item-gi :span="24" :label="t('fields.remark')" path="remark">
          <n-input type="textarea" :placeholder="t('fields.remark')" :autosize="{ minRows: 2 }"
            v-model:value="sflow.remark" />
        </n-form-item-gi>
        
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

<style scoped> 
</style>