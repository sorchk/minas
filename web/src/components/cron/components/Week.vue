<template>
  <n-card>
    <n-radio-group v-model:value="radioValue" name="selectWeek">
      <n-grid x-gap="20" y-gap="20" :cols="1">
        <n-gi>
          <n-flex align="center">
            <n-radio :value="1" /> 每周，允许的通配符[, - * ? / L #]
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center"> <n-radio :value="2" /> 不指定 </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="3" />
            周期从
            <n-radio-group v-model:value="cycle01" name="radiobuttongroup1" size="small" @update:value="radioChange(3)">
              <n-radio-button v-for="song in weekList" :key="song.value" :value="song.value"
                :disabled="song.value === 1" :label="song.label" />
            </n-radio-group>
            <div>到</div>
            <n-radio-group v-model:value="cycle02" name="radiobuttongroup2" size="small" @update:value="radioChange(3)">
              <n-radio-button v-for="song in weekList" :key="song.value" :value="song.value"
                :disabled="song.value < cycle01 && song.value !== 1" :label="song.label" />
            </n-radio-group>
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="4" />
            周期从第
            <n-input-number class="nin" v-model:value="average01" :min="1" :max="4" clearable
              @update:value="radioChange(4)" @focus="radioChange(4)" />
            周的
            <n-radio-group v-model:value="average02" name="radiobuttongroup3" size="small"
              @update:value="radioChange(4)">
              <n-radio-button v-for="song in weekList" :key="song.value" :value="song.value" :label="song.label" />
            </n-radio-group>
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="5" /> 本月最后一个
            <n-radio-group v-model:value="weekday" name="radiobuttongroup4" size="small" @update:value="radioChange(5)">
              <n-radio-button v-for="song in weekList" :key="song.value" :value="song.value" :label="song.label" />
            </n-radio-group>
          </n-flex>
        </n-gi>
        <n-gi>
          <n-flex align="center">
            <n-radio :value="6" />
            指定周
            <div style="flex: 1">
              <n-checkbox-group v-model:value="checkboxList" @update:value="radioChange(6)">
                <n-space>
                  <n-checkbox v-for="item in weekList" :key="item.value" :value="item.value" :label="item.label" />
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
export default { name: "CronGenWeek" };
</script>
<script lang="ts" setup>
import { checkNumber, zeroFill } from "../utils";
import { computed, onMounted, ref, watch } from 'vue';
import { NGi, NCard, NGrid, NTabs, NTabPane, NTag, NRadioButton, NRadioGroup, NSpace, NCheckbox, NCheckboxGroup, NRadio, NFlex, NInput, NInputNumber } from "naive-ui";


const props = defineProps({
  modelValue: {
    required: true,
    type: String,
  },
});

const emit = defineEmits<{
  (e: "update:modelValue", v: string): void;
}>();

const radioValue = ref<number>(2);
const weekday = ref<number>(2);
const cycle01 = ref<number>(2);
const cycle02 = ref<number>(3);
const average01 = ref<number>(1);
const average02 = ref<number>(2);
const checkboxList = ref<number[]>([]);
const weekList = [
  { value: 1, label: "周一" },
  { value: 2, label: "周二" },
  { value: 3, label: "周三" },
  { value: 4, label: "周四" },
  { value: 5, label: "周五" },
  { value: 6, label: "周六" },
  { value: 0, label: "周日" },
];

const value = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});
console.log(value.value);

watch(
  () => value.value,
  () => assign()
);
onMounted(() => assign());

// 赋值
const assign = () => {
  if (value.value === "*") {
    radioValue.value = 1;
  } else if (value.value === "?") {
    radioValue.value = 2;
  } else if (value.value.indexOf("-") > -1) {
    let indexArr = value.value.split("-");
    cycle01.value = isNaN(+indexArr[0]) ? 0 : +indexArr[0];
    cycle02.value = +indexArr[1];
    radioValue.value = 3;
  } else if (value.value.indexOf("#") > -1) {
    let indexArr = value.value.split("#");
    average01.value = isNaN(+indexArr[1]) ? 1 : +indexArr[1];
    average02.value = +indexArr[0];
    radioValue.value = 4;
  } else if (value.value.indexOf("L") > -1) {
    let indexArr = value.value.split("L");
    weekday.value = isNaN(+indexArr[0]) ? 1 : +indexArr[0];
    radioValue.value = 5;
  } else {
    checkboxList.value = value.value.split(",").map((i) => +i);
    radioValue.value = 6;
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
      value.value = "?";
    } else if (v === 3) {
      value.value = cycleTotal.value;
    } else if (v === 4) {
      value.value = averageTotal.value;
    } else if (v === 5) {
      value.value = weekdayCheck.value;
    } else if (v === 6) {
      value.value = checkboxString.value;
    } else {
      value.value = "*";
    }
  }
);

// 计算两个周期值
const cycleTotal = computed(() => {
  return (
    checkNumber(cycle01.value, 0, 6) + "-" + checkNumber(cycle02.value, 0, 6)
  );
});

watch(
  () => cycleTotal.value,
  (v) => (value.value = v)
);

// 计算平均用到的值
const averageTotal = computed(() => {
  return (
    checkNumber(average02.value, 0, 6) +
    "#" +
    checkNumber(average01.value, 1, 4)
  );
});
watch(
  () => averageTotal.value,
  (v) => (value.value = v)
);

// 最近的工作日（格式）
const weekdayCheck = computed(() => {
  return checkNumber(weekday.value, 0, 6) + "L";
});
watch(
  () => weekdayCheck.value,
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
