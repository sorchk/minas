import ajax, { Result, SearchArgs, SearchResult, SetStatusArgs } from '@/api/ajax'
import { t } from '@/locales'


// C创建 U修改 D删除 L列表 V查看 G下载 I导入 E导出 A审批 R撤销 S设置 T标签
// C上传 U修改(复制、移动、重命名) D删除 L列表 G下载 S设置 A锁定 R解锁
export interface SchTask {
    id: string;
    name: string;
    type: string;
    cron: string;
    last_status: number;
    last_run_time: string;
    next_run_time: string;
    log_keep_num: number;
    script: string;
    remark: string;
    is_disable: number;
    created_at: string;
    updated_at: string;
    created_by: number;
    updated_by: number;
}
export const runStatusMapping = {
    "-2": { info: t('schtask.status.abnormal_exit'), type: 'error' },
    "-1": { info: t('schtask.status.failed'), type: 'error' },
    "0": { info: t('schtask.status.waiting'), type: 'info' },
    "1": { info: t('schtask.status.waiting'), type: 'success' },
    "2": { info: t('schtask.status.waiting'), type: 'success' },
} as any
export const typeMapping = {
    "SHELL": t('schtask.type.shell'),
    "FILE_BACKUP": t('schtask.type.file_backup'),
    "FILE_CLEAN": t('schtask.type.file_clean'),
    "JOB_TASK": t('schtask.type.job_task'),
} as any
export const typeOptions = [
    { "label": t('schtask.type.shell'), "value": "SHELL" },
    { "label": t('schtask.type.file_backup'), "value": "FILE_BACKUP" },
    { "label": t('schtask.type.file_clean'), "value": "FILE_CLEAN" },
    // { "label": t('schtask.type.job_task'), "value": "JOB_TASK" },
]
//增量备份 只备份一份，完整备份按时间新建目录 多份备份
export const backupTypeOptions = [
    { "label": t('schtask.backup_type.incremental'), "value": "1" },
    { "label": t('schtask.backup_type.one_way_mirror'), "value": "2" },
    // { "label": t('schtask.backup_type.two_way_sync'), "value": "3" },
    { "label": t('schtask.backup_type.full_backup'), "value": "4" },
]
const baseUrl = '/sch/schtask';
export class SchTaskApi {

    exec(id: any) {
        return ajax.post<Result<any>>(baseUrl + '/exec/' + id, {})
    }
    save(user: SchTask) {
        return ajax.post<Result<any>>(baseUrl + '/save', user)
    }

    search(args: SearchArgs) {
        args.columns = JSON.stringify("id,type,name,cron,last_status,last_run_time,next_run_time,log_keep_num,is_disable,created_at,updated_at,created_by,updated_by".split(","))
        return ajax.search<SchTask>(baseUrl + '/list', args)
    }

    load(id: string) {
        return ajax.get<SchTask>(baseUrl + '/load/' + id)
    }
    disable(id: string) {
        return ajax.post<Result<any>>(baseUrl + '/disable/' + id)
    }
    enable(id: string) {
        return ajax.post<Result<any>>(baseUrl + '/enable/' + id)
    }

    delete(id: string) {
        return ajax.post<Result<any>>(baseUrl + '/delete/' + id)
    }
    listDir(path: string) {
        return ajax.get<any>(baseUrl + '/list-dir', { path })
    }
}

export default new SchTaskApi