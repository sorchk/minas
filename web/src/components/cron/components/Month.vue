<template>
  <n-card>
    <n-radio-group v-model:value="radioValue" name="selectMonth">
      <n-grid x-gap="20" y-gap="20" :cols="1">
        <n-gi>
          <n-flex align="center">
            <n-radio :value="1" /> 每月，允许的通配符[, - * /]
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="2" />
            周期从
            <n-input-number class="nin" v-model:value="cycle01" :min="1" :max="11" clearable
              @update:value="radioChange(2)" @focus="radioChange(2)" />
            到
            <n-input-number class="nin" v-model:value="cycle02" :min="cycle01 ? cycle01 + 1 : 2" :max="12" clearable
              @update:value="radioChange(2)" @focus="radioChange(2)" />月
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="3" />
            周期从
            <n-input-number class="nin" v-model:value="average01" :min="1" :max="11" clearable
              @update:value="radioChange(3)" @focus="radioChange(3)" />
            月开始，每
            <n-input-number class="nin" v-model:value="average02" :min="1" :max="12 - average01 || 0" clearable
              @update:value="radioChange(3)" @focus="radioChange(3)" />月执行一次
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex>
            <n-radio :value="4" />
            指定月
            <div style="flex: 1">
              <n-checkbox-group v-model:value="checkboxList" @update:value="radioChange(4)">
                <n-space>
                  <n-checkbox v-for="item in 12" :key="item" :value="item" :label="zeroFill(item)" />
                </n-space>
              </n-checkbox-group>
            </div>
          </n-flex>
        </n-gi>
      </n-grid>
    </n-radio-group>
  </n-card>
</template>
<script lang="ts">
export default { name: "CronGenMonth" };
</script>
<script lang="ts" setup>
import { checkNumber, zeroFill } from "../utils";
import { computed, onMounted, ref, watch } from 'vue';
import { NGi, NCard, NGrid, NTabs, NTabPane, NTag, NRadioGroup, NSpace, NCheckbox, NCheckboxGroup, NRadio, NFlex, NInput, NInputNumber } from "naive-ui";


const props = defineProps({
  modelValue: {
    required: true,
    type: String,
  },
});

const emit = defineEmits<{
  (e: "update:modelValue", v: string): void;
}>();

const radioValue = ref<number>(1);
const cycle01 = ref<number>(1);
const cycle02 = ref<number>(2);
const average01 = ref<number>(1);
const average02 = ref<number>(1);
const checkboxList = ref<number[]>([]);

const value = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});
watch(
  () => value.value,
  () => assign()
);
onMounted(() => assign());

// 赋值
const assign = () => {
  if (value.value === "*") {
    radioValue.value = 1;
  } else if (value.value.indexOf("-") > -1) {
    const indexArr = value.value.split("-");
    cycle01.value = isNaN(+indexArr[0]) ? 0 : +indexArr[0];
    cycle02.value = +indexArr[1];
    radioValue.value = 2;
  } else if (value.value.indexOf("/") > -1) {
    const indexArr = value.value.split("/");
    average01.value = isNaN(+indexArr[0]) ? 0 : +indexArr[0];
    average02.value = +indexArr[1];
    radioValue.value = 3;
  } else {
    radioValue.value = 4;
    checkboxList.value = value.value.split(",").map((i) => +i);
  }
};

// 单选按钮值变化时
const radioChange = (v: number) => {
  if (radioValue.value !== v) {
    radioValue.value = v;
  }
};

watch(
  () => radioValue.value,
  (v) => {
    if (v === 1) {
      value.value = "*";
    } else if (v === 2) {
      value.value = cycleTotal.value;
    } else if (v === 3) {
      value.value = averageTotal.value;
    } else if (v === 4) {
      value.value = checkboxString.value;
    } else {
      value.value = "*";
    }
  }
);

// 计算两个周期值
const cycleTotal = computed(() => {
  const cycle1 = checkNumber(cycle01.value, 0, 58);
  const cycle2 = checkNumber(cycle02.value, cycle1 ? cycle1 + 1 : 1, 59);
  return cycle1 + "-" + cycle2;
});
watch(
  () => cycleTotal.value,
  (v) => (value.value = v)
);

// 计算平均用到的值
const averageTotal = computed(() => {
  const average1 = checkNumber(average01.value, 0, 58);
  const average2 = checkNumber(average02.value, 1, 59 - average1 || 0);
  return average1 + "/" + average2;
});
watch(
  () => averageTotal.value,
  (v) => (value.value = v)
);

// 计算勾选的checkbox值合集
const checkboxString = computed(() => {
  const str = checkboxList.value.join();
  return str == "" ? "*" : str;
});
watch(
  () => checkboxString.value,
  (v) => (value.value = v)
);
</script>
