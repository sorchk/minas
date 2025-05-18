<template>
  <s-dag-designer :title="state.title" :initNodes="state.initNodes" :content="state.content" :groups="state.groups"
    :edgeFields="edgeFields" :result="state.result" openGroup="base" @exec="handleTest" @save="handleSave"  @back="listHandler" @change="handleContentChange"/>
</template>

<script lang="ts" setup>
 
import { useRoute, useRouter } from 'vue-router'
import { onMounted, reactive, getCurrentInstance, onBeforeUnmount, ref } from "vue";
// 导入 Naive UI 组件
import { useMessage, NButton, useDialog } from "naive-ui";
import { componentGroups, edgeFields, components } from "./procs";
import { xxtea } from "@/utils/xxtea";
import sflowApi, { SFlow, statusMapping } from '@/api/sflow'
import dagflowApi from '@/api/dagflow' 

// 页面状态
const route = useRoute()
const router = useRouter()
const dialog = useDialog()
// 确保 NButton 被正确导入
// 如果使用了 auto-import 插件，可能不需要显式导入，但需要确保配置正确

const instance = getCurrentInstance();
if (instance && instance.proxy && instance.proxy.$registerDagNode) {
  instance.proxy?.$registerDagNode(components);
}
const message = useMessage();
const hasUnsavedChanges = ref(false);

const state = reactive({
  title: '作业',
  dagId: 0,
  result:{} as any,
  content: {} as any,
  groups: componentGroups,
  edgeFields: edgeFields,
  initNodes: [{
    id: "0",
    shape: "start",
    x: function () {
      return 120;
    },
    y: function () {
      return 30;
    }
  }, {
    id: "999999",
    shape: "end",
    x: function (area: { width: number, height: number }) {
      return area.width - 400;
    },
    y: function (area: { width: number, height: number }) {
      return area.height - 120;
    }
  }]
});

// 监听内容变化
const handleContentChange = () => {
  hasUnsavedChanges.value = true;
  console.log("handleContentChange:---");
};

// 页面退出提示
const promptUnsavedChanges = (next?: Function) => {
  if (hasUnsavedChanges.value) {
    dialog.warning({
      title: '未保存的更改',
      content: '您有未保存的更改，确定要离开吗？',
      positiveText: '确定离开',
      negativeText: '取消',
      onPositiveClick: () => {
        hasUnsavedChanges.value = false;
        if (next) next();
      }
    });
    return false;
  }
  return true;
};

// 监听路由变化
const originalPush = router.push;
router.push = function(location) {
  if (promptUnsavedChanges(() => originalPush.call(router, location))) {
    return originalPush.call(router, location);
  }
  return Promise.resolve(false);
};

// 监听浏览器关闭/刷新
const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (hasUnsavedChanges.value) {
    e.preventDefault();
    e.returnValue = '';
    return '';
  }
};

onMounted(() => {
  state.dagId = Number(route.params.id as string);
  console.log("dagId:", state.dagId)
  handleLoad();
  window.addEventListener('beforeunload', handleBeforeUnload);
});

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload);
});

// 返回列表页
const listHandler = () => {
  if (promptUnsavedChanges(() => router.push({ name: 'sflow_list' }))) {
    router.push({ name: 'sflow_list' });
  }
};

// 保存dag图数据
const handleSave = async (json: any) => {
  console.log("handleSave:", json);
  const form = {
    id: state.dagId,
    content: xxtea.encryptAuto(JSON.stringify(json))
  }
  const data = await sflowApi.saveContent( form);
  if (data) {
    hasUnsavedChanges.value = false;
    message.success('保存成功！');
  } else {
    message.error('保存失败！');
  }
}
const handleLoad = async () => {
  sflowApi.load(state.dagId).then((res) => {
    if (res.code === 200 && res.data) {
      state.title = "数据作业 -- " + res.data.name + " - (" + res.data.id + ")"
      document.title = state.title + " -- 设计器"
      const content = xxtea.decryptAuto(res.data.content);
      state.content = JSON.parse(content || '[]');
      console.log("load:", state.content)
    } else {
      message.error('读取数据失败');
    }
  }).catch((err) => {
    console.error(err);
    message.error('读取数据失败');
  });
}
const handleTest = () => {
  dagflowApi.debug(state.dagId).then((res) => {
    console.log("handleTest:", res);
    if (res.data) {
      state.result = res.data;
      console.log("data:", state.result);
    } 
  }).catch((err) => {
    console.error(err);
    message.error('读取数据失败');
  });
}
</script>

