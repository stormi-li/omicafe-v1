package main

import (
	"fmt"

	"github.com/stormi-li/omicafe-v1"
)

func main() {
	cafe := omicafe.NewFileCache("cache", 100*1024)
	cafe.Set("//afa/afa",[]byte("fsfs"))
	fmt.Println(string(cafe.Get("//afa/afa")))
	fmt.Println(cafe.CurrentSize())
	fmt.Println(cafe.CurrentSize())
}
