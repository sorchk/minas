<template>
  <div ref="nodeDiv" :class="className" class="node status">
    <i :class="icon" />
    <label class="label" :style="`width:${autoWidth}px;`">{{ label }}</label>
  </div>
</template>
<script lang="ts" setup name="task-node">
import { computed, inject, onBeforeUnmount, onMounted, reactive } from 'vue';
import consts from '../graph/consts';
const props = defineProps({
  data: {
    type: Object
  }
});
//计算状态样式
const className = computed(() => {
  const classNameTmp = {
    'disabled': !state.status && !state.enabled,
    'success': state.status == 'completed',
    'failed': state.status == 'failed',
    'running': state.status == 'running',
    'pending': state.status == 'pending',

  }
  return classNameTmp
})
//标题
const label = computed(() => {
  return state.label || props.data?.name || '未命名'
});
//图标
const icon = computed(() => {
  return props.data?.icon || 'fa fa-microchip'
});
//状态数据
const state = reactive({
  label: props.data?.name,
  status: '',
  enabled: true,
  readonly: false,
  canvas: false as any,
  autoWidth: consts.nodeDefaultWidth
});
const autoWidth = computed(() => {
  return state.autoWidth;
})
const getNode = inject('getNode') as Function;
const fixWidth = () => {
  const node = getNode();
  const size = node.size();
  const autoWidth = 50 + getTextWidth(state.label, 'normal 12px Arial');
  state.autoWidth = Math.min(consts.nodeMaxWidth, Math.max(consts.nodeMinWidth, autoWidth)) - 50;
  node.size({ width: state.autoWidth + 50, height: size.height });
}
//数据改变
const change = function (data: any) {
  console.log("node data change:", data)
  state.label = data.form.label;
  state.status = data.status || ''
  state.enabled = data.form.enabled != false;
  state.readonly = data.readonly || false;
  fixWidth();
}


const getTextWidth = (text: string, font: any) => {
  if (!state.canvas) {
    state.canvas = document.createElement("canvas");
  }
  if (state.canvas) {
    const context = state.canvas.getContext("2d")
    if (context) {
      context.font = font
      const metrics = context.measureText(text)
      return metrics.width
    }
  }
  return 0;
}

// 页面加载时
onMounted(() => {
  const node = getNode();
  //初始化数据
  if (!node.data) {
    node.data = {};
  }
  if (!node.data.form) {
    node.data.form = {};
  }
  state.label = node.data.form.label;
  state.status = node.data.status || '';
  state.enabled = node.data.form.enabled != false;
  state.readonly = node.data.readonly || false;
  // 监听数据改变事件
  node.on('change:data', (data: any) => {
    change(data.current);
  });
  fixWidth();
});
// 页面销毁时
onBeforeUnmount(() => {
});
</script>
