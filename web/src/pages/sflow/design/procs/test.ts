import testNode from "../shape/TestNode.vue";

export const test_components = [
    {
        id: "simple",
        name: "测试组件",
        proc: "simple",
        group: "base",
        inherit: 'vue-shape',
        component: {
            components: {
                testNode
            },
            template: `
              <test-node/>`
        },
        form: {"label": "测试组件"},
        fields: [
            {prop: "fullname", label: "姓名", value: "", placeholder: "请输入姓名"},
            {prop: "remark", label: "备注", value: "", type: "textarea", placeholder: "请输入备注"},
            {prop: "pwd", label: "密码", type: "password", value: "", required: true, placeholder: "请输入密码"},
            {
                prop: "age", label: "年龄", type: "number",
                value: "5",
                required: true,
                placeholder: "请选择年龄",
                controlsPosition: '',
                controls: true,
                "step": 1,
                "precision": 0,
                "stepStrictly": true
            }, {
                label: '性别',
                prop: 'sex',
                type: 'radio',
                value: 1,
                button: true,
                data: [{
                    label: '男',
                    value: 0
                }, {
                    label: '女',
                    value: 1
                }, {
                    label: '未知',
                    value: ''
                }]
            },
            {prop: "state", label: "状态", type: "switch", required: true, active: "启用", inactive: "禁用", value: true},
            {
                prop: "subject", label: "年级", type: "select",
                value: "2",
                required: true,
                placeholder: "请选择年级",
                "data": [{"label": "一年级", "value": "1"}, {"label": "二年级", "value": "2"}, {
                    "label": "三年级",
                    "value": "3"
                }],
                "labelName": "label",
                "valueName": "value",
                "multiple": false,
                "multiple-limit": "0"
            },
            {
                prop: "orgtype", label: "机构类型机构类型", type: "select",
                value: "",
                required: true,
                placeholder: "请选择机构类型",
                "url": "/api/blade-system/dict/dictionary?code=org_category",
                "labelName": "label",
                "valueName": "value",
                "multiple": true,
                "multiple-limit": "0"
            }
        ],
        shape: {
            height: 40
        }
    },
    {
        id: "pipeline",
        name: "算法组件",
        description: "",
        inherit: 'simple',
        group: "base",
        fields: [
            {prop: "fullname", label: "姓名", value: "", placeholder: "请输入姓名"},
            {prop: "el", label: "EL表达式", value: "", placeholder: "请输入"}
        ]
    },
    {
        id: "tests",
        name: "数据库组件",
        description: "",
        group: "base",
        type: "database",
        icon:"fa fa-microchip",
        fields: [
            {prop: "el", label: "EL表达式", value: "", placeholder: "请输入"}
        ]
    },
    {
        id: "fork",
        name: "判断节点",
        description: "",
        group: "base",
        fields: [
            {prop: "el", label: "EL表达式", value: "", placeholder: "请输入"}
        ]
    }];
