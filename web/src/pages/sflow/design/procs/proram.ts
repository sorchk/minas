
const apiURL = window.location.href.split("/");
const apiId = apiURL[apiURL.length - 1];
export const program_components = [
    {
        id: "Depute",
        name: "委托任务",
        icon: 'fa fa-handshake-o',
        description: "",
        group: "base",
        fields: [
            {
                prop: "id",
                label: '委托任务',
                type: 'selectTree',
                url: '/rs/srv/api/data_processor/select',
                props: {
                    label: 'name',
                    children: 'children',
                    isLeaf: function (data) {
                        return data.type == 'file';
                    },
                    disabled: function (data) {
                        return data.id == apiId || (data.type == 'dir' && (!data.children || data.children.length < 1));
                    }
                },
                dataRoot: '',
                labelName: "name",
                valueName: "id",
                multiple: false,
                value: ''
            }, {
                prop: "params",
                label: '参数',
                labelMap: '参数名=参数值',
                width: '200=280',
                type: 'map',
                value: [],
                placeholder: "例如:test=1",
                help: '任务内取参数值为${params.参数键}'
            }, {
                prop: "ignoreSubErr",
                label: '忽略任务错误',
                type: 'switch',
                title: '是|否',
                help: '如果是则返回委托任务错误信息，不抛出异常，程序将继续运行',
                value: false
            }, {
                prop: "isReturn",
                label: '是否返回数据',
                type: 'switch',
                title: '是|否',
                help: '在日志级别为数据（DATA）以上级别时大量返回数据将造成日志存储错误，也会造成内存占用过高',
                value: false
            }
        ]
    },
    {
        id: "ForEachLoop",
        name: "遍历集合",
        icon: 'fa fa-th',
        description: "",
        group: "base",
        type: 'group',
        fields: [
            {
                prop: 'data',
                label: '循环数据',
                type: 'input',
                placeholder: '循环数据需要为List，如果为Map则仅执行一次',
                help: '取List长度为${forCount},循环内取当前数据为${forData},取当前索引为${forIndex},是否是第一个记录${forIsFirst},是否是最后一个记录${forIsLast}',
                value: ''
            }, {
                prop: "forBreakEl",
                label: '跳出(break)',
                type: 'input',
                value: '',
                placeholder: '填写跳出循环的el表达式',
                help: ''
            }, {
                prop: "forContinueEl",
                label: '跳过(continue)',
                type: 'input',
                value: '',
                placeholder: '填写跳过当前索引的el表达式',
                help: ''
            }, {
                prop: "params",
                label: '参数',
                labelMap: '参数名=参数值',
                width: '200=280',
                type: 'map',
                value: [],
                placeholder: "例如:test=1",
                help: '循环内取参数值为${参数键}'
            }, {
                prop: "ignoreSubErr",
                label: '忽略循环内错误',
                type: 'switch',
                title: '是|否',
                help: '如果是则返回循环内误信息，否则不抛出异常，程序将继续运行',
                value: false
            }, {
                prop: "isReturn",
                label: '是否返回数据',
                type: 'switch',
                title: '是|否',
                help: '在日志级别为数据（DATA）以上级别时大量返回数据将造成日志存储错误，也会造成内存占用过高',
                value: false
            }
        ]
    },
    {
        id: "ForLoop",
        name: "循环",
        icon: 'fa fa-retweet',
        description: "",
        group: "base",
        type: 'group',
        fields: [
            {
                prop: "forBegin",
                label: '开始',
                type: 'input',
                placeholder: '循环开始索引，默认值为0',
                value: '0',
                help: '${forBegin}'
            }, {
                prop: "forStep",
                label: '步长',
                type: 'input',
                value: '1',
                placeholder: '循环步长，默认值为1',
                help: '${forStep}'
            }, {
                prop: "forEnd",
                label: '结束(不含)',
                type: 'input',
                value: '',
                placeholder: '结束的索引号，必须填写，否则循环不执行',
                help: '${forEnd}'
            }, {
                prop: "forBreakEl",
                label: '跳出(break)',
                type: 'input',
                value: '',
                placeholder: '填写跳出循环的el表达式',
                help: ''
            }, {
                prop: "forContinueEl",
                label: '跳过(continue)',
                type: 'input',
                value: '',
                placeholder: '填写跳过当前索引的el表达式',
                help: ''
            }, {
                prop: "params",
                label: '参数',
                labelMap: '参数名=参数值',
                width: '200=280',
                type: 'map',
                value: [],
                placeholder: "例如:test=1",
                help: '循环内取参数值为${参数键}'
            }, {
                prop: "ignoreSubErr",
                label: '忽略循环内错误',
                type: 'switch',
                title: '是|否',
                help: '如果是则返回循环内误信息，否则不抛出异常，程序将继续运行',
                value: false
            }, {
                prop: "isReturn",
                label: '是否返回数据',
                type: 'switch',
                title: '是|否',
                help: '在日志级别为数据（DATA）以上级别时大量返回数据将造成日志存储错误，也会造成内存占用过高',
                value: false
            }
        ]
    }, {
        id: "DataState",
        name: "状态机",
        icon: "fa fa-archive",
        description: "",
        group: "base",
        fields: [
            {
                label: '操作类型',
                prop: 'op',
                type: 'radio',
                value: 0,
                button: true,
                data: [{
                    value: 0,
                    label: "读状态"
                }, {
                    value: 1,
                    label: "写状态"
                }, {
                    value: 2,
                    label: "注册状态机"
                }]
            },
            {
                prop: "key",
                label: "标识",
                placeholder: "请输入状态的唯一标识",
                value: '',
                help: '状态的唯一标识'
            },
            {
                prop: "value",
                label: "状态值",
                placeholder: "",
                value: '',
                vif: function (form) {
                    return form.op == 1;
                }
            },
            {
                prop: "type",
                label: "状态值类型",
                type: 'radio',
                value: 0,
                button: true,
                vif: function (form) {
                    return form.op == 2;
                },
                data: [{
                    value: 0,
                    label: "字符串"
                }, {
                    value: 1,
                    label: "数字"
                }, {
                    value: 2,
                    label: "日期"
                }]
            },
            {
                prop: "format",
                label: "格式",
                placeholder: "请设置日期时间数据格式，默认：YYYY-MM-dd HH:mm:ss.SSS",
                value: '',
                vif: function (form) {
                    return form.op == 2 && form.type == 2;
                },
                help: '例如：YYYY-MM-dd HH:mm:ss.SSS 或 YYYY年MM月dd日 HH时mm分ss秒'
            },
            {
                prop: "format",
                label: "格式",
                placeholder: "请设置字符串数据格式，默认空",
                value: '',
                vif: function (form) {
                    return form.op == 2 && form.type == 2;
                },
                help: '例如：%08d 长度不够8位的自动前补0'
            },
            {
                prop: "value",
                label: "初始状态",
                placeholder: "",
                value: '',
                vif: function (form) {
                    return form.op == 2;
                },
            }
        ]
    }, {
        id: "CheckPoint",
        name: "检查点",
        icon: "fa fa-map-pin",
        description: "",
        group: "base",
        fields: [
            {
                label: '操作类型',
                prop: 'op',
                type: 'radio',
                value: 0,
                button: true,
                data: [{
                    value: 0,
                    label: "读检查点"
                }, {
                    value: 1,
                    label: "写检查点"
                }, {
                    value: 2,
                    label: "注册检查点"
                }]
            },
            {
                prop: "checkPoint",
                label: "检查点",
                placeholder: "请输入检查点标识",
                value: '',
                help: '读取检查点时需要使用'
            },
            {
                prop: "value",
                label: "存档数据",
                type: "map",
                keyLabel: "数据键",
                valueLabel: "数据值",
                placeholder: "",
                value: [],
                vif: function (form) {
                    return form.op == 1;
                }
            },
            {
                prop: "value",
                label: "初始检查点数据",
                type: "map",
                keyLabel: "数据键",
                valueLabel: "数据值",
                placeholder: "",
                value: [],
                vif: function (form) {
                    return form.op == 2;
                },
            }
        ]
    }, {
        id: "DirtyData",
        name: "脏数据操作",
        icon: "fa fa-optin-monster",
        description: "",
        group: "base",
        fields: [
            {
                label: '操作类型',
                prop: 'op',
                type: 'radio',
                value: 1,
                button: true,
                data: [{
                    value: 0,
                    label: "读取列表"
                }, {
                    value: 1,
                    label: "写入"
                }, {
                    value: 2,
                    label: "读取一条"
                }, {
                    value: 3,
                    label: "修改状态"
                }]
            },
            {
                prop: "dataTable",
                label: "数据表",
                placeholder: "请输入数据表名",
                value: '',
                vif: function (form) {
                    return form.op == 1 || form.op == 0;
                },
                help: '一般直接使用数据库的表名，如果跨多个库的脚本可以使用数据库名.模式.表名，也可以不是表名自己定义唯一标识也可以'
            },
            {
                prop: "id",
                label: "脏数据ID",
                placeholder: "请输入脏数据的ID，一般使用变量",
                value: '',
                vif: function (form) {
                    return form.op == 2 || form.op == 3;
                },
                help: ''
            },
            {
                prop: "value",
                label: "脏数据",
                placeholder: "输入脏数据的el表达式",
                value: '',
                vif: function (form) {
                    return form.op == 1;
                },
                help: '输入脏数据EL表达式，也可以是{key:value}格式的jsonMap字符串数据'
            },
            {
                prop: "status",
                label: "状态",
                type: "number",
                value: 0,
                vif: function (form) {
                    return form.op == 3;
                },
            }
        ]
    }, {
        id: "Counter",
        name: "计数器",
        icon: "fa fa-info",
        description: "",
        group: "base",
        fields: [
            {
                label: '操作类型',
                prop: 'op',
                type: 'radio',
                value: 0,
                button: true,
                data: [{
                    value: 0,
                    label: "读取"
                }, {
                    value: 1,
                    label: "更新"
                }, {
                    value: 2,
                    label: "初始化"
                }, {
                    value: 3,
                    label: "读取多个"
                }]
            },
            {
                prop: "key",
                label: "计数器标识",
                placeholder: "请输入计数器标识",
                value: '',
                help: '读取计数器值时需要使用',
                vif: function (form) {
                    return form.op != 3;
                }
            },
            {
                prop: "keys",
                label: "计数器标识",
                type: "inputs",
                placeholder: "请输入计数器标识",
                value: [],
                help: '返回数据为map,以计数据标识为key',
                vif: function (form) {
                    return form.op == 3;
                }
            },
            {
                prop: "logInfo",
                label: "日志信息",
                placeholder: "允许空，空则不输出日志",
                value: '',
                help: '${countData}为当前计数器值，日志级别为INFO'
            },
            {
                prop: "step",
                label: "步长",
                type: "number",
                value: 1,
                min: 1,
                vif: function (form) {
                    return form.op == 1;
                }
            }
        ]
    }, {
        id: "Log",
        name: "日志",
        icon: "fa fa-commenting",
        description: "",
        group: "base",
        fields: [
            {
                label: '日志级别',
                prop: 'logLevel',
                type: 'radio',
                value: "info",
                button: true,
                data: [{
                    value: 'debug',
                    label: "调试"
                }, {
                    value: 'info',
                    label: "信息"
                }, {
                    value: 'warn',
                    label: "警告"
                }, {
                    value: 'error',
                    label: "错误"
                }, {
                    value: 'none',
                    label: "无"
                }],
                help: '日志级别为无时不输出日志'
            },
            {
                prop: "logInfo",
                label: "日志信息",
                type: "editor",
                placeholder: "允许空，空则不输出日志，支持EL表达式",
                value: ''
            },
        ]
    }
]
