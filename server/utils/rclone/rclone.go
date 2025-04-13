package rclone

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/utils"
	"server/utils/config"
	"server/utils/data"
	"server/utils/logger"
	"strconv"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
)

// const SocketPath = "/tmp/rclone.sock"
// ContentType HTTP请求Content-Type头的值，用于与RClone API通信
const ContentType = "application/json"

// client HTTP客户端实例，用于向RClone发送请求
var (
	client http.Client
)

// rcIp RClone API服务监听的IP地址
const rcIp = "localhost"

// GetRConfigPath 获取RClone配置文件路径
// 根据全局配置决定使用哪个配置文件路径
// 返回值: RClone配置文件的绝对路径
func GetRConfigPath() string {
	if config.CONF.RClone.ConfigPath == "" {
		file := filepath.Join(utils.GetWorkDirectory(), "/data/rclone.conf")
		if !fileutil.IsExist(file) {
			fileutil.CreateFile(file)
		}
		return file
	} else if config.CONF.RClone.ConfigPath == "default" {
		return ""
	}
	return config.CONF.RClone.ConfigPath
}

// GetRClonePath 获取RClone可执行文件路径
// 返回值: RClone可执行文件的路径，如果配置为空则返回默认值"rclone"
func GetRClonePath() string {
	if config.CONF.RClone.BinPath == "" {
		return "rclone"
	}
	return config.CONF.RClone.BinPath
}

// getRcAddr 获取RClone API服务地址
// 返回值: 格式为"IP:端口"的RClone API服务地址
func getRcAddr() string {
	if config.CONF.RClone.Port == 0 {
		return rcIp + ":5572"
	}
	return rcIp + ":" + convertor.ToString(config.CONF.RClone.Port)
}

// StartRClone 启动RClone守护进程
// 以API服务模式启动RClone
// 返回值:
//   - int: 启动的RClone进程ID
//   - error: 启动过程中的错误，成功则为nil
func StartRClone() (int, error) {
	cmdStr := []string{GetRClonePath()}
	if GetRConfigPath() != "" {
		cmdStr = append(cmdStr, "--config", GetRConfigPath())
	}
	cmdStr = append(cmdStr, "rcd", "--rc-addr="+getRcAddr(), "--rc-no-auth")
	tcmd := exec.Command(cmdStr[0], cmdStr[1:]...)

	tcmd.Stdout = os.Stdout
	tcmd.Stderr = os.Stderr
	logger.LOG.Info(tcmd.String())
	err := tcmd.Start()
	if err != nil {
		return 0, err
	} else {
		time.Sleep(1 * time.Second)
		pid := tcmd.Process.Pid
		return pid, nil
	}

}

// CoreObscure 使用RClone API加密字符串
// 参数 plainText: 需要加密的明文
// 返回值:
//   - string: 加密后的字符串
//   - error: 加密过程中的错误，成功则为nil
func CoreObscure(plainText string) (string, error) {
	params := data.Map{"clear": plainText}
	code, data, err := Api("core/obscure", params)
	if err != nil {
		return plainText, err
	}
	if code != 200 {
		return plainText, errors.New("获取RClone加密字符串失败: code:" + strconv.Itoa(code) + ", msg:" + getResultError(data))
	}
	result := getResultByKey(data, "obscured").(string)
	return result, nil
}

// getResultError 从API响应中提取错误消息
// 参数 data: API响应数据
// 返回值: 错误消息字符串
func getResultError(data interface{}) string {
	msg := getResultByKey(data, "error").(string)
	return msg
}

// getResultByKey 从API响应中获取指定键的值
// 参数 data: API响应数据
// 参数 key: 要获取的键名
// 返回值: 对应键的值
func getResultByKey(data interface{}, key string) interface{} {
	msg := data.(map[string]interface{})[key]
	return msg
}

// ConfigDelete 删除RClone配置
// 参数 name: 要删除的配置名称
// 返回值: 删除过程中的错误，成功则为nil
func ConfigDelete(name string) error {
	params := data.Map{"name": name}
	code, data, err := Api("config/delete", params)
	if err != nil {
		return err
	}
	if code != 200 {
		// b, err := json.Marshal(data)
		// if err != nil {
		// 	return err
		// }
		return errors.New("删除RClone配置失败: code:" + strconv.Itoa(code) + ", msg:" + getResultError(data))
	}
	return nil
}

// ApiProxy 代理请求到RClone API
// 参数 url: API端点URL
// 参数 requestBody: 请求体
// 返回值:
//   - int: HTTP状态码
//   - interface{}: 解析后的响应内容
//   - error: 请求过程中的错误，成功则为nil
func ApiProxy(url string, requestBody io.Reader) (int, interface{}, error) {
	code, body, err := Post(url, requestBody)
	if err != nil {
		return code, nil, err
	}
	var result interface{}
	err = json.Unmarshal(body, &result)
	return code, result, err

}

// Api 发送请求到RClone API并解析响应
// 参数 url: API端点URL
// 参数 params: 请求参数
// 返回值:
//   - int: HTTP状态码
//   - interface{}: 解析后的响应内容
//   - error: 请求过程中的错误，成功则为nil
func Api(url string, params data.Map) (int, interface{}, error) {
	if params == nil {
		params = data.Map{}
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return 0, nil, err
	}

	requestBody := bytes.NewReader(jsonData)
	code, body, err := Post(url, requestBody)
	if err != nil {
		return code, nil, err
	}
	logger.LOG.Infof("POST:%s, BODY:%s\n", url, string(body))
	var result interface{}
	err = json.Unmarshal(body, &result)
	return code, result, err

}

// Post 发送HTTP POST请求到RClone API
// 参数 url: API端点URL
// 参数 requestBody: 请求体
// 返回值:
//   - int: HTTP状态码
//   - []byte: 响应体内容
//   - error: 请求过程中的错误，成功则为nil
func Post(url string, requestBody io.Reader) (int, []byte, error) {
	if url[0] == '/' {
		url = url[1:]
	}
	resp, err := client.Post("http://"+getRcAddr()+"/"+url, ContentType, requestBody)
	if err != nil {
		return 0, nil, err
	}
	code := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return code, nil, err
	}
	defer resp.Body.Close()
	return code, body, err
}
