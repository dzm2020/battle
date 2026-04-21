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
	b, err := config.LoadCombatBundleFromDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	errs := config.ValidateCombatBundle(b)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}
	fmt.Println("ok")
}
