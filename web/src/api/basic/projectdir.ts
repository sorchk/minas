import ajax, { Result, SearchArgs, SearchResult } from '@/api/ajax'

// 定义项目目录类型接口
export interface ProjectDirItem {
  id: number;
  name: string;
  parent_id: number | undefined;
  remark: string;
  is_disable: number;
  has_child: boolean;
  children?: ProjectDirItem[];
  created_at?: string;
  updated_at?: string;
  created_by?: number;
  updated_by?: number;
}

const basePath = '/basic/project-dir';

export class ProjectDirApi {
  /**
   * 获取项目目录列表（分页）
   * @param args 分页查询参数
   * @returns 包含列表数据和总数的响应
   */
  search(args: SearchArgs) {
    return ajax.get<SearchResult<ProjectDirItem>>(`${basePath}/list`, args)
  }

  /**
   * 获取项目目录树
   * @returns 树形结构的目录数据
   */
  getTree() {
    return ajax.get<Result<ProjectDirItem[]>>(`${basePath}/tree`)
  }

  /**
   * 获取指定父目录下的子目录
   * @param parentId 父目录ID，0表示顶级目录
   * @returns 子目录列表
   */
  getChildren(parentId: number) {
    return ajax.get<Result<ProjectDirItem[]>>(`${basePath}/children`, { parent_id: parentId })
  }

  /**
   * 根据ID加载项目目录
   * @param id 目录ID
   * @returns 目录详情
   */
  find(id: number) {
    return ajax.get<Result<ProjectDirItem>>(`${basePath}/load/${id}`)
  }

  /**
   * 保存项目目录（创建或更新）
   * @param data 目录数据
   * @returns 保存结果
   */
  save(data: ProjectDirItem) {
    return ajax.post<Result<any>>(`${basePath}/save`, data)
  }

  /**
   * 启用项目目录
   * @param id 目录ID
   * @returns 操作结果
   */
  enable(id: number) {
    return ajax.post<Result<any>>(`${basePath}/enable/${id}`)
  }

  /**
   * 禁用项目目录
   * @param id 目录ID
   * @returns 操作结果
   */
  disable(id: number) {
    return ajax.post<Result<any>>(`${basePath}/disable/${id}`)
  }

  /**
   * 删除项目目录
   * @param id 目录ID
   * @returns 操作结果
   */
  delete(id: number) {
    return ajax.post<Result<any>>(`${basePath}/delete/${id}`)
  }
}

export default new ProjectDirApi