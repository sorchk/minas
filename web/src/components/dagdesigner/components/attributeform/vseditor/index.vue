<template>
  <div ref="codeContainer" :style="`height: ${height}px`" />
</template>

<script lang="ts" setup name="vs-editor">
import * as monaco from 'monaco-editor'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import cssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import htmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import { nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

// 定义子组件向父组件传值/事件
const emit = defineEmits(["update:modelValue"]);
// 定义父组件传过来的值
const props = defineProps({
  modelValue: { type: String, default: "" },
  lang: { type: String, default: '' },
  readOnly: { type: Boolean, default: false },
  minimapEnabled: { type: Boolean, default: true },
  theme: { type: String, default: 'vs-dark' },
  automaticLayout: { type: Boolean, default: true },
  lineHeight: { type: Number, default: 20 },
  tabSize: { type: Number, default: 2 },
  height: { type: String, default: "300" },
});
const codeContainer = ref()
//监控数据变化
self.MonacoEnvironment = {
  getWorker(_: string, label: string) {
    if (label === 'json') {
      return new jsonWorker()
    }
    if (label === 'css' || label === 'scss' || label === 'less') {
      return new cssWorker()
    }
    if (label === 'html' || label === 'handlebars' || label === 'razor') {
      return new htmlWorker()
    }
    if (['typescript', 'javascript'].includes(label)) {
      return new tsWorker()
    }
    return new EditorWorker()
  },
}
let editor: monaco.editor.IStandaloneCodeEditor;

const state = reactive({
  data: "",
});
// 页面加载时
onMounted(() => {
  state.data = props.modelValue;
  editorInit();
});
watch(() => props.modelValue,
  (newValue) => {
    if (newValue != state.data) {
      editor.setValue(newValue);
    }
  });
onBeforeUnmount(() => {
  returnValue();
  editor.dispose()
})
const returnValue = () => {
  emit("update:modelValue", editor.getValue());
}
const format = () => {
  nextTick(() => {
    console.log("editor:", editor)
    if (props.readOnly) {
      editor.updateOptions({ readOnly: false });
    }
    editor.getAction('editor.action.formatDocument').run();//自动格式化代码
    setTimeout(function () {
      state.data = editor.getValue();
      console.log("format:", state.data)
    }, 200);
    if (props.readOnly) {
      setTimeout(function () {
        editor.updateOptions({ readOnly: true });
      }, 200);
    }
  })
}
const getValue = () => {
  state.data = editor.getValue();
  return state.data;
}
const editorInit = () => {
  nextTick(() => {
    monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
      noSemanticValidation: true,
      noSyntaxValidation: false
    });
    monaco.languages.typescript.javascriptDefaults.setCompilerOptions({
      target: monaco.languages.typescript.ScriptTarget.ES2016,
      allowNonTsExtensions: true
    });
    if (!editor) {
      editor = monaco.editor.create(codeContainer.value, {
        value: '', // 编辑器初始显示文字
        language: props.lang, // 语言支持自行查阅demo
        automaticLayout: props.automaticLayout, // 自适应布局
        theme: props.theme, // 官方自带三种主题vs, hc-black, or vs-dark
        foldingStrategy: 'indentation',
        renderLineHighlight: 'all', // 行亮
        selectOnLineNumbers: true, // 显示行号
        tabSize: props.tabSize, // 缩进
        lineHeight: props.lineHeight, // 行高
        minimap: {
          enabled: props.minimapEnabled,
        },
        readOnly: props.readOnly, // 只读
        fontSize: 16, // 字体大小
        scrollBeyondLastLine: false, // 取消代码后面一大段空白
        overviewRulerBorder: false, // 不要滚动条的边框
      });
    }
    editor.setValue(state.data);
    format();
  })
}
defineExpose({ getValue, format })
</script>

<style scoped></style>
