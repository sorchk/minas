import DagDesigner from "./index.vue";
import { registerNode } from "./graph/shape";
import { App } from "vue";

import './index.scss';

const components = [
    DagDesigner,
]
const DagDesignerSDK = {
    registerNode,
}

declare module "@vue/runtime-core" {
    interface ComponentCustomProperties {
      $registerDagNode?: (components: any) => void;
    }
  }
const install = (app: App) => {
    app.config.globalProperties.$registerDagNode = registerNode;
    components.forEach((component:any) => {
        app.component(component.name, component)
    })
    // window.axios = axios
}
DagDesigner.install = install;
export default {
    install,
    DagDesigner,
    DagDesignerSDK
}
