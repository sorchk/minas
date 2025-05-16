export const crawler_components = [
    {
        id: "StaticCrawler",
        name: "静态爬虫",
        icon: 'iconfont icon-zhongduancanshuchaxun',
        description: "静态爬虫，支持同步渲染页面",
        group: "crawler",
        fields: [
            { prop: "url", label: "请求地址", value: "", placeholder: "http://test.com/getInfo" },
            {
                prop: "heardes",
                label: "请求头",
                type: "map",
                value: [],
                keyLabel: "请求头",
                valueLabel: "值",
                placeholder: "例如:123或${param.id}"
            },
            { prop: "timeout", label: "超时", value: "", placeholder: "设置请求超时，默认10000单位毫秒", help: '单位毫秒' },
            { prop: "cookie", label: "Cookie", value: "", placeholder: "设置Cookie，默认空", help: '设置Cookie' },
            {
                prop: "datas",
                label: "解析数据",
                type: "crawlerXpath",
                value: []
            },
        ]
    },
    {
        id: "DynamicCrawler",
        name: "动态爬虫",
        icon: 'iconfont icon-zhongduancanshu',
        description: "动态爬虫，支持动态异步渲染页面",
        group: "crawler",
        fields: [
            {
                label: '浏览器驱动',
                prop: 'method',
                type: 'radio',
                value: "Chrome",
                button: true,
                data: [{
                    value: 'Chrome',
                    label: "Chrome"
                }, {
                    value: 'Firefox',
                    label: "Firefox"
                }]
            },
            { prop: "url", label: "请求地址", value: "", placeholder: "http://test.com/getInfo" },
            {
                prop: "heardes",
                label: "请求头",
                type: "map",
                value: [],
                keyLabel: "请求头",
                valueLabel: "值",
                placeholder: "例如:123或${param.id}"
            },
            { prop: "timeout", label: "超时", value: "", placeholder: "设置请求超时，默认30000单位毫秒", help: '' },
            { prop: "cookie", label: "Cookie", value: "", placeholder: "设置Cookie，默认空", help: '' },
            {
                prop: "waitFullPage", label: "等待整页内容", type: "switch", active: "等待", inactive: "跳过", value: false,
                help: "可能会消耗更长时间，程序会自动将滚动条滚动到页面底部，等待动态内容加载完成,再进行数据解析。默认不需要开启，程序会自动等待xpath元素加载完成"
            }, {
                prop: "waitAllRequest", label: "等待所有请求", type: "switch", active: "等待", inactive: "跳过", value: false,
                help: "可能会消耗更长时间，等待所有请求完成后再进行数据解析。默认不需要开启，程序会自动等待xpath元素加载完成"
            },
            {
                prop: "datas",
                label: "解析数据",
                type: "crawlerXpath",
                value: []
            },
        ]
    },
    {
        id: "FileDownload",
        name: "文件下载",
        icon: 'fa fa-cloud-download',
        description: "文件下载",
        group: "crawler",
        fields: [
            { prop: "url", label: "请求地址", value: "", placeholder: "http://test.com/getInfo" },
            {
                prop: "heardes",
                label: "请求头",
                type: "map",
                value: [],
                keyLabel: "请求头",
                valueLabel: "值",
                placeholder: "例如:123或${param.id}"
            },
            { prop: "timeout", label: "超时", value: "", placeholder: "设置请求超时，默认600000单位毫秒", help: '' },
            { prop: "cookie", label: "Cookie", value: "", placeholder: "设置Cookie，默认空", help: '' },
            { prop: "savePath", label: "保存路径", value: "", placeholder: "设置保存路径，默认空,使用\\tmp", help: '' },
            { prop: "fileName", label: "文件名", value: "", placeholder: "设置文件名，默认空，从url提取，无法提取则随机生成", help: '' },
            {
                prop: "overwrite", label: "覆盖文件", type: "switch", active: "覆盖", inactive: "跳过", value: true,
                help: ""
            },
            {
                prop: "followRedirects", label: "跟随重定向", type: "switch", active: "是", inactive: "否", value: true,
                help: ""
            },
            { prop: "maxRedirects", label: "最大重定向次数", type: "number", value: "5", placeholder: "设置文件名，默认5，从url提取，无法提取则随机生成", help: '' },


        ]
    },
    {
        id: "Sleep",
        name: "休眠",
        icon: 'iconfont icon-wenducanshu-05',
        description: "程序休眠，等待指定时间",
        group: "crawler",
        fields: [
            { prop: "min", label: "最小休眠时间", value: "0", placeholder: "设置请求超时，默认0单位毫秒,随机休眠的最小时间", help: '' },
            { prop: "max", label: "最大休眠时间", value: "0", placeholder: "设置请求休眠，默认0单位毫秒,随机休眠的最大时间", help: '' },
        ]
    }
];
