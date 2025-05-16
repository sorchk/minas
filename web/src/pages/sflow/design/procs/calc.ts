export const calc_components = [{
  id: "StreamStatistics",
  name: "统计",
  icon: 'fa fa-calculator',
  description: "",
  group: "calc",
  fields: [{
      prop: "data",
      label: "数据",
      placeholder: "",
      help: ''
    },
    {
      prop: "parallel",
      label: "并行计算",
      type: "switch",
      active: "是",
      inactive: "否",
      value: false
    },
    {
      label: '结果类型',
      prop: 'result_type',
      type: 'radio',
      value: "int",
      button: true,
      data: [{
        value: 'int',
        label: "整数"
      }, {
        value: 'long',
        label: "长整数"
      }, {
        value: 'double',
        label: "小数"
      }]
    }, {
      label: '统计方式',
      prop: 'stat_type',
      type: 'radio',
      value: "count",
      button: true,
      data: [{
        value: 'count',
        label: "计数"
      }, {
        value: 'average',
        label: "平均值"
      }, {
        value: 'sum',
        label: "求和"
      }, {
        value: 'max',
        label: "取最大"
      }, {
        value: 'min',
        label: "取最小"
      }]
    },
    {
      prop: "script_text",
      label: "统计脚本",
      type: 'editor',
      lang: 'javascript',
      value: '',
      placeholder: "",
      help: '例如：function test (x){&#xa;    return x.age;&#xa;}'
    }
  ]
}];
