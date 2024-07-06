package main

import (
	"cache_system/cache"
	"fmt"
	"unsafe"
)

func main() {
	/*
		cache := cache_server.NewMemCache()
		cache.SetMaxMemory("100MB")
		cache.Set("int", 1, 10*time.Second)
		fmt.Println(cache.Get("int"))
		cache.Set("bool", false, 10*time.Second)
		fmt.Println(cache.Get("bool"))
		cache.Set("data", map[string]interface{}{"a": 1}, 10*time.Second)
		fmt.Println(cache.Get("data"))
		cache.Del("int")
		fmt.Println(cache.Get("int"))

		fmt.Println(cache.Keys())
		cache.Flush()
		fmt.Println(cache.Keys())
	*/

	// 都是转换成interface 然后在计算size的 所以值一样，计算的是interface 结构的大小
	cache.GetValSize(1)              // 16
	cache.GetValSize(false)          // 16
	cache.GetValSize("张三张三张三张三张三张三") // 16
	cache.GetValSize(map[string]string{
		"name": "张三张三张三张三张三张三",
		"addr": "北京市",
	}) // 16

	// 字符串没办法计算准确的长度
	fmt.Println(unsafe.Sizeof(1))
	fmt.Println(unsafe.Sizeof(false))
	fmt.Println(unsafe.Sizeof("张三张三张三张三张三张三"))
	fmt.Println(unsafe.Sizeof(map[string]string{
		"name": "张三张三张三张三张三张三",
		"addr": "北京市",
	}))
}
