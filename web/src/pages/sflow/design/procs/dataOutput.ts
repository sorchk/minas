
export const output_components = [
    {
        id: "SqlDataPersistence",
        name: "Sql数据库存储",
        icon: 'fa fa-database',
        description: "",
        group: "output",
        fields: [{
            prop: "ds",
            label: '数据源',
            type: 'select',
            url: '/rs/meta/ds/select?type=sql',
            root: 'data',
            labelName: "name",
            valueName: "id",
            value: '',
            placeholder: '请选择'
        }, {
            prop: "sql",
            label: 'sql语句',
            type: 'editor',
            value: '',
            lang: 'sql'
        }, {
            prop: "data",
            label: '参数数据',
            type: 'input',
            value: '',
            help: '可空,例如:task_1,为取任务1的结果数据,如果结果为List则循环执行sql',
            placeholder: '可空,例如:task_1'
        }, {
            prop: "params",
            label: '参数',
            value: [],
            type: 'inputs'
        }]
    },
    {
        id: "MongoDataPersistence",
        name: "Mongodb存储",
        icon: 'fa fa-envira',
        description: "",
        group: "output",
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
            prop: "op",
            label: '操作类型',
            type: 'radio',
            data: [{
                value: 'save',
                label: "保存"
            }, {
                value: 'insert',
                label: "添加"
            }, {
                value: 'update',
                label: "更新"
            }, {
                value: 'updateMany',
                label: "更新（多记录）"
            }, {
                value: 'delete',
                label: "删除"
            }],
            value: ''
        }, {
            prop: "data",
            label: '参数数据',
            help: '可空,例如:task_1,为取任务1的结果数据,如果结果为List则循环执行sql',
            placeholder: '可空,例如:task_1',
            value: ''
        }, {
            prop: "filterJson",
            label: '条件模板',
            type: 'editor',
            lang: 'json',
            value: ''
        }, {
            prop: "dataJson",
            label: '数据模板',
            type: 'editor',
            lang: 'json',
            value: ''
        }, {
            prop: 'ignoreKeyDuplicate',
            label: '忽略主键重复错误',
            type: 'switch',
            active: "是", inactive: "否",
            value: true
        }, {
            prop: 'overall',
            label: '保存时删除不存在的字段',
            type: 'switch',
            active: "是", inactive: "否",
            value: false
        }]
    }, {
        id: "TableDataSave",
        name: "表增删改(测试)",
        icon: 'fa fa-table',
        description: "",
        group: "output",
        fields: [
            {
                label: '操作类型',
                prop: 'op',
                type: 'radio',
                value: 'save',
                button: true,
                data: [{
                    value: 'insert',
                    label: "添加"
                }, {
                    value: 'update',
                    label: "修改"
                }, {
                    value: 'save',
                    label: "添加或修改"
                }, {
                    value: 'delete',
                    label: "删除"
                }]
            }, {
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
                prop: "params",
                label: '数据',
                type: "map",
                value: [],
                placeholder: "例如:123或${param.id}"
            }, {
                prop: "pks",
                label: '主键',
                type: 'inputs',
                value: [],
                vif: function (form: any) {
                    return form.op != 'insert';
                }
            }]
    },
];
