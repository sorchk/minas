import type { App } from 'vue';
import { wavesDirective, dragDirective, resizeDirective } from '@/directive/customDirective';
import { clickOutSideDirective } from "@/directive/clickOutSide";

/**
 * 导出指令方法：v-xxx
 * @methods authDirective 用户权限指令，用法：v-auth
 * @methods wavesDirective 按钮波浪指令，用法：v-waves
 * @methods dragDirective 自定义拖动指令，用法：v-drag
 */
export function directive(app: App) {
	// 按钮波浪指令
	wavesDirective(app);
	// 元素大小变化
	resizeDirective(app);
	// 自定义拖动指令
	dragDirective(app);
	//点击空白事件指令
	clickOutSideDirective(app);
}
