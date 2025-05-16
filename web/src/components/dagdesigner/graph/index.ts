import {
    Cell,
    EdgeView,
    FunctionExt,
    Graph,
    Node,
    NodeView,
    Shape
} from "@antv/x6";
import { Clipboard } from '@antv/x6-plugin-clipboard';
import { Dnd } from '@antv/x6-plugin-dnd';
import { History } from '@antv/x6-plugin-history';
import { Keyboard } from '@antv/x6-plugin-keyboard';
import { Scroller } from '@antv/x6-plugin-scroller';
import { Selection } from '@antv/x6-plugin-selection';
import { Snapline } from '@antv/x6-plugin-snapline';
import { Transform } from '@antv/x6-plugin-transform';
import '@antv/x6-vue-shape';
import { Options } from "@antv/x6/lib/graph/options";
import { cloneDeep } from "lodash";
import consts from "./consts";
import { NodeGroup, dagEdgeEntity } from "./shape";
export class DagGraph {
    public container: HTMLElement;
    public graph: Graph;
    public idIndex: number = 1;
    public maxIndex: number = 999999;
    public initNodes: Array<any> = new Array();
    public selectCell: Cell | null = null;
    public dnd: Dnd;
    public newId = () => {
        this.idIndex++;
        return this.idIndex;
    }
    public getGraph(): Graph {
        return this.graph;
    }
    public getDnd(): Dnd {
        return this.dnd;
    }
    constructor(container: HTMLElement, initNodes?: Array<any>) {
        this.container = container;
        this.graph = createGraph(this);
        if (initNodes) {
            this.initNodes = initNodes;
        }
        const that = this;
        this.dnd = new Dnd({
            target: this.graph,
            scaled: true,
            getDropNode: function (node) {
                const newId = that.newId() + "";
                const newNode = that.graph.createNode({
                    id: newId,
                    shape: node.shape,
                });
                return newNode;
            }
        })
        initEvent(this);
        initAnimate(this);

    }
    // 销毁
    public destroy() {
        this.graph?.dispose();
    }
    public initGraphStart() {
        this.idIndex = 1;
        this.maxIndex = 999999;
        const graph = this.graph;
        if (!graph) {
            return;
        }
        //初始化索引
        graph.clearCells();
        console.log("initNodes:", this.initNodes)
        for (let i in this.initNodes) {
            let node = this.initNodes[i];
            if (node.x && typeof (node.x) == 'function') {
                node.x = node.x({ width: graph.getGraphArea().width, height: graph.getGraphArea().height });
            }
            if (node.y && typeof (node.y) == 'function') {
                node.y = node.y({ width: graph.getGraphArea().width, height: graph.getGraphArea().height });
            }
            graph.addNode(node);
        }
        graph.zoomTo(1); // 缩放比例
        graph.centerContent(); // 内容居中
    }
    public initGraphShape(data: any) {
        let graph = this.graph;
        const cells: any[] = [];
        data.cells.forEach((item: any) => {
            if (item.shape === 'dag-edge') {
                cells.push(graph.createEdge(item));
            } else {
                delete item.component;
                cells.push(graph.createNode(item));
            }
            //更新id最大索引号
            let id = item.id;
            let isNumber = !isNaN(id);
            if (isNumber) {
                let num = parseInt(id);
                if (this.idIndex < num && num < this.maxIndex) {
                    this.idIndex = num;
                }
            }
        });
        graph.resetCells(cells);
        //居中显示
        const num = 1 - graph.zoom();
        num > 1 ? graph.zoom(num * -1) : graph.zoom(num);
        graph.centerContent();
    }
}
const showAllPorts = (dagGraph: DagGraph, show?: boolean) => {
    const ports = dagGraph.container?.querySelectorAll('.x6-port-body');
    showPorts(ports, show);
}
const showPorts = (ports?: NodeListOf<Element>, show?: boolean) => {
    if (ports) {
        for (let i = 0, len = ports.length; i < len; i = i + 1) {
            const port: Element = ports.item(i);
            port.setAttribute("visibility", show ? "visible" : "hidden");
            port.setAttribute("style", "visibility:" + show ? "visible" : "hidden");
        }
    }
}

const initEvent = (dagGraph: DagGraph) => {
    const graph = dagGraph.graph;
    if (!graph) {
        return;
    }

    // graph.on("cell:added", (node: any) => {
    //     console.log("--added--------id------", node.cell.id, ((node.cell.id + "").length > 10))
    //     if ((node.cell.id + "").length > 10) {
    //         const id = dagGraph.newId();
    //         console.log("--added--------id------", id)
    //         node.cell.setProp("id", id)
    //         node.cell.id = id
    //     }
    //     console.log("--added--------------", node)
    // });
    // 节点鼠标移入
    graph.on('node:mouseenter', FunctionExt.debounce(() => {
        // 显示连接点
        showAllPorts(dagGraph, true);
    }),
        500
    );
    // 节点鼠标移出
    graph.on('node:mouseleave', () => {
        // 隐藏连接点 
        showAllPorts(dagGraph, false);
    });

    // graph.on('node:removed', ({node}) => {
    //     const data = node.store.data;
    //     if (data.type === 'taskNode') {
    //         const posIndex = this.nodeData.findIndex((item) => item.id === data.id);
    //         this.nodeData.splice(posIndex, 1);
    //     }
    // });
    graph.on('selection:changed', (args) => {
        args.added.forEach(cell => {
            dagGraph.selectCell = cell;
        });
    });

    graph.on('node:selected', ({ node }) => {
        const zIndex = (node.zIndex || 1000);
        const children = node.getChildren();
        if (children) {
            children.forEach((item: Cell) => {
                const mIndex = (item.zIndex || 1000);
                if (mIndex < zIndex) {
                    item.setZIndex(zIndex + 1);
                } else {
                    item.setZIndex(mIndex);
                }
            })
        }
        const edges = graph.getIncomingEdges(node)
        edges?.forEach((edge) => {
            edge.attr('line/strokeDasharray', 5)
            edge.attr('line/stroke', '#1890ff')
            edge.attr('line/style/animation', 'running-line 30s infinite linear')
            edge.zIndex = 3001;
        })
        const outEdges = graph.getOutgoingEdges(node)
        outEdges?.forEach((edge) => {
            edge.attr('line/strokeDasharray', 5)
            edge.attr('line/stroke', '#de6720')
            edge.attr('line/style/animation', 'running-line 30s infinite linear')
            edge.zIndex = 3001;
        })
    })
    graph.on('node:unselected', ({ node }) => {
        const edges = graph.getIncomingEdges(node)
        edges?.forEach((edge) => {
            edge.removeAttrByPath('line/strokeDasharray')
            edge.attr('line/stroke', '#A2B1C3')
            edge.removeAttrByPath('line/style/animation')
            edge.zIndex = 3000;
        })
        const outEdges = graph.getOutgoingEdges(node)
        outEdges?.forEach((edge) => {
            edge.removeAttrByPath('line/strokeDasharray')
            edge.attr('line/stroke', '#A2B1C3')
            edge.removeAttrByPath('line/style/animation')
            edge.zIndex = 3000;
        })
    })

    graph.on('node:collapse', ({ node, e }: any) => {
        node.toggleCollapse();
        const collapsed = node.isCollapsed();
        const cells = node.getDescendants();
        cells.forEach((n: any) => {
            if (collapsed) {
                n.hide();
            } else {
                n.show();
            }
        });
        e.stopPropagation();
    });
    //节点移动 父子状态
    let ctrlPressed = false;
    graph.on('node:embedding', ({ e }: { e: any }) => {
        ctrlPressed = e.metaKey || e.ctrlKey;
    });
    graph.on('node:embedded', () => {
        ctrlPressed = false;
    });
    graph.on('node:change:size', ({ node, options }) => {
        if (options.skipParentHandler) {
            return;
        }
        const children = node.getChildren();
        if (children && children.length) {
            const group = node as NodeGroup;
            group.setExpandSize(group.getSize())
        }
    });
    graph.on('node:change:position', ({ node, options }) => {
        if (options.skipParentHandler || ctrlPressed) {
            return;
        }
        const children = node.getChildren();
        if (children && children.length) {
            node.prop('originPosition', node.getPosition());
        }
        const parent = node.getParent();
        if (parent && parent.isNode()) {
            let originSize = parent.prop('originSize');
            if (originSize == null) {
                originSize = parent.getSize();
                parent.prop('originSize', originSize);
            }
            let originPosition = parent.prop('originPosition');
            if (originPosition == null) {
                originPosition = parent.getPosition();
                parent.prop('originPosition', originPosition);
            }
            let x = originPosition.x;
            let y = originPosition.y;
            let cornerX = originPosition.x + originSize.width;
            let cornerY = originPosition.y + originSize.height;
            let hasChange = false;


            const children = parent.getChildren();
            if (children) {
                children.forEach((child) => {
                    const bbox = child.getBBox().inflate(20);
                    const corner = bbox.getCorner();
                    if (bbox.x < x) {
                        x = bbox.x;
                        hasChange = true;
                    }
                    if (bbox.y < y) {
                        y = bbox.y;
                        hasChange = true;
                    }
                    if (corner.x > cornerX) {
                        cornerX = corner.x;
                        hasChange = true;
                    }
                    if (corner.y > cornerY) {
                        cornerY = corner.y;
                        hasChange = true;
                    }
                });
            }
            if (hasChange) {
                parent.prop(
                    {
                        position: { x, y },
                        size: { width: cornerX - x, height: cornerY - y },
                    },
                    { skipParentHandler: true },
                );

                const group = parent as NodeGroup;
                group.setExpandSize({ width: cornerX - x, height: cornerY - y })
            }
        }
    });



    graph.on('edge:connected', ({ edge }) => {
        edge.attr({
            line: {
                strokeDasharray: ''
            }
        })
    });
    graph.on('edge:selected', ({ cell }) => {
        // 显示连接点
        showAllPorts(dagGraph, true);
        cell.addTools([
            {
                name: 'source-arrowhead',
                args: {
                    attrs: {
                        fill: 'green',
                    },
                },
            },
            {
                name: 'target-arrowhead',
                args: {
                    attrs: {
                        fill: 'red',
                    },
                },
            },
        ])
    });
    graph.on('edge:unselected', ({ cell }) => {
        // 隐藏连接点
        showAllPorts(dagGraph, false);
        cell.removeTools()
    });

}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const initAnimate = (dagGraph: DagGraph) => {
    const graph = dagGraph.graph;
    if (!graph) {
        return;
    }
    let flash = (cell: Cell) => {
        const cellView = graph.findViewByCell(cell);
        if (cellView) {
            cellView.highlight();
            setTimeout(() => cellView.unhighlight(), 350);
        }
    };
    graph.on("signal", (cell: Cell) => {
        if (cell.isEdge()) {
            const view = graph.findViewByCell(cell) as EdgeView;
            if (view) {
                // const token = Vector.create("circle", {
                //     r: 6,
                //     fill: "#feb662"
                // });
                // const target = cell.getTargetCell();
                // setTimeout(() => {
                //     view.sendToken(token.node, 1200, () => {
                //         if (target) {
                //             graph.trigger("signal", target);
                //         }
                //     });
                // }, 350);
            }
        } else {
            flash(cell);
            const edges = graph.model.getConnectedEdges(cell, {
                outgoing: true
            });
            edges.forEach((edge) => graph.trigger("signal", edge));
        }
    });
    graph.on("node:mousedown", ({ cell }) => {
        graph.trigger("signal", cell);
    });
}




const createGraph = (dagGraph: DagGraph): Graph => {
    const graph = new Graph({
        container: dagGraph.container,
        autoResize: true,
        grid: false,
        interacting: {
            nodeMovable: true,
            edgeMovable: true
        },
        // 设置滚轮缩放画布 
        mousewheel: {
            enabled: true,
            zoomAtMousePosition: true,
            modifiers: ['ctrl', 'meta'],
            minScale: 0.5,
            maxScale: 3,
        },
        // 配置全局连线规则
        connecting: connectingConfig(dagGraph),
        highlighting: {
            magnetAvailable: {
                name: 'stroke',
                args: {
                    padding: 8,
                    attrs: {
                        strokeWidth: 8,
                        stroke: '#52c41a',
                    },
                },
            },
            default: {
                name: "opacity",
                args: {
                    padding: 5,
                    attrs: {
                        "stroke-width": 2,
                        stroke: "#F00" //rgba(223,234,255)  #e52e1a
                    }
                }
            },
        },
        // 节点拖拽到另一节点，形成父子节点关系
        embedding: {
            enabled: true,
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            findParent: function ({ node, view }: {
                node: Node;
                view: NodeView;
            }): Array<any> {
                const bbox = node.getBBox();
                return this.getNodes().filter((cell: Cell) => {
                    // 只有 data.parent 为 true 的节点才是父节点
                    const data = cell.getData();
                    if (data && data.parent) {
                        const targetBBox = cell.getBBox();
                        return bbox.isIntersectWithRect(targetBBox);
                    }
                    return false;
                }) || new Array();
            }
        },
    });
    // graph.addNode = (node, options = {}) => {
    //     console.log("addNode:", node, options)
    //     return graph.model.addNode(node, options);
    // }
    // graph.createNode = (metadata: any) => {
    //     console.log("createNode:", metadata)
    //     const newNode = graph.model.createNode(metadata);

    //     console.log("newNode:", newNode)

    //     return newNode;
    // };
    // #region 使用插件
    graph
        .use(
            new Transform({
                resizing: {
                    enabled: (node: Node) => {
                        if (node.shape == 'start' || node.shape == 'end') {
                            return false;
                        } else if (node.data?.parent) {
                            const group = node as NodeGroup;
                            console.log("group.isCollapsed():", group.isCollapsed())
                            if (group.isCollapsed()) {
                                return false;
                            } else {
                                return true;
                            }
                        } else {
                            return false;
                        }
                    },
                    minWidth: (node: Node) => {
                        if (node.data?.parent) {
                            return consts.groupMinWidth;
                        } else {
                            return consts.nodeMinWidth;
                        }
                    },
                    minHeight: (node: Node) => {
                        if (node.data?.parent) {
                            return consts.groupMinHeight;
                        } else {
                            return consts.nodeMinHeight;
                        }
                    },
                    maxWidth: (node: Node) => {
                        if (node.data?.parent) {
                            return consts.groupMaxWidth;
                        } else {
                            return consts.nodeMaxWidth;
                        }
                    },
                    maxHeight: (node: Node) => {
                        if (node.data?.parent) {
                            return consts.groupMaxHeight;
                        } else {
                            return consts.nodeMaxHeight;
                        }
                    },
                    orthogonal: false
                },
                rotating: false,
            }),
        )
        .use(new Snapline({
            enabled: true,
        }))
        .use(new Keyboard({
            enabled: true,
            global: true,
        }))
        .use(new Clipboard({
            enabled: true,
        }))
        .use(new History({
            enabled: true,
        }))
        .use(
            new Scroller({
                enabled: true,
                pageVisible: true,
                pageBreak: true,
                pannable: true,
                // modifiers: "ctrl",
                // autoResize: true
            }),
        ).use(
            new Selection({
                enabled: true, //开启点击
                multiple: false,
                modifiers: ['ctrl', 'meta'],
                rubberband: false,
                strict: true,
                movable: true,
                showNodeSelectionBox: false,
                showEdgeSelectionBox: false
            }),
        )
    // #endregion


    // #region 快捷键与事件
    graph.bindKey(['meta+c', 'ctrl+c'], () => {
        const cells = graph.getSelectedCells()
        if (cells.length) {
            graph.copy(cells)
        }
        return false
    })
    graph.bindKey(['meta+x', 'ctrl+x'], () => {
        const cells = graph.getSelectedCells()
        if (cells.length) {
            graph.cut(cells)
        }
        return false
    })
    graph.bindKey(['meta+v', 'ctrl+v'], () => {
        if (!graph.isClipboardEmpty()) {
            const cells = graph.paste({ offset: 32 })
            graph.cleanSelection()
            graph.select(cells)
        }
        return false
    })

    // undo redo
    graph.bindKey(['meta+z', 'ctrl+z'], () => {
        if (graph.canUndo()) {
            graph.undo()
        }
        return false
    })
    graph.bindKey(['meta+shift+z', 'meta+y', 'ctrl+shift+z', 'ctrl+y'], () => {
        if (graph.canRedo()) {
            graph.redo()
        }
        return false
    })

    // select all
    graph.bindKey(['meta+a', 'ctrl+a'], () => {
        const nodes = graph.getNodes()
        if (nodes) {
            graph.select(nodes)
        }
    })

    // delete
    graph.bindKey('ctrl+delete', () => {
        const cells = graph.getSelectedCells()
        if (cells.length) {
            graph.removeCells(cells)
        }
    })

    // zoom
    graph.bindKey(['ctrl+1', 'meta+1', 'ctrl++'], () => {
        const zoom = graph.zoom()
        if (zoom < 1.5) {
            graph.zoom(0.1)
        }
    })
    graph.bindKey(['ctrl+2', 'meta+2', 'ctrl--'], () => {
        const zoom = graph.zoom()
        if (zoom > 0.5) {
            graph.zoom(-0.1)
        }
    })

    //格式化居中
    graph.bindKey('ctrl+shift+f', () => {
        graph.centerContent();
        graph.zoomTo(1)
        return false
    })

    return graph;
}


const connectingConfig = (dagGraph: DagGraph): Partial<Options.Connecting> => {
    return {
        anchor: "center",
        connector: 'algo-connector',
        connectionPoint: "anchor",
        snap: true,
        allowLoop: false,
        allowBlank: false,
        allowMulti: false,
        highlight: true,
        allowNode: false,
        allowEdge: false,
        allowPort: true,
        edgeAnchor: 'ratio',
        router: 'normal',
        createEdge() {
            let edge = cloneDeep(dagEdgeEntity);
            edge.id = dagGraph.newId() + "";
            return new Shape.Edge(edge);
        },
        validateMagnet({ magnet }) {
            //允许从out端口拖拽出连线
            return magnet.getAttribute("port-group") == "out";
        },
        validateConnection({
            sourceView,
            targetView,
            sourceMagnet,
            targetMagnet,
            targetPort,
            sourcePort,
            targetCell,
            sourceCell,
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            edge,
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            edgeView
        }: Options.ValidateConnectionArgs) {
            // 不能自身链接
            if (sourceView === targetView) {
                return false;
            }
            // 不能已输入桩为起点
            if (!sourceMagnet || sourceMagnet.getAttribute("port-group") !== "out") {
                return false;
            }
            // 只能连接到输入链接桩
            if (!targetMagnet || targetMagnet.getAttribute("port-group") !== "in") {
                return false;
            }

            let args;
            if (targetCell?.isNode()) {
                let ports = targetCell?.getPorts() || [];
                for (let i in ports) {
                    let port = ports[i];
                    if (port.id === targetPort) {
                        args = port.args;
                        break;
                    }
                }
            }
            if (args) {
                let maxIn = args.maxIn as number || 0;
                if (maxIn > 0) {
                    let edges = dagGraph.graph?.getIncomingEdges(targetCell || '');
                    let conned = 0;
                    if (edges) {
                        for (let j in edges) {
                            let e = edges[j] as any;
                            if (e.target.port === targetPort) {
                                conned++;
                            }
                        }
                    }
                    //连接数超
                    if (conned >= maxIn) {
                        return false;
                    }
                }
            }
            let portArgs: any;
            if (sourceCell?.isNode()) {
                let ports = sourceCell?.getPorts();
                console.log(ports);
                for (let i in ports) {
                    let port = ports[i];
                    if (port.id === sourcePort) {
                        portArgs = port.args;
                        break;
                    }
                }
            }

            if (portArgs) {
                let maxIns: number = portArgs.maxIn || 0;
                if (maxIns > 0) {
                    let edges = dagGraph.graph?.getOutgoingEdges(sourceCell || '');
                    let conneds = 0;
                    if (edges) {
                        for (let j in edges) {
                            let es = edges[j] as any;
                            if (es.source.port === sourcePort) {
                                conneds++;
                            }
                        }
                    }
                    let num = maxIns + 1
                    //连接数超
                    if (conneds >= num) {
                        console.log("超出连接数：" + maxIns)
                        return false;
                    }
                }
            }
            return true;
        }
    };
}
export default DagGraph;