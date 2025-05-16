<template>
    <div style="width: 100%;">
        <el-form-item class="mb5">
            <el-button size="small" type="primary" round @click="onAdd">添加{{ name }}</el-button>
            <el-link class="ml10" type="primary" style="margin-bottom: 2px;" round @click="onCopy">复制</el-link>
            <el-link class="ml10" v-if="state.canPaste" type="primary" style="margin-bottom: 2px;" round
                @click="onPaste">粘贴</el-link>
        </el-form-item>
        <el-row class="mb10" :gutter="5" style="width: 100%;" v-for="(item, i) in  state.values " :key="i">
            <el-col :span="5">
                <el-form-item class="mb5" :label="'数据项[' + i + ']：'">
                    <el-input size="small" v-model="item.key" :placeholder="'请输入' + name"></el-input>
                </el-form-item>
            </el-col>
            <el-col :span="3"><el-form-item class="mb5" label="类型：">
                    <el-select size="small" v-model="item.type" placeholder="请选择" :multiple="false" :clearable="false"
                        :collapse-tags="false" :multiple-limit="1">
                        <el-option v-for=" item  in  dataTypes " :key="item.value" :label="item.label"
                            :value="item.value" />
                    </el-select>
                </el-form-item>
            </el-col>
            <el-col :span="2">
                <el-form-item class="mb5" label="必须：">
                    <el-switch size="small" v-model="item.required" value="true" :active-text="item['active']"
                        :inactive-text="item['inactive']"></el-switch>
                </el-form-item>
            </el-col>
            <el-col :span="14"><el-form-item class="mb5" label="xpath表达式：">
                    <el-input size="small" v-model="item.xpath" placeholder="在浏览器按F12，选择元素，选择需要提取的页面元素，右键复制-》复制xpath或复制完整xpath"></el-input></el-form-item>
            </el-col>
            <el-col :span="11"><el-form-item class="mb5" label="正则提取：">
                    <el-input size="small" v-model="item.regex" placeholder="请输入正则表达式，使用分组提取数据"></el-input></el-form-item>
            </el-col>
            <el-col :span="12"><el-form-item class="mb5" label="输出模板：">
                    <el-input size="small" v-model="item.template" placeholder="例如：名称是$2；账号是$1。说明：$0为整个匹配结果，$1为第一个分组，依次类推$2 $3 ..."></el-input></el-form-item>
            </el-col>
            <el-col :span="1">
                <el-button size="small" type="danger" circle @click="onDel(i)">
                    <el-icon>
                        <Delete />
                    </el-icon>
                </el-button>
            </el-col>
        </el-row>
    </div>
</template>
<script lang="ts" setup name="crawler-xpath">
import { onBeforeUnmount, onMounted, reactive, watch } from 'vue'
import { CrawlerXpathMap } from "../../../interface";
import useClipboard from "vue-clipboard3";
import { useMessage, NButton, NGrid, NGridItem, NForm, NFormItem, NInput, NIcon, NA } from 'naive-ui';

const { toClipboard } = useClipboard();
// 引入组件
// 定义子组件向父组件传值/事件
const emit = defineEmits(['update:modelValue']);


const message = useMessage();
// 定义父组件传过来的值
const props = defineProps({
    modelValue: {
        type: <any>Array, default: () => {
            return new Array<any>()
        }
    },
    name: { type: String, default: '' }
});
const state = reactive({
    values: new Array<CrawlerXpathMap>(),
    canPaste: false
});
const dataTypes = [{
    label: '字符串',
    value: 'string'
}, {
    label: 'HTML代码',
    value: 'html'
},{
    label: '字符串集合',
    value: 'strings'
}, {
    label: 'HTML代码集合',
    value: 'htmls'
}, {
    label: '数字',
    value: 'number'
}];
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

const COPY_PREFIX = "dag:params:"
const canPaste = async () => {
    const data = await navigator.clipboard.readText();
    return data.indexOf(COPY_PREFIX) == 0;
}
const onPaste = async () => {
    try {
        const data = await navigator.clipboard.readText();
        if (data.indexOf(COPY_PREFIX) == 0) {
            state.values = JSON.parse(data.substring(COPY_PREFIX.length));
           
      message.success('已粘贴数据')
            
        }
    } catch (e) { }
}
const onCopy = async () => {
    try {
        if (state.values.length > 0) {
            const data = COPY_PREFIX + JSON.stringify(state.values);
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
    const map = new CrawlerXpathMap();
    state.values.push(map);
    emit("update:modelValue", state.values);
}
const onDel = (i: number) => {
    state.values.splice(i, 1);
    emit("update:modelValue", state.values);
}
</script>

<style scoped></style>
