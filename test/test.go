package main

import (
	"fmt"
	"strconv"

	"github.com/stormi-li/omicafe-v1"
)

func main() {
	cafe := omicafe.NewFileCache("cache", 5*1024)
	cafe.Set("//afa/afa", []byte("fsfs"))
	for i := 0; i < 1000; i++ {
		cafe.Set("//afa/afa"+strconv.Itoa(i), []byte("fsffsdfsdfs"))
	}
	data, has := cafe.Get("//afa/afa")
	fmt.Println(has, string(data))
	cafe.Get("//afa/afa-1")
	fmt.Println(cafe.GetCacheHitCount())
	fmt.Println(cafe.GetCacheClearCount())
	fmt.Println(cafe.GetCacheMissCount())
}
