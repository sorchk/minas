package rclone

import (
	"encoding/json"
	"fmt"
	"server/utils/data"
	"server/utils/logger"
	"time"
)

// FileObj 文件对象结构体
// 表示从RClone获取的文件或目录条目
type FileObj struct {
	Name    string `json:"Name"`    // 文件或目录名称
	Path    string `json:"Path"`    // 完整路径
	ModTime string `json:"ModTime"` // 修改时间
	IsDir   bool   `json:"IsDir"`   // 是否是目录
	Size    int64  `json:"Size"`    // 文件大小（字节）
}

// List 列出指定目录中的所有文件和目录
// 支持过滤器和时间范围过滤
// 参数:
//   - dir: 要列出内容的目录路径
//   - includes: 包含的文件模式列表（通配符格式）
//   - excludes: 排除的文件模式列表（通配符格式）
//   - minAge: 最小文件年龄，只返回比这个时间更早的文件
//   - maxAge: 最大文件年龄，只返回比这个时间更新的文件
//
// 返回值:
//   - []FileObj: 符合条件的文件和目录对象列表
//   - error: 操作过程中的错误，成功则为nil
func List(dir string, includes []string, excludes []string, minAge time.Duration, maxAge time.Duration) ([]FileObj, error) {
	// 构建过滤器参数
	filter := data.Map{"IncludeRule": includes, "ExcludeRule": excludes}
	if minAge > 0 {
		filter["MinAge"] = minAge
	}
	if maxAge > 0 {
		filter["MaxAge"] = maxAge
	}

	// 构建API请求参数
	params := data.Map{"fs": dir, "remote": "",
		"_filter": filter,
		"_async":  false, // 同步执行，等待结果
		"opt": data.Map{
			"noMimeType": true, // 不获取MIME类型信息以提高性能
		},
	}

	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}

	// 记录操作到日志
	logger.LOG.Infof("----------List------------:%s\n", params)

	// 调用RClone API执行列表操作
	code, body, err := Api("operations/list", params)
	if err != nil {
		return nil, err
	} else {
		if code == 200 {
			// 从返回数据中提取文件列表
			list := body.(map[string]interface{})["list"]
			// 将结果转换为FileObj结构体切片
			jsonData, _ := json.Marshal(list)
			var fileList []FileObj
			json.Unmarshal(jsonData, &fileList)
			return fileList, nil
		} else {
			// 返回错误信息
			return nil, fmt.Errorf("%v", body.(map[string]interface{})["error"])
		}
	}
}

// ListDir 仅列出指定目录中的子目录（不包括文件）
// 参数与List函数类似，但仅返回目录
// 参数:
//   - dir: 要列出内容的目录路径
//   - includes: 包含的目录模式列表（通配符格式）
//   - excludes: 排除的目录模式列表（通配符格式）
//   - minAge: 最小目录年龄，只返回比这个时间更早的目录
//   - maxAge: 最大目录年龄，只返回比这个时间更新的目录
//
// 返回值:
//   - []FileObj: 符合条件的目录对象列表
//   - error: 操作过程中的错误，成功则为nil
func ListDir(dir string, includes []string, excludes []string, minAge time.Duration, maxAge time.Duration) ([]FileObj, error) {
	// 构建过滤器参数
	filter := data.Map{"IncludeRule": includes, "ExcludeRule": excludes}
	if minAge > 0 {
		filter["MinAge"] = minAge
	}
	if maxAge > 0 {
		filter["MaxAge"] = maxAge
	}

	// 构建API请求参数
	params := data.Map{"fs": dir, "remote": "",
		"_filter": filter,
		"_async":  false, // 同步执行，等待结果
		"opt": data.Map{
			"noMimeType": true, // 不获取MIME类型信息以提高性能
			"dirsOnly":   true, // 仅列出目录
		},
	}

	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}

	// 记录操作到日志
	logger.LOG.Infof("----------ListDir------------:%s\n", params)

	// 调用RClone API执行列表操作
	code, body, err := Api("operations/list", params)
	if err != nil {
		return nil, err
	} else {
		if code == 200 {
			// 从返回数据中提取目录列表
			list := body.(map[string]interface{})["list"]
			// 将结果转换为FileObj结构体切片
			jsonData, _ := json.Marshal(list)
			var fileList []FileObj
			json.Unmarshal(jsonData, &fileList)
			return fileList, nil
		} else {
			// 返回错误信息
			return nil, fmt.Errorf("%v", body.(map[string]interface{})["error"])
		}
	}
}

// ListFile 仅列出指定目录中的文件（不包括目录）
// 参数与List函数类似，但仅返回文件
// 参数:
//   - dir: 要列出内容的目录路径
//   - includes: 包含的文件模式列表（通配符格式）
//   - excludes: 排除的文件模式列表（通配符格式）
//   - minAge: 最小文件年龄，只返回比这个时间更早的文件
//   - maxAge: 最大文件年龄，只返回比这个时间更新的文件
//
// 返回值:
//   - []FileObj: 符合条件的文件对象列表
//   - error: 操作过程中的错误，成功则为nil
func ListFile(dir string, includes []string, excludes []string, minAge time.Duration, maxAge time.Duration) ([]FileObj, error) {
	// 构建过滤器参数
	filter := data.Map{"IncludeRule": includes, "ExcludeRule": excludes}
	if minAge > 0 {
		filter["MinAge"] = minAge
	}
	if maxAge > 0 {
		filter["MaxAge"] = maxAge
	}

	// 构建API请求参数
	params := data.Map{"fs": dir, "remote": "",
		"_filter": filter,
		"_async":  false, // 同步执行，等待结果
		"opt": data.Map{
			"noMimeType": true, // 不获取MIME类型信息以提高性能
			"filesOnly":  true, // 仅列出文件
		},
	}

	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}

	// 记录操作到日志
	logger.LOG.Infof("----------ListFile------------:%s\n", params)

	// 调用RClone API执行列表操作
	code, body, err := Api("operations/list", params)
	if err != nil {
		return nil, err
	} else {
		if code == 200 {
			// 从返回数据中提取文件列表
			list := body.(map[string]interface{})["list"]
			// 将结果转换为FileObj结构体切片
			jsonData, _ := json.Marshal(list)
			var fileList []FileObj
			json.Unmarshal(jsonData, &fileList)
			return fileList, nil
		} else {
			// 返回错误信息
			return nil, fmt.Errorf("%v", body.(map[string]interface{})["error"])
		}
	}
}
