/**
 * 定义接口来定义对象的类型
 */

export class KVMap {
	public key: string = '';
	public value: string = '';
}
export class ParamValue {
	public name: string = '';
	public value: string = '';
}
export class ParamsMap {
	public key: string = '';
	public value: string = '';
	public type: string = '';
	public required: boolean = false;
	public defval: string = '';
	public remark: string = '';
	public active: string = '';
	public inactive: string = '';
}
export class CrawlerXpathMap {
	public key: string = '';
	public xpath: string = '';
	public type: string = 'string';
	public template: string = '';
	public regex: string = '';
	public required: boolean = false;
	public active: string = '';
	public inactive: string = '';
}