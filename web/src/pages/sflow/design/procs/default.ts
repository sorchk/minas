import startNode from "../shape/StartNode.vue";
import endNode from "../shape/EndNode.vue";
import { createVNode } from "vue";

export const default_components = [
    {
        id: "start",
        name: "",
        description: "",
        inherit: 'vue-shape',
        width: 60,
        height: 60,
        component: {
            components: {
                startNode
            },
            render: () => {
                return createVNode(startNode);
            }
            // ,
            // template: `
            //   <start-node/>`
        },
        group: 'base',
        hidden: true,
        form: {},
        fields: [
            {
                prop: "logEnable", label: "日志", type: 'switch', active: "启用", inactive: "禁用", value: false
            },
            {
                prop: "params", label: "请求参数", type: "params", value: []
            },
            {
                prop: "headers", label: "请求头", type: "params", value: []
            }
        ],
        ports: {
            groups: {
                out: {
                    position: "bottom",
                    attrs: {
                        circle: {
                            r: 4,
                            magnet: true,
                            fill: '#fff',
                            fillOpacity: "0.15",
                            stroke: '#31bd01',
                            strokeWidth: 2,
                            style: {
                                visibility: "visible"
                            }
                        }
                    }
                }
            },
            items: [
                {
                    group: "out"
                }
            ]
        }
    },
    {
        id: 'end',
        inherit: 'vue-shape',
        width: 60,
        height: 60,
        component: {
            components: {
                endNode
            },
            render: () => {
                return createVNode(endNode);
            }
            // ,
            // template: `
            //   <end-node/>`
        },
        group: 'base',
        hidden: true,
        fields: [
            {
                prop: "logEnable", label: "日志", type: 'switch', active: "启用", inactive: "禁用", value: false
            }, {
                prop: "resultType",
                label: '结果类型',
                type: 'radio',
                data: [{
                    value: 'none',
                    label: "不限制"
                }, {
                    value: 'null',
                    label: "空(null)"
                }, {
                    value: 'list',
                    label: "集合"
                }, {
                    value: 'map',
                    label: "单记录"
                }, {
                    value: 'number',
                    label: "数字"
                }, {
                    value: 'string',
                    label: "字符串"
                }],
                value: 'null'
            },
        ],
        ports: {
            groups: {
                in: {
                    position: "top",
                    attrs: {
                        circle: {
                            r: 4,
                            magnet: true,
                            fill: '#fff',
                            fillOpacity: "0.15",
                            stroke: '#3078fa',
                            strokeWidth: 2,
                            style: {
                                visibility: "visible"
                            }
                        }
                    }
                }
            },
            items: [
                {
                    group: "in"
                }
            ]
        }
    }
];
