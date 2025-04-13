import ajax, { Result, SearchArgs, SearchResult, SetStatusArgs } from '@/api/ajax'


export interface AuthUser {
    error: number,
    token: string;
    id: string;
    username: string;
    name: string;
    avatar: string;
    perms: string[];
    ip: string;
}

export interface ModifyPasswordArgs {
    password: string;
    newpassword: string;
}

const baseUrl = '/basic/user';
export class UserApi {
    login(args: any) {
        return ajax.post('/login', args)
    }

    save(user: any) {
        return ajax.post<Result<any>>('/user/save', user)
    }

    find(id: string) {
        return ajax.get<any>(baseUrl+'/load', { id })
    }

    profile() {
        return ajax.get<any>(baseUrl+'/profile')
    }

    search(args: SearchArgs) {
        return ajax.get<SearchResult<any>>(baseUrl+'/list', args)
    }

    setStatus(args: SetStatusArgs) {
        return ajax.post<Result<any>>(baseUrl+'/set-status', args)
    }

    delete(id: string, name: string) {
        return ajax.post<Result<any>>(baseUrl+'/delete', { id, name })
    }
    modifyPassword(args: ModifyPasswordArgs) {
        return ajax.post<Result<any>>(baseUrl+'/modify-password', args)
    }

    modifyProfile(user: any) {
        return ajax.post<Result<any>>(baseUrl+'/modify-profile', user)
    }
    enableMfa(args: any) {
        return ajax.post<Result<any>>(baseUrl+'/mfa/enable', args)
    }
    disableMfa() {
        return ajax.post<Result<any>>(baseUrl+'/mfa/disable')
    }
}

export default new UserApi