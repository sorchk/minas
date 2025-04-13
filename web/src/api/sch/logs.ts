import ajax, { Result, SearchArgs, SearchResult, SetStatusArgs } from '@/api/ajax'

// C创建 U修改 D删除 L列表 V查看 G下载 I导入 E导出 A审批 R撤销 S设置 T标签
// C上传 U修改(复制、移动、重命名) D删除 L列表 G下载 S设置 A锁定 R解锁
export interface SchLog {
    id: string;
    task_id: string;
    status: string;
    log_text: string;
    start_time: string;
    end_time: string;
}
export const statusMapping = {
    "-2": { info: '中断', type: 'warning' },
    "-1": { info: '失败', type: 'error' },
    "0": { info: '运行中', type: 'info' },
    "1": { info: '成功', type: 'success' },
} as any

const baseUrl = '/sch/log';
export class SchLogApi {
    search(args: SearchArgs) {
        return ajax.search<SchLog>(baseUrl + '/list', args)
    }

    load(id: string) {
        return ajax.get<SchLog>(baseUrl + '/load/' + id)
    }
}

export default new SchLogApi