import ajax, { Result, SearchArgs, SearchResult } from '@/api/ajax'
  
const baseUrl = '/dagflow';
// 作业流程API类
export class SFlowApi {
    // API基础路径
    
    /**
     * 执行作业流程
     * @param id 作业流程ID
     */
    execute(id: number) {
        return ajax.post(baseUrl + '/execute/' + id, {})
    }
    /**
     * 验证作业流程
     * @param id 作业流程ID
     */
   validate(id: number) {
    return ajax.post(baseUrl + '/validate/' + id, {})
}
  
    /**
     * 执行作业流程
     * @param id 作业流程ID
     */
    debug(id: number) {
        return ajax.post(baseUrl + '/debug/' + id, {})
    }


    /**
     * 执行作业流程
     */
    handlers() {
        return ajax.get(baseUrl + '/handlers', {})
    }
}

// 导出API实例
export default new SFlowApi()