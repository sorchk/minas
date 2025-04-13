package data

import (
	"strings"
)

// Map 是string到interface{}的映射类型，提供了额外的操作方法
type Map map[string]interface{}

// Get 获取指定键的值
// 参数 key: 要获取的键名
// 返回值: 键对应的值，如果键不存在则返回nil
func (m Map) Get(key string) interface{} {
	return m[key]
}

// TryGet 尝试获取指定键的值
// 参数 key: 要获取的键名
// 返回值: 键对应的值和一个表示键是否存在的布尔值
func (m Map) TryGet(key string) (v interface{}, ok bool) {
	v, ok = m[key]
	return
}

// Set 设置指定键的值
// 参数 key: 要设置的键名
// 参数 value: 要设置的值
// 返回值: 返回Map本身，支持链式调用
func (m Map) Set(key string, value interface{}) Map {
	m[key] = value
	return m
}

// Remove 移除指定键值对
// 参数 key: 要移除的键名
func (m Map) Remove(key string) {
	delete(m, key)
}

// Keys 获取Map中所有键的列表
// 返回值: 包含所有键名的字符串切片
func (m Map) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Contains 检查Map中是否包含指定的键
// 参数 key: 要检查的键名
// 返回值: 如果键存在则返回true，否则返回false
func (m Map) Contains(key string) bool {
	_, ok := m[key]
	return ok
}

// Merge 将src中的所有键值对合并到m中
// 如果键已存在且值为Map类型，则递归合并
// 参数 src: 要合并的源Map
func (m Map) Merge(src Map) {
	for k, sv := range src {
		if tv, ok := m[k]; ok {
			m1, ok1 := m.tryConvert(tv)
			m2, ok2 := m.tryConvert(sv)
			if ok1 && ok2 {
				m1.Merge(m2)
			}
			continue
		}

		if tm, ok := m.tryConvert(sv); ok {
			m[k] = tm
		} else {
			m[k] = sv
		}
	}
}

// Cover 将src中的所有键值对合并到m中，并替换已有值
// 与Merge不同，Cover会替换非Map类型的值
// 参数 src: 要合并的源Map
func (m Map) Cover(src Map) {
	for k, sv := range src {
		if tv, ok := m[k]; ok {
			m1, ok1 := m.tryConvert(tv)
			m2, ok2 := m.tryConvert(sv)
			if ok1 && ok2 {
				m1.Cover(m2)
			} else if ok2 {
				m[k] = m2
			} else {
				m[k] = sv
			}
			continue
		}

		if tm, ok := m.tryConvert(sv); ok {
			m[k] = tm
		} else {
			m[k] = sv
		}
	}
}

// Find 根据点号分隔的键路径查找嵌套值
// 例如："user.address.city" 将查找 m["user"]["address"]["city"]
// 参数 key: 点号分隔的键路径
// 返回值: 找到的值，如果路径不存在则返回nil
func (m Map) Find(key string) interface{} {
	if v, ok := m[key]; ok {
		return v
	}

	for i := strings.LastIndex(key, "."); i != -1; i = strings.LastIndex(key[:i], ".") {
		if v, ok := m[key[:i]]; ok {
			if tmp, ok := m.tryConvert(v); ok {
				return tmp.Find(key[i+1:])
			}
		}
	}
	return nil
}

// tryConvert 尝试将接口转换为Map类型
// 参数 i: 要转换的接口值
// 返回值: 转换后的Map和一个表示是否成功的布尔值
func (m Map) tryConvert(i interface{}) (Map, bool) {
	switch v := i.(type) {
	case Map:
		return v, true
	case map[string]interface{}:
		return Map(v), true
		//case []interface{}:
		//	for i, value := range v {
		//		if tm, ok := m.tryConvert(value); ok {
		//			v[i] = tm
		//		}
		//	}
	}
	return nil, false
}
