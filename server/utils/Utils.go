package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/convertor"
)

// GetInt64 从map中获取int64类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 参数 def: 当键不存在或类型转换失败时的默认值
// 返回值: 获取到的int64值或默认值
func GetInt64(data map[string]any, key string, def int64) int64 {
	v, err := convertor.ToInt(data[key])
	if err != nil {
		return def
	} else {
		return v
	}
}

// GetInt 从map中获取int类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 参数 def: 当键不存在或类型转换失败时的默认值
// 返回值: 获取到的int值或默认值
func GetInt(data map[string]any, key string, def int) int {
	return int(GetInt64(data, key, int64(def)))
}

// GetUint 从map中获取uint类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 参数 def: 当键不存在或类型转换失败时的默认值
// 返回值: 获取到的uint值或默认值
func GetUint(data map[string]any, key string, def uint) uint {
	return uint(GetInt64(data, key, int64(def)))
}

// GetString 从map中获取字符串类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 返回值: 获取到的字符串值
func GetString(data map[string]any, key string) string {
	return convertor.ToString(data[key])
}

// GetBool 从map中获取布尔类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 参数 def: 当键不存在或类型转换失败时的默认值
// 返回值: 获取到的布尔值或默认值
func GetBool(data map[string]any, key string, def bool) bool {
	v, err := convertor.ToBool(convertor.ToString(data[key]))
	if err != nil {
		return def
	} else {
		return v
	}
}

// GetFloat64 从map中获取float64类型的值
// 参数 data: 数据源map
// 参数 key: 要获取的键名
// 参数 def: 当键不存在或类型转换失败时的默认值
// 返回值: 获取到的float64值或默认值
func GetFloat64(data map[string]any, key string, def float64) float64 {
	v, err := convertor.ToFloat(data[key])
	if err != nil {
		return def
	} else {
		return v
	}
}

// ToUint 将字符串转换为uint类型
// 参数 x: 要转换的字符串
// 返回值: 转换后的uint值，转换失败时返回0
func ToUint(x string) uint {
	n, _ := strconv.Atoi(x)
	return uint(n)
}

// ToUintDef 将字符串转换为uint类型，带有默认值
// 参数 x: 要转换的字符串
// 参数 def: 转换失败时的默认值
// 返回值: 转换后的uint值或默认值
func ToUintDef(x string, def uint) uint {
	n, err := strconv.Atoi(x)
	if err != nil {
		return def
	} else {
		return uint(n)
	}
}

// Md5 计算字符串的MD5哈希值
// 参数 src: 源字符串
// 返回值: MD5哈希值的十六进制表示
func Md5(src string) string {
	//hash := Config().Md5.Hash
	m := md5.New()
	m.Write([]byte(src))
	//res := hex.EncodeToString(m.Sum([]byte(hash)))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// IsEmpty 判断一个值是否为空
// 支持判断nil、空字符串、空数组、空map和空切片
// 参数 obj: 要判断的值
// 返回值: 如果为空则返回true，否则返回false
func IsEmpty(obj interface{}) bool {
	if obj == nil {
		// log.Println("obj is nil")
		return true
	}
	valueOfObj := reflect.ValueOf(obj)
	if valueOfObj.Kind() == reflect.Ptr {
		valueOfObj = valueOfObj.Elem()
	}
	if valueOfObj.Kind() == reflect.String {
		// log.Printf("obj is string: %v\n", obj)
		return strings.TrimSpace(fmt.Sprintf("%v", obj)) == ""
	} else if valueOfObj.Kind() == reflect.Array {
		// log.Printf("obj is array: %v\n", obj)
		return valueOfObj.Len() == 0
	} else if valueOfObj.Kind() == reflect.Map {
		// log.Printf("obj is map: %v\n", obj)
		return valueOfObj.Len() == 0
	} else if valueOfObj.Kind() == reflect.Slice {
		// log.Printf("obj is slice: %v\n", obj)
		// log.Printf("obj is len: %v\n", valueOfObj.Len())
		return valueOfObj.Len() == 0
	} else {
		// log.Printf("obj is unknown: %v\n", obj)
		return false
	}
}

// GetWorkDirectory 获取当前工作目录
// 返回值: 当前工作目录的路径
func GetWorkDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

// GetExeFileDirectory 获取可执行文件所在目录
// 返回值: 可执行文件所在目录的路径
func GetExeFileDirectory() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}
