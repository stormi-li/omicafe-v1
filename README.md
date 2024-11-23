# Omicafe 文件缓存框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  
## 简介

**Omicafe** 是一个简单易用的本地文件缓存框架，支持以键值对的形式缓存文件。适合需要在本地文件系统中缓存小到中型数据的应用，提供高效的缓存管理。


## 功能

- **键值对缓存**：支持将数据缓存为键值对，便于管理。
- **文件存储**：数据直接存储在本地文件系统中，支持配置缓存目录和大小限制。
## 教程
### 安装
```shell
go get github.com/stormi-li/omicafe-v1
```
### 使用
```go
package main

import (
	"fmt"

	"github.com/stormi-li/omicafe-v1"
)

func main() {
	// 创建文件缓存实例，指定缓存目录和大小上限
	cafe := omicafe.NewFileCache("cache", 100*1024*1024) // 缓存目录: cache, 大小限制: 100 MB

	// 设置缓存数据
	key := "name"
	value := []byte("stormi-li")
	cafe.Set(key, value) // 存储键值对到缓存

	// 获取缓存数据
	data, has := cafe.Get(key) // 从缓存中读取数据
	if has {
		fmt.Println("缓存命中:", string(data))
	} else {
		fmt.Println("缓存未命中")
	}

	// 删除缓存数据
	cafe.Del(key) // 删除指定键的数据
	fmt.Println("缓存数据已删除")
}
```