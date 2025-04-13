<template>
  <x-page-header />
  <n-space class="page-body" vertical :size="12">
    <x-panel :title="t('fields.profile')" :subtitle="t('tips.profile')" divider="bottom"
      :collapsed="panel !== 'profile'">
      <template #action>
        <n-button secondary strong class="toggle" size="small" @click="togglePanel('profile')">{{ panel === 'profile' ?
    t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">
        <n-form inline :model="profile" :rules="profileRules" ref="profileForm">
          <n-grid cols="1 640:2" :x-gap="24">
            <n-form-item-gi :label="t('fields.username')" path="name">
              <n-input :placeholder="t('fields.username')" v-model:value="profile.name" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.login_name')" path="account">
              <n-input :placeholder="t('fields.login_name')" v-model:value="profile.account" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.email')" path="email">
              <n-input :placeholder="t('fields.email')" v-model:value="profile.email" />
            </n-form-item-gi>
          </n-grid>
        </n-form>
        <n-button type="primary" :disabled="profileSubmiting" :loading="profileSubmiting" @click="modifyProfile">
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.update') }}
        </n-button>
      </div>
    </x-panel>
    <x-panel :title="t('fields.password')" :subtitle="t('tips.password')" divider="bottom"
      :collapsed="panel !== 'password'">
      <template #action>
        <n-button secondary strong size="small" class="toggle" @click="togglePanel('password')">{{ panel === 'password'
    ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">
        <n-form :model="password" ref="passwordForm" :rules="passwordRules">
          <n-grid cols="1 640:3" :x-gap="24">
            <n-form-item-gi path="old" :label="t('fields.password_old')">
              <n-input v-model:value="password.old" type="password" :placeholder="t('fields.password_old')" />
            </n-form-item-gi>
            <n-form-item-gi first path="new" :label="t('fields.password_new')">
              <n-input v-model:value="password.new" type="password" :placeholder="t('fields.password_new')" />
            </n-form-item-gi>
            <n-form-item-gi first path="confirm" :label="t('fields.password_confirm')">
              <n-input :disabled="!password.new" v-model:value="password.confirm" type="password"
                :placeholder="t('fields.password_confirm')" />
            </n-form-item-gi>
          </n-grid>
        </n-form>
        <n-button type="primary" :disabled="passwordSubmiting" :loading="passwordSubmiting" @click="modifyPassword">
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.update') }}
        </n-button>
      </div>
    </x-panel>

    <x-panel :title="t('fields.mfa')" :subtitle="t('tips.mfa')" divider="bottom" :collapsed="panel !== 'mfa'">
      <template #action>
        <n-button secondary strong size="small" class="toggle" @click="togglePanel('mfa')">{{ panel === 'mfa' ?
    t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">

        <n-popconfirm v-if="!showMfaForm" @positive-click="handlePositiveClick" @negative-click="handleNegativeClick">
          <template #trigger>
            <n-switch v-model:value="profile.mfa_enable" :checked-value="1" :unchecked-value="0">
              <template #checked>
                {{ t('enums.mfa_enabled') }}
              </template>
              <template #unchecked>
                {{ t('enums.mfa_disabled') }}
              </template>
            </n-switch>
          </template>
          {{ profile.mfa_enable == 0 ? t('prompts.mfa_disable') : t('prompts.mfa_enable') }}
        </n-popconfirm>

        <n-form style="margin-top: 12px" v-if="showMfaForm" :model="mfa" ref="mfaForm" :rules="mfaRules">
          <n-grid cols="1 640:1" :x-gap="24">
            <n-form-item-gi :label="t('tips.mfa_tool')">
              <ul>
                <li>Google Authenticator</li>
                <li>Microsoft Authenticator</li>
                <li>Yubico Authenticator</li>
                <li>Bitwarden</li>
                <li>1Password</li>
                <li>LastPass</li>
                <li>Authenticator</li>
              </ul>
            </n-form-item-gi>
            <n-form-item-gi v-if="mfa.otpauth" :label="t('tips.mfa_code')">
              <n-qr-code :value="mfa.otpauth" />
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.mfa_code')">
              {{ mfa.mfa_code }}
            </n-form-item-gi>
            <n-form-item-gi :label="t('fields.title')">
              <n-input v-model:value="mfa.mfa_title" type="text" :placeholder="t('fields.title')" />
            </n-form-item-gi>
            <n-form-item-gi path="mfa_verify_code" :label="t('fields.mfa_verify_code')">
              <n-input v-model:value="mfa.mfa_verify_code" type="text" :placeholder="t('fields.mfa_verify_code')" />
            </n-form-item-gi>
          </n-grid>
        </n-form>
        <n-space v-if="showMfaForm">
          <n-button class="mr10" type="primary" @click="cancelSaveMfa">
            <template #icon>
              <n-icon>
                <close-icon />
              </n-icon>
            </template>
            {{ t('buttons.cancel') }}
          </n-button>

          <n-button type="primary" @click="saveMfaCode">
            <template #icon>
              <n-icon>
                <save-icon />
              </n-icon>
            </template>
            {{ t('buttons.save') }}
          </n-button>
        </n-space>
      </div>
    </x-panel>
    <x-panel :title="t('fields.preferences')" :subtitle="t('tips.preference')" :collapsed="panel !== 'preference'">
      <template #action>
        <n-button secondary strong class="toggle" size="small" @click="togglePanel('preference')">{{ panel ===
    'preference' ? t('buttons.collapse') : t('buttons.expand') }}</n-button>
      </template>
      <div style="padding: 4px 0 0 12px">
        <n-form inline :model="preference" ref="preferenceForm" label-placement="left">
          <n-form-item :label="t('fields.language')" path="locale">
            <n-radio-group v-model:value="preference.locale">
              <n-radio-button value="zh">中文</n-radio-button>
              <n-radio-button value="en">English</n-radio-button>
            </n-radio-group>
          </n-form-item>
          <n-form-item :label="t('fields.theme')" path="theme">
            <n-radio-group v-model:value="preference.theme">
              <n-radio-button value="light">{{ t('enums.light') }}</n-radio-button>
              <n-radio-button value="dark">{{ t('enums.dark') }}</n-radio-button>
            </n-radio-group>
          </n-form-item>
        </n-form>
        <n-button type="primary" @click="savePreference">
          <template #icon>
            <n-icon>
              <save-icon />
            </n-icon>
          </template>
          {{ t('buttons.save') }}
        </n-button>
      </div>
    </x-panel>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NPopconfirm,
  NQrCode,
  NSwitch,
  NIcon,
  NForm,
  NFormItem,
  NFormItemGi,
  NGrid,
  NRadioButton,
  NRadioGroup,
} from "naive-ui";
import {
  SaveOutline as SaveIcon,
  CloseOutline as CloseIcon,
} from "@vicons/ionicons5";
import XPageHeader from "@/components/PageHeader.vue";
import XPanel from "@/components/Panel.vue";
import userApi from "@/api/basic/user";
import { useForm, emailRule, requiredRule, customRule, lengthRule } from "@/utils/form";
import { Mutations } from "@/store/mutations";
import { useStore } from "vuex";
import { useI18n } from 'vue-i18n'
import { Base32 } from "@/utils/base32";
import bcrypt from 'bcryptjs'
import { xxtea } from "@/utils/xxtea";

const { t } = useI18n()
const panel = ref('')
function togglePanel(name: string) {
  if (panel.value === name) {
    panel.value = ''
  } else {
    panel.value = name
  }
}

// profile
const profile = ref({} as any)
const profileRules: any = {
  name: requiredRule(),
  loginName: requiredRule(),
  email: [requiredRule(), emailRule()],
};
const profileForm = ref();
const { submit: modifyProfile, submiting: profileSubmiting } = useForm(profileForm, () => userApi.modifyProfile(profile.value))

// password
const password = reactive({
  showDlg: false,
  old: "",
  new: "",
  confirm: "",
})
const passwordForm = ref()
const passwordRules = {
  old: requiredRule(),
  new: [
    requiredRule(),
    lengthRule(6, 15),
    customRule((_: any, value: string) => value !== password.old, t('tips.password_new_rule')),
  ],
  confirm: [
    requiredRule(),
    customRule((_: any, value: string) => value === password.new, t('tips.password_confirm_rule')),
  ],
};
const { submit: modifyPassword, submiting: passwordSubmiting } = useForm(
  passwordForm,
  () => userApi.modifyPassword({ password: xxtea.encryptAuto(password.old), newpassword: bcrypt.hashSync(password.new, 10) }),
  () => {
    password.showDlg = false
    window.message.info(t('texts.action_success'));
  }
);

const mfaForm = ref()
const mfa = reactive({
  mfa_enable: 1,
  mfa_code: "",
  mfa_title: "",
  mfa_verify_code: "",
  otpauth: ""
})
const mfaRules = {
  mfa_verify_code: [
    requiredRule()
  ],
}
const cancelSaveMfa = () => {
  showMfaForm.value = false;
  profile.value.mfa_enable = 0
}
const saveMfaCode = () => {
  //保持开启多因素认证状态
  userApi.enableMfa(mfa).then((res) => {
    if (res.code == 200) {
      profile.value.mfa_enable = 1
      showMfaForm.value = false;
    } else {
      profile.value.mfa_enable = 0
    }
  })

}
const cannelMfaCode = () => {
  userApi.disableMfa().then((res) => {
    if (res.code == 200) {
      profile.value.mfa_enable = 0
      showMfaForm.value = false;
    } else {
      profile.value.mfa_enable = 1
    }
  })
}
// preference
const store = useStore();
const preference = reactive({
  locale: store.state.preference.locale,
  theme: store.state.preference.theme || 'light',
})
const preferenceForm = ref()
function savePreference() {
  store.commit(Mutations.SetPreference, preference)
  window.message.info(t('texts.action_success'));
  setTimeout(() => location.reload(), 100)
}
const handlePositiveClick = () => {
  handleUpdateMfa(profile.value.mfa_enable)
}
const handleNegativeClick = () => {
  console.log(profile.value.mfa_enable)
  if (profile.value.mfa_enable == 0) {
    profile.value.mfa_enable = 1
  } else {
    profile.value.mfa_enable = 0
  }
}
const base32 = new Base32({ padding: false });
const showMfaForm = ref(false)
const generateMfaCode = () => {
  var chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
  var result = '';
  var charactersLength = chars.length;
  for (var i = 0; i < 15; i++) {
    result += chars.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
const genMfa = () => {
  mfa.mfa_enable = 1
  if (!mfa.mfa_title) {
    mfa.mfa_title = "Minas"
  }
  mfa.mfa_code = base32.encode(generateMfaCode()) + "";
  mfa.otpauth = "otpauth://totp/" + mfa.mfa_title + ":" + profile.value.account + "?secret=" + mfa.mfa_code + "&issuer=Minas"
  showMfaForm.value = true
}
const handleUpdateMfa = (value: number) => {
  if (value == 0) {
    //取消多因素认证
    cannelMfaCode();
  } else {
    //打开多因素认证表单
    genMfa();
  }
}
async function fetchData() {
  let r = await userApi.profile();
  profile.value = r.data as any;
}

onMounted(fetchData);
</script>

<style scoped>
.toggle {
  width: 75px;
}
</style>