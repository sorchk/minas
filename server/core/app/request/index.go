package request

import (
	"encoding/json"
	"fmt"
	"server/utils"
	"server/utils/xxtea"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserID 从Gin上下文中获取当前用户ID
// 参数 ctx: Gin上下文
// 返回值: 当前用户ID，如果未找到则返回0
func GetUserID(ctx *gin.Context) uint {
	UserID, _ := ctx.Get("UserID")
	UserIDUint, ok := UserID.(uint)
	if !ok {
		UserIDUint = 0
	}
	return UserIDUint
}

// Pagination 分页查询基础结构体
// 包含常用的分页和排序参数
type Pagination struct {
	PageNo   int    // 页码，从1开始
	PageSize int    // 每页数量
	KeyWords string // 关键字，用于搜索
	DescAsc  string // 排序方向，"asc"升序或"desc"降序
	OrderBy  string // 排序字段
	Type     int    // 查询类型
}

// PageQuery 高级分页查询结构体
// 支持复杂查询条件和多字段排序
type PageQuery struct {
	Columns []string      // 查询的列
	Table   string        // 查询的表
	Filters []QueryFilter // 查询条件
	Sorts   []string      // 排序字段
	Page    int           // 页码
	Size    int           // 每页大小
}

// QueryFilter 查询条件结构体
// 支持嵌套条件和多种操作符
type QueryFilter struct {
	Column   string        // 列名
	Operator string        // 操作符，例如"=?", "like"等
	Value    any           // 值
	DataType string        // 数据类型
	Or       bool          // 是否为OR条件
	Filters  []QueryFilter // 嵌套条件
}

// getEncJson 从Gin上下文中获取加密的JSON数据并解密
// 参数 ctx: Gin上下文
// 参数 key: 查询参数的键
// 参数 v: 解密后的数据将存储到该变量中
// 返回值: 错误信息，如果没有错误则返回nil
func getEncJson(ctx *gin.Context, key string, v any) error {
	data := ctx.Query(key)
	if utils.IsEmpty(data) {
		return nil
	}
	data = xxtea.DecryptAuto(data, key)
	if utils.IsEmpty(data) {
		return nil
	}
	return json.Unmarshal([]byte(data), v)
}

// NewEqualFilter 创建一个等于条件的查询过滤器
// 参数 column: 列名
// 参数 value: 值
// 返回值: 查询过滤器
func NewEqualFilter(column string, value any) QueryFilter {
	return QueryFilter{Column: column, Operator: "=?", Value: value}
}

// NewLikeFilter 创建一个模糊匹配条件的查询过滤器
// 参数 column: 列名
// 参数 value: 值
// 返回值: 查询过滤器
func NewLikeFilter(column string, value any) QueryFilter {
	return QueryFilter{Column: column, Operator: " like ?", Value: fmt.Sprintf("%%%s%%", value)}
}

// AddFilter 向PageQuery中添加查询过滤器
// 参数 filter: 查询过滤器
func (query *PageQuery) AddFilter(filter ...QueryFilter) {
	query.Filters = append(query.Filters, filter...)
}

// GetPageQuery 从Gin上下文中获取分页查询参数
// 参数 ctx: Gin上下文
// 返回值: 分页查询对象
func GetPageQuery(ctx *gin.Context) PageQuery {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	var filtersArray []QueryFilter
	getEncJson(ctx, "filters", &filtersArray)
	var sortsArray = []string{}
	getEncJson(ctx, "sorts", &sortsArray)
	var columnsArray = []string{}
	getEncJson(ctx, "columns", &columnsArray)

	pageQuery := PageQuery{
		Page:    page,
		Size:    size,
		Filters: filtersArray,
		Sorts:   sortsArray,
		Columns: columnsArray,
	}
	return pageQuery
}
