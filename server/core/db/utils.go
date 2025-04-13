package db

import (
	"fmt"
	"reflect"
	"server/utils/global"

	"gorm.io/gorm/schema"
)

// 解析字段名称
func ParseColumnName(field reflect.StructField) string {
	tagSetting := schema.ParseTagSetting(field.Tag.Get("gorm"), ";")
	// 从 GORM 标签设置中获取 "COLUMN" 的值，这是数据库中对应的列名
	name, ok := tagSetting["COLUMN"]
	if ok {
		return name
	}
	return global.DB.Config.NamingStrategy.ColumnName("", field.Name)
}

func getUserCacheKey(userID uint) string {
	return fmt.Sprintf("db:user:%d", userID)
}
func SetUserNameCache(userID uint, name string) {
	key := getUserCacheKey(userID)
	global.CACHE.Set(key, name)
}
func GetUserNameCache(userID uint) string {
	key := getUserCacheKey(userID)
	name, err := global.CACHE.Get(key)
	if err == nil && name != "" {
		return name
	} else {
		err = global.DB.Table("users").Where("id = ?", userID).Select("name").Scan(&name).Error
		if err == nil {
			SetUserNameCache(userID, name)
			return name
		}
		return name
	}
}
