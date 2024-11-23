package main

import (
	"fmt"

	"github.com/stormi-li/omicafe-v1"
)

func main() {
	cafe := omicafe.NewFileCache("cache", 100*1024*1024)
	key := "name"
	value := []byte("stormi-li")
	cafe.Set(key, value)
	data, has := cafe.Get(key)
	fmt.Println(has, string(data))
	cafe.Del(key)
}
