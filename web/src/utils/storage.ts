import Cookies from 'js-cookie';
/**
 * window.localStorage 浏览器永久缓存
 * @method set 设置永久缓存
 * @method get 获取永久缓存
 * @method remove 移除永久缓存
 * @method clear 移除全部永久缓存
 */
export const Local = {
    formatKey(key:string){
        return "minas:"+key;
    },
    // 设置永久缓存
    set(key: string, val: any) {
        window.localStorage.setItem(Local.formatKey(key), JSON.stringify(val));
    },
    // 获取永久缓存
    get(key: string) {
        let json = <string>window.localStorage.getItem(Local.formatKey(key));
        return JSON.parse(json);
    },
    // 移除永久缓存
    remove(key: string) {
        window.localStorage.removeItem(Local.formatKey(key));
    },
    // 移除全部永久缓存
    clear() {
        window.localStorage.clear();
    },
};

/**
 * window.sessionStorage 浏览器临时缓存
 * @method set 设置临时缓存
 * @method get 获取临时缓存
 * @method remove 移除临时缓存
 * @method clear 移除全部临时缓存
 */
export const Session = {
    // 设置临时缓存
    set(key: string, val: any) {
        if (key === 'token') return Cookies.set(Local.formatKey(key), val);
        window.sessionStorage.setItem(Local.formatKey(key), JSON.stringify(val));
    },
    // 获取临时缓存
    get(key: string) {
        if (key === 'token') return Cookies.get(Local.formatKey(key));
        let json = <string>window.sessionStorage.getItem(Local.formatKey(key));
        return JSON.parse(json);
    },
    // 移除临时缓存
    remove(key: string) {
        if (key === 'token') return Cookies.remove(Local.formatKey(key));
        window.sessionStorage.removeItem(Local.formatKey(key));
    },
    // 移除全部临时缓存
    clear() {
        Cookies.remove(Local.formatKey('token'));
        window.sessionStorage.clear();
    },
};
