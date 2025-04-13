<template>
  <div class="terminal" v-resize="onResize" v-loading="loading" ref="terminal" element-loading-text="拼命连接中"></div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, watch, reactive, onUnmounted } from "vue";
import { debounce } from 'lodash'
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import termApi from "@/api/term/index";
// import { AttachAddon } from "xterm-addon-attach";
import "xterm/css/xterm.css";
import { number } from "echarts";
onMounted(() => {
  console.log("open")
})
onUnmounted(() => {
  console.log("close")
})
const terminal = ref(null);
const props = defineProps({
  termId: { type: Number },
});
const fitAddon = new FitAddon();
// const attachAddon = new AttachAddon();

let first = ref(true);
let loading = ref(true);
let terminalSocket = ref();
let term = ref();


const runRealTerminal = () => {
  loading.value = false;
  terminalSocket.value.send(connMsg());
}

const onWSReceive = (message: any) => {
  // 首次接收消息,发送给后端，进行同步适配
  if (first.value === true) {
    first.value = false;
    resizeRemoteTerminal();
  }
  const index = message.data.indexOf(":");
  let type = parseInt(message.data.substring(index == -1 ? 0 : index));
  let data = "";
  if (index != -1) {
    data = message.data.substring(index + 1);
    type = parseInt(message.data.substring(0, index));
  } else {
    type = parseInt(message.data.substring(0));
  }
  switch (type) {
    case 0:
      //unknown
      break;
    case 1:
      //conn
      break;
    case 2:
      //close
      closeRealTerminal();
      break;
    case 3:
      //msg
      term.value.element && term.value.focus()
      term.value.write(data);

      // const msgs = data.split(",");
      // msgs.forEach((item: any) => {
      //   const msg = parseInt(item);
      //   term.value.write(String.fromCharCode(msg));
      // })

      break;
    case 4:
      // resize
      break;
    case 5:
      //error
      const error = data;
      errorRealTerminal({ message: error });
      break;
  }
}
const connMsg = () => {
  const data = state.token;
  return "1:" + data;
}
const closeMsg = () => {
  return "2:";
}
const dataMsg = (data: string) => {
  return "3:" + data;
}
const resizeMsg = (cols: number, rows: number) => {
  return "4:" + cols + "," + rows;
}
const heartbeatMsg = () => {
  return "0:" + new Date().getTime();
}

const errorRealTerminal = (ex: any) => {
  let message = ex.message;
  console.log("message:", message)
  if (!message) message = 'disconnected'
  term.value.write(`\x1b[31m${message}\x1b[m\r\n`)
}
const closeRealTerminal = () => {
  console.log("close");
  term.value.write(`\x1b[31mclose\x1b[m\r\n`)
  terminalSocket.value && terminalSocket.value.send(closeMsg());
}
const state = reactive({
  clientId: "",
  token: "",
})
const heartbeatTimer = ref();
const baseURL = import.meta.env.VITE_API_URL;
const createWS = async () => {
  const tokenData = (await termApi.auth(props.termId + "")).data as any;
  state.clientId = tokenData.clientId;
  state.token = tokenData.token;
  const protocol = window.location.protocol;
  const host = window.location.hostname;
  const port = window.location.port;
  const url = (protocol == 'https:' ? 'wss' : 'ws') + '://' + host + ':' + port + baseURL + '/term/ws/' + state.clientId + '/' + state.token;
  terminalSocket.value = new WebSocket(url);
  terminalSocket.value.onopen = runRealTerminal;
  terminalSocket.value.onmessage = onWSReceive;
  terminalSocket.value.onclose = closeRealTerminal;
  terminalSocket.value.onerror = errorRealTerminal;

  heartbeatTimer.value = setInterval(() => {
    if (isWsOpen()) {
      terminalSocket.value!.send(heartbeatMsg());
    }
  }, 1000 * 10);
}
const initWS = async () => {
  if (!terminalSocket.value) {
    await createWS();
  }
  if (terminalSocket.value && terminalSocket.value.readyState > 1) {
    terminalSocket.value.close();
    await createWS();
  }
}
// 发送给后端,调整后端终端大小,和前端保持一致,不然前端只是范围变大了,命令还是会换行
const resizeRemoteTerminal = () => {
  const { cols, rows } = term.value
  if (isWsOpen()) {
    terminalSocket.value.send(resizeMsg(cols, rows));
  }
}
const initTerm = () => {
  term.value = new Terminal({
    lineHeight: 1.2,
    fontSize: 14,
    fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
    theme: {
      background: '#181d28',
    },
    // 光标闪烁
    cursorBlink: true,
    cursorStyle: 'underline',
    scrollback: 1000,
    tabStopWidth: 4,
  });
  term.value.open(terminal.value);
  term.value.loadAddon(fitAddon);
  // term.value.loadAddon(attachAddon);
  // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
  setTimeout(() => {
    fitAddon.fit();
  }, 5);
}
// 是否连接中0 1 2 3
const isWsOpen = () => {
  const readyState = terminalSocket.value && terminalSocket.value.readyState;
  return readyState === 1
}
const fitTerm = () => {
  fitAddon.fit();
  resizeRemoteTerminal();
}
const onResize = debounce(() => fitTerm(), 800);


const termData = () => {
  // 输入与粘贴的情况,onData不能重复绑定,不然会发送多次
  term.value.onData((data: any) => {
    if (isWsOpen()) {
      // JSON.stringify({"operate": "command", "command": data})
      terminalSocket.value.send(dataMsg(data));
    }
  });
}





const open = async () => {
  console.log("initWS.", props)
  await initWS();
  initTerm();
  console.log("initTerm.")
  termData();
  console.log("termData.")
}
const close = () => {
  console.log("closeTerminal.")
  terminalSocket.value && terminalSocket.value.close();
}
onMounted(() => {
  open();
})
onBeforeUnmount(() => {
  close();
})
// 暴露变量
defineExpose({
  open, close,
});
</script>
<style lang="scss" scoped>
.terminal {
  padding: 0;
  margin: 0;
  width: 100%;
  height: 100%;
}
</style>
