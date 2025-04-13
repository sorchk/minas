import { Icon } from "@iconify/vue";
import { NIcon } from "naive-ui";
import { h } from "vue";

export function renderIcon(icon: string, props?: import("naive-ui").IconProps) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) });
}
