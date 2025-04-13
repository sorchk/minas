package rclone

import (
	"bytes"
	"encoding/json"
	"server/utils/data"
	"time"
)

// JobStatus RClone异步作业状态结构体
// 用于表示RClone后台任务的执行状态和结果
type JobStatus struct {
	Id        uint        // 作业唯一标识符
	Duration  float64     // 作业执行时长（秒）
	StartTime time.Time   // 作业开始时间
	EndTime   time.Time   // 作业结束时间
	Finished  bool        // 作业是否已完成
	Success   bool        // 作业是否成功完成
	Error     string      // 如果作业失败，存储错误信息
	Output    interface{} // 作业输出结果
}

// Status 获取RClone异步作业的状态
// 通过作业ID查询当前作业的执行状态
// 参数 jobid: 要查询的作业ID
// 返回值:
//   - int: HTTP状态码
//   - JobStatus: 作业状态信息
//   - error: 操作过程中的错误，成功则为nil
func Status(jobid string) (int, JobStatus, error) {
	// 创建查询参数
	params := data.Map{"jobid": jobid}
	// 如果配置了特定的配置文件路径，则将其添加到参数中
	if GetRConfigPath() != "" {
		params["_config"] = data.Map{"config": GetRConfigPath()}
	}

	var result JobStatus
	// 将参数序列化为JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return 0, result, err
	}

	// 创建HTTP请求体
	requestBody := bytes.NewReader(jsonData)
	// 发送POST请求到job/status接口
	code, body, err := Post("job/status", requestBody)
	if err != nil {
		return code, result, err
	}

	// 解析响应数据到JobStatus结构体
	err = json.Unmarshal(body, &result)
	return code, result, err
}
