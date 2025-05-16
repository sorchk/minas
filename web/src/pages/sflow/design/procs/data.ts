export const data_components = [{
    id: "BuildMap",
    name: "构造Map",
    icon: 'fa fa-paragraph',
    description: "",
    group: "data",
    fields: [{
        prop: 'buildmap',
        label: 'Map',
        type: 'map',
        value: [],
        placeholder: "例如:test=123或${param.id}"
    }]
}, {
    id: "DataAppend",
    name: "追加数据",
    icon: 'fa fa-paragraph',
    description: "",
    group: "data",
    fields: [{
        prop: "mdata",
        label: '主数据',
        type: 'input',
        placeholder: '输入任务节点编号,例如task_1',
        help: '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据',
        value: ''
    }, {
        prop: "sdata",
        label: '从数据',
        type: 'input',
        placeholder: '输入任务节点编号,例如task_1',
        help: '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据',
        value: ''
    }, {
        prop: "joins",
        label: '连接关系',
        labelMap: '主字段=从字段',
        width: '250=250',
        type: 'map',
        placeholder: '例如uid=例如userId',
        help: '',
        value: []
    }, {
        prop: "appendField",
        label: '追加字段',
        type: 'input',
        placeholder: '',
        help: '',
        value: ''
    }]
}, {
    id: "DataFieldAppend",
    name: "追加数据字段",
    icon: 'fa fa-plus',
    description: "",
    group: "data",
    fields: [{
        prop: "data",
        label: '数据',
        type: 'input',
        placeholder: '默认为上个任务节点的结果数据',
        help: '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据',
        value: ''
    }, {
        prop: 'field',
        label: '添加字段',
        labelMap: '字段=值',
        help: '值支持EL表达式,EL表达式可以从上面的参数中取数据',
        width: '120=70%;',
        type: 'map',
        value: [],
        placeholder: "例如:test=${param.test=='1'}"
    }]
}, {
    id: "DataFieldFilter",
    name: "过滤数据字段",
    icon: 'fa fa-filter',
    description: "",
    group: "data",
    fields: [{
        prop: "data",
        label: '数据',
        type: 'input',
        placeholder: '默认为上个任务节点的结果数据',
        value: '',
        help: '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop: 'isInclude',
        label: '筛选方式',
        type: 'switch',
        title: '包含|排除',
        value: 'true'
    }, {
        prop: 'field',
        label: '筛选字段',
        type: 'inputs',
        value: [],
        placeholder: '输入要包含或排除的字段名,例如:remark'
    }]
},{
    id: "DataFieldMapping",
    name: "数据字段映射",
    icon: 'fa fa-paragraph',
    description: "",
    group: "data",
    fields: [ {
        prop : "data",
        label : '数据',
        type : 'input',
        placeholder : '默认为上个任务节点的结果数据',
        value : '',
        help:'输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "fieldMappings",
        label : '字段映射',
        labelMap : '字段=重命名',
        width:'250=250',
        type : 'map',
        value: [],
        placeholder : '例如:name=label'
    } ]
},{
    id: "DataJoin",
    name: "数据链接",
    icon: 'fa fa-gg',
    description: "",
    group: "data",
    fields: [{
        prop : "mdata",
        label : '主数据',
        type : 'input',
        placeholder : '输入任务节点编号,例如task_1',
        value : '',
        help : '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "sdata",
        label : '从数据',
        type : 'input',
        placeholder : '输入任务节点编号,例如task_1',
        value : '',
        help : '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "joins",
        label : '连接关系',
        labelMap : '主字段=从字段',
        width : '250=250',
        type : 'map',
        placeholder : '例如uid=例如userId',
        value: [],
        help : ''
    } ]
},{
    id: "DataMerge",
    name: "数据合并",
    icon: 'fa fa-paragraph',
    description: "",
    group: "data",
    fields: [{
        prop : "resultType",
        label : '结果类型',
        type : 'radio',
        data : [ {
            value : 'list',
            label : "集合"
        }, {
            value : 'map',
            label : "单记录"
        } ],
        value : 'list'
    }, {
        prop : "data",
        label : '数据',
        type : 'inputs',
        placeholder : '默认为上个任务节点的结果数据',
        help : '输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据',
        value : []
    } ]
},{
    id: "DataList2Map",
    name: "List转Map",
    icon: 'fa fa-ellipsis-h',
    description: "",
    group: "data",
    fields: [{
        prop : "data",
        label : '数据',
        type : 'input',
        placeholder : '默认为上个任务节点的结果数据,数据必须为List格式',
        value : '',
        help:'输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "idField",
        label : '编号字段名',
        type : 'input',
        value : '',
        placeholder : '默认为id,此字段的值将作为Map的key,不允许空'
    }, {
        prop : 'isPk',
        label : '主键(唯一键)',
        type : 'switch',
        title : '是|否',
        value : 'false'
    }]
},{
    id: "DataList2Tree",
    name: "List转Tree",
    icon: 'fa fa-sitemap',
    description: "",
    group: "data",
    fields: [{
        id : "data",
        label : '数据',
        type : 'input',
        value : '',
        placeholder : '默认为上个任务节点的结果数据',
        help:'输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "rootId",
        label : '根节点编号',
        type : 'input',
        value : '',
        placeholder : '默认空,自动判断'
    }, {
        prop : "idField",
        label : '编号字段名',
        type : 'input',
        value : '',
        placeholder : '默认为id'
    }, {
        prop : "parentIdField",
        label : '父级编号字段名',
        type : 'input',
        value : '',
        placeholder : '默认为parentId'
    }, {
        prop : "chiledrenField",
        label : '子节点集合名',
        type : 'input',
        value : '',
        placeholder : '默认为chiledren'
    }]
},{
    id: "DataRow2Col",
    name: "行转列",
    icon: 'fa fa-level-down',
    description: "",
    group: "data",
    fields: [{
        prop : "data",
        label : '数据',
        type : 'input',
        value : '',
        placeholder : '默认为上个任务节点的结果数据,数据必须为List格式',
        help:'输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    }, {
        prop : "groupFields",
        label : '分组字段',
        type : 'inputs',
        value : '',
        placeholder : '允许空，如果此字段名为空，则数据将被扁平化为一条记录',
        help : '分组分类的字段名，多个使用连接符连接'
    } , {
        prop : "flatFields",
        label : '扁平化字段',
        type : 'inputs',
        value : '',
        placeholder : '允许空，如果此字段名为空，则数据将被格式化为KV形式的数据记录',
        help : '此字段值将作为扁平后的数据的字段名，多个使用连接符连接'
    } , {
        prop : "valueField",
        label : '值字段',
        type : 'input',
        value : '',
        placeholder : '不能为空,值字段名',
        help : '此字段名的值当做扁平化后的数据值'
    } , {
        prop : "connector",
        label : '连接符号',
        type : 'input',
        placeholder : '默认为_',
        value:'_',
        help : '多字段时使用连接符号连接字段名或数据'
    }  ]
},{
    id: "DataTypeConvert",
    name: "数据类型转换",
    icon: 'fa fa-exchange',
    description: "",
    group: "data",
    fields: [{
        prop : "data",
        label : '数据',
        type : 'input',
        placeholder : '默认为上个任务节点的结果数据',
        value : '',
        help:'输入任务节点编号,例如task_1,获取此任务节点数据,默认为上个任务节点的数据'
    },  {
        prop : 'field',
        label : '类型及格式',
        labelMap : '字段=类型及格式',
        width:'200=280',
        type : 'map',
        value : [],
        placeholder : "例如:test=date-yyyy-MM-dd HH"
    }]
}
]
