import ajax, { Result, SearchArgs, SearchResult } from '@/api/ajax'

// 作业流程日志数据接口
export interface SFlowLog {
    id: string;
    sflow_id: string;
    status: string;
    log_text: string;
    start_time: string;
    end_time: string;
}

// 状态映射
export const statusMapping = {
    "-2": { info: '中断', type: 'warning' },
    "-1": { info: '失败', type: 'error' },
    "0": { info: '运行中', type: 'info' },
    "1": { info: '成功', type: 'success' },
} as any

// API基础路径
const baseUrl = '/sflow/log';

// 作业流程日志API类
export class SFlowLogApi {
    /**
     * 搜索作业流程日志列表
     * @param args 搜索参数
     */
    search(args: SearchArgs) {
        return ajax.search<SFlowLog>(baseUrl + '/list', args)
    }

    /**
     * 加载单个作业流程日志详情
     * @param id 日志ID
     */
    load(id: string) {
        return ajax.get<SFlowLog>(baseUrl + '/load/' + id)
    }
}

// 导出API实例
export default new SFlowLogApi()