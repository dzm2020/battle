package log

import "fmt"

func Debug(format string, a ...any) {
	format = fmt.Sprintf("%s\n", format)
	fmt.Printf(format, a...)
}

func Info(format string, a ...any) {
	format = fmt.Sprintf("%s\n", format)
	fmt.Printf(format, a...)
}
