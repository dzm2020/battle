// Command configvalidate 读取目录下 buffs.json / skills.json / units.json 并做交叉引用校验。
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
	config.Load(dir)
	fmt.Println("ok")
}
