import ajax, { Result, SearchArgs, SearchResult, SetStatusArgs } from '@/api/ajax'



// C创建 U修改 D删除 L列表 V查看 G下载 I导入 E导出 A审批 R撤销 S设置 T标签
// C上传 U修改(复制、移动、重命名) D删除 L列表 G下载 S设置 A锁定 R解锁
export interface WebDav {
    id: string;
    name: string;
    account: string;
    token: string;
    home: string;
    is_disable: number;
    perms: string;
    remark: string;
    created_at: string;
    updated_at: string;
    created_by: number;
    updated_by: number;
}



const baseUrl = '/nas/webdav';
export class WebDavApi {
    save(user: WebDav) {
        return ajax.post<Result<any>>(baseUrl + '/save', user)
    }

    search(args: SearchArgs) {
        return ajax.get<Array<WebDav>>(baseUrl + '/list', args)
    }

    load(id: string) {
        return ajax.get<WebDav>(baseUrl + '/load/' + id)
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

export default new WebDavApi