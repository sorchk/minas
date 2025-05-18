package utils

import (
	"fmt"
	"server/dagflow/core/el"
	"server/dagflow/model"
	"strconv"
)

//	func GetVal(node model.TaskNode, key string) (any, error) {
//		if val, ok := node.Properties[key]; ok {
//			return val, nil
//		}
//		return nil, fmt.Errorf("key not found")
//	}
func GetStr(node model.TaskNode, key string, defval string, data map[string]any) string {
	if val, ok := node.Properties[key]; ok {
		result, err := GetEl(val, data)
		if err != nil || result == nil {
			return defval
		}
		if strVal, ok := result.(string); ok {
			return strVal
		}
		return fmt.Sprintf("%v", result)
	}
	return defval
}
func GetBool(node model.TaskNode, key string, defval bool, data map[string]any) bool {
	if val, ok := node.Properties[key].(bool); ok {
		return val
	} else if val, ok := node.Properties[key]; ok {
		result, err := GetEl(val, data)
		if err != nil || result == nil {
			return defval
		}
		v, err := strconv.ParseBool(result.(string))
		if err == nil {
			return v
		} else {
			return defval
		}
	}
	return defval
}
func GetInt(node model.TaskNode, key string, defval int, data map[string]any) int {
	if val, ok := node.Properties[key].(int); ok {
		return val
	} else if val, ok := node.Properties[key]; ok {
		result, err := GetEl(val, data)
		if err != nil || result == nil {
			return defval
		}
		v, err := strconv.Atoi(result.(string))
		if err == nil {
			return v
		}
	}
	return defval
}
func GetList(node model.TaskNode, varkey string, data map[string]any) ([]any, error) {
	if vars, ok := node.Properties[varkey].([]any); ok {
		newList := make([]any, 0)
		for _, v := range vars {
			result, err := GetEl(v, data)
			if err != nil {
				return nil, err
			}
			newList = append(newList, result)
		}
		return newList, nil
	}
	return nil, fmt.Errorf("配置属性:%s 不存在或不是List类型", varkey)
}
func GetMap(node model.TaskNode, varkey string, data map[string]any) (map[string]any, error) {
	newMap := make(map[string]any)
	if vars, ok := node.Properties[varkey].([]interface{}); ok {
		for _, v := range vars {
			if varMap, ok := v.(map[string]interface{}); ok {
				key := varMap["key"].(string)
				val, err := GetEl(varMap["value"], data)
				if err == nil {
					newMap[key] = val
				} else {
					return nil, err
				}
			}
		}
		return newMap, nil
	} else {
		return nil, fmt.Errorf("配置属性:%s 不存在或不是Key=Value类型", varkey)
	}
}
func GetEl(expression any, data map[string]any) (any, error) {
	if data == nil {
		return expression, nil
	}
	if len(data) == 0 {
		return expression, nil
	}
	if expression != nil && fmt.Sprintf("%T", expression) == "string" {
		exprStr, ok := expression.(string)
		if !ok {
			return expression, nil
		}
		return el.Evaluate(exprStr, data)
	}
	return expression, nil
}
