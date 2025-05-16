
import { cloneDeep } from "lodash";
import { script_components } from "./script";
import { system_components } from "./system";
import { notice_components } from "./notice";
import { output_components } from "./dataOutput";
import { input_components } from "./dataInput";
import { check_components } from "./check";
import { calc_components } from "./calc";
import { default_components } from "./default";
import { program_components } from "./proram";
import { data_components } from "./data";
import { crawler_components } from "./crawler";

/**
 * 节点通用字段
 * @type {[{prop: string, label: string, placeholder: string, value: string},{prop: string, label: string, placeholder: string, value: string},{prop: string, label: string, placeholder: string, value: string},{inactive: string, prop: string, active: string, label: string, type: string, value: boolean},{prop: string, label: string, placeholder: string, value: string}]}
 */
const nodeFields = [
    { prop: "label", label: "标签", value: "", placeholder: "请输入标签名称" },
    { prop: "datakey", label: "数据KEY", value: "", placeholder: "允许空,修改后不能使用task_id取数据,只能用设置的值取数据" },
    {
        prop: "cacheTime",
        label: "缓存时间",
        value: "",
        placeholder: "单位秒，如果空或小于1则不启用缓存。 10800=3小时，86400=24小时，604800=7天，2592000=30天"
    },
    {
        prop: "enabled", label: "是否启用", type: 'switch', active: "启用", inactive: "禁用", value: true
    },
    {
        prop: "logEnable", label: "日志", type: 'switch', active: "启用", inactive: "禁用", value: false
    },
    {
        prop: "ignoreSimpleException",
        label: "忽略异常",
        value: "",
        placeholder: "忽略处理器异常(SimpleException)，这里填写el表达式，需要结果是布尔值",
        help: `变量说明：\${taskInput}本节点以前的数据（Map）、\${message}异常消息（文本）、
        \${className}异常类全名（文本）、
        \${name}异常类名（文本）、\${stackTrace}异常栈（数组）
        常用表达式：包含\${fn:containsIgnoreCase(message,'subtext')}
        、contains、开头\${fn:startsWith(message, 'prefix')}`
    }
];
/**
 * 添加节点字段
 * @param item
 */
export function addFields(node: any, fields: any) {
    if (node.group && !node.hidden) {
        node.fields = cloneDeep(fields).concat(node.fields);
    }
}

/**
 * 连线通用字段
 * @type {[{prop: string, label: string, placeholder: string, value: string},{prop: string, label: string, placeholder: string, value: string}]}
 */
export const edgeFields = [
    { prop: "label", label: "名称", placeholder: "请输入任务名称", value: "" },
    { prop: "expr", label: "EL表达式", value: "", placeholder: "请输入EL，真则向下执行，反之跳过后面的任务节点" }
];

// dag-node节点链接桩群组的配置数据
export const componentGroups = [
    {
        id: "base",
        name: "基础组件",
        icon: "fa fa-folder-o",
        remark: "循环，遍历，迭代，子程序",
        items: new Array()
    },
    {
        id: "script",
        name: "脚本",
        remark: "javascript，python，el，mvel",
        items: new Array()
    },
    {
        id: "check",
        name: "校验",
        remark: "布尔验证，正则验证，常用规则，规则引擎",
        items: new Array()
    },
    {
        id: "input",
        name: "输入",
        remark: "参数，文件，数据库，接口，消息队列",
        items: new Array()
    },
    {
        id: "output",
        name: "输出",
        remark: "json输出，文件，数据库，接口，消息队列",
        items: new Array()
    },
    {
        id: "crawler",
        name: "爬虫",
        remark: "静态爬虫，动态爬虫，变量爬虫，文件下载，休眠",
        items: []
    },
    {
        id: "data",
        name: "数据处理",
        remark: "行转列，排序，数据过滤，字段筛选，list转map，map转tree，数据格式化",
        items: new Array()
    },
    {
        id: "calc",
        name: "计算",
        remark: "求和，平均数，最大，最小，中值，百分比",
        items: new Array()
    },
    {
        id: "notice",
        name: "通知",
        remark: "接口，消息队列，邮件，短信",
        items: new Array()
    },
    {
        id: "system",
        name: "系统",
        remark: "系统，文件，网络，其他",
        items: new Array()
    }
];
export const components =  new Array();
components.push(...default_components);
// components.push(...test_components);
components.push(...script_components);
components.push(...input_components);
components.push(...output_components);
components.push(...notice_components);
components.push(...system_components);
components.push(...check_components);
components.push(...calc_components);
components.push(...program_components);
components.push(...data_components);
components.push(...crawler_components);

//修复组件数据
for (let j in components) {
    const node = components[j];
    if (node.hidden == null || node.hidden == undefined) {
        node.hidden = false;
    }
    addFields(node, nodeFields)
}
//组件分组
for (let i in componentGroups) {
    const item = componentGroups[i];
    item.items = new Array();
    for (let j in components) {
        const node = components[j];
        if (node.group == item.id) {
            item.items.push(node);
        }
    }
}