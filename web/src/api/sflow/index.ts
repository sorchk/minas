import ajax, { Result, SearchArgs, SearchResult } from '@/api/ajax'

// 作业流程数据接口
export interface SFlow {
    id: string;
    name: string;
    content: string; 
    type: string; 
    log_level: number;
    last_status: number;
    last_run_time: string;
    log_keep_num: number;
    project_dir_id: string;
    remark: string;
    is_disable: number;
    created_by: number;
    created_at: string;
    updated_by: number;
    updated_at: string;
    created_by_name: string;
    updated_by_name: string;
}

// 状态映射
export const statusMapping = {
    "-1": { info: '失败', type: 'error' },
    "0": { info: '未执行', type: 'default' },
    "1": { info: '成功', type: 'success' },
} as any

// API基础路径
const baseUrl = '/sflow/sflow';

// 作业流程API类
export class SFlowApi {
    saveContent(data: { id: number; content: any; }) {
        return ajax.post<SFlow>(baseUrl + '/saveContent', data)
    }
    /**
     * 搜索作业流程列表
     * @param args 搜索参数
     */
    search(args: SearchArgs) {
        return ajax.search<SFlow>(baseUrl + '/list', args)
    }

    /**
     * 加载单个作业流程详情
     * @param id 作业流程ID
     */
    load(id: number) {
        return ajax.get<SFlow>(baseUrl + '/load/' + id)
    }

    /**
     * 保存作业流程（新建或更新）
     * @param data 作业流程数据
     */
    save(data: SFlow) {
        return ajax.post<SFlow>(baseUrl + '/save', data)
    }

    /**
     * 删除作业流程
     * @param id 作业流程ID
     */
    delete(id: number) {
        return ajax.post(baseUrl + '/delete/' + id, {})
    }

    /**
     * 启用作业流程
     * @param id 作业流程ID
     */
    enable(id: number) {
        return ajax.post(baseUrl + '/enable/' + id, {})
    }

    /**
     * 禁用作业流程
     * @param id 作业流程ID
     */
    disable(id: number) {
        return ajax.post(baseUrl + '/disable/' + id, {})
    }

    /**
     * 执行作业流程
     * @param id 作业流程ID
     */
    exec(id: number) {
        return ajax.post(baseUrl + '/exec/' + id, {})
    }
}

// 导出API实例
export default new SFlowApi()