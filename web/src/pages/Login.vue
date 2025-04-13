<template>
  <div :class="['container', isMobile ? '' : 'pc']">
    <div class="form">
      <h1 class="title">Minas</h1>
      <n-form v-show="!state.mfaFormShow" :model="model" ref="form" :rules="rules" label-placement="left" @keydown.enter="login">
        <n-form-item path="username">
          <n-input round v-model:value="model.username" :placeholder="t('fields.login_name')" clearable>
            <template #prefix>
              <n-icon>
                <person-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-form-item path="password">
          <n-input round v-model:value="model.password" type="password" :placeholder="t('fields.password')" clearable>
            <template #prefix>
              <n-icon>
                <lock-closed-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-button round block type="primary" :disabled="state.loading" :loading="state.loading"
          @click.prevent="login">{{
    t('buttons.sign_in') }}</n-button>
      </n-form>
      <n-form v-show="state.mfaFormShow" :model="model" ref="mfaFormRef" :rules="mfaRules" label-placement="left"
        @keydown.enter="login">
        <n-form-item path="mfa_verify_code">
          <n-input round v-model:value="model.mfa_verify_code" :placeholder="t('fields.mfa_verify_code')" clearable>
            <template #prefix>
              <n-icon>
                <lock-closed-outline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        <n-button round block type="primary" :disabled="state.loading" :loading="state.loading"
          @click.prevent="login">{{
    t('buttons.sign_in') }}</n-button>
      </n-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NForm, NFormItem, NInput, NButton, NIcon, useMessage } from "naive-ui";
import { PersonOutline, LockClosedOutline } from "@vicons/ionicons5";
import userApi from "@/api/basic/user";
import systemApi from "@/api/system";
import { useStore } from "vuex";
import { useIsMobile } from "@/utils";
import { Mutations } from "@/store/mutations";
import { useForm, requiredRule } from "@/utils/form";
import { useI18n } from 'vue-i18n'
import { Md5 } from 'ts-md5';
import { xxtea } from '@/utils/xxtea'

const message = useMessage()
const { t } = useI18n()
const router = useRouter();
const route = useRoute();
const store = useStore();
const isMobile = useIsMobile()
const form = ref();
const model = reactive({} as any);
const rules = {
  username: requiredRule(),
  password: requiredRule(),
};
const mfaRules = {
  mfa_verify_code: requiredRule(),
};
const mfaFormRef = ref();
const mfaForm = ref(false);
const state = reactive({
  loading: false,
  mfaFormShow: false
} as any);
const login = () => {
  state.loading = true;
  const params = { username: model.username, password: xxtea.encryptAuto(model.password), login_type: "password" } as any;
  if (state.mfaFormShow) {
    params.mfa_verify_code = model.mfa_verify_code
    params.login_type = "mfa"
  }
  userApi.login(params).then((user: any) => {
    console.log("user:", user)
    if (user.code == 200 && user?.data?.token) {
      store.commit(Mutations.SetUser, user?.data);
      let redirect = decodeURIComponent(<string>route.query.redirect || "/");
      router.push({ path: redirect });
    } else if (user.code === 4401) {
      //需要mfa认证
      state.mfaFormShow = true
    }
    state.loading = false;
  }).catch((err: any) => {
    state.loading = false;
  })
}

async function checkState() {
  const r = await systemApi.checkState();
  if (r.data.fresh) {
    router.push("/init")
  }
}

onMounted(checkState);
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  border-radius: 5px;
  box-shadow: 1px 1px 10px #ddd;
  display: flex;
  justify-content: center;
  align-items: center;

  .form {
    flex: 60%;
    padding: 20px;

    .title {
      margin-top: -10px;
      text-align: center;
    }
  }
}

.pc {
  width: 300px;
  margin: calc(50vh - 200px) auto;
}
</style>