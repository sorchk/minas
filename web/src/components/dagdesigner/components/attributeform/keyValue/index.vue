<template>
  <div style="width: 100%;">
    <n-form-item>
      <n-button type="primary" size="small" ghost @click="onAdd">添加{{ name }}</n-button>
      <n-a class="ml10" @click="onCopy">复制</n-a>
      <n-a class="ml10" v-if="state.canPaste" @click="onPaste">粘贴</n-a>
    </n-form-item>
    <n-form-item v-for="(item, i) in state.values" :key="i" :label="'键[' + i + ']：'">
      <n-grid :cols="24" :x-gap="5">
        <n-grid-item :span="6">
          <n-input size="small" v-model:value="item.key" placeholder="请输入"></n-input>
        </n-grid-item>
        <n-grid-item :span="1" style="margin: 0 5px 0 5px;">
          值：
        </n-grid-item>
        <n-grid-item :span="15">
          <n-input size="small" v-model:value="item.value" :placeholder="placeholder"></n-input>
        </n-grid-item>
        <n-grid-item :span="1">
          <n-button size="small" circle type="error" @click="onDel(i)">
            <template #icon><n-icon><DeleteOutlined /></n-icon></template>
          </n-button>
        </n-grid-item>
      </n-grid>
    </n-form-item>
  </div>
</template>

<script setup lang="ts" >
import { onBeforeUnmount, onMounted, reactive, watch } from 'vue';
import useClipboard from "vue-clipboard3";
import { DeleteOutlined } from "@vicons/antd";
import { useMessage, NButton, NGrid, NGridItem, NForm, NFormItem, NInput, NIcon, NA } from 'naive-ui';

const { toClipboard } = useClipboard();
import { KVMap } from "../../../interface";

// 定义子组件向父组件传值/事件
const emit = defineEmits(["update:modelValue"]);
// 定义父组件传过来的值
const props = defineProps({
  modelValue: { type: <any>Array, default: () => new Array<any>() },
  name: { type: String, default: '' },
  placeholder: { type: String, default: '' },
});
const state = reactive({
  values: new Array<KVMap>(),
  canPaste: false
});
const message = useMessage();

// 页面加载时
onMounted(async () => {
  state.values = props.modelValue;
  state.canPaste = await canPaste();
});
watch(() => props.modelValue,
  (newValue) => {
    if (newValue != state.values) {
      state.values = newValue;
    }
  });

const COPY_PREFIX = "dag:keyValues:"
const canPaste = async () => {
  const data = await navigator.clipboard.readText();
  return data.indexOf(COPY_PREFIX) == 0;
}

const onPaste = async () => {
  try {
    const data = await navigator.clipboard.readText();
    if (data.indexOf(COPY_PREFIX) == 0) {
      const values = data.substring(COPY_PREFIX.length).split(", ");
      state.values = new Array<KVMap>();
      values.forEach(item => {
        const kv = item.split(":", 2);
        state.values.push({ key: kv[0], value: kv[1] });
      })
      message.success('已粘贴数据')
    }
  } catch (e) { }
}
const onCopy = async () => {
  try {
    const values = new Array<string>();
    if (state.values.length > 0) {
      state.values.forEach(item => {
        values.push(item.key + ':' + item.value);
      })
      const data = COPY_PREFIX + values.join(', ');
      toClipboard(data)
      message.success('已复制到剪切板');
      // 复制成功
    }
  } catch (e) {
    // 复制失败
    console.log('复制失败:', e)
    message.error('复制失败')
  }
}
// 页面销毁时
onBeforeUnmount(() => {
  emit("update:modelValue", state.values);
});
//定义方法
const onAdd = () => {
  const map = new KVMap();
  state.values.push(map);
  emit("update:modelValue", state.values);
}
const onDel = (i: number) => {
  state.values.splice(i, 1);
  emit("update:modelValue", state.values);
}
</script>

<style scoped>
.ml10 {
  margin-left: 10px;
}
</style>
