export const system_components = [{
    id: "Shell",
    name: "系统命令行",
    description: "",
    group: "system",
    icon: 'fa fa-terminal',
    fields: [
        {
            prop: "shell",
            label: "程序及参数",
            type: "inputs",
            placeholder: "",
            value: [],
            help: '程序及参数'
        },
        {
            prop: "envp",
            label: "环境变量",
            type: "map",
            value: [],
            placeholder: "例如:123或${param.id}"
        },

        { prop: "dir", label: "运行目录", value: "", placeholder: "运行目录" },
        {
            prop: "waitResult",
            label: "等待返回结果",
            type: "switch",
            value: false,
            active: "是",
            inactive: "否",
            placeholder: "",
            help: '是否等待返回结果,如果否此节点数据永远是null。'
        }, {
            label: '结果编码格式',
            prop: 'encoding',
            type: 'radio',
            value: "UTF-8",
            button: true,
            data: [{
                label: 'UTF-8',
                value: 'UTF-8'
            }, {
                label: 'GBK',
                value: 'GBK'
            }],
            help: '默认UTF-8，windows系统GBK'
        },
        {
            prop: "timeout", label: "运行超时", type: "number",
            value: "5000",
            placeholder: "单位毫秒",
            controlsPosition: '',
            controls: true,
            step: 500,
            precision: 0,
            stepStrictly: true,
            help: '单位毫秒,默认一直等待到完成，或出现错误终止。如果输入时间将在超时后终止进程'
        }
    ]
}, {
    id: "NetTest",
    name: "网络测试",
    icon: 'fa fa-heartbeat',
    description: "",
    group: "system",
    fields: [
        {
            label: '测试类型',
            prop: 'type',
            type: 'radio',
            value: "http",
            button: true,
            data: [{
                value: 'http',
                label: "http"
            }, {
                value: 'ping',
                label: "ping"
            }, {
                value: 'telnet',
                label: "telnet"
            }]
        },
        {
            prop: "url",
            label: "测试地址",
            placeholder: "",
            value: '',
            help: 'http格式为：请求方法 协议://ip:端口/路径，ping格式为ip，telnet格式为ip:端口'
        },
        {
            prop: "timeout", label: "运行超时", type: "number",
            value: "1000",
            placeholder: "单位毫秒，默认1000毫秒",
            controlsPosition: '',
            controls: true,
            step: 500,
            precision: 0,
            stepStrictly: true
        },
        {
            prop: "retries", label: "重试次数", type: "number",
            value: "3",
            placeholder: "默认失败后重试3次",
            controlsPosition: '',
            controls: true,
            step: 1,
            precision: 0,
            stepStrictly: true
        },
        {
            prop: "thread", label: "线程数", type: "number",
            value: "3",
            placeholder: "默认3",
            controlsPosition: '',
            controls: true,
            step: 1,
            precision: 0,
            stepStrictly: true
        }
    ]
}, {
    id: "File",
    name: "文件操作",
    icon: "fa fa-archive",
    description: "",
    group: "system",
    fields: [
        {
            label: '操作类型',
            prop: 'optype',
            type: 'radio',
            value: "ls",
            button: true,
            data: [{
                value: 'ls',
                label: "列目录"
            }, {
                value: 'cp',
                label: "复制"
            }, {
                value: 'mv',
                label: "移动"
            }, {
                value: 'rn',
                label: "重命名"
            }, {
                value: 'rm',
                label: "删除"
            }, {
                value: 'wt',
                label: "写入文本"
            }, {
                value: 'rt',
                label: "读文本"
            }, {
                value: 'wo',
                label: "写入对象"
            }]
        },
        {
            prop: "src",
            label: "源路径",
            placeholder: "",
            value: '',
            vif: function (form) {
                return form.optype == 'cp' || form.optype == 'rn' || form.optype == 'mv' || form.optype == 'cp';
            },
            help: '复制、移动和重命名需要输入原路径'
        },
        {
            prop: "content",
            label: "内容",
            type: 'textarea',
            rows: 8,
            placeholder: "",
            value: '',
            vif: function (form) {
                return form.optype == 'wt';
            },
            help: '写入文件的内容'
        },
        {
            prop: "content",
            label: "内容",
            type: 'input',
            placeholder: "",
            value: '',
            vif: function (form) {
                return form.optype == 'wo';
            },
            help: '写入文件的对象表达式'
        },
        {
            prop: "target",
            label: "目标路径",
            placeholder: "",
            value: '',
            help: '操作目标文件或目录的路径'
        }
    ]
}, {
    id: "MinioFile",
    name: "Minio文件操作",
    icon: "fa fa-archive",
    description: "",
    group: "system",
    fields: [
        {
            prop: "endpoint",
            label: "地址",
            type: 'input',
            placeholder: "请输入服务器地址，例如：https://localhost:9000/",
            value: '',
            help: ''
        },
        {
            prop: "accessKey",
            label: "访问账号",
            type: 'input',
            placeholder: "请输入访问账号",
            value: '',
            help: ''
        },
        {
            prop: "secretKey",
            label: "密钥",
            type: 'password',
            placeholder: "请输入访问密钥",
            value: '',
            help: ''
        },
        {
            label: '操作类型',
            prop: 'optype',
            type: 'radio',
            value: "ls",
            button: true,
            data: [{
                value: 'ls',
                label: "列目录"
            }, {
                value: 'rm',
                label: "删除"
            }, {
                value: 'wt',
                label: "写入文本"
            }, {
                value: 'rt',
                label: "读文本"
            }, {
                value: 'wo',
                label: "写入对象"
            }, {
                value: 'ro',
                label: "读对象"
            }, {
                value: 'info',
                label: "读文件元数据"
            }, {
                value: 'exist',
                label: "文件是否存在"
            }]
        },
        {
            prop: "max",
            label: "文件数",
            type: 'number',
            placeholder: "",
            value: '100',
            vif: function (form) {
                return form.optype == 'ls';
            },
            help: '最大显示文件数'
        },
        {
            prop: "bufferSize",
            label: "缓冲区大小",
            type: 'number',
            placeholder: "",
            value: '4096',
            vif: function (form) {
                return form.optype == 'ro' || form.optype == 'wo' || form.optype == 'rt' || form.optype == 'wt';
            },
            help: '缓冲区大小，默认4k'
        },
        {
            prop: "content",
            label: "内容",
            type: 'textarea',
            rows: 8,
            placeholder: "",
            value: '',
            vif: function (form) {
                return form.optype == 'wt';
            },
            help: '写入文件的内容'
        },
        {
            prop: "content",
            label: "内容",
            type: 'input',
            placeholder: "",
            value: '',
            vif: function (form) {
                return form.optype == 'wo';
            },
            help: '写入文件的对象表达式'
        },
        {
            prop: "buckets",
            label: "文件桶",
            placeholder: "",
            value: '',
            help: ''
        },
        {
            prop: "target",
            label: "路径",
            placeholder: "",
            value: '',
            help: '操作文件的路径,不包括文件桶'
        },
        {
            prop: "contentType",
            label: "文件类型",
            placeholder: "允许空,系统将自动识别，文本、文件对象、文件路径的文件类型，字节、其它对象无法识别",
            value: '',
            vif: function (form) {
                return form.optype == 'wo' || form.optype == 'wt';
            },
        },
        {
            prop: "meta",
            label: "文件元数据",
            type: "map",
            value: [],
            keyLabel: "属性名",
            valueLabel: "属性值",
            placeholder: "例如:123或${param.id}",
            help: '文件元数据',
            vif: function (form) {
                return form.optype == 'wo' || form.optype == 'wt';
            },

        }
    ]
}];

