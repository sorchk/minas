export const check_components = [
    {
        id: "BoolCheck",
        name: "布尔检查",
        icon: 'fa fa-paragraph',
        description: "",
        group: "check",
        fields: [
            {
                prop: "el", label: "布尔表达式", type: "inputs", value: [], placeholder: "请输入，支持EL表达式，结果是布尔值",
                help: "如果输入空或者非布尔表达，则默认返回false"
            },
            {
                prop: "isAll", label: "检测", type: "switch", active: "全部", inactive: "任意", value: true,
                help: "全部：要求所有布尔表达式都返回真则通过，任意：要求任意一个布尔表达式返回真则通过"
            },
            { prop: "showError", label: "抛异常", type: "switch", active: "是", inactive: "否", value: true },
            { prop: "errorInfo", label: "错误信息", value: "", placeholder: "错误信息" }
        ]
    },
    {
        id: "RegexCheck",
        name: "正则表达式",
        icon: 'fa fa-asterisk',
        description: "",
        group: "check",
        fields: [
            { prop: "data", label: "被测试数据", value: "", placeholder: "请输入被测试数据，支持EL表达式，结果是文本值" },
            { prop: "regex", label: "正则表达式", value: "", placeholder: "请输入正则表达式" },
            {
                prop: "isCase", label: "区分大小写", type: "switch", active: "是", inactive: "否", value: true
            },
            {
                prop: "isAll", label: "匹配方式", type: "switch", active: "完全匹配", inactive: "包含", value: true
            },
            {
                prop: "showError", label: "抛异常", type: "switch", active: "是", inactive: "否", value: true
            },
            { prop: "errorInfo", label: "错误信息", value: "", placeholder: "错误信息" }
        ]
    }];
