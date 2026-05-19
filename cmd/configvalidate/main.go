// Command configvalidate 读取目录下战斗 JSON 配表并做加载校验。
package main

import (
	"fmt"
	"os"

	"battle/internal/battle/config"
)

func main() {
	dir := "."
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	if err := config.Load(dir); err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")
}
