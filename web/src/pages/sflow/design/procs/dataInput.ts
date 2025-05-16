export const input_components = [{
    id: "CsvLoad",
    name: "CSV文件",
    description: "",
    group: "input",
    icon: 'fa fa-file',
    fields: [{
        prop: "csvpath",
        label: 'CSV文件路径',
        type: 'input',
        value: ''
    }, {
        prop: "encoding",
        label: '文件编码',
        type: 'radio',
        data: [{
            value: '',
            label: "自动"
        }, {
            value: 'UTF-8',
            label: "UTF-8"
        }, {
            value: 'GBK',
            label: "GBK"
        }, {
            value: 'GB2312',
            label: "GB2312"
        }, {
            value: 'ISO-8859-1',
            label: "ISO-8859-1"
        }],
        value: ''
    }, {
        prop: "title",
        label: '标题',
        type: 'input',
        placeholder: '可空,默认为csv第一行',
        value: ''
    }, {
        prop: "titleLineNum",
        label: '标题行号',
        type: 'input',
        value: "",
        placeholder: '可空,默认为1'
    }, {
        prop: "dataLineNum",
        label: '数据行号',
        type: 'input',
        value: "",
        placeholder: '可空,默认为标题行号加1'
    }, {
        prop: "rowSeparator",
        label: '行分隔符',
        type: 'input',
        value: "",
        placeholder: '可空,默认\n'
    }, {
        prop: "colSeparator",
        label: '列隔符',
        type: 'input',
        value: "",
        placeholder: '可空,默认,'
    }]
}, {
    id: "SqlDataQuery",
    name: "Sql查询",
    icon: 'fa fa-search',
    description: "",
    group: "input",
    fields: [{
        prop: "ds",
        label: '数据源',
        type: 'select',
        url: '/rs/meta/ds/select?type=sql',
        data: [{
            value: 'test',
            label: "test"
        }],
        labelName: "name",
        valueName: "id",
        multiple: false,
        value: '',
        placeholder: '请选择'
    }, {
        prop: "resultType",
        label: '结果类型',
        type: 'radio',
        data: [{
            value: 'list',
            label: "集合"
        }, {
            value: 'map',
            label: "单记录"
        }, {
            value: 'int',
            label: "整数"
        }, {
            value: 'string',
            label: "字符串"
        }, {
            value: 'string[]',
            label: "字符串数组"
        }],
        value: 'list'
    }, {
        prop: "sql",
        label: '查询sql语句',
        type: 'editor',
        lang: 'sql',
        value: ''
    }, {
        prop: "params",
        label: '参数',
        type: 'inputs',
        value: []
    }, {
        prop: "sort",
        label: '排序字段',
        type: 'input',
        value: ''
    }, {
        prop: 'sortType',
        label: '排序方式',
        type: 'switch',
        title: '升序|降序',
        value: 'true'
    }, {
        prop: 'pageType',
        label: '分页',
        type: 'radio',
        data: [{
            value: '',
            label: "不分页"
        }, {
            value: 'sqlserver',
            label: "sqlserver"
        }, {
            value: 'mysql',
            label: "mysql"
        }],
        value: ''
    }, {
        prop: 'pageSize',
        label: '每页显示几条',
        type: 'input',
        value: ''
    }, {
        prop: 'pageNo',
        label: '页码',
        type: 'input',
        value: ''
    }, {
        prop: 'countSql',
        label: '分页计数Sql',
        type: 'editor',
        lang: 'sql',
        value: ''
    }]
}, {
    id: "TrinoSqlQuery",
    name: "TrinoSql查询",
    icon: 'fa fa-quora',
    description: "",
    group: "input",
    fields: [{
        prop: "resultType",
        label: '结果类型',
        type: 'radio',
        data: [{
            value: 'list',
            label: "集合"
        }, {
            value: 'map',
            label: "单记录"
        }, {
            value: 'int',
            label: "整数"
        }, {
            value: 'string',
            label: "字符串"
        }, {
            value: 'array',
            label: "单记录数组"
        }],
        value: 'list'
    }, {
        prop: "sql",
        label: '查询sql语句',
        type: 'editor',
        height: 400,
        lang: 'sql',
        value: ''
    }, {
        prop: "params",
        label: '参数',
        type: 'inputs',
        value: []
    }, {
        prop: 'offset',
        label: '开始记录',
        type: 'input',
        value: '0',
        help: '从第几条记录开始返回数据，默认0从头。也就是分页的(页码-1)*分页大小'
    }, {
        prop: 'limit',
        label: '查询记录数',
        type: 'input',
        value: '0',
        help: '查询最大返回多少条数据，默认0返回所有。也就是分页大写'
    }]
}, {
    id: "SqlQuery",
    name: "SQL查询(测试)",
    icon: 'fa fa-quora',
    description: "",
    group: "input",
    fields: [{
        prop: "ds",
        label: '数据源',
        type: 'select',
        url: '/rs/meta/ds/select?type=sql',
        data: [{
            value: 'test',
            label: "test"
        }],
        labelName: "name",
        valueName: "id",
        multiple: false,
        value: '',
        placeholder: '请选择'
    }, {
        prop: "resultType",
        label: '结果类型',
        type: 'radio',
        data: [{
            value: 'list',
            label: "集合"
        }, {
            value: 'map',
            label: "单记录"
        }, {
            value: 'int',
            label: "整数"
        }, {
            value: 'string',
            label: "字符串"
        }, {
            value: 'string[]',
            label: "字符串数组"
        }],
        value: 'list'
    }, {
        prop: "sql",
        label: '查询sql语句',
        type: 'editor',
        lang: 'sql',
        value: ''
    }, {
        prop: "params",
        label: '参数',
        type: 'inputs',
        value: []
    }]
}, {
    id: "TableDataQuery",
    name: "表查询(测试)",
    icon: 'fa fa-table',
    description: "",
    group: "input",
    fields: [{
        prop: "ds",
        label: '数据源',
        type: 'select',
        url: '/rs/meta/ds/select?type=sql',
        data: [{
            value: 'test',
            label: "test"
        }],
        labelName: "name",
        valueName: "id",
        multiple: false,
        value: '',
        placeholder: '请选择'
    }, {
        prop: "table",
        label: '表名',
        type: 'input',
        value: ''
    }, {
        prop: "filters",
        label: '过滤条件',
        type: "map",
        value: [],
        placeholder: "例如:123或${param.id}"
    }, {
        prop: "sorts",
        label: '排序字段',
        type: 'inputs',
        value: [],
    }, {
        prop: 'size',
        label: '每页显示几条',
        type: 'number',
        value: 0,
        help: '0为不分页'
    }]
}, {
    id: "MongoDataQuery",
    name: "Mongodb查询",
    icon: 'fa fa-envira',
    description: "",
    group: "input",
    fields: [{
        prop: "ds",
        label: '数据源',
        type: 'select',
        url: '/rs/meta/ds/select?subtype=mongodb',
        labelName: "name",
        valueName: "id",
        multiple: false,
        value: '',
        placeholder: '请选择'
    }, {
        prop: "collection",
        label: '表（集合）',
        value: ''
    }, {
        prop: "resultType",
        label: '结果类型',
        type: 'radio',
        data: [{
            value: 'list',
            label: "集合"
        }, {
            value: 'map',
            label: "单记录"
        }, {
            value: 'int',
            label: "整数"
        }],
        value: 'list'
    }, {
        prop: "filterJson",
        label: '查询条件模板',
        type: 'editor',
        lang: 'json',
        value: ''
    }, {
        prop: "sort",
        label: '排序字段',
        value: ''
    }, {
        prop: 'sortType',
        label: '排序方式',
        type: 'switch',
        active: "升序",
        inactive: "降序",
        value: true
    }, {
        prop: 'pageType',
        label: '分页',
        type: 'switch',
        active: "分页", inactive: "不分页",
        value: false
    }, {
        prop: 'pageSize',
        label: '每页显示几条',
        type: 'number',
        value: ''
    }, {
        prop: 'pageNo',
        label: '页码',
        type: 'number',
        value: ''
    }]
}, {
    id: "MongoCommand",
    name: "MongoCommand",
    icon: 'fa fa-leaf',
    description: "",
    group: "input",
    fields: [{
        prop: "ds",
        label: '数据源',
        type: 'select',
        url: '/rs/meta/ds/select?subtype=mongodb',
        labelName: "name",
        valueName: "id",
        multiple: false,
        value: '',
        placeholder: '请选择'
    }, {
        prop: "cmd",
        label: '执行命令',
        type: 'editor',
        lang: 'json',
        height: '500',
        value: ''
    }]
}, {
    id: "HttpRequest",
    name: "http接口",
    icon: "fa fa-internet-explorer",
    description: "",
    group: "input",
    fields: [
        { prop: "url", label: "请求地址", value: "", placeholder: "http://test.com/getInfo" },
        {
            label: '请求方法',
            prop: 'method',
            type: 'radio',
            value: "GET",
            button: true,
            data: [{
                value: 'GET',
                label: "GET"
            }, {
                value: 'POST',
                label: "POST"
            }, {
                value: 'PUT',
                label: "PUT"
            }, {
                value: 'DELETE',
                label: "DELETE"
            }]
        },
        {
            label: '数据类型',
            prop: 'contentType',
            type: 'radio',
            value: "none",
            button: true,
            data: [{
                value: 'none',
                label: "无"
            }, {
                value: 'application/json',
                label: "json"
            }, {
                value: 'application/x-www-form-urlencoded',
                label: "x-www-form-urlencoded"
            }, {
                value: 'multipart/form-data',
                label: "form-data"
            }]
        },
        {
            label: '编码',
            prop: 'encoding',
            type: 'radio',
            value: "UTF-8",
            span: 24,
            button: true,
            data: [{
                value: 'UTF-8',
                label: "UTF-8"
            }, {
                value: 'GBK',
                label: "GBK"
            }, {
                value: 'ISO-8859-1',
                label: "ISO-8859-1"
            }]
        },
        {
            prop: "heardes",
            label: "请求头",
            type: "map",
            value: [],
            keyLabel: "请求头",
            valueLabel: "值",
            placeholder: "例如:123或${param.id}"
        },
        {
            prop: "params",
            label: "请求参数",
            type: "map", value: [],
            keyLabel: "参数名",
            valueLabel: "参数值",
            placeholder: "例如:123或${param.id}"
        }, {
            prop: "body",
            label: '请求消息体',
            type: 'editor',
            lang: 'json',
            vif: function (form) {
                return form.contentType == 'application/json';
            },
            value: '',
            help: '需要设置数据类型为json'
        },
        { prop: "timeout", label: "超时", value: "", placeholder: "设置请求超时，默认10000毫秒", help: '单位毫秒' },
        { prop: "supportCookie", label: "支持会话", type: "switch", active: "是", inactive: "否", value: false }
    ]
}
];
