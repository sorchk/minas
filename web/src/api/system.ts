import ajax, { Result } from '@/api/ajax'


const baseUrl = '/system';
export class SystemApi {
    checkState() {
        return ajax.get<any>(baseUrl+'/check-state')
    }

    init(user: any) {
        return ajax.post<Result<any>>(baseUrl+'/init', user)
    }

    version() {
        return ajax.get<any>(baseUrl+'/version')
    }

}

export default new SystemApi
