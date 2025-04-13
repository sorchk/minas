package local

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// Item 缓存项，包含实际存储的对象和过期时间
type Item struct {
	Object     interface{} // 存储的实际对象，可以是任意类型
	Expiration int64       // 过期时间戳（纳秒级）
}

// Expired 检查缓存项是否已过期
// 返回true表示已过期，false表示未过期或永不过期
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false // 过期时间为0表示永不过期
	}
	return time.Now().UnixNano() > item.Expiration // 当前时间超过过期时间则已过期
}

const (
	// NoExpiration 表示缓存项永不过期
	NoExpiration time.Duration = -1
	// DefaultExpiration 表示使用创建缓存时指定的默认过期时间
	DefaultExpiration time.Duration = 0
)

// Cache 缓存对象，是对内部cache结构的封装
type Cache struct {
	*cache // 嵌入内部cache结构，继承其所有方法
}

// cache 内部缓存实现结构，包含了实际的缓存逻辑
type cache struct {
	defaultExpiration time.Duration             // 默认过期时间
	items             map[string]Item           // 缓存项存储映射
	mu                sync.RWMutex              // 读写锁，保证并发安全
	onEvicted         func(string, interface{}) // 当缓存项被移除时的回调函数
	janitor           *janitor                  // 清理过期缓存的后台goroutine
}

// Set 向缓存中添加项，如果已存在则替换
// k: 缓存键名
// x: 缓存的值，可以是任意类型
// d: 过期时间，如果是DefaultExpiration则使用缓存默认过期时间，如果是NoExpiration则永不过期
func (c *cache) Set(k string, x interface{}, d time.Duration) {
	c.mu.Lock()
	c.set(k, x, d)
	// 注意: 这里未使用defer解锁是为了性能考虑
	// 在Go 1.x版本中，defer会增加约200ns的性能开销
	c.mu.Unlock()
}

// set 内部设置缓存项的方法，需要在外部加锁保护
func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration // 使用默认过期时间
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano() // 计算过期时间戳
	}
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

// SetDefault 向缓存中添加项并使用默认过期时间
func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

// Add 仅当缓存中不存在该键或已过期时，添加缓存项
// 如果键已存在且未过期，则返回错误
func (c *cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

// Replace 仅当缓存中已存在该键且未过期时，替换缓存项
// 如果键不存在或已过期，则返回错误
func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

// GetItem 获取缓存项（包含值和过期信息）
// 如果缓存项不存在或已过期，返回false
func (c *cache) GetItem(k string) (Item, bool) {
	c.mu.RLock()
	// "内联"get和Expired函数的逻辑
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return item, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return item, false
		}
	}
	c.mu.RUnlock()
	return item, true
}

// Get 从缓存获取项
// 返回缓存的值和是否找到标志
// 如果缓存项不存在或已过期，返回nil和false
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	// "内联"get和Expired函数的逻辑提高性能
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

// GetWithExpiration 从缓存获取项及其过期时间
// 返回缓存的值、过期时间和是否找到标志
// 如果项永不过期，过期时间为零值time.Time{}
func (c *cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.RLock()
	// "内联"get和Expired函数的逻辑
	item, found := c.items[k]
	if !found {
		c.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, time.Time{}, false
		}

		// 返回项和过期时间
		c.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiration), true
	}

	// 如果过期时间<=0（即未设置过期时间）则返回项和零值时间
	c.mu.RUnlock()
	return item.Object, time.Time{}, true
}

// get 内部获取缓存项的方法，不加锁，调用者需确保已加锁
func (c *cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	// "内联"Expired函数的逻辑
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}

// IncrementInt 将int类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int，返回错误
func (c *cache) IncrementInt(k string, n int, d time.Duration) (int, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d) // 不存在或已过期时创建新项
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementInt8 将int8类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int8，返回错误
func (c *cache) IncrementInt8(k string, n int8, d time.Duration) (int8, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int8)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementInt16 将int16类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int16，返回错误
func (c *cache) IncrementInt16(k string, n int16, d time.Duration) (int16, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int16)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementInt32 将int32类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int32，返回错误
func (c *cache) IncrementInt32(k string, n int32, d time.Duration) (int32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementInt64 将int64类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int64，返回错误
func (c *cache) IncrementInt64(k string, n int64, d time.Duration) (int64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUint 将uint类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint，返回错误
func (c *cache) IncrementUint(k string, n uint, d time.Duration) (uint, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUintptr 将uintptr类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uintptr，返回错误
func (c *cache) IncrementUintptr(k string, n uintptr, d time.Duration) (uintptr, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uintptr)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUint8 将uint8类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint8，返回错误
func (c *cache) IncrementUint8(k string, n uint8, d time.Duration) (uint8, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint8)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUint16 将uint16类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint16，返回错误
func (c *cache) IncrementUint16(k string, n uint16, d time.Duration) (uint16, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint16)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUint32 将uint32类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint32，返回错误
func (c *cache) IncrementUint32(k string, n uint32, d time.Duration) (uint32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementUint64 将uint64类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint64，返回错误
func (c *cache) IncrementUint64(k string, n uint64, d time.Duration) (uint64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementFloat32 将float32类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是float32，返回错误
func (c *cache) IncrementFloat32(k string, n float32, d time.Duration) (float32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(float32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// IncrementFloat64 将float64类型的缓存项值增加n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是float64，返回错误
func (c *cache) IncrementFloat64(k string, n float64, d time.Duration) (float64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(float64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float64", k)
	}
	nv := rv + n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementInt 将int类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int，返回错误
func (c *cache) DecrementInt(k string, n int, d time.Duration) (int, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementInt8 将int8类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int8，返回错误
func (c *cache) DecrementInt8(k string, n int8, d time.Duration) (int8, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int8)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int8", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementInt16 将int16类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int16，返回错误
func (c *cache) DecrementInt16(k string, n int16, d time.Duration) (int16, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int16)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int16", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementInt32 将int32类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int32，返回错误
func (c *cache) DecrementInt32(k string, n int32, d time.Duration) (int32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int32", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementInt64 将int64类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是int64，返回错误
func (c *cache) DecrementInt64(k string, n int64, d time.Duration) (int64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(int64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an int64", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUint 将uint类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint，返回错误
func (c *cache) DecrementUint(k string, n uint, d time.Duration) (uint, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUintptr 将uintptr类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uintptr，返回错误
func (c *cache) DecrementUintptr(k string, n uintptr, d time.Duration) (uintptr, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uintptr)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uintptr", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUint8 将uint8类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint8，返回错误
func (c *cache) DecrementUint8(k string, n uint8, d time.Duration) (uint8, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint8)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint8", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUint16 将uint16类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint16，返回错误
func (c *cache) DecrementUint16(k string, n uint16, d time.Duration) (uint16, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint16)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint16", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUint32 将uint32类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint32，返回错误
func (c *cache) DecrementUint32(k string, n uint32, d time.Duration) (uint32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint32", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementUint64 将uint64类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是uint64，返回错误
func (c *cache) DecrementUint64(k string, n uint64, d time.Duration) (uint64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(uint64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an uint64", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementFloat32 将float32类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是float32，返回错误
func (c *cache) DecrementFloat32(k string, n float32, d time.Duration) (float32, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(float32)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float32", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// DecrementFloat64 将float64类型的缓存项值减少n
// 如果项不存在或已过期，则用n值创建新项
// 如果项存在但类型不是float64，返回错误
func (c *cache) DecrementFloat64(k string, n float64, d time.Duration) (float64, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.set(k, n, d)
		c.mu.Unlock()
		return n, nil
	}
	rv, ok := v.Object.(float64)
	if !ok {
		c.mu.Unlock()
		return 0, fmt.Errorf("The value for %s is not an float64", k)
	}
	nv := rv - n
	c.set(k, nv, d)
	c.mu.Unlock()
	return nv, nil
}

// Delete 从缓存删除项
// 如果设置了onEvicted回调函数，则在删除后调用此函数
func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}
}

// delete 内部删除缓存项的方法，返回被删除的值和是否触发了回调
func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

// 键值对结构，用于保存被删除的项信息
type keyAndValue struct {
	key   string
	value interface{}
}

// DeleteExpired 删除所有已过期的缓存项
// 遍历所有缓存项，检查是否过期，过期则删除
func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		// "内联"expired函数的逻辑
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	c.mu.Unlock()
	// 在释放锁后执行回调，避免在锁内执行用户代码
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

// OnEvicted 设置一个可选的函数，当项从缓存中被删除时调用
// 传入被删除项的键名和值作为参数
func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

// Save 使用Gob序列化将缓存项保存到io.Writer
//
// 注意：此方法已弃用，建议使用c.Items()和NewFrom()
func (c *cache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&c.items)
	return
}

// SaveFile 将缓存项保存到指定文件
// 如果文件不存在则创建，存在则覆盖
//
// 注意：此方法已弃用，建议使用c.Items()和NewFrom()
func (c *cache) SaveFile(fname string) error {
	fp, err := os.Create(fname)
	if err != nil {
		return err
	}
	err = c.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

// Load 从io.Reader加载（Gob序列化的）缓存项
// 如果当前缓存中已存在相同键且未过期的项，则不加载该项
//
// 注意：此方法已弃用，建议使用c.Items()和NewFrom()
func (c *cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	err := dec.Decode(&items)
	if err == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		for k, v := range items {
			ov, found := c.items[k]
			if !found || ov.Expired() {
				c.items[k] = v
			}
		}
	}
	return err
}

// LoadFile 从指定文件加载缓存项
// 如果当前缓存中已存在相同键的项，则不加载该项
//
// 注意：此方法已弃用，建议使用c.Items()和NewFrom()
func (c *cache) LoadFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	err = c.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

// Items 将所有未过期的缓存项复制到新的map并返回
// 在返回前先过滤掉所有已过期的项
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		// "内联"Expired函数的逻辑
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		m[k] = v
	}
	return m
}

// ItemCount 返回缓存中的项数量
// 此数量可能包括已过期但尚未清理的项
func (c *cache) ItemCount() int {
	c.mu.RLock()
	n := len(c.items)
	c.mu.RUnlock()
	return n
}

// Flush 清空缓存中的所有项
func (c *cache) Flush() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()
}

// janitor 清理器结构，负责定期删除过期的缓存项
type janitor struct {
	Interval time.Duration // 清理间隔
	stop     chan bool     // 停止信号通道
}

// Run 启动清理器的工作循环
// 按照指定的间隔定期调用DeleteExpired删除过期项
func (j *janitor) Run(c *cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

// stopJanitor 停止清理器的工作
func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

// runJanitor 启动一个清理器来定期清理过期的缓存项
func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

// newCache 创建一个新的缓存实例
func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1 // 默认设为永不过期
	}
	c := &cache{
		defaultExpiration: de,
		items:             m,
	}
	return c
}

// newCacheWithJanitor 创建一个带清理器的缓存实例
func newCacheWithJanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)
	// 这个技巧确保了janitor协程（如果启用的话）不会阻止返回的Cache对象被垃圾回收
	// 当Cache被垃圾回收时，终结器会停止janitor协程，然后c可以被回收
	C := &Cache{c}
	if ci > 0 {
		runJanitor(c, ci)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

// New 创建一个新的缓存实例，指定默认过期时间和清理间隔
// 如果过期时间小于1（或NoExpiration），则缓存项默认永不过期，必须手动删除
// 如果清理间隔小于1，则不会自动删除过期项，需要手动调用DeleteExpired()
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}

// NewFrom 创建一个新的缓存实例，指定默认过期时间和清理间隔，同时使用提供的items作为初始缓存项
// 如果过期时间小于1（或NoExpiration），则缓存项默认永不过期
// 如果清理间隔小于1，则不会自动删除过期项
//
// 这个方法适用于从已序列化的缓存恢复（如使用gob.Encode()序列化的c.Items()），
// 或者当预期缓存会达到一定规模时，预先分配足够大小的map以提高启动性能
//
// 注意：只有cache的方法会同步访问这个map，所以创建缓存后不建议保留对这个map的任何引用
// 如有必要，可以通过c.Items()在之后访问该map（但同样存在上述警告）
//
// 关于序列化的注意事项: 当使用gob等序列化工具时，确保在编码通过c.Items()获取的map前，
// 使用gob.Register()注册缓存中存储的各种类型；同样，在解码包含items map的blob前
// 也要注册相同的类型
func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}
