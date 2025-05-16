<template>
  <n-modal v-model:show="state.visible"  
           :style="{ width: '80%' }"
           preset="card"
           :mask-closable="false"
           :close-on-esc="false"
           :segmented="{ content: true, footer: 'soft' }"
           :transform-origin="'center'"
           :trap-focus="false"
           :on-close="handleClose"
           :bordered="false"
           :fullscreen="state.fullscreen">
    <template #header-extra>
      <n-space>
        <n-button v-if="canPaste" quaternary circle type="info" title="粘贴数据" @click="handlePaste">
          <template #icon><n-icon><ClipboardOutline /></n-icon></template>
        </n-button>
        <n-button quaternary circle type="primary" title="复制数据" @click="handleCopy">
          <template #icon><n-icon><CopyOutlined /></n-icon></template>
        </n-button>
        <n-button quaternary circle type="default" :title="state.fullscreen ? '还原' : '最大化'" @click="handleFullSceen">
          <template #icon>
            <n-icon>
              <component :is="state.fullscreen ? 'FullscreenExitOutlined' : 'FullscreenOutlined'" />
            </n-icon>
          </template>
        </n-button>
      </n-space>
    </template>
    
    <n-form>
      <div v-for="item in state.fields" :key="item.prop">
        <n-form-item v-if="!item.vif || item.vif(state.form)" :path="item.prop" :label="item.label"
                    :label-width="setPx(item['labelWidth'] || 120)" size="small">
          <template #label>
            <div>
              <span :title="item.placeholder">{{ item.label }}</span>
              <n-tooltip v-if="!!item.help" trigger="hover" :placement="'top-start'">
                <template #trigger>
                  <n-icon size="16"><QuestionCircleOutlined /></n-icon>
                </template>
                {{ item.help }}
              </n-tooltip>
            </div>
          </template>
          
          <!-- 参数表单 -->
          <dag-params v-if="item.type == 'params'" :name="item['label']" v-model="state.form[item.prop]"></dag-params>
          
          <!-- 输入框 -->
          <n-input v-if="item.type == 'input' || item.type == 'text' || item.type == 'password' || item.type == 'hidden' || item.type == 'textarea' || !item.type"
                  v-model:value="state.form[item.prop]" 
                  :type="item.type === 'password' ? 'password' : item.type === 'textarea' ? 'textarea' : 'text'"
                  :rows="item['rows']" 
                  :disabled="item['disabled']"
                  :clearable="item['clearable']"
                  :show-password-on="item['showPassword']"
                  :autosize="item.type === 'textarea' ? { minRows: 2, maxRows: 6 } : undefined"
                  :readonly="item['readonly']"
                  :placeholder="item.placeholder">
          </n-input>
          
          <!-- 输入列表 -->
          <dag-inputs v-if="item.type == 'inputs'" :name="item['label']" v-model="state.form[item.prop]"
            :placeholder="item.placeholder"></dag-inputs>
          
          <!-- 键值对 -->
          <dag-key-value v-if="item.type == 'map'" :name="item['label']" v-model="state.form[item.prop]"
            :placeholder="item.placeholder"></dag-key-value>
          
          <!-- 数字输入 -->
          <n-input-number v-if="item.type == 'number'" 
                        v-model:value="state.form[item.prop]"
                        :precision="item.precision" 
                        :min="item.min" 
                        :max="item.max"
                        :step="item.step">
          </n-input-number>

          <!-- 单选组 -->
          <n-radio-group v-if="item.type == 'radio'" v-model:value="state.form[item.prop]" size="small">
            <n-radio v-for="v in item.data" :key="v.value" :value="v['value']">
              {{ v['label'] }}
            </n-radio>
          </n-radio-group>

          <!-- 编辑器 -->
          <VsEditor v-if="item.type == 'editor'" v-model="state.form[item.prop]" :lang="item.lang" style="width:100%;">
          </VsEditor>

          <!-- 开关 -->
          <n-switch v-if="item.type == 'switch'" v-model:value="state.form[item.prop]">
            <template #checked>{{ item.active }}</template>
            <template #unchecked>{{ item.inactive }}</template>
          </n-switch>
          
          <!-- 选择器 -->
          <n-select v-if="item.type == 'select'" 
                  v-model:value="state.form[item.prop]"
                  :placeholder="item.placeholder" 
                  :multiple="item.multiple" 
                  :clearable="item.clearable"
                  :max-tag-count="item.collapseTags ? 2 : undefined"
                  :options="item.data && item.data.map((v:any)=> ({
                    label: v[item.labelName || 'label'],
                    value: v[item.valueName || 'value'],
                    disabled: v['disabled']
                  }))">
          </n-select>
          
          <!-- 树形选择 -->
          <n-tree-select v-if="item.type == 'selectTree'" 
                      v-model:value="state.form[item.prop]"
                      :placeholder="item.placeholder" 
                      :options="item.data && formatTreeData(item.data, item.labelName || 'label', item.valueName || 'value')" 
                      :key-field="item.valueName || 'value'"
                      :label-field="item.labelName || 'label'"
                      :children-field="'children'"
                      clearable>
          </n-tree-select>
        </n-form-item>
      </div>
    </n-form>
    
    <template #footer>
      <n-space justify="end">
        <n-button size="small" @click="handleClose">关 闭</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>

import rsapi from "@/api";
import { Cell, Edge } from '@antv/x6';
import { reactive, computed } from 'vue';
import { useMessage, NIcon, NTreeSelect, NSwitch, NSelect, NRadio, NRadioButton, NRadioGroup, NInput, NInputNumber, NTooltip, NModal, NSpace, NForm, NFormItem, NButton } from 'naive-ui';
import { 
  CopyOutlined, 
  FullscreenOutlined,
  FullscreenExitOutlined,
  QuestionCircleOutlined 
} from '@vicons/antd';

import { 
  ClipboardOutline
} from '@vicons/ionicons5';

import DagInputs from "../attributeform/inputs/index.vue";
import DagKeyValue from "../attributeform/keyValue/index.vue";
import VsEditor from "../attributeform/vseditor/index.vue";
import DagParams from "./params/index.vue";

import { deepClone } from '../../graph/shape';

const message = useMessage();

interface Field {
  prop: string,
  type: string,
  label: string,
  valueName: string,
  labelName: string,
  placeholder: string,
  active: string,
  inactive: string,
  lang: string,
  resize: boolean,
  readonly: boolean,
  autosize: boolean,
  showPassword: boolean,
  disabled: boolean,
  rows: number,
  precision: number,
  max: number,
  min: number,
  step: number,
  multiple: boolean,
  clearable: boolean,
  collapseTags: boolean,
  multipleLimit: number,
  filterable: boolean,
  props: any,
  data: any
}

export interface CurrentPageData {
  cell: Cell | Edge | null | undefined;
  fields: Array<any>;
  form: any;
  proc: any,
  labelWidth: number;
  visible: boolean;
  fullscreen: boolean;
  title: String;
}

const state = reactive<CurrentPageData>({
  cell: null,
  fields: new Array<Field>(),
  form: {} as any,
  proc: {} as any,
  labelWidth: 120,
  visible: false,
  fullscreen: false,
  title: '',
});

// 格式化树形数据 
const formatTreeData = (data: any[], labelField: string, valueField: string): Array<{ label: string; key: string; disabled?: boolean; children?: any[] }> => {
  if (!Array.isArray(data)) return [];
  
  return data.map(item => ({
    label: item[labelField],
    key: item[valueField],
    disabled: item.disabled,
    children: item.children ? formatTreeData(item.children, labelField, valueField) : undefined
  }));
};

//打开属性窗口
const openDialog = function (cell: Cell, procs: any) {
  state.title = '配置属性';
  if (cell) {
    state.cell = cell;
    state.title = '配置属性 [' + state.cell?.shape + ']  task_' + state.cell?.id
    state.proc = procs[cell.shape];
    state.fields = deepClone(state.proc.fields || new Array<Field>());
    console.log("openDialog cell:", cell)
    //如果无数据则初始化
    if (!cell.getData()) {
      cell.setData({ form: {} });
    }
    const data = cell.getData();
    //从节点加载表单数据
    state.form = data.form || {};
    //修复空标签
    if (!state.form.label) {
      state.form.label = state.proc.name;
    }
    init();
    state.visible = true;
    console.log("openDialog cell state:", state)
  } else {
    message.error('无效节点')
  }
};

const canPaste = computed(() => {
  const copyText = localStorage.getItem("dag_copy_form");
  if (copyText) {
    try {
      const data = JSON.parse(copyText);
      if (data && data.type == state.cell?.shape && data.form) {
        return true;
      }
    } catch (e) { }
  }
  return false;
});

const handleFullSceen = () => {
  state.fullscreen = !state.fullscreen;
}

const handlePaste = () => {
  const copyText = localStorage.getItem("dag_copy_form");
  if (copyText) {
    try {
      const data = JSON.parse(copyText);
      if (data && state.cell && data.type == state.cell?.shape && data.form) {
        state.form = data.form
        const cellData = state.cell.getData();
        cellData['form'] = data.form;
        message.success('已粘贴数据')
      }
    } catch (e) { }
  }
}

const handleCopy = async () => {
  try {
    const data = JSON.stringify({ type: state.cell?.shape, form: state.form });
    localStorage.setItem("dag_copy_form", data);
    message.success('已复制数据')
    // 复制成功
  } catch (e) {
    // 复制失败
    console.log('复制失败:', e)
    message.error('复制失败')
  }
}

//关闭属性窗口
const handleClose = () => {
  //调用父级方法 更新节点数据
  const data = state.cell?.getData();
  const label = state.form.label || state.proc.name

  //更新基础组件名称 
  if (state.cell?.isEdge()) {
    state.cell.setLabels(label);
  } else if (state.cell?.data.parent === true) {
    state.cell?.attr('text/text', label);
  }
  
  //更新数据
  const options = {
    silent: false,
    overwrite: false,
    deep: true,
  };
  
  state.cell?.trigger("change:data", {
    cell: state.cell,
    current: data,
    options: options
  });
  state.visible = false;
}

//初始化表单数据
const init = function () {
  for (let i in state.fields) {
    let field = state.fields[i] as any;
    let fieldValue = state.form[field.prop];
    if (!field['labelName']) {
      state.fields[i]['labelName'] = 'label';
    }
    if (!field['valueName']) {
      state.fields[i]['valueName'] = 'value';
    }
    let value = field['value'];
    if (value === null || value === undefined || value === 'undefined') {
      message.error(state.cell?.shape + '[' + field.prop + ']' + ' 组件定义错误，必须定义value属性的值，可以是空字符串、0、或空数组')
      return;
    }
    if (fieldValue === null || fieldValue === undefined || fieldValue === 'undefined') {
      state.form[field.prop] = value;
    }
    //从url读取组件数据
    let url: string = field['url'];
    if (url !== null && url !== undefined && url !== '' && url.length > 0) {
      let dataRoot: string = field['dataRoot'];
      rsapi.get(url, {}).then((res: any) => {
        let data = new Array();
        console.log("字典数据：", res)
        if (dataRoot) {
          data = res[dataRoot];
        } else {
          data = res;
        }
        if (field.type == 'selectTree') {
          const appendDisabled = (item: any) => {
            if (!item['isLeaf']) {
              item['disabled'] = true;
            }
          }
          data.forEach(appendDisabled);
          state.fields[i].data = rsapi.listToTree(data, {
            id: "id",
            parentId: "parentId",
            children: "children",
            disabled: "disabled",
            isLeaf: "isLeaf",
            enabledParent: true
          });
        } else {
          state.fields[i].data = data;
        }
      }).catch((e: Error) => {
        message.error('读取数据错误:' + e)
      });
    }
  }
}

//计算标签宽度
const setPx = (val: any, defval = '') => {
  if (val === null || val === undefined) {
    val = defval;
  }
  if (val === null || val === undefined) {
    return '';
  }
  val = val + '';
  if (val.indexOf('%') === -1) {
    val = val + 'px';
  }
  return val;
};

// 暴露变量
defineExpose({
  openDialog, handleClose, setPx
});
</script>

<style scoped>
/* 移除Element Plus样式，添加Naive UI样式 */
</style>