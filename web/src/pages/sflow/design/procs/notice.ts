export const notice_components = [{
    id: "EmailNotify",
    name: "邮件通知",
    icon: "fa fa-envelope-o",
    description: "",
    group: "notice",
    fields: [
        {prop: "tomail", label: "接收人", value: "", placeholder: "sorc@qq.com:昵称"},
        {prop: "cc", label: "抄送", type: 'inputs', value: [], placeholder: "sorc@qq.com:昵称"},
        {prop: "bcc", label: "密送", type: 'inputs', value: [], placeholder: "sorc@qq.com:昵称"},
        {prop: "subject", label: "标题", value: "", placeholder: ""},
        {prop: "content", label: "正文", type: 'editor', lang: 'html', value: "", placeholder: ""},
        {
            prop: "files",
            label: "附件",
            type: 'map', value: [],
            keyLabel: "文件名",
            valueLabel: "文件路径",
            placeholder: "例如:相片.jpg=D:/pic/my.jpg"
        },
        {prop: "frommail", label: "发送人", value: "", placeholder: "sorc@qq.com:昵称"},
        {prop: "ssl", label: "ssl支持", type: "switch", active: "是", inactive: "否", value: false},
        {prop: "host", label: "服务器地址", value: "", placeholder: "请输入"},
        {prop: "port", label: "端口", value: "", placeholder: "默认为25，支持ssl默认为465"},
        {prop: "user", label: "认证帐号", value: "", placeholder: "默认为发件人邮箱地址"},
        {prop: "passwd", label: "认证密码", value: "", placeholder: "如果为空，则不进行认证"}
    ]
}, {
    id: "MqNotify",
    name: "消息队列通知",
    icon: "fa fa-exchange",
    description: "",
    group: "notice",
    fields: [
        {
            prop: "mqid", label: "队列服务器",
            type: 'select', value: "",
            url: "/rs/srv/mqtopic/select",
            labelName: "name",
            valueName: "id",
            multiple: false,
            data: [{
                value: 'test',
                label: "test"
            }], placeholder: "请选择"
        },
        {prop: "msg", label: "消息内容", type: 'editor', lang: 'json', value: "", placeholder: ""}
    ]
}, {
    id: "ApiNotify",
    name: "接口通知",
    icon: "fa fa-internet-explorer",
    description: "",
    group: "notice",
    fields: [
        {prop: "url", label: "请求地址", value: "", placeholder: "http://test.com/getInfo"},
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
            value: "application/x-www-form-urlencoded",
            button: true,
            data: [{
                value: 'application/json',
                label: "json"
            }, {
                value: 'application/x-www-form-urlencoded',
                label: "x-www-form-urlencoded"
            }, {
                value: 'multipart/form-data',
                label: "form-data"
            }, {
                value: 'none',
                label: "无"
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
            type: "map", value: [],
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
        {prop: "timeout", label: "超时", value: "", placeholder: "设置请求超时，默认10000毫秒", help: '单位毫秒'},
        {prop: "supportCookie", label: "支持会话", type: "switch", active: "是", inactive: "否", value: false}
    ]
}];
