export const script_components = [
    {
        id: "EL",
        name: "EL表达式",
        description: "",
        group: "script",
        icon: 'fa fa-etsy',
        fields: [
            {
                label: '表达式语言',
                prop: 'elType',
                type: 'radio',
                value: "juel",
                button: true,
                data: [{
                    label: 'MVEL',
                    value: 'mvel'
                }, {
                    label: 'MVEL模板',
                    value: 'mvel_template'
                }, {
                    label: 'JUEL',
                    value: 'juel'
                }]
            },
            {
                prop: "elVars",
                label: "定义变量",
                type: "map", value: [],
                placeholder: "例如:123或${param.id}"
            },
            { prop: "el", label: "EL表达式", value: "", placeholder: "请输入" }
        ]
    }, {
        id: "JavaScript",
        name: "JavaScript脚本",
        icon: 'fa fa-code',
        description: "",
        group: "script",
        fields: [
            {
                prop: "scriptVars",
                label: "定义变量",
                type: "map",
                value: [],
                placeholder: "例如:123或${param.id}"
            },
            { prop: "scriptText", label: "脚本", type: "editor", value: "", lang: 'javascript', placeholder: "请输入" },
            {
                prop: "compilable",
                label: "编译脚本",
                type: "switch",
                active: "是",
                inactive: "否",
                value: true,
                help: '编译脚本将提供运行速度，需要脚本语言本身支持，否则不编译。'
            }
        ]

    },
    {
        id: "RegexData",
        name: "正则提取数据",
        icon: 'fa fa-asterisk',
        description: "",
        group: "script",
        fields: [
            { prop: "data", label: "被查找数据", value: "", placeholder: "请输入被查找数据，支持EL表达式，结果是文本值" },
            { prop: "regex", label: "正则表达式", value: "", placeholder: "请输入正则表达式" },
            {
                prop: "isCase", label: "区分大小写", type: "switch", active: "是", inactive: "否", value: true
            },
            { prop: "template", label: "数据模板", value: "", placeholder: "请输入数据模板，例如$1=$2。默认为$1提取第一个分组的数据" }
        ]
    }];
