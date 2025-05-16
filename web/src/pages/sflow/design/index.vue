<template>
  <s-dag-designer :title="state.title" :initNodes="state.initNodes" :content="state.content" :groups="state.groups"
    :edgeFields="edgeFields" openGroup="base" @exec="handleTest" @save="handleSave" />
</template>

<script lang="ts" setup>

import rsapi from "@/api";
import { useRoute, useRouter } from 'vue-router'
import { onMounted, reactive, getCurrentInstance } from "vue";
// 导入 Naive UI 组件
import { useMessage, NButton } from "naive-ui";
import { componentGroups, edgeFields, components } from "./procs";
import { xxtea } from "@/utils/xxtea";
import sflowApi, { SFlow, statusMapping } from '@/api/sflow'

// 页面状态
const route = useRoute()
const router = useRouter()
// 确保 NButton 被正确导入
// 如果使用了 auto-import 插件，可能不需要显式导入，但需要确保配置正确

const instance = getCurrentInstance();
if (instance && instance.proxy && instance.proxy.$registerDagNode) {
  instance.proxy?.$registerDagNode(components);
}
const message = useMessage();
const state = reactive({
  title: '作业',
  dagId: 0,
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
onMounted(() => {
  state.dagId = Number(route.params.id as string);
  console.log("dagId:", state.dagId)
  handleLoad();
})
// 保存dag图数据
const handleSave = async (json: any) => {
  console.log("handleSave:", json);
  const form = {
    id: state.dagId,
    content: xxtea.encryptAuto(JSON.stringify(json), state.dagId)
  }
  const data = await rsapi.post("/sflow/sflow/saveContent", form);
  if (data) {
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
      const content = xxtea.decryptAuto(res.data.content, res.data.id);
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
  rsapi.get("/sapi/run/start/" + state.dagId);
  message.success('已提交后台运行,详细请查看运行日志');
}
</script>

