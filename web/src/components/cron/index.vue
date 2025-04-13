<template>
  <n-grid x-gap="12" y-gap="12" :cols="1">
    <n-gi>
      <n-tabs type="segment" animated>
        <n-tab-pane name="second" tab="秒">
          <Second v-model="cronGenValue.second" />
        </n-tab-pane>
        <n-tab-pane name="min" tab="分">
          <Min v-model="cronGenValue.min" />
        </n-tab-pane>
        <n-tab-pane name="hour" tab="时">
          <Hour v-model="cronGenValue.hour" />
        </n-tab-pane>
        <n-tab-pane name="day" tab="日">
          <Day v-model="cronGenValue.day" />
        </n-tab-pane>
        <n-tab-pane name="month" tab="月">
          <Month v-model="cronGenValue.month" />
        </n-tab-pane>
        <n-tab-pane name="week" tab="周">
          <Week v-model="cronGenValue.week" />
        </n-tab-pane>
      </n-tabs>
    </n-gi>
    <n-gi>
      <n-card title="时间表达式" :segmented="{ content: true }">
        <template #header-extra>
          <div class="step1">
            Cron 表达式：
            <n-tag type="success" :bordered="false">
              {{ cronGenValueString }}
            </n-tag>
          </div>
        </template>
        <div class="cronGenValue step2">
          <div v-for="(item, key, index) in cronGenValue" :key="index">
            {{ map[key] }}：
            <n-tag v-if="item" :bordered="false" size="small">
              {{ item }}
            </n-tag>
          </div>
        </div>
      </n-card>
    </n-gi>
    <n-gi>
      <n-card title="最近5次运行时间" :segmented="{ content: true }">
        <Result v-model="cronGenValueString" class="step3" />
      </n-card>
    </n-gi>
  </n-grid>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import Result from './components/Result.vue';
import Second from './components/Second.vue';
import Min from './components/Min.vue';
import Hour from './components/Hour.vue';
import Day from './components/Day.vue';
import Month from './components/Month.vue';
import Week from './components/Week.vue';

import { NGi, NCard, NGrid, NTabs, NTabPane, NTag } from "naive-ui";

const emit = defineEmits(['update:modelValue', 'change']);
const props = defineProps({
  modelValue: { type: String },
});
onMounted(() => {
  console.log("props:", props.modelValue)
});

interface CronGenValue {
  second: string;
  min: string;
  hour: string;
  day: string;
  month: string;
  week: string;
  // year: string;
}

const value = computed({
  get: () => props.modelValue,
  set: (v) => emit("update:modelValue", v),
});

// 初始值
const defaultValue: CronGenValue = {
  second: "0",
  min: "0",
  hour: "0",
  day: "*",
  month: "*",
  week: "?",
};

// 反解析 表达式
const getValue = () => {
  console.log("props:", props.modelValue)
  console.log("value:", value)
  if (value.value) {
    const arr = value.value.split(" ");
    if (arr.length >= 5) {
      //5 位以上是合法表达式
      const obj: CronGenValue = {
        second: arr[0],
        min: arr[1],
        hour: arr[2],
        day: arr[3],
        month: arr[4],
        week: arr[5],
      };
      return { ...obj };
    } else {
      return { ...defaultValue };
    }
  } else {
    return { ...defaultValue };
  }
};

const cronGenValue = ref<CronGenValue>(getValue());

const map: CronGenValue = {
  second: "秒",
  min: "分钟",
  hour: "小时",
  day: "日",
  month: "月",
  week: "周",
};

watch(
  () => cronGenValue.value.day,
  (v) => {
    if (v !== "?") {
      cronGenValue.value.week = "?";
    }
  }
);

watch(
  () => cronGenValue.value.week,
  (v) => {
    if (v !== "?") {
      cronGenValue.value.day = "?";
    }
  }
);

const cronGenValueString = computed(() => {
  const { second, min, hour, day, month, week } = cronGenValue.value;
  return `${second} ${min} ${hour} ${day} ${month} ${week}`;
});

watch(
  () => cronGenValueString.value,
  (v) => {
    console.log("11111111i:", v)
    emit("update:modelValue", v);
    emit('change', v)
  }
);
</script>
<style scoped lang="css">
.nan {
  width: 60px;
}
</style>