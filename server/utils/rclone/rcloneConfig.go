package rclone

import (
	"encoding/json"
	"errors"
	"server/utils/data"
	"strconv"
)

// ListConfigProviders 获取RClone支持的所有存储提供商列表
// 调用RClone API获取支持的存储类型信息
// 返回值:
//   - interface{}: 包含所有支持的存储提供商信息的数据结构
//   - error: 操作过程中的错误，成功则为nil
func ListConfigProviders() (interface{}, error) {
	_, body, err := Api("config/providers", nil)
	if err != nil {
		return nil, err
	} else {
		return body, nil
	}
}

// DumpConfig 导出当前RClone配置信息
// 获取所有已配置的远程存储设置
// 返回值:
//   - map[string]interface{}: 配置信息的映射表
//   - error: 操作过程中的错误，成功则为nil
func DumpConfig() (map[string]interface{}, error) {
	// 创建参数映射
	params := data.Map{}
	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}
	// 调用路径配置API（这里忽略返回值，可能是为了确保配置已加载）
	Api("config/paths", params)
	// 调用导出配置API
	_, body, err := Api("config/dump", params)
	if err != nil {
		return nil, err
	} else {
		return body.(map[string]interface{}), nil
	}
}

// ConfigCreate 创建新的RClone配置项
// 参数:
//   - name: 配置名称，用于标识此配置
//   - rtype: 远程存储类型，如"s3"、"ftp"等
//   - config: 配置的JSON字符串，包含连接参数
//
// 返回值: 操作过程中的错误，成功则为nil
func ConfigCreate(name string, rtype string, config string) error {
	// 将JSON字符串解析为映射表
	jsonData := map[string]interface{}{}
	err := json.Unmarshal([]byte(config), &jsonData)
	if err != nil {
		return err
	}
	// 创建API请求参数
	params := data.Map{"name": name, "parameters": jsonData, "type": rtype, "opt": data.Map{"obscure": true}}
	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}

	// 调用创建配置API
	code, data, err := Api("config/create", params)
	if err != nil {
		return err
	}
	// 检查API返回状态码
	if code != 200 {
		return errors.New("创建RClone配置失败: code:" + strconv.Itoa(code) + ", msg:" + getResultError(data))
	}
	return nil
}

// ConfigUpdate 更新现有的RClone配置项
// 参数:
//   - name: 配置名称，用于标识要更新的配置
//   - rtype: 远程存储类型，如"s3"、"ftp"等
//   - config: 更新后的配置JSON字符串
//
// 返回值: 操作过程中的错误，成功则为nil
func ConfigUpdate(name string, rtype string, config string) error {
	// 将JSON字符串解析为映射表
	jsonData := map[string]interface{}{}
	err := json.Unmarshal([]byte(config), &jsonData)
	if err != nil {
		return err
	}
	// 创建API请求参数
	params := data.Map{"name": name, "parameters": jsonData, "opt": data.Map{"obscure": true}}
	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}
	// 调用更新配置API
	code, data, err := Api("config/update", params)
	if err != nil {
		return err
	}
	// 检查API返回状态码
	if code != 200 {
		if code == 500 {
			msg := getResultError(data)
			// 如果配置不存在，则尝试创建它
			if msg == "couldn't find type field in config" {
				return ConfigCreate(name, rtype, config)
			}
		} else {
			return errors.New("更新RClone配置失败: code:" + strconv.Itoa(code) + ", msg:" + getResultError(data))
		}

	}
	return nil
}
