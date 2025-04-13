import ajax, { Result } from '@/api/ajax'
const baseUrl = '/term';
export class TermApi {
    auth(termId: string) {
        return ajax.get<Result<any>>(baseUrl + '/auth/token/' + termId)
    }

}

export default new TermApi