package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type memCache struct {
	// 最大内存
	maxMemorySize int64
	// 最大内存字符串
	maxMemorySizeStr string
	// 当前已使用内存
	currMemorySize int64
	// 缓存键值对
	values map[string]*memCacheValue
	// 读写锁 才能支持并发操作
	locker sync.RWMutex
	// 定期清除过期缓存
	clearExpireItemTimeInterval time.Duration
}

type memCacheValue struct {
	// value值
	val interface{}
	// 过期时间
	expireTime time.Time
	// 有效时长
	expire time.Duration
	// value值的大小
	size int64
}

func NewMemCache() Cache {
	mc := &memCache{
		values:                      make(map[string]*memCacheValue),
		clearExpireItemTimeInterval: time.Second,
	}
	go mc.clearExpireItem()
	return mc
}

// sizi : 1KB 100KB 1MB 2MB 1GB
func (mc *memCache) SetMaxMemory(size string) bool {
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	return false
}

// 将value 写入缓存
func (mc *memCache) Set(key string, val interface{}, expire time.Duration) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	v := &memCacheValue{
		val:        val,
		expireTime: time.Now().Add(expire),
		size:       GetValueSize(val),
		expire:     expire,
	}
	mc.del(key)
	mc.add(key, v)
	if mc.currMemorySize > mc.maxMemorySize {
		mc.del(key)
		log.Println(fmt.Sprintf("max memory size %d", mc.maxMemorySize))
		return false
	}
	return true
}

func (mc *memCache) get(key string) (*memCacheValue, bool) {
	val, ok := mc.values[key]
	return val, ok
}

func (mc *memCache) del(key string) {
	tmp, ok := mc.get(key)
	if ok && tmp != nil {
		mc.currMemorySize -= tmp.size
		delete(mc.values, key)
	}
}

func (mc *memCache) add(key string, val *memCacheValue) {
	mc.values[key] = val
	mc.currMemorySize += val.size
}

// 根据key值获取value
func (mc *memCache) Get(key string) (interface{}, bool) {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	v, ok := mc.get(key)
	if ok {
		// 判断缓存是否过期
		if v.expire != 0 && v.expireTime.Before(time.Now()) {
			mc.del(key)
			return nil, false
		}
		return v.val, ok
	}
	return nil, false
}

// 删除key值
func (mc *memCache) Del(key string) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.del(key)
	return true
}

// 判断Key是否存在
func (mc *memCache) Exists(key string) bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	_, ok := mc.values[key]
	return ok
}

// 清空所有key
func (mc *memCache) Flush() bool {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.values = make(map[string]*memCacheValue, 0) // 复制为空map，go的垃圾回收会自动回收
	mc.currMemorySize = 0
	return true
}

// 获取缓存中所有key的数量
func (mc *memCache) Keys() int64 {
	mc.locker.Lock()
	defer mc.locker.Unlock()
	return int64(len(mc.values))
}

// 定期去清空已经过期的key
func (mc *memCache) clearExpireItem() {
	timeTicket := time.NewTicker(mc.clearExpireItemTimeInterval)
	defer timeTicket.Stop()
	for {
		select {
		case <-timeTicket.C:
			for key, item := range mc.values {
				if item.expire != 0 && time.Now().After(item.expireTime) {
					mc.locker.Lock()
					mc.del(key)
					mc.locker.Unlock()
				}
			}
		}
	}
}

func GetValSize(val interface{}) int64 {
	bytes, _ := json.Marshal(val)
	size := int64((len(bytes)))
	return size
}
