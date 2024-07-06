package cache

import (
	"testing"
	"time"
)

func TestCacheOP(t *testing.T) {
	testData := []struct {
		key    string
		val    interface{}
		expire time.Duration
	}{
		{"key1", 666, time.Second * 10},
		{"qqweq", false, time.Second * 11},
		{"ddda", true, time.Second * 12},
		{"sfgas", map[string]interface{}{"a": 3, "b": false}, time.Second * 13},
		{"sdaf", "kangkang", time.Second * 14},
		{"wwws", "中文", time.Second * 15},
	}
	c := NewMemCache()
	c.SetMaxMemory("10MB")
	for _, item := range testData {
		c.Set(item.key, item.val, item.expire)
		val, ok := c.Get(item.key)
		if !ok {
			t.Error("缓存取值失败")
		}
		if item.key != "sfgas" && val != item.val {
			t.Error("缓存取值数据与预期不一致")
		}
		_, ok1 := val.(map[string]interface{})
		if item.key == "sfgas" && !ok1 {
			t.Error("缓存取值数据与预期不一致")
		}
	}
	if int64(len(testData)) != c.Keys() {
		t.Error("缓存数量不一致")
	}
	c.Del(testData[0].key)
	c.Del(testData[1].key)
	if int64(len(testData)) != c.Keys()+2 {
		t.Error("缓存数量不一致")
	}
	time.Sleep(time.Second * 16)
	if c.Keys() != 0 {
		t.Error("缓存清空失败")
	}
}
