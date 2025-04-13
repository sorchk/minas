import ajax, { Result, SearchArgs, SearchResult, SetStatusArgs } from '@/api/ajax'



// C创建 U修改 D删除 L列表 V查看 G下载 I导入 E导出 A审批 R撤销 S设置 T标签
// C上传 U修改(复制、移动、重命名) D删除 L列表 G下载 S设置 A锁定 R解锁
export interface ExternalNas {
    id: string;
    name: string;
    type: string;
    rc_name: string;
    is_adv: number;
    config: any;
    is_sync: boolean;
    is_disable: number;
    remark: string;
    created_at: string;
    updated_at: string;
    created_by: number;
    updated_by: number;
}



const baseUrl = '/nas/external';
export class ExternalNasApi {
    save(user: ExternalNas) {
        return ajax.post<Result<any>>(baseUrl + '/save', user)
    }

    search(args: SearchArgs) {
        return ajax.get<Array<ExternalNas>>(baseUrl + '/list', args)
    }

    load(id: string) {
        return ajax.get<ExternalNas>(baseUrl + '/load/' + id)
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
    listDir(nas: string, path: string) {
        return ajax.get<any>(baseUrl + '/list-dir', { nas, path })
    }
    rcloneApi(uri: string, params: any) {
        return ajax.post<Result<any>>(baseUrl + '/rclone/api/' + uri, params)
    }
}

export default new ExternalNasApi