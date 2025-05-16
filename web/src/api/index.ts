 
import ajax from '@/api/ajax'
import axios from "axios";
import { baseUrl } from '@/config';
import { useMessage, useDialog } from "naive-ui";
const request = ajax.orequest;
export class RsApi {
    public path: string;

    constructor(module: string) {
        this.path = module;
    }

    list(params?: any): Promise<any> {
        if (params) {
            if (params.pageNum) {
                params.page = params.pageNum - 1;
                delete params.pageNum;
            }
            if (params.pageSize) {
                params.size = params.pageSize;
                delete params.pageSize;
            }
        }
        return request({
            url: '/rs' + this.path + '/list',
            method: 'get',
            params: params || {}
        });
    }
    postList(params?: any): Promise<any> {
        if (params) {
            if (params.pageNum) {
                params.page = params.pageNum - 1;
                delete params.pageNum;
            }
            if (params.pageSize) {
                params.size = params.pageSize;
                delete params.pageSize;
            }
        }
        return request({
            url: '/rs' + this.path + '/list',
            method: 'post',
            data: params || {}
        });
    }
    page(params: any): Promise<Array<any>> {
        if (params.pageNum) {
            params.page = params.pageNum - 1;
            delete params.pageNum;
        }
        if (params.pageSize) {
            params.size = params.pageSize;
            delete params.pageSize;
        }
        return request({
            url: '/rs' + this.path,
            method: 'get',
            params: params || {}
        }).then(response => response.data);
    }
    postCount(params?: object): any {
        return request({
            url: '/rs' + this.path + '/count',
            method: 'post',
            data: params || {}
        }).then(response => response.data);
    }
    count(params?: object): any {
        return request({
            url: '/rs' + this.path + '/count',
            method: 'get',
            params: params || {}
        }).then(response => response.data);
    }

    load(id: string | null): any {
        if (id == null || undefined) {
            return new Error("参数id不允许空");
        }
        return request({
            url: '/rs' + this.path + '/' + id,
            method: 'get'
        }).then(response => response.data);
    }

    view(id: string | null): any {
        if (id == null || undefined) {
            return new Error("参数id不允许空");
        }
        return request({
            url: '/rs' + this.path + '/view/' + id,
            method: 'get'
        }).then(response => response.data);
    }

    tree(params?: object): Promise<Array<any>> {
        return request({
            url: '/rs' + this.path + '/tree',
            method: 'get',
            params: params || {}
        }).then(response => response.data);
    }

    select(params?: object): Promise<Array<any>> {
        return request({
            url: '/rs' + this.path + '/select',
            method: 'get',
            params: params || {}
        }).then(response => response.data);
    }

    top(num: number): Promise<Array<any>> {
        return request({
            url: '/rs' + this.path + '/top/' + num,
            method: 'get'
        }).then(response => response.data);
    }

    save(params?: object): any {
        return request({
            url: '/rs' + this.path + '/save',
            method: 'post',
            data: params || {}
        }).then(response => response.data);
    }

    add(params?: object): any {
        return request({
            url: '/rs' + this.path + '/add',
            method: 'post',
            data: params || {}
        }).then(response => response.data);
    }

    

    remove(id: string | null, params?: object): any {
        if (id == null || undefined) {
            return new Error("参数id不允许空");
        }
        return request({
            url: '/rs' + this.path + '/remove/' + id,
            method: 'post',
            data: params || {}
        });
    }

    move(id: string, target: string, type: string, snlist: object): any {
        return request({
            url: '/rs' + this.path + '/move/' + id + '/' + target + '/' + type,
            method: 'post',
            data: snlist
        });
    }

    copy(id: string): any {
        return request({
            url: '/rs' + this.path + '/copy/' + id,
            method: 'post'
        });
    }

    enable(id: string): any {
        return request({
            url: '/rs' + this.path + '/enable/' + id,
            method: 'post'
        });
    }

    disable(id: string): any {
        return request({
            url: '/rs' + this.path + '/disable/' + id,
            method: 'post'
        });
    }

    start(id: string): any {
        return request({
            url: '/rs' + this.path + '/start/' + id,
            method: 'post'
        });
    }

    stop(id: string): any {
        return request({
            url: '/rs' + this.path + '/stop/' + id,
            method: 'post'
        });
    }

    revise(id: string): any {
        return request({
            url: '/rs' + this.path + '/revise/' + id,
            method: 'post'
        });
    }

    publish(id: string): any {
        return request({
            url: '/rs' + this.path + '/publish/' + id,
            method: 'post'
        });
    }

    compare(id: string): any {
        return request({
            url: '/rs' + this.path + '/compare/' + id,
            method: 'get'
        });
    }

    test(id: string, params?: any): any {
        return request({
            url: '/rs' + this.path + '/test/' + id,
            method: 'post',
            data: params || {}
        });
    }

}

export const rsapi = {
    build(module: string) {
        return new RsApi(module);
    },
    fileSizeShow(byteSize: number) {
        const unitArr = new Array("B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB");
        const index = Math.floor(Math.log(byteSize) / Math.log(1024));
        let size = byteSize / Math.pow(1024, index);
        return size.toFixed(2) + unitArr[index];
    },
    
    get(url: string, params?: object): any {
        return request({
            url: url,
            method: 'get',
            params: params || {}
        });
    },
    post(url: string, params?: object): any {
        return request({
            url: url,
            method: 'post',
            data: params || {}
        });
    },
     
    treeToList(data: any, start: number, childrenKey: string, snKey: string): Array<any> {
        if (!childrenKey) {
            childrenKey = "children";
        }
        if (snKey === null || snKey === undefined) {
            snKey = "sn";
        }
        if (!start) {
            start = 1;
        }
        let list = [];
        if (data) {
            for (let i in data) {
                let item = data[i];
                if (snKey) {
                    item[snKey] = start + list.length;
                }
                list.push(item);
                if (item[childrenKey]) {
                    list.push(...this.treeToList(item[childrenKey], start + list.length, childrenKey, snKey))
                }
            }
        }
        return list;
    },
    listToTree(list: any, treeModel?: any): any {
        const keyMapping: any = {};
        treeModel = this.extend(treeModel || {}, {
            id: "id",
            parentId: "parentId",
            children: "children",
            disabled: "disabled",
            isLeaf: "",
            enabledParent: false
        });
        const tmpList = this.deepClone(list);
        //数据的主键数据与数据映射
        for (let i in tmpList) {
            let item = tmpList[i];
            keyMapping[item[treeModel.id]] = item;
        }
        const results = [];
        for (let i in tmpList) {
            let item = tmpList[i];
            let parentId = item[treeModel.parentId];
            if (parentId !== null && parentId !== undefined && parentId !== '' && parentId != 0) {
                let parent = keyMapping[parentId];
                if (parent) {
                    if (!parent[treeModel.children]) {
                        parent[treeModel.children] = [];
                    }
                    parent[treeModel.children].push(item);
                    if (treeModel.enabledParent) {
                        parent[treeModel.disabled] = false;
                    }
                }
            } else {
                results.push(item);
            }
        }
        return results;
    },
    listToMapping(list: Array<any>, key: string): any {
        const result = {} as any;
        if (list) {
            list.forEach(item => {
                result[item[key]] = item;
            })
        }
        return result;
    },
    groupToMapping(list: Array<any>, groupKey: string): any {
        const result = {} as any;
        if (list) {
            list.forEach(item => {
                if(!result[item[groupKey]]){
                    result[item[groupKey]]=new Array();
                }
                result[item[groupKey]].push(item);
            })
        }
        return result;
    },
    leftOutJoinString(leftList: Array<any>, rightList: Array<any>, leftKey: string, rightKey: string, newKey: string,strKey:string): any {
        const rightMapping = this.listToMapping(rightList, rightKey);
        for (let i in leftList) {
            const leftItem = leftList[i];
            const rightItem = rightMapping[leftItem[leftKey]];
            if (rightItem) {
                leftItem[newKey] = rightItem[strKey];
            }
        }
        return leftList;
    },
    innerJoinString(leftList: Array<any>, rightList: Array<any>, leftKey: string, rightKey: string, newKey: string,strKey:string): any {
        const rightMapping = this.listToMapping(rightList, rightKey);
        const result=new Array();
        for (let i in leftList) {
            const leftItem = leftList[i];
            const rightItem = rightMapping[leftItem[leftKey]];
            if (rightItem) {
                leftItem[newKey] = rightItem[strKey];
                result.push(leftItem);
            }
        }
        return result;
    },
    leftOutJoinObject(leftList: Array<any>, rightList: Array<any>, leftKey: string, rightKey: string, newKey: string): any {
        const rightMapping = this.listToMapping(rightList, rightKey);
        for (let i in leftList) {
            const leftItem = leftList[i];
            const rightItem = rightMapping[leftItem[leftKey]];
            if (rightItem) {
                leftItem[newKey] = rightItem;
            }
        }
        return leftList;
    },
    clone(data: any): any {
        if (data === null || data === undefined) {
            return data;
        }
        let obj: any;
        if (Array.isArray(data)) {
            obj = [];
        } else if ((typeof data) === 'object') {
            obj = {};
        } else {
            //不再具有下一层次
            return data;
        }
        if (Array.isArray(data)) {
            for (let i = 0, len = data.length; i < len; i++) {
                obj.push(data[i]);
            }
        } else if ((typeof data) === 'object') {
            for (let key in data) {
                obj[key] = data[key];
            }
        }
        return obj;
    },
    /**
     * 对象深拷贝
     */
    deepClone(data: any): any {
        if (data === null || data === undefined) {
            return data;
        }
        let obj: any;
        if (Array.isArray(data)) {
            obj = [];
        } else if ((typeof data) === 'object') {
            obj = {};
        } else {
            //不再具有下一层次
            return data;
        }
        if (Array.isArray(data)) {
            for (let i = 0, len = data.length; i < len; i++) {
                obj.push(this.deepClone(data[i]));
            }
        } else if ((typeof data) === 'object') {
            for (let key in data) {
                obj[key] = this.deepClone(data[key]);
            }
        }
        return obj;
    },
    extend(objA: any, objB: any): object {
        if (objB === null || objB === undefined) {
            return objA;
        } else if (objA === null || objA === undefined) {
            objA = this.deepClone(objB);
            return objA;
        } else {
            for (let key in objB) {
                let bv = objB[key];
                if (bv === null || bv === undefined) {
                    continue;
                }
                let av = objA[key];
                if (av === null || av === undefined) {
                    objA[key] = this.deepClone(bv);
                }
            }
            return objA;
        }
    },
    // getObjType(obj: object): string {
    //     let toString = Object.prototype.toString;
    //     let map = {
    //         '[object Boolean]': 'boolean',
    //         '[object Number]': 'number',
    //         '[object String]': 'string',
    //         '[object Function]': 'function',
    //         '[object Array]': 'array',
    //         '[object Date]': 'date',
    //         '[object RegExp]': 'regExp',
    //         '[object Undefined]': 'undefined',
    //         '[object Null]': 'null',
    //         '[object Object]': 'object'
    //     };
    //     if (obj instanceof Element) {
    //         return 'element';
    //     }
    //     return map[toString.call(obj)];
    // }
}

export default rsapi;
