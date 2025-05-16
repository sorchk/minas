<template>
  <div ref="contentMenu" v-show="state.isShow" id="x6-menu-wrap" class="contextmenu">
    <!-- <li v-if="state.isSelected && state.canCopy" @click="copyNode">复制</li> -->
    <li v-if="state.isSelected" @click="deleteCell">删除</li>
    <li v-if="state.isSelected" @click="onAttrs">属性</li>
  </div>
</template>
<script setup lang="ts">
const emit = defineEmits(['attrs']);
import { Cell, Graph } from '@antv/x6';
import {
  computed,
  onMounted,
  reactive, ref
} from 'vue';
// 定义父组件传过来的值
const props = defineProps({
  getGraph: {
    type: Function, default: () => { return () => { } }
  },
  funcs: {
    type: Function, default: () => { return () => { return {}; } }
  }


});
const state = reactive({
  isShow: false,
  isSelected: false,
  canPaste: false,
  canCopy: false,
});
// 定义变量内容
const contentMenu = ref();
const graph = computed((): Graph => {
  return props.getGraph();
});
// 页面加载时
onMounted(() => {
  if (graph.value) {
    initEvent();
  }
});
//初始化事件
const initEvent = () => {
  //点击右键
  graph.value.on("cell:contextmenu", ({ e, cell }: { e: any, cell: Cell }) => {
    const shape = cell.shape;
    if (shape === "start" || shape === "end") return;
    state.isShow = true;
    state.isSelected = true;
    graph.value.resetSelection(cell);
    state.canCopy = graph.value.isClipboardEnabled();
    if (cell.shape == "dag-edge") {
      state.canCopy = false;
    }
    contentMenu.value.style.top = e.clientY + 15 + "px";
    contentMenu.value.style.left = e.clientX + 15 + "px";
  });
  //点击右键
  graph.value.on("blank:contextmenu", ({ e }: { e: any }) => {
    graph.value.cleanSelection();
    state.isShow = true;
    state.isSelected = false;
    state.canCopy = false;
    contentMenu.value.style.top = e.clientY + 5 + "px";
    contentMenu.value.style.left = e.clientX + 5 + "px";
  });
  //点击空白
  graph.value.on("blank:click", () => {
    cancelShow();
  });
  //点击元素
  graph.value.on("cell:click", () => {
    cancelShow();
  });
}
//复制
// const copyNode = () => {
//   const cells = graph.value.getSelectedCells();
//   if (cells.length) {
//     graph.value.copy(cells);
//   }
//   const id = props.funcs().newId();
//   graph.value.paste({ offset: 20, nodeProps: { id: id } });
//   graph.value.cleanSelection();
//   cancelShow();
// }
// 删除边或节点和边
const deleteCell = () => {
  const cell = graph.value.getSelectedCells();
  graph.value.removeCells(cell);
  cancelShow();
}
//打开属性
const onAttrs = () => {
  emit('attrs');
  cancelShow();
}
//关闭右键菜单
const cancelShow = () => {
  state.isShow = false;
}
defineExpose({ initEvent })
</script>