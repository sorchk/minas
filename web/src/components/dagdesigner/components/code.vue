<template>
  <n-drawer v-model:show="active"
         title="编辑代码"
         :width="drawerWidth"
         placement="right"
         :mask-closable="false"
         :close-on-esc="false"
         :on-close="handleClose">
    <n-space vertical :size="12" style="height: 100%;">
      <vseditor lang="json" style="height: calc(100% - 50px);" ref="codeEditorRef" v-model="state.code"></vseditor>
      <n-space justify="end" style="margin-top: 8px;">
        <n-button size="small" @click="handleFormat">格式化代码</n-button>
        <n-button size="small" @click="handleUpdate">更 新</n-button>
        <n-button size="small" @click="handleClose">关 闭</n-button>
      </n-space>
    </n-space>
  </n-drawer>
</template>
  
<script lang="ts" setup>
import { useDialog } from 'naive-ui'; 
import { useMessage, NSpace, NDrawer, NButton, NGrid, NGridItem, NForm, NFormItem, NInput, NIcon, NA } from 'naive-ui';

import { reactive, ref, nextTick, computed, Ref } from 'vue';
import vseditor from "./attributeform/vseditor/index.vue";
const codeEditorRef = ref();
const active = ref(false)
const emit = defineEmits(["update:modelValue", "upgrade"]);
export interface CurrentPageData {
    oldCode: string;
    code: string;
    visible: Ref<boolean>;
}
const state = reactive<CurrentPageData>({
    oldCode: '',
    code: '',
    visible: ref(false),
});

// 计算抽屉宽度为窗口宽度的80%
const drawerWidth = computed(() => {
    return window.innerWidth * 0.8;
});

const dialog = useDialog();
//打开属性窗口
const openDialog = function (code: string) {
    state.oldCode = code;
    state.code = code;
    state.visible = true;
    
    // 使用 nextTick 确保在DOM更新后再设置编辑器内容
    nextTick(() => {
        if (codeEditorRef.value) {
            codeEditorRef.value.setValue(code);
            // 确保编辑器刷新并聚焦
            setTimeout(() => {
                codeEditorRef.value.refresh && codeEditorRef.value.refresh();
                codeEditorRef.value.focus && codeEditorRef.value.focus();
            }, 100);
        }
    });
}
const handleUpdate = function () {
    state.oldCode = codeEditorRef.value.getValue();
    emit("upgrade", state.oldCode);
}
const handleFormat = function () {
    codeEditorRef.value.format();
}
//关闭属性窗口
const handleClose = function () {
    if (state.oldCode != codeEditorRef.value.getValue()) {
        dialog.warning({
            title: '警告',
            content: '代码发生变化，关闭前是否需要更新?',
            positiveText: '是',
            negativeText: '否',
            onPositiveClick: () => {
                handleUpdate();
                state.visible = false;
            },
            onNegativeClick: () => {
                state.visible = false;
            }
        });
    } else {
        state.visible = false;
    }
}
// 暴露变量
defineExpose({
    openDialog, handleClose,
});
</script>