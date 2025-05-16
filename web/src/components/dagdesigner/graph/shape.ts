import { Graph, KeyValue, Node, Path } from "@antv/x6";
import { register } from '@antv/x6-vue-shape';
import consts from "./consts";
import { createVNode } from "vue";
import taskNode from "../shape/TaskNode.vue";
/**
 * 对象深克隆
 * @param obj 源对象
 * @returns 克隆后的对象
 */
export const deepClone = (obj: any) => {
    let newObj: any;
    try {
        newObj = obj.push ? [] : {};
    } catch (error) {
        newObj = {};
    }
    for (let attr in obj) {
        if (obj[attr] && typeof obj[attr] === 'object') {
            newObj[attr] = deepClone(obj[attr]);
        } else {
            newObj[attr] = obj[attr];
        }
    }
    return newObj;
}
export const dagEdgeEntity = {
    id: '',
    shape: 'dag-edge',
    attrs: {
        line: {
            stroke: '#1890ff',
            strokeWidth: 1,
            strokeDasharray: '5 5',
            targetMarker: {
                "name": "classic",
                "size": 6
            }
        }
    },
    defaultLabel: {
        markup: [
            {
                tagName: "rect",
                selector: "body"
            },
            {
                tagName: "text",
                selector: "label"
            }
        ],
        attrs: {
            label: {
                fontSize: 14,
                textAnchor: "middle",
                textVerticalAnchor: "middle",
                pointerEvents: "none",
                fill: "#2CFEFF"
            }
        },
        position: {
            distance: 0.5
        }
    },
    zIndex: -1
};
export const commonPortsGroups = {
    in: {
        position: "top",
        attrs: {
            circle: {
                r: 4,
                magnet: true,
                fill: '#fff',
                fillOpacity: "0.15",
                stroke: '#31bd01',
                strokeWidth: 2,
                visibility: "hidden"
            }
        }
    },
    out: {
        position: "bottom",
        attrs: {
            circle: {
                r: 4,
                magnet: true,
                fill: '#fff',
                fillOpacity: "0.15",
                stroke: '#3078fa',
                strokeWidth: 2,
                visibility: "hidden"
            },
        }
    }
};


Graph.registerEdge(
    'dag-edge',
    {
        inherit: 'edge',
        attrs: {
            line: {
                stroke: '#C2C8D5',
                strokeWidth: 1,
                targetMarker: {
                    "name": "classic",
                    "size": 6
                }
            }
        }
    },
    true
);
Graph.registerConnector(
    'algo-connector',
    (s, e) => {
        const offset = 2
        const deltaY = Math.abs(e.y - s.y);
        const control = Math.floor((deltaY / 3) * 2);

        const v1 = { x: s.x, y: s.y + offset + control };
        const v2 = { x: e.x, y: e.y - offset - control };

        return Path.normalize(
            `M ${s.x} ${s.y}
       L ${s.x} ${s.y + offset}
       C ${v1.x} ${v1.y} ${v2.x} ${v2.y} ${e.x} ${e.y - offset}
       L ${e.x} ${e.y}
      `);
    },
    true
);

export class NodeGroup extends Node {
    public collapsed = false;
    private expandSize: { width: number; height: number; }
        = { width: consts.groupDefaultWidth, height: consts.groupDefaultHeight };

    protected postprocess() {
        this.collapsed = this.prop('collapsed')
        if (this.prop('expandSize')) {
            this.expandSize = this.prop('expandSize');
        }
    }

    isCollapsed() {
        return this.collapsed;
    }
    setExpandSize(size: { width: number; height: number; }) {
        console.log(this.isCollapsed(), size)
        if (!this.isCollapsed) {
            this.expandSize = size;
            this.prop('expandSize', this.expandSize)
        }
    }
    toggleCollapse(collapsed?: boolean) {
        this.postprocess();
        const target = (collapsed == null || collapsed == undefined) ? !this.collapsed : collapsed;
        this.prop('collapsed', target)
        if (target) {
            this.attr('buttonSign', { d: 'M 1 5 9 5 M 5 1 5 9' });
            this.setExpandSize(this.getSize())
            this.resize(consts.nodeDefaultWidth, consts.nodeDefaultHeight);
        } else {
            this.attr('buttonSign', { d: 'M 2 5 8 5' });
            if (this.expandSize) {
                this.resize(this.expandSize.width, this.expandSize.height)
            }
        }
        this.collapsed = target;
    }
}


NodeGroup.config({
    shape: 'rect',
    markup: [
        {
            tagName: 'rect',
            selector: 'body',
        },
        {
            tagName: 'text',
            selector: 'text',
        },
        {
            tagName: 'g',
            selector: 'buttonGroup',
            children: [
                {
                    tagName: 'rect',
                    selector: 'button',
                    attrs: {
                        'pointer-events': 'visiblePainted',
                    },
                },
                {
                    tagName: 'path',
                    selector: 'buttonSign',
                    attrs: {
                        fill: 'none',
                        'pointer-events': 'none',
                    },
                },
            ],
        },
    ],
    attrs: {
        body: {
            refWidth: '100%',
            refHeight: '100%',
            strokeWidth: 1,
            fill: 'rgba(95,149,255,0.05)',
            stroke: '#5F95FF',
            rx: 8,
            ry: 8,
        },
        text: {
            fontSize: 12,
            fill: '#666',
            refX: 30,
            refY: 10,
        },
        buttonGroup: {
            refX: 8,
            refY: 8,
        },
        button: {
            height: 14,
            width: 16,
            rx: 2,
            ry: 2,
            fill: '#f5f5f5',
            stroke: '#ccc',
            cursor: 'pointer',
            event: 'node:collapse',
        },
        buttonSign: {
            refX: 3,
            refY: 2,
            stroke: '#808080',
        },
        label: {
            fontSize: 12,
            fill: '#fff',
            refX: 32,
            refY: 10,
        },
    },
});
Graph.registerNode('groupNode', NodeGroup);

// register({
//     shape: 'task-node',
//     width: 120,
//     height: 50,
//     component: taskNode,
// })
export function registerNode(procs: Array<any>) {
    for (let i in procs) {
        let item = procs[i];
        if (!item.group) {
            //跳过分组
            continue;
        }
        let entity: KeyValue<any> = {}
        let component = null;
        if (item.type === 'group') {
            entity = {
                inherit: "groupNode",
                width: (item.width) || consts.groupDefaultWidth,
                height: (item.height) || consts.groupDefaultHeight,
                attrs: {
                    text: { text: item.name || '' }
                },
                ports: null,
            };
        } else {
            entity = {
                inherit: item.inherit ? item.inherit : "vue-shape",
                width: (item.width) || consts.nodeDefaultWidth,
                height: (item.height) || consts.nodeDefaultHeight,
                attrs: {
                    label: { text: item.name || '', y: 0, x: 0 }
                },
                ports: null,
            };

            if (item.component) {
                component = item.component;
            } else {
                component = {
                    setup() {
                        return { data: item }
                    },
                    components: {
                        taskNode
                    },
                    render: () => {
                        return createVNode(taskNode, { "data": item });
                    }
                    //    , template: `
                    // <task-node :data="data" />`
                };
            }
            entity.attrs.label.text = '';
        }
        entity.data = {
            parent: item.type === 'group',
            form: {
                label: item.name
            }
        };
        if (item.ports) {
            entity.ports = item.ports;
        } else {
            entity.ports = {
                groups: commonPortsGroups,
                items: [{
                    group: "in"
                }, {
                    group: "out"
                }]
            };
        }
        //Graph.registerNode(item.id, entity, true);
        if (item.type === 'group') {
            Graph.registerNode(item.id, entity, true);
        } else {
            register({
                shape: item.id,
                width: (item.width) || consts.nodeDefaultWidth,
                height: (item.height) || consts.nodeDefaultHeight,
                ports: entity.ports,
                component: component,
                data: entity.data
            })
        }

    }
}
